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
	"github.com/keep-network/tbtc/abi"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var dlLogger = log.Logger("keep-contract-DepositLog")

type DepositLog struct {
	contract          *abi.DepositLog
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

func NewDepositLog(
	contractAddress common.Address,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethutil.NonceManager,
	miningWaiter *ethutil.MiningWaiter,
	transactionMutex *sync.Mutex,
) (*DepositLog, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	transactorOptions := bind.NewKeyedTransactor(
		accountKey.PrivateKey,
	)

	randomBeaconContract, err := abi.NewDepositLog(
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

	contractABI, err := ethereumabi.JSON(strings.NewReader(abi.DepositLogABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &DepositLog{
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
func (dl *DepositLog) LogExitedCourtesyCall(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logExitedCourtesyCall",
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogExitedCourtesyCall(
		transactorOptions,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logExitedCourtesyCall",
		)
	}

	dlLogger.Infof(
		"submitted transaction logExitedCourtesyCall with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogExitedCourtesyCall(
				transactorOptions,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logExitedCourtesyCall",
				)
			}

			dlLogger.Infof(
				"submitted transaction logExitedCourtesyCall with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogExitedCourtesyCall(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logExitedCourtesyCall",
		&result,
	)

	return err
}

func (dl *DepositLog) LogExitedCourtesyCallGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logExitedCourtesyCall",
		dl.contractABI,
		dl.transactor,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogFraudDuringSetup(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logFraudDuringSetup",
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogFraudDuringSetup(
		transactorOptions,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logFraudDuringSetup",
		)
	}

	dlLogger.Infof(
		"submitted transaction logFraudDuringSetup with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogFraudDuringSetup(
				transactorOptions,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logFraudDuringSetup",
				)
			}

			dlLogger.Infof(
				"submitted transaction logFraudDuringSetup with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogFraudDuringSetup(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logFraudDuringSetup",
		&result,
	)

	return err
}

func (dl *DepositLog) LogFraudDuringSetupGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logFraudDuringSetup",
		dl.contractABI,
		dl.transactor,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogFunderRequestedAbort(
	_abortOutputScript []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logFunderRequestedAbort",
		"params: ",
		fmt.Sprint(
			_abortOutputScript,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogFunderRequestedAbort(
		transactorOptions,
		_abortOutputScript,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logFunderRequestedAbort",
			_abortOutputScript,
		)
	}

	dlLogger.Infof(
		"submitted transaction logFunderRequestedAbort with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogFunderRequestedAbort(
				transactorOptions,
				_abortOutputScript,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logFunderRequestedAbort",
					_abortOutputScript,
				)
			}

			dlLogger.Infof(
				"submitted transaction logFunderRequestedAbort with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogFunderRequestedAbort(
	_abortOutputScript []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logFunderRequestedAbort",
		&result,
		_abortOutputScript,
	)

	return err
}

func (dl *DepositLog) LogFunderRequestedAbortGasEstimate(
	_abortOutputScript []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logFunderRequestedAbort",
		dl.contractABI,
		dl.transactor,
		_abortOutputScript,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogFunded(
	_txid [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logFunded",
		"params: ",
		fmt.Sprint(
			_txid,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogFunded(
		transactorOptions,
		_txid,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logFunded",
			_txid,
		)
	}

	dlLogger.Infof(
		"submitted transaction logFunded with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogFunded(
				transactorOptions,
				_txid,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logFunded",
					_txid,
				)
			}

			dlLogger.Infof(
				"submitted transaction logFunded with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogFunded(
	_txid [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logFunded",
		&result,
		_txid,
	)

	return err
}

func (dl *DepositLog) LogFundedGasEstimate(
	_txid [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logFunded",
		dl.contractABI,
		dl.transactor,
		_txid,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogLiquidated(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logLiquidated",
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogLiquidated(
		transactorOptions,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logLiquidated",
		)
	}

	dlLogger.Infof(
		"submitted transaction logLiquidated with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogLiquidated(
				transactorOptions,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logLiquidated",
				)
			}

			dlLogger.Infof(
				"submitted transaction logLiquidated with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogLiquidated(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logLiquidated",
		&result,
	)

	return err
}

func (dl *DepositLog) LogLiquidatedGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logLiquidated",
		dl.contractABI,
		dl.transactor,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogRedemptionRequested(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
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

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogRedemptionRequested(
		transactorOptions,
		_requester,
		_digest,
		_utxoValue,
		_redeemerOutputScript,
		_requestedFee,
		_outpoint,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
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

	dlLogger.Infof(
		"submitted transaction logRedemptionRequested with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogRedemptionRequested(
				transactorOptions,
				_requester,
				_digest,
				_utxoValue,
				_redeemerOutputScript,
				_requestedFee,
				_outpoint,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
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

			dlLogger.Infof(
				"submitted transaction logRedemptionRequested with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogRedemptionRequested(
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
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
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

func (dl *DepositLog) LogRedemptionRequestedGasEstimate(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logRedemptionRequested",
		dl.contractABI,
		dl.transactor,
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
func (dl *DepositLog) LogSetupFailed(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logSetupFailed",
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogSetupFailed(
		transactorOptions,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logSetupFailed",
		)
	}

	dlLogger.Infof(
		"submitted transaction logSetupFailed with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogSetupFailed(
				transactorOptions,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logSetupFailed",
				)
			}

			dlLogger.Infof(
				"submitted transaction logSetupFailed with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogSetupFailed(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logSetupFailed",
		&result,
	)

	return err
}

func (dl *DepositLog) LogSetupFailedGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logSetupFailed",
		dl.contractABI,
		dl.transactor,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogRegisteredPubkey(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logRegisteredPubkey",
		"params: ",
		fmt.Sprint(
			_signingGroupPubkeyX,
			_signingGroupPubkeyY,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogRegisteredPubkey(
		transactorOptions,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logRegisteredPubkey",
			_signingGroupPubkeyX,
			_signingGroupPubkeyY,
		)
	}

	dlLogger.Infof(
		"submitted transaction logRegisteredPubkey with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogRegisteredPubkey(
				transactorOptions,
				_signingGroupPubkeyX,
				_signingGroupPubkeyY,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logRegisteredPubkey",
					_signingGroupPubkeyX,
					_signingGroupPubkeyY,
				)
			}

			dlLogger.Infof(
				"submitted transaction logRegisteredPubkey with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogRegisteredPubkey(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logRegisteredPubkey",
		&result,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)

	return err
}

func (dl *DepositLog) LogRegisteredPubkeyGasEstimate(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logRegisteredPubkey",
		dl.contractABI,
		dl.transactor,
		_signingGroupPubkeyX,
		_signingGroupPubkeyY,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogCourtesyCalled(

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logCourtesyCalled",
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogCourtesyCalled(
		transactorOptions,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logCourtesyCalled",
		)
	}

	dlLogger.Infof(
		"submitted transaction logCourtesyCalled with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogCourtesyCalled(
				transactorOptions,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logCourtesyCalled",
				)
			}

			dlLogger.Infof(
				"submitted transaction logCourtesyCalled with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogCourtesyCalled(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logCourtesyCalled",
		&result,
	)

	return err
}

func (dl *DepositLog) LogCourtesyCalledGasEstimate() (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logCourtesyCalled",
		dl.contractABI,
		dl.transactor,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogCreated(
	_keepAddress common.Address,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logCreated",
		"params: ",
		fmt.Sprint(
			_keepAddress,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogCreated(
		transactorOptions,
		_keepAddress,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logCreated",
			_keepAddress,
		)
	}

	dlLogger.Infof(
		"submitted transaction logCreated with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogCreated(
				transactorOptions,
				_keepAddress,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logCreated",
					_keepAddress,
				)
			}

			dlLogger.Infof(
				"submitted transaction logCreated with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogCreated(
	_keepAddress common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logCreated",
		&result,
		_keepAddress,
	)

	return err
}

func (dl *DepositLog) LogCreatedGasEstimate(
	_keepAddress common.Address,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logCreated",
		dl.contractABI,
		dl.transactor,
		_keepAddress,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogGotRedemptionSignature(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logGotRedemptionSignature",
		"params: ",
		fmt.Sprint(
			_digest,
			_r,
			_s,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogGotRedemptionSignature(
		transactorOptions,
		_digest,
		_r,
		_s,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logGotRedemptionSignature",
			_digest,
			_r,
			_s,
		)
	}

	dlLogger.Infof(
		"submitted transaction logGotRedemptionSignature with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogGotRedemptionSignature(
				transactorOptions,
				_digest,
				_r,
				_s,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logGotRedemptionSignature",
					_digest,
					_r,
					_s,
				)
			}

			dlLogger.Infof(
				"submitted transaction logGotRedemptionSignature with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogGotRedemptionSignature(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logGotRedemptionSignature",
		&result,
		_digest,
		_r,
		_s,
	)

	return err
}

func (dl *DepositLog) LogGotRedemptionSignatureGasEstimate(
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logGotRedemptionSignature",
		dl.contractABI,
		dl.transactor,
		_digest,
		_r,
		_s,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogRedeemed(
	_txid [32]uint8,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logRedeemed",
		"params: ",
		fmt.Sprint(
			_txid,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogRedeemed(
		transactorOptions,
		_txid,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logRedeemed",
			_txid,
		)
	}

	dlLogger.Infof(
		"submitted transaction logRedeemed with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogRedeemed(
				transactorOptions,
				_txid,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logRedeemed",
					_txid,
				)
			}

			dlLogger.Infof(
				"submitted transaction logRedeemed with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogRedeemed(
	_txid [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logRedeemed",
		&result,
		_txid,
	)

	return err
}

func (dl *DepositLog) LogRedeemedGasEstimate(
	_txid [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logRedeemed",
		dl.contractABI,
		dl.transactor,
		_txid,
	)

	return result, err
}

// Transaction submission.
func (dl *DepositLog) LogStartedLiquidation(
	_wasFraud bool,

	transactionOptions ...ethutil.TransactionOptions,
) (*types.Transaction, error) {
	dlLogger.Debug(
		"submitting transaction logStartedLiquidation",
		"params: ",
		fmt.Sprint(
			_wasFraud,
		),
	)

	dl.transactionMutex.Lock()
	defer dl.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *dl.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := dl.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := dl.contract.LogStartedLiquidation(
		transactorOptions,
		_wasFraud,
	)
	if err != nil {
		return transaction, dl.errorResolver.ResolveError(
			err,
			dl.transactorOptions.From,
			nil,
			"logStartedLiquidation",
			_wasFraud,
		)
	}

	dlLogger.Infof(
		"submitted transaction logStartedLiquidation with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go dl.miningWaiter.ForceMining(
		transaction,
		func(newGasPrice *big.Int) (*types.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := dl.contract.LogStartedLiquidation(
				transactorOptions,
				_wasFraud,
			)
			if err != nil {
				return transaction, dl.errorResolver.ResolveError(
					err,
					dl.transactorOptions.From,
					nil,
					"logStartedLiquidation",
					_wasFraud,
				)
			}

			dlLogger.Infof(
				"submitted transaction logStartedLiquidation with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	dl.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (dl *DepositLog) CallLogStartedLiquidation(
	_wasFraud bool,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := ethutil.CallAtBlock(
		dl.transactorOptions.From,
		blockNumber, nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"logStartedLiquidation",
		&result,
		_wasFraud,
	)

	return err
}

func (dl *DepositLog) LogStartedLiquidationGasEstimate(
	_wasFraud bool,
) (uint64, error) {
	var result uint64

	result, err := ethutil.EstimateGas(
		dl.callerOptions.From,
		dl.contractAddress,
		"logStartedLiquidation",
		dl.contractABI,
		dl.transactor,
		_wasFraud,
	)

	return result, err
}

// ----- Const Methods ------

func (dl *DepositLog) ApprovedToLog(
	_caller common.Address,
) (bool, error) {
	var result bool
	result, err := dl.contract.ApprovedToLog(
		dl.callerOptions,
		_caller,
	)

	if err != nil {
		return result, dl.errorResolver.ResolveError(
			err,
			dl.callerOptions.From,
			nil,
			"approvedToLog",
			_caller,
		)
	}

	return result, err
}

func (dl *DepositLog) ApprovedToLogAtBlock(
	_caller common.Address,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := ethutil.CallAtBlock(
		dl.callerOptions.From,
		blockNumber,
		nil,
		dl.contractABI,
		dl.caller,
		dl.errorResolver,
		dl.contractAddress,
		"approvedToLog",
		&result,
		_caller,
	)

	return result, err
}

// ------ Events -------

type depositLogCreatedFunc func(
	_depositContractAddress common.Address,
	_keepAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchCreated(
	success depositLogCreatedFunc,
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

		subscription, err := dl.subscribeCreated(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeCreated(
	success depositLogCreatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogCreated)
	eventSubscription, err := dl.contract.WatchCreated(
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
					event._depositContractAddress,
					event._keepAddress,
					event._timestamp,
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

type depositLogExitedCourtesyCallFunc func(
	_depositContractAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchExitedCourtesyCall(
	success depositLogExitedCourtesyCallFunc,
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

		subscription, err := dl.subscribeExitedCourtesyCall(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeExitedCourtesyCall(
	success depositLogExitedCourtesyCallFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogExitedCourtesyCall)
	eventSubscription, err := dl.contract.WatchExitedCourtesyCall(
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
					event._depositContractAddress,
					event._timestamp,
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

type depositLogGotRedemptionSignatureFunc func(
	_depositContractAddress common.Address,
	_digest [32]uint8,
	_r [32]uint8,
	_s [32]uint8,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchGotRedemptionSignature(
	success depositLogGotRedemptionSignatureFunc,
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

		subscription, err := dl.subscribeGotRedemptionSignature(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeGotRedemptionSignature(
	success depositLogGotRedemptionSignatureFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogGotRedemptionSignature)
	eventSubscription, err := dl.contract.WatchGotRedemptionSignature(
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
					event._depositContractAddress,
					event._digest,
					event._r,
					event._s,
					event._timestamp,
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

type depositLogLiquidatedFunc func(
	_depositContractAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchLiquidated(
	success depositLogLiquidatedFunc,
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

		subscription, err := dl.subscribeLiquidated(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeLiquidated(
	success depositLogLiquidatedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogLiquidated)
	eventSubscription, err := dl.contract.WatchLiquidated(
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
					event._depositContractAddress,
					event._timestamp,
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

type depositLogRedeemedFunc func(
	_depositContractAddress common.Address,
	_txid [32]uint8,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchRedeemed(
	success depositLogRedeemedFunc,
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

		subscription, err := dl.subscribeRedeemed(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeRedeemed(
	success depositLogRedeemedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogRedeemed)
	eventSubscription, err := dl.contract.WatchRedeemed(
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
					event._depositContractAddress,
					event._txid,
					event._timestamp,
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

type depositLogCourtesyCalledFunc func(
	_depositContractAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchCourtesyCalled(
	success depositLogCourtesyCalledFunc,
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

		subscription, err := dl.subscribeCourtesyCalled(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeCourtesyCalled(
	success depositLogCourtesyCalledFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogCourtesyCalled)
	eventSubscription, err := dl.contract.WatchCourtesyCalled(
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
					event._depositContractAddress,
					event._timestamp,
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

type depositLogFundedFunc func(
	_depositContractAddress common.Address,
	_txid [32]uint8,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchFunded(
	success depositLogFundedFunc,
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

		subscription, err := dl.subscribeFunded(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeFunded(
	success depositLogFundedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogFunded)
	eventSubscription, err := dl.contract.WatchFunded(
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
					event._depositContractAddress,
					event._txid,
					event._timestamp,
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

type depositLogFunderAbortRequestedFunc func(
	_depositContractAddress common.Address,
	_abortOutputScript []uint8,
	blockNumber uint64,
)

func (dl *DepositLog) WatchFunderAbortRequested(
	success depositLogFunderAbortRequestedFunc,
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

		subscription, err := dl.subscribeFunderAbortRequested(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeFunderAbortRequested(
	success depositLogFunderAbortRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogFunderAbortRequested)
	eventSubscription, err := dl.contract.WatchFunderAbortRequested(
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
					event._depositContractAddress,
					event._abortOutputScript,
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

type depositLogRedemptionRequestedFunc func(
	_depositContractAddress common.Address,
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,
	blockNumber uint64,
)

func (dl *DepositLog) WatchRedemptionRequested(
	success depositLogRedemptionRequestedFunc,
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

		subscription, err := dl.subscribeRedemptionRequested(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeRedemptionRequested(
	success depositLogRedemptionRequestedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogRedemptionRequested)
	eventSubscription, err := dl.contract.WatchRedemptionRequested(
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
					event._depositContractAddress,
					event._requester,
					event._digest,
					event._utxoValue,
					event._redeemerOutputScript,
					event._requestedFee,
					event._outpoint,
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

type depositLogRegisteredPubkeyFunc func(
	_depositContractAddress common.Address,
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchRegisteredPubkey(
	success depositLogRegisteredPubkeyFunc,
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

		subscription, err := dl.subscribeRegisteredPubkey(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeRegisteredPubkey(
	success depositLogRegisteredPubkeyFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogRegisteredPubkey)
	eventSubscription, err := dl.contract.WatchRegisteredPubkey(
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
					event._depositContractAddress,
					event._signingGroupPubkeyX,
					event._signingGroupPubkeyY,
					event._timestamp,
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

type depositLogSetupFailedFunc func(
	_depositContractAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchSetupFailed(
	success depositLogSetupFailedFunc,
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

		subscription, err := dl.subscribeSetupFailed(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeSetupFailed(
	success depositLogSetupFailedFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogSetupFailed)
	eventSubscription, err := dl.contract.WatchSetupFailed(
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
					event._depositContractAddress,
					event._timestamp,
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

type depositLogStartedLiquidationFunc func(
	_depositContractAddress common.Address,
	_wasFraud bool,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchStartedLiquidation(
	success depositLogStartedLiquidationFunc,
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

		subscription, err := dl.subscribeStartedLiquidation(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeStartedLiquidation(
	success depositLogStartedLiquidationFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogStartedLiquidation)
	eventSubscription, err := dl.contract.WatchStartedLiquidation(
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
					event._depositContractAddress,
					event._wasFraud,
					event._timestamp,
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

type depositLogFraudDuringSetupFunc func(
	_depositContractAddress common.Address,
	_timestamp *big.Int,
	blockNumber uint64,
)

func (dl *DepositLog) WatchFraudDuringSetup(
	success depositLogFraudDuringSetupFunc,
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

		subscription, err := dl.subscribeFraudDuringSetup(
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
				dlLogger.Warning(
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

func (dl *DepositLog) subscribeFraudDuringSetup(
	success depositLogFraudDuringSetupFunc,
	fail func(err error) error,
	_depositContractAddressFilter []common.Address,
) (subscription.EventSubscription, error) {
	eventChan := make(chan *abi.DepositLogFraudDuringSetup)
	eventSubscription, err := dl.contract.WatchFraudDuringSetup(
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
					event._depositContractAddress,
					event._timestamp,
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
