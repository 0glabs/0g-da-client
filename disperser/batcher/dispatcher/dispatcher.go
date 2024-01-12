package dispatcher

import (
	"context"
	"fmt"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-merkletree"
	"github.com/zero-gravity-labs/zerog-data-avail/common"
	"github.com/zero-gravity-labs/zerog-data-avail/common/storage_node"
	"github.com/zero-gravity-labs/zerog-data-avail/core"
	"github.com/zero-gravity-labs/zerog-data-avail/disperser"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	zg_core "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
)

type Config struct {
	EthClientURL      string
	PrivateKeyString  string
	StorageNodeConfig storage_node.ClientConfig
}

type dispatcher struct {
	*Config

	Flow     *contract.FlowContract
	Nodes    []*node.Client
	KVNode   *kv.Client
	StreamId eth_common.Hash

	logger common.Logger
}

func NewDispatcher(cfg *Config, logger common.Logger) (*dispatcher, error) {
	client := blockchain.MustNewWeb3(cfg.EthClientURL, cfg.PrivateKeyString)
	contractAddr := eth_common.HexToAddress(cfg.StorageNodeConfig.FlowContractAddress)
	flow, err := contract.NewFlowContract(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("NewDispatcher: failed to create flow contract: %v", err)
	}

	return &dispatcher{
		Config:   cfg,
		logger:   logger,
		Flow:     flow,
		Nodes:    node.MustNewClients(cfg.StorageNodeConfig.StorageNodeURLs),
		KVNode:   kv.NewClient(node.MustNewClient(cfg.StorageNodeConfig.KVNodeURL), nil),
		StreamId: cfg.StorageNodeConfig.KVStreamId,
	}, nil
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func DumpEncodedBlobs(blobs []*core.EncodedBlob) ([]byte, []uint, []uint, error) {
	res := make([]byte, 0)
	chunkOffsets := make([]uint, 0)
	offset := 0
	for _, blob := range blobs {
		empty := core.SegmentSize - offset%core.SegmentSize
		// check padding, we can just do this check before looping over chunks
		if empty < core.SegmentSize {
			if len(blob.Bundles[0].Coeffs)*core.CoeffSize > empty {
				offset += empty
				zeros := make([]byte, empty)
				res = append(res, zeros...)
			}
		}
		chunkOffsets = append(chunkOffsets, uint(offset))
		for _, chunk := range blob.Bundles {
			bytes := chunk.CoeffsToBytes()
			res = append(res, bytes...)
			offset += len(bytes)
		}
	}
	proofOffsets := make([]uint, 0)
	empty := core.SegmentSize - offset%core.SegmentSize
	for _, blob := range blobs {
		proofOffsets = append(proofOffsets, uint(offset))
		for _, chunk := range blob.Bundles {
			// check padding
			if core.ProofSize > empty && empty > 0 {
				offset += empty
				zeros := make([]byte, empty)
				res = append(res, zeros...)
			}
			empty = (empty + core.SegmentSize - core.CoeffSize) % core.SegmentSize
			// insert proof
			bytes := chunk.ProofToBytes()
			res = append(res, bytes[:]...)
			offset += len(bytes)
		}
	}
	return res, chunkOffsets, proofOffsets, nil
}

func GetSegRoot(path string) (eth_common.Hash, error) {
	file, err := zg_core.Open(path)
	defer file.Close()
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("Failed to open file %v: %v", path, err)
	}
	tree, err := zg_core.MerkleTree(file)
	if err != nil {
		return eth_common.Hash{}, fmt.Errorf("Failed to create file merkle tree %v: %v", file, err)
	}
	return tree.Root(), nil
}

func (c *dispatcher) DisperseBatch(ctx context.Context, batchHeaderHash [32]byte, batchHeader *core.BatchHeader, blobs []*core.EncodedBlob, proofs []*merkletree.Proof) (eth_common.Hash, error) {
	uploader := transfer.NewUploader(c.Flow, c.Nodes)
	encoded, chunkOffsets, proofOffsets, err := DumpEncodedBlobs(blobs)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "NewClient: cannot get chainId: %w")
	}

	// encoded blobs
	encodedBlobsData := zg_core.NewDataInMemory(encoded)
	tree, err := zg_core.MerkleTree(encodedBlobsData)
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "failed to get encoded data merkle tree")
	}
	batchHeader.DataRoot = tree.Root()

	// kv
	batcher := c.KVNode.Batcher()
	serializedBatchHeader, err := batchHeader.Serialize()
	if err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize batch header")
	}
	batcher.Set(c.StreamId, batchHeaderHash[:], serializedBatchHeader)
	for blobIndex := range blobs {
		key, err := (&core.KVBlobInfoKey{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       uint32(blobIndex),
		}).Serialize()
		if err != nil {
			return eth_common.Hash{}, errors.WithMessage(err, "Failed to serialize blob info key")
		}
		value, err := (&core.KVBlobInfo{
			BlobHeader:  blobs[blobIndex].BlobHeader,
			MerkleProof: proofs[blobIndex],
			ChunkOffset: chunkOffsets[blobIndex],
			ProofOffset: proofOffsets[blobIndex],
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
	kvData := zg_core.NewDataInMemory(rawKVData)

	// upload batchly
	if err := uploader.BatchUpload([]zg_core.IterableData{encodedBlobsData, kvData}, []transfer.UploadOption{
		// encoded blobs options
		{
			Tags:     hexutil.MustDecode("0x"),
			Force:    true,
			Disperse: true,
		},
		// kv options
		{
			Tags:     batcher.BuildTags(),
			Force:    true,
			Disperse: false,
		}}); err != nil {
		return eth_common.Hash{}, fmt.Errorf("Failed to upload file: %v", err)
	}

	return tree.Root(), nil
}
