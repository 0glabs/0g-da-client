package batcher

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/common/geth"
	"github.com/0glabs/0g-data-avail/common/storage_node"
	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	pb "github.com/0glabs/0g-data-avail/disperser/api/grpc/signer"
	"github.com/0glabs/0g-data-avail/disperser/contract"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hashicorp/go-multierror"
	"github.com/wealdtech/go-merkletree"
)

var errNoSignedResults = errors.New("no signed results")

type SignatureSizeNotifier struct {
	mu sync.Mutex

	Notify chan struct{}
	// threshold is the size of the total encoded blob results in bytes that triggers the notifier
	threshold uint64
	// active is set to false after the notifier is triggered to prevent it from triggering again for the same batch
	// This is reset when CreateBatch is called and the encoded results have been consumed
	active bool
}

func NewSignatureSizeNotifier(notify chan struct{}, threshold uint64) *SignatureSizeNotifier {
	return &SignatureSizeNotifier{
		Notify:    notify,
		threshold: threshold,
		active:    true,
	}
}

type SignerConfig struct {
	// EncodingRequestTimeout is the timeout for each encoding request
	SigningRequestTimeout time.Duration

	// EncodingQueueLimit is the maximum number of encoding requests that can be queued
	EncodingQueueLimit int

	MaxNumRetriesPerBlob uint

	MaxNumRetriesSign uint
}

type SignInfo struct {
	headerHash [32]byte
	batch      *batch
	ts         uint64
	proofs     []*merkletree.Proof
	reties     uint

	epoch    *big.Int
	quorumId *big.Int
	signers  []*SignerInfo
}

type SignResult struct {
	signatures []*core.Signature
	signer     *SignerInfo
}

type SignResultOrStatus struct {
	SignResult
	Err error
}

type SignerInfo struct {
	Signer eth_common.Address
	Socket string
	PkG1   *core.G1Point
	PkG2   *core.G2Point
}

type BatchCommitRootSubmission struct {
	submissions []*core.CommitRootSubmission
	headerHash  [32]byte
	batch       *batch
	ts          uint64
	proofs      []*merkletree.Proof
}

type SliceSigner struct {
	SignerConfig

	mu sync.RWMutex

	Pool                  common.WorkerPool
	SignatureSizeNotifier *SignatureSizeNotifier
	SignerChan            chan *SignInfo

	EncodingStreamer *EncodingStreamer
	Finalizer        Finalizer

	pendingBatches       []*SignInfo
	pendingBatchesToSign []*SignInfo
	pendingSubmissions   map[uint64]*BatchCommitRootSubmission
	signedBatching       map[uint64]uint64
	signedBatches        map[uint64][]uint64
	signedBlobSize       uint64

	daContract   *contract.DAContract
	signerClient disperser.SignerClient

	retryOption contract.RetryOption

	blobStore disperser.BlobStore
	metrics   *Metrics

	logger common.Logger
}

func NewEncodedSliceSigner(
	ethConfig geth.EthClientConfig,
	storageNodeConfig storage_node.ClientConfig,
	config SignerConfig,
	workerPool common.WorkerPool,
	signatureSizeNotifier *SignatureSizeNotifier,
	signerClient disperser.SignerClient,
	blobStore disperser.BlobStore,
	metrics *Metrics,
	logger common.Logger) (*SliceSigner, error) {

	client := blockchain.MustNewWeb3(ethConfig.RPCURL, ethConfig.PrivateKeyString)
	daEntranceAddress := eth_common.HexToAddress(storageNodeConfig.DAEntranceContractAddress)
	daSignersAddress := eth_common.HexToAddress(storageNodeConfig.DASignersContractAddress)
	daContract, err := contract.NewDAContract(daEntranceAddress, daSignersAddress, client)
	if err != nil {
		return nil, fmt.Errorf("signer: failed to create DAEntrance contract: %v", err)
	}

	return &SliceSigner{
		SignerConfig:          config,
		Pool:                  workerPool,
		SignatureSizeNotifier: signatureSizeNotifier,
		SignerChan:            make(chan *SignInfo),
		daContract:            daContract,
		signerClient:          signerClient,
		retryOption: contract.RetryOption{
			Rounds:   ethConfig.ReceiptPollingRounds,
			Interval: ethConfig.ReceiptPollingInterval,
		},
		blobStore: blobStore,
		metrics:   metrics,
		logger:    logger,

		pendingBatches:       make([]*SignInfo, 0),
		pendingBatchesToSign: make([]*SignInfo, 0),
		pendingSubmissions:   make(map[uint64]*BatchCommitRootSubmission),
		signedBatching:       make(map[uint64]uint64),
		signedBatches:        make(map[uint64][]uint64),
		signedBlobSize:       0,
	}, nil
}

