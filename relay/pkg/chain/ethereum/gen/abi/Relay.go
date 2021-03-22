// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RelayABI is the input ABI used to generate the binding from.
const RelayABI = "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_genesisHeader\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_periodStart\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_first\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_last\",\"type\":\"bytes32\"}],\"name\":\"Extension\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_from\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_to\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_gcd\",\"type\":\"bytes32\"}],\"name\":\"NewTip\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"HEIGHT_INTERVAL\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_anchor\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_headers\",\"type\":\"bytes\"}],\"name\":\"addHeaders\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_oldPeriodStartHeader\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_oldPeriodEndHeader\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_headers\",\"type\":\"bytes\"}],\"name\":\"addHeadersWithRetarget\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_offset\",\"type\":\"uint256\"}],\"name\":\"findAncestor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"}],\"name\":\"findHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getBestKnownDigest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentEpochDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLastReorgCommonAncestor\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPrevEpochDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRelayGenesis\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ancestor\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_descendant\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"isAncestor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_ancestor\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_currentBest\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_newBest\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_limit\",\"type\":\"uint256\"}],\"name\":\"markNewHeaviest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Relay is an auto generated Go binding around an Ethereum contract.
type Relay struct {
	RelayCaller     // Read-only binding to the contract
	RelayTransactor // Write-only binding to the contract
	RelayFilterer   // Log filterer for contract events
}

// RelayCaller is an auto generated read-only Go binding around an Ethereum contract.
type RelayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RelayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RelayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RelaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RelaySession struct {
	Contract     *Relay            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RelayCallerSession struct {
	Contract *RelayCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RelayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RelayTransactorSession struct {
	Contract     *RelayTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RelayRaw is an auto generated low-level Go binding around an Ethereum contract.
type RelayRaw struct {
	Contract *Relay // Generic contract binding to access the raw methods on
}

// RelayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RelayCallerRaw struct {
	Contract *RelayCaller // Generic read-only contract binding to access the raw methods on
}

// RelayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RelayTransactorRaw struct {
	Contract *RelayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRelay creates a new instance of Relay, bound to a specific deployed contract.
func NewRelay(address common.Address, backend bind.ContractBackend) (*Relay, error) {
	contract, err := bindRelay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Relay{RelayCaller: RelayCaller{contract: contract}, RelayTransactor: RelayTransactor{contract: contract}, RelayFilterer: RelayFilterer{contract: contract}}, nil
}

// NewRelayCaller creates a new read-only instance of Relay, bound to a specific deployed contract.
func NewRelayCaller(address common.Address, caller bind.ContractCaller) (*RelayCaller, error) {
	contract, err := bindRelay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RelayCaller{contract: contract}, nil
}

// NewRelayTransactor creates a new write-only instance of Relay, bound to a specific deployed contract.
func NewRelayTransactor(address common.Address, transactor bind.ContractTransactor) (*RelayTransactor, error) {
	contract, err := bindRelay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RelayTransactor{contract: contract}, nil
}

// NewRelayFilterer creates a new log filterer instance of Relay, bound to a specific deployed contract.
func NewRelayFilterer(address common.Address, filterer bind.ContractFilterer) (*RelayFilterer, error) {
	contract, err := bindRelay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RelayFilterer{contract: contract}, nil
}

// bindRelay binds a generic wrapper to an already deployed contract.
func bindRelay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RelayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.RelayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.RelayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Relay *RelayCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Relay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Relay *RelayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Relay *RelayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Relay.Contract.contract.Transact(opts, method, params...)
}

// HEIGHTINTERVAL is a free data retrieval call binding the contract method 0x70d53c18.
//
// Solidity: function HEIGHT_INTERVAL() constant returns(uint32)
func (_Relay *RelayCaller) HEIGHTINTERVAL(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "HEIGHT_INTERVAL")
	return *ret0, err
}

