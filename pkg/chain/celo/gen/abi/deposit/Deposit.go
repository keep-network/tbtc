// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/event"
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

// DepositABI is the input ABI used to generate the binding from.
const DepositABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"auctionValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"collateralizationPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentState\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"exitCourtesyCall\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fundingInfo\",\"outputs\":[{\"internalType\":\"bytes8\",\"name\":\"utxoValueBytes\",\"type\":\"bytes8\"},{\"internalType\":\"uint256\",\"name\":\"fundedAt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"utxoOutpoint\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"}],\"name\":\"getOwnerRedemptionTbtcRequirement\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"}],\"name\":\"getRedemptionTbtcRequirement\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"inActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes8\",\"name\":\"_previousOutputValueBytes\",\"type\":\"bytes8\"},{\"internalType\":\"bytes8\",\"name\":\"_newOutputValueBytes\",\"type\":\"bytes8\"}],\"name\":\"increaseRedemptionFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"initialCollateralizedPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_factory\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractITBTCSystem\",\"name\":\"_tbtcSystem\",\"type\":\"address\"},{\"internalType\":\"contractTBTCToken\",\"name\":\"_tbtcToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC721\",\"name\":\"_tbtcDepositToken\",\"type\":\"address\"},{\"internalType\":\"contractFeeRebateToken\",\"name\":\"_feeRebateToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vendingMachineAddress\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_lotSizeSatoshis\",\"type\":\"uint64\"}],\"name\":\"initializeDeposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keepAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lotSizeSatoshis\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lotSizeTbtc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyCourtesyCall\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyCourtesyCallExpired\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyFundingTimedOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyRedemptionProofTimedOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyRedemptionSignatureTimedOut\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifySignerSetupFailed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"notifyUndercollateralizedLiquidation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_txVersion\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"_txInputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_txOutputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"_txLocktime\",\"type\":\"bytes4\"},{\"internalType\":\"uint8\",\"name\":\"_fundingOutputIndex\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"_merkleProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_txIndexInBlock\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_bitcoinHeaders\",\"type\":\"bytes\"}],\"name\":\"provideBTCFundingProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_signedDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_preimage\",\"type\":\"bytes\"}],\"name\":\"provideECDSAFraudProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_signedDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_preimage\",\"type\":\"bytes\"}],\"name\":\"provideFundingECDSAFraudProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"_txVersion\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"_txInputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_txOutputVector\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"_txLocktime\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"_merkleProof\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_txIndexInBlock\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_bitcoinHeaders\",\"type\":\"bytes\"}],\"name\":\"provideRedemptionProof\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"provideRedemptionSignature\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"purchaseSignerBondsAtAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"remainingTerm\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_abortOutputScript\",\"type\":\"bytes\"}],\"name\":\"requestFunderAbort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes8\",\"name\":\"_outputValueBytes\",\"type\":\"bytes8\"},{\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"}],\"name\":\"requestRedemption\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"retrieveSignerPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"severelyUndercollateralizedThresholdPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"signerFeeTbtc\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes8\",\"name\":\"_outputValueBytes\",\"type\":\"bytes8\"},{\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"},{\"internalType\":\"addresspayable\",\"name\":\"_finalRecipient\",\"type\":\"address\"}],\"name\":\"transferAndRequestRedemption\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"undercollateralizedThresholdPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"utxoValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawFunds\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawableAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Deposit is an auto generated Go binding around an Ethereum contract.
type Deposit struct {
	DepositCaller     // Read-only binding to the contract
	DepositTransactor // Write-only binding to the contract
	DepositFilterer   // Log filterer for contract events
}

// DepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSession struct {
	Contract     *Deposit          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCallerSession struct {
	Contract *DepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// DepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositTransactorSession struct {
	Contract     *DepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositRaw struct {
	Contract *Deposit // Generic contract binding to access the raw methods on
}

// DepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCallerRaw struct {
	Contract *DepositCaller // Generic read-only contract binding to access the raw methods on
}

// DepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositTransactorRaw struct {
	Contract *DepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeposit creates a new instance of Deposit, bound to a specific deployed contract.
func NewDeposit(address common.Address, backend bind.ContractBackend) (*Deposit, error) {
	contract, err := bindDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// NewDepositCaller creates a new read-only instance of Deposit, bound to a specific deployed contract.
func NewDepositCaller(address common.Address, caller bind.ContractCaller) (*DepositCaller, error) {
	contract, err := bindDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCaller{contract: contract}, nil
}

// NewDepositTransactor creates a new write-only instance of Deposit, bound to a specific deployed contract.
func NewDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositTransactor, error) {
	contract, err := bindDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositTransactor{contract: contract}, nil
}

// NewDepositFilterer creates a new log filterer instance of Deposit, bound to a specific deployed contract.
func NewDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositFilterer, error) {
	contract, err := bindDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositFilterer{contract: contract}, nil
}

// bindDeposit binds a generic wrapper to an already deployed contract.
func bindDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ParseDepositABI parses the ABI
func ParseDepositABI() (*abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.DepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transact(opts, method, params...)
}

// AuctionValue is a free data retrieval call binding the contract method 0x13f654df.
//
// Solidity: function auctionValue() view returns(uint256)
func (_Deposit *DepositCaller) AuctionValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "auctionValue")
	return *ret0, err
}

// AuctionValue is a free data retrieval call binding the contract method 0x13f654df.
//
// Solidity: function auctionValue() view returns(uint256)
func (_Deposit *DepositSession) AuctionValue() (*big.Int, error) {
	return _Deposit.Contract.AuctionValue(&_Deposit.CallOpts)
}

