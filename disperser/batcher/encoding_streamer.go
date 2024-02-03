package batcher

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/wealdtech/go-merkletree"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/core"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser"
)

const encodingInterval = 2 * time.Second

var errNoEncodedResults = errors.New("no encoded results")

type EncodedSizeNotifier struct {
	mu sync.Mutex

	Notify chan struct{}
	// threshold is the size of the total encoded blob results in bytes that triggers the notifier
	threshold uint64
	// active is set to false after the notifier is triggered to prevent it from triggering again for the same batch
	// This is reset when CreateBatch is called and the encoded results have been consumed
	active bool
}

type StreamerConfig struct {

	// SRSOrder is the order of the SRS used for encoding
	SRSOrder int
	// EncodingRequestTimeout is the timeout for each encoding request
	EncodingRequestTimeout time.Duration

	// EncodingQueueLimit is the maximum number of encoding requests that can be queued
	EncodingQueueLimit int
}

type EncodingStreamer struct {
	StreamerConfig

	mu sync.RWMutex

	EncodedBlobstore     *encodedBlobStore
	ReferenceBlockNumber uint
	Pool                 common.WorkerPool
	EncodedSizeNotifier  *EncodedSizeNotifier

	blobStore disperser.BlobStore
	// chainState            core.IndexedChainState
	encoderClient disperser.EncoderClient
	// assignmentCoordinator core.AssignmentCoordinator

	encodingCtxCancelFuncs []context.CancelFunc

	metrics *EncodingStreamerMetrics
	logger  common.Logger
}

type batch struct {
	EncodedBlobs []*core.EncodedBlob
	BlobMetadata []*disperser.BlobMetadata
	BlobHeaders  []*core.BlobHeader
	BatchHeader  *core.BatchHeader
	State        *core.IndexedOperatorState
	MerkleTree   *merkletree.MerkleTree
	TxHash       eth_common.Hash
}

func NewEncodedSizeNotifier(notify chan struct{}, threshold uint64) *EncodedSizeNotifier {
	return &EncodedSizeNotifier{
		Notify:    notify,
		threshold: threshold,
		active:    true,
	}
}

func NewEncodingStreamer(
	config StreamerConfig,
	blobStore disperser.BlobStore,
	encoderClient disperser.EncoderClient,
	encodedSizeNotifier *EncodedSizeNotifier,
	workerPool common.WorkerPool,
	metrics *EncodingStreamerMetrics,
	logger common.Logger) (*EncodingStreamer, error) {
	if config.EncodingQueueLimit <= 0 {
		return nil, fmt.Errorf("EncodingQueueLimit should be greater than 0")
	}
	return &EncodingStreamer{
		StreamerConfig:         config,
		EncodedBlobstore:       newEncodedBlobStore(logger),
		ReferenceBlockNumber:   uint(0),
		Pool:                   workerPool,
		EncodedSizeNotifier:    encodedSizeNotifier,
		blobStore:              blobStore,
		encoderClient:          encoderClient,
		encodingCtxCancelFuncs: make([]context.CancelFunc, 0),
		metrics:                metrics,
		logger:                 logger,
	}, nil
}

func (e *EncodingStreamer) Start(ctx context.Context) error {
	encoderChan := make(chan EncodingResultOrStatus)

	// goroutine for handling blob encoding responses
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case response := <-encoderChan:
				err := e.ProcessEncodedBlobs(ctx, response)
				if err != nil {
					if strings.Contains(err.Error(), context.Canceled.Error()) {
						// ignore canceled errors because canceled encoding requests are normal
						continue
					}
					e.logger.Error("error processing encoded blobs", "err", err)
				}
			}
		}
	}()

	// goroutine for making blob encoding requests
	go func() {
		ticker := time.NewTicker(encodingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := e.RequestEncoding(ctx, encoderChan)
				if err != nil {
					e.logger.Warn("error requesting encoding", "err", err)
				}
			}
		}
	}()

	return nil
}