// HEIGHTINTERVAL is a free data retrieval call binding the contract method 0x70d53c18.
//
// Solidity: function HEIGHT_INTERVAL() constant returns(uint32)
func (_Relay *RelaySession) HEIGHTINTERVAL() (uint32, error) {
	return _Relay.Contract.HEIGHTINTERVAL(&_Relay.CallOpts)
}

// HEIGHTINTERVAL is a free data retrieval call binding the contract method 0x70d53c18.
//
// Solidity: function HEIGHT_INTERVAL() constant returns(uint32)
func (_Relay *RelayCallerSession) HEIGHTINTERVAL() (uint32, error) {
	return _Relay.Contract.HEIGHTINTERVAL(&_Relay.CallOpts)
}

// FindAncestor is a free data retrieval call binding the contract method 0x30017b3b.
//
// Solidity: function findAncestor(bytes32 _digest, uint256 _offset) constant returns(bytes32)
func (_Relay *RelayCaller) FindAncestor(opts *bind.CallOpts, _digest [32]byte, _offset *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "findAncestor", _digest, _offset)
	return *ret0, err
}

// FindAncestor is a free data retrieval call binding the contract method 0x30017b3b.
//
// Solidity: function findAncestor(bytes32 _digest, uint256 _offset) constant returns(bytes32)
func (_Relay *RelaySession) FindAncestor(_digest [32]byte, _offset *big.Int) ([32]byte, error) {
	return _Relay.Contract.FindAncestor(&_Relay.CallOpts, _digest, _offset)
}

// FindAncestor is a free data retrieval call binding the contract method 0x30017b3b.
//
// Solidity: function findAncestor(bytes32 _digest, uint256 _offset) constant returns(bytes32)
func (_Relay *RelayCallerSession) FindAncestor(_digest [32]byte, _offset *big.Int) ([32]byte, error) {
	return _Relay.Contract.FindAncestor(&_Relay.CallOpts, _digest, _offset)
}

// FindHeight is a free data retrieval call binding the contract method 0x60b5c390.
//
// Solidity: function findHeight(bytes32 _digest) constant returns(uint256)
func (_Relay *RelayCaller) FindHeight(opts *bind.CallOpts, _digest [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "findHeight", _digest)
	return *ret0, err
}

// FindHeight is a free data retrieval call binding the contract method 0x60b5c390.
//
// Solidity: function findHeight(bytes32 _digest) constant returns(uint256)
func (_Relay *RelaySession) FindHeight(_digest [32]byte) (*big.Int, error) {
	return _Relay.Contract.FindHeight(&_Relay.CallOpts, _digest)
}

// FindHeight is a free data retrieval call binding the contract method 0x60b5c390.
//
// Solidity: function findHeight(bytes32 _digest) constant returns(uint256)
func (_Relay *RelayCallerSession) FindHeight(_digest [32]byte) (*big.Int, error) {
	return _Relay.Contract.FindHeight(&_Relay.CallOpts, _digest)
}

// GetBestKnownDigest is a free data retrieval call binding the contract method 0x1910d973.
//
// Solidity: function getBestKnownDigest() constant returns(bytes32)
func (_Relay *RelayCaller) GetBestKnownDigest(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "getBestKnownDigest")
	return *ret0, err
}

// GetBestKnownDigest is a free data retrieval call binding the contract method 0x1910d973.
//
// Solidity: function getBestKnownDigest() constant returns(bytes32)
func (_Relay *RelaySession) GetBestKnownDigest() ([32]byte, error) {
	return _Relay.Contract.GetBestKnownDigest(&_Relay.CallOpts)
}

// GetBestKnownDigest is a free data retrieval call binding the contract method 0x1910d973.
//
// Solidity: function getBestKnownDigest() constant returns(bytes32)
func (_Relay *RelayCallerSession) GetBestKnownDigest() ([32]byte, error) {
	return _Relay.Contract.GetBestKnownDigest(&_Relay.CallOpts)
}

// GetCurrentEpochDifficulty is a free data retrieval call binding the contract method 0x113764be.
//
// Solidity: function getCurrentEpochDifficulty() constant returns(uint256)
func (_Relay *RelayCaller) GetCurrentEpochDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "getCurrentEpochDifficulty")
	return *ret0, err
}