// AuctionValue is a free data retrieval call binding the contract method 0x13f654df.
//
// Solidity: function auctionValue() view returns(uint256)
func (_Deposit *DepositCallerSession) AuctionValue() (*big.Int, error) {
	return _Deposit.Contract.AuctionValue(&_Deposit.CallOpts)
}

// CollateralizationPercentage is a free data retrieval call binding the contract method 0x6e4668be.
//
// Solidity: function collateralizationPercentage() view returns(uint256)
func (_Deposit *DepositCaller) CollateralizationPercentage(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "collateralizationPercentage")
	return *ret0, err
}

// CollateralizationPercentage is a free data retrieval call binding the contract method 0x6e4668be.
//
// Solidity: function collateralizationPercentage() view returns(uint256)
func (_Deposit *DepositSession) CollateralizationPercentage() (*big.Int, error) {
	return _Deposit.Contract.CollateralizationPercentage(&_Deposit.CallOpts)
}

// CollateralizationPercentage is a free data retrieval call binding the contract method 0x6e4668be.
//
// Solidity: function collateralizationPercentage() view returns(uint256)
func (_Deposit *DepositCallerSession) CollateralizationPercentage() (*big.Int, error) {
	return _Deposit.Contract.CollateralizationPercentage(&_Deposit.CallOpts)
}

// CurrentState is a free data retrieval call binding the contract method 0x0c3f6acf.
//
// Solidity: function currentState() view returns(uint256)
func (_Deposit *DepositCaller) CurrentState(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "currentState")
	return *ret0, err
}

// CurrentState is a free data retrieval call binding the contract method 0x0c3f6acf.
//
// Solidity: function currentState() view returns(uint256)
func (_Deposit *DepositSession) CurrentState() (*big.Int, error) {
	return _Deposit.Contract.CurrentState(&_Deposit.CallOpts)
}

// CurrentState is a free data retrieval call binding the contract method 0x0c3f6acf.
//
// Solidity: function currentState() view returns(uint256)
func (_Deposit *DepositCallerSession) CurrentState() (*big.Int, error) {
	return _Deposit.Contract.CurrentState(&_Deposit.CallOpts)
}

// FundingInfo is a free data retrieval call binding the contract method 0xdba49153.
//
// Solidity: function fundingInfo() view returns(bytes8 utxoValueBytes, uint256 fundedAt, bytes utxoOutpoint)
func (_Deposit *DepositCaller) FundingInfo(opts *bind.CallOpts) (struct {
	UtxoValueBytes [8]byte
	FundedAt       *big.Int
	UtxoOutpoint   []byte
}, error) {
	ret := new(struct {
		UtxoValueBytes [8]byte
		FundedAt       *big.Int
		UtxoOutpoint   []byte
	})
	out := ret
	err := _Deposit.contract.Call(opts, out, "fundingInfo")
	return *ret, err
}

// FundingInfo is a free data retrieval call binding the contract method 0xdba49153.
//
// Solidity: function fundingInfo() view returns(bytes8 utxoValueBytes, uint256 fundedAt, bytes utxoOutpoint)
func (_Deposit *DepositSession) FundingInfo() (struct {
	UtxoValueBytes [8]byte
	FundedAt       *big.Int
	UtxoOutpoint   []byte
}, error) {
	return _Deposit.Contract.FundingInfo(&_Deposit.CallOpts)
}

// FundingInfo is a free data retrieval call binding the contract method 0xdba49153.
//
// Solidity: function fundingInfo() view returns(bytes8 utxoValueBytes, uint256 fundedAt, bytes utxoOutpoint)
func (_Deposit *DepositCallerSession) FundingInfo() (struct {
	UtxoValueBytes [8]byte
	FundedAt       *big.Int
	UtxoOutpoint   []byte
}, error) {
	return _Deposit.Contract.FundingInfo(&_Deposit.CallOpts)
}

// GetOwnerRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd8d02330.
//
// Solidity: function getOwnerRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositCaller) GetOwnerRedemptionTbtcRequirement(opts *bind.CallOpts, _redeemer common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "getOwnerRedemptionTbtcRequirement", _redeemer)
	return *ret0, err
}

// GetOwnerRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd8d02330.
//
// Solidity: function getOwnerRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositSession) GetOwnerRedemptionTbtcRequirement(_redeemer common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetOwnerRedemptionTbtcRequirement(&_Deposit.CallOpts, _redeemer)
}

// GetOwnerRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd8d02330.
//
// Solidity: function getOwnerRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositCallerSession) GetOwnerRedemptionTbtcRequirement(_redeemer common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetOwnerRedemptionTbtcRequirement(&_Deposit.CallOpts, _redeemer)
}

// GetRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd02fd958.
//
// Solidity: function getRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositCaller) GetRedemptionTbtcRequirement(opts *bind.CallOpts, _redeemer common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "getRedemptionTbtcRequirement", _redeemer)
	return *ret0, err
}

// GetRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd02fd958.
//
// Solidity: function getRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositSession) GetRedemptionTbtcRequirement(_redeemer common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetRedemptionTbtcRequirement(&_Deposit.CallOpts, _redeemer)
}

// GetRedemptionTbtcRequirement is a free data retrieval call binding the contract method 0xd02fd958.
//
// Solidity: function getRedemptionTbtcRequirement(address _redeemer) view returns(uint256)
func (_Deposit *DepositCallerSession) GetRedemptionTbtcRequirement(_redeemer common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetRedemptionTbtcRequirement(&_Deposit.CallOpts, _redeemer)
}

// InActive is a free data retrieval call binding the contract method 0xf97a02fa.
//
// Solidity: function inActive() view returns(bool)
func (_Deposit *DepositCaller) InActive(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "inActive")
	return *ret0, err
}

