// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package da_entrance

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// BN254G2Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G2Point struct {
	X [2]*big.Int
	Y [2]*big.Int
}

// IDAEntranceCommitRootSubmission is an auto generated low-level Go binding around an user-defined struct.
type IDAEntranceCommitRootSubmission struct {
	DataRoot          [32]byte
	Epoch             *big.Int
	QuorumId          *big.Int
	ErasureCommitment BN254G1Point
	QuorumBitmap      []byte
	AggPkG2           BN254G2Point
	Signature         BN254G1Point
}

// IDASampleSampleRange is an auto generated low-level Go binding around an user-defined struct.
type IDASampleSampleRange struct {
	StartEpoch uint64
	EndEpoch   uint64
}

// IDASampleSampleTask is an auto generated low-level Go binding around an user-defined struct.
type IDASampleSampleTask struct {
	SampleHash      [32]byte
	PodasTarget     *big.Int
	RestSubmissions uint64
}

// SampleResponse is an auto generated low-level Go binding around an user-defined struct.
type SampleResponse struct {
	SampleSeed   [32]byte
	Epoch        uint64
	QuorumId     uint64
	LineIndex    uint32
	SublineIndex uint32
	Quality      *big.Int
	DataRoot     [32]byte
	BlobRoots    [3][32]byte
	Proof        [][32]byte
	Data         []byte
}

// DAEntranceMetaData contains all meta data concerning the DAEntrance contract.
var DAEntranceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sampleRound\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quality\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"lineIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sublineIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"DAReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blobPrice\",\"type\":\"uint256\"}],\"name\":\"DataUpload\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"}],\"name\":\"ErasureCommitmentVerified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sampleRound\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sampleHeight\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"sampleSeed\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"podasTarget\",\"type\":\"uint256\"}],\"name\":\"NewSampleRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DA_SIGNERS\",\"outputs\":[{\"internalType\":\"contractIDASigners\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_PODAS_TARGET\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PARAMS_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLICE_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLICE_NUMERATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activedReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"blobPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_quorumId\",\"type\":\"uint256\"}],\"name\":\"commitmentExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpochReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentSampleSeed\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"donate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochWindowSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextSampleHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"dest\",\"type\":\"address\"}],\"name\":\"payments\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"podasTarget\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"roundSubmissions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"samplePeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sampleRange\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"startEpoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"endEpoch\",\"type\":\"uint64\"}],\"internalType\":\"structIDASample.SampleRange\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sampleRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sampleTask\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sampleHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"podasTarget\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"restSubmissions\",\"type\":\"uint64\"}],\"internalType\":\"structIDASample.SampleTask\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceFeeRateBps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_baseReward\",\"type\":\"uint256\"}],\"name\":\"setBaseReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_blobPrice\",\"type\":\"uint256\"}],\"name\":\"setBlobPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_epochWindowSize\",\"type\":\"uint64\"}],\"name\":\"setEpochWindowSize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_rewardRatio\",\"type\":\"uint64\"}],\"name\":\"setRewardRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_targetRoundSubmissions\",\"type\":\"uint64\"}],\"name\":\"setRoundSubmissions\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"samplePeriod_\",\"type\":\"uint64\"}],\"name\":\"setSamplePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"bps\",\"type\":\"uint256\"}],\"name\":\"setServiceFeeRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"treasury_\",\"type\":\"address\"}],\"name\":\"setTreasury\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_dataRoots\",\"type\":\"bytes32[]\"}],\"name\":\"submitOriginalData\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sampleSeed\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"quorumId\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"lineIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"sublineIndex\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"quality\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[3]\",\"name\":\"blobRoots\",\"type\":\"bytes32[3]\"},{\"internalType\":\"bytes32[]\",\"name\":\"proof\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structSampleResponse\",\"name\":\"rep\",\"type\":\"tuple\"}],\"name\":\"submitSamplingResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"erasureCommitment\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"quorumBitmap\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"aggPkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"signature\",\"type\":\"tuple\"}],\"internalType\":\"structIDAEntrance.CommitRootSubmission[]\",\"name\":\"_submissions\",\"type\":\"tuple[]\"}],\"name\":\"submitVerifiedCommitRoots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sync\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_times\",\"type\":\"uint256\"}],\"name\":\"syncFixedTimes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetRoundSubmissions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetRoundSubmissionsNext\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalBaseReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"treasury\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_quorumId\",\"type\":\"uint256\"}],\"name\":\"verifiedErasureCommitment\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"payee\",\"type\":\"address\"}],\"name\":\"withdrawPayments\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DAEntranceABI is the input ABI used to generate the binding from.
// Deprecated: Use DAEntranceMetaData.ABI instead.
var DAEntranceABI = DAEntranceMetaData.ABI

// DAEntrance is an auto generated Go binding around an Ethereum contract.
type DAEntrance struct {
	DAEntranceCaller     // Read-only binding to the contract
	DAEntranceTransactor // Write-only binding to the contract
	DAEntranceFilterer   // Log filterer for contract events
}

// DAEntranceCaller is an auto generated read-only Go binding around an Ethereum contract.
type DAEntranceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAEntranceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DAEntranceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAEntranceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DAEntranceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DAEntranceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DAEntranceSession struct {
	Contract     *DAEntrance       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DAEntranceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DAEntranceCallerSession struct {
	Contract *DAEntranceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// DAEntranceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DAEntranceTransactorSession struct {
	Contract     *DAEntranceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DAEntranceRaw is an auto generated low-level Go binding around an Ethereum contract.
type DAEntranceRaw struct {
	Contract *DAEntrance // Generic contract binding to access the raw methods on
}

// DAEntranceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DAEntranceCallerRaw struct {
	Contract *DAEntranceCaller // Generic read-only contract binding to access the raw methods on
}

// DAEntranceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DAEntranceTransactorRaw struct {
	Contract *DAEntranceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDAEntrance creates a new instance of DAEntrance, bound to a specific deployed contract.
func NewDAEntrance(address common.Address, backend bind.ContractBackend) (*DAEntrance, error) {
	contract, err := bindDAEntrance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DAEntrance{DAEntranceCaller: DAEntranceCaller{contract: contract}, DAEntranceTransactor: DAEntranceTransactor{contract: contract}, DAEntranceFilterer: DAEntranceFilterer{contract: contract}}, nil
}

// NewDAEntranceCaller creates a new read-only instance of DAEntrance, bound to a specific deployed contract.
func NewDAEntranceCaller(address common.Address, caller bind.ContractCaller) (*DAEntranceCaller, error) {
	contract, err := bindDAEntrance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DAEntranceCaller{contract: contract}, nil
}

// NewDAEntranceTransactor creates a new write-only instance of DAEntrance, bound to a specific deployed contract.
func NewDAEntranceTransactor(address common.Address, transactor bind.ContractTransactor) (*DAEntranceTransactor, error) {
	contract, err := bindDAEntrance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DAEntranceTransactor{contract: contract}, nil
}

// NewDAEntranceFilterer creates a new log filterer instance of DAEntrance, bound to a specific deployed contract.
func NewDAEntranceFilterer(address common.Address, filterer bind.ContractFilterer) (*DAEntranceFilterer, error) {
	contract, err := bindDAEntrance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DAEntranceFilterer{contract: contract}, nil
}

// bindDAEntrance binds a generic wrapper to an already deployed contract.
func bindDAEntrance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DAEntranceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DAEntrance *DAEntranceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DAEntrance.Contract.DAEntranceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DAEntrance *DAEntranceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.Contract.DAEntranceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DAEntrance *DAEntranceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DAEntrance.Contract.DAEntranceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DAEntrance *DAEntranceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DAEntrance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DAEntrance *DAEntranceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DAEntrance *DAEntranceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DAEntrance.Contract.contract.Transact(opts, method, params...)
}