// GetCurrentEpochDifficulty is a free data retrieval call binding the contract method 0x113764be.
//
// Solidity: function getCurrentEpochDifficulty() constant returns(uint256)
func (_Relay *RelaySession) GetCurrentEpochDifficulty() (*big.Int, error) {
	return _Relay.Contract.GetCurrentEpochDifficulty(&_Relay.CallOpts)
}

// GetCurrentEpochDifficulty is a free data retrieval call binding the contract method 0x113764be.
//
// Solidity: function getCurrentEpochDifficulty() constant returns(uint256)
func (_Relay *RelayCallerSession) GetCurrentEpochDifficulty() (*big.Int, error) {
	return _Relay.Contract.GetCurrentEpochDifficulty(&_Relay.CallOpts)
}

// GetLastReorgCommonAncestor is a free data retrieval call binding the contract method 0xc58242cd.
//
// Solidity: function getLastReorgCommonAncestor() constant returns(bytes32)
func (_Relay *RelayCaller) GetLastReorgCommonAncestor(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "getLastReorgCommonAncestor")
	return *ret0, err
}

// GetLastReorgCommonAncestor is a free data retrieval call binding the contract method 0xc58242cd.
//
// Solidity: function getLastReorgCommonAncestor() constant returns(bytes32)
func (_Relay *RelaySession) GetLastReorgCommonAncestor() ([32]byte, error) {
	return _Relay.Contract.GetLastReorgCommonAncestor(&_Relay.CallOpts)
}

// GetLastReorgCommonAncestor is a free data retrieval call binding the contract method 0xc58242cd.
//
// Solidity: function getLastReorgCommonAncestor() constant returns(bytes32)
func (_Relay *RelayCallerSession) GetLastReorgCommonAncestor() ([32]byte, error) {
	return _Relay.Contract.GetLastReorgCommonAncestor(&_Relay.CallOpts)
}

// GetPrevEpochDifficulty is a free data retrieval call binding the contract method 0x2b97be24.
//
// Solidity: function getPrevEpochDifficulty() constant returns(uint256)
func (_Relay *RelayCaller) GetPrevEpochDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "getPrevEpochDifficulty")
	return *ret0, err
}

// GetPrevEpochDifficulty is a free data retrieval call binding the contract method 0x2b97be24.
//
// Solidity: function getPrevEpochDifficulty() constant returns(uint256)
func (_Relay *RelaySession) GetPrevEpochDifficulty() (*big.Int, error) {
	return _Relay.Contract.GetPrevEpochDifficulty(&_Relay.CallOpts)
}

// GetPrevEpochDifficulty is a free data retrieval call binding the contract method 0x2b97be24.
//
// Solidity: function getPrevEpochDifficulty() constant returns(uint256)
func (_Relay *RelayCallerSession) GetPrevEpochDifficulty() (*big.Int, error) {
	return _Relay.Contract.GetPrevEpochDifficulty(&_Relay.CallOpts)
}

// GetRelayGenesis is a free data retrieval call binding the contract method 0xe3d8d8d8.
//
// Solidity: function getRelayGenesis() constant returns(bytes32)
func (_Relay *RelayCaller) GetRelayGenesis(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "getRelayGenesis")
	return *ret0, err
}

// GetRelayGenesis is a free data retrieval call binding the contract method 0xe3d8d8d8.
//
// Solidity: function getRelayGenesis() constant returns(bytes32)
func (_Relay *RelaySession) GetRelayGenesis() ([32]byte, error) {
	return _Relay.Contract.GetRelayGenesis(&_Relay.CallOpts)
}

// GetRelayGenesis is a free data retrieval call binding the contract method 0xe3d8d8d8.
//
// Solidity: function getRelayGenesis() constant returns(bytes32)
func (_Relay *RelayCallerSession) GetRelayGenesis() ([32]byte, error) {
	return _Relay.Contract.GetRelayGenesis(&_Relay.CallOpts)
}