func (s *SliceSigner) Start(ctx context.Context) error {
	// goroutine for making blob signing requests
	go func() {
		ticker := time.NewTicker(encodingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := s.doSigning(ctx)
				if err != nil {
					s.logger.Warn("[signer] error requesting signing", "err", err)
				}
			}
		}
	}()

	// goroutine for wait tx finalized
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				signInfo := s.getPendingBatch()
				if signInfo != nil {
					err := s.waitBatchTxFinalized(ctx, signInfo)
					if err != nil {
						s.logger.Warn("[signer] error wait batch tx finalized", "err", err)
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case batchInfo := <-s.SignerChan:
				s.putPendingBatches(batchInfo)
			}
		}
	}()

	return nil

}

func (s *SliceSigner) putPendingBatches(info *SignInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pendingBatches = append(s.pendingBatches, info)
	s.logger.Info(`[signer] received pending batch`, "queue size", len(s.pendingBatches))
}

func (s *SliceSigner) getPendingBatch() *SignInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pendingBatches) == 0 {
		return nil
	}
	info := s.pendingBatches[0]
	s.pendingBatches = s.pendingBatches[1:]
	s.logger.Info(`[signer] retrieved one pending batch`, "queue size", len(s.pendingBatches))
	return info
}

func (s *SliceSigner) waitBatchTxFinalized(ctx context.Context, batchInfo *SignInfo) error {
	dataUploadEvents, blockNumber, err := s.waitForReceipt(batchInfo.batch.TxHash)
	s.logger.Debug("[signer] wait for receipt", "epochs", dataUploadEvents, "block number", blockNumber)

	if err != nil {
		// batch is not confirmed
		_ = s.handleFailure(ctx, batchInfo.batch.BlobMetadata, FailBatchReceipt)
		s.EncodingStreamer.RemoveBatchingStatus(batchInfo.ts)
		return err
	}

	for i := 1; i < len(dataUploadEvents); i++ {
		if dataUploadEvents[i].Epoch != dataUploadEvents[i-1].Epoch {
			_ = s.handleFailure(ctx, batchInfo.batch.BlobMetadata, FailBatchEpochMismatch)
			s.EncodingStreamer.RemoveBatchingStatus(batchInfo.ts)
			return fmt.Errorf("error epoch in one batch is mismatch: %w", err)
		}
	}

	epoch := dataUploadEvents[0].Epoch
	quorumId := dataUploadEvents[0].QuorumId
	signers, err := s.getSigners(epoch, quorumId)
	if err != nil {
		// if signInfo.reties < s.MaxNumRetriesSign {
		// 	s.mu.Lock()
		// 	defer s.mu.Unlock()

		// 	signInfo.reties += 1
		// 	s.pendingBatchesToSign = append(s.pendingBatchesToSign, signInfo)
		// } else {
		_ = s.handleFailure(ctx, batchInfo.batch.BlobMetadata, FailGetSigners)
		s.EncodingStreamer.RemoveBatchingStatus(batchInfo.ts)
		// }

		return fmt.Errorf("error getting signers from contract: %w", err)
	}

	// update epoch
	batchInfo.epoch = epoch
	batchInfo.quorumId = quorumId
	batchInfo.signers = signers

	s.mu.Lock()
	defer s.mu.Unlock()
	s.pendingBatchesToSign = append(s.pendingBatchesToSign, batchInfo)
	s.logger.Info(`[signer] batch finalized`, "batch ts", batchInfo.ts)
	return nil
}

