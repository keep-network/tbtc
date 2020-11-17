// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	ethereumabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/subscription"
	abi "github.com/keep-network/tbtc/pkg/chain/ethereum/gen/abi/system"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var tbtcsLogger = log.Logger("keep-contract-TBTCSystem")

type TBTCSystem struct {
	contract          *abi.TBTCSystem
	contractAddress   common.Address
	contractABI       *ethereumabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *ethutil.ErrorResolver
	nonceManager      *ethutil.NonceManager
	miningWaiter      *ethutil.MiningWaiter

	transactionMutex *sync.Mutex
}

func NewTBTCSystem(
	contractAddress common.Address,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethutil.NonceManager,
	miningWaiter *ethutil.MiningWaiter,
	transactionMutex *sync.Mutex,
) (*TBTCSystem, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	transactorOptions := bind.NewKeyedTransactor(
		accountKey.PrivateKey,
	)

	randomBeaconContract, err := abi.NewTBTCSystem(
		contractAddress,
		backend,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract at address: %s [%v]",
			contractAddress.String(),
			err,
		)
	}

	contractABI, err := ethereumabi.JSON(strings.NewReader(abi.TBTCSystemABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TBTCSystem{
		contract:          randomBeaconContract,
		contractAddress:   contractAddress,
		contractABI:       &contractABI,
		caller:            backend,
		transactor:        backend,
		callerOptions:     callerOptions,
		transactorOptions: transactorOptions,
		errorResolver:     ethutil.NewErrorResolver(backend, &contractABI, &contractAddress),
		nonceManager:      nonceManager,
		miningWaiter:      miningWaiter,
		transactionMutex:  transactionMutex,
	}, nil
}

// ----- Non-const Methods ------

// Transaction submission.
func (tbtcs *TBTCSystem) BeginKeepFactoriesUpdate(
	_keepStakedFactory common.Address,
	_fullyBackedFactory common.Address,
	_factorySelector common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginKeepFactoriesUpdate",
		"params: ",
		fmt.Sprint(
			_keepStakedFactory,
			_fullyBackedFactory,
			_factorySelector,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.BeginKeepFactoriesUpdate(
		transactorOptions,
		_keepStakedFactory,
		_fullyBackedFactory,
		_factorySelector,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"beginKeepFactoriesUpdate",
			_keepStakedFactory,
			_fullyBackedFactory,
			_factorySelector,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction beginKeepFactoriesUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.BeginKeepFactoriesUpdate(
				transactorOptions,
				_keepStakedFactory,
				_fullyBackedFactory,
				_factorySelector,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginKeepFactoriesUpdate",
					_keepStakedFactory,
					_fullyBackedFactory,
					_factorySelector,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginKeepFactoriesUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallBeginKeepFactoriesUpdate(
	_keepStakedFactory common.Address,
	_fullyBackedFactory common.Address,
	_factorySelector common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"beginKeepFactoriesUpdate",
		&result,
		_keepStakedFactory,
		_fullyBackedFactory,
		_factorySelector,
	)

	return err
}

func (tbtcs *TBTCSystem) BeginKeepFactoriesUpdateGasEstimate(
	_keepStakedFactory common.Address,
	_fullyBackedFactory common.Address,
	_factorySelector common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"beginKeepFactoriesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
		_keepStakedFactory,
		_fullyBackedFactory,
		_factorySelector,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) RenounceOwnership(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction renounceOwnership",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.RenounceOwnership(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"renounceOwnership",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction renounceOwnership with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.RenounceOwnership(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"renounceOwnership",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction renounceOwnership with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallRenounceOwnership(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"renounceOwnership",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) RenounceOwnershipGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"renounceOwnership",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeEthBtcPriceFeedAddition(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction finalizeEthBtcPriceFeedAddition",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.FinalizeEthBtcPriceFeedAddition(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"finalizeEthBtcPriceFeedAddition",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction finalizeEthBtcPriceFeedAddition with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.FinalizeEthBtcPriceFeedAddition(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeEthBtcPriceFeedAddition",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeEthBtcPriceFeedAddition with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallFinalizeEthBtcPriceFeedAddition(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"finalizeEthBtcPriceFeedAddition",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) FinalizeEthBtcPriceFeedAdditionGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeEthBtcPriceFeedAddition",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogExitedCourtesyCall(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logExitedCourtesyCall",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogExitedCourtesyCall(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logExitedCourtesyCall",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logExitedCourtesyCall with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogExitedCourtesyCall(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logExitedCourtesyCall",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logExitedCourtesyCall with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogExitedCourtesyCall(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logExitedCourtesyCall",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) LogExitedCourtesyCallGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logExitedCourtesyCall",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) BeginSignerFeeDivisorUpdate(
	_signerFeeDivisor uint16,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginSignerFeeDivisorUpdate",
		"params: ",
		fmt.Sprint(
			_signerFeeDivisor,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.BeginSignerFeeDivisorUpdate(
		transactorOptions,
		_signerFeeDivisor,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"beginSignerFeeDivisorUpdate",
			_signerFeeDivisor,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction beginSignerFeeDivisorUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.BeginSignerFeeDivisorUpdate(
				transactorOptions,
				_signerFeeDivisor,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginSignerFeeDivisorUpdate",
					_signerFeeDivisor,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginSignerFeeDivisorUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallBeginSignerFeeDivisorUpdate(
	_signerFeeDivisor uint16,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"beginSignerFeeDivisorUpdate",
		&result,
		_signerFeeDivisor,
	)

	return err
}

func (tbtcs *TBTCSystem) BeginSignerFeeDivisorUpdateGasEstimate(
	_signerFeeDivisor uint16,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"beginSignerFeeDivisorUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
		_signerFeeDivisor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogFunderRequestedAbort(
	_abortOutputScript []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logFunderRequestedAbort",
		"params: ",
		fmt.Sprint(
			_abortOutputScript,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogFunderRequestedAbort(
		transactorOptions,
		_abortOutputScript,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logFunderRequestedAbort",
			_abortOutputScript,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logFunderRequestedAbort with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogFunderRequestedAbort(
				transactorOptions,
				_abortOutputScript,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFunderRequestedAbort",
					_abortOutputScript,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFunderRequestedAbort with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogFunderRequestedAbort(
	_abortOutputScript []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logFunderRequestedAbort",
		&result,
		_abortOutputScript,
	)

	return err
}

func (tbtcs *TBTCSystem) LogFunderRequestedAbortGasEstimate(
	_abortOutputScript []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logFunderRequestedAbort",
		tbtcs.contractABI,
		tbtcs.transactor,
		_abortOutputScript,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogGotRedemptionSignature(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logGotRedemptionSignature",
		"params: ",
		fmt.Sprint(
			_digest,
			_r,
			_s,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogGotRedemptionSignature(
		transactorOptions,
		_digest,
		_r,
		_s,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logGotRedemptionSignature",
			_digest,
			_r,
			_s,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logGotRedemptionSignature with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogGotRedemptionSignature(
				transactorOptions,
				_digest,
				_r,
				_s,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logGotRedemptionSignature",
					_digest,
					_r,
					_s,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logGotRedemptionSignature with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogGotRedemptionSignature(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logGotRedemptionSignature",
		&result,
		_digest,
		_r,
		_s,
	)

	return err
}

func (tbtcs *TBTCSystem) LogGotRedemptionSignatureGasEstimate(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logGotRedemptionSignature",
		tbtcs.contractABI,
		tbtcs.transactor,
		_digest,
		_r,
		_s,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) EmergencyPauseNewDeposits(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction emergencyPauseNewDeposits",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.EmergencyPauseNewDeposits(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"emergencyPauseNewDeposits",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction emergencyPauseNewDeposits with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.EmergencyPauseNewDeposits(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"emergencyPauseNewDeposits",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction emergencyPauseNewDeposits with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallEmergencyPauseNewDeposits(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"emergencyPauseNewDeposits",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) EmergencyPauseNewDepositsGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"emergencyPauseNewDeposits",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeKeepFactoriesUpdate(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction finalizeKeepFactoriesUpdate",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.FinalizeKeepFactoriesUpdate(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"finalizeKeepFactoriesUpdate",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction finalizeKeepFactoriesUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.FinalizeKeepFactoriesUpdate(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeKeepFactoriesUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeKeepFactoriesUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallFinalizeKeepFactoriesUpdate(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"finalizeKeepFactoriesUpdate",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) FinalizeKeepFactoriesUpdateGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeKeepFactoriesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogStartedLiquidation(
	_wasFraud bool,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logStartedLiquidation",
		"params: ",
		fmt.Sprint(
			_wasFraud,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogStartedLiquidation(
		transactorOptions,
		_wasFraud,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logStartedLiquidation",
			_wasFraud,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logStartedLiquidation with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogStartedLiquidation(
				transactorOptions,
				_wasFraud,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logStartedLiquidation",
					_wasFraud,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logStartedLiquidation with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogStartedLiquidation(
	_wasFraud bool,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logStartedLiquidation",
		&result,
		_wasFraud,
	)

	return err
}

func (tbtcs *TBTCSystem) LogStartedLiquidationGasEstimate(
	_wasFraud bool,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logStartedLiquidation",
		tbtcs.contractABI,
		tbtcs.transactor,
		_wasFraud,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogFunded(
	_txid [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logFunded",
		"params: ",
		fmt.Sprint(
			_txid,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogFunded(
		transactorOptions,
		_txid,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logFunded",
			_txid,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logFunded with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogFunded(
				transactorOptions,
				_txid,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFunded",
					_txid,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFunded with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogFunded(
	_txid [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logFunded",
		&result,
		_txid,
	)

	return err
}

func (tbtcs *TBTCSystem) LogFundedGasEstimate(
	_txid [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logFunded",
		tbtcs.contractABI,
		tbtcs.transactor,
		_txid,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogSetupFailed(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logSetupFailed",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogSetupFailed(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logSetupFailed",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logSetupFailed with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogSetupFailed(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logSetupFailed",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logSetupFailed with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogSetupFailed(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logSetupFailed",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) LogSetupFailedGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logSetupFailed",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogFraudDuringSetup(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logFraudDuringSetup",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogFraudDuringSetup(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logFraudDuringSetup",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logFraudDuringSetup with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogFraudDuringSetup(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFraudDuringSetup",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFraudDuringSetup with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogFraudDuringSetup(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logFraudDuringSetup",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) LogFraudDuringSetupGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logFraudDuringSetup",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeLotSizesUpdate(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction finalizeLotSizesUpdate",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.FinalizeLotSizesUpdate(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"finalizeLotSizesUpdate",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction finalizeLotSizesUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.FinalizeLotSizesUpdate(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeLotSizesUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeLotSizesUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallFinalizeLotSizesUpdate(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"finalizeLotSizesUpdate",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) FinalizeLotSizesUpdateGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeLotSizesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) BeginLotSizesUpdate(
	_lotSizes []uint64,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginLotSizesUpdate",
		"params: ",
		fmt.Sprint(
			_lotSizes,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.BeginLotSizesUpdate(
		transactorOptions,
		_lotSizes,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"beginLotSizesUpdate",
			_lotSizes,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction beginLotSizesUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.BeginLotSizesUpdate(
				transactorOptions,
				_lotSizes,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginLotSizesUpdate",
					_lotSizes,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginLotSizesUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallBeginLotSizesUpdate(
	_lotSizes []uint64,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"beginLotSizesUpdate",
		&result,
		_lotSizes,
	)

	return err
}

func (tbtcs *TBTCSystem) BeginLotSizesUpdateGasEstimate(
	_lotSizes []uint64,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"beginLotSizesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
		_lotSizes,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogLiquidated(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logLiquidated",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogLiquidated(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logLiquidated",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logLiquidated with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogLiquidated(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logLiquidated",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logLiquidated with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogLiquidated(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logLiquidated",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) LogLiquidatedGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logLiquidated",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) ResumeNewDeposits(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction resumeNewDeposits",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.ResumeNewDeposits(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"resumeNewDeposits",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction resumeNewDeposits with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.ResumeNewDeposits(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"resumeNewDeposits",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction resumeNewDeposits with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallResumeNewDeposits(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"resumeNewDeposits",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) ResumeNewDepositsGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"resumeNewDeposits",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogRegisteredPubkey(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRegisteredPubkey",
		"params: ",
		fmt.Sprint(
			_signingGroupPubkeyX,
			_signingGroupPubkeyY,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogRegisteredPubkey(
		transactorOptions,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logRegisteredPubkey",
			_signingGroupPubkeyX,
			_signingGroupPubkeyY,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logRegisteredPubkey with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogRegisteredPubkey(
				transactorOptions,
				_signingGroupPubkeyX,
				_signingGroupPubkeyY,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logRegisteredPubkey",
					_signingGroupPubkeyX,
					_signingGroupPubkeyY,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logRegisteredPubkey with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogRegisteredPubkey(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logRegisteredPubkey",
		&result,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)

	return err
}

func (tbtcs *TBTCSystem) LogRegisteredPubkeyGasEstimate(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logRegisteredPubkey",
		tbtcs.contractABI,
		tbtcs.transactor,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) RefreshMinimumBondableValue(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction refreshMinimumBondableValue",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.RefreshMinimumBondableValue(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"refreshMinimumBondableValue",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction refreshMinimumBondableValue with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.RefreshMinimumBondableValue(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"refreshMinimumBondableValue",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction refreshMinimumBondableValue with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallRefreshMinimumBondableValue(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"refreshMinimumBondableValue",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) RefreshMinimumBondableValueGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"refreshMinimumBondableValue",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) Initialize(
	_defaultKeepFactory common.Address,
	_depositFactory common.Address,
	_masterDepositAddress common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachine common.Address,
	_keepThreshold uint16,
	_keepSize uint16,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction initialize",
		"params: ",
		fmt.Sprint(
			_defaultKeepFactory,
			_depositFactory,
			_masterDepositAddress,
			_tbtcToken,
			_tbtcDepositToken,
			_feeRebateToken,
			_vendingMachine,
			_keepThreshold,
			_keepSize,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.Initialize(
		transactorOptions,
		_defaultKeepFactory,
		_depositFactory,
		_masterDepositAddress,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachine,
		_keepThreshold,
		_keepSize,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"initialize",
			_defaultKeepFactory,
			_depositFactory,
			_masterDepositAddress,
			_tbtcToken,
			_tbtcDepositToken,
			_feeRebateToken,
			_vendingMachine,
			_keepThreshold,
			_keepSize,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction initialize with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.Initialize(
				transactorOptions,
				_defaultKeepFactory,
				_depositFactory,
				_masterDepositAddress,
				_tbtcToken,
				_tbtcDepositToken,
				_feeRebateToken,
				_vendingMachine,
				_keepThreshold,
				_keepSize,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"initialize",
					_defaultKeepFactory,
					_depositFactory,
					_masterDepositAddress,
					_tbtcToken,
					_tbtcDepositToken,
					_feeRebateToken,
					_vendingMachine,
					_keepThreshold,
					_keepSize,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction initialize with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallInitialize(
	_defaultKeepFactory common.Address,
	_depositFactory common.Address,
	_masterDepositAddress common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachine common.Address,
	_keepThreshold uint16,
	_keepSize uint16,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"initialize",
		&result,
		_defaultKeepFactory,
		_depositFactory,
		_masterDepositAddress,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachine,
		_keepThreshold,
		_keepSize,
	)

	return err
}

func (tbtcs *TBTCSystem) InitializeGasEstimate(
	_defaultKeepFactory common.Address,
	_depositFactory common.Address,
	_masterDepositAddress common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachine common.Address,
	_keepThreshold uint16,
	_keepSize uint16,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"initialize",
		tbtcs.contractABI,
		tbtcs.transactor,
		_defaultKeepFactory,
		_depositFactory,
		_masterDepositAddress,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachine,
		_keepThreshold,
		_keepSize,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogRedemptionRequested(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRedemptionRequested",
		"params: ",
		fmt.Sprint(
			_requester,
			_digest,
			_utxoValue,
			_redeemerOutputScript,
			_requestedFee,
			_outpoint,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogRedemptionRequested(
		transactorOptions,
		_requester,
		_digest,
		_utxoValue,
		_redeemerOutputScript,
		_requestedFee,
		_outpoint,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logRedemptionRequested",
			_requester,
			_digest,
			_utxoValue,
			_redeemerOutputScript,
			_requestedFee,
			_outpoint,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logRedemptionRequested with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogRedemptionRequested(
				transactorOptions,
				_requester,
				_digest,
				_utxoValue,
				_redeemerOutputScript,
				_requestedFee,
				_outpoint,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logRedemptionRequested",
					_requester,
					_digest,
					_utxoValue,
					_redeemerOutputScript,
					_requestedFee,
					_outpoint,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logRedemptionRequested with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogRedemptionRequested(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logRedemptionRequested",
		&result,
		_requester,
		_digest,
		_utxoValue,
		_redeemerOutputScript,
		_requestedFee,
		_outpoint,
	)

	return err
}

func (tbtcs *TBTCSystem) LogRedemptionRequestedGasEstimate(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logRedemptionRequested",
		tbtcs.contractABI,
		tbtcs.transactor,
		_requester,
		_digest,
		_utxoValue,
		_redeemerOutputScript,
		_requestedFee,
		_outpoint,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) TransferOwnership(
	newOwner common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction transferOwnership",
		"params: ",
		fmt.Sprint(
			newOwner,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.TransferOwnership(
		transactorOptions,
		newOwner,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"transferOwnership",
			newOwner,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction transferOwnership with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.TransferOwnership(
				transactorOptions,
				newOwner,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"transferOwnership",
					newOwner,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction transferOwnership with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallTransferOwnership(
	newOwner common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"transferOwnership",
		&result,
		newOwner,
	)

	return err
}

func (tbtcs *TBTCSystem) TransferOwnershipGasEstimate(
	newOwner common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"transferOwnership",
		tbtcs.contractABI,
		tbtcs.transactor,
		newOwner,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeSignerFeeDivisorUpdate(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction finalizeSignerFeeDivisorUpdate",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.FinalizeSignerFeeDivisorUpdate(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"finalizeSignerFeeDivisorUpdate",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction finalizeSignerFeeDivisorUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.FinalizeSignerFeeDivisorUpdate(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeSignerFeeDivisorUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeSignerFeeDivisorUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallFinalizeSignerFeeDivisorUpdate(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"finalizeSignerFeeDivisorUpdate",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) FinalizeSignerFeeDivisorUpdateGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeSignerFeeDivisorUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) BeginCollateralizationThresholdsUpdate(
	_initialCollateralizedPercent uint16,
	_undercollateralizedThresholdPercent uint16,
	_severelyUndercollateralizedThresholdPercent uint16,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginCollateralizationThresholdsUpdate",
		"params: ",
		fmt.Sprint(
			_initialCollateralizedPercent,
			_undercollateralizedThresholdPercent,
			_severelyUndercollateralizedThresholdPercent,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.BeginCollateralizationThresholdsUpdate(
		transactorOptions,
		_initialCollateralizedPercent,
		_undercollateralizedThresholdPercent,
		_severelyUndercollateralizedThresholdPercent,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"beginCollateralizationThresholdsUpdate",
			_initialCollateralizedPercent,
			_undercollateralizedThresholdPercent,
			_severelyUndercollateralizedThresholdPercent,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction beginCollateralizationThresholdsUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.BeginCollateralizationThresholdsUpdate(
				transactorOptions,
				_initialCollateralizedPercent,
				_undercollateralizedThresholdPercent,
				_severelyUndercollateralizedThresholdPercent,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginCollateralizationThresholdsUpdate",
					_initialCollateralizedPercent,
					_undercollateralizedThresholdPercent,
					_severelyUndercollateralizedThresholdPercent,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginCollateralizationThresholdsUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallBeginCollateralizationThresholdsUpdate(
	_initialCollateralizedPercent uint16,
	_undercollateralizedThresholdPercent uint16,
	_severelyUndercollateralizedThresholdPercent uint16,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"beginCollateralizationThresholdsUpdate",
		&result,
		_initialCollateralizedPercent,
		_undercollateralizedThresholdPercent,
		_severelyUndercollateralizedThresholdPercent,
	)

	return err
}

func (tbtcs *TBTCSystem) BeginCollateralizationThresholdsUpdateGasEstimate(
	_initialCollateralizedPercent uint16,
	_undercollateralizedThresholdPercent uint16,
	_severelyUndercollateralizedThresholdPercent uint16,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"beginCollateralizationThresholdsUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
		_initialCollateralizedPercent,
		_undercollateralizedThresholdPercent,
		_severelyUndercollateralizedThresholdPercent,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) BeginEthBtcPriceFeedAddition(
	_ethBtcPriceFeed common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginEthBtcPriceFeedAddition",
		"params: ",
		fmt.Sprint(
			_ethBtcPriceFeed,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.BeginEthBtcPriceFeedAddition(
		transactorOptions,
		_ethBtcPriceFeed,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"beginEthBtcPriceFeedAddition",
			_ethBtcPriceFeed,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction beginEthBtcPriceFeedAddition with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.BeginEthBtcPriceFeedAddition(
				transactorOptions,
				_ethBtcPriceFeed,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginEthBtcPriceFeedAddition",
					_ethBtcPriceFeed,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginEthBtcPriceFeedAddition with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallBeginEthBtcPriceFeedAddition(
	_ethBtcPriceFeed common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"beginEthBtcPriceFeedAddition",
		&result,
		_ethBtcPriceFeed,
	)

	return err
}

func (tbtcs *TBTCSystem) BeginEthBtcPriceFeedAdditionGasEstimate(
	_ethBtcPriceFeed common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"beginEthBtcPriceFeedAddition",
		tbtcs.contractABI,
		tbtcs.transactor,
		_ethBtcPriceFeed,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogCourtesyCalled(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logCourtesyCalled",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogCourtesyCalled(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logCourtesyCalled",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logCourtesyCalled with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogCourtesyCalled(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logCourtesyCalled",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logCourtesyCalled with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogCourtesyCalled(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logCourtesyCalled",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) LogCourtesyCalledGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logCourtesyCalled",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogCreated(
	_keepAddress common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logCreated",
		"params: ",
		fmt.Sprint(
			_keepAddress,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogCreated(
		transactorOptions,
		_keepAddress,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logCreated",
			_keepAddress,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logCreated with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogCreated(
				transactorOptions,
				_keepAddress,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logCreated",
					_keepAddress,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logCreated with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogCreated(
	_keepAddress common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logCreated",
		&result,
		_keepAddress,
	)

	return err
}

func (tbtcs *TBTCSystem) LogCreatedGasEstimate(
	_keepAddress common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logCreated",
		tbtcs.contractABI,
		tbtcs.transactor,
		_keepAddress,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeCollateralizationThresholdsUpdate(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction finalizeCollateralizationThresholdsUpdate",
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.FinalizeCollateralizationThresholdsUpdate(
		transactorOptions,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"finalizeCollateralizationThresholdsUpdate",
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction finalizeCollateralizationThresholdsUpdate with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.FinalizeCollateralizationThresholdsUpdate(
				transactorOptions,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeCollateralizationThresholdsUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeCollateralizationThresholdsUpdate with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallFinalizeCollateralizationThresholdsUpdate(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"finalizeCollateralizationThresholdsUpdate",
		&result,
	)

	return err
}

func (tbtcs *TBTCSystem) FinalizeCollateralizationThresholdsUpdateGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeCollateralizationThresholdsUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogRedeemed(
	_txid [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRedeemed",
		"params: ",
		fmt.Sprint(
			_txid,
		),
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.LogRedeemed(
		transactorOptions,
		_txid,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			nil,
			"logRedeemed",
			_txid,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction logRedeemed with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.LogRedeemed(
				transactorOptions,
				_txid,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logRedeemed",
					_txid,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logRedeemed with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallLogRedeemed(
	_txid [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"logRedeemed",
		&result,
		_txid,
	)

	return err
}

func (tbtcs *TBTCSystem) LogRedeemedGasEstimate(
	_txid [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logRedeemed",
		tbtcs.contractABI,
		tbtcs.transactor,
		_txid,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) RequestNewKeep(
	_requestedLotSizeSatoshis uint64,
	_maxSecuredLifetime *big.Int,
	value *big.Int,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction requestNewKeep",
		"params: ",
		fmt.Sprint(
			_requestedLotSizeSatoshis,
			_maxSecuredLifetime,
		),
		"value: ", value,
	)

	tbtcs.transactionMutex.Lock()
	defer tbtcs.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *tbtcs.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := tbtcs.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := tbtcs.contract.RequestNewKeep(
		transactorOptions,
		_requestedLotSizeSatoshis,
		_maxSecuredLifetime,
	)
	if err != nil {
		return transaction, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.transactorOptions.From,
			value,
			"requestNewKeep",
			_requestedLotSizeSatoshis,
			_maxSecuredLifetime,
		)
	}

	tbtcsLogger.Infof(
		"submitted transaction requestNewKeep with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := tbtcs.contract.RequestNewKeep(
				transactorOptions,
				_requestedLotSizeSatoshis,
				_maxSecuredLifetime,
			)
			if err != nil {
				return transaction, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					value,
					"requestNewKeep",
					_requestedLotSizeSatoshis,
					_maxSecuredLifetime,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction requestNewKeep with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	tbtcs.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (tbtcs *TBTCSystem) CallRequestNewKeep(
	_requestedLotSizeSatoshis uint64,
	_maxSecuredLifetime *big.Int,
	value *big.Int,
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tbtcs.transactorOptions.From,
		blockNumber, value,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"requestNewKeep",
		&result,
		_requestedLotSizeSatoshis,
		_maxSecuredLifetime,
	)

	return result, err
}

func (tbtcs *TBTCSystem) RequestNewKeepGasEstimate(
	_requestedLotSizeSatoshis uint64,
	_maxSecuredLifetime *big.Int,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"requestNewKeep",
		tbtcs.contractABI,
		tbtcs.transactor,
		_requestedLotSizeSatoshis,
		_maxSecuredLifetime,
	)

	return result, err
}

// ----- Const Methods ------

func (tbtcs *TBTCSystem) KeepThreshold() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.KeepThreshold(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"keepThreshold",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) KeepThresholdAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"keepThreshold",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) FetchRelayPreviousDifficulty() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.FetchRelayPreviousDifficulty(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"fetchRelayPreviousDifficulty",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) FetchRelayPreviousDifficultyAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"fetchRelayPreviousDifficulty",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingPauseTerm() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingPauseTerm(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingPauseTerm",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingPauseTermAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingPauseTerm",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) KeepSize() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.KeepSize(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"keepSize",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) KeepSizeAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"keepSize",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) FetchRelayCurrentDifficulty() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.FetchRelayCurrentDifficulty(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"fetchRelayCurrentDifficulty",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) FetchRelayCurrentDifficultyAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"fetchRelayCurrentDifficulty",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetGovernanceTimeDelay() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetGovernanceTimeDelay(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getGovernanceTimeDelay",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetGovernanceTimeDelayAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getGovernanceTimeDelay",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetInitialCollateralizedPercent() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.GetInitialCollateralizedPercent(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getInitialCollateralizedPercent",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetInitialCollateralizedPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getInitialCollateralizedPercent",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetMinimumLotSize() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetMinimumLotSize(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getMinimumLotSize",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetMinimumLotSizeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getMinimumLotSize",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) ApprovedToLog(
	_caller common.Address,
) (bool, error) {
	var result bool
	result, err := tbtcs.contract.ApprovedToLog(
		tbtcs.callerOptions,
		_caller,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"approvedToLog",
			_caller,
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) ApprovedToLogAtBlock(
	_caller common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"approvedToLog",
		&result,
		_caller,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetAllowNewDeposits() (bool, error) {
	var result bool
	result, err := tbtcs.contract.GetAllowNewDeposits(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getAllowNewDeposits",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetAllowNewDepositsAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getAllowNewDeposits",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingEthBtcPriceFeedAdditionTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingEthBtcPriceFeedAdditionTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingEthBtcPriceFeedAdditionTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingEthBtcPriceFeedAdditionTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingEthBtcPriceFeedAdditionTime",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingCollateralizationThresholdsUpdateTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingCollateralizationThresholdsUpdateTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingCollateralizationThresholdsUpdateTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingCollateralizationThresholdsUpdateTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingCollateralizationThresholdsUpdateTime",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) IsAllowedLotSize(
	_requestedLotSizeSatoshis uint64,
) (bool, error) {
	var result bool
	result, err := tbtcs.contract.IsAllowedLotSize(
		tbtcs.callerOptions,
		_requestedLotSizeSatoshis,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"isAllowedLotSize",
			_requestedLotSizeSatoshis,
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) IsAllowedLotSizeAtBlock(
	_requestedLotSizeSatoshis uint64,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"isAllowedLotSize",
		&result,
		_requestedLotSizeSatoshis,
	)

	return result, err
}

func (tbtcs *TBTCSystem) PriceFeed() (common.Address, error) {
	var result common.Address
	result, err := tbtcs.contract.PriceFeed(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"priceFeed",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) PriceFeedAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"priceFeed",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetNewDepositFeeEstimate() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetNewDepositFeeEstimate(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getNewDepositFeeEstimate",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetNewDepositFeeEstimateAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getNewDepositFeeEstimate",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) IsOwner() (bool, error) {
	var result bool
	result, err := tbtcs.contract.IsOwner(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"isOwner",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) IsOwnerAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"isOwner",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetSeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.GetSeverelyUndercollateralizedThresholdPercent(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getSeverelyUndercollateralizedThresholdPercent",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetSeverelyUndercollateralizedThresholdPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getSeverelyUndercollateralizedThresholdPercent",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetSignerFeeDivisor() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.GetSignerFeeDivisor(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getSignerFeeDivisor",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetSignerFeeDivisorAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getSignerFeeDivisor",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetUndercollateralizedThresholdPercent() (uint16, error) {
	var result uint16
	result, err := tbtcs.contract.GetUndercollateralizedThresholdPercent(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getUndercollateralizedThresholdPercent",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetUndercollateralizedThresholdPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getUndercollateralizedThresholdPercent",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) Relay() (common.Address, error) {
	var result common.Address
	result, err := tbtcs.contract.Relay(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"relay",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) RelayAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"relay",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetAllowedLotSizes() ([]uint64, error) {
	var result []uint64
	result, err := tbtcs.contract.GetAllowedLotSizes(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getAllowedLotSizes",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetAllowedLotSizesAtBlock(
	blockNumber *big.Int,
) ([]uint64, error) {
	var result []uint64

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getAllowedLotSizes",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetPriceFeedGovernanceTimeDelay() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetPriceFeedGovernanceTimeDelay(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getPriceFeedGovernanceTimeDelay",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetPriceFeedGovernanceTimeDelayAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getPriceFeedGovernanceTimeDelay",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) FetchBitcoinPrice() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.FetchBitcoinPrice(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"fetchBitcoinPrice",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) FetchBitcoinPriceAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"fetchBitcoinPrice",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetMaximumLotSize() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetMaximumLotSize(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getMaximumLotSize",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetMaximumLotSizeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getMaximumLotSize",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetKeepFactoriesUpgradeabilityPeriod() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetKeepFactoriesUpgradeabilityPeriod(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getKeepFactoriesUpgradeabilityPeriod",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetKeepFactoriesUpgradeabilityPeriodAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getKeepFactoriesUpgradeabilityPeriod",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingLotSizesUpdateTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingLotSizesUpdateTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingLotSizesUpdateTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingLotSizesUpdateTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingLotSizesUpdateTime",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingSignerFeeDivisorUpdateTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingSignerFeeDivisorUpdateTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingSignerFeeDivisorUpdateTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingSignerFeeDivisorUpdateTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingSignerFeeDivisorUpdateTime",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) Owner() (common.Address, error) {
	var result common.Address
	result, err := tbtcs.contract.Owner(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"owner",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) OwnerAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"owner",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingKeepFactoriesUpdateTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingKeepFactoriesUpdateTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingKeepFactoriesUpdateTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingKeepFactoriesUpdateTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingKeepFactoriesUpdateTime",
		&result,
	)

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingKeepFactoriesUpgradeabilityTime() (*big.Int, error) {
	var result *big.Int
	result, err := tbtcs.contract.GetRemainingKeepFactoriesUpgradeabilityTime(
		tbtcs.callerOptions,
	)

	if err != nil {
		return result, tbtcs.errorResolver.ResolveError(
			err,
			tbtcs.callerOptions.From,
			nil,
			"getRemainingKeepFactoriesUpgradeabilityTime",
		)
	}

	return result, err
}

func (tbtcs *TBTCSystem) GetRemainingKeepFactoriesUpgradeabilityTimeAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := ethutil.CallAtBlock(
		tbtcs.callerOptions.From,
		blockNumber,
		nil,
		tbtcs.contractABI,
		tbtcs.caller,
		tbtcs.errorResolver,
		tbtcs.contractAddress,
		"getRemainingKeepFactoriesUpgradeabilityTime",
		&result,
	)

	return result, err
}

// ------ Events -------

type tBTCSystemKeepFactoriesUpdateStartedFunc func(
	KeepStakedFactory common.Address,
	FullyBackedFactory common.Address,
	FactorySelector common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastKeepFactoriesUpdateStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemKeepFactoriesUpdateStarted, error) {
	iterator, err := tbtcs.contract.FilterKeepFactoriesUpdateStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past KeepFactoriesUpdateStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemKeepFactoriesUpdateStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchKeepFactoriesUpdateStarted(
	success tBTCSystemKeepFactoriesUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeKeepFactoriesUpdateStarted(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event KeepFactoriesUpdateStarted terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeKeepFactoriesUpdateStarted(
	success tBTCSystemKeepFactoriesUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemKeepFactoriesUpdateStarted)
	eventSubscription, err := tbtcs.contract.WatchKeepFactoriesUpdateStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for KeepFactoriesUpdateStarted events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.KeepStakedFactory,
					event.FullyBackedFactory,
					event.FactorySelector,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemKeepFactoriesUpdatedFunc func(
	KeepStakedFactory common.Address,
	FullyBackedFactory common.Address,
	FactorySelector common.Address,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastKeepFactoriesUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemKeepFactoriesUpdated, error) {
	iterator, err := tbtcs.contract.FilterKeepFactoriesUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past KeepFactoriesUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemKeepFactoriesUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchKeepFactoriesUpdated(
	success tBTCSystemKeepFactoriesUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeKeepFactoriesUpdated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event KeepFactoriesUpdated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeKeepFactoriesUpdated(
	success tBTCSystemKeepFactoriesUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemKeepFactoriesUpdated)
	eventSubscription, err := tbtcs.contract.WatchKeepFactoriesUpdated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for KeepFactoriesUpdated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.KeepStakedFactory,
					event.FullyBackedFactory,
					event.FactorySelector,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemSetupFailedFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastSetupFailedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemSetupFailed, error) {
	iterator, err := tbtcs.contract.FilterSetupFailed(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past SetupFailed events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemSetupFailed, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchSetupFailed(
	success tBTCSystemSetupFailedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeSetupFailed(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event SetupFailed terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeSetupFailed(
	success tBTCSystemSetupFailedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemSetupFailed)
	eventSubscription, err := tbtcs.contract.WatchSetupFailed(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for SetupFailed events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemFunderAbortRequestedFunc func(
	DepositContractAddress common.Address,
	AbortOutputScript []uint8,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastFunderAbortRequestedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemFunderAbortRequested, error) {
	iterator, err := tbtcs.contract.FilterFunderAbortRequested(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past FunderAbortRequested events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemFunderAbortRequested, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchFunderAbortRequested(
	success tBTCSystemFunderAbortRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeFunderAbortRequested(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event FunderAbortRequested terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeFunderAbortRequested(
	success tBTCSystemFunderAbortRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemFunderAbortRequested)
	eventSubscription, err := tbtcs.contract.WatchFunderAbortRequested(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for FunderAbortRequested events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.AbortOutputScript,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemLotSizesUpdatedFunc func(
	LotSizes []uint64,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastLotSizesUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemLotSizesUpdated, error) {
	iterator, err := tbtcs.contract.FilterLotSizesUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past LotSizesUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemLotSizesUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchLotSizesUpdated(
	success tBTCSystemLotSizesUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeLotSizesUpdated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event LotSizesUpdated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeLotSizesUpdated(
	success tBTCSystemLotSizesUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemLotSizesUpdated)
	eventSubscription, err := tbtcs.contract.WatchLotSizesUpdated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for LotSizesUpdated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.LotSizes,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemRedemptionRequestedFunc func(
	DepositContractAddress common.Address,
	Requester common.Address,
	Digest [32]uint8,
	UtxoValue *big.Int,
	RedeemerOutputScript []uint8,
	RequestedFee *big.Int,
	Outpoint []uint8,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastRedemptionRequestedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) ([]*abi.TBTCSystemRedemptionRequested, error) {
	iterator, err := tbtcs.contract.FilterRedemptionRequested(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
		_requesterFilter,
		_digestFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RedemptionRequested events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemRedemptionRequested, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchRedemptionRequested(
	success tBTCSystemRedemptionRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeRedemptionRequested(
			success,
			failCallback,
			_depositContractAddressFilter,
			_requesterFilter,
			_digestFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event RedemptionRequested terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeRedemptionRequested(
	success tBTCSystemRedemptionRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemRedemptionRequested)
	eventSubscription, err := tbtcs.contract.WatchRedemptionRequested(
		nil,
		eventChan,
		_depositContractAddressFilter,
		_requesterFilter,
		_digestFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RedemptionRequested events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Requester,
					event.Digest,
					event.UtxoValue,
					event.RedeemerOutputScript,
					event.RequestedFee,
					event.Outpoint,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemRegisteredPubkeyFunc func(
	DepositContractAddress common.Address,
	SigningGroupPubkeyX [32]uint8,
	SigningGroupPubkeyY [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastRegisteredPubkeyEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemRegisteredPubkey, error) {
	iterator, err := tbtcs.contract.FilterRegisteredPubkey(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past RegisteredPubkey events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemRegisteredPubkey, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchRegisteredPubkey(
	success tBTCSystemRegisteredPubkeyFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeRegisteredPubkey(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event RegisteredPubkey terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeRegisteredPubkey(
	success tBTCSystemRegisteredPubkeyFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemRegisteredPubkey)
	eventSubscription, err := tbtcs.contract.WatchRegisteredPubkey(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for RegisteredPubkey events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.SigningGroupPubkeyX,
					event.SigningGroupPubkeyY,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemCourtesyCalledFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastCourtesyCalledEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemCourtesyCalled, error) {
	iterator, err := tbtcs.contract.FilterCourtesyCalled(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past CourtesyCalled events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemCourtesyCalled, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchCourtesyCalled(
	success tBTCSystemCourtesyCalledFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeCourtesyCalled(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event CourtesyCalled terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeCourtesyCalled(
	success tBTCSystemCourtesyCalledFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemCourtesyCalled)
	eventSubscription, err := tbtcs.contract.WatchCourtesyCalled(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for CourtesyCalled events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemCreatedFunc func(
	DepositContractAddress common.Address,
	KeepAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastCreatedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) ([]*abi.TBTCSystemCreated, error) {
	iterator, err := tbtcs.contract.FilterCreated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
		_keepAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Created events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemCreated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchCreated(
	success tBTCSystemCreatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeCreated(
			success,
			failCallback,
			_depositContractAddressFilter,
			_keepAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event Created terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeCreated(
	success tBTCSystemCreatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemCreated)
	eventSubscription, err := tbtcs.contract.WatchCreated(
		nil,
		eventChan,
		_depositContractAddressFilter,
		_keepAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Created events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.KeepAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemEthBtcPriceFeedAdditionStartedFunc func(
	PriceFeed common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastEthBtcPriceFeedAdditionStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemEthBtcPriceFeedAdditionStarted, error) {
	iterator, err := tbtcs.contract.FilterEthBtcPriceFeedAdditionStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past EthBtcPriceFeedAdditionStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemEthBtcPriceFeedAdditionStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchEthBtcPriceFeedAdditionStarted(
	success tBTCSystemEthBtcPriceFeedAdditionStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeEthBtcPriceFeedAdditionStarted(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event EthBtcPriceFeedAdditionStarted terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeEthBtcPriceFeedAdditionStarted(
	success tBTCSystemEthBtcPriceFeedAdditionStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemEthBtcPriceFeedAdditionStarted)
	eventSubscription, err := tbtcs.contract.WatchEthBtcPriceFeedAdditionStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for EthBtcPriceFeedAdditionStarted events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.PriceFeed,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemFundedFunc func(
	DepositContractAddress common.Address,
	Txid [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastFundedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) ([]*abi.TBTCSystemFunded, error) {
	iterator, err := tbtcs.contract.FilterFunded(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
		_txidFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Funded events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemFunded, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchFunded(
	success tBTCSystemFundedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeFunded(
			success,
			failCallback,
			_depositContractAddressFilter,
			_txidFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event Funded terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeFunded(
	success tBTCSystemFundedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemFunded)
	eventSubscription, err := tbtcs.contract.WatchFunded(
		nil,
		eventChan,
		_depositContractAddressFilter,
		_txidFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Funded events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Txid,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemGotRedemptionSignatureFunc func(
	DepositContractAddress common.Address,
	Digest [32]uint8,
	R [32]uint8,
	S [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastGotRedemptionSignatureEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) ([]*abi.TBTCSystemGotRedemptionSignature, error) {
	iterator, err := tbtcs.contract.FilterGotRedemptionSignature(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
		_digestFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past GotRedemptionSignature events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemGotRedemptionSignature, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchGotRedemptionSignature(
	success tBTCSystemGotRedemptionSignatureFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeGotRedemptionSignature(
			success,
			failCallback,
			_depositContractAddressFilter,
			_digestFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event GotRedemptionSignature terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeGotRedemptionSignature(
	success tBTCSystemGotRedemptionSignatureFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemGotRedemptionSignature)
	eventSubscription, err := tbtcs.contract.WatchGotRedemptionSignature(
		nil,
		eventChan,
		_depositContractAddressFilter,
		_digestFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for GotRedemptionSignature events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Digest,
					event.R,
					event.S,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemLiquidatedFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastLiquidatedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemLiquidated, error) {
	iterator, err := tbtcs.contract.FilterLiquidated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Liquidated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemLiquidated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchLiquidated(
	success tBTCSystemLiquidatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeLiquidated(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event Liquidated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeLiquidated(
	success tBTCSystemLiquidatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemLiquidated)
	eventSubscription, err := tbtcs.contract.WatchLiquidated(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Liquidated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemOwnershipTransferredFunc func(
	PreviousOwner common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastOwnershipTransferredEvents(
	startBlock uint64,
	endBlock *uint64,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) ([]*abi.TBTCSystemOwnershipTransferred, error) {
	iterator, err := tbtcs.contract.FilterOwnershipTransferred(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		previousOwnerFilter,
		newOwnerFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past OwnershipTransferred events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemOwnershipTransferred, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchOwnershipTransferred(
	success tBTCSystemOwnershipTransferredFunc,
	fail func(err error) error,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeOwnershipTransferred(
			success,
			failCallback,
			previousOwnerFilter,
			newOwnerFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event OwnershipTransferred terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeOwnershipTransferred(
	success tBTCSystemOwnershipTransferredFunc,
	fail func(err error) error,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemOwnershipTransferred)
	eventSubscription, err := tbtcs.contract.WatchOwnershipTransferred(
		nil,
		eventChan,
		previousOwnerFilter,
		newOwnerFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for OwnershipTransferred events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.PreviousOwner,
					event.NewOwner,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemRedeemedFunc func(
	DepositContractAddress common.Address,
	Txid [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastRedeemedEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) ([]*abi.TBTCSystemRedeemed, error) {
	iterator, err := tbtcs.contract.FilterRedeemed(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
		_txidFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Redeemed events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemRedeemed, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchRedeemed(
	success tBTCSystemRedeemedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeRedeemed(
			success,
			failCallback,
			_depositContractAddressFilter,
			_txidFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event Redeemed terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeRedeemed(
	success tBTCSystemRedeemedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemRedeemed)
	eventSubscription, err := tbtcs.contract.WatchRedeemed(
		nil,
		eventChan,
		_depositContractAddressFilter,
		_txidFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for Redeemed events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Txid,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemCollateralizationThresholdsUpdateStartedFunc func(
	InitialCollateralizedPercent uint16,
	UndercollateralizedThresholdPercent uint16,
	SeverelyUndercollateralizedThresholdPercent uint16,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastCollateralizationThresholdsUpdateStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemCollateralizationThresholdsUpdateStarted, error) {
	iterator, err := tbtcs.contract.FilterCollateralizationThresholdsUpdateStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past CollateralizationThresholdsUpdateStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemCollateralizationThresholdsUpdateStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchCollateralizationThresholdsUpdateStarted(
	success tBTCSystemCollateralizationThresholdsUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeCollateralizationThresholdsUpdateStarted(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event CollateralizationThresholdsUpdateStarted terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeCollateralizationThresholdsUpdateStarted(
	success tBTCSystemCollateralizationThresholdsUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemCollateralizationThresholdsUpdateStarted)
	eventSubscription, err := tbtcs.contract.WatchCollateralizationThresholdsUpdateStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for CollateralizationThresholdsUpdateStarted events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.InitialCollateralizedPercent,
					event.UndercollateralizedThresholdPercent,
					event.SeverelyUndercollateralizedThresholdPercent,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemCollateralizationThresholdsUpdatedFunc func(
	InitialCollateralizedPercent uint16,
	UndercollateralizedThresholdPercent uint16,
	SeverelyUndercollateralizedThresholdPercent uint16,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastCollateralizationThresholdsUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemCollateralizationThresholdsUpdated, error) {
	iterator, err := tbtcs.contract.FilterCollateralizationThresholdsUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past CollateralizationThresholdsUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemCollateralizationThresholdsUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchCollateralizationThresholdsUpdated(
	success tBTCSystemCollateralizationThresholdsUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeCollateralizationThresholdsUpdated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event CollateralizationThresholdsUpdated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeCollateralizationThresholdsUpdated(
	success tBTCSystemCollateralizationThresholdsUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemCollateralizationThresholdsUpdated)
	eventSubscription, err := tbtcs.contract.WatchCollateralizationThresholdsUpdated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for CollateralizationThresholdsUpdated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.InitialCollateralizedPercent,
					event.UndercollateralizedThresholdPercent,
					event.SeverelyUndercollateralizedThresholdPercent,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemEthBtcPriceFeedAddedFunc func(
	PriceFeed common.Address,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastEthBtcPriceFeedAddedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemEthBtcPriceFeedAdded, error) {
	iterator, err := tbtcs.contract.FilterEthBtcPriceFeedAdded(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past EthBtcPriceFeedAdded events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemEthBtcPriceFeedAdded, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchEthBtcPriceFeedAdded(
	success tBTCSystemEthBtcPriceFeedAddedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeEthBtcPriceFeedAdded(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event EthBtcPriceFeedAdded terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeEthBtcPriceFeedAdded(
	success tBTCSystemEthBtcPriceFeedAddedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemEthBtcPriceFeedAdded)
	eventSubscription, err := tbtcs.contract.WatchEthBtcPriceFeedAdded(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for EthBtcPriceFeedAdded events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.PriceFeed,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemExitedCourtesyCallFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastExitedCourtesyCallEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemExitedCourtesyCall, error) {
	iterator, err := tbtcs.contract.FilterExitedCourtesyCall(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past ExitedCourtesyCall events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemExitedCourtesyCall, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchExitedCourtesyCall(
	success tBTCSystemExitedCourtesyCallFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeExitedCourtesyCall(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event ExitedCourtesyCall terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeExitedCourtesyCall(
	success tBTCSystemExitedCourtesyCallFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemExitedCourtesyCall)
	eventSubscription, err := tbtcs.contract.WatchExitedCourtesyCall(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for ExitedCourtesyCall events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemSignerFeeDivisorUpdatedFunc func(
	SignerFeeDivisor uint16,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastSignerFeeDivisorUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemSignerFeeDivisorUpdated, error) {
	iterator, err := tbtcs.contract.FilterSignerFeeDivisorUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past SignerFeeDivisorUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemSignerFeeDivisorUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchSignerFeeDivisorUpdated(
	success tBTCSystemSignerFeeDivisorUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeSignerFeeDivisorUpdated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event SignerFeeDivisorUpdated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeSignerFeeDivisorUpdated(
	success tBTCSystemSignerFeeDivisorUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemSignerFeeDivisorUpdated)
	eventSubscription, err := tbtcs.contract.WatchSignerFeeDivisorUpdated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for SignerFeeDivisorUpdated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.SignerFeeDivisor,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemStartedLiquidationFunc func(
	DepositContractAddress common.Address,
	WasFraud bool,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastStartedLiquidationEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemStartedLiquidation, error) {
	iterator, err := tbtcs.contract.FilterStartedLiquidation(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past StartedLiquidation events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemStartedLiquidation, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchStartedLiquidation(
	success tBTCSystemStartedLiquidationFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeStartedLiquidation(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event StartedLiquidation terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeStartedLiquidation(
	success tBTCSystemStartedLiquidationFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemStartedLiquidation)
	eventSubscription, err := tbtcs.contract.WatchStartedLiquidation(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for StartedLiquidation events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.WasFraud,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemAllowNewDepositsUpdatedFunc func(
	AllowNewDeposits bool,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastAllowNewDepositsUpdatedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemAllowNewDepositsUpdated, error) {
	iterator, err := tbtcs.contract.FilterAllowNewDepositsUpdated(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past AllowNewDepositsUpdated events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemAllowNewDepositsUpdated, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchAllowNewDepositsUpdated(
	success tBTCSystemAllowNewDepositsUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeAllowNewDepositsUpdated(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event AllowNewDepositsUpdated terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeAllowNewDepositsUpdated(
	success tBTCSystemAllowNewDepositsUpdatedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemAllowNewDepositsUpdated)
	eventSubscription, err := tbtcs.contract.WatchAllowNewDepositsUpdated(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for AllowNewDepositsUpdated events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.AllowNewDeposits,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemFraudDuringSetupFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastFraudDuringSetupEvents(
	startBlock uint64,
	endBlock *uint64,
	_depositContractAddressFilter []common.Address,
) ([]*abi.TBTCSystemFraudDuringSetup, error) {
	iterator, err := tbtcs.contract.FilterFraudDuringSetup(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_depositContractAddressFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past FraudDuringSetup events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemFraudDuringSetup, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchFraudDuringSetup(
	success tBTCSystemFraudDuringSetupFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeFraudDuringSetup(
			success,
			failCallback,
			_depositContractAddressFilter,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event FraudDuringSetup terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeFraudDuringSetup(
	success tBTCSystemFraudDuringSetupFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemFraudDuringSetup)
	eventSubscription, err := tbtcs.contract.WatchFraudDuringSetup(
		nil,
		eventChan,
		_depositContractAddressFilter,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for FraudDuringSetup events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemLotSizesUpdateStartedFunc func(
	LotSizes []uint64,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastLotSizesUpdateStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemLotSizesUpdateStarted, error) {
	iterator, err := tbtcs.contract.FilterLotSizesUpdateStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past LotSizesUpdateStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemLotSizesUpdateStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchLotSizesUpdateStarted(
	success tBTCSystemLotSizesUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeLotSizesUpdateStarted(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event LotSizesUpdateStarted terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeLotSizesUpdateStarted(
	success tBTCSystemLotSizesUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemLotSizesUpdateStarted)
	eventSubscription, err := tbtcs.contract.WatchLotSizesUpdateStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for LotSizesUpdateStarted events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.LotSizes,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

type tBTCSystemSignerFeeDivisorUpdateStartedFunc func(
	SignerFeeDivisor uint16,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (tbtcs *TBTCSystem) PastSignerFeeDivisorUpdateStartedEvents(
	startBlock uint64,
	endBlock *uint64,
) ([]*abi.TBTCSystemSignerFeeDivisorUpdateStarted, error) {
	iterator, err := tbtcs.contract.FilterSignerFeeDivisorUpdateStarted(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past SignerFeeDivisorUpdateStarted events: [%v]",
			err,
		)
	}

	events := make([]*abi.TBTCSystemSignerFeeDivisorUpdateStarted, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (tbtcs *TBTCSystem) WatchSignerFeeDivisorUpdateStarted(
	success tBTCSystemSignerFeeDivisorUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	errorChan := make(chan error)
	unsubscribeChan := make(chan struct{})

	// Delay which must be preserved before a new resubscription attempt.
	// There is no sense to resubscribe immediately after the fail of current
	// subscription because the publisher must have some time to recover.
	retryDelay := 5 * time.Second

	watch := func() {
		failCallback := func(err error) error {
			fail(err)
			errorChan <- err // trigger resubscription signal
			return err
		}

		subscription, err := tbtcs.subscribeSignerFeeDivisorUpdateStarted(
			success,
			failCallback,
		)
		if err != nil {
			errorChan <- err // trigger resubscription signal
			return
		}

		// wait for unsubscription signal
		<-unsubscribeChan
		subscription.Unsubscribe()
	}

	// trigger the resubscriber goroutine
	go func() {
		go watch() // trigger first subscription

		for {
			select {
			case <-errorChan:
				tbtcsLogger.Warning(
					"subscription to event SignerFeeDivisorUpdateStarted terminated with error; " +
						"resubscription attempt will be performed after the retry delay",
				)
				time.Sleep(retryDelay)
				go watch()
			case <-unsubscribeChan:
				// shutdown the resubscriber goroutine on unsubscribe signal
				return
			}
		}
	}()

	// closing the unsubscribeChan will trigger a unsubscribe signal and
	// run unsubscription for all subscription instances
	unsubscribeCallback := func() {
		close(unsubscribeChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}

func (tbtcs *TBTCSystem) subscribeSignerFeeDivisorUpdateStarted(
	success tBTCSystemSignerFeeDivisorUpdateStartedFunc,
	fail func(err error) error,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.TBTCSystemSignerFeeDivisorUpdateStarted)
	eventSubscription, err := tbtcs.contract.WatchSignerFeeDivisorUpdateStarted(
		nil,
		eventChan,
	)
	if err != nil {
		close(eventChan)
		return eventSubscription, fmt.Errorf(
			"error creating watch for SignerFeeDivisorUpdateStarted events: [%v]",
			err,
		)
	}

	var subscriptionMutex = &sync.Mutex{}

	go func() {
		for {
			select {
			case event, subscribed := <-eventChan:
				subscriptionMutex.Lock()
				// if eventChan has been closed, it means we have unsubscribed
				if !subscribed {
					subscriptionMutex.Unlock()
					return
				}
				success(
					event.SignerFeeDivisor,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
				subscriptionMutex.Unlock()
			case ee := <-eventSubscription.Err():
				fail(ee)
				return
			}
		}
	}()

	unsubscribeCallback := func() {
		subscriptionMutex.Lock()
		defer subscriptionMutex.Unlock()

		eventSubscription.Unsubscribe()
		close(eventChan)
	}

	return subscription.NewEventSubscription(unsubscribeCallback), nil
}