func (e *EncodingStreamer) RequestEncoding(ctx context.Context, encoderChan chan EncodingResultOrStatus) error {
	stageTimer := time.Now()
	// pull new blobs and send to encoder
	e.logger.Info("[encodingstreamer] requesting processing blobs..")
	metadatas, err := e.blobStore.GetBlobMetadataByStatus(ctx, disperser.Processing)
	if err != nil {
		return fmt.Errorf("error getting blob metadatas: %w", err)
	}
	// filter requested/encoded blobs
	n := 0
	for _, metadata := range metadatas {
		if !e.EncodedBlobstore.HasEncodingRequested(metadata.GetBlobKey()) {
			metadatas[n] = metadata
			n++
		}
	}
	metadatas = metadatas[:n]
	if len(metadatas) == 0 {
		e.logger.Info("no new metadatas to encode")
		return nil
	}

	e.logger.Info("[encodingstreamer] metadata in processing status", "numMetadata", len(metadatas))

	waitingQueueSize := e.Pool.WaitingQueueSize()
	numMetadatastoProcess := e.EncodingQueueLimit - waitingQueueSize
	if numMetadatastoProcess > len(metadatas) {
		numMetadatastoProcess = len(metadatas)
	}
	if numMetadatastoProcess <= 0 {
		// encoding queue is full
		e.logger.Warn("[RequestEncoding] worker pool queue is full. skipping this round of encoding requests", "waitingQueueSize", waitingQueueSize, "encodingQueueLimit", e.EncodingQueueLimit)
		return nil
	}
	// only process subset of blobs so it doesn't exceed the EncodingQueueLimit
	// TODO: this should be done at the request time and keep the cursor so that we don't fetch the same metadata every time
	metadatas = metadatas[:numMetadatastoProcess]

	e.logger.Trace("[encodingstreamer] new metadatas to encode", "numMetadata", len(metadatas), "duration", time.Since(stageTimer))

	stageTimer = time.Now()
	blobs, err := e.blobStore.GetBlobsByMetadata(ctx, metadatas)
	if err != nil {
		return fmt.Errorf("error getting blobs from blob store: %w", err)
	}
	e.logger.Trace("[RequestEncoding] retrieved blobs to encode", "numBlobs", len(blobs), "duration", time.Since(stageTimer))

	e.logger.Trace("[RequestEncoding] encoding blobs...", "numBlobs", len(blobs))

	for i := range metadatas {
		metadata := metadatas[i]

		e.RequestEncodingForBlob(ctx, metadata, blobs[metadata.GetBlobKey()], encoderChan)
	}

	return nil
}

type pendingRequestInfo struct {
	// BlobQuorumInfo *core.BlobQuorumInfo
	EncodingParams core.EncodingParams
	// Assignments    map[core.OperatorID]core.Assignment
}

func (e *EncodingStreamer) RequestEncodingForBlob(ctx context.Context, metadata *disperser.BlobMetadata, blob *core.Blob, encoderChan chan EncodingResultOrStatus) {

	// Validate the encoding parameters for each quorum

	blobKey := metadata.GetBlobKey()

	blobLength := core.GetBlobLength(metadata.RequestMetadata.BlobSize)

	chunkLength, chunkNum := core.SplitToChunks(blobLength)

	params, err := core.GetEncodingParams(chunkLength, chunkNum)
	if err != nil {
		e.logger.Error("[RequestEncodingForBlob] error getting encoding params", "err", err)
		return
	}

	err = core.ValidateEncodingParams(params, int(blobLength), e.SRSOrder)
	if err != nil {
		e.logger.Error("[RequestEncodingForBlob] invalid encoding params", "err", err)
		// Cancel the blob
		err := e.blobStore.MarkBlobFailed(ctx, blobKey)
		if err != nil {
			e.logger.Error("[RequestEncodingForBlob] error marking blob failed", "err", err)
		}
		return
	}

	pending := pendingRequestInfo{
		EncodingParams: params,
	}

	encodingCtx, cancel := context.WithTimeout(ctx, e.EncodingRequestTimeout)
	e.Pool.Submit(func() {
		defer cancel()
		commits, chunks, err := e.encoderClient.EncodeBlob(encodingCtx, blob.Data, pending.EncodingParams)
		if err != nil {
			encoderChan <- EncodingResultOrStatus{Err: err, EncodingResult: EncodingResult{
				BlobMetadata: metadata,
			}}
			return
		}

		encoderChan <- EncodingResultOrStatus{
			EncodingResult: EncodingResult{
				BlobMetadata:         metadata,
				ReferenceBlockNumber: 0,
				Commitment:           commits,
				Chunks:               chunks,
			},
			Err: nil,
		}
	})
	e.EncodedBlobstore.PutEncodingRequest(blobKey)
	e.logger.Trace("requested encoding for blob", "blob key", blobKey)
}