func (s *SliceSigner) waitForReceipt(txHash eth_common.Hash) ([]*contract.DataUploadEvent, uint32, error) {
	if txHash.Cmp(eth_common.Hash{}) == 0 {
		return nil, 0, errors.New("empty transaction hash")
	}
	s.logger.Info("[signer] Waiting batch be confirmed", "transaction hash", txHash)
	// data is not duplicate, there is a new transaction
	var blockNumber uint64
	submitEventHash := eth_common.HexToHash(contract.DataUploadEventHash)
	var submissions []*contract.DataUploadEvent

	for {
		receipt, err := s.daContract.WaitForReceipt(txHash, true, s.retryOption)
		if err != nil {
			return nil, 0, err
		}

		blockNumber = receipt.BlockNumber
		if blockNumber > s.Finalizer.LatestFinalizedBlock() {
			s.logger.Debug("[signer] waiting batch be confirmed")
			time.Sleep(time.Second * 5)
			continue
		}

		// parse submission from event log
		for _, v := range receipt.Logs {
			if v.Topics[0] == submitEventHash {
				log := contract.ConvertToGethLog(v)

				submission, err := s.daContract.ParseDataUpload(*log)
				if err != nil {
					return nil, 0, err
				}

				submissions = append(submissions, &contract.DataUploadEvent{
					DataRoot: submission.DataRoot,
					Epoch:    submission.Id,
					QuorumId: submission.QuorumId,
				})
			}
		}
		break
	}

	return submissions, uint32(blockNumber), nil
}

func (s *SliceSigner) getSigners(epoch *big.Int, quorumId *big.Int) ([]*SignerInfo, error) {
	signerAddresses, err := s.daContract.GetQuorum(nil, epoch, quorumId)
	if err != nil {
		return nil, err
	}

	signers, err := s.daContract.GetSigner(nil, signerAddresses)
	if err != nil {
		return nil, err
	}

	signersInfo := make([]*SignerInfo, 0, len(signers))
	for idx, signer := range signers {
		pubkeyG1 := core.NewG1Point(signer.PkG1.X, signer.PkG1.Y)

		pubkeyG2 := new(bn254.G2Affine)
		pubkeyG2.X.A0.SetBigInt(signer.PkG2.X[1])
		pubkeyG2.X.A1.SetBigInt(signer.PkG2.X[0])
		pubkeyG2.Y.A0.SetBigInt(signer.PkG2.Y[1])
		pubkeyG2.X.A1.SetBigInt(signer.PkG2.X[0])

		signersInfo[idx] = &SignerInfo{
			Signer: signer.Signer,
			Socket: signer.Socket,
			PkG1:   pubkeyG1,
			PkG2:   &core.G2Point{G2Affine: pubkeyG2},
		}
	}

	return signersInfo, nil
}

func (s *SliceSigner) handleFailure(ctx context.Context, blobMetadatas []*disperser.BlobMetadata, reason FailReason) error {
	var result *multierror.Error
	for _, metadata := range blobMetadatas {
		err := s.blobStore.HandleBlobFailure(ctx, metadata, s.MaxNumRetriesPerBlob)
		if err != nil {
			s.logger.Error("[signer] HandleSingleBatch: error handling blob failure", "err", err)
			// Append the error
			result = multierror.Append(result, err)
		}
		s.metrics.UpdateCompletedBlob(int(metadata.RequestMetadata.BlobSize), disperser.Failed)
	}
	s.metrics.UpdateBatchError(reason, len(blobMetadatas))

	// Return the error(s)
	return result.ErrorOrNil()
}

func (s *SliceSigner) getPendingBatchToSign() *SignInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pendingBatchesToSign) == 0 {
		return nil
	}
	info := s.pendingBatchesToSign[0]
	s.pendingBatchesToSign = s.pendingBatchesToSign[1:]
	s.logger.Info(`[signer] retrieved one batch to sign`, "queue size", len(s.pendingBatchesToSign))
	return info
}

