package dispatcher

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/common/storage_node"
	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	"github.com/0glabs/0g-data-avail/disperser/batcher/transactor"
	"github.com/0glabs/0g-data-avail/disperser/contract"
	"github.com/0glabs/0g-data-avail/disperser/contract/da_entrance"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	zg_core "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type Config struct {
	EthClientURL      string
	PrivateKeyString  string
	StorageNodeConfig storage_node.ClientConfig
	UploadTaskSize    uint
}

type dispatcher struct {
	*Config

	daContract     *contract.DAContract
	UploadTaskSize uint

	transactor *transactor.Transactor

	logger common.Logger
}

func NewDispatcher(cfg *Config, transactor *transactor.Transactor, logger common.Logger) (*dispatcher, error) {
	client := blockchain.MustNewWeb3(cfg.EthClientURL, cfg.PrivateKeyString)
	daEntranceAddress := eth_common.HexToAddress(cfg.StorageNodeConfig.DAEntranceContractAddress)
	daSignersAddress := eth_common.HexToAddress(cfg.StorageNodeConfig.DASignersContractAddress)
	daContract, err := contract.NewDAContract(daEntranceAddress, daSignersAddress, client)
	if err != nil {
		return nil, fmt.Errorf("NewDispatcher: failed to create DAEntrance contract: %v", err)
	}

	return &dispatcher{
		Config:         cfg,
		logger:         logger,
		daContract:     daContract,
		UploadTaskSize: cfg.StorageNodeConfig.UploadTaskSize,
		transactor:     transactor,
	}, nil
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func DumpEncodedBlobs(extendedMatrix []*core.ExtendedMatrix) ([]byte, error) {
	blobLocations := make([]*core.BlobLocation, len(extendedMatrix))
	for i, matrix := range extendedMatrix {
		rows := matrix.GetRows()
		cols := matrix.GetCols()
		blobLocations[i] = &core.BlobLocation{
			Rows:           uint(rows),
			Cols:           uint(cols),
			SegmentIndexes: make([]uint, rows),
			Offsets:        make([]uint, rows),
		}
	}
	segmentNum := core.AllocateRows(blobLocations)
	res := make([]byte, segmentNum*core.SegmentSize)
	for i, location := range blobLocations {
		for j := range location.SegmentIndexes {
			offset := location.SegmentIndexes[j]*core.SegmentSize + location.Offsets[j]
			coeffs := extendedMatrix[i].GetRowInBytes(j)
			commitment := extendedMatrix[i].Commitments[j][:]
			copy(res[offset:], coeffs)
			copy(res[offset+uint(len(coeffs)):], commitment)
		}
	}
	return res, nil
}

func (c *dispatcher) DisperseBatch(ctx context.Context, batchHeaderHash [32]byte, batchHeader *core.BatchHeader, blobCommitments []*core.BlobCommitments, blobHeaders []*core.BlobHeader) (eth_common.Hash, error) {
	encoded := make([]*zg_core.DataInMemory, len(blobCommitments))
	for i, commit := range blobCommitments {
		// encoded blobs
		encodedBlobsData, err := zg_core.NewDataInMemory(commit.EncodedData)
		if err != nil {
			return eth_common.Hash{}, errors.WithMessage(err, "failed to build encoded blobs data")
		}

		encoded[i] = encodedBlobsData
	}

	// encoded, err := DumpEncodedBlobs(extendedMatrix)
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "failed to dump encoded blobs")
	// }

	// // encoded blobs
	// encodedBlobsData, err := zg_core.NewDataInMemory(encoded)
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "failed to build encoded blobs data")
	// }

	// // calculate data root
	// tree, err := zg_core.MerkleTree(encodedBlobsData)
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "failed to get encoded data merkle tree")
	// }
	// batchHeader.DataRoot = tree.Root()

	// kv
	// batcher info
	// batcher := c.KVNode.Batcher()
	// blobDisperseInfos := make([]core.BlobDisperseInfo, len(extendedMatrix))
	// for i, matrix := range extendedMatrix {
	// 	blobDisperseInfos[i] = core.BlobDisperseInfo{
	// 		BlobLength: matrix.Length,
	// 		Rows:       uint(matrix.GetRows()),
	// 		Cols:       uint(matrix.GetCols()),
	// 	}
	// }
	// kvBatchInfo := core.KVBatchInfo{
	// 	BatchHeader:       batchHeader,
	// 	BlobDisperseInfos: blobDisperseInfos,
	// }
	// serializedBatchInfo, err := json.Marshal(kvBatchInfo)
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize batch info")
	// }
	// batcher.Set(c.StreamId, batchHeaderHash[:], serializedBatchInfo)
	// // blob info
	// for blobIndex := range extendedMatrix {
	// 	key, err := json.Marshal((&core.KVBlobInfoKey{
	// 		BatchHeaderHash: batchHeaderHash,
	// 		BlobIndex:       uint32(blobIndex),
	// 	}))
	// 	if err != nil {
	// 		return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize kv blob info key")
	// 	}

	// 	value, err := json.Marshal(core.KVBlobInfo{
	// 		BlobHeader: blobHeaders[blobIndex],
	// 		MerkleProof: &core.MerkleProof{
	// 			Hashes: proofs[blobIndex].Hashes,
	// 			Index:  proofs[blobIndex].Index,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize blob info")
	// 	}
	// 	batcher.Set(c.StreamId, key, value)
	// }
	// streamData, err := batcher.Build()
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "Failed to build stream data")
	// }
	// rawKVData, err := streamData.Encode()
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "Failed to encode stream data")
	// }
	// kvData, err := zg_core.NewDataInMemory(rawKVData)
	// if err != nil {
	// 	return eth_common.Hash{}, errors.WithMessage(err, "failed to build kv data")
	// }

	n := len(encoded)
	trees := make([]*merkle.Tree, n)
	dataRoots := make([]eth_common.Hash, n)
	for i := 0; i < n; i++ {
		data := encoded[i]

		c.logger.Info("[dispatcher] Data prepared to upload", "size", data.Size(), "chunks", data.NumChunks(), "segments", data.NumSegments())

		// Calculate file merkle root.
		tree, err := zg_core.MerkleTree(data)
		if err != nil {
			return eth_common.Hash{}, errors.WithMessage(err, "Failed to create data merkle tree")
		}
		c.logger.Info("[dispatcher] Data merkle root calculated", "root", tree.Root())
		trees[i] = tree
		dataRoots[i] = trees[i].Root()

		if eth_common.BytesToHash(blobCommitments[i].StorageRoot) != dataRoots[i] {
			return eth_common.Hash{}, fmt.Errorf("failed to storage root not match: %v", err)
		}
	}

	// upload batchly
	txHash, dataRoots, err := c.transactor.BatchUpload(c.daContract, dataRoots)
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("failed to submit blob data roots: %v", err)
	}

	leafs := make([][]byte, len(dataRoots))
	for i, dataRoot := range dataRoots {

		leafs[i] = dataRoot[:]
	}
	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("failed to get batch data root: %v", err)
	}
	batchHeader.DataRoot = eth_common.Hash(tree.Root())

	return txHash, nil
}