// DASIGNERS is a free data retrieval call binding the contract method 0x807f063a.
//
// Solidity: function DA_SIGNERS() view returns(address)
func (_DAEntrance *DAEntranceCaller) DASIGNERS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "DA_SIGNERS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DASIGNERS is a free data retrieval call binding the contract method 0x807f063a.
//
// Solidity: function DA_SIGNERS() view returns(address)
func (_DAEntrance *DAEntranceSession) DASIGNERS() (common.Address, error) {
	return _DAEntrance.Contract.DASIGNERS(&_DAEntrance.CallOpts)
}

// DASIGNERS is a free data retrieval call binding the contract method 0x807f063a.
//
// Solidity: function DA_SIGNERS() view returns(address)
func (_DAEntrance *DAEntranceCallerSession) DASIGNERS() (common.Address, error) {
	return _DAEntrance.Contract.DASIGNERS(&_DAEntrance.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DAEntrance.Contract.DEFAULTADMINROLE(&_DAEntrance.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DAEntrance.Contract.DEFAULTADMINROLE(&_DAEntrance.CallOpts)
}

// MAXPODASTARGET is a free data retrieval call binding the contract method 0x8bdcc712.
//
// Solidity: function MAX_PODAS_TARGET() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) MAXPODASTARGET(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "MAX_PODAS_TARGET")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPODASTARGET is a free data retrieval call binding the contract method 0x8bdcc712.
//
// Solidity: function MAX_PODAS_TARGET() view returns(uint256)
func (_DAEntrance *DAEntranceSession) MAXPODASTARGET() (*big.Int, error) {
	return _DAEntrance.Contract.MAXPODASTARGET(&_DAEntrance.CallOpts)
}

// MAXPODASTARGET is a free data retrieval call binding the contract method 0x8bdcc712.
//
// Solidity: function MAX_PODAS_TARGET() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) MAXPODASTARGET() (*big.Int, error) {
	return _DAEntrance.Contract.MAXPODASTARGET(&_DAEntrance.CallOpts)
}

// PARAMSADMINROLE is a free data retrieval call binding the contract method 0xb15d20da.
//
// Solidity: function PARAMS_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceCaller) PARAMSADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "PARAMS_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PARAMSADMINROLE is a free data retrieval call binding the contract method 0xb15d20da.
//
// Solidity: function PARAMS_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceSession) PARAMSADMINROLE() ([32]byte, error) {
	return _DAEntrance.Contract.PARAMSADMINROLE(&_DAEntrance.CallOpts)
}

// PARAMSADMINROLE is a free data retrieval call binding the contract method 0xb15d20da.
//
// Solidity: function PARAMS_ADMIN_ROLE() view returns(bytes32)
func (_DAEntrance *DAEntranceCallerSession) PARAMSADMINROLE() ([32]byte, error) {
	return _DAEntrance.Contract.PARAMSADMINROLE(&_DAEntrance.CallOpts)
}

// SLICEDENOMINATOR is a free data retrieval call binding the contract method 0x602b245a.
//
// Solidity: function SLICE_DENOMINATOR() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) SLICEDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "SLICE_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SLICEDENOMINATOR is a free data retrieval call binding the contract method 0x602b245a.
//
// Solidity: function SLICE_DENOMINATOR() view returns(uint256)
func (_DAEntrance *DAEntranceSession) SLICEDENOMINATOR() (*big.Int, error) {
	return _DAEntrance.Contract.SLICEDENOMINATOR(&_DAEntrance.CallOpts)
}

// SLICEDENOMINATOR is a free data retrieval call binding the contract method 0x602b245a.
//
// Solidity: function SLICE_DENOMINATOR() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) SLICEDENOMINATOR() (*big.Int, error) {
	return _DAEntrance.Contract.SLICEDENOMINATOR(&_DAEntrance.CallOpts)
}

// SLICENUMERATOR is a free data retrieval call binding the contract method 0x5626a47b.
//
// Solidity: function SLICE_NUMERATOR() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) SLICENUMERATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "SLICE_NUMERATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SLICENUMERATOR is a free data retrieval call binding the contract method 0x5626a47b.
//
// Solidity: function SLICE_NUMERATOR() view returns(uint256)
func (_DAEntrance *DAEntranceSession) SLICENUMERATOR() (*big.Int, error) {
	return _DAEntrance.Contract.SLICENUMERATOR(&_DAEntrance.CallOpts)
}

// SLICENUMERATOR is a free data retrieval call binding the contract method 0x5626a47b.
//
// Solidity: function SLICE_NUMERATOR() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) SLICENUMERATOR() (*big.Int, error) {
	return _DAEntrance.Contract.SLICENUMERATOR(&_DAEntrance.CallOpts)
}

// ActivedReward is a free data retrieval call binding the contract method 0x12577052.
//
// Solidity: function activedReward() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) ActivedReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "activedReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActivedReward is a free data retrieval call binding the contract method 0x12577052.
//
// Solidity: function activedReward() view returns(uint256)
func (_DAEntrance *DAEntranceSession) ActivedReward() (*big.Int, error) {
	return _DAEntrance.Contract.ActivedReward(&_DAEntrance.CallOpts)
}

// ActivedReward is a free data retrieval call binding the contract method 0x12577052.
//
// Solidity: function activedReward() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) ActivedReward() (*big.Int, error) {
	return _DAEntrance.Contract.ActivedReward(&_DAEntrance.CallOpts)
}

// BaseReward is a free data retrieval call binding the contract method 0x76ad03bc.
//
// Solidity: function baseReward() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) BaseReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "baseReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseReward is a free data retrieval call binding the contract method 0x76ad03bc.
//
// Solidity: function baseReward() view returns(uint256)
func (_DAEntrance *DAEntranceSession) BaseReward() (*big.Int, error) {
	return _DAEntrance.Contract.BaseReward(&_DAEntrance.CallOpts)
}

// BaseReward is a free data retrieval call binding the contract method 0x76ad03bc.
//
// Solidity: function baseReward() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) BaseReward() (*big.Int, error) {
	return _DAEntrance.Contract.BaseReward(&_DAEntrance.CallOpts)
}

// BlobPrice is a free data retrieval call binding the contract method 0x3d00448a.
//
// Solidity: function blobPrice() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) BlobPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "blobPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlobPrice is a free data retrieval call binding the contract method 0x3d00448a.
//
// Solidity: function blobPrice() view returns(uint256)
func (_DAEntrance *DAEntranceSession) BlobPrice() (*big.Int, error) {
	return _DAEntrance.Contract.BlobPrice(&_DAEntrance.CallOpts)
}

// BlobPrice is a free data retrieval call binding the contract method 0x3d00448a.
//
// Solidity: function blobPrice() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) BlobPrice() (*big.Int, error) {
	return _DAEntrance.Contract.BlobPrice(&_DAEntrance.CallOpts)
}