func (e *EncodingStreamer) ProcessEncodedBlobs(ctx context.Context, result EncodingResultOrStatus) error {
	if result.Err != nil {
		e.EncodedBlobstore.DeleteEncodingRequest(result.BlobMetadata.GetBlobKey())
		return fmt.Errorf("error encoding blob: %w, blob hash: %v", result.Err, result.BlobMetadata.BlobHash)
	}

	err := e.EncodedBlobstore.PutEncodingResult(&result.EncodingResult)
	if err != nil {
		return fmt.Errorf("failed to putEncodedBlob: %w", err)
	}

	e.logger.Trace("blob encoded", "blob key", result.BlobMetadata.GetBlobKey())

	count, encodedSize := e.EncodedBlobstore.GetEncodedResultSize()
	e.metrics.UpdateEncodedBlobs(count, encodedSize)
	if e.EncodedSizeNotifier.threshold > 0 && encodedSize >= e.EncodedSizeNotifier.threshold {
		e.EncodedSizeNotifier.mu.Lock()

		if e.EncodedSizeNotifier.active {
			e.logger.Info("encoded size threshold reached", "size", encodedSize)
			e.EncodedSizeNotifier.Notify <- struct{}{}
			// make sure this doesn't keep triggering before encoded blob store is reset
			e.EncodedSizeNotifier.active = false
		}
		e.EncodedSizeNotifier.mu.Unlock()
	}

	return nil
}

// CreateBatch makes a batch from all blobs in the encoded blob store.
// If successful, it returns a batch, and updates the reference block number for next batch to use.
// Otherwise, it returns an error and keeps the blobs in the encoded blob store.
// This function is meant to be called periodically in a single goroutine as it resets the state of the encoded blob store.
func (e *EncodingStreamer) CreateBatch() (*batch, uint64, error) {
	// Get all encoded blobs
	ts := uint64(time.Now().Nanosecond())
	encodedResults := e.EncodedBlobstore.GetNewEncodingResults(ts)

	// Reset the notifier
	e.EncodedSizeNotifier.mu.Lock()
	e.EncodedSizeNotifier.active = true
	e.EncodedSizeNotifier.mu.Unlock()

	if len(encodedResults) == 0 {
		return nil, ts, errNoEncodedResults
	}

	encodedBlobByKey := make(map[disperser.BlobKey]*core.EncodedBlob)
	blobHeaderByKey := make(map[disperser.BlobKey]*core.BlobHeader)
	metadataByKey := make(map[disperser.BlobKey]*disperser.BlobMetadata)
	blobKeys := make([]disperser.BlobKey, 0)
	for i := range encodedResults {
		// each result represent an encoded result per blob
		result := encodedResults[i]

		blobKey := result.BlobMetadata.GetBlobKey()
		blobKeys = append(blobKeys, blobKey)
		if _, ok := encodedBlobByKey[blobKey]; !ok {
			metadataByKey[blobKey] = result.BlobMetadata
		}
		blobHeader := &core.BlobHeader{
			BlobCommitments: *result.Commitment,
		}
		blobHeaderByKey[blobKey] = blobHeader
		encodedBlobByKey[blobKey] = &core.BlobMessage{
			BlobHeader: blobHeader,
			Bundles:    result.Chunks,
		}
	}

	// sort blobs by size
	sort.SliceStable(blobKeys, func(i, j int) bool {
		return blobHeaderByKey[blobKeys[i]].Length < blobHeaderByKey[blobKeys[j]].Length
	})

	// Transform maps to slices so orders in different slices match
	encodedBlobs := make([]*core.EncodedBlob, len(metadataByKey))
	blobHeaders := make([]*core.BlobHeader, len(metadataByKey))
	metadatas := make([]*disperser.BlobMetadata, len(metadataByKey))
	i := 0
	for _, key := range blobKeys {
		encodedBlobs[i] = encodedBlobByKey[key]
		blobHeaders[i] = blobHeaderByKey[key]
		metadatas[i] = metadataByKey[key]
		i++
	}

	// Populate the batch header
	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: e.ReferenceBlockNumber,
		BatchRoot:            [32]byte{},
	}

	tree, err := batchHeader.SetBatchRoot(blobHeaders)
	if err != nil {
		return nil, ts, err
	}

	e.ReferenceBlockNumber = 0

	return &batch{
		EncodedBlobs: encodedBlobs,
		BatchHeader:  batchHeader,
		BlobHeaders:  blobHeaders,
		BlobMetadata: metadatas,
		MerkleTree:   tree,
	}, ts, nil
}

func (e *EncodingStreamer) RemoveEncodedBlob(metadata *disperser.BlobMetadata) {
	e.EncodedBlobstore.DeleteEncodingResult(metadata.GetBlobKey())
}

func (e *EncodingStreamer) RemoveBatchingStatus(ts uint64) {
	e.EncodedBlobstore.DeleteBatchingStatus(ts)
}
