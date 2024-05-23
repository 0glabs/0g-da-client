package contract

import (
	"math/big"
	"time"

	"github.com/0glabs/0g-data-avail/disperser/contract/da_entrance"
	"github.com/0glabs/0g-data-avail/disperser/contract/da_signers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth_common "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/interfaces"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
)

const DataUploadEventHash = "0x14875b7d8815073fac5ab87ee8f53431fc634767659e261c78df3440973003cb"

var Web3LogEnabled bool

var CustomGasPrice uint64
var CustomGasLimit uint64

type RetryOption struct {
	Rounds   uint
	Interval time.Duration
}

type DataUploadEvent struct {
	DataRoot [32]byte
	Epoch    *big.Int
	QuorumId *big.Int
}

type DAContract struct {
	*da_entrance.DAEntrance
	*da_signers.DASigners
	client  *web3go.Client
	account eth_common.Address // account to send transaction
	signer  bind.SignerFn
}

func defaultSigner(clientWithSigner *web3go.Client) (interfaces.Signer, error) {
	sm, err := clientWithSigner.GetSignerManager()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get signer manager from client")
	}

	if sm == nil {
		return nil, errors.New("Signer not specified")
	}

	signers := sm.List()
	if len(signers) == 0 {
		return nil, errors.WithMessage(err, "Account not configured in signer manager")
	}

	return signers[0], nil
}

func NewDAContract(daEntranceAddress, daSignersAddress eth_common.Address, clientWithSigner *web3go.Client) (*DAContract, error) {
	backend, signer := clientWithSigner.ToClientForContract()

	default_signer, err := defaultSigner(clientWithSigner)
	if err != nil {
		return nil, err
	}

	flow, err := da_entrance.NewDAEntrance(daEntranceAddress, backend)
	if err != nil {
		return nil, err
	}

	signers, err := da_signers.NewDASigners(daSignersAddress, backend)
	if err != nil {
		return nil, err
	}

	return &DAContract{
		DAEntrance: flow,
		DASigners:  signers,
		client:     clientWithSigner,
		account:    default_signer.Address(),
		signer:     signer,
	}, nil
}

func (c *DAContract) SubmitVerifiedCommitRoots(submissions []da_entrance.IDAEntranceCommitRootSubmission, waitForReceipt bool) (eth_common.Hash, *types.Receipt, error) {
	opts, err := c.CreateTransactOpts()
	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	tx, err := c.DAEntrance.SubmitVerifiedCommitRoots(opts, submissions)

	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to send transaction to submit verified commit roots")
	}

	if waitForReceipt {
		// Wait for successful execution
		receipt, err := c.WaitForReceipt(tx.Hash(), true)
		return tx.Hash(), receipt, err
	}
	return tx.Hash(), nil, nil
}

func (c *DAContract) SubmitOriginalData(dataRoots []eth_common.Hash, waitForReceipt bool) (eth_common.Hash, *types.Receipt, error) {
	params := make([][32]byte, 0, len(dataRoots))
	for i, dataRoot := range dataRoots {
		params[i] = dataRoot
	}

	// Submit log entry to smart contract.
	opts, err := c.CreateTransactOpts()
	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	tx, err := c.DAEntrance.SubmitOriginalData(opts, params)

	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to send transaction to submit original data")
	}

	if waitForReceipt {
		// Wait for successful execution
		receipt, err := c.WaitForReceipt(tx.Hash(), true)
		return tx.Hash(), receipt, err
	}
	return tx.Hash(), nil, nil
}

func (c *DAContract) CreateTransactOpts() (*bind.TransactOpts, error) {
	var gasPrice *big.Int
	if CustomGasPrice > 0 {
		gasPrice = new(big.Int).SetUint64(CustomGasPrice)
	}

	return &bind.TransactOpts{
		From:     c.account,
		GasPrice: gasPrice,
		GasLimit: CustomGasLimit,
		Signer:   c.signer,
	}, nil
}

func (c *DAContract) WaitForReceipt(txHash eth_common.Hash, successRequired bool, opts ...RetryOption) (*types.Receipt, error) {
	return WaitForReceipt(c.client, txHash, successRequired, opts...)
}

func WaitForReceipt(client *web3go.Client, txHash eth_common.Hash, successRequired bool, opts ...RetryOption) (receipt *types.Receipt, err error) {
	var opt RetryOption
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		// default infinite wait
		opt.Rounds = 0
		opt.Interval = time.Second * 3
	}

	var tries uint
	for receipt == nil {
		if tries > opt.Rounds+1 && opt.Rounds != 0 {
			return nil, errors.New("no receipt after max retries")
		}
		time.Sleep(opt.Interval)
		if receipt, err = client.Eth.TransactionReceipt(txHash); err != nil {
			return nil, err
		}
		tries++
	}

	if receipt.Status == nil {
		return nil, errors.New("Status not found in receipt")
	}

	switch *receipt.Status {
	case gethTypes.ReceiptStatusSuccessful:
		return receipt, nil
	case gethTypes.ReceiptStatusFailed:
		if !successRequired {
			return receipt, nil
		}

		if receipt.TxExecErrorMsg == nil {
			return nil, errors.New("Transaction execution failed")
		}

		return nil, errors.Errorf("Transaction execution failed, %v", *receipt.TxExecErrorMsg)
	default:
		return nil, errors.Errorf("Unknown receipt status %v", *receipt.Status)
	}
}

func ConvertToGethLog(log *types.Log) *gethTypes.Log {
	if log == nil {
		return nil
	}

	return &gethTypes.Log{
		Address:     log.Address,
		Topics:      log.Topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash,
		TxIndex:     log.TxIndex,
		BlockHash:   log.BlockHash,
		Index:       log.Index,
		Removed:     log.Removed,
	}
}