// InActive is a free data retrieval call binding the contract method 0xf97a02fa.
//
// Solidity: function inActive() view returns(bool)
func (_Deposit *DepositSession) InActive() (bool, error) {
	return _Deposit.Contract.InActive(&_Deposit.CallOpts)
}

// InActive is a free data retrieval call binding the contract method 0xf97a02fa.
//
// Solidity: function inActive() view returns(bool)
func (_Deposit *DepositCallerSession) InActive() (bool, error) {
	return _Deposit.Contract.InActive(&_Deposit.CallOpts)
}

// InitialCollateralizedPercent is a free data retrieval call binding the contract method 0x76ef5510.
//
// Solidity: function initialCollateralizedPercent() view returns(uint16)
func (_Deposit *DepositCaller) InitialCollateralizedPercent(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "initialCollateralizedPercent")
	return *ret0, err
}

// InitialCollateralizedPercent is a free data retrieval call binding the contract method 0x76ef5510.
//
// Solidity: function initialCollateralizedPercent() view returns(uint16)
func (_Deposit *DepositSession) InitialCollateralizedPercent() (uint16, error) {
	return _Deposit.Contract.InitialCollateralizedPercent(&_Deposit.CallOpts)
}

// InitialCollateralizedPercent is a free data retrieval call binding the contract method 0x76ef5510.
//
// Solidity: function initialCollateralizedPercent() view returns(uint16)
func (_Deposit *DepositCallerSession) InitialCollateralizedPercent() (uint16, error) {
	return _Deposit.Contract.InitialCollateralizedPercent(&_Deposit.CallOpts)
}

// KeepAddress is a free data retrieval call binding the contract method 0x6c3b0114.
//
// Solidity: function keepAddress() view returns(address)
func (_Deposit *DepositCaller) KeepAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "keepAddress")
	return *ret0, err
}

// KeepAddress is a free data retrieval call binding the contract method 0x6c3b0114.
//
// Solidity: function keepAddress() view returns(address)
func (_Deposit *DepositSession) KeepAddress() (common.Address, error) {
	return _Deposit.Contract.KeepAddress(&_Deposit.CallOpts)
}

// KeepAddress is a free data retrieval call binding the contract method 0x6c3b0114.
//
// Solidity: function keepAddress() view returns(address)
func (_Deposit *DepositCallerSession) KeepAddress() (common.Address, error) {
	return _Deposit.Contract.KeepAddress(&_Deposit.CallOpts)
}

// LotSizeSatoshis is a free data retrieval call binding the contract method 0x90a2f687.
//
// Solidity: function lotSizeSatoshis() view returns(uint64)
func (_Deposit *DepositCaller) LotSizeSatoshis(opts *bind.CallOpts) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "lotSizeSatoshis")
	return *ret0, err
}

// LotSizeSatoshis is a free data retrieval call binding the contract method 0x90a2f687.
//
// Solidity: function lotSizeSatoshis() view returns(uint64)
func (_Deposit *DepositSession) LotSizeSatoshis() (uint64, error) {
	return _Deposit.Contract.LotSizeSatoshis(&_Deposit.CallOpts)
}

// LotSizeSatoshis is a free data retrieval call binding the contract method 0x90a2f687.
//
// Solidity: function lotSizeSatoshis() view returns(uint64)
func (_Deposit *DepositCallerSession) LotSizeSatoshis() (uint64, error) {
	return _Deposit.Contract.LotSizeSatoshis(&_Deposit.CallOpts)
}

// LotSizeTbtc is a free data retrieval call binding the contract method 0x946fbf4c.
//
// Solidity: function lotSizeTbtc() view returns(uint256)
func (_Deposit *DepositCaller) LotSizeTbtc(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "lotSizeTbtc")
	return *ret0, err
}

// LotSizeTbtc is a free data retrieval call binding the contract method 0x946fbf4c.
//
// Solidity: function lotSizeTbtc() view returns(uint256)
func (_Deposit *DepositSession) LotSizeTbtc() (*big.Int, error) {
	return _Deposit.Contract.LotSizeTbtc(&_Deposit.CallOpts)
}

// LotSizeTbtc is a free data retrieval call binding the contract method 0x946fbf4c.
//
// Solidity: function lotSizeTbtc() view returns(uint256)
func (_Deposit *DepositCallerSession) LotSizeTbtc() (*big.Int, error) {
	return _Deposit.Contract.LotSizeTbtc(&_Deposit.CallOpts)
}

// RemainingTerm is a free data retrieval call binding the contract method 0x35bc0ebe.
//
// Solidity: function remainingTerm() view returns(uint256)
func (_Deposit *DepositCaller) RemainingTerm(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "remainingTerm")
	return *ret0, err
}

// RemainingTerm is a free data retrieval call binding the contract method 0x35bc0ebe.
//
// Solidity: function remainingTerm() view returns(uint256)
func (_Deposit *DepositSession) RemainingTerm() (*big.Int, error) {
	return _Deposit.Contract.RemainingTerm(&_Deposit.CallOpts)
}

// RemainingTerm is a free data retrieval call binding the contract method 0x35bc0ebe.
//
// Solidity: function remainingTerm() view returns(uint256)
func (_Deposit *DepositCallerSession) RemainingTerm() (*big.Int, error) {
	return _Deposit.Contract.RemainingTerm(&_Deposit.CallOpts)
}

// SeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x0d5889f4.
//
// Solidity: function severelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositCaller) SeverelyUndercollateralizedThresholdPercent(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "severelyUndercollateralizedThresholdPercent")
	return *ret0, err
}

// SeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x0d5889f4.
//
// Solidity: function severelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositSession) SeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	return _Deposit.Contract.SeverelyUndercollateralizedThresholdPercent(&_Deposit.CallOpts)
}

// SeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x0d5889f4.
//
// Solidity: function severelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositCallerSession) SeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	return _Deposit.Contract.SeverelyUndercollateralizedThresholdPercent(&_Deposit.CallOpts)
}

// SignerFeeTbtc is a free data retrieval call binding the contract method 0x058d3703.
//
// Solidity: function signerFeeTbtc() view returns(uint256)
func (_Deposit *DepositCaller) SignerFeeTbtc(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "signerFeeTbtc")
	return *ret0, err
}

// SignerFeeTbtc is a free data retrieval call binding the contract method 0x058d3703.
//
// Solidity: function signerFeeTbtc() view returns(uint256)
func (_Deposit *DepositSession) SignerFeeTbtc() (*big.Int, error) {
	return _Deposit.Contract.SignerFeeTbtc(&_Deposit.CallOpts)
}

// SignerFeeTbtc is a free data retrieval call binding the contract method 0x058d3703.
//
// Solidity: function signerFeeTbtc() view returns(uint256)
func (_Deposit *DepositCallerSession) SignerFeeTbtc() (*big.Int, error) {
	return _Deposit.Contract.SignerFeeTbtc(&_Deposit.CallOpts)
}

// UndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x85df153d.
//
// Solidity: function undercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositCaller) UndercollateralizedThresholdPercent(opts *bind.CallOpts) (uint16, error) {
	var (
		ret0 = new(uint16)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "undercollateralizedThresholdPercent")
	return *ret0, err
}

// UndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x85df153d.
//
// Solidity: function undercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositSession) UndercollateralizedThresholdPercent() (uint16, error) {
	return _Deposit.Contract.UndercollateralizedThresholdPercent(&_Deposit.CallOpts)
}

// UndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x85df153d.
//
// Solidity: function undercollateralizedThresholdPercent() view returns(uint16)
func (_Deposit *DepositCallerSession) UndercollateralizedThresholdPercent() (uint16, error) {
	return _Deposit.Contract.UndercollateralizedThresholdPercent(&_Deposit.CallOpts)
}

// UtxoValue is a free data retrieval call binding the contract method 0x87a90d80.
//
// Solidity: function utxoValue() view returns(uint256)
func (_Deposit *DepositCaller) UtxoValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "utxoValue")
	return *ret0, err
}

// UtxoValue is a free data retrieval call binding the contract method 0x87a90d80.
//
// Solidity: function utxoValue() view returns(uint256)
func (_Deposit *DepositSession) UtxoValue() (*big.Int, error) {
	return _Deposit.Contract.UtxoValue(&_Deposit.CallOpts)
}

// UtxoValue is a free data retrieval call binding the contract method 0x87a90d80.
//
// Solidity: function utxoValue() view returns(uint256)
func (_Deposit *DepositCallerSession) UtxoValue() (*big.Int, error) {
	return _Deposit.Contract.UtxoValue(&_Deposit.CallOpts)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_Deposit *DepositCaller) WithdrawableAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Deposit.contract.Call(opts, out, "withdrawableAmount")
	return *ret0, err
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_Deposit *DepositSession) WithdrawableAmount() (*big.Int, error) {
	return _Deposit.Contract.WithdrawableAmount(&_Deposit.CallOpts)
}

// WithdrawableAmount is a free data retrieval call binding the contract method 0x951303f5.
//
// Solidity: function withdrawableAmount() view returns(uint256)
func (_Deposit *DepositCallerSession) WithdrawableAmount() (*big.Int, error) {
	return _Deposit.Contract.WithdrawableAmount(&_Deposit.CallOpts)
}

// ExitCourtesyCall is a paid mutator transaction binding the contract method 0x287b32e5.
//
// Solidity: function exitCourtesyCall() returns()
func (_Deposit *DepositTransactor) ExitCourtesyCall(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "exitCourtesyCall")
}

// ExitCourtesyCall is a paid mutator transaction binding the contract method 0x287b32e5.
//
// Solidity: function exitCourtesyCall() returns()
func (_Deposit *DepositSession) ExitCourtesyCall() (*types.Transaction, error) {
	return _Deposit.Contract.ExitCourtesyCall(&_Deposit.TransactOpts)
}

// ExitCourtesyCall is a paid mutator transaction binding the contract method 0x287b32e5.
//
// Solidity: function exitCourtesyCall() returns()
func (_Deposit *DepositTransactorSession) ExitCourtesyCall() (*types.Transaction, error) {
	return _Deposit.Contract.ExitCourtesyCall(&_Deposit.TransactOpts)
}

// IncreaseRedemptionFee is a paid mutator transaction binding the contract method 0x9894d734.
//
// Solidity: function increaseRedemptionFee(bytes8 _previousOutputValueBytes, bytes8 _newOutputValueBytes) returns()
func (_Deposit *DepositTransactor) IncreaseRedemptionFee(opts *bind.TransactOpts, _previousOutputValueBytes [8]byte, _newOutputValueBytes [8]byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "increaseRedemptionFee", _previousOutputValueBytes, _newOutputValueBytes)
}

// IncreaseRedemptionFee is a paid mutator transaction binding the contract method 0x9894d734.
//
// Solidity: function increaseRedemptionFee(bytes8 _previousOutputValueBytes, bytes8 _newOutputValueBytes) returns()
func (_Deposit *DepositSession) IncreaseRedemptionFee(_previousOutputValueBytes [8]byte, _newOutputValueBytes [8]byte) (*types.Transaction, error) {
	return _Deposit.Contract.IncreaseRedemptionFee(&_Deposit.TransactOpts, _previousOutputValueBytes, _newOutputValueBytes)
}

