package batcher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/common/geth"
	"github.com/0glabs/0g-da-client/core"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/contract"
	"github.com/0glabs/0g-da-client/disperser/signer"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/go-merkletree"
)

const (
	QuantizationFactor = uint(1)
	indexerWarmupDelay = 2 * time.Second
)

type TimeoutConfig struct {
	EncodingTimeout   time.Duration
	ChainReadTimeout  time.Duration
	ChainWriteTimeout time.Duration
	SigningTimeout    time.Duration
}

type Config struct {
	PullInterval             time.Duration
	FinalizerInterval        time.Duration
	EncoderSocket            string
	SRSOrder                 int
	NumConnections           int
	EncodingRequestQueueSize int
	// BatchSizeMBLimit is the maximum size of a batch in MB
	BatchSizeMBLimit     uint
	MaxNumRetriesPerBlob uint
	ConfirmerNum         uint

	DAEntranceContractAddress     string
	DASignersContractAddress      string
	EncodingInterval              time.Duration
	SigningInterval               time.Duration
	MaxNumRetriesForSign          uint
	FinalizedBlockCount           uint
	ExpirationPollIntervalSec     uint64
	SignedPullInterval            time.Duration
	VerifiedCommitRootsTxGasLimit uint64
}

type Batcher struct {
	Config
	TimeoutConfig

	Queue         disperser.BlobStore
	Dispatcher    disperser.Dispatcher
	EncoderClient disperser.EncoderClient

	EncodingStreamer *EncodingStreamer
	Metrics          *Metrics

	finalizer   Finalizer
	confirmer   *Confirmer
	sliceSigner *SliceSigner
	logger      common.Logger
}

func NewBatcher(
	config Config,
	timeoutConfig TimeoutConfig,
	ethConfig geth.EthClientConfig,
	queue disperser.BlobStore,
	dispatcher disperser.Dispatcher,
	encoderClient disperser.EncoderClient,
	finalizer Finalizer,
	confirmer *Confirmer,
	daContract *contract.DAContract,
	logger common.Logger,
	metrics *Metrics,
	blobKeyCache *disperser.BlobKeyCache,
) (*Batcher, error) {
	batchTrigger := NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		uint64(config.BatchSizeMBLimit)*1024*1024, // convert to bytes
	)
	streamerConfig := StreamerConfig{
		SRSOrder:               config.SRSOrder,
		EncodingRequestTimeout: timeoutConfig.EncodingTimeout,
		EncodingQueueLimit:     config.EncodingRequestQueueSize,
		EncodingInterval:       config.EncodingInterval,
	}
	encodingWorkerPool := workerpool.New(config.NumConnections)
	encodingStreamer, err := NewEncodingStreamer(streamerConfig, queue, encoderClient, batchTrigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
	if err != nil {
		return nil, err
	}

	signerClient, err := signer.NewSignerClient(timeoutConfig.SigningTimeout)
	if err != nil {
		return nil, err
	}

	signerTrigger := NewSignatureSizeNotifier(
		make(chan struct{}, 1),
		100,
	)
	signerConfig := SignerConfig{
		SigningRequestTimeout: timeoutConfig.SigningTimeout,
		MaxNumRetriesPerBlob:  config.MaxNumRetriesPerBlob,
		MaxNumRetriesSign:     config.MaxNumRetriesForSign,
		SigningInterval:       config.SigningInterval,
	}
	signingWorkerPool := workerpool.New(config.NumConnections)
	sliceSigner, err := NewEncodedSliceSigner(
		ethConfig,
		signerConfig,
		signingWorkerPool,
		signerTrigger,
		signerClient,
		queue,
		daContract,
		metrics,
		logger,
		blobKeyCache,
	)
	if err != nil {
		return nil, err
	}

	return &Batcher{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		Queue:         queue,
		Dispatcher:    dispatcher,
		EncoderClient: encoderClient,

		EncodingStreamer: encodingStreamer,
		Metrics:          metrics,

		finalizer:   finalizer,
		confirmer:   confirmer,
		sliceSigner: sliceSigner,
		logger:      logger,
	}, nil
}

