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

// DAEntranceMetaData contains all meta data concerning the DAEntrance contract.
var DAEntranceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"}],\"name\":\"DataUpload\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"}],\"name\":\"ErasureCommitmentVerified\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DA_SIGNERS\",\"outputs\":[{\"internalType\":\"contractIDASigners\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLICE_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SLICE_NUMERATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_dataRoots\",\"type\":\"bytes32[]\"}],\"name\":\"submitOriginalData\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"quorumId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"erasureCommitment\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"quorumBitmap\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"X\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"Y\",\"type\":\"uint256[2]\"}],\"internalType\":\"structBN254.G2Point\",\"name\":\"aggPkG2\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"signature\",\"type\":\"tuple\"}],\"internalType\":\"structIDAEntrance.CommitRootSubmission[]\",\"name\":\"_submissions\",\"type\":\"tuple[]\"}],\"name\":\"submitVerifiedCommitRoots\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_dataRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_quorumId\",\"type\":\"uint256\"}],\"name\":\"verifiedErasureCommitment\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"X\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"Y\",\"type\":\"uint256\"}],\"internalType\":\"structBN254.G1Point\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
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

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DAEntrance *DAEntranceTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DAEntrance.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DAEntrance *DAEntranceSession) Receive() (*types.Transaction, error) {
	return _DAEntrance.Contract.Receive(&_DAEntrance.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_DAEntrance *DAEntranceTransactorSession) Receive() (*types.Transaction, error) {
	return _DAEntrance.Contract.Receive(&_DAEntrance.TransactOpts)
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
	DataRoot [32]byte
	Epoch    *big.Int
	QuorumId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDataUpload is a free log retrieval operation binding the contract event 0xf0bf37f8713754493879920443065424c575888634675f146c115709bbb59acb.
//
// Solidity: event DataUpload(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
func (_DAEntrance *DAEntranceFilterer) FilterDataUpload(opts *bind.FilterOpts) (*DAEntranceDataUploadIterator, error) {

	logs, sub, err := _DAEntrance.contract.FilterLogs(opts, "DataUpload")
	if err != nil {
		return nil, err
	}
	return &DAEntranceDataUploadIterator{contract: _DAEntrance.contract, event: "DataUpload", logs: logs, sub: sub}, nil
}

// WatchDataUpload is a free log subscription operation binding the contract event 0xf0bf37f8713754493879920443065424c575888634675f146c115709bbb59acb.
//
// Solidity: event DataUpload(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
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

// ParseDataUpload is a log parse operation binding the contract event 0xf0bf37f8713754493879920443065424c575888634675f146c115709bbb59acb.
//
// Solidity: event DataUpload(bytes32 dataRoot, uint256 epoch, uint256 quorumId)
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
