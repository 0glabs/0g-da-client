package batcher

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const maxRetries = 3
const baseDelay = 1 * time.Second

// Finalizer runs periodically to finalize blobs that have been confirmed
type Finalizer interface {
	Start(ctx context.Context)
	FinalizeBlobs(ctx context.Context) error
	LatestFinalizedBlock() uint64
}

type finalizer struct {
	mu sync.RWMutex

	timeout                    time.Duration
	loopInterval               time.Duration
	blobStore                  disperser.BlobStore
	ethClient                  common.EthClient
	rpcClient                  common.RPCEthClient
	maxNumRetriesPerBlob       uint
	logger                     common.Logger
	latestFinalizedBlock       uint64
	defaultFinalizedBlockCount uint64
	kvStore                    *disperser.Store
	ExpirationPollIntervalSec  uint64
}

func NewFinalizer(timeout time.Duration, batcherConfig Config, blobStore disperser.BlobStore, ethClient common.EthClient, rpcClient common.RPCEthClient, logger common.Logger, kvStore *disperser.Store) Finalizer {
	return &finalizer{
		timeout:                    timeout,
		loopInterval:               batcherConfig.FinalizerInterval,
		blobStore:                  blobStore,
		ethClient:                  ethClient,
		rpcClient:                  rpcClient,
		maxNumRetriesPerBlob:       batcherConfig.MaxNumRetriesPerBlob,
		logger:                     logger,
		latestFinalizedBlock:       0,
		defaultFinalizedBlockCount: uint64(batcherConfig.FinalizedBlockCount),
		kvStore:                    kvStore,
		ExpirationPollIntervalSec:  batcherConfig.ExpirationPollIntervalSec,
	}
}

func (f *finalizer) Start(ctx context.Context) {
	go f.expireLoop()

	go func() {
		for {
			f.updateFinalizedBlockNumber(ctx)
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		ticker := time.NewTicker(f.loopInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := f.FinalizeBlobs(ctx); err != nil {
					f.logger.Error("[finalizer] failed to finalize blobs", "err", err)
				}
			}
		}
	}()
}

func (f *finalizer) updateFinalizedBlockNumber(ctx context.Context) {
	var header = types.Header{}
	var err error
	for i := 0; i < maxRetries; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, f.timeout)
		defer cancel()
		err = f.rpcClient.CallContext(ctxWithTimeout, &header, "eth_getBlockByNumber", "finalized", false)
		if err == nil {
			break
		}

		retrySec := math.Pow(2, float64(i))
		f.logger.Error("[finalizer] Finalizer: error getting latest finalized block", "err", err, "retrySec", retrySec)
		time.Sleep(time.Duration(retrySec) * baseDelay)
	}

	var blockNumber uint64
	if err != nil {
		f.logger.Error("[finalizer] error getting latest finalized block", "err", err)

		ctxWithTimeout, cancel := context.WithTimeout(ctx, f.timeout)
		defer cancel()

		err := f.rpcClient.CallContext(ctxWithTimeout, &header, "eth_getBlockByNumber", "latest", false)
		if err != nil {
			f.logger.Error("[finalizer] error getting latest block", "err", err)
		} else {
			blockNumber = header.Number.Uint64() - f.defaultFinalizedBlockCount
		}
	} else {
		blockNumber = uint64(header.Number.Uint64())
	}

	f.mu.Lock()
	if blockNumber > f.latestFinalizedBlock {
		f.latestFinalizedBlock = blockNumber
		f.logger.Info("[finalizer] latest finalized block number updated", "number", f.latestFinalizedBlock)
	}
	f.mu.Unlock()
}

func (f *finalizer) LatestFinalizedBlock() uint64 {
	f.mu.RLock()
	blockNumber := f.latestFinalizedBlock
	f.mu.RUnlock()

	return blockNumber
}

// FinalizeBlobs checks the latest finalized block and marks blobs in `confirmed` state as `finalized` if their confirmation
// block number is less than or equal to the latest finalized block number.
// If it failes to process some blobs, it will log the error, skip the failed blobs, and will not return an error. The function should be invoked again to retry.
func (f *finalizer) FinalizeBlobs(ctx context.Context) error {
	f.mu.RLock()
	finalizedBlokNumber := f.latestFinalizedBlock
	f.mu.RUnlock()

	metadatas, err := f.blobStore.GetBlobMetadataByStatus(ctx, disperser.Confirmed)
	if err != nil {
		return fmt.Errorf("FinalizeBlobs: error getting blob headers: %w", err)
	}

	f.logger.Info("[finalizer] FinalizeBlobs: finalizing blobs", "numBlobs", len(metadatas), "finalizedBlockNumber", finalizedBlokNumber)

	finalizedMetadatas := make([]*disperser.BlobMetadata, 0)
	for _, m := range metadatas {
		blobKey := m.GetBlobKey()
		confirmationMetadata, err := f.blobStore.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			f.logger.Error("[finalizer] FinalizeBlobs: error getting confirmed metadata", "blobKey", blobKey.String(), "err", err)
			continue
		}

		// Leave as confirmed if the confirmation block is after the latest finalized block (not yet finalized)
		if uint64(confirmationMetadata.ConfirmationInfo.ConfirmationBlockNumber) > finalizedBlokNumber {
			continue
		}

		// confirmation block number may have changed due to reorg
		confirmationBlockNumber, err := f.getTransactionBlockNumber(ctx, confirmationMetadata.ConfirmationInfo.ConfirmationTxnHash)
		if errors.Is(err, ethereum.NotFound) {
			// The confirmed block is finalized, but the transaction is not found. It means the transaction should be considered forked/invalid and the blob should be considered as failed.
			err := f.blobStore.HandleBlobFailure(ctx, m, f.maxNumRetriesPerBlob)
			if err != nil {
				f.logger.Error("[finalizer] FinalizeBlobs: error marking blob as failed", "blobKey", blobKey.String(), "err", err)
			}
			continue
		}
		if err != nil {
			f.logger.Error("[finalizer] FinalizeBlobs: error getting transaction block number", "err", err)
			continue
		}

		// Leave as confirmed if the reorged confirmation block is after the latest finalized block (not yet finalized)
		if uint64(confirmationBlockNumber) > finalizedBlokNumber {
			continue
		}

		confirmationMetadata.ConfirmationInfo.ConfirmationBlockNumber = uint32(confirmationBlockNumber)
		err = f.blobStore.MarkBlobFinalized(ctx, blobKey)
		if err != nil {
			f.logger.Error("[finalizer] FinalizeBlobs: error marking blob as finalized", "blobKey", blobKey.String(), "err", err)
			continue
		}

		finalizedMetadatas = append(finalizedMetadatas, m)
	}

	f.PersistConfirmedBlobs(ctx, finalizedMetadatas)
	f.logger.Info("[finalizer] FinalizeBlobs: successfully processed all finalized blobs")
	return nil
}