func (c *dispatcher) SubmitAggregateSignatures(ctx context.Context, rootSubmission []*core.CommitRootSubmission) (eth_common.Hash, error) {
	submissions := make([]da_entrance.IDAEntranceCommitRootSubmission, len(rootSubmission))
	for i, s := range rootSubmission {
		submissions[i] = da_entrance.IDAEntranceCommitRootSubmission{
			DataRoot:   s.DataRoot,
			Epoch:      s.Epoch,
			QuorumId:   s.QuorumId,
			CommitRoot: s.CommitRoot,

			QuorumBitmap: s.QuorumBitmap,
			AggPkG2: da_entrance.BN254G2Point{
				X: [2]*big.Int{
					s.AggPkG2.X.A1.BigInt(new(big.Int)),
					s.AggPkG2.X.A0.BigInt(new(big.Int)),
				},
				Y: [2]*big.Int{
					s.AggPkG2.Y.A1.BigInt(new(big.Int)),
					s.AggPkG2.Y.A0.BigInt(new(big.Int)),
				},
			},
			Signature: da_entrance.BN254G1Point{
				X: s.AggSigs.X.BigInt(new(big.Int)),
				Y: s.AggSigs.Y.BigInt(new(big.Int)),
			},
		}
	}

	txHash, err := c.transactor.SubmitVerifiedCommitRoots(c.daContract, submissions)
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("failed to submit verified commit roots: %v", err)
	}

	return txHash, nil
}