// CommitmentExists is a free data retrieval call binding the contract method 0x6a53525d.
//
// Solidity: function commitmentExists(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns(bool)
func (_DAEntrance *DAEntranceCaller) CommitmentExists(opts *bind.CallOpts, _dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (bool, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "commitmentExists", _dataRoot, _epoch, _quorumId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CommitmentExists is a free data retrieval call binding the contract method 0x6a53525d.
//
// Solidity: function commitmentExists(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns(bool)
func (_DAEntrance *DAEntranceSession) CommitmentExists(_dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (bool, error) {
	return _DAEntrance.Contract.CommitmentExists(&_DAEntrance.CallOpts, _dataRoot, _epoch, _quorumId)
}

// CommitmentExists is a free data retrieval call binding the contract method 0x6a53525d.
//
// Solidity: function commitmentExists(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns(bool)
func (_DAEntrance *DAEntranceCallerSession) CommitmentExists(_dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (bool, error) {
	return _DAEntrance.Contract.CommitmentExists(&_DAEntrance.CallOpts, _dataRoot, _epoch, _quorumId)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_DAEntrance *DAEntranceSession) CurrentEpoch() (*big.Int, error) {
	return _DAEntrance.Contract.CurrentEpoch(&_DAEntrance.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) CurrentEpoch() (*big.Int, error) {
	return _DAEntrance.Contract.CurrentEpoch(&_DAEntrance.CallOpts)
}

// CurrentEpochReward is a free data retrieval call binding the contract method 0xe68e035b.
//
// Solidity: function currentEpochReward() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) CurrentEpochReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "currentEpochReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpochReward is a free data retrieval call binding the contract method 0xe68e035b.
//
// Solidity: function currentEpochReward() view returns(uint256)
func (_DAEntrance *DAEntranceSession) CurrentEpochReward() (*big.Int, error) {
	return _DAEntrance.Contract.CurrentEpochReward(&_DAEntrance.CallOpts)
}

// CurrentEpochReward is a free data retrieval call binding the contract method 0xe68e035b.
//
// Solidity: function currentEpochReward() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) CurrentEpochReward() (*big.Int, error) {
	return _DAEntrance.Contract.CurrentEpochReward(&_DAEntrance.CallOpts)
}

// CurrentSampleSeed is a free data retrieval call binding the contract method 0x9fae1fa4.
//
// Solidity: function currentSampleSeed() view returns(bytes32)
func (_DAEntrance *DAEntranceCaller) CurrentSampleSeed(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "currentSampleSeed")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CurrentSampleSeed is a free data retrieval call binding the contract method 0x9fae1fa4.
//
// Solidity: function currentSampleSeed() view returns(bytes32)
func (_DAEntrance *DAEntranceSession) CurrentSampleSeed() ([32]byte, error) {
	return _DAEntrance.Contract.CurrentSampleSeed(&_DAEntrance.CallOpts)
}

// CurrentSampleSeed is a free data retrieval call binding the contract method 0x9fae1fa4.
//
// Solidity: function currentSampleSeed() view returns(bytes32)
func (_DAEntrance *DAEntranceCallerSession) CurrentSampleSeed() ([32]byte, error) {
	return _DAEntrance.Contract.CurrentSampleSeed(&_DAEntrance.CallOpts)
}

// EpochWindowSize is a free data retrieval call binding the contract method 0x3e898337.
//
// Solidity: function epochWindowSize() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) EpochWindowSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "epochWindowSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochWindowSize is a free data retrieval call binding the contract method 0x3e898337.
//
// Solidity: function epochWindowSize() view returns(uint256)
func (_DAEntrance *DAEntranceSession) EpochWindowSize() (*big.Int, error) {
	return _DAEntrance.Contract.EpochWindowSize(&_DAEntrance.CallOpts)
}

// EpochWindowSize is a free data retrieval call binding the contract method 0x3e898337.
//
// Solidity: function epochWindowSize() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) EpochWindowSize() (*big.Int, error) {
	return _DAEntrance.Contract.EpochWindowSize(&_DAEntrance.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DAEntrance *DAEntranceCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DAEntrance *DAEntranceSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DAEntrance.Contract.GetRoleAdmin(&_DAEntrance.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DAEntrance *DAEntranceCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DAEntrance.Contract.GetRoleAdmin(&_DAEntrance.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_DAEntrance *DAEntranceCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_DAEntrance *DAEntranceSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _DAEntrance.Contract.GetRoleMember(&_DAEntrance.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_DAEntrance *DAEntranceCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _DAEntrance.Contract.GetRoleMember(&_DAEntrance.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_DAEntrance *DAEntranceCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_DAEntrance *DAEntranceSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _DAEntrance.Contract.GetRoleMemberCount(&_DAEntrance.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _DAEntrance.Contract.GetRoleMemberCount(&_DAEntrance.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DAEntrance *DAEntranceCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DAEntrance *DAEntranceSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DAEntrance.Contract.HasRole(&_DAEntrance.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DAEntrance *DAEntranceCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DAEntrance.Contract.HasRole(&_DAEntrance.CallOpts, role, account)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_DAEntrance *DAEntranceCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_DAEntrance *DAEntranceSession) Initialized() (bool, error) {
	return _DAEntrance.Contract.Initialized(&_DAEntrance.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_DAEntrance *DAEntranceCallerSession) Initialized() (bool, error) {
	return _DAEntrance.Contract.Initialized(&_DAEntrance.CallOpts)
}

// NextSampleHeight is a free data retrieval call binding the contract method 0x7397eb33.
//
// Solidity: function nextSampleHeight() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) NextSampleHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "nextSampleHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextSampleHeight is a free data retrieval call binding the contract method 0x7397eb33.
//
// Solidity: function nextSampleHeight() view returns(uint256)
func (_DAEntrance *DAEntranceSession) NextSampleHeight() (*big.Int, error) {
	return _DAEntrance.Contract.NextSampleHeight(&_DAEntrance.CallOpts)
}

// NextSampleHeight is a free data retrieval call binding the contract method 0x7397eb33.
//
// Solidity: function nextSampleHeight() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) NextSampleHeight() (*big.Int, error) {
	return _DAEntrance.Contract.NextSampleHeight(&_DAEntrance.CallOpts)
}

// Payments is a free data retrieval call binding the contract method 0xe2982c21.
//
// Solidity: function payments(address dest) view returns(uint256)
func (_DAEntrance *DAEntranceCaller) Payments(opts *bind.CallOpts, dest common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "payments", dest)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Payments is a free data retrieval call binding the contract method 0xe2982c21.
//
// Solidity: function payments(address dest) view returns(uint256)
func (_DAEntrance *DAEntranceSession) Payments(dest common.Address) (*big.Int, error) {
	return _DAEntrance.Contract.Payments(&_DAEntrance.CallOpts, dest)
}

// Payments is a free data retrieval call binding the contract method 0xe2982c21.
//
// Solidity: function payments(address dest) view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) Payments(dest common.Address) (*big.Int, error) {
	return _DAEntrance.Contract.Payments(&_DAEntrance.CallOpts, dest)
}

// PodasTarget is a free data retrieval call binding the contract method 0xc8d3b359.
//
// Solidity: function podasTarget() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) PodasTarget(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "podasTarget")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PodasTarget is a free data retrieval call binding the contract method 0xc8d3b359.
//
// Solidity: function podasTarget() view returns(uint256)
func (_DAEntrance *DAEntranceSession) PodasTarget() (*big.Int, error) {
	return _DAEntrance.Contract.PodasTarget(&_DAEntrance.CallOpts)
}

// PodasTarget is a free data retrieval call binding the contract method 0xc8d3b359.
//
// Solidity: function podasTarget() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) PodasTarget() (*big.Int, error) {
	return _DAEntrance.Contract.PodasTarget(&_DAEntrance.CallOpts)
}

// RewardRatio is a free data retrieval call binding the contract method 0x646033bc.
//
// Solidity: function rewardRatio() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) RewardRatio(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "rewardRatio")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardRatio is a free data retrieval call binding the contract method 0x646033bc.
//
// Solidity: function rewardRatio() view returns(uint256)
func (_DAEntrance *DAEntranceSession) RewardRatio() (*big.Int, error) {
	return _DAEntrance.Contract.RewardRatio(&_DAEntrance.CallOpts)
}

// RewardRatio is a free data retrieval call binding the contract method 0x646033bc.
//
// Solidity: function rewardRatio() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) RewardRatio() (*big.Int, error) {
	return _DAEntrance.Contract.RewardRatio(&_DAEntrance.CallOpts)
}

// RoundSubmissions is a free data retrieval call binding the contract method 0xff187748.
//
// Solidity: function roundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) RoundSubmissions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "roundSubmissions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RoundSubmissions is a free data retrieval call binding the contract method 0xff187748.
//
// Solidity: function roundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceSession) RoundSubmissions() (*big.Int, error) {
	return _DAEntrance.Contract.RoundSubmissions(&_DAEntrance.CallOpts)
}

// RoundSubmissions is a free data retrieval call binding the contract method 0xff187748.
//
// Solidity: function roundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) RoundSubmissions() (*big.Int, error) {
	return _DAEntrance.Contract.RoundSubmissions(&_DAEntrance.CallOpts)
}

// SamplePeriod is a free data retrieval call binding the contract method 0x98920f57.
//
// Solidity: function samplePeriod() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) SamplePeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "samplePeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SamplePeriod is a free data retrieval call binding the contract method 0x98920f57.
//
// Solidity: function samplePeriod() view returns(uint256)
func (_DAEntrance *DAEntranceSession) SamplePeriod() (*big.Int, error) {
	return _DAEntrance.Contract.SamplePeriod(&_DAEntrance.CallOpts)
}

// SamplePeriod is a free data retrieval call binding the contract method 0x98920f57.
//
// Solidity: function samplePeriod() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) SamplePeriod() (*big.Int, error) {
	return _DAEntrance.Contract.SamplePeriod(&_DAEntrance.CallOpts)
}

// SampleRound is a free data retrieval call binding the contract method 0x168a062c.
//
// Solidity: function sampleRound() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) SampleRound(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "sampleRound")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SampleRound is a free data retrieval call binding the contract method 0x168a062c.
//
// Solidity: function sampleRound() view returns(uint256)
func (_DAEntrance *DAEntranceSession) SampleRound() (*big.Int, error) {
	return _DAEntrance.Contract.SampleRound(&_DAEntrance.CallOpts)
}

// SampleRound is a free data retrieval call binding the contract method 0x168a062c.
//
// Solidity: function sampleRound() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) SampleRound() (*big.Int, error) {
	return _DAEntrance.Contract.SampleRound(&_DAEntrance.CallOpts)
}

// ServiceFee is a free data retrieval call binding the contract method 0x8abdf5aa.
//
// Solidity: function serviceFee() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) ServiceFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "serviceFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ServiceFee is a free data retrieval call binding the contract method 0x8abdf5aa.
//
// Solidity: function serviceFee() view returns(uint256)
func (_DAEntrance *DAEntranceSession) ServiceFee() (*big.Int, error) {
	return _DAEntrance.Contract.ServiceFee(&_DAEntrance.CallOpts)
}

// ServiceFee is a free data retrieval call binding the contract method 0x8abdf5aa.
//
// Solidity: function serviceFee() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) ServiceFee() (*big.Int, error) {
	return _DAEntrance.Contract.ServiceFee(&_DAEntrance.CallOpts)
}

// ServiceFeeRateBps is a free data retrieval call binding the contract method 0xc0575111.
//
// Solidity: function serviceFeeRateBps() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) ServiceFeeRateBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "serviceFeeRateBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ServiceFeeRateBps is a free data retrieval call binding the contract method 0xc0575111.
//
// Solidity: function serviceFeeRateBps() view returns(uint256)
func (_DAEntrance *DAEntranceSession) ServiceFeeRateBps() (*big.Int, error) {
	return _DAEntrance.Contract.ServiceFeeRateBps(&_DAEntrance.CallOpts)
}

// ServiceFeeRateBps is a free data retrieval call binding the contract method 0xc0575111.
//
// Solidity: function serviceFeeRateBps() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) ServiceFeeRateBps() (*big.Int, error) {
	return _DAEntrance.Contract.ServiceFeeRateBps(&_DAEntrance.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DAEntrance *DAEntranceCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DAEntrance *DAEntranceSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DAEntrance.Contract.SupportsInterface(&_DAEntrance.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DAEntrance *DAEntranceCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DAEntrance.Contract.SupportsInterface(&_DAEntrance.CallOpts, interfaceId)
}

// TargetRoundSubmissions is a free data retrieval call binding the contract method 0x2fc0534b.
//
// Solidity: function targetRoundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) TargetRoundSubmissions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "targetRoundSubmissions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetRoundSubmissions is a free data retrieval call binding the contract method 0x2fc0534b.
//
// Solidity: function targetRoundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceSession) TargetRoundSubmissions() (*big.Int, error) {
	return _DAEntrance.Contract.TargetRoundSubmissions(&_DAEntrance.CallOpts)
}

// TargetRoundSubmissions is a free data retrieval call binding the contract method 0x2fc0534b.
//
// Solidity: function targetRoundSubmissions() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) TargetRoundSubmissions() (*big.Int, error) {
	return _DAEntrance.Contract.TargetRoundSubmissions(&_DAEntrance.CallOpts)
}

// TargetRoundSubmissionsNext is a free data retrieval call binding the contract method 0x257a3aa3.
//
// Solidity: function targetRoundSubmissionsNext() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) TargetRoundSubmissionsNext(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "targetRoundSubmissionsNext")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetRoundSubmissionsNext is a free data retrieval call binding the contract method 0x257a3aa3.
//
// Solidity: function targetRoundSubmissionsNext() view returns(uint256)
func (_DAEntrance *DAEntranceSession) TargetRoundSubmissionsNext() (*big.Int, error) {
	return _DAEntrance.Contract.TargetRoundSubmissionsNext(&_DAEntrance.CallOpts)
}

// TargetRoundSubmissionsNext is a free data retrieval call binding the contract method 0x257a3aa3.
//
// Solidity: function targetRoundSubmissionsNext() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) TargetRoundSubmissionsNext() (*big.Int, error) {
	return _DAEntrance.Contract.TargetRoundSubmissionsNext(&_DAEntrance.CallOpts)
}

// TotalBaseReward is a free data retrieval call binding the contract method 0x7f1b5e43.
//
// Solidity: function totalBaseReward() view returns(uint256)
func (_DAEntrance *DAEntranceCaller) TotalBaseReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "totalBaseReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBaseReward is a free data retrieval call binding the contract method 0x7f1b5e43.
//
// Solidity: function totalBaseReward() view returns(uint256)
func (_DAEntrance *DAEntranceSession) TotalBaseReward() (*big.Int, error) {
	return _DAEntrance.Contract.TotalBaseReward(&_DAEntrance.CallOpts)
}

// TotalBaseReward is a free data retrieval call binding the contract method 0x7f1b5e43.
//
// Solidity: function totalBaseReward() view returns(uint256)
func (_DAEntrance *DAEntranceCallerSession) TotalBaseReward() (*big.Int, error) {
	return _DAEntrance.Contract.TotalBaseReward(&_DAEntrance.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_DAEntrance *DAEntranceCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_DAEntrance *DAEntranceSession) Treasury() (common.Address, error) {
	return _DAEntrance.Contract.Treasury(&_DAEntrance.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_DAEntrance *DAEntranceCallerSession) Treasury() (common.Address, error) {
	return _DAEntrance.Contract.Treasury(&_DAEntrance.CallOpts)
}

// VerifiedErasureCommitment is a free data retrieval call binding the contract method 0x9da3a69b.
//
// Solidity: function verifiedErasureCommitment(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns((uint256,uint256))
func (_DAEntrance *DAEntranceCaller) VerifiedErasureCommitment(opts *bind.CallOpts, _dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (BN254G1Point, error) {
	var out []interface{}
	err := _DAEntrance.contract.Call(opts, &out, "verifiedErasureCommitment", _dataRoot, _epoch, _quorumId)

	if err != nil {
		return *new(BN254G1Point), err
	}

	out0 := *abi.ConvertType(out[0], new(BN254G1Point)).(*BN254G1Point)

	return out0, err

}

// VerifiedErasureCommitment is a free data retrieval call binding the contract method 0x9da3a69b.
//
// Solidity: function verifiedErasureCommitment(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns((uint256,uint256))
func (_DAEntrance *DAEntranceSession) VerifiedErasureCommitment(_dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (BN254G1Point, error) {
	return _DAEntrance.Contract.VerifiedErasureCommitment(&_DAEntrance.CallOpts, _dataRoot, _epoch, _quorumId)
}

// VerifiedErasureCommitment is a free data retrieval call binding the contract method 0x9da3a69b.
//
// Solidity: function verifiedErasureCommitment(bytes32 _dataRoot, uint256 _epoch, uint256 _quorumId) view returns((uint256,uint256))
func (_DAEntrance *DAEntranceCallerSession) VerifiedErasureCommitment(_dataRoot [32]byte, _epoch *big.Int, _quorumId *big.Int) (BN254G1Point, error) {
	return _DAEntrance.Contract.VerifiedErasureCommitment(&_DAEntrance.CallOpts, _dataRoot, _epoch, _quorumId)
}

// Donate is a paid mutator transaction binding the contract method 0xed88c68e.
//
// Solidity: function donate() payable returns()
func (_DAEntrance *DAEntranceTransactor) Donate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "donate")
}

// Donate is a paid mutator transaction binding the contract method 0xed88c68e.
//
// Solidity: function donate() payable returns()
func (_DAEntrance *DAEntranceSession) Donate() (*types.Transaction, error) {
	return _DAEntrance.Contract.Donate(&_DAEntrance.TransactOpts)
}

// Donate is a paid mutator transaction binding the contract method 0xed88c68e.
//
// Solidity: function donate() payable returns()
func (_DAEntrance *DAEntranceTransactorSession) Donate() (*types.Transaction, error) {
	return _DAEntrance.Contract.Donate(&_DAEntrance.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.GrantRole(&_DAEntrance.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.GrantRole(&_DAEntrance.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DAEntrance *DAEntranceTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DAEntrance *DAEntranceSession) Initialize() (*types.Transaction, error) {
	return _DAEntrance.Contract.Initialize(&_DAEntrance.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DAEntrance *DAEntranceTransactorSession) Initialize() (*types.Transaction, error) {
	return _DAEntrance.Contract.Initialize(&_DAEntrance.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.RenounceRole(&_DAEntrance.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.RenounceRole(&_DAEntrance.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.RevokeRole(&_DAEntrance.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DAEntrance *DAEntranceTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.RevokeRole(&_DAEntrance.TransactOpts, role, account)
}

// SampleRange is a paid mutator transaction binding the contract method 0x6efc2555.
//
// Solidity: function sampleRange() returns((uint64,uint64))
func (_DAEntrance *DAEntranceTransactor) SampleRange(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "sampleRange")
}

// SampleRange is a paid mutator transaction binding the contract method 0x6efc2555.
//
// Solidity: function sampleRange() returns((uint64,uint64))
func (_DAEntrance *DAEntranceSession) SampleRange() (*types.Transaction, error) {
	return _DAEntrance.Contract.SampleRange(&_DAEntrance.TransactOpts)
}

// SampleRange is a paid mutator transaction binding the contract method 0x6efc2555.
//
// Solidity: function sampleRange() returns((uint64,uint64))
func (_DAEntrance *DAEntranceTransactorSession) SampleRange() (*types.Transaction, error) {
	return _DAEntrance.Contract.SampleRange(&_DAEntrance.TransactOpts)
}

// SampleTask is a paid mutator transaction binding the contract method 0x988ea94e.
//
// Solidity: function sampleTask() returns((bytes32,uint256,uint64))
func (_DAEntrance *DAEntranceTransactor) SampleTask(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "sampleTask")
}

// SampleTask is a paid mutator transaction binding the contract method 0x988ea94e.
//
// Solidity: function sampleTask() returns((bytes32,uint256,uint64))
func (_DAEntrance *DAEntranceSession) SampleTask() (*types.Transaction, error) {
	return _DAEntrance.Contract.SampleTask(&_DAEntrance.TransactOpts)
}

// SampleTask is a paid mutator transaction binding the contract method 0x988ea94e.
//
// Solidity: function sampleTask() returns((bytes32,uint256,uint64))
func (_DAEntrance *DAEntranceTransactorSession) SampleTask() (*types.Transaction, error) {
	return _DAEntrance.Contract.SampleTask(&_DAEntrance.TransactOpts)
}

// SetBaseReward is a paid mutator transaction binding the contract method 0x0373a23a.
//
// Solidity: function setBaseReward(uint256 _baseReward) returns()
func (_DAEntrance *DAEntranceTransactor) SetBaseReward(opts *bind.TransactOpts, _baseReward *big.Int) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setBaseReward", _baseReward)
}

// SetBaseReward is a paid mutator transaction binding the contract method 0x0373a23a.
//
// Solidity: function setBaseReward(uint256 _baseReward) returns()
func (_DAEntrance *DAEntranceSession) SetBaseReward(_baseReward *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetBaseReward(&_DAEntrance.TransactOpts, _baseReward)
}

// SetBaseReward is a paid mutator transaction binding the contract method 0x0373a23a.
//
// Solidity: function setBaseReward(uint256 _baseReward) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetBaseReward(_baseReward *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetBaseReward(&_DAEntrance.TransactOpts, _baseReward)
}

// SetBlobPrice is a paid mutator transaction binding the contract method 0x23dd60a6.
//
// Solidity: function setBlobPrice(uint256 _blobPrice) returns()
func (_DAEntrance *DAEntranceTransactor) SetBlobPrice(opts *bind.TransactOpts, _blobPrice *big.Int) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setBlobPrice", _blobPrice)
}

// SetBlobPrice is a paid mutator transaction binding the contract method 0x23dd60a6.
//
// Solidity: function setBlobPrice(uint256 _blobPrice) returns()
func (_DAEntrance *DAEntranceSession) SetBlobPrice(_blobPrice *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetBlobPrice(&_DAEntrance.TransactOpts, _blobPrice)
}

// SetBlobPrice is a paid mutator transaction binding the contract method 0x23dd60a6.
//
// Solidity: function setBlobPrice(uint256 _blobPrice) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetBlobPrice(_blobPrice *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetBlobPrice(&_DAEntrance.TransactOpts, _blobPrice)
}

// SetEpochWindowSize is a paid mutator transaction binding the contract method 0xb1be17ab.
//
// Solidity: function setEpochWindowSize(uint64 _epochWindowSize) returns()
func (_DAEntrance *DAEntranceTransactor) SetEpochWindowSize(opts *bind.TransactOpts, _epochWindowSize uint64) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setEpochWindowSize", _epochWindowSize)
}

// SetEpochWindowSize is a paid mutator transaction binding the contract method 0xb1be17ab.
//
// Solidity: function setEpochWindowSize(uint64 _epochWindowSize) returns()
func (_DAEntrance *DAEntranceSession) SetEpochWindowSize(_epochWindowSize uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetEpochWindowSize(&_DAEntrance.TransactOpts, _epochWindowSize)
}

// SetEpochWindowSize is a paid mutator transaction binding the contract method 0xb1be17ab.
//
// Solidity: function setEpochWindowSize(uint64 _epochWindowSize) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetEpochWindowSize(_epochWindowSize uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetEpochWindowSize(&_DAEntrance.TransactOpts, _epochWindowSize)
}

// SetRewardRatio is a paid mutator transaction binding the contract method 0x3bab2a70.
//
// Solidity: function setRewardRatio(uint64 _rewardRatio) returns()
func (_DAEntrance *DAEntranceTransactor) SetRewardRatio(opts *bind.TransactOpts, _rewardRatio uint64) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setRewardRatio", _rewardRatio)
}

// SetRewardRatio is a paid mutator transaction binding the contract method 0x3bab2a70.
//
// Solidity: function setRewardRatio(uint64 _rewardRatio) returns()
func (_DAEntrance *DAEntranceSession) SetRewardRatio(_rewardRatio uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetRewardRatio(&_DAEntrance.TransactOpts, _rewardRatio)
}

// SetRewardRatio is a paid mutator transaction binding the contract method 0x3bab2a70.
//
// Solidity: function setRewardRatio(uint64 _rewardRatio) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetRewardRatio(_rewardRatio uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetRewardRatio(&_DAEntrance.TransactOpts, _rewardRatio)
}

// SetRoundSubmissions is a paid mutator transaction binding the contract method 0x88521ec7.
//
// Solidity: function setRoundSubmissions(uint64 _targetRoundSubmissions) returns()
func (_DAEntrance *DAEntranceTransactor) SetRoundSubmissions(opts *bind.TransactOpts, _targetRoundSubmissions uint64) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setRoundSubmissions", _targetRoundSubmissions)
}