// IncreaseRedemptionFee is a paid mutator transaction binding the contract method 0x9894d734.
//
// Solidity: function increaseRedemptionFee(bytes8 _previousOutputValueBytes, bytes8 _newOutputValueBytes) returns()
func (_Deposit *DepositTransactorSession) IncreaseRedemptionFee(_previousOutputValueBytes [8]byte, _newOutputValueBytes [8]byte) (*types.Transaction, error) {
	return _Deposit.Contract.IncreaseRedemptionFee(&_Deposit.TransactOpts, _previousOutputValueBytes, _newOutputValueBytes)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _factory) returns()
func (_Deposit *DepositTransactor) Initialize(opts *bind.TransactOpts, _factory common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "initialize", _factory)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _factory) returns()
func (_Deposit *DepositSession) Initialize(_factory common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts, _factory)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _factory) returns()
func (_Deposit *DepositTransactorSession) Initialize(_factory common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts, _factory)
}

// InitializeDeposit is a paid mutator transaction binding the contract method 0xa81e63f7.
//
// Solidity: function initializeDeposit(address _tbtcSystem, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachineAddress, uint64 _lotSizeSatoshis) payable returns()
func (_Deposit *DepositTransactor) InitializeDeposit(opts *bind.TransactOpts, _tbtcSystem common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachineAddress common.Address, _lotSizeSatoshis uint64) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "initializeDeposit", _tbtcSystem, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachineAddress, _lotSizeSatoshis)
}

// InitializeDeposit is a paid mutator transaction binding the contract method 0xa81e63f7.
//
// Solidity: function initializeDeposit(address _tbtcSystem, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachineAddress, uint64 _lotSizeSatoshis) payable returns()
func (_Deposit *DepositSession) InitializeDeposit(_tbtcSystem common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachineAddress common.Address, _lotSizeSatoshis uint64) (*types.Transaction, error) {
	return _Deposit.Contract.InitializeDeposit(&_Deposit.TransactOpts, _tbtcSystem, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachineAddress, _lotSizeSatoshis)
}

// InitializeDeposit is a paid mutator transaction binding the contract method 0xa81e63f7.
//
// Solidity: function initializeDeposit(address _tbtcSystem, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachineAddress, uint64 _lotSizeSatoshis) payable returns()
func (_Deposit *DepositTransactorSession) InitializeDeposit(_tbtcSystem common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachineAddress common.Address, _lotSizeSatoshis uint64) (*types.Transaction, error) {
	return _Deposit.Contract.InitializeDeposit(&_Deposit.TransactOpts, _tbtcSystem, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachineAddress, _lotSizeSatoshis)
}

// NotifyCourtesyCall is a paid mutator transaction binding the contract method 0x96aab311.
//
// Solidity: function notifyCourtesyCall() returns()
func (_Deposit *DepositTransactor) NotifyCourtesyCall(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyCourtesyCall")
}

// NotifyCourtesyCall is a paid mutator transaction binding the contract method 0x96aab311.
//
// Solidity: function notifyCourtesyCall() returns()
func (_Deposit *DepositSession) NotifyCourtesyCall() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyCourtesyCall(&_Deposit.TransactOpts)
}

// NotifyCourtesyCall is a paid mutator transaction binding the contract method 0x96aab311.
//
// Solidity: function notifyCourtesyCall() returns()
func (_Deposit *DepositTransactorSession) NotifyCourtesyCall() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyCourtesyCall(&_Deposit.TransactOpts)
}

// NotifyCourtesyCallExpired is a paid mutator transaction binding the contract method 0x91d165e3.
//
// Solidity: function notifyCourtesyCallExpired() returns()
func (_Deposit *DepositTransactor) NotifyCourtesyCallExpired(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyCourtesyCallExpired")
}

// NotifyCourtesyCallExpired is a paid mutator transaction binding the contract method 0x91d165e3.
//
// Solidity: function notifyCourtesyCallExpired() returns()
func (_Deposit *DepositSession) NotifyCourtesyCallExpired() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyCourtesyCallExpired(&_Deposit.TransactOpts)
}

// NotifyCourtesyCallExpired is a paid mutator transaction binding the contract method 0x91d165e3.
//
// Solidity: function notifyCourtesyCallExpired() returns()
func (_Deposit *DepositTransactorSession) NotifyCourtesyCallExpired() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyCourtesyCallExpired(&_Deposit.TransactOpts)
}

// NotifyFundingTimedOut is a paid mutator transaction binding the contract method 0x259b1ea3.
//
// Solidity: function notifyFundingTimedOut() returns()
func (_Deposit *DepositTransactor) NotifyFundingTimedOut(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyFundingTimedOut")
}

// NotifyFundingTimedOut is a paid mutator transaction binding the contract method 0x259b1ea3.
//
// Solidity: function notifyFundingTimedOut() returns()
func (_Deposit *DepositSession) NotifyFundingTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyFundingTimedOut(&_Deposit.TransactOpts)
}

// NotifyFundingTimedOut is a paid mutator transaction binding the contract method 0x259b1ea3.
//
// Solidity: function notifyFundingTimedOut() returns()
func (_Deposit *DepositTransactorSession) NotifyFundingTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyFundingTimedOut(&_Deposit.TransactOpts)
}

// NotifyRedemptionProofTimedOut is a paid mutator transaction binding the contract method 0xba346839.
//
// Solidity: function notifyRedemptionProofTimedOut() returns()
func (_Deposit *DepositTransactor) NotifyRedemptionProofTimedOut(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyRedemptionProofTimedOut")
}