// IsAncestor is a free data retrieval call binding the contract method 0xb985621a.
//
// Solidity: function isAncestor(bytes32 _ancestor, bytes32 _descendant, uint256 _limit) constant returns(bool)
func (_Relay *RelayCaller) IsAncestor(opts *bind.CallOpts, _ancestor [32]byte, _descendant [32]byte, _limit *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Relay.contract.Call(opts, out, "isAncestor", _ancestor, _descendant, _limit)
	return *ret0, err
}

// IsAncestor is a free data retrieval call binding the contract method 0xb985621a.
//
// Solidity: function isAncestor(bytes32 _ancestor, bytes32 _descendant, uint256 _limit) constant returns(bool)
func (_Relay *RelaySession) IsAncestor(_ancestor [32]byte, _descendant [32]byte, _limit *big.Int) (bool, error) {
	return _Relay.Contract.IsAncestor(&_Relay.CallOpts, _ancestor, _descendant, _limit)
}

// IsAncestor is a free data retrieval call binding the contract method 0xb985621a.
//
// Solidity: function isAncestor(bytes32 _ancestor, bytes32 _descendant, uint256 _limit) constant returns(bool)
func (_Relay *RelayCallerSession) IsAncestor(_ancestor [32]byte, _descendant [32]byte, _limit *big.Int) (bool, error) {
	return _Relay.Contract.IsAncestor(&_Relay.CallOpts, _ancestor, _descendant, _limit)
}

// AddHeaders is a paid mutator transaction binding the contract method 0x65da41b9.
//
// Solidity: function addHeaders(bytes _anchor, bytes _headers) returns(bool)
func (_Relay *RelayTransactor) AddHeaders(opts *bind.TransactOpts, _anchor []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "addHeaders", _anchor, _headers)
}

// AddHeaders is a paid mutator transaction binding the contract method 0x65da41b9.
//
// Solidity: function addHeaders(bytes _anchor, bytes _headers) returns(bool)
func (_Relay *RelaySession) AddHeaders(_anchor []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.Contract.AddHeaders(&_Relay.TransactOpts, _anchor, _headers)
}

// AddHeaders is a paid mutator transaction binding the contract method 0x65da41b9.
//
// Solidity: function addHeaders(bytes _anchor, bytes _headers) returns(bool)
func (_Relay *RelayTransactorSession) AddHeaders(_anchor []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.Contract.AddHeaders(&_Relay.TransactOpts, _anchor, _headers)
}

// AddHeadersWithRetarget is a paid mutator transaction binding the contract method 0x7fa637fc.
//
// Solidity: function addHeadersWithRetarget(bytes _oldPeriodStartHeader, bytes _oldPeriodEndHeader, bytes _headers) returns(bool)
func (_Relay *RelayTransactor) AddHeadersWithRetarget(opts *bind.TransactOpts, _oldPeriodStartHeader []byte, _oldPeriodEndHeader []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "addHeadersWithRetarget", _oldPeriodStartHeader, _oldPeriodEndHeader, _headers)
}

// AddHeadersWithRetarget is a paid mutator transaction binding the contract method 0x7fa637fc.
//
// Solidity: function addHeadersWithRetarget(bytes _oldPeriodStartHeader, bytes _oldPeriodEndHeader, bytes _headers) returns(bool)
func (_Relay *RelaySession) AddHeadersWithRetarget(_oldPeriodStartHeader []byte, _oldPeriodEndHeader []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.Contract.AddHeadersWithRetarget(&_Relay.TransactOpts, _oldPeriodStartHeader, _oldPeriodEndHeader, _headers)
}

// AddHeadersWithRetarget is a paid mutator transaction binding the contract method 0x7fa637fc.
//
// Solidity: function addHeadersWithRetarget(bytes _oldPeriodStartHeader, bytes _oldPeriodEndHeader, bytes _headers) returns(bool)
func (_Relay *RelayTransactorSession) AddHeadersWithRetarget(_oldPeriodStartHeader []byte, _oldPeriodEndHeader []byte, _headers []byte) (*types.Transaction, error) {
	return _Relay.Contract.AddHeadersWithRetarget(&_Relay.TransactOpts, _oldPeriodStartHeader, _oldPeriodEndHeader, _headers)
}