func (f *finalizer) PersistConfirmedBlobs(ctx context.Context, metadatas []*disperser.BlobMetadata) error {
	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	for _, metadata := range metadatas {
		retrieveMetadata := disperser.BlobRetrieveMetadata{
			DataRoot:    metadata.ConfirmationInfo.DataRoot,
			Epoch:       metadata.ConfirmationInfo.Epoch,
			QuorumId:    metadata.ConfirmationInfo.QuorumId,
			BlockNumber: metadata.ConfirmationInfo.ConfirmationBlockNumber,
		}
		val, err := retrieveMetadata.Serialize()
		if err != nil {
			return errors.WithMessage(err, "failed to serialize retrieve metadata")
		}

		key := []byte(metadata.GetBlobKey().String())
		keys = append(keys, key)
		values = append(values, val)
	}

	_, err := f.kvStore.StoreMetadataBatch(ctx, keys, values)
	if err != nil {
		return errors.WithMessage(err, "failed to save retrieve metadata to kv db")
	}

	f.logger.Info("[finalizer] removing confirmed blobs")
	for _, metadata := range metadatas {
		f.logger.Info("[finalizer] removing blob", "blob key", metadata.GetBlobKey().String())
		err := f.blobStore.RemoveBlob(ctx, metadata)
		if err != nil {
			f.logger.Warn("[finalizer] failed to remove blob", "error", err)
		}
	}
	f.logger.Info("[finalizer] confirmed blobs removed")
	return nil
}

func (f *finalizer) getTransactionBlockNumber(ctx context.Context, hash gcommon.Hash) (uint64, error) {
	var ctxWithTimeout context.Context
	var cancel context.CancelFunc
	var txReceipt *types.Receipt
	var err error
	for i := 0; i < maxRetries; i++ {
		ctxWithTimeout, cancel = context.WithTimeout(ctx, f.timeout)
		defer cancel()
		txReceipt, err = f.ethClient.TransactionReceipt(ctxWithTimeout, hash)
		if err == nil {
			break
		}
		if errors.Is(err, ethereum.NotFound) {
			// If the transaction is not found, it means the transaction has been reorged out of the chain.
			return 0, err
		}

		retrySec := math.Pow(2, float64(i))
		f.logger.Error("[finalizer] Finalizer: error getting transaction", "err", err, "retrySec", retrySec, "hash", hash.Hex())
		time.Sleep(time.Duration(retrySec) * baseDelay)
	}

	if err != nil {
		return 0, fmt.Errorf("Finalizer: error getting transaction receipt after retries: %w", err)
	}

	return txReceipt.BlockNumber.Uint64(), nil
}

// The expireLoop is a loop that is run once per configured second(s) while the node
// is running. It scans for expired blobs and removes them from the local database.
func (f *finalizer) expireLoop() {
	f.logger.Info("[finalizer] start expireLoop goroutine in background to periodically remove expired blobs on the node")
	ticker := time.NewTicker(time.Duration(f.ExpirationPollIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// We cap the time the deletion function can run, to make sure there is no overlapping
		// between loops and the garbage collection doesn't take too much resource.
		// The heuristic is to cap the GC time to a percentage of the poll interval, but at
		// least have 1 second.
		timeLimitSec := uint64(math.Max(float64(f.ExpirationPollIntervalSec)*gcPercentageTime, 1.0))
		numBlobsDeleted, err := f.kvStore.DeleteExpiredEntries(time.Now().Unix(), timeLimitSec)
		f.logger.Info("[finalizer] complete an expiration cycle to remove expired blobs", "num expired blobs found and removed", numBlobsDeleted)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				f.logger.Error("Expiration cycle exited with ContextDeadlineExceed, meaning more expired blobs need to be removed, which will continue in next cycle", "time limit (sec)", timeLimitSec)
			} else {
				f.logger.Error("Expiration cycle encountered error when removing expired blobs, which will be retried in next cycle", "err", err)
			}
		}
	}
}
