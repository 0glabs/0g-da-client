package batcher

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/common/geth"
	"github.com/0glabs/0g-da-client/disperser"
	"github.com/0glabs/0g-da-client/disperser/contract"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
)

const (
	// The percentage of time in garbage collection in a GC cycle.
	gcPercentageTime = 0.1
)

type Confirmer struct {
	mu sync.RWMutex

	Queue            disperser.BlobStore
	EncodingStreamer *EncodingStreamer
	SliceSigner      *SliceSigner

	daContract  *contract.DAContract
	ConfirmChan chan *BatchInfo

	pendingBatches       []*BatchInfo
	MaxNumRetriesPerBlob uint

	routines uint

	retryOption contract.RetryOption

	logger  common.Logger
	Metrics *Metrics
}

type BatchInfo struct {
	headerHash [][32]byte
	batch      []*batch
	ts         []uint64
	proofs     [][]*merkletree.Proof
	signedTs   uint64
	txHash     eth_common.Hash
	epochs     []*big.Int
	quorumIds  []*big.Int
}

func NewConfirmer(ethConfig geth.EthClientConfig, batcherConfig Config, queue disperser.BlobStore, daContract *contract.DAContract, logger common.Logger, metrics *Metrics) (*Confirmer, error) {
	if ethConfig.TxGasLimit > 0 {
		blockchain.CustomGasLimit = uint64(ethConfig.TxGasLimit)
	}

	return &Confirmer{
		Queue:                queue,
		daContract:           daContract,
		ConfirmChan:          make(chan *BatchInfo),
		pendingBatches:       make([]*BatchInfo, 0),
		routines:             batcherConfig.ConfirmerNum,
		MaxNumRetriesPerBlob: batcherConfig.MaxNumRetriesPerBlob,
		retryOption: contract.RetryOption{
			Rounds:   ethConfig.ReceiptPollingRounds,
			Interval: ethConfig.ReceiptPollingInterval,
		},
		logger:  logger,
		Metrics: metrics,
	}, nil
}

func (c *Confirmer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batchInfo := <-c.ConfirmChan:
				c.putPendingBatches(batchInfo)
			}
		}
	}()

	for i := 0; i < int(c.routines); i++ {
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					batchInfo := c.getPendingBatch()
					if batchInfo != nil {
						if err := c.ConfirmBatch(ctx, batchInfo); err != nil {
							c.logger.Error("[confirmer] failed to confirm batch", "err", err)
						}
					}
				}
			}
		}()
	}
}

func (c *Confirmer) putPendingBatches(info *BatchInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.pendingBatches = append(c.pendingBatches, info)
	c.logger.Info(`[confirmer] received pending batch`, "queue size", len(c.pendingBatches))
}

func (c *Confirmer) getPendingBatch() *BatchInfo {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.pendingBatches) == 0 {
		return nil
	}
	info := c.pendingBatches[0]
	c.pendingBatches = c.pendingBatches[1:]
	c.logger.Info(`[confirmer] retrieved one pending batch`, "queue size", len(c.pendingBatches))
	return info
}