func (b *Batcher) Start(ctx context.Context) error {
	// Wait for few seconds for indexer to index blockchain
	// This won't be needed when we switch to using Graph node
	time.Sleep(indexerWarmupDelay)
	err := b.EncodingStreamer.Start(ctx)
	if err != nil {
		return err
	}
	batchTrigger := b.EncodingStreamer.EncodedSizeNotifier
	submitAggregateSignaturesTrigger := b.sliceSigner.SignatureSizeNotifier

	b.sliceSigner.EncodingStreamer = b.EncodingStreamer
	b.sliceSigner.Finalizer = b.finalizer
	b.sliceSigner.Start(ctx)

	// confirmer
	b.confirmer.EncodingStreamer = b.EncodingStreamer
	b.confirmer.SliceSigner = b.sliceSigner
	b.confirmer.Start(ctx)
	// finalizer
	b.finalizer.Start(ctx)

	go func() {
		ticker := time.NewTicker(b.PullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if ts, blobMetadata, err := b.HandleSingleBatch(ctx); err != nil {
					for _, metadata := range blobMetadata {
						meta, err := b.Queue.GetBlobMetadata(ctx, metadata.GetBlobKey())
						if err != nil {
							b.logger.Error("[batcher] failed to get blob metadata", "key", metadata.GetBlobKey(), "err", err)
						} else {
							if meta.BlobStatus == disperser.Failed {
								b.logger.Info("[batcher] handle encode batch reach max retries", "key", metadata.GetBlobKey())
								b.EncodingStreamer.RemoveEncodedBlob(metadata)
								b.Queue.RemoveBlob(ctx, metadata)
							}
						}
					}

					b.EncodingStreamer.RemoveBatchingStatus(ts)
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Debug("[batcher] no encoded results to make a batch with")
					} else {
						b.logger.Error("[batcher] failed to process a batch", "err", err)
					}
				}
			case <-batchTrigger.Notify:
				ticker.Stop()
				if ts, blobMetadata, err := b.HandleSingleBatch(ctx); err != nil {
					for _, metadata := range blobMetadata {
						meta, err := b.Queue.GetBlobMetadata(ctx, metadata.GetBlobKey())
						if err != nil {
							b.logger.Error("[batcher] failed to get blob metadata", "key", metadata.GetBlobKey(), "err", err)
						} else {
							if meta.BlobStatus == disperser.Failed {
								b.logger.Info("[batcher] handle encode batch reach max retries", "key", metadata.GetBlobKey())
								b.EncodingStreamer.RemoveEncodedBlob(metadata)
								b.Queue.RemoveBlob(ctx, metadata)
							}
						}
					}

					b.EncodingStreamer.RemoveBatchingStatus(ts)
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Debug("[batcher] no encoded results to make a batch with(Notified)")
					} else {
						b.logger.Error("[batcher] failed to process a batch(Notified)", "err", err)
					}
				}
				ticker.Reset(b.PullInterval)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(b.SignedPullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				if err := b.HandleSignedBatch(ctx); err != nil {
					if errors.Is(err, errNoSignedResults) {
						b.logger.Debug("[batcher] no signed results to make a batch with")
					} else {
						b.logger.Error("[batcher] failed to process a signed batch", "err", err)
					}
				}

			case <-submitAggregateSignaturesTrigger.Notify:
				ticker.Stop()
				if err := b.HandleSignedBatch(ctx); err != nil {
					if errors.Is(err, errNoSignedResults) {
						b.logger.Debug("[batcher] no signed results to make a batch with(Notified)")
					} else {
						b.logger.Error("[batcher] failed to process a signed batch(Notified)", "err", err)
					}
				}

				ticker.Reset((b.PullInterval))
			}
		}
	}()

	return nil
}

func serializeProof(proof *merkletree.Proof) []byte {
	proofBytes := make([]byte, 0)
	for _, hash := range proof.Hashes {
		proofBytes = append(proofBytes, hash[:]...)
	}
	return proofBytes
}