// SetRoundSubmissions is a paid mutator transaction binding the contract method 0x88521ec7.
//
// Solidity: function setRoundSubmissions(uint64 _targetRoundSubmissions) returns()
func (_DAEntrance *DAEntranceSession) SetRoundSubmissions(_targetRoundSubmissions uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetRoundSubmissions(&_DAEntrance.TransactOpts, _targetRoundSubmissions)
}

// SetRoundSubmissions is a paid mutator transaction binding the contract method 0x88521ec7.
//
// Solidity: function setRoundSubmissions(uint64 _targetRoundSubmissions) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetRoundSubmissions(_targetRoundSubmissions uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetRoundSubmissions(&_DAEntrance.TransactOpts, _targetRoundSubmissions)
}

// SetSamplePeriod is a paid mutator transaction binding the contract method 0x1192de9a.
//
// Solidity: function setSamplePeriod(uint64 samplePeriod_) returns()
func (_DAEntrance *DAEntranceTransactor) SetSamplePeriod(opts *bind.TransactOpts, samplePeriod_ uint64) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setSamplePeriod", samplePeriod_)
}

// SetSamplePeriod is a paid mutator transaction binding the contract method 0x1192de9a.
//
// Solidity: function setSamplePeriod(uint64 samplePeriod_) returns()
func (_DAEntrance *DAEntranceSession) SetSamplePeriod(samplePeriod_ uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetSamplePeriod(&_DAEntrance.TransactOpts, samplePeriod_)
}