// NotifyRedemptionProofTimedOut is a paid mutator transaction binding the contract method 0xba346839.
//
// Solidity: function notifyRedemptionProofTimedOut() returns()
func (_Deposit *DepositSession) NotifyRedemptionProofTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyRedemptionProofTimedOut(&_Deposit.TransactOpts)
}

// NotifyRedemptionProofTimedOut is a paid mutator transaction binding the contract method 0xba346839.
//
// Solidity: function notifyRedemptionProofTimedOut() returns()
func (_Deposit *DepositTransactorSession) NotifyRedemptionProofTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyRedemptionProofTimedOut(&_Deposit.TransactOpts)
}

// NotifyRedemptionSignatureTimedOut is a paid mutator transaction binding the contract method 0x2b0bc981.
//
// Solidity: function notifyRedemptionSignatureTimedOut() returns()
func (_Deposit *DepositTransactor) NotifyRedemptionSignatureTimedOut(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyRedemptionSignatureTimedOut")
}

// NotifyRedemptionSignatureTimedOut is a paid mutator transaction binding the contract method 0x2b0bc981.
//
// Solidity: function notifyRedemptionSignatureTimedOut() returns()
func (_Deposit *DepositSession) NotifyRedemptionSignatureTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyRedemptionSignatureTimedOut(&_Deposit.TransactOpts)
}

// NotifyRedemptionSignatureTimedOut is a paid mutator transaction binding the contract method 0x2b0bc981.
//
// Solidity: function notifyRedemptionSignatureTimedOut() returns()
func (_Deposit *DepositTransactorSession) NotifyRedemptionSignatureTimedOut() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyRedemptionSignatureTimedOut(&_Deposit.TransactOpts)
}

// NotifySignerSetupFailed is a paid mutator transaction binding the contract method 0x4f706e44.
//
// Solidity: function notifySignerSetupFailed() returns()
func (_Deposit *DepositTransactor) NotifySignerSetupFailed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifySignerSetupFailed")
}

// NotifySignerSetupFailed is a paid mutator transaction binding the contract method 0x4f706e44.
//
// Solidity: function notifySignerSetupFailed() returns()
func (_Deposit *DepositSession) NotifySignerSetupFailed() (*types.Transaction, error) {
	return _Deposit.Contract.NotifySignerSetupFailed(&_Deposit.TransactOpts)
}

// NotifySignerSetupFailed is a paid mutator transaction binding the contract method 0x4f706e44.
//
// Solidity: function notifySignerSetupFailed() returns()
func (_Deposit *DepositTransactorSession) NotifySignerSetupFailed() (*types.Transaction, error) {
	return _Deposit.Contract.NotifySignerSetupFailed(&_Deposit.TransactOpts)
}

// NotifyUndercollateralizedLiquidation is a paid mutator transaction binding the contract method 0xb4bd2e7a.
//
// Solidity: function notifyUndercollateralizedLiquidation() returns()
func (_Deposit *DepositTransactor) NotifyUndercollateralizedLiquidation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "notifyUndercollateralizedLiquidation")
}

// NotifyUndercollateralizedLiquidation is a paid mutator transaction binding the contract method 0xb4bd2e7a.
//
// Solidity: function notifyUndercollateralizedLiquidation() returns()
func (_Deposit *DepositSession) NotifyUndercollateralizedLiquidation() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyUndercollateralizedLiquidation(&_Deposit.TransactOpts)
}

// NotifyUndercollateralizedLiquidation is a paid mutator transaction binding the contract method 0xb4bd2e7a.
//
// Solidity: function notifyUndercollateralizedLiquidation() returns()
func (_Deposit *DepositTransactorSession) NotifyUndercollateralizedLiquidation() (*types.Transaction, error) {
	return _Deposit.Contract.NotifyUndercollateralizedLiquidation(&_Deposit.TransactOpts)
}

// ProvideBTCFundingProof is a paid mutator transaction binding the contract method 0xd9f74b0e.
//
// Solidity: function provideBTCFundingProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, uint8 _fundingOutputIndex, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositTransactor) ProvideBTCFundingProof(opts *bind.TransactOpts, _txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _fundingOutputIndex uint8, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "provideBTCFundingProof", _txVersion, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideBTCFundingProof is a paid mutator transaction binding the contract method 0xd9f74b0e.
//
// Solidity: function provideBTCFundingProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, uint8 _fundingOutputIndex, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositSession) ProvideBTCFundingProof(_txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _fundingOutputIndex uint8, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideBTCFundingProof(&_Deposit.TransactOpts, _txVersion, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideBTCFundingProof is a paid mutator transaction binding the contract method 0xd9f74b0e.
//
// Solidity: function provideBTCFundingProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, uint8 _fundingOutputIndex, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositTransactorSession) ProvideBTCFundingProof(_txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _fundingOutputIndex uint8, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideBTCFundingProof(&_Deposit.TransactOpts, _txVersion, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideECDSAFraudProof is a paid mutator transaction binding the contract method 0xd5eef971.
//
// Solidity: function provideECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositTransactor) ProvideECDSAFraudProof(opts *bind.TransactOpts, _v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "provideECDSAFraudProof", _v, _r, _s, _signedDigest, _preimage)
}

// ProvideECDSAFraudProof is a paid mutator transaction binding the contract method 0xd5eef971.
//
// Solidity: function provideECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositSession) ProvideECDSAFraudProof(_v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideECDSAFraudProof(&_Deposit.TransactOpts, _v, _r, _s, _signedDigest, _preimage)
}

