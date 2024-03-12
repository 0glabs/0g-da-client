package dispatcher

import (
	"context"
	"fmt"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/common/storage_node"
	"github.com/0glabs/0g-data-avail/core"
	"github.com/0glabs/0g-data-avail/disperser"
	"github.com/0glabs/0g-data-avail/disperser/batcher/transactor"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	zg_core "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
)

type Config struct {
	EthClientURL      string
	PrivateKeyString  string
	StorageNodeConfig storage_node.ClientConfig
	UploadTaskSize    uint
}

type dispatcher struct {
	*Config

	Flow           *contract.FlowContract
	Nodes          []*node.Client
	KVNode         *kv.Client
	StreamId       eth_common.Hash
	UploadTaskSize uint

	transactor *transactor.Transactor

	logger common.Logger
}

func NewDispatcher(cfg *Config, transactor *transactor.Transactor, logger common.Logger) (*dispatcher, error) {
	client := blockchain.MustNewWeb3(cfg.EthClientURL, cfg.PrivateKeyString)
	contractAddr := eth_common.HexToAddress(cfg.StorageNodeConfig.FlowContractAddress)
	flow, err := contract.NewFlowContract(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("NewDispatcher: failed to create flow contract: %v", err)
	}

	return &dispatcher{
		Config:         cfg,
		logger:         logger,
		Flow:           flow,
		Nodes:          node.MustNewClients(cfg.StorageNodeConfig.StorageNodeURLs),
		KVNode:         kv.NewClient(node.MustNewClient(cfg.StorageNodeConfig.KVNodeURL), nil),
		StreamId:       cfg.StorageNodeConfig.KVStreamId,
		UploadTaskSize: cfg.StorageNodeConfig.UploadTaskSize,
		transactor:     transactor,
	}, nil
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func DumpEncodedBlobs(blobs []*core.EncodedBlob) ([]byte, error) {
	blobLocations := make([]*core.BlobLocation, len(blobs))
	for i, blob := range blobs {
		chunkLength := uint(len(blob.Bundles[0].Coeffs))
		chunkNum := uint(len(blob.Bundles))
		blobLocations[i] = &core.BlobLocation{
			ChunkLength:    chunkLength,
			ChunkNum:       chunkNum,
			SegmentIndexes: make([]uint, chunkNum),
			Offsets:        make([]uint, chunkNum),
		}
	}
	segmentNum := core.AllocateChunks(blobLocations)
	res := make([]byte, segmentNum*core.SegmentSize)
	for i, location := range blobLocations {
		for j := range location.SegmentIndexes {
			offset := location.SegmentIndexes[j]*core.SegmentSize + location.Offsets[j]
			coeffs := blobs[i].Bundles[j].CoeffsToBytes()
			proof := blobs[i].Bundles[j].ProofToBytes()
			copy(res[offset:], coeffs)
			copy(res[offset+uint(len(coeffs)):], proof[:])
		}
	}
	return res, nil
}

func (c *dispatcher) DisperseBatch(ctx context.Context, batchHeaderHash [32]byte, batchHeader *core.BatchHeader, blobs []*core.EncodedBlob, proofs []*merkletree.Proof) (eth_common.Hash, error) {
	uploader := transfer.NewUploader(c.Flow, c.Nodes)
	encoded, err := DumpEncodedBlobs(blobs)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "NewClient: cannot get chainId: %w")
	}

	// encoded blobs
	encodedBlobsData, err := zg_core.NewDataInMemory(encoded)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "failed to build encoded blobs data")
	}

	// calculate data root
	tree, err := zg_core.MerkleTree(encodedBlobsData)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "failed to get encoded data merkle tree")
	}
	batchHeader.DataRoot = tree.Root()

	// kv
	// batcher info
	batcher := c.KVNode.Batcher()
	blobDisperseInfos := make([]core.BlobDisperseInfo, len(blobs))
	for i, blob := range blobs {
		blobDisperseInfos[i] = core.BlobDisperseInfo{
			BlobLength:   blob.BlobHeader.Length,
			BlobChunkNum: uint(len(blob.Bundles)),
		}
	}
	kvBatchInfo := core.KVBatchInfo{
		BatchHeader:       batchHeader,
		BlobDisperseInfos: blobDisperseInfos,
	}
	serializedBatchInfo, err := kvBatchInfo.Serialize()
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize batch info")
	}
	batcher.Set(c.StreamId, batchHeaderHash[:], serializedBatchInfo)
	// blob info
	for blobIndex := range blobs {
		key := (&core.KVBlobInfoKey{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       uint32(blobIndex),
		}).Bytes()

		value, err := (&core.KVBlobInfo{
			BlobHeader:  blobs[blobIndex].BlobHeader,
			MerkleProof: proofs[blobIndex],
		}).Serialize()
		if err != nil {
			return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize blob info")
		}
		batcher.Set(c.StreamId, key, value)
	}
	streamData, err := batcher.Build()
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to build stream data")
	}
	rawKVData, err := streamData.Encode()
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to encode stream data")
	}
	kvData, err := zg_core.NewDataInMemory(rawKVData)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "failed to build kv data")
	}

	// upload batchly
	txHash, dataRoots, err := c.transactor.BatchUpload(uploader, []zg_core.IterableData{encodedBlobsData, kvData}, []transfer.UploadOption{
		// encoded blobs options
		{
			Tags:     hexutil.MustDecode("0x"),
			Force:    true,
			Disperse: true,
			TaskSize: c.UploadTaskSize,
		},
		// kv options
		{
			Tags:     batcher.BuildTags(),
			Force:    true,
			Disperse: false,
			TaskSize: c.UploadTaskSize,
		}})
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("failed to upload file: %v", err)
	}
	batchHeader.DataRoot = dataRoots[0]

	return txHash, nil
}