func (s *SliceSigner) doSigning(ctx context.Context) error {
	signInfo := s.getPendingBatchToSign()
	signerCounter := len(signInfo.signers)
	requestData := s.assignEncodedBlobs(signerCounter, signInfo.batch, signInfo.epoch.Uint64())

	update := make(chan SignResultOrStatus, signerCounter)
	for i := 0; i < signerCounter; i++ {
		encodingCtx, cancel := context.WithTimeout(ctx, s.SigningRequestTimeout)
		s.Pool.Submit(func() {
			defer cancel()

			// Todo: assume this is no empty EncodedSlice
			// n := 0
			// for _, req := range reqs[i] {
			// 	if len(req.EncodedSlice) != 0 {
			// 		reqs[i][n] = req
			// 		n++
			// 	}
			// }

			reply, err := s.signerClient.BatchSign(encodingCtx, signInfo.signers[i].Socket, requestData[i])
			if err != nil {
				update <- SignResultOrStatus{
					Err:        err,
					SignResult: SignResult{},
				}
				return
			}

			update <- SignResultOrStatus{
				Err: nil,
				SignResult: SignResult{
					signatures: reply,
					signer:     signInfo.signers[i],
				},
			}
		})

		s.logger.Trace("requested sign for batch", "ts", signInfo.ts)
	}

	err := s.aggregateSignature(ctx, signInfo, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *SliceSigner) assignEncodedBlobs(signerCounter int, batch *batch, epoch uint64) [][]*pb.SignRequest {
	requestData := make([][]*pb.SignRequest, 0, signerCounter)
	for signerIdx := 0; signerIdx < signerCounter; signerIdx++ {
		requestData[signerIdx] = make([]*pb.SignRequest, 0, len(batch.EncodedBlobs))

		for blobIdx, encodedBlob := range batch.EncodedBlobs {
			requestData[signerIdx][blobIdx] = &pb.SignRequest{
				Epoch:             epoch,
				ErasureCommitment: encodedBlob.ErasureCommitment,
				StorageRoot:       encodedBlob.StorageRoot,
				EncodedSlice:      make([][]byte, 0),
			}
		}
	}

	for blobIdx, encodedBlobs := range batch.EncodedBlobs {
		for sliceIdx, encodedSlice := range encodedBlobs.EncodedSlice {
			signerIdx := sliceIdx % signerCounter
			requestData[signerIdx][blobIdx].EncodedSlice = append(requestData[signerIdx][blobIdx].EncodedSlice, encodedSlice)
		}
	}

	return requestData
}

func (s *SliceSigner) aggregateSignature(ctx context.Context, signInfo *SignInfo, update chan SignResultOrStatus) error {
	signerCounter := len(signInfo.signers)
	signersIndex := make(map[eth_common.Address]int)
	for signerIdx, s := range signInfo.signers {
		signersIndex[s.Signer] = signerIdx
	}

	blobSize := len(signInfo.batch.EncodedBlobs)
	erasureCommitments := make([][32]byte, 0, blobSize)
	storageRoots := make([][32]byte, 0, blobSize)
	messages := make([][32]byte, 0, blobSize)
	for blobIdx, encodedBlobs := range signInfo.batch.EncodedBlobs {
		var dataRoot [32]byte
		copy(dataRoot[:], encodedBlobs.StorageRoot[:32])
		var commitRoot [32]byte
		copy(commitRoot[:], encodedBlobs.ErasureCommitment[:32])

		erasureCommitments[blobIdx] = commitRoot
		storageRoots[blobIdx] = dataRoot
		messages[blobIdx] = encodedBlobs.GetHash()
	}

	aggSigs := make([]*core.Signature, blobSize)
	aggPubKeys := make([]*core.G2Point, blobSize)

	bitmapLen := signerCounter / 8
	if signerCounter%8 != 0 {
		bitmapLen++
	}
	signatureCounts := make([]int, blobSize)
	quorumBitmap := make([][]byte, 0, blobSize)

	for i := 0; i < signerCounter; i++ {
		recv := <-update
		signer := recv.signer
		signatures := recv.signatures

		if recv.Err != nil {
			s.logger.Warn("error returned from messageChan", "socket", signer.Socket, "err", recv.Err)
			continue
		}

		for blobIdx, sig := range signatures {
			message := messages[blobIdx]

			// Verify Signature
			ok := sig.Verify(signer.PkG2, message)
			if !ok {
				s.logger.Error("signature is not valid", "signerAddress", signer.Signer, "socket", signer.Socket, "pubkey", hexutil.Encode(signer.PkG2.Serialize()))
				continue
			}

			if aggSigs[blobIdx] == nil {
				aggSigs[blobIdx] = &core.Signature{G1Point: sig.Clone()}
				aggPubKeys[blobIdx] = signer.PkG2.Clone()

				signatureCounts[blobIdx] = 0
				quorumBitmap[blobIdx] = make([]byte, bitmapLen)
			} else {
				aggSigs[blobIdx].Add(sig.G1Point)
				aggPubKeys[blobIdx].Add(signer.PkG2)
			}

			signatureCounts[blobIdx]++
			slot := signersIndex[signer.Signer] / 8
			offset := signersIndex[signer.Signer] % 8
			quorumBitmap[blobIdx][slot] |= 1 << offset
		}
	}

	threshold := int(math.Ceil(float64(signerCounter) * 2 / 3))
	valid := true
	rootSubmissions := make([]*core.CommitRootSubmission, 0)
	for blobIdx, sig := range aggSigs {
		if signatureCounts[blobIdx] < threshold {
			valid = false
			break
		}

		rootSubmissions = append(rootSubmissions, &core.CommitRootSubmission{
			DataRoot:     storageRoots[blobIdx],
			CommitRoot:   erasureCommitments[blobIdx],
			Epoch:        signInfo.epoch,
			QuorumId:     signInfo.quorumId,
			QuorumBitmap: quorumBitmap[blobIdx],
			AggPkG2:      aggPubKeys[blobIdx],
			AggSigs:      sig,
		})
	}

	if valid {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.pendingSubmissions[signInfo.ts] = &BatchCommitRootSubmission{
			submissions: rootSubmissions,
			headerHash:  signInfo.headerHash,
			batch:       signInfo.batch,
			ts:          signInfo.ts,
			proofs:      signInfo.proofs,
		}
		s.signedBlobSize += uint64(len(signInfo.batch.EncodedBlobs))

		s.metrics.UpdateSignedBlobs(len(s.pendingSubmissions), s.signedBlobSize)

		if s.SignatureSizeNotifier.threshold > 0 && s.signedBlobSize > s.SignatureSizeNotifier.threshold {
			s.SignatureSizeNotifier.mu.Lock()

			if s.SignatureSizeNotifier.active {
				s.SignatureSizeNotifier.Notify <- struct{}{}
				s.SignatureSizeNotifier.active = false
			}
			s.SignatureSizeNotifier.mu.Unlock()
		}
	} else {
		if signInfo.reties < s.MaxNumRetriesSign {
			s.mu.Lock()
			defer s.mu.Unlock()

			signInfo.reties += 1
			s.pendingBatchesToSign = append(s.pendingBatchesToSign, signInfo)
		} else {
			_ = s.handleFailure(ctx, signInfo.batch.BlobMetadata, FailAggregateSignatures)
			s.EncodingStreamer.RemoveBatchingStatus(signInfo.ts)
			return errors.New("failed aggregate signatures")
		}
	}

	return nil
}

func (s *SliceSigner) GetCommitRootSubmissionBatch() ([]*BatchCommitRootSubmission, uint64, error) {
	ts := uint64(time.Now().Nanosecond())

	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pendingSubmissions) == 0 {
		return nil, ts, errNoSignedResults
	}

	// Reset the notifier
	s.SignatureSizeNotifier.mu.Lock()
	s.SignatureSizeNotifier.active = true
	s.SignatureSizeNotifier.mu.Unlock()

	fetched := make([]*BatchCommitRootSubmission, 0)
	if _, ok := s.signedBatches[ts]; !ok {
		s.signedBatches[ts] = make([]uint64, 0)
	}
	for id, signedResult := range s.pendingSubmissions {
		if _, ok := s.signedBatching[id]; !ok {
			fetched = append(fetched, signedResult)
			s.signedBatching[id] = ts
			s.signedBatches[ts] = append(s.signedBatches[ts], id)
		}
	}
	s.logger.Trace("consumed signed results", "fetched", len(fetched))

	if len(fetched) == 0 {
		return nil, ts, errNoSignedResults
	}

	// n := len(s.pendingSubmissions)

	// info := make([]BatchedCommitRootSubmission, 0)
	// info = append(info, s.pendingSubmissions[:n]...)
	// s.pendingSubmissions = s.pendingSubmissions[n:]
	// c.logger.Info(`[confirmer] retrieved one pending batch`, "queue size", len(c.pendingBatches))
	return fetched, ts, nil

}

func (s *SliceSigner) RemoveSignedBlob(ts uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.pendingSubmissions[ts]
	if !ok {
		return
	}

	s.signedBlobSize -= uint64(len(s.pendingSubmissions[ts].batch.EncodedBlobs))
	delete(s.pendingSubmissions, ts)
	// remove from batching status
	delete(s.signedBatching, ts)
}

func (s *SliceSigner) RemoveBatchingStatus(ts uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if batch, ok := s.signedBatches[ts]; ok {
		for _, id := range batch {
			delete(s.signedBatching, id)
		}
	}
	delete(s.signedBatches, ts)
}
