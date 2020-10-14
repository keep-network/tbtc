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

// DepositLogABI is the input ABI used to generate the binding from.
const DepositLogABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"CourtesyCalled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_keepAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Created\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"ExitedCourtesyCall\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"FraudDuringSetup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Funded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_abortOutputScript\",\"type\":\"bytes\"}],\"name\":\"FunderAbortRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"GotRedemptionSignature\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Liquidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_utxoValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requestedFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_outpoint\",\"type\":\"bytes\"}],\"name\":\"RedemptionRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyX\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyY\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"RegisteredPubkey\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"SetupFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_wasFraud\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"StartedLiquidation\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"}],\"name\":\"approvedToLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logCourtesyCalled\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_keepAddress\",\"type\":\"address\"}],\"name\":\"logCreated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logExitedCourtesyCall\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logFraudDuringSetup\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"}],\"name\":\"logFunded\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_abortOutputScript\",\"type\":\"bytes\"}],\"name\":\"logFunderRequestedAbort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"logGotRedemptionSignature\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logLiquidated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"}],\"name\":\"logRedeemed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_utxoValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_requestedFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_outpoint\",\"type\":\"bytes\"}],\"name\":\"logRedemptionRequested\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyX\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyY\",\"type\":\"bytes32\"}],\"name\":\"logRegisteredPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logSetupFailed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_wasFraud\",\"type\":\"bool\"}],\"name\":\"logStartedLiquidation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DepositLog is an auto generated Go binding around an Ethereum contract.
type DepositLog struct {
	DepositLogCaller     // Read-only binding to the contract
	DepositLogTransactor // Write-only binding to the contract
	DepositLogFilterer   // Log filterer for contract events
}

// DepositLogCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositLogCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositLogTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositLogTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositLogFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositLogFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositLogSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositLogSession struct {
	Contract     *DepositLog       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositLogCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositLogCallerSession struct {
	Contract *DepositLogCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// DepositLogTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositLogTransactorSession struct {
	Contract     *DepositLogTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// DepositLogRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositLogRaw struct {
	Contract *DepositLog // Generic contract binding to access the raw methods on
}

// DepositLogCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositLogCallerRaw struct {
	Contract *DepositLogCaller // Generic read-only contract binding to access the raw methods on
}

// DepositLogTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositLogTransactorRaw struct {
	Contract *DepositLogTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositLog creates a new instance of DepositLog, bound to a specific deployed contract.