// MarkNewHeaviest is a paid mutator transaction binding the contract method 0x74c3a3a9.
//
// Solidity: function markNewHeaviest(bytes32 _ancestor, bytes _currentBest, bytes _newBest, uint256 _limit) returns(bool)
func (_Relay *RelayTransactor) MarkNewHeaviest(opts *bind.TransactOpts, _ancestor [32]byte, _currentBest []byte, _newBest []byte, _limit *big.Int) (*types.Transaction, error) {
	return _Relay.contract.Transact(opts, "markNewHeaviest", _ancestor, _currentBest, _newBest, _limit)
}

// MarkNewHeaviest is a paid mutator transaction binding the contract method 0x74c3a3a9.
//
// Solidity: function markNewHeaviest(bytes32 _ancestor, bytes _currentBest, bytes _newBest, uint256 _limit) returns(bool)
func (_Relay *RelaySession) MarkNewHeaviest(_ancestor [32]byte, _currentBest []byte, _newBest []byte, _limit *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.MarkNewHeaviest(&_Relay.TransactOpts, _ancestor, _currentBest, _newBest, _limit)
}

// MarkNewHeaviest is a paid mutator transaction binding the contract method 0x74c3a3a9.
//
// Solidity: function markNewHeaviest(bytes32 _ancestor, bytes _currentBest, bytes _newBest, uint256 _limit) returns(bool)
func (_Relay *RelayTransactorSession) MarkNewHeaviest(_ancestor [32]byte, _currentBest []byte, _newBest []byte, _limit *big.Int) (*types.Transaction, error) {
	return _Relay.Contract.MarkNewHeaviest(&_Relay.TransactOpts, _ancestor, _currentBest, _newBest, _limit)
}

