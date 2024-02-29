package batcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/go-multierror"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/go-merkletree"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/core"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser"
)

const (
	QuantizationFactor = uint(1)
	indexerWarmupDelay = 2 * time.Second
)

type QuorumInfo struct {
	Assignments        map[core.OperatorID]core.Assignment
	Info               core.AssignmentInfo
	QuantizationFactor uint
}

type TimeoutConfig struct {
	EncodingTimeout   time.Duration
	ChainReadTimeout  time.Duration
	ChainWriteTimeout time.Duration
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
}

type Batcher struct {
	Config
	TimeoutConfig

	Queue         disperser.BlobStore
	Dispatcher    disperser.Dispatcher
	EncoderClient disperser.EncoderClient

	EncodingStreamer *EncodingStreamer
	Metrics          *Metrics

	finalizer Finalizer
	confirmer *Confirmer
	logger    common.Logger
}

func NewBatcher(
	config Config,
	timeoutConfig TimeoutConfig,
	queue disperser.BlobStore,
	dispatcher disperser.Dispatcher,
	encoderClient disperser.EncoderClient,
	finalizer Finalizer,
	confirmer *Confirmer,
	logger common.Logger,
	metrics *Metrics,
) (*Batcher, error) {
	batchTrigger := NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		uint64(config.BatchSizeMBLimit)*1024*1024, // convert to bytes
	)
	streamerConfig := StreamerConfig{
		SRSOrder:               config.SRSOrder,
		EncodingRequestTimeout: timeoutConfig.EncodingTimeout,
		EncodingQueueLimit:     config.EncodingRequestQueueSize,
	}
	encodingWorkerPool := workerpool.New(config.NumConnections)
	encodingStreamer, err := NewEncodingStreamer(streamerConfig, queue, encoderClient, batchTrigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, logger)
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

		finalizer: finalizer,
		confirmer: confirmer,
		logger:    logger,
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
	// confirmer
	b.confirmer.EncodingStreamer = b.EncodingStreamer
	b.confirmer.Start(ctx)
	// finalizer
	if !b.Queue.MetadataHashAsBlobKey() {
		b.finalizer.Start(ctx)
	}

	go func() {
		ticker := time.NewTicker(b.PullInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if ts, err := b.HandleSingleBatch(ctx); err != nil {
					b.EncodingStreamer.RemoveBatchingStatus(ts)
					if errors.Is(err, errNoEncodedResults) {
						b.logger.Debug("[batcher] no encoded results to make a batch with")
					} else {
						b.logger.Error("[batcher] failed to process a batch", "err", err)
					}
				}
			case <-batchTrigger.Notify:
				ticker.Stop()
				if ts, err := b.HandleSingleBatch(ctx); err != nil {
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

func (b *Batcher) HandleSingleBatch(ctx context.Context) (uint64, error) {
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
		return ts, err
	}
	log.Info("[batcher] CreateBatch took", "duration", time.Since(stageTimer), "blobNum", len(batch.EncodedBlobs))

	// Get the batch header hash
	log.Trace("[batcher] Getting batch header hash...")
	headerHash, err := batch.BatchHeader.GetBatchHeaderHash()
	if err != nil {
		_ = b.handleFailure(ctx, batch.BlobMetadata, FailBatchHeaderHash)
		return ts, fmt.Errorf("HandleSingleBatch: error getting batch header hash: %w", err)
	}

	proofs := make([]*merkletree.Proof, 0)
	// Prepare data writes to kv stream
	for blobIndex := range batch.BlobMetadata {
		var blobHeader *core.BlobHeader
		// generate inclusion proof
		if blobIndex >= len(batch.BlobHeaders) {
			return ts, fmt.Errorf("HandleSingleBatch: error preparing kv data: blob header at index %d not found in batch", blobIndex)
		}
		blobHeader = batch.BlobHeaders[blobIndex]

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			return ts, fmt.Errorf("HandleSingleBatch: failed to get blob header hash: %w", err)
		}
		merkleProof, err := batch.MerkleTree.GenerateProof(blobHeaderHash[:], 0)
		if err != nil {
			return ts, fmt.Errorf("HandleSingleBatch: failed to generate blob header inclusion proof: %w", err)
		}
		proofs = append(proofs, merkleProof)
	}

	// Dispatch encoded batch
	log.Info("[batcher] Dispatching encoded batch...")
	stageTimer = time.Now()
	batch.TxHash, err = b.Dispatcher.DisperseBatch(ctx, headerHash, batch.BatchHeader, batch.EncodedBlobs, proofs)
	if err != nil {
		return ts, err
	}
	log.Info("[batcher] DisperseBatch took", "duration", time.Since(stageTimer))

	b.confirmer.ConfirmChan <- &BatchInfo{
		headerHash: headerHash,
		batch:      batch,
		proofs:     proofs,
		ts:         ts,
	}
	return ts, nil
}
