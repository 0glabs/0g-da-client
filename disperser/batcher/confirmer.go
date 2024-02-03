package batcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/geth"
	"github.com/zero-gravity-labs/zerog-data-avail/common/storage_node"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

type Confirmer struct {
	mu sync.RWMutex

	Queue            disperser.BlobStore
	EncodingStreamer *EncodingStreamer

	Flow        *contract.FlowContract
	ConfirmChan chan *BatchInfo

	pendingBatches       []*BatchInfo
	MaxNumRetriesPerBlob uint

	routines uint

	retryOption blockchain.RetryOption

	logger  common.Logger
	Metrics *Metrics
}

type BatchInfo struct {
	headerHash [32]byte
	batch      *batch
	ts         uint64
	proofs     []*merkletree.Proof
}

func NewConfirmer(ethConfig geth.EthClientConfig, storageNodeConfig storage_node.ClientConfig, queue disperser.BlobStore, maxNumRetriesPerBlob uint, routines uint, logger common.Logger, metrics *Metrics) (*Confirmer, error) {
	client := blockchain.MustNewWeb3(ethConfig.RPCURL, ethConfig.PrivateKeyString)
	contractAddr := eth_common.HexToAddress(storageNodeConfig.FlowContractAddress)
	flow, err := contract.NewFlowContract(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("NewConfirmer: failed to create flow contract: %v", err)
	}

	if ethConfig.TxGasLimit > 0 {
		blockchain.CustomGasLimit = uint64(ethConfig.TxGasLimit)
	}

	return &Confirmer{
		Queue:          queue,
		Flow:           flow,
		ConfirmChan:    make(chan *BatchInfo),
		pendingBatches: make([]*BatchInfo, 0),
		routines:       routines,
		retryOption: blockchain.RetryOption{
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
	c.logger.Info(`[confirmer] retreived one pending batch`, "queue size", len(c.pendingBatches))
	return info
}

func (c *Confirmer) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	for _, metadata := range blobMetadatas {
		err := c.Queue.HandleBlobFailure(ctx, metadata, c.MaxNumRetriesPerBlob)
		if err != nil {
			c.logger.Error("HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}
		c.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
	}
	c.Metrics.UpdateBatchError(reason, len(blobMetadatas))

	// Return the error(s)
	return result.ErrorOrNil()
}

func (c *Confirmer) waitForReceipt(txHash eth_common.Hash) (uint32, uint32, error) {
	if txHash.Cmp(eth_common.Hash{}) == 0 {
		return 0, 0, errors.New("empty transaction hash")
	}
	c.logger.Info("[confirmer] Waiting batch be confirmed", "transaction hash", txHash)
	// data is not duplicate, there is a new transaction
	receipt, err := c.Flow.WaitForReceipt(txHash, true, c.retryOption)
	if err != nil {
		return 0, 0, err
	}

	submitEventHash := eth_common.HexToHash(transfer.SubmitEventHash)
	var submission *contract.FlowSubmit
	// parse submission from event log
	for _, v := range receipt.Logs {
		if v.Topics[0] == submitEventHash {
			log := blockchain.ConvertToGethLog(v)

			if submission, err = c.Flow.ParseSubmit(*log); err != nil {
				return 0, 0, err
			}

			break
		}
	}
	return uint32(submission.SubmissionIndex.Uint64()), uint32(receipt.BlockNumber), nil
}

func (c *Confirmer) ConfirmBatch(ctx context.Context, batchInfo *BatchInfo) error {
	batch := batchInfo.batch
	proofs := batchInfo.proofs

	txSeq, blockNumber, err := c.waitForReceipt(batch.TxHash)
	if err != nil {
		// batch is not confirmed
		c.EncodingStreamer.RemoveBatchingStatus(batchInfo.ts)
		return err
	}

	batchID := txSeq
	c.logger.Info("[confirmer] batch confirmed.", "batch ID", batchID, "transaction hash", batch.TxHash)
	// Mark the blobs as complete
	c.logger.Info("[confirmer] Marking blobs as complete...")
	stageTimer := time.Now()
	blobsToRetry := make([]*disperser.BlobMetadata, 0)
	var updateConfirmationInfoErr error
	for blobIndex, metadata := range batch.BlobMetadata {
		confirmationInfo := &disperser.ConfirmationInfo{
			BatchHeaderHash:         batchInfo.headerHash,
			BlobIndex:               uint32(blobIndex),
			ReferenceBlockNumber:    uint32(batch.BatchHeader.ReferenceBlockNumber),
			BatchRoot:               batch.BatchHeader.BatchRoot[:],
			BlobInclusionProof:      serializeProof(proofs[blobIndex]),
			BlobCommitment:          &batch.BlobHeaders[blobIndex].BlobCommitments,
			BatchID:                 uint32(batchID),
			ConfirmationTxnHash:     batch.TxHash,
			ConfirmationBlockNumber: blockNumber,
		}
		c.logger.Trace("confirming blob", "blob key", metadata.GetBlobKey())
		if _, updateConfirmationInfoErr = c.Queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo); updateConfirmationInfoErr == nil {
			c.Metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Confirmed)
			// remove encoded blob from storage so we don't disperse it again
			c.EncodingStreamer.RemoveEncodedBlob(metadata)
			c.logger.Trace("blob confirmed", "blob key", metadata.GetBlobKey())
		}
		if updateConfirmationInfoErr != nil {
			c.logger.Error("HandleSingleBatch: error updating blob confirmed metadata", "err", updateConfirmationInfoErr)
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
	c.Metrics.IncrementBatchCount(batchSize)
	return nil
}