// RelayExtensionIterator is returned from FilterExtension and is used to iterate over the raw logs and unpacked data for Extension events raised by the Relay contract.
type RelayExtensionIterator struct {
	Event *RelayExtension // Event containing the contract specifics and raw log

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
func (it *RelayExtensionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayExtension)
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
		it.Event = new(RelayExtension)
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
func (it *RelayExtensionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayExtensionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayExtension represents a Extension event raised by the Relay contract.
type RelayExtension struct {
	First [32]byte
	Last  [32]byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterExtension is a free log retrieval operation binding the contract event 0xf90e4f1d9cd0dd55e339411cbc9b152482307c3a23ed64715e4a2858f641a3f5.
//
// Solidity: event Extension(bytes32 indexed _first, bytes32 indexed _last)
func (_Relay *RelayFilterer) FilterExtension(opts *bind.FilterOpts, _first [][32]byte, _last [][32]byte) (*RelayExtensionIterator, error) {

	var _firstRule []interface{}
	for _, _firstItem := range _first {
		_firstRule = append(_firstRule, _firstItem)
	}
	var _lastRule []interface{}
	for _, _lastItem := range _last {
		_lastRule = append(_lastRule, _lastItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "Extension", _firstRule, _lastRule)
	if err != nil {
		return nil, err
	}
	return &RelayExtensionIterator{contract: _Relay.contract, event: "Extension", logs: logs, sub: sub}, nil
}

// WatchExtension is a free log subscription operation binding the contract event 0xf90e4f1d9cd0dd55e339411cbc9b152482307c3a23ed64715e4a2858f641a3f5.
//
// Solidity: event Extension(bytes32 indexed _first, bytes32 indexed _last)
func (_Relay *RelayFilterer) WatchExtension(opts *bind.WatchOpts, sink chan<- *RelayExtension, _first [][32]byte, _last [][32]byte) (event.Subscription, error) {

	var _firstRule []interface{}
	for _, _firstItem := range _first {
		_firstRule = append(_firstRule, _firstItem)
	}
	var _lastRule []interface{}
	for _, _lastItem := range _last {
		_lastRule = append(_lastRule, _lastItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "Extension", _firstRule, _lastRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayExtension)
				if err := _Relay.contract.UnpackLog(event, "Extension", log); err != nil {
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

// ParseExtension is a log parse operation binding the contract event 0xf90e4f1d9cd0dd55e339411cbc9b152482307c3a23ed64715e4a2858f641a3f5.
//
// Solidity: event Extension(bytes32 indexed _first, bytes32 indexed _last)
func (_Relay *RelayFilterer) ParseExtension(log types.Log) (*RelayExtension, error) {
	event := new(RelayExtension)
	if err := _Relay.contract.UnpackLog(event, "Extension", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RelayNewTipIterator is returned from FilterNewTip and is used to iterate over the raw logs and unpacked data for NewTip events raised by the Relay contract.
type RelayNewTipIterator struct {
	Event *RelayNewTip // Event containing the contract specifics and raw log

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
func (it *RelayNewTipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RelayNewTip)
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
		it.Event = new(RelayNewTip)
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
func (it *RelayNewTipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RelayNewTipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RelayNewTip represents a NewTip event raised by the Relay contract.
type RelayNewTip struct {
	From [32]byte
	To   [32]byte
	Gcd  [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewTip is a free log retrieval operation binding the contract event 0x3cc13de64df0f0239626235c51a2da251bbc8c85664ecce39263da3ee03f606c.
//
// Solidity: event NewTip(bytes32 indexed _from, bytes32 indexed _to, bytes32 indexed _gcd)
func (_Relay *RelayFilterer) FilterNewTip(opts *bind.FilterOpts, _from [][32]byte, _to [][32]byte, _gcd [][32]byte) (*RelayNewTipIterator, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _gcdRule []interface{}
	for _, _gcdItem := range _gcd {
		_gcdRule = append(_gcdRule, _gcdItem)
	}

	logs, sub, err := _Relay.contract.FilterLogs(opts, "NewTip", _fromRule, _toRule, _gcdRule)
	if err != nil {
		return nil, err
	}
	return &RelayNewTipIterator{contract: _Relay.contract, event: "NewTip", logs: logs, sub: sub}, nil
}

// WatchNewTip is a free log subscription operation binding the contract event 0x3cc13de64df0f0239626235c51a2da251bbc8c85664ecce39263da3ee03f606c.
//
// Solidity: event NewTip(bytes32 indexed _from, bytes32 indexed _to, bytes32 indexed _gcd)
func (_Relay *RelayFilterer) WatchNewTip(opts *bind.WatchOpts, sink chan<- *RelayNewTip, _from [][32]byte, _to [][32]byte, _gcd [][32]byte) (event.Subscription, error) {

	var _fromRule []interface{}
	for _, _fromItem := range _from {
		_fromRule = append(_fromRule, _fromItem)
	}
	var _toRule []interface{}
	for _, _toItem := range _to {
		_toRule = append(_toRule, _toItem)
	}
	var _gcdRule []interface{}
	for _, _gcdItem := range _gcd {
		_gcdRule = append(_gcdRule, _gcdItem)
	}

	logs, sub, err := _Relay.contract.WatchLogs(opts, "NewTip", _fromRule, _toRule, _gcdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RelayNewTip)
				if err := _Relay.contract.UnpackLog(event, "NewTip", log); err != nil {
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

// ParseNewTip is a log parse operation binding the contract event 0x3cc13de64df0f0239626235c51a2da251bbc8c85664ecce39263da3ee03f606c.
//
// Solidity: event NewTip(bytes32 indexed _from, bytes32 indexed _to, bytes32 indexed _gcd)
func (_Relay *RelayFilterer) ParseNewTip(log types.Log) (*RelayNewTip, error) {
	event := new(RelayNewTip)
	if err := _Relay.contract.UnpackLog(event, "NewTip", log); err != nil {
		return nil, err
	}
	return event, nil
}