func NewDepositLog(address common.Address, backend bind.ContractBackend) (*DepositLog, error) {
	contract, err := bindDepositLog(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositLog{DepositLogCaller: DepositLogCaller{contract: contract}, DepositLogTransactor: DepositLogTransactor{contract: contract}, DepositLogFilterer: DepositLogFilterer{contract: contract}}, nil
}

// NewDepositLogCaller creates a new read-only instance of DepositLog, bound to a specific deployed contract.
func NewDepositLogCaller(address common.Address, caller bind.ContractCaller) (*DepositLogCaller, error) {
	contract, err := bindDepositLog(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositLogCaller{contract: contract}, nil
}

// NewDepositLogTransactor creates a new write-only instance of DepositLog, bound to a specific deployed contract.
func NewDepositLogTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositLogTransactor, error) {
	contract, err := bindDepositLog(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositLogTransactor{contract: contract}, nil
}

// NewDepositLogFilterer creates a new log filterer instance of DepositLog, bound to a specific deployed contract.
func NewDepositLogFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositLogFilterer, error) {
	contract, err := bindDepositLog(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositLogFilterer{contract: contract}, nil
}

// bindDepositLog binds a generic wrapper to an already deployed contract.
func bindDepositLog(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositLogABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositLog *DepositLogRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DepositLog.Contract.DepositLogCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositLog *DepositLogRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.Contract.DepositLogTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositLog *DepositLogRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositLog.Contract.DepositLogTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositLog *DepositLogCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DepositLog.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositLog *DepositLogTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositLog *DepositLogTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositLog.Contract.contract.Transact(opts, method, params...)
}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) constant returns(bool)
func (_DepositLog *DepositLogCaller) ApprovedToLog(opts *bind.CallOpts, _caller common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DepositLog.contract.Call(opts, out, "approvedToLog", _caller)
	return *ret0, err
}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) constant returns(bool)
func (_DepositLog *DepositLogSession) ApprovedToLog(_caller common.Address) (bool, error) {
	return _DepositLog.Contract.ApprovedToLog(&_DepositLog.CallOpts, _caller)
}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) constant returns(bool)
func (_DepositLog *DepositLogCallerSession) ApprovedToLog(_caller common.Address) (bool, error) {
	return _DepositLog.Contract.ApprovedToLog(&_DepositLog.CallOpts, _caller)
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_DepositLog *DepositLogTransactor) LogCourtesyCalled(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logCourtesyCalled")
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_DepositLog *DepositLogSession) LogCourtesyCalled() (*types.Transaction, error) {
	return _DepositLog.Contract.LogCourtesyCalled(&_DepositLog.TransactOpts)
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_DepositLog *DepositLogTransactorSession) LogCourtesyCalled() (*types.Transaction, error) {
	return _DepositLog.Contract.LogCourtesyCalled(&_DepositLog.TransactOpts)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_DepositLog *DepositLogTransactor) LogCreated(opts *bind.TransactOpts, _keepAddress common.Address) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logCreated", _keepAddress)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_DepositLog *DepositLogSession) LogCreated(_keepAddress common.Address) (*types.Transaction, error) {
	return _DepositLog.Contract.LogCreated(&_DepositLog.TransactOpts, _keepAddress)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_DepositLog *DepositLogTransactorSession) LogCreated(_keepAddress common.Address) (*types.Transaction, error) {
	return _DepositLog.Contract.LogCreated(&_DepositLog.TransactOpts, _keepAddress)
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_DepositLog *DepositLogTransactor) LogExitedCourtesyCall(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logExitedCourtesyCall")
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_DepositLog *DepositLogSession) LogExitedCourtesyCall() (*types.Transaction, error) {
	return _DepositLog.Contract.LogExitedCourtesyCall(&_DepositLog.TransactOpts)
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_DepositLog *DepositLogTransactorSession) LogExitedCourtesyCall() (*types.Transaction, error) {
	return _DepositLog.Contract.LogExitedCourtesyCall(&_DepositLog.TransactOpts)
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_DepositLog *DepositLogTransactor) LogFraudDuringSetup(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logFraudDuringSetup")
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_DepositLog *DepositLogSession) LogFraudDuringSetup() (*types.Transaction, error) {
	return _DepositLog.Contract.LogFraudDuringSetup(&_DepositLog.TransactOpts)
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_DepositLog *DepositLogTransactorSession) LogFraudDuringSetup() (*types.Transaction, error) {
	return _DepositLog.Contract.LogFraudDuringSetup(&_DepositLog.TransactOpts)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_DepositLog *DepositLogTransactor) LogFunded(opts *bind.TransactOpts, _txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logFunded", _txid)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_DepositLog *DepositLogSession) LogFunded(_txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogFunded(&_DepositLog.TransactOpts, _txid)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_DepositLog *DepositLogTransactorSession) LogFunded(_txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogFunded(&_DepositLog.TransactOpts, _txid)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_DepositLog *DepositLogTransactor) LogFunderRequestedAbort(opts *bind.TransactOpts, _abortOutputScript []byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logFunderRequestedAbort", _abortOutputScript)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_DepositLog *DepositLogSession) LogFunderRequestedAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogFunderRequestedAbort(&_DepositLog.TransactOpts, _abortOutputScript)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_DepositLog *DepositLogTransactorSession) LogFunderRequestedAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogFunderRequestedAbort(&_DepositLog.TransactOpts, _abortOutputScript)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_DepositLog *DepositLogTransactor) LogGotRedemptionSignature(opts *bind.TransactOpts, _digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logGotRedemptionSignature", _digest, _r, _s)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_DepositLog *DepositLogSession) LogGotRedemptionSignature(_digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogGotRedemptionSignature(&_DepositLog.TransactOpts, _digest, _r, _s)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_DepositLog *DepositLogTransactorSession) LogGotRedemptionSignature(_digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogGotRedemptionSignature(&_DepositLog.TransactOpts, _digest, _r, _s)
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_DepositLog *DepositLogTransactor) LogLiquidated(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logLiquidated")
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_DepositLog *DepositLogSession) LogLiquidated() (*types.Transaction, error) {
	return _DepositLog.Contract.LogLiquidated(&_DepositLog.TransactOpts)
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_DepositLog *DepositLogTransactorSession) LogLiquidated() (*types.Transaction, error) {
	return _DepositLog.Contract.LogLiquidated(&_DepositLog.TransactOpts)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_DepositLog *DepositLogTransactor) LogRedeemed(opts *bind.TransactOpts, _txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logRedeemed", _txid)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_DepositLog *DepositLogSession) LogRedeemed(_txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRedeemed(&_DepositLog.TransactOpts, _txid)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_DepositLog *DepositLogTransactorSession) LogRedeemed(_txid [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRedeemed(&_DepositLog.TransactOpts, _txid)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_DepositLog *DepositLogTransactor) LogRedemptionRequested(opts *bind.TransactOpts, _requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logRedemptionRequested", _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_DepositLog *DepositLogSession) LogRedemptionRequested(_requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRedemptionRequested(&_DepositLog.TransactOpts, _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_DepositLog *DepositLogTransactorSession) LogRedemptionRequested(_requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRedemptionRequested(&_DepositLog.TransactOpts, _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_DepositLog *DepositLogTransactor) LogRegisteredPubkey(opts *bind.TransactOpts, _signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logRegisteredPubkey", _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_DepositLog *DepositLogSession) LogRegisteredPubkey(_signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRegisteredPubkey(&_DepositLog.TransactOpts, _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_DepositLog *DepositLogTransactorSession) LogRegisteredPubkey(_signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _DepositLog.Contract.LogRegisteredPubkey(&_DepositLog.TransactOpts, _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_DepositLog *DepositLogTransactor) LogSetupFailed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logSetupFailed")
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_DepositLog *DepositLogSession) LogSetupFailed() (*types.Transaction, error) {
	return _DepositLog.Contract.LogSetupFailed(&_DepositLog.TransactOpts)
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_DepositLog *DepositLogTransactorSession) LogSetupFailed() (*types.Transaction, error) {
	return _DepositLog.Contract.LogSetupFailed(&_DepositLog.TransactOpts)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_DepositLog *DepositLogTransactor) LogStartedLiquidation(opts *bind.TransactOpts, _wasFraud bool) (*types.Transaction, error) {
	return _DepositLog.contract.Transact(opts, "logStartedLiquidation", _wasFraud)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_DepositLog *DepositLogSession) LogStartedLiquidation(_wasFraud bool) (*types.Transaction, error) {
	return _DepositLog.Contract.LogStartedLiquidation(&_DepositLog.TransactOpts, _wasFraud)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_DepositLog *DepositLogTransactorSession) LogStartedLiquidation(_wasFraud bool) (*types.Transaction, error) {
	return _DepositLog.Contract.LogStartedLiquidation(&_DepositLog.TransactOpts, _wasFraud)
}

// DepositLogCourtesyCalledIterator is returned from FilterCourtesyCalled and is used to iterate over the raw logs and unpacked data for CourtesyCalled events raised by the DepositLog contract.
type DepositLogCourtesyCalledIterator struct {
	Event *DepositLogCourtesyCalled // Event containing the contract specifics and raw log

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
func (it *DepositLogCourtesyCalledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogCourtesyCalled)
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
		it.Event = new(DepositLogCourtesyCalled)
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
func (it *DepositLogCourtesyCalledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogCourtesyCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogCourtesyCalled represents a CourtesyCalled event raised by the DepositLog contract.
type DepositLogCourtesyCalled struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterCourtesyCalled is a free log retrieval operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterCourtesyCalled(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogCourtesyCalledIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "CourtesyCalled", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogCourtesyCalledIterator{contract: _DepositLog.contract, event: "CourtesyCalled", logs: logs, sub: sub}, nil
}

// WatchCourtesyCalled is a free log subscription operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchCourtesyCalled(opts *bind.WatchOpts, sink chan<- *DepositLogCourtesyCalled, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "CourtesyCalled", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogCourtesyCalled)
				if err := _DepositLog.contract.UnpackLog(event, "CourtesyCalled", log); err != nil {
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

// ParseCourtesyCalled is a log parse operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseCourtesyCalled(log types.Log) (*DepositLogCourtesyCalled, error) {
	event := new(DepositLogCourtesyCalled)
	if err := _DepositLog.contract.UnpackLog(event, "CourtesyCalled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogCreatedIterator is returned from FilterCreated and is used to iterate over the raw logs and unpacked data for Created events raised by the DepositLog contract.
type DepositLogCreatedIterator struct {
	Event *DepositLogCreated // Event containing the contract specifics and raw log

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
func (it *DepositLogCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogCreated)
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
		it.Event = new(DepositLogCreated)
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
func (it *DepositLogCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogCreated represents a Created event raised by the DepositLog contract.
type DepositLogCreated struct {
	DepositContractAddress common.Address
	KeepAddress            common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterCreated is a free log retrieval operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterCreated(opts *bind.FilterOpts, _depositContractAddress []common.Address, _keepAddress []common.Address) (*DepositLogCreatedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _keepAddressRule []interface{}
	for _, _keepAddressItem := range _keepAddress {
		_keepAddressRule = append(_keepAddressRule, _keepAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "Created", _depositContractAddressRule, _keepAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogCreatedIterator{contract: _DepositLog.contract, event: "Created", logs: logs, sub: sub}, nil
}

// WatchCreated is a free log subscription operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchCreated(opts *bind.WatchOpts, sink chan<- *DepositLogCreated, _depositContractAddress []common.Address, _keepAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _keepAddressRule []interface{}
	for _, _keepAddressItem := range _keepAddress {
		_keepAddressRule = append(_keepAddressRule, _keepAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "Created", _depositContractAddressRule, _keepAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogCreated)
				if err := _DepositLog.contract.UnpackLog(event, "Created", log); err != nil {
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

// ParseCreated is a log parse operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseCreated(log types.Log) (*DepositLogCreated, error) {
	event := new(DepositLogCreated)
	if err := _DepositLog.contract.UnpackLog(event, "Created", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogExitedCourtesyCallIterator is returned from FilterExitedCourtesyCall and is used to iterate over the raw logs and unpacked data for ExitedCourtesyCall events raised by the DepositLog contract.
type DepositLogExitedCourtesyCallIterator struct {
	Event *DepositLogExitedCourtesyCall // Event containing the contract specifics and raw log

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
func (it *DepositLogExitedCourtesyCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogExitedCourtesyCall)
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
		it.Event = new(DepositLogExitedCourtesyCall)
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
func (it *DepositLogExitedCourtesyCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogExitedCourtesyCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogExitedCourtesyCall represents a ExitedCourtesyCall event raised by the DepositLog contract.
type DepositLogExitedCourtesyCall struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterExitedCourtesyCall is a free log retrieval operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterExitedCourtesyCall(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogExitedCourtesyCallIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "ExitedCourtesyCall", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogExitedCourtesyCallIterator{contract: _DepositLog.contract, event: "ExitedCourtesyCall", logs: logs, sub: sub}, nil
}

// WatchExitedCourtesyCall is a free log subscription operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchExitedCourtesyCall(opts *bind.WatchOpts, sink chan<- *DepositLogExitedCourtesyCall, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "ExitedCourtesyCall", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogExitedCourtesyCall)
				if err := _DepositLog.contract.UnpackLog(event, "ExitedCourtesyCall", log); err != nil {
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

// ParseExitedCourtesyCall is a log parse operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseExitedCourtesyCall(log types.Log) (*DepositLogExitedCourtesyCall, error) {
	event := new(DepositLogExitedCourtesyCall)
	if err := _DepositLog.contract.UnpackLog(event, "ExitedCourtesyCall", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogFraudDuringSetupIterator is returned from FilterFraudDuringSetup and is used to iterate over the raw logs and unpacked data for FraudDuringSetup events raised by the DepositLog contract.
type DepositLogFraudDuringSetupIterator struct {
	Event *DepositLogFraudDuringSetup // Event containing the contract specifics and raw log

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
func (it *DepositLogFraudDuringSetupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogFraudDuringSetup)
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
		it.Event = new(DepositLogFraudDuringSetup)
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
func (it *DepositLogFraudDuringSetupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogFraudDuringSetupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogFraudDuringSetup represents a FraudDuringSetup event raised by the DepositLog contract.
type DepositLogFraudDuringSetup struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFraudDuringSetup is a free log retrieval operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterFraudDuringSetup(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogFraudDuringSetupIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "FraudDuringSetup", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogFraudDuringSetupIterator{contract: _DepositLog.contract, event: "FraudDuringSetup", logs: logs, sub: sub}, nil
}

// WatchFraudDuringSetup is a free log subscription operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchFraudDuringSetup(opts *bind.WatchOpts, sink chan<- *DepositLogFraudDuringSetup, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "FraudDuringSetup", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogFraudDuringSetup)
				if err := _DepositLog.contract.UnpackLog(event, "FraudDuringSetup", log); err != nil {
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

// ParseFraudDuringSetup is a log parse operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseFraudDuringSetup(log types.Log) (*DepositLogFraudDuringSetup, error) {
	event := new(DepositLogFraudDuringSetup)
	if err := _DepositLog.contract.UnpackLog(event, "FraudDuringSetup", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogFundedIterator is returned from FilterFunded and is used to iterate over the raw logs and unpacked data for Funded events raised by the DepositLog contract.
type DepositLogFundedIterator struct {
	Event *DepositLogFunded // Event containing the contract specifics and raw log

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
func (it *DepositLogFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogFunded)
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
		it.Event = new(DepositLogFunded)
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
func (it *DepositLogFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogFunded represents a Funded event raised by the DepositLog contract.
type DepositLogFunded struct {
	DepositContractAddress common.Address
	Txid                   [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFunded is a free log retrieval operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterFunded(opts *bind.FilterOpts, _depositContractAddress []common.Address, _txid [][32]byte) (*DepositLogFundedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "Funded", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogFundedIterator{contract: _DepositLog.contract, event: "Funded", logs: logs, sub: sub}, nil
}

// WatchFunded is a free log subscription operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *DepositLogFunded, _depositContractAddress []common.Address, _txid [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "Funded", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogFunded)
				if err := _DepositLog.contract.UnpackLog(event, "Funded", log); err != nil {
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

// ParseFunded is a log parse operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseFunded(log types.Log) (*DepositLogFunded, error) {
	event := new(DepositLogFunded)
	if err := _DepositLog.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogFunderAbortRequestedIterator is returned from FilterFunderAbortRequested and is used to iterate over the raw logs and unpacked data for FunderAbortRequested events raised by the DepositLog contract.
type DepositLogFunderAbortRequestedIterator struct {
	Event *DepositLogFunderAbortRequested // Event containing the contract specifics and raw log

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
func (it *DepositLogFunderAbortRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogFunderAbortRequested)
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
		it.Event = new(DepositLogFunderAbortRequested)
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
func (it *DepositLogFunderAbortRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogFunderAbortRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogFunderAbortRequested represents a FunderAbortRequested event raised by the DepositLog contract.
type DepositLogFunderAbortRequested struct {
	DepositContractAddress common.Address
	AbortOutputScript      []byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFunderAbortRequested is a free log retrieval operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_DepositLog *DepositLogFilterer) FilterFunderAbortRequested(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogFunderAbortRequestedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "FunderAbortRequested", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogFunderAbortRequestedIterator{contract: _DepositLog.contract, event: "FunderAbortRequested", logs: logs, sub: sub}, nil
}

// WatchFunderAbortRequested is a free log subscription operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_DepositLog *DepositLogFilterer) WatchFunderAbortRequested(opts *bind.WatchOpts, sink chan<- *DepositLogFunderAbortRequested, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "FunderAbortRequested", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogFunderAbortRequested)
				if err := _DepositLog.contract.UnpackLog(event, "FunderAbortRequested", log); err != nil {
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

// ParseFunderAbortRequested is a log parse operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_DepositLog *DepositLogFilterer) ParseFunderAbortRequested(log types.Log) (*DepositLogFunderAbortRequested, error) {
	event := new(DepositLogFunderAbortRequested)
	if err := _DepositLog.contract.UnpackLog(event, "FunderAbortRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogGotRedemptionSignatureIterator is returned from FilterGotRedemptionSignature and is used to iterate over the raw logs and unpacked data for GotRedemptionSignature events raised by the DepositLog contract.
type DepositLogGotRedemptionSignatureIterator struct {
	Event *DepositLogGotRedemptionSignature // Event containing the contract specifics and raw log

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
func (it *DepositLogGotRedemptionSignatureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogGotRedemptionSignature)
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
		it.Event = new(DepositLogGotRedemptionSignature)
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
func (it *DepositLogGotRedemptionSignatureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogGotRedemptionSignatureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogGotRedemptionSignature represents a GotRedemptionSignature event raised by the DepositLog contract.
type DepositLogGotRedemptionSignature struct {
	DepositContractAddress common.Address
	Digest                 [32]byte
	R                      [32]byte
	S                      [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterGotRedemptionSignature is a free log retrieval operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterGotRedemptionSignature(opts *bind.FilterOpts, _depositContractAddress []common.Address, _digest [][32]byte) (*DepositLogGotRedemptionSignatureIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "GotRedemptionSignature", _depositContractAddressRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogGotRedemptionSignatureIterator{contract: _DepositLog.contract, event: "GotRedemptionSignature", logs: logs, sub: sub}, nil
}

// WatchGotRedemptionSignature is a free log subscription operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchGotRedemptionSignature(opts *bind.WatchOpts, sink chan<- *DepositLogGotRedemptionSignature, _depositContractAddress []common.Address, _digest [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "GotRedemptionSignature", _depositContractAddressRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogGotRedemptionSignature)
				if err := _DepositLog.contract.UnpackLog(event, "GotRedemptionSignature", log); err != nil {
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

// ParseGotRedemptionSignature is a log parse operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseGotRedemptionSignature(log types.Log) (*DepositLogGotRedemptionSignature, error) {
	event := new(DepositLogGotRedemptionSignature)
	if err := _DepositLog.contract.UnpackLog(event, "GotRedemptionSignature", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogLiquidatedIterator is returned from FilterLiquidated and is used to iterate over the raw logs and unpacked data for Liquidated events raised by the DepositLog contract.
type DepositLogLiquidatedIterator struct {
	Event *DepositLogLiquidated // Event containing the contract specifics and raw log

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
func (it *DepositLogLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogLiquidated)
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
		it.Event = new(DepositLogLiquidated)
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
func (it *DepositLogLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogLiquidated represents a Liquidated event raised by the DepositLog contract.
type DepositLogLiquidated struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLiquidated is a free log retrieval operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterLiquidated(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogLiquidatedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "Liquidated", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogLiquidatedIterator{contract: _DepositLog.contract, event: "Liquidated", logs: logs, sub: sub}, nil
}

// WatchLiquidated is a free log subscription operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchLiquidated(opts *bind.WatchOpts, sink chan<- *DepositLogLiquidated, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "Liquidated", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogLiquidated)
				if err := _DepositLog.contract.UnpackLog(event, "Liquidated", log); err != nil {
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

// ParseLiquidated is a log parse operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseLiquidated(log types.Log) (*DepositLogLiquidated, error) {
	event := new(DepositLogLiquidated)
	if err := _DepositLog.contract.UnpackLog(event, "Liquidated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the DepositLog contract.
type DepositLogRedeemedIterator struct {
	Event *DepositLogRedeemed // Event containing the contract specifics and raw log

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
func (it *DepositLogRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogRedeemed)
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
		it.Event = new(DepositLogRedeemed)
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
func (it *DepositLogRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogRedeemed represents a Redeemed event raised by the DepositLog contract.
type DepositLogRedeemed struct {
	DepositContractAddress common.Address
	Txid                   [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterRedeemed(opts *bind.FilterOpts, _depositContractAddress []common.Address, _txid [][32]byte) (*DepositLogRedeemedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "Redeemed", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogRedeemedIterator{contract: _DepositLog.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *DepositLogRedeemed, _depositContractAddress []common.Address, _txid [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "Redeemed", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogRedeemed)
				if err := _DepositLog.contract.UnpackLog(event, "Redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseRedeemed(log types.Log) (*DepositLogRedeemed, error) {
	event := new(DepositLogRedeemed)
	if err := _DepositLog.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogRedemptionRequestedIterator is returned from FilterRedemptionRequested and is used to iterate over the raw logs and unpacked data for RedemptionRequested events raised by the DepositLog contract.
type DepositLogRedemptionRequestedIterator struct {
	Event *DepositLogRedemptionRequested // Event containing the contract specifics and raw log

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
func (it *DepositLogRedemptionRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogRedemptionRequested)
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
		it.Event = new(DepositLogRedemptionRequested)
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
func (it *DepositLogRedemptionRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogRedemptionRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogRedemptionRequested represents a RedemptionRequested event raised by the DepositLog contract.
type DepositLogRedemptionRequested struct {
	DepositContractAddress common.Address
	Requester              common.Address
	Digest                 [32]byte
	UtxoValue              *big.Int
	RedeemerOutputScript   []byte
	RequestedFee           *big.Int
	Outpoint               []byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRedemptionRequested is a free log retrieval operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_DepositLog *DepositLogFilterer) FilterRedemptionRequested(opts *bind.FilterOpts, _depositContractAddress []common.Address, _requester []common.Address, _digest [][32]byte) (*DepositLogRedemptionRequestedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "RedemptionRequested", _depositContractAddressRule, _requesterRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogRedemptionRequestedIterator{contract: _DepositLog.contract, event: "RedemptionRequested", logs: logs, sub: sub}, nil
}

// WatchRedemptionRequested is a free log subscription operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_DepositLog *DepositLogFilterer) WatchRedemptionRequested(opts *bind.WatchOpts, sink chan<- *DepositLogRedemptionRequested, _depositContractAddress []common.Address, _requester []common.Address, _digest [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "RedemptionRequested", _depositContractAddressRule, _requesterRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogRedemptionRequested)
				if err := _DepositLog.contract.UnpackLog(event, "RedemptionRequested", log); err != nil {
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

// ParseRedemptionRequested is a log parse operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_DepositLog *DepositLogFilterer) ParseRedemptionRequested(log types.Log) (*DepositLogRedemptionRequested, error) {
	event := new(DepositLogRedemptionRequested)
	if err := _DepositLog.contract.UnpackLog(event, "RedemptionRequested", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogRegisteredPubkeyIterator is returned from FilterRegisteredPubkey and is used to iterate over the raw logs and unpacked data for RegisteredPubkey events raised by the DepositLog contract.
type DepositLogRegisteredPubkeyIterator struct {
	Event *DepositLogRegisteredPubkey // Event containing the contract specifics and raw log

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
func (it *DepositLogRegisteredPubkeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogRegisteredPubkey)
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
		it.Event = new(DepositLogRegisteredPubkey)
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
func (it *DepositLogRegisteredPubkeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogRegisteredPubkeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogRegisteredPubkey represents a RegisteredPubkey event raised by the DepositLog contract.
type DepositLogRegisteredPubkey struct {
	DepositContractAddress common.Address
	SigningGroupPubkeyX    [32]byte
	SigningGroupPubkeyY    [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRegisteredPubkey is a free log retrieval operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterRegisteredPubkey(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogRegisteredPubkeyIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "RegisteredPubkey", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogRegisteredPubkeyIterator{contract: _DepositLog.contract, event: "RegisteredPubkey", logs: logs, sub: sub}, nil
}

// WatchRegisteredPubkey is a free log subscription operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchRegisteredPubkey(opts *bind.WatchOpts, sink chan<- *DepositLogRegisteredPubkey, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "RegisteredPubkey", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogRegisteredPubkey)
				if err := _DepositLog.contract.UnpackLog(event, "RegisteredPubkey", log); err != nil {
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

// ParseRegisteredPubkey is a log parse operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseRegisteredPubkey(log types.Log) (*DepositLogRegisteredPubkey, error) {
	event := new(DepositLogRegisteredPubkey)
	if err := _DepositLog.contract.UnpackLog(event, "RegisteredPubkey", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogSetupFailedIterator is returned from FilterSetupFailed and is used to iterate over the raw logs and unpacked data for SetupFailed events raised by the DepositLog contract.
type DepositLogSetupFailedIterator struct {
	Event *DepositLogSetupFailed // Event containing the contract specifics and raw log

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
func (it *DepositLogSetupFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogSetupFailed)
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
		it.Event = new(DepositLogSetupFailed)
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
func (it *DepositLogSetupFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogSetupFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogSetupFailed represents a SetupFailed event raised by the DepositLog contract.
type DepositLogSetupFailed struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetupFailed is a free log retrieval operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterSetupFailed(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogSetupFailedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "SetupFailed", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogSetupFailedIterator{contract: _DepositLog.contract, event: "SetupFailed", logs: logs, sub: sub}, nil
}

// WatchSetupFailed is a free log subscription operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchSetupFailed(opts *bind.WatchOpts, sink chan<- *DepositLogSetupFailed, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "SetupFailed", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogSetupFailed)
				if err := _DepositLog.contract.UnpackLog(event, "SetupFailed", log); err != nil {
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

// ParseSetupFailed is a log parse operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseSetupFailed(log types.Log) (*DepositLogSetupFailed, error) {
	event := new(DepositLogSetupFailed)
	if err := _DepositLog.contract.UnpackLog(event, "SetupFailed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DepositLogStartedLiquidationIterator is returned from FilterStartedLiquidation and is used to iterate over the raw logs and unpacked data for StartedLiquidation events raised by the DepositLog contract.
type DepositLogStartedLiquidationIterator struct {
	Event *DepositLogStartedLiquidation // Event containing the contract specifics and raw log

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
func (it *DepositLogStartedLiquidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositLogStartedLiquidation)
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
		it.Event = new(DepositLogStartedLiquidation)
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
func (it *DepositLogStartedLiquidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositLogStartedLiquidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositLogStartedLiquidation represents a StartedLiquidation event raised by the DepositLog contract.
type DepositLogStartedLiquidation struct {
	DepositContractAddress common.Address
	WasFraud               bool
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterStartedLiquidation is a free log retrieval operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) FilterStartedLiquidation(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*DepositLogStartedLiquidationIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.FilterLogs(opts, "StartedLiquidation", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &DepositLogStartedLiquidationIterator{contract: _DepositLog.contract, event: "StartedLiquidation", logs: logs, sub: sub}, nil
}

// WatchStartedLiquidation is a free log subscription operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) WatchStartedLiquidation(opts *bind.WatchOpts, sink chan<- *DepositLogStartedLiquidation, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _DepositLog.contract.WatchLogs(opts, "StartedLiquidation", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositLogStartedLiquidation)
				if err := _DepositLog.contract.UnpackLog(event, "StartedLiquidation", log); err != nil {
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

// ParseStartedLiquidation is a log parse operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_DepositLog *DepositLogFilterer) ParseStartedLiquidation(log types.Log) (*DepositLogStartedLiquidation, error) {
	event := new(DepositLogStartedLiquidation)
	if err := _DepositLog.contract.UnpackLog(event, "StartedLiquidation", log); err != nil {
		return nil, err
	}
	return event, nil
}