// SetSamplePeriod is a paid mutator transaction binding the contract method 0x1192de9a.
//
// Solidity: function setSamplePeriod(uint64 samplePeriod_) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetSamplePeriod(samplePeriod_ uint64) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetSamplePeriod(&_DAEntrance.TransactOpts, samplePeriod_)
}

// SetServiceFeeRate is a paid mutator transaction binding the contract method 0x9b1d3091.
//
// Solidity: function setServiceFeeRate(uint256 bps) returns()
func (_DAEntrance *DAEntranceTransactor) SetServiceFeeRate(opts *bind.TransactOpts, bps *big.Int) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setServiceFeeRate", bps)
}

// SetServiceFeeRate is a paid mutator transaction binding the contract method 0x9b1d3091.
//
// Solidity: function setServiceFeeRate(uint256 bps) returns()
func (_DAEntrance *DAEntranceSession) SetServiceFeeRate(bps *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetServiceFeeRate(&_DAEntrance.TransactOpts, bps)
}

// SetServiceFeeRate is a paid mutator transaction binding the contract method 0x9b1d3091.
//
// Solidity: function setServiceFeeRate(uint256 bps) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetServiceFeeRate(bps *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetServiceFeeRate(&_DAEntrance.TransactOpts, bps)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury_) returns()
func (_DAEntrance *DAEntranceTransactor) SetTreasury(opts *bind.TransactOpts, treasury_ common.Address) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "setTreasury", treasury_)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury_) returns()
func (_DAEntrance *DAEntranceSession) SetTreasury(treasury_ common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetTreasury(&_DAEntrance.TransactOpts, treasury_)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury_) returns()
func (_DAEntrance *DAEntranceTransactorSession) SetTreasury(treasury_ common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.SetTreasury(&_DAEntrance.TransactOpts, treasury_)
}

