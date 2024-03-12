package transactor

import (
	"sync"
	"time"

	"github.com/0glabs/0g-data-avail/common"
	zg_core "github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/transfer"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Transactor struct {
	mu sync.Mutex

	logger common.Logger
}

func NewTransactor(logger common.Logger) *Transactor {
	return &Transactor{
		logger: logger,
	}
}

func (t *Transactor) SubmitLogEntry(uploader *transfer.Uploader, datas []zg_core.IterableData, tags [][]byte) (eth_common.Hash, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Append log on blockchain
	var txHash eth_common.Hash
	var err error
	if txHash, _, err = uploader.SubmitLogEntry(datas, tags, false); err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to submit log entry")
	}
	return txHash, nil
}

func (t *Transactor) BatchUpload(uploader *transfer.Uploader, datas []zg_core.IterableData, opts []transfer.UploadOption) (eth_common.Hash, []eth_common.Hash, error) {
	stageTimer := time.Now()

	n := len(datas)
	trees := make([]*merkle.Tree, n)
	toSubmitDatas := make([]zg_core.IterableData, 0)
	toSubmitTags := make([][]byte, 0)
	dataRoots := make([]eth_common.Hash, n)
	for i := 0; i < n; i++ {
		data := datas[i]
		opt := opts[i]

		t.logger.Info("[transactor] Data prepared to upload", "size", data.Size(), "chunks", data.NumChunks(), "segments", data.NumSegments())

		// Calculate file merkle root.
		tree, err := zg_core.MerkleTree(data)
		if err != nil {
			return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to create data merkle tree")
		}
		t.logger.Info("[transactor] Data merkle root calculated", "root", tree.Root())
		trees[i] = tree
		dataRoots[i] = trees[i].Root()

		toSubmitDatas = append(toSubmitDatas, data)
		toSubmitTags = append(toSubmitTags, opt.Tags)
	}

	txHash, err := t.SubmitLogEntry(uploader, toSubmitDatas, toSubmitTags)
	if err != nil {
		return eth_common.Hash{}, nil, err
	}

	for i := 0; i < n; i++ {
		// Upload file to storage node
		if err := uploader.UploadFile(datas[i], trees[i], 0, opts[i].Disperse, opts[i].TaskSize); err != nil {
			return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to upload file")
		}
	}

	t.logger.Info("[transactor] batch upload took", "duration", time.Since(stageTimer))

	return txHash, dataRoots, nil
}