// ProvideECDSAFraudProof is a paid mutator transaction binding the contract method 0xd5eef971.
//
// Solidity: function provideECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositTransactorSession) ProvideECDSAFraudProof(_v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideECDSAFraudProof(&_Deposit.TransactOpts, _v, _r, _s, _signedDigest, _preimage)
}

// ProvideFundingECDSAFraudProof is a paid mutator transaction binding the contract method 0x2d099442.
//
// Solidity: function provideFundingECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositTransactor) ProvideFundingECDSAFraudProof(opts *bind.TransactOpts, _v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "provideFundingECDSAFraudProof", _v, _r, _s, _signedDigest, _preimage)
}

// ProvideFundingECDSAFraudProof is a paid mutator transaction binding the contract method 0x2d099442.
//
// Solidity: function provideFundingECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositSession) ProvideFundingECDSAFraudProof(_v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideFundingECDSAFraudProof(&_Deposit.TransactOpts, _v, _r, _s, _signedDigest, _preimage)
}

// ProvideFundingECDSAFraudProof is a paid mutator transaction binding the contract method 0x2d099442.
//
// Solidity: function provideFundingECDSAFraudProof(uint8 _v, bytes32 _r, bytes32 _s, bytes32 _signedDigest, bytes _preimage) returns()
func (_Deposit *DepositTransactorSession) ProvideFundingECDSAFraudProof(_v uint8, _r [32]byte, _s [32]byte, _signedDigest [32]byte, _preimage []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideFundingECDSAFraudProof(&_Deposit.TransactOpts, _v, _r, _s, _signedDigest, _preimage)
}

// ProvideRedemptionProof is a paid mutator transaction binding the contract method 0xd459c416.
//
// Solidity: function provideRedemptionProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositTransactor) ProvideRedemptionProof(opts *bind.TransactOpts, _txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "provideRedemptionProof", _txVersion, _txInputVector, _txOutputVector, _txLocktime, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideRedemptionProof is a paid mutator transaction binding the contract method 0xd459c416.
//
// Solidity: function provideRedemptionProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositSession) ProvideRedemptionProof(_txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideRedemptionProof(&_Deposit.TransactOpts, _txVersion, _txInputVector, _txOutputVector, _txLocktime, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideRedemptionProof is a paid mutator transaction binding the contract method 0xd459c416.
//
// Solidity: function provideRedemptionProof(bytes4 _txVersion, bytes _txInputVector, bytes _txOutputVector, bytes4 _txLocktime, bytes _merkleProof, uint256 _txIndexInBlock, bytes _bitcoinHeaders) returns()
func (_Deposit *DepositTransactorSession) ProvideRedemptionProof(_txVersion [4]byte, _txInputVector []byte, _txOutputVector []byte, _txLocktime [4]byte, _merkleProof []byte, _txIndexInBlock *big.Int, _bitcoinHeaders []byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideRedemptionProof(&_Deposit.TransactOpts, _txVersion, _txInputVector, _txOutputVector, _txLocktime, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
}

// ProvideRedemptionSignature is a paid mutator transaction binding the contract method 0xc4159559.
//
// Solidity: function provideRedemptionSignature(uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Deposit *DepositTransactor) ProvideRedemptionSignature(opts *bind.TransactOpts, _v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "provideRedemptionSignature", _v, _r, _s)
}

// ProvideRedemptionSignature is a paid mutator transaction binding the contract method 0xc4159559.
//
// Solidity: function provideRedemptionSignature(uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Deposit *DepositSession) ProvideRedemptionSignature(_v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideRedemptionSignature(&_Deposit.TransactOpts, _v, _r, _s)
}

// ProvideRedemptionSignature is a paid mutator transaction binding the contract method 0xc4159559.
//
// Solidity: function provideRedemptionSignature(uint8 _v, bytes32 _r, bytes32 _s) returns()
func (_Deposit *DepositTransactorSession) ProvideRedemptionSignature(_v uint8, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _Deposit.Contract.ProvideRedemptionSignature(&_Deposit.TransactOpts, _v, _r, _s)
}

// PurchaseSignerBondsAtAuction is a paid mutator transaction binding the contract method 0x2c735daa.
//
// Solidity: function purchaseSignerBondsAtAuction() returns()
func (_Deposit *DepositTransactor) PurchaseSignerBondsAtAuction(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "purchaseSignerBondsAtAuction")
}

// PurchaseSignerBondsAtAuction is a paid mutator transaction binding the contract method 0x2c735daa.
//
// Solidity: function purchaseSignerBondsAtAuction() returns()
func (_Deposit *DepositSession) PurchaseSignerBondsAtAuction() (*types.Transaction, error) {
	return _Deposit.Contract.PurchaseSignerBondsAtAuction(&_Deposit.TransactOpts)
}

// PurchaseSignerBondsAtAuction is a paid mutator transaction binding the contract method 0x2c735daa.
//
// Solidity: function purchaseSignerBondsAtAuction() returns()
func (_Deposit *DepositTransactorSession) PurchaseSignerBondsAtAuction() (*types.Transaction, error) {
	return _Deposit.Contract.PurchaseSignerBondsAtAuction(&_Deposit.TransactOpts)
}

// RequestFunderAbort is a paid mutator transaction binding the contract method 0x0049ce75.
//
// Solidity: function requestFunderAbort(bytes _abortOutputScript) returns()
func (_Deposit *DepositTransactor) RequestFunderAbort(opts *bind.TransactOpts, _abortOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "requestFunderAbort", _abortOutputScript)
}

// RequestFunderAbort is a paid mutator transaction binding the contract method 0x0049ce75.
//
// Solidity: function requestFunderAbort(bytes _abortOutputScript) returns()
func (_Deposit *DepositSession) RequestFunderAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.Contract.RequestFunderAbort(&_Deposit.TransactOpts, _abortOutputScript)
}