func (b *Batcher) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	for _, metadata := range blobMetadatas {
		err := b.Queue.HandleBlobFailure(ctx, metadata, b.MaxNumRetriesPerBlob)
		if err != nil {
			b.logger.Error("[batcher] HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}
		b.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
	}
	b.Metrics.UpdateBatchError(reason, len(blobMetadatas))

	// Return the error(s)
	return result.ErrorOrNil()
}

func (b *Batcher) HandleSingleBatch(ctx context.Context) (uint64, []*disperser.BlobMetadata, error) {
	log := b.logger
	// start a timer
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(f float64) {
		b.Metrics.ObserveLatency("total", f*1000) // make milliseconds
	}))
	defer timer.ObserveDuration()

	stageTimer := time.Now()
	log.Info("[batcher] Creating batch", "ts", stageTimer)
	batch, ts, err := b.EncodingStreamer.CreateBatch()
	if err != nil {
		return ts, nil, err
	}
	log.Info("[batcher] CreateBatch took", "duration", time.Since(stageTimer), "blobNum", len(batch.EncodedBlobs))

	// Get the batch header hash
	log.Trace("[batcher] Getting batch header hash...")
	headerHash, err := batch.BatchHeader.GetBatchHeaderHash()
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchHeaderHash)
		return ts, batch.BlobMetadata, fmt.Errorf("HandleSingleBatch: error getting batch header hash: %w", err)
	}

	proofs := make([]*merkletree.Proof, 0)
	// Prepare data writes to kv stream
	for blobIndex := range batch.BlobMetadata {
		var blobHeader *core.BlobHeader
		// generate inclusion proof
		if blobIndex >= len(batch.BlobHeaders) {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchBlobIndex)
			return ts, batch.BlobMetadata, fmt.Errorf("HandleSingleBatch: error preparing kv data: blob header at index %d not found in batch", blobIndex)
		}
		blobHeader = batch.BlobHeaders[blobIndex]

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchBlobHeaderHash)
			return ts, batch.BlobMetadata, fmt.Errorf("HandleSingleBatch: failed to get blob header hash: %w", err)
		}
		merkleProof, err := batch.MerkleTree.GenerateProof(blobHeaderHash[:], 0)
		if err != nil {
			_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchProof)
			return ts, batch.BlobMetadata, fmt.Errorf("HandleSingleBatch: failed to generate blob header inclusion proof: %w", err)
		}
		proofs = append(proofs, merkleProof)
	}

	// Dispatch encoded batch
	log.Info("[batcher] Dispatching encoded batch...")
	stageTimer = time.Now()
	batch.TxHash, err = b.Dispatcher.DisperseBatch(ctx, headerHash, batch.BatchHeader, batch.EncodedBlobs, batch.BlobHeaders)
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchSubmitRoot)
		return ts, batch.BlobMetadata, err
	}
	log.Info("[batcher] DisperseBatch took", "duration", time.Since(stageTimer))

	b.sliceSigner.SignerChan <- &SignInfo{
		headerHash: headerHash,
		batch:      batch,
		proofs:     proofs,
		ts:         ts,
		reties:     0,
	}

	return ts, nil, nil
}

func (b *Batcher) HandleSignedBatch(ctx context.Context) error {
	log := b.logger

	s, signedTs, err := b.sliceSigner.GetCommitRootSubmissionBatch()
	if err != nil {
		b.sliceSigner.RemoveBatchingStatus(signedTs)
		return err
	}

	log.Info("[batcher] Create signed batch", "batch size", len(s), "signed ts", signedTs)

	submissions := make([]*core.CommitRootSubmission, 0)
	headerHash := make([][32]byte, 0)
	batch := make([]*batch, 0)
	ts := make([]uint64, 0)
	proofs := make([][]*merkletree.Proof, 0)
	epochs := make([]*big.Int, 0)
	quorumIds := make([]*big.Int, 0)
	for _, item := range s {
		submissions = append(submissions, item.submissions...)

		headerHash = append(headerHash, item.headerHash)
		batch = append(batch, item.batch)
		ts = append(ts, item.ts)
		proofs = append(proofs, item.proofs)

		epochs = append(epochs, item.epoch)
		quorumIds = append(quorumIds, item.quorumId)
	}

	stageTimer := time.Now()
	var txHash *eth_common.Hash
	if len(submissions) > 0 {
		hash, err := b.Dispatcher.SubmitAggregateSignatures(ctx, submissions)
		if err != nil {
			for idx, item := range batch {
				_ = b.handleFailure(ctx, item.BlobMetadata, FailSubmitAggregateSignatures)
				for _, metadata := range item.BlobMetadata {
					meta, err := b.Queue.GetBlobMetadata(ctx, metadata.GetBlobKey())
					if err != nil {
						log.Error("[batcher] failed to get blob metadata", "key", metadata.GetBlobKey(), "err", err)
					} else {
						if meta.BlobStatus == disperser.Failed {
							log.Info("[batcher] submit aggregateSignatures reach max retries", "key", metadata.GetBlobKey())
							b.EncodingStreamer.RemoveEncodedBlob(metadata)
							b.sliceSigner.RemoveSignedBlob(ts[idx])
							b.Queue.RemoveBlob(ctx, metadata)
						}
					}
				}

				b.EncodingStreamer.RemoveBatchingStatus(ts[idx])
			}
			b.sliceSigner.RemoveBatchingStatus(signedTs)
			if len(s) > 1 {
				b.sliceSigner.SignedBatchSize = uint(len(s)) / 2
			} else {
				b.sliceSigner.SignedBatchSize = 1
			}
			return err
		}

		txHash = &hash
	}

	b.logger.Info("[batcher] submit aggregate signatures", "duration", time.Since(stageTimer))

	b.sliceSigner.SignedBatchSize = 0
	b.confirmer.ConfirmChan <- &BatchInfo{
		headerHash: headerHash,
		batch:      batch,
		ts:         ts,
		proofs:     proofs,
		signedTs:   signedTs,
		txHash:     txHash,
		epochs:     epochs,
		quorumIds:  quorumIds,
	}

	return nil
}
