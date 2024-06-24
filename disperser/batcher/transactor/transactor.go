package transactor

import (
	"sync"
	"time"

	"github.com/0glabs/0g-data-avail/common"
	"github.com/0glabs/0g-data-avail/disperser/contract"
	"github.com/0glabs/0g-data-avail/disperser/contract/da_entrance"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type Transactor struct {
	mu sync.Mutex

	gasLimit uint64
	logger   common.Logger
}

func NewTransactor(gasLimit uint64, logger common.Logger) *Transactor {
	return &Transactor{
		gasLimit: gasLimit,
		logger:   logger,
	}
}

func (t *Transactor) SubmitLogEntry(daContract *contract.DAContract, dataRoots []eth_common.Hash) (eth_common.Hash, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Append log on blockchain
	var txHash eth_common.Hash
	var err error
	if txHash, _, err = daContract.SubmitOriginalData(dataRoots, false); err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to submit log entry")
	}
	return txHash, nil
}

func (t *Transactor) BatchUpload(daContract *contract.DAContract, dataRoots []eth_common.Hash) (eth_common.Hash, error) {
	stageTimer := time.Now()

	txHash, err := t.SubmitLogEntry(daContract, dataRoots)
	if err != nil {
		return eth_common.Hash{}, err
	}

	t.logger.Info("[transactor] batch upload took", "duration", time.Since(stageTimer))

	return txHash, nil
}

func (t *Transactor) SubmitVerifiedCommitRoots(daContract *contract.DAContract, submissions []da_entrance.IDAEntranceCommitRootSubmission) (eth_common.Hash, error) {
	stageTimer := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	var txHash eth_common.Hash
	var err error

	if txHash, _, err = daContract.SubmitVerifiedCommitRoots(submissions, t.gasLimit, false); err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to submit verified commit roots")
	}

	t.logger.Debug("[transactor] submit verified commit roots took", "duration", time.Since(stageTimer))

	return txHash, nil
}