// SubmitOriginalData is a paid mutator transaction binding the contract method 0xd4ae59c9.
//
// Solidity: function submitOriginalData(bytes32[] _dataRoots) payable returns()
func (_DAEntrance *DAEntranceTransactor) SubmitOriginalData(opts *bind.TransactOpts, _dataRoots [][32]byte) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "submitOriginalData", _dataRoots)
}

// SubmitOriginalData is a paid mutator transaction binding the contract method 0xd4ae59c9.
//
// Solidity: function submitOriginalData(bytes32[] _dataRoots) payable returns()
func (_DAEntrance *DAEntranceSession) SubmitOriginalData(_dataRoots [][32]byte) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitOriginalData(&_DAEntrance.TransactOpts, _dataRoots)
}

// SubmitOriginalData is a paid mutator transaction binding the contract method 0xd4ae59c9.
//
// Solidity: function submitOriginalData(bytes32[] _dataRoots) payable returns()
func (_DAEntrance *DAEntranceTransactorSession) SubmitOriginalData(_dataRoots [][32]byte) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitOriginalData(&_DAEntrance.TransactOpts, _dataRoots)
}

// SubmitSamplingResponse is a paid mutator transaction binding the contract method 0xf6902775.
//
// Solidity: function submitSamplingResponse((bytes32,uint64,uint64,uint32,uint32,uint256,bytes32,bytes32[3],bytes32[],bytes) rep) returns()
func (_DAEntrance *DAEntranceTransactor) SubmitSamplingResponse(opts *bind.TransactOpts, rep SampleResponse) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "submitSamplingResponse", rep)
}

// SubmitSamplingResponse is a paid mutator transaction binding the contract method 0xf6902775.
//
// Solidity: function submitSamplingResponse((bytes32,uint64,uint64,uint32,uint32,uint256,bytes32,bytes32[3],bytes32[],bytes) rep) returns()
func (_DAEntrance *DAEntranceSession) SubmitSamplingResponse(rep SampleResponse) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitSamplingResponse(&_DAEntrance.TransactOpts, rep)
}

// SubmitSamplingResponse is a paid mutator transaction binding the contract method 0xf6902775.
//
// Solidity: function submitSamplingResponse((bytes32,uint64,uint64,uint32,uint32,uint256,bytes32,bytes32[3],bytes32[],bytes) rep) returns()
func (_DAEntrance *DAEntranceTransactorSession) SubmitSamplingResponse(rep SampleResponse) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitSamplingResponse(&_DAEntrance.TransactOpts, rep)
}

// SubmitVerifiedCommitRoots is a paid mutator transaction binding the contract method 0xeafed6ce.
//
// Solidity: function submitVerifiedCommitRoots((bytes32,uint256,uint256,(uint256,uint256),bytes,(uint256[2],uint256[2]),(uint256,uint256))[] _submissions) returns()
func (_DAEntrance *DAEntranceTransactor) SubmitVerifiedCommitRoots(opts *bind.TransactOpts, _submissions []IDAEntranceCommitRootSubmission) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "submitVerifiedCommitRoots", _submissions)
}

// SubmitVerifiedCommitRoots is a paid mutator transaction binding the contract method 0xeafed6ce.
//
// Solidity: function submitVerifiedCommitRoots((bytes32,uint256,uint256,(uint256,uint256),bytes,(uint256[2],uint256[2]),(uint256,uint256))[] _submissions) returns()
func (_DAEntrance *DAEntranceSession) SubmitVerifiedCommitRoots(_submissions []IDAEntranceCommitRootSubmission) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitVerifiedCommitRoots(&_DAEntrance.TransactOpts, _submissions)
}

// SubmitVerifiedCommitRoots is a paid mutator transaction binding the contract method 0xeafed6ce.
//
// Solidity: function submitVerifiedCommitRoots((bytes32,uint256,uint256,(uint256,uint256),bytes,(uint256[2],uint256[2]),(uint256,uint256))[] _submissions) returns()
func (_DAEntrance *DAEntranceTransactorSession) SubmitVerifiedCommitRoots(_submissions []IDAEntranceCommitRootSubmission) (*types.Transaction, error) {
	return _DAEntrance.Contract.SubmitVerifiedCommitRoots(&_DAEntrance.TransactOpts, _submissions)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_DAEntrance *DAEntranceTransactor) Sync(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "sync")
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_DAEntrance *DAEntranceSession) Sync() (*types.Transaction, error) {
	return _DAEntrance.Contract.Sync(&_DAEntrance.TransactOpts)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_DAEntrance *DAEntranceTransactorSession) Sync() (*types.Transaction, error) {
	return _DAEntrance.Contract.Sync(&_DAEntrance.TransactOpts)
}