// RequestFunderAbort is a paid mutator transaction binding the contract method 0x0049ce75.
//
// Solidity: function requestFunderAbort(bytes _abortOutputScript) returns()
func (_Deposit *DepositTransactorSession) RequestFunderAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.Contract.RequestFunderAbort(&_Deposit.TransactOpts, _abortOutputScript)
}

// RequestRedemption is a paid mutator transaction binding the contract method 0x994aa931.
//
// Solidity: function requestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript) returns()
func (_Deposit *DepositTransactor) RequestRedemption(opts *bind.TransactOpts, _outputValueBytes [8]byte, _redeemerOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "requestRedemption", _outputValueBytes, _redeemerOutputScript)
}

// RequestRedemption is a paid mutator transaction binding the contract method 0x994aa931.
//
// Solidity: function requestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript) returns()
func (_Deposit *DepositSession) RequestRedemption(_outputValueBytes [8]byte, _redeemerOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.Contract.RequestRedemption(&_Deposit.TransactOpts, _outputValueBytes, _redeemerOutputScript)
}

// RequestRedemption is a paid mutator transaction binding the contract method 0x994aa931.
//
// Solidity: function requestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript) returns()
func (_Deposit *DepositTransactorSession) RequestRedemption(_outputValueBytes [8]byte, _redeemerOutputScript []byte) (*types.Transaction, error) {
	return _Deposit.Contract.RequestRedemption(&_Deposit.TransactOpts, _outputValueBytes, _redeemerOutputScript)
}

// RetrieveSignerPubkey is a paid mutator transaction binding the contract method 0xea3db250.
//
// Solidity: function retrieveSignerPubkey() returns()
func (_Deposit *DepositTransactor) RetrieveSignerPubkey(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "retrieveSignerPubkey")
}

// RetrieveSignerPubkey is a paid mutator transaction binding the contract method 0xea3db250.
//
// Solidity: function retrieveSignerPubkey() returns()
func (_Deposit *DepositSession) RetrieveSignerPubkey() (*types.Transaction, error) {
	return _Deposit.Contract.RetrieveSignerPubkey(&_Deposit.TransactOpts)
}

// RetrieveSignerPubkey is a paid mutator transaction binding the contract method 0xea3db250.
//
// Solidity: function retrieveSignerPubkey() returns()
func (_Deposit *DepositTransactorSession) RetrieveSignerPubkey() (*types.Transaction, error) {
	return _Deposit.Contract.RetrieveSignerPubkey(&_Deposit.TransactOpts)
}

// TransferAndRequestRedemption is a paid mutator transaction binding the contract method 0xfb7c592a.
//
// Solidity: function transferAndRequestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript, address _finalRecipient) returns()
func (_Deposit *DepositTransactor) TransferAndRequestRedemption(opts *bind.TransactOpts, _outputValueBytes [8]byte, _redeemerOutputScript []byte, _finalRecipient common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "transferAndRequestRedemption", _outputValueBytes, _redeemerOutputScript, _finalRecipient)
}

// TransferAndRequestRedemption is a paid mutator transaction binding the contract method 0xfb7c592a.
//
// Solidity: function transferAndRequestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript, address _finalRecipient) returns()
func (_Deposit *DepositSession) TransferAndRequestRedemption(_outputValueBytes [8]byte, _redeemerOutputScript []byte, _finalRecipient common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferAndRequestRedemption(&_Deposit.TransactOpts, _outputValueBytes, _redeemerOutputScript, _finalRecipient)
}

// TransferAndRequestRedemption is a paid mutator transaction binding the contract method 0xfb7c592a.
//
// Solidity: function transferAndRequestRedemption(bytes8 _outputValueBytes, bytes _redeemerOutputScript, address _finalRecipient) returns()
func (_Deposit *DepositTransactorSession) TransferAndRequestRedemption(_outputValueBytes [8]byte, _redeemerOutputScript []byte, _finalRecipient common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferAndRequestRedemption(&_Deposit.TransactOpts, _outputValueBytes, _redeemerOutputScript, _finalRecipient)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x24600fc3.
//
// Solidity: function withdrawFunds() returns()
func (_Deposit *DepositTransactor) WithdrawFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "withdrawFunds")
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x24600fc3.
//
// Solidity: function withdrawFunds() returns()
func (_Deposit *DepositSession) WithdrawFunds() (*types.Transaction, error) {
	return _Deposit.Contract.WithdrawFunds(&_Deposit.TransactOpts)
}

// WithdrawFunds is a paid mutator transaction binding the contract method 0x24600fc3.
//
// Solidity: function withdrawFunds() returns()
func (_Deposit *DepositTransactorSession) WithdrawFunds() (*types.Transaction, error) {
	return _Deposit.Contract.WithdrawFunds(&_Deposit.TransactOpts)
}

// TryParseLog attempts to parse a log. Returns the parsed log, evenName and whether it was succesfull
func (_Deposit *DepositFilterer) TryParseLog(log types.Log) (eventName string, event interface{}, ok bool, err error) {
	eventName, ok, err = _Deposit.contract.LogEventName(log)
	if err != nil || !ok {
		return "", nil, false, err
	}

	switch eventName {
	}
	if err != nil {
		return "", nil, false, err
	}

	return eventName, event, ok, nil
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Deposit *DepositTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Deposit.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Deposit *DepositSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Deposit.Contract.Fallback(&_Deposit.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Deposit *DepositTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Deposit.Contract.Fallback(&_Deposit.TransactOpts, calldata)
}