func (c *Confirmer) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	for _, metadata := range blobMetadatas {
		err := c.Queue.HandleBlobFailure(ctx, metadata, c.MaxNumRetriesPerBlob)
		if err != nil {
			c.logger.Error("[confirmer] HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}
		c.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
	}
	c.Metrics.UpdateBatchError(reason, len(blobMetadatas))

	// Return the error(s)
	return result.ErrorOrNil()
}

func (c *Confirmer) waitForReceipt(txHash eth_common.Hash) (uint32, error) {
	if txHash.Cmp(eth_common.Hash{}) == 0 {
		return 0, errors.New("empty transaction hash")
	}
	c.logger.Info("[confirmer] Waiting signing batch be confirmed", "transaction hash", txHash)
	// data is not duplicate, there is a new transaction
	receipt, err := c.daContract.WaitForReceipt(txHash, true, c.retryOption)
	if err != nil {
		return 0, err
	}

	blockNumber := receipt.BlockNumber
	c.logger.Debug("[confirmer] waiting signed tx to be confirmed", "receipt block", blockNumber)

	return uint32(blockNumber), nil
}

func (c *Confirmer) ConfirmBatch(ctx context.Context, batchInfo *BatchInfo) error {
	blockNumber, err := c.waitForReceipt(batchInfo.txHash)
	if err != nil {
		// batch is not confirmed
		for idx := range batchInfo.ts {
			_ = c.handleFailure(ctx, batchInfo.batch[idx].BlobMetadata, FailConfirmBatch)
			// c.EncodingStreamer.RemoveBatchingStatus(ts)
		}

		c.SliceSigner.RemoveBatchingStatus(batchInfo.signedTs)
		return err
	}

	for idx, batch := range batchInfo.batch {
		proofs := batchInfo.proofs[idx]

		epoch := batchInfo.epochs[idx].Uint64()
		quorumId := batchInfo.quorumIds[idx].Uint64()

		batchID := batchInfo.ts[idx]
		c.logger.Info("[confirmer] batch confirmed.", "batch ID", batchID, "transaction hash", batch.TxHash)
		// Mark the blobs as complete
		c.logger.Info("[confirmer] Marking blobs as complete...")
		stageTimer := time.Now()
		blobsToRetry := make([]*disperser.BlobMetadata, 0)
		var updateConfirmationInfoErr error
		for blobIndex, metadata := range batch.BlobMetadata {
			confirmationInfo := &disperser.ConfirmationInfo{
				BatchHeaderHash:         batchInfo.headerHash[idx],
				BlobIndex:               uint32(blobIndex),
				ReferenceBlockNumber:    0,
				BatchRoot:               batch.BatchHeader.BatchRoot[:],
				BlobInclusionProof:      serializeProof(proofs[blobIndex]),
				CommitmentRoot:          batch.BlobHeaders[blobIndex].CommitmentRoot,
				DataRoot:                batch.EncodedBlobs[blobIndex].StorageRoot,
				Epoch:                   epoch,
				QuorumId:                quorumId,
				Length:                  uint32(batch.BlobHeaders[blobIndex].Length),
				BatchID:                 uint32(batchID),
				SubmissionTxnHash:       batch.TxHash,
				ConfirmationTxnHash:     batchInfo.txHash,
				ConfirmationBlockNumber: blockNumber,
			}
			c.logger.Trace("[confirmer] confirming blob", "blob key", metadata.GetBlobKey())
			_, updateConfirmationInfoErr := c.Queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
			if updateConfirmationInfoErr == nil {
				c.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Confirmed)
				// remove encoded blob from storage so we don't disperse it again
				c.EncodingStreamer.RemoveEncodedBlob(metadata)
				c.logger.Trace("[confirmer] blob confirmed", "blob key", metadata.GetBlobKey())

			} else {
				c.logger.Error("[confirmer] HandleSingleBatch: error updating blob confirmed metadata", "err", updateConfirmationInfoErr)
				blobsToRetry = append(blobsToRetry, batch.BlobMetadata[blobIndex])
			}
			requestTime := time.Unix(0, int64(metadata.RequestMetadata.RequestedAt))
			c.Metrics.ObserveLatency("E2E", float64(time.Since(requestTime).Milliseconds()))
		}

		if len(blobsToRetry) > 0 {
			_ = c.handleFailure(ctx, blobsToRetry, FailUpdateConfirmationInfo)
			if len(blobsToRetry) == len(batch.BlobMetadata) {
				return fmt.Errorf("HandleSingleBatch: failed to update blob confirmed metadata for all blobs in batch: %w", updateConfirmationInfoErr)
			}
		}

		c.logger.Info("[confirmer] Update confirmation info took", "duration", time.Since(stageTimer))
		c.Metrics.ObserveLatency("UpdateConfirmationInfo", float64(time.Since(stageTimer).Milliseconds()))
		batchSize := int64(0)
		for _, blobMeta := range batch.BlobMetadata {
			batchSize += int64(blobMeta.RequestMetadata.BlobSize)
		}

		c.SliceSigner.RemoveSignedBlob(batchInfo.ts[idx])
		c.EncodingStreamer.RemoveBatchingStatus(batchInfo.ts[idx])
		c.Metrics.IncrementBatchCount(batchSize)
	}

	c.SliceSigner.RemoveBatchingStatus(batchInfo.signedTs)
	return nil
}
