package contract

import (
	"math/big"
	"time"

	"github.com/0glabs/0g-da-client/disperser/contract/da_entrance"
	"github.com/0glabs/0g-da-client/disperser/contract/da_signers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth_common "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/interfaces"
	"github.com/openweb3/web3go/signers"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const DataUploadEventHash = "0x57b8b1a6583dc6ce934dfba3d66f2a8e1591b6e171bb2e0921cc64640277087b"

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

func NewDAContract(daEntranceAddress, daSignersAddress eth_common.Address, rpcURL, privateKeyString string) (*DAContract, error) {
	clientWithSigner := MustNewWeb3(rpcURL, privateKeyString)
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

func (c *DAContract) SubmitVerifiedCommitRoots(submissions []da_entrance.IDAEntranceCommitRootSubmission, gasLimit uint64, waitForReceipt bool, estimateGas bool) (*types.Transaction, *types.Receipt, error) {
	opts, err := c.CreateTransactOpts()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	if estimateGas {
		opts.NoSend = estimateGas
	} else {
		opts.GasLimit = gasLimit
	}

	tx, err := c.DAEntrance.SubmitVerifiedCommitRoots(opts, submissions)

	if err != nil {
		return nil, nil, errors.WithMessage(err, "Failed to send transaction to submit verified commit roots")
	}

	if waitForReceipt {
		// Wait for successful execution
		receipt, err := c.WaitForReceipt(tx.Hash(), true)
		return tx, receipt, err
	}
	return tx, nil, nil
}

func (c *DAContract) SubmitOriginalData(dataRoots []eth_common.Hash, waitForReceipt bool) (eth_common.Hash, *types.Receipt, error) {
	params := make([][32]byte, len(dataRoots))
	for i, dataRoot := range dataRoots {
		params[i] = dataRoot
	}

	blobPrice, err := c.BlobPrice(nil)
	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to get blob price")
	}

	// Submit log entry to smart contract.
	opts, err := c.CreateTransactOpts()
	if err != nil {
		return eth_common.Hash{}, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	opts.Value = new(big.Int)
	txValue := new(big.Int).SetUint64(uint64(len(dataRoots)))
	opts.Value.Mul(blobPrice, txValue)

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

func MustNewWeb3(url, key string) *web3go.Client {
	client, err := NewWeb3(url, key)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Fatal("Failed to connect to fullnode")
	}

	return client
}

func NewWeb3(url, key string) (*web3go.Client, error) {
	sm := signers.MustNewSignerManagerByPrivateKeyStrings([]string{key})

	option := new(web3go.ClientOption).
		WithTimout(60 * time.Second).
		WithSignerManager(sm)

	if Web3LogEnabled {
		option = option.WithLooger(logrus.StandardLogger().Out)
	}

	return web3go.NewClientWithOption(url, *option)
}
