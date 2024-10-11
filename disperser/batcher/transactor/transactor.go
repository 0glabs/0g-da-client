package transactor

import (
	"sync"
	"time"

	"github.com/0glabs/0g-da-client/common"
	"github.com/0glabs/0g-da-client/disperser/contract"
	"github.com/0glabs/0g-da-client/disperser/contract/da_entrance"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go/types"
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
	var tx *types.Transaction
	var err error
	if tx, _, err = daContract.SubmitOriginalData(dataRoots, false, 0); err != nil {
		t.logger.Debug("[transactor] estimate SubmitLogEntry tx failed")
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to estimate SubmitLogEntry tx")
	}

	gasLimit := tx.Gas() + tx.Gas()*3/10
	if tx, _, err = daContract.SubmitOriginalData(dataRoots, false, gasLimit); err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to submit log entry")
	}

	return tx.Hash(), nil
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

	var tx *types.Transaction
	var err error

	var gasLimit uint64
	if t.gasLimit == 0 {
		if tx, _, err = daContract.SubmitVerifiedCommitRoots(submissions, 0, false, true); err != nil {
			return eth_common.Hash{}, errors.WithMessage(err, "Failed to estimate SubmitVerifiedCommitRoots")
		}

		gasLimit = tx.Gas() + tx.Gas()*3/10
		t.logger.Info("[transactor] estimate gas", "gas limit", tx.Gas())
	} else {
		gasLimit = t.gasLimit
	}

	if tx, _, err = daContract.SubmitVerifiedCommitRoots(submissions, gasLimit, false, false); err != nil {
		return eth_common.Hash{}, errors.WithMessage(err, "Failed to submit verified commit roots")
	}

	t.logger.Debug("[transactor] submit verified commit roots took", "duration", time.Since(stageTimer))

	return tx.Hash(), nil
}