// SyncFixedTimes is a paid mutator transaction binding the contract method 0xe051d1ea.
//
// Solidity: function syncFixedTimes(uint256 _times) returns()
func (_DAEntrance *DAEntranceTransactor) SyncFixedTimes(opts *bind.TransactOpts, _times *big.Int) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "syncFixedTimes", _times)
}

// SyncFixedTimes is a paid mutator transaction binding the contract method 0xe051d1ea.
//
// Solidity: function syncFixedTimes(uint256 _times) returns()
func (_DAEntrance *DAEntranceSession) SyncFixedTimes(_times *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SyncFixedTimes(&_DAEntrance.TransactOpts, _times)
}

// SyncFixedTimes is a paid mutator transaction binding the contract method 0xe051d1ea.
//
// Solidity: function syncFixedTimes(uint256 _times) returns()
func (_DAEntrance *DAEntranceTransactorSession) SyncFixedTimes(_times *big.Int) (*types.Transaction, error) {
	return _DAEntrance.Contract.SyncFixedTimes(&_DAEntrance.TransactOpts, _times)
}

// WithdrawPayments is a paid mutator transaction binding the contract method 0x31b3eb94.
//
// Solidity: function withdrawPayments(address payee) returns()
func (_DAEntrance *DAEntranceTransactor) WithdrawPayments(opts *bind.TransactOpts, payee common.Address) (*types.Transaction, error) {
	return _DAEntrance.contract.Transact(opts, "withdrawPayments", payee)
}

// WithdrawPayments is a paid mutator transaction binding the contract method 0x31b3eb94.
//
// Solidity: function withdrawPayments(address payee) returns()
func (_DAEntrance *DAEntranceSession) WithdrawPayments(payee common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.WithdrawPayments(&_DAEntrance.TransactOpts, payee)
}

// WithdrawPayments is a paid mutator transaction binding the contract method 0x31b3eb94.
//
// Solidity: function withdrawPayments(address payee) returns()
func (_DAEntrance *DAEntranceTransactorSession) WithdrawPayments(payee common.Address) (*types.Transaction, error) {
	return _DAEntrance.Contract.WithdrawPayments(&_DAEntrance.TransactOpts, payee)
}

