package dispatcher

import (
	"context"
	"fmt"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	zg_core "github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
	"github.com/zero-gravity-labs/zgda/common"
	"github.com/zero-gravity-labs/zgda/common/storage_node"
	"github.com/zero-gravity-labs/zgda/core"
	"github.com/zero-gravity-labs/zgda/disperser"
)

type Config struct {
	EthClientURL      string
	PrivateKeyString  string
	StorageNodeConfig storage_node.ClientConfig
}

type dispatcher struct {
	*Config

	Flow  *contract.FlowContract
	Nodes []*node.Client

	logger common.Logger
}

func NewDispatcher(cfg *Config, logger common.Logger) (*dispatcher, error) {
	client := blockchain.MustNewWeb3(cfg.EthClientURL, cfg.PrivateKeyString)
	contractAddr := eth_common.HexToAddress(cfg.StorageNodeConfig.FlowContractAddress)
	flow, err := contract.NewFlowContract(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("NewDispatcher: failed to create flow contract: %v", err)
	}
	nodes := node.MustNewClients(cfg.StorageNodeConfig.StorageNodeURLs)

	return &dispatcher{
		Config: cfg,
		logger: logger,
		Flow:   flow,
		Nodes:  nodes,
	}, nil
}

var _ disperser.Dispatcher = (*dispatcher)(nil)

func DumpEncodedBlobs(blobs []*core.EncodedBlob) ([]byte, error) {
	res := make([]byte, 0)
	for _, blob := range blobs {
		for _, chunk := range blob.Bundles {
			chunkData, err := chunk.Serialize()
			if err != nil {
				return nil, err
			}
			res = append(res, chunkData...)
		}
	}
	return res, nil
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

func (c *dispatcher) DisperseBatch(ctx context.Context, blobs []*core.EncodedBlob) (*node.FileInfo, error) {
	uploader := transfer.NewUploader(c.Flow, c.Nodes)
	opt := transfer.UploadOption{
		Tags:     hexutil.MustDecode("0x"),
		Force:    false,
		Disperse: true,
	}
	encoded, err := DumpEncodedBlobs(blobs)
	if err != nil {
		return nil, fmt.Errorf("NewClient: cannot get chainId: %w", err)
	}

	data := zg_core.NewDataInMemory(encoded)
	tree, err := zg_core.MerkleTree(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to create file merkle tree: %v", err)
	}

	if err := uploader.Upload(data, opt); err != nil {
		return nil, fmt.Errorf("Failed to upload file: %v", err)
	}

	info, err := c.Nodes[0].ZeroGStorage().GetFileInfo(tree.Root())
	if err != nil {
		return nil, fmt.Errorf("Failed to get file info: %v", err)
	}
	// check finalization
	if !info.Finalized {
		return nil, fmt.Errorf("Tx not finalized: %v", info)
	}
	return info, nil
}