// DAEntranceDARewardIterator is returned from FilterDAReward and is used to iterate over the raw logs and unpacked data for DAReward events raised by the DAEntrance contract.
type DAEntranceDARewardIterator struct {
	Event *DAEntranceDAReward // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceDARewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceDAReward)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceDAReward)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceDARewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceDARewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceDAReward represents a DAReward event raised by the DAEntrance contract.
type DAEntranceDAReward struct {
	Beneficiary  common.Address
	SampleRound  *big.Int
	Epoch        *big.Int
	QuorumId     *big.Int
	DataRoot     [32]byte
	Quality      *big.Int
	LineIndex    *big.Int
	SublineIndex *big.Int
	Reward       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDAReward is a free log retrieval operation binding the contract event 0xc3898eb7106c1cb2f727da316a76320c0035f5692950aa7f6b65d20a5efaedc5.
//
// Solidity: event DAReward(address indexed beneficiary, uint256 indexed sampleRound, uint256 indexed epoch, uint256 quorumId, bytes32 dataRoot, uint256 quality, uint256 lineIndex, uint256 sublineIndex, uint256 reward)
func (_DAEntrance *DAEntranceFilterer) FilterDAReward(opts *bind.FilterOpts, beneficiary []common.Address, sampleRound []*big.Int, epoch []*big.Int) (*DAEntranceDARewardIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}
	var sampleRoundRule []interface{}
	for _, sampleRoundItem := range sampleRound {
		sampleRoundRule = append(sampleRoundRule, sampleRoundItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "DAReward", beneficiaryRule, sampleRoundRule, epochRule)
	if err != nil {
		return nil, err
	}
	return &DAEntranceDARewardIterator{contract: _DAEntrance.contract, event: "DAReward", logs: logs, sub: sub}, nil
}

// WatchDAReward is a free log subscription operation binding the contract event 0xc3898eb7106c1cb2f727da316a76320c0035f5692950aa7f6b65d20a5efaedc5.
//
// Solidity: event DAReward(address indexed beneficiary, uint256 indexed sampleRound, uint256 indexed epoch, uint256 quorumId, bytes32 dataRoot, uint256 quality, uint256 lineIndex, uint256 sublineIndex, uint256 reward)
func (_DAEntrance *DAEntranceFilterer) WatchDAReward(opts *bind.WatchOpts, sink chan<- *DAEntranceDAReward, beneficiary []common.Address, sampleRound []*big.Int, epoch []*big.Int) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}
	var sampleRoundRule []interface{}
	for _, sampleRoundItem := range sampleRound {
		sampleRoundRule = append(sampleRoundRule, sampleRoundItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "DAReward", beneficiaryRule, sampleRoundRule, epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceDAReward)
				if err := _DAEntrance.contract.UnpackLog(event, "DAReward", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDAReward is a log parse operation binding the contract event 0xc3898eb7106c1cb2f727da316a76320c0035f5692950aa7f6b65d20a5efaedc5.
//
// Solidity: event DAReward(address indexed beneficiary, uint256 indexed sampleRound, uint256 indexed epoch, uint256 quorumId, bytes32 dataRoot, uint256 quality, uint256 lineIndex, uint256 sublineIndex, uint256 reward)
func (_DAEntrance *DAEntranceFilterer) ParseDAReward(log types.Log) (*DAEntranceDAReward, error) {
	event := new(DAEntranceDAReward)
	if err := _DAEntrance.contract.UnpackLog(event, "DAReward", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceDataUploadIterator is returned from FilterDataUpload and is used to iterate over the raw logs and unpacked data for DataUpload events raised by the DAEntrance contract.
type DAEntranceDataUploadIterator struct {
	Event *DAEntranceDataUpload // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceDataUploadIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceDataUpload)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceDataUpload)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceDataUploadIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceDataUploadIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceDataUpload represents a DataUpload event raised by the DAEntrance contract.
type DAEntranceDataUpload struct {
	Sender    common.Address
	DataRoot  [32]byte
	Epoch     *big.Int
	QuorumId  *big.Int
	BlobPrice *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDataUpload is a free log retrieval operation binding the contract event 0x57b8b1a6583dc6ce934dfba3d66f2a8e1591b6e171bb2e0921cc64640277087b.
//
// Solidity: event DataUpload(address sender, bytes32 dataRoot, uint256 epoch, uint256 quorumId, uint256 blobPrice)
func (_DAEntrance *DAEntranceFilterer) FilterDataUpload(opts *bind.FilterOpts) (*DAEntranceDataUploadIterator, error) {

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "DataUpload")
	if err != nil {
		return nil, err
	}
	return &DAEntranceDataUploadIterator{contract: _DAEntrance.contract, event: "DataUpload", logs: logs, sub: sub}, nil
}

// WatchDataUpload is a free log subscription operation binding the contract event 0x57b8b1a6583dc6ce934dfba3d66f2a8e1591b6e171bb2e0921cc64640277087b.
//
// Solidity: event DataUpload(address sender, bytes32 dataRoot, uint256 epoch, uint256 quorumId, uint256 blobPrice)
func (_DAEntrance *DAEntranceFilterer) WatchDataUpload(opts *bind.WatchOpts, sink chan<- *DAEntranceDataUpload) (event.Subscription, error) {

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "DataUpload")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceDataUpload)
				if err := _DAEntrance.contract.UnpackLog(event, "DataUpload", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDataUpload is a log parse operation binding the contract event 0x57b8b1a6583dc6ce934dfba3d66f2a8e1591b6e171bb2e0921cc64640277087b.
//
// Solidity: event DataUpload(address sender, bytes32 dataRoot, uint256 epoch, uint256 quorumId, uint256 blobPrice)
func (_DAEntrance *DAEntranceFilterer) ParseDataUpload(log types.Log) (*DAEntranceDataUpload, error) {
	event := new(DAEntranceDataUpload)
	if err := _DAEntrance.contract.UnpackLog(event, "DataUpload", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceErasureCommitmentVerifiedIterator is returned from FilterErasureCommitmentVerified and is used to iterate over the raw logs and unpacked data for ErasureCommitmentVerified events raised by the DAEntrance contract.
type DAEntranceErasureCommitmentVerifiedIterator struct {
	Event *DAEntranceErasureCommitmentVerified // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceErasureCommitmentVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceErasureCommitmentVerified)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceErasureCommitmentVerified)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceErasureCommitmentVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceErasureCommitmentVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceErasureCommitmentVerified represents a ErasureCommitmentVerified event raised by the DAEntrance contract.
type DAEntranceErasureCommitmentVerified struct {
	DataRoot [32]byte
	Epoch    *big.Int
	QuorumId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterErasureCommitmentVerified is a free log retrieval operation binding the contract event 0x0f1b20d87bebd11dddaaab51f01cf2726880cb3f8073b636dbafa2aa8cacd256.
//
// Solidity: event ErasureCommitmentVerified(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
func (_DAEntrance *DAEntranceFilterer) FilterErasureCommitmentVerified(opts *bind.FilterOpts) (*DAEntranceErasureCommitmentVerifiedIterator, error) {

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "ErasureCommitmentVerified")
	if err != nil {
		return nil, err
	}
	return &DAEntranceErasureCommitmentVerifiedIterator{contract: _DAEntrance.contract, event: "ErasureCommitmentVerified", logs: logs, sub: sub}, nil
}

// WatchErasureCommitmentVerified is a free log subscription operation binding the contract event 0x0f1b20d87bebd11dddaaab51f01cf2726880cb3f8073b636dbafa2aa8cacd256.
//
// Solidity: event ErasureCommitmentVerified(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
func (_DAEntrance *DAEntranceFilterer) WatchErasureCommitmentVerified(opts *bind.WatchOpts, sink chan<- *DAEntranceErasureCommitmentVerified) (event.Subscription, error) {

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "ErasureCommitmentVerified")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceErasureCommitmentVerified)
				if err := _DAEntrance.contract.UnpackLog(event, "ErasureCommitmentVerified", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseErasureCommitmentVerified is a log parse operation binding the contract event 0x0f1b20d87bebd11dddaaab51f01cf2726880cb3f8073b636dbafa2aa8cacd256.
//
// Solidity: event ErasureCommitmentVerified(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
func (_DAEntrance *DAEntranceFilterer) ParseErasureCommitmentVerified(log types.Log) (*DAEntranceErasureCommitmentVerified, error) {
	event := new(DAEntranceErasureCommitmentVerified)
	if err := _DAEntrance.contract.UnpackLog(event, "ErasureCommitmentVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceNewSampleRoundIterator is returned from FilterNewSampleRound and is used to iterate over the raw logs and unpacked data for NewSampleRound events raised by the DAEntrance contract.
type DAEntranceNewSampleRoundIterator struct {
	Event *DAEntranceNewSampleRound // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceNewSampleRoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceNewSampleRound)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceNewSampleRound)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceNewSampleRoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceNewSampleRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceNewSampleRound represents a NewSampleRound event raised by the DAEntrance contract.
type DAEntranceNewSampleRound struct {
	SampleRound  *big.Int
	SampleHeight *big.Int
	SampleSeed   [32]byte
	PodasTarget  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterNewSampleRound is a free log retrieval operation binding the contract event 0xdfb5db5886e81f083727f21152a2a83457e99364e9f104e1aa10bbd6d9b4b95f.
//
// Solidity: event NewSampleRound(uint256 indexed sampleRound, uint256 sampleHeight, bytes32 sampleSeed, uint256 podasTarget)
func (_DAEntrance *DAEntranceFilterer) FilterNewSampleRound(opts *bind.FilterOpts, sampleRound []*big.Int) (*DAEntranceNewSampleRoundIterator, error) {

	var sampleRoundRule []interface{}
	for _, sampleRoundItem := range sampleRound {
		sampleRoundRule = append(sampleRoundRule, sampleRoundItem)
	}

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "NewSampleRound", sampleRoundRule)
	if err != nil {
		return nil, err
	}
	return &DAEntranceNewSampleRoundIterator{contract: _DAEntrance.contract, event: "NewSampleRound", logs: logs, sub: sub}, nil
}

// WatchNewSampleRound is a free log subscription operation binding the contract event 0xdfb5db5886e81f083727f21152a2a83457e99364e9f104e1aa10bbd6d9b4b95f.
//
// Solidity: event NewSampleRound(uint256 indexed sampleRound, uint256 sampleHeight, bytes32 sampleSeed, uint256 podasTarget)
func (_DAEntrance *DAEntranceFilterer) WatchNewSampleRound(opts *bind.WatchOpts, sink chan<- *DAEntranceNewSampleRound, sampleRound []*big.Int) (event.Subscription, error) {

	var sampleRoundRule []interface{}
	for _, sampleRoundItem := range sampleRound {
		sampleRoundRule = append(sampleRoundRule, sampleRoundItem)
	}

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "NewSampleRound", sampleRoundRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceNewSampleRound)
				if err := _DAEntrance.contract.UnpackLog(event, "NewSampleRound", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewSampleRound is a log parse operation binding the contract event 0xdfb5db5886e81f083727f21152a2a83457e99364e9f104e1aa10bbd6d9b4b95f.
//
// Solidity: event NewSampleRound(uint256 indexed sampleRound, uint256 sampleHeight, bytes32 sampleSeed, uint256 podasTarget)
func (_DAEntrance *DAEntranceFilterer) ParseNewSampleRound(log types.Log) (*DAEntranceNewSampleRound, error) {
	event := new(DAEntranceNewSampleRound)
	if err := _DAEntrance.contract.UnpackLog(event, "NewSampleRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the DAEntrance contract.
type DAEntranceRoleAdminChangedIterator struct {
	Event *DAEntranceRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceRoleAdminChanged represents a RoleAdminChanged event raised by the DAEntrance contract.
type DAEntranceRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DAEntrance *DAEntranceFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*DAEntranceRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &DAEntranceRoleAdminChangedIterator{contract: _DAEntrance.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DAEntrance *DAEntranceFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *DAEntranceRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceRoleAdminChanged)
				if err := _DAEntrance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DAEntrance *DAEntranceFilterer) ParseRoleAdminChanged(log types.Log) (*DAEntranceRoleAdminChanged, error) {
	event := new(DAEntranceRoleAdminChanged)
	if err := _DAEntrance.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the DAEntrance contract.
type DAEntranceRoleGrantedIterator struct {
	Event *DAEntranceRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceRoleGranted represents a RoleGranted event raised by the DAEntrance contract.
type DAEntranceRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DAEntranceRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DAEntranceRoleGrantedIterator{contract: _DAEntrance.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *DAEntranceRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceRoleGranted)
				if err := _DAEntrance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) ParseRoleGranted(log types.Log) (*DAEntranceRoleGranted, error) {
	event := new(DAEntranceRoleGranted)
	if err := _DAEntrance.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DAEntranceRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the DAEntrance contract.
type DAEntranceRoleRevokedIterator struct {
	Event *DAEntranceRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DAEntranceRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DAEntranceRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DAEntranceRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DAEntranceRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DAEntranceRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DAEntranceRoleRevoked represents a RoleRevoked event raised by the DAEntrance contract.
type DAEntranceRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DAEntranceRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DAEntranceRoleRevokedIterator{contract: _DAEntrance.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *DAEntranceRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DAEntrance.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DAEntranceRoleRevoked)
				if err := _DAEntrance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DAEntrance *DAEntranceFilterer) ParseRoleRevoked(log types.Log) (*DAEntranceRoleRevoked, error) {
	event := new(DAEntranceRoleRevoked)
	if err := _DAEntrance.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
