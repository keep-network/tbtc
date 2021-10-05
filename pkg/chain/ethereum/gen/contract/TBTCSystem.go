// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	hostchainabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"

	"github.com/ipfs/go-log"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/chain/ethlike"
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
	contractABI       *hostchainabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *chainutil.ErrorResolver
	nonceManager      *ethlike.NonceManager
	miningWaiter      *chainutil.MiningWaiter
	blockCounter      *ethlike.BlockCounter

	transactionMutex *sync.Mutex
}

func NewTBTCSystem(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*TBTCSystem, error) {
	callerOptions := &bind.CallOpts{
		From: accountKey.Address,
	}

	// FIXME Switch to bind.NewKeyedTransactorWithChainID when
	// FIXME celo-org/celo-blockchain merges in changes from upstream
	// FIXME ethereum/go-ethereum beyond v1.9.25.
	transactorOptions, err := chainutil.NewKeyedTransactorWithChainID(
		accountKey.PrivateKey,
		chainId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate transactor: [%v]", err)
	}

	contract, err := abi.NewTBTCSystem(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.TBTCSystemABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &TBTCSystem{
		contract:          contract,
		contractAddress:   contractAddress,
		contractABI:       &contractABI,
		caller:            backend,
		transactor:        backend,
		callerOptions:     callerOptions,
		transactorOptions: transactorOptions,
		errorResolver:     chainutil.NewErrorResolver(backend, &contractABI, &contractAddress),
		nonceManager:      nonceManager,
		miningWaiter:      miningWaiter,
		blockCounter:      blockCounter,
		transactionMutex:  transactionMutex,
	}, nil
}

// ----- Non-const Methods ------

// Transaction submission.
func (tbtcs *TBTCSystem) BeginCollateralizationThresholdsUpdate(
	_initialCollateralizedPercent uint16,
	_undercollateralizedThresholdPercent uint16,
	_severelyUndercollateralizedThresholdPercent uint16,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginCollateralizationThresholdsUpdate",
		" params: ",
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
		"submitted transaction beginCollateralizationThresholdsUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.BeginCollateralizationThresholdsUpdate(
				newTransactorOptions,
				_initialCollateralizedPercent,
				_undercollateralizedThresholdPercent,
				_severelyUndercollateralizedThresholdPercent,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
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
				"submitted transaction beginCollateralizationThresholdsUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginEthBtcPriceFeedAddition",
		" params: ",
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
		"submitted transaction beginEthBtcPriceFeedAddition with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.BeginEthBtcPriceFeedAddition(
				newTransactorOptions,
				_ethBtcPriceFeed,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginEthBtcPriceFeedAddition",
					_ethBtcPriceFeed,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginEthBtcPriceFeedAddition with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) BeginKeepFactoriesUpdate(
	_keepStakedFactory common.Address,
	_fullyBackedFactory common.Address,
	_factorySelector common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginKeepFactoriesUpdate",
		" params: ",
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
		"submitted transaction beginKeepFactoriesUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.BeginKeepFactoriesUpdate(
				newTransactorOptions,
				_keepStakedFactory,
				_fullyBackedFactory,
				_factorySelector,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
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
				"submitted transaction beginKeepFactoriesUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) BeginLotSizesUpdate(
	_lotSizes []uint64,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginLotSizesUpdate",
		" params: ",
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
		"submitted transaction beginLotSizesUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.BeginLotSizesUpdate(
				newTransactorOptions,
				_lotSizes,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginLotSizesUpdate",
					_lotSizes,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginLotSizesUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) BeginSignerFeeDivisorUpdate(
	_signerFeeDivisor uint16,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction beginSignerFeeDivisorUpdate",
		" params: ",
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
		"submitted transaction beginSignerFeeDivisorUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.BeginSignerFeeDivisorUpdate(
				newTransactorOptions,
				_signerFeeDivisor,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"beginSignerFeeDivisorUpdate",
					_signerFeeDivisor,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction beginSignerFeeDivisorUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) EmergencyPauseNewDeposits(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction emergencyPauseNewDeposits with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.EmergencyPauseNewDeposits(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"emergencyPauseNewDeposits",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction emergencyPauseNewDeposits with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"emergencyPauseNewDeposits",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeCollateralizationThresholdsUpdate(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction finalizeCollateralizationThresholdsUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.FinalizeCollateralizationThresholdsUpdate(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeCollateralizationThresholdsUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeCollateralizationThresholdsUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeCollateralizationThresholdsUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeEthBtcPriceFeedAddition(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction finalizeEthBtcPriceFeedAddition with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.FinalizeEthBtcPriceFeedAddition(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeEthBtcPriceFeedAddition",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeEthBtcPriceFeedAddition with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeEthBtcPriceFeedAddition",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeKeepFactoriesUpdate(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction finalizeKeepFactoriesUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.FinalizeKeepFactoriesUpdate(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeKeepFactoriesUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeKeepFactoriesUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeKeepFactoriesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeLotSizesUpdate(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction finalizeLotSizesUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.FinalizeLotSizesUpdate(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeLotSizesUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeLotSizesUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeLotSizesUpdate",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) FinalizeSignerFeeDivisorUpdate(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction finalizeSignerFeeDivisorUpdate with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.FinalizeSignerFeeDivisorUpdate(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"finalizeSignerFeeDivisorUpdate",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction finalizeSignerFeeDivisorUpdate with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"finalizeSignerFeeDivisorUpdate",
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

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction initialize",
		" params: ",
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
		"submitted transaction initialize with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.Initialize(
				newTransactorOptions,
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
				return nil, tbtcs.errorResolver.ResolveError(
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
				"submitted transaction initialize with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogCourtesyCalled(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction logCourtesyCalled with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogCourtesyCalled(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logCourtesyCalled",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logCourtesyCalled with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logCreated",
		" params: ",
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
		"submitted transaction logCreated with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogCreated(
				newTransactorOptions,
				_keepAddress,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logCreated",
					_keepAddress,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logCreated with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogExitedCourtesyCall(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction logExitedCourtesyCall with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogExitedCourtesyCall(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logExitedCourtesyCall",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logExitedCourtesyCall with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logExitedCourtesyCall",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogFraudDuringSetup(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction logFraudDuringSetup with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogFraudDuringSetup(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFraudDuringSetup",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFraudDuringSetup with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logFraudDuringSetup",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogFunded(
	_txid [32]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logFunded",
		" params: ",
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
		"submitted transaction logFunded with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogFunded(
				newTransactorOptions,
				_txid,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFunded",
					_txid,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFunded with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogFunderRequestedAbort(
	_abortOutputScript []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logFunderRequestedAbort",
		" params: ",
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
		"submitted transaction logFunderRequestedAbort with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogFunderRequestedAbort(
				newTransactorOptions,
				_abortOutputScript,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logFunderRequestedAbort",
					_abortOutputScript,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logFunderRequestedAbort with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logGotRedemptionSignature",
		" params: ",
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
		"submitted transaction logGotRedemptionSignature with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogGotRedemptionSignature(
				newTransactorOptions,
				_digest,
				_r,
				_s,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
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
				"submitted transaction logGotRedemptionSignature with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogLiquidated(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction logLiquidated with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogLiquidated(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logLiquidated",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logLiquidated with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logLiquidated",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogRedeemed(
	_txid [32]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRedeemed",
		" params: ",
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
		"submitted transaction logRedeemed with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogRedeemed(
				newTransactorOptions,
				_txid,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logRedeemed",
					_txid,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logRedeemed with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogRedemptionRequested(
	_requester common.Address,
	_digest [32]uint8,
	_utxoValue *big.Int,
	_redeemerOutputScript []uint8,
	_requestedFee *big.Int,
	_outpoint []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRedemptionRequested",
		" params: ",
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
		"submitted transaction logRedemptionRequested with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogRedemptionRequested(
				newTransactorOptions,
				_requester,
				_digest,
				_utxoValue,
				_redeemerOutputScript,
				_requestedFee,
				_outpoint,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
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
				"submitted transaction logRedemptionRequested with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogRegisteredPubkey(
	_signingGroupPubkeyX [32]uint8,
	_signingGroupPubkeyY [32]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logRegisteredPubkey",
		" params: ",
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
		"submitted transaction logRegisteredPubkey with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogRegisteredPubkey(
				newTransactorOptions,
				_signingGroupPubkeyX,
				_signingGroupPubkeyY,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logRegisteredPubkey",
					_signingGroupPubkeyX,
					_signingGroupPubkeyY,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logRegisteredPubkey with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) LogSetupFailed(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction logSetupFailed with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogSetupFailed(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logSetupFailed",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logSetupFailed with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"logSetupFailed",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) LogStartedLiquidation(
	_wasFraud bool,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction logStartedLiquidation",
		" params: ",
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
		"submitted transaction logStartedLiquidation with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.LogStartedLiquidation(
				newTransactorOptions,
				_wasFraud,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"logStartedLiquidation",
					_wasFraud,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction logStartedLiquidation with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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
func (tbtcs *TBTCSystem) RefreshMinimumBondableValue(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction refreshMinimumBondableValue with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.RefreshMinimumBondableValue(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"refreshMinimumBondableValue",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction refreshMinimumBondableValue with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"refreshMinimumBondableValue",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) RenounceOwnership(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.RenounceOwnership(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"renounceOwnership",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction renounceOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"renounceOwnership",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) RequestNewKeep(
	_requestedLotSizeSatoshis uint64,
	_maxSecuredLifetime *big.Int,
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction requestNewKeep",
		" params: ",
		fmt.Sprint(
			_requestedLotSizeSatoshis,
			_maxSecuredLifetime,
		),
		" value: ", value,
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
		"submitted transaction requestNewKeep with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.RequestNewKeep(
				newTransactorOptions,
				_requestedLotSizeSatoshis,
				_maxSecuredLifetime,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					value,
					"requestNewKeep",
					_requestedLotSizeSatoshis,
					_maxSecuredLifetime,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction requestNewKeep with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
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

// Transaction submission.
func (tbtcs *TBTCSystem) ResumeNewDeposits(

	transactionOptions ...chainutil.TransactionOptions,
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
		"submitted transaction resumeNewDeposits with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.ResumeNewDeposits(
				newTransactorOptions,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"resumeNewDeposits",
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction resumeNewDeposits with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"resumeNewDeposits",
		tbtcs.contractABI,
		tbtcs.transactor,
	)

	return result, err
}

// Transaction submission.
func (tbtcs *TBTCSystem) TransferOwnership(
	newOwner common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	tbtcsLogger.Debug(
		"submitting transaction transferOwnership",
		" params: ",
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
		"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go tbtcs.miningWaiter.ForceMining(
		transaction,
		transactorOptions,
		func(newTransactorOptions *bind.TransactOpts) (*types.Transaction, error) {
			// If original transactor options has a non-zero gas limit, that
			// means the client code set it on their own. In that case, we
			// should rewrite the gas limit from the original transaction
			// for each resubmission. If the gas limit is not set by the client
			// code, let the the submitter re-estimate the gas limit on each
			// resubmission.
			if transactorOptions.GasLimit != 0 {
				newTransactorOptions.GasLimit = transactorOptions.GasLimit
			}

			transaction, err := tbtcs.contract.TransferOwnership(
				newTransactorOptions,
				newOwner,
			)
			if err != nil {
				return nil, tbtcs.errorResolver.ResolveError(
					err,
					tbtcs.transactorOptions.From,
					nil,
					"transferOwnership",
					newOwner,
				)
			}

			tbtcsLogger.Infof(
				"submitted transaction transferOwnership with id: [%s] and nonce [%v]",
				transaction.Hash(),
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

	err := chainutil.CallAtBlock(
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

	result, err := chainutil.EstimateGas(
		tbtcs.callerOptions.From,
		tbtcs.contractAddress,
		"transferOwnership",
		tbtcs.contractABI,
		tbtcs.transactor,
		newOwner,
	)

	return result, err
}

// ----- Const Methods ------

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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

	err := chainutil.CallAtBlock(
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

// ------ Events -------

func (tbtcs *TBTCSystem) AllowNewDepositsUpdated(
	opts *ethlike.SubscribeOpts,
) *TbtcsAllowNewDepositsUpdatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsAllowNewDepositsUpdatedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsAllowNewDepositsUpdatedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemAllowNewDepositsUpdatedFunc func(
	AllowNewDeposits bool,
	blockNumber uint64,
)

func (andus *TbtcsAllowNewDepositsUpdatedSubscription) OnEvent(
	handler tBTCSystemAllowNewDepositsUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemAllowNewDepositsUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.AllowNewDeposits,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := andus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (andus *TbtcsAllowNewDepositsUpdatedSubscription) Pipe(
	sink chan *abi.TBTCSystemAllowNewDepositsUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(andus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := andus.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - andus.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past AllowNewDepositsUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := andus.contract.PastAllowNewDepositsUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past AllowNewDepositsUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := andus.contract.watchAllowNewDepositsUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchAllowNewDepositsUpdated(
	sink chan *abi.TBTCSystemAllowNewDepositsUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchAllowNewDepositsUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event AllowNewDepositsUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event AllowNewDepositsUpdated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) CollateralizationThresholdsUpdateStarted(
	opts *ethlike.SubscribeOpts,
) *TbtcsCollateralizationThresholdsUpdateStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsCollateralizationThresholdsUpdateStartedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsCollateralizationThresholdsUpdateStartedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemCollateralizationThresholdsUpdateStartedFunc func(
	InitialCollateralizedPercent uint16,
	UndercollateralizedThresholdPercent uint16,
	SeverelyUndercollateralizedThresholdPercent uint16,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (ctuss *TbtcsCollateralizationThresholdsUpdateStartedSubscription) OnEvent(
	handler tBTCSystemCollateralizationThresholdsUpdateStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemCollateralizationThresholdsUpdateStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.InitialCollateralizedPercent,
					event.UndercollateralizedThresholdPercent,
					event.SeverelyUndercollateralizedThresholdPercent,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ctuss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ctuss *TbtcsCollateralizationThresholdsUpdateStartedSubscription) Pipe(
	sink chan *abi.TBTCSystemCollateralizationThresholdsUpdateStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ctuss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ctuss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ctuss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past CollateralizationThresholdsUpdateStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ctuss.contract.PastCollateralizationThresholdsUpdateStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past CollateralizationThresholdsUpdateStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ctuss.contract.watchCollateralizationThresholdsUpdateStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchCollateralizationThresholdsUpdateStarted(
	sink chan *abi.TBTCSystemCollateralizationThresholdsUpdateStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchCollateralizationThresholdsUpdateStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event CollateralizationThresholdsUpdateStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event CollateralizationThresholdsUpdateStarted failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) CollateralizationThresholdsUpdated(
	opts *ethlike.SubscribeOpts,
) *TbtcsCollateralizationThresholdsUpdatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsCollateralizationThresholdsUpdatedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsCollateralizationThresholdsUpdatedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemCollateralizationThresholdsUpdatedFunc func(
	InitialCollateralizedPercent uint16,
	UndercollateralizedThresholdPercent uint16,
	SeverelyUndercollateralizedThresholdPercent uint16,
	blockNumber uint64,
)

func (ctus *TbtcsCollateralizationThresholdsUpdatedSubscription) OnEvent(
	handler tBTCSystemCollateralizationThresholdsUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemCollateralizationThresholdsUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.InitialCollateralizedPercent,
					event.UndercollateralizedThresholdPercent,
					event.SeverelyUndercollateralizedThresholdPercent,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ctus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ctus *TbtcsCollateralizationThresholdsUpdatedSubscription) Pipe(
	sink chan *abi.TBTCSystemCollateralizationThresholdsUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ctus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ctus.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ctus.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past CollateralizationThresholdsUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ctus.contract.PastCollateralizationThresholdsUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past CollateralizationThresholdsUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ctus.contract.watchCollateralizationThresholdsUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchCollateralizationThresholdsUpdated(
	sink chan *abi.TBTCSystemCollateralizationThresholdsUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchCollateralizationThresholdsUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event CollateralizationThresholdsUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event CollateralizationThresholdsUpdated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) CourtesyCalled(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsCourtesyCalledSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsCourtesyCalledSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsCourtesyCalledSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemCourtesyCalledFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (ccs *TbtcsCourtesyCalledSubscription) OnEvent(
	handler tBTCSystemCourtesyCalledFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemCourtesyCalled)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ccs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ccs *TbtcsCourtesyCalledSubscription) Pipe(
	sink chan *abi.TBTCSystemCourtesyCalled,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ccs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ccs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ccs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past CourtesyCalled events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ccs.contract.PastCourtesyCalledEvents(
					fromBlock,
					nil,
					ccs._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past CourtesyCalled events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ccs.contract.watchCourtesyCalled(
		sink,
		ccs._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchCourtesyCalled(
	sink chan *abi.TBTCSystemCourtesyCalled,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchCourtesyCalled(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event CourtesyCalled had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event CourtesyCalled failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) Created(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) *TbtcsCreatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsCreatedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
		_keepAddressFilter,
	}
}

type TbtcsCreatedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
	_keepAddressFilter            []common.Address
}

type tBTCSystemCreatedFunc func(
	DepositContractAddress common.Address,
	KeepAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (cs *TbtcsCreatedSubscription) OnEvent(
	handler tBTCSystemCreatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemCreated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.KeepAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := cs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (cs *TbtcsCreatedSubscription) Pipe(
	sink chan *abi.TBTCSystemCreated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(cs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := cs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - cs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past Created events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := cs.contract.PastCreatedEvents(
					fromBlock,
					nil,
					cs._depositContractAddressFilter,
					cs._keepAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past Created events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := cs.contract.watchCreated(
		sink,
		cs._depositContractAddressFilter,
		cs._keepAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchCreated(
	sink chan *abi.TBTCSystemCreated,
	_depositContractAddressFilter []common.Address,
	_keepAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchCreated(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
			_keepAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event Created had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event Created failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) EthBtcPriceFeedAdded(
	opts *ethlike.SubscribeOpts,
) *TbtcsEthBtcPriceFeedAddedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsEthBtcPriceFeedAddedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsEthBtcPriceFeedAddedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemEthBtcPriceFeedAddedFunc func(
	PriceFeed common.Address,
	blockNumber uint64,
)

func (ebpfas *TbtcsEthBtcPriceFeedAddedSubscription) OnEvent(
	handler tBTCSystemEthBtcPriceFeedAddedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemEthBtcPriceFeedAdded)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.PriceFeed,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ebpfas.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ebpfas *TbtcsEthBtcPriceFeedAddedSubscription) Pipe(
	sink chan *abi.TBTCSystemEthBtcPriceFeedAdded,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ebpfas.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ebpfas.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ebpfas.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past EthBtcPriceFeedAdded events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ebpfas.contract.PastEthBtcPriceFeedAddedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past EthBtcPriceFeedAdded events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ebpfas.contract.watchEthBtcPriceFeedAdded(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchEthBtcPriceFeedAdded(
	sink chan *abi.TBTCSystemEthBtcPriceFeedAdded,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchEthBtcPriceFeedAdded(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event EthBtcPriceFeedAdded had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event EthBtcPriceFeedAdded failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) EthBtcPriceFeedAdditionStarted(
	opts *ethlike.SubscribeOpts,
) *TbtcsEthBtcPriceFeedAdditionStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsEthBtcPriceFeedAdditionStartedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsEthBtcPriceFeedAdditionStartedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemEthBtcPriceFeedAdditionStartedFunc func(
	PriceFeed common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (ebpfass *TbtcsEthBtcPriceFeedAdditionStartedSubscription) OnEvent(
	handler tBTCSystemEthBtcPriceFeedAdditionStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemEthBtcPriceFeedAdditionStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.PriceFeed,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ebpfass.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ebpfass *TbtcsEthBtcPriceFeedAdditionStartedSubscription) Pipe(
	sink chan *abi.TBTCSystemEthBtcPriceFeedAdditionStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ebpfass.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ebpfass.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ebpfass.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past EthBtcPriceFeedAdditionStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ebpfass.contract.PastEthBtcPriceFeedAdditionStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past EthBtcPriceFeedAdditionStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ebpfass.contract.watchEthBtcPriceFeedAdditionStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchEthBtcPriceFeedAdditionStarted(
	sink chan *abi.TBTCSystemEthBtcPriceFeedAdditionStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchEthBtcPriceFeedAdditionStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event EthBtcPriceFeedAdditionStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event EthBtcPriceFeedAdditionStarted failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) ExitedCourtesyCall(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsExitedCourtesyCallSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsExitedCourtesyCallSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsExitedCourtesyCallSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemExitedCourtesyCallFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (eccs *TbtcsExitedCourtesyCallSubscription) OnEvent(
	handler tBTCSystemExitedCourtesyCallFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemExitedCourtesyCall)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := eccs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (eccs *TbtcsExitedCourtesyCallSubscription) Pipe(
	sink chan *abi.TBTCSystemExitedCourtesyCall,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(eccs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := eccs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - eccs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past ExitedCourtesyCall events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := eccs.contract.PastExitedCourtesyCallEvents(
					fromBlock,
					nil,
					eccs._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past ExitedCourtesyCall events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := eccs.contract.watchExitedCourtesyCall(
		sink,
		eccs._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchExitedCourtesyCall(
	sink chan *abi.TBTCSystemExitedCourtesyCall,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchExitedCourtesyCall(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event ExitedCourtesyCall had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event ExitedCourtesyCall failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) FraudDuringSetup(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsFraudDuringSetupSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsFraudDuringSetupSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsFraudDuringSetupSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemFraudDuringSetupFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (fdss *TbtcsFraudDuringSetupSubscription) OnEvent(
	handler tBTCSystemFraudDuringSetupFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemFraudDuringSetup)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := fdss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (fdss *TbtcsFraudDuringSetupSubscription) Pipe(
	sink chan *abi.TBTCSystemFraudDuringSetup,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(fdss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := fdss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - fdss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past FraudDuringSetup events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := fdss.contract.PastFraudDuringSetupEvents(
					fromBlock,
					nil,
					fdss._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past FraudDuringSetup events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := fdss.contract.watchFraudDuringSetup(
		sink,
		fdss._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchFraudDuringSetup(
	sink chan *abi.TBTCSystemFraudDuringSetup,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchFraudDuringSetup(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event FraudDuringSetup had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event FraudDuringSetup failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) Funded(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) *TbtcsFundedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsFundedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
		_txidFilter,
	}
}

type TbtcsFundedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
	_txidFilter                   [][32]uint8
}

type tBTCSystemFundedFunc func(
	DepositContractAddress common.Address,
	Txid [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (fs *TbtcsFundedSubscription) OnEvent(
	handler tBTCSystemFundedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemFunded)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Txid,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := fs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (fs *TbtcsFundedSubscription) Pipe(
	sink chan *abi.TBTCSystemFunded,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(fs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := fs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - fs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past Funded events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := fs.contract.PastFundedEvents(
					fromBlock,
					nil,
					fs._depositContractAddressFilter,
					fs._txidFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past Funded events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := fs.contract.watchFunded(
		sink,
		fs._depositContractAddressFilter,
		fs._txidFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchFunded(
	sink chan *abi.TBTCSystemFunded,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchFunded(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
			_txidFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event Funded had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event Funded failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) FunderAbortRequested(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsFunderAbortRequestedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsFunderAbortRequestedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsFunderAbortRequestedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemFunderAbortRequestedFunc func(
	DepositContractAddress common.Address,
	AbortOutputScript []uint8,
	blockNumber uint64,
)

func (fars *TbtcsFunderAbortRequestedSubscription) OnEvent(
	handler tBTCSystemFunderAbortRequestedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemFunderAbortRequested)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.AbortOutputScript,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := fars.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (fars *TbtcsFunderAbortRequestedSubscription) Pipe(
	sink chan *abi.TBTCSystemFunderAbortRequested,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(fars.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := fars.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - fars.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past FunderAbortRequested events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := fars.contract.PastFunderAbortRequestedEvents(
					fromBlock,
					nil,
					fars._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past FunderAbortRequested events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := fars.contract.watchFunderAbortRequested(
		sink,
		fars._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchFunderAbortRequested(
	sink chan *abi.TBTCSystemFunderAbortRequested,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchFunderAbortRequested(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event FunderAbortRequested had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event FunderAbortRequested failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) GotRedemptionSignature(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) *TbtcsGotRedemptionSignatureSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsGotRedemptionSignatureSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
		_digestFilter,
	}
}

type TbtcsGotRedemptionSignatureSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
	_digestFilter                 [][32]uint8
}

type tBTCSystemGotRedemptionSignatureFunc func(
	DepositContractAddress common.Address,
	Digest [32]uint8,
	R [32]uint8,
	S [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (grss *TbtcsGotRedemptionSignatureSubscription) OnEvent(
	handler tBTCSystemGotRedemptionSignatureFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemGotRedemptionSignature)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Digest,
					event.R,
					event.S,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := grss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (grss *TbtcsGotRedemptionSignatureSubscription) Pipe(
	sink chan *abi.TBTCSystemGotRedemptionSignature,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(grss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := grss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - grss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past GotRedemptionSignature events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := grss.contract.PastGotRedemptionSignatureEvents(
					fromBlock,
					nil,
					grss._depositContractAddressFilter,
					grss._digestFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past GotRedemptionSignature events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := grss.contract.watchGotRedemptionSignature(
		sink,
		grss._depositContractAddressFilter,
		grss._digestFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchGotRedemptionSignature(
	sink chan *abi.TBTCSystemGotRedemptionSignature,
	_depositContractAddressFilter []common.Address,
	_digestFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchGotRedemptionSignature(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
			_digestFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event GotRedemptionSignature had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event GotRedemptionSignature failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) KeepFactoriesUpdateStarted(
	opts *ethlike.SubscribeOpts,
) *TbtcsKeepFactoriesUpdateStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsKeepFactoriesUpdateStartedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsKeepFactoriesUpdateStartedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemKeepFactoriesUpdateStartedFunc func(
	KeepStakedFactory common.Address,
	FullyBackedFactory common.Address,
	FactorySelector common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (kfuss *TbtcsKeepFactoriesUpdateStartedSubscription) OnEvent(
	handler tBTCSystemKeepFactoriesUpdateStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemKeepFactoriesUpdateStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.KeepStakedFactory,
					event.FullyBackedFactory,
					event.FactorySelector,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := kfuss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (kfuss *TbtcsKeepFactoriesUpdateStartedSubscription) Pipe(
	sink chan *abi.TBTCSystemKeepFactoriesUpdateStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(kfuss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := kfuss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - kfuss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past KeepFactoriesUpdateStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := kfuss.contract.PastKeepFactoriesUpdateStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past KeepFactoriesUpdateStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := kfuss.contract.watchKeepFactoriesUpdateStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchKeepFactoriesUpdateStarted(
	sink chan *abi.TBTCSystemKeepFactoriesUpdateStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchKeepFactoriesUpdateStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event KeepFactoriesUpdateStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event KeepFactoriesUpdateStarted failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) KeepFactoriesUpdated(
	opts *ethlike.SubscribeOpts,
) *TbtcsKeepFactoriesUpdatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsKeepFactoriesUpdatedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsKeepFactoriesUpdatedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemKeepFactoriesUpdatedFunc func(
	KeepStakedFactory common.Address,
	FullyBackedFactory common.Address,
	FactorySelector common.Address,
	blockNumber uint64,
)

func (kfus *TbtcsKeepFactoriesUpdatedSubscription) OnEvent(
	handler tBTCSystemKeepFactoriesUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemKeepFactoriesUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.KeepStakedFactory,
					event.FullyBackedFactory,
					event.FactorySelector,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := kfus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (kfus *TbtcsKeepFactoriesUpdatedSubscription) Pipe(
	sink chan *abi.TBTCSystemKeepFactoriesUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(kfus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := kfus.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - kfus.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past KeepFactoriesUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := kfus.contract.PastKeepFactoriesUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past KeepFactoriesUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := kfus.contract.watchKeepFactoriesUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchKeepFactoriesUpdated(
	sink chan *abi.TBTCSystemKeepFactoriesUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchKeepFactoriesUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event KeepFactoriesUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event KeepFactoriesUpdated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) Liquidated(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsLiquidatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsLiquidatedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsLiquidatedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemLiquidatedFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (ls *TbtcsLiquidatedSubscription) OnEvent(
	handler tBTCSystemLiquidatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemLiquidated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ls.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ls *TbtcsLiquidatedSubscription) Pipe(
	sink chan *abi.TBTCSystemLiquidated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ls.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ls.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ls.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past Liquidated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ls.contract.PastLiquidatedEvents(
					fromBlock,
					nil,
					ls._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past Liquidated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ls.contract.watchLiquidated(
		sink,
		ls._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchLiquidated(
	sink chan *abi.TBTCSystemLiquidated,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchLiquidated(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event Liquidated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event Liquidated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) LotSizesUpdateStarted(
	opts *ethlike.SubscribeOpts,
) *TbtcsLotSizesUpdateStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsLotSizesUpdateStartedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsLotSizesUpdateStartedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemLotSizesUpdateStartedFunc func(
	LotSizes []uint64,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (lsuss *TbtcsLotSizesUpdateStartedSubscription) OnEvent(
	handler tBTCSystemLotSizesUpdateStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemLotSizesUpdateStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.LotSizes,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := lsuss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (lsuss *TbtcsLotSizesUpdateStartedSubscription) Pipe(
	sink chan *abi.TBTCSystemLotSizesUpdateStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(lsuss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := lsuss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - lsuss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past LotSizesUpdateStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := lsuss.contract.PastLotSizesUpdateStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past LotSizesUpdateStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := lsuss.contract.watchLotSizesUpdateStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchLotSizesUpdateStarted(
	sink chan *abi.TBTCSystemLotSizesUpdateStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchLotSizesUpdateStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event LotSizesUpdateStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event LotSizesUpdateStarted failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) LotSizesUpdated(
	opts *ethlike.SubscribeOpts,
) *TbtcsLotSizesUpdatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsLotSizesUpdatedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsLotSizesUpdatedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemLotSizesUpdatedFunc func(
	LotSizes []uint64,
	blockNumber uint64,
)

func (lsus *TbtcsLotSizesUpdatedSubscription) OnEvent(
	handler tBTCSystemLotSizesUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemLotSizesUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.LotSizes,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := lsus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (lsus *TbtcsLotSizesUpdatedSubscription) Pipe(
	sink chan *abi.TBTCSystemLotSizesUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(lsus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := lsus.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - lsus.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past LotSizesUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := lsus.contract.PastLotSizesUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past LotSizesUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := lsus.contract.watchLotSizesUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchLotSizesUpdated(
	sink chan *abi.TBTCSystemLotSizesUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchLotSizesUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event LotSizesUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event LotSizesUpdated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) OwnershipTransferred(
	opts *ethlike.SubscribeOpts,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) *TbtcsOwnershipTransferredSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsOwnershipTransferredSubscription{
		tbtcs,
		opts,
		previousOwnerFilter,
		newOwnerFilter,
	}
}

type TbtcsOwnershipTransferredSubscription struct {
	contract            *TBTCSystem
	opts                *ethlike.SubscribeOpts
	previousOwnerFilter []common.Address
	newOwnerFilter      []common.Address
}

type tBTCSystemOwnershipTransferredFunc func(
	PreviousOwner common.Address,
	NewOwner common.Address,
	blockNumber uint64,
)

func (ots *TbtcsOwnershipTransferredSubscription) OnEvent(
	handler tBTCSystemOwnershipTransferredFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemOwnershipTransferred)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.PreviousOwner,
					event.NewOwner,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := ots.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (ots *TbtcsOwnershipTransferredSubscription) Pipe(
	sink chan *abi.TBTCSystemOwnershipTransferred,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(ots.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := ots.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - ots.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past OwnershipTransferred events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := ots.contract.PastOwnershipTransferredEvents(
					fromBlock,
					nil,
					ots.previousOwnerFilter,
					ots.newOwnerFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past OwnershipTransferred events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := ots.contract.watchOwnershipTransferred(
		sink,
		ots.previousOwnerFilter,
		ots.newOwnerFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchOwnershipTransferred(
	sink chan *abi.TBTCSystemOwnershipTransferred,
	previousOwnerFilter []common.Address,
	newOwnerFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchOwnershipTransferred(
			&bind.WatchOpts{Context: ctx},
			sink,
			previousOwnerFilter,
			newOwnerFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event OwnershipTransferred had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event OwnershipTransferred failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) Redeemed(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) *TbtcsRedeemedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsRedeemedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
		_txidFilter,
	}
}

type TbtcsRedeemedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
	_txidFilter                   [][32]uint8
}

type tBTCSystemRedeemedFunc func(
	DepositContractAddress common.Address,
	Txid [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (rs *TbtcsRedeemedSubscription) OnEvent(
	handler tBTCSystemRedeemedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemRedeemed)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Txid,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rs *TbtcsRedeemedSubscription) Pipe(
	sink chan *abi.TBTCSystemRedeemed,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past Redeemed events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rs.contract.PastRedeemedEvents(
					fromBlock,
					nil,
					rs._depositContractAddressFilter,
					rs._txidFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past Redeemed events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rs.contract.watchRedeemed(
		sink,
		rs._depositContractAddressFilter,
		rs._txidFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchRedeemed(
	sink chan *abi.TBTCSystemRedeemed,
	_depositContractAddressFilter []common.Address,
	_txidFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchRedeemed(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
			_txidFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event Redeemed had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event Redeemed failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) RedemptionRequested(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) *TbtcsRedemptionRequestedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsRedemptionRequestedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
		_requesterFilter,
		_digestFilter,
	}
}

type TbtcsRedemptionRequestedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
	_requesterFilter              []common.Address
	_digestFilter                 [][32]uint8
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

func (rrs *TbtcsRedemptionRequestedSubscription) OnEvent(
	handler tBTCSystemRedemptionRequestedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemRedemptionRequested)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Requester,
					event.Digest,
					event.UtxoValue,
					event.RedeemerOutputScript,
					event.RequestedFee,
					event.Outpoint,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rrs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rrs *TbtcsRedemptionRequestedSubscription) Pipe(
	sink chan *abi.TBTCSystemRedemptionRequested,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rrs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rrs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rrs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past RedemptionRequested events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rrs.contract.PastRedemptionRequestedEvents(
					fromBlock,
					nil,
					rrs._depositContractAddressFilter,
					rrs._requesterFilter,
					rrs._digestFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past RedemptionRequested events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rrs.contract.watchRedemptionRequested(
		sink,
		rrs._depositContractAddressFilter,
		rrs._requesterFilter,
		rrs._digestFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchRedemptionRequested(
	sink chan *abi.TBTCSystemRedemptionRequested,
	_depositContractAddressFilter []common.Address,
	_requesterFilter []common.Address,
	_digestFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchRedemptionRequested(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
			_requesterFilter,
			_digestFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event RedemptionRequested had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event RedemptionRequested failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) RegisteredPubkey(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsRegisteredPubkeySubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsRegisteredPubkeySubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsRegisteredPubkeySubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemRegisteredPubkeyFunc func(
	DepositContractAddress common.Address,
	SigningGroupPubkeyX [32]uint8,
	SigningGroupPubkeyY [32]uint8,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (rps *TbtcsRegisteredPubkeySubscription) OnEvent(
	handler tBTCSystemRegisteredPubkeyFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemRegisteredPubkey)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.SigningGroupPubkeyX,
					event.SigningGroupPubkeyY,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := rps.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (rps *TbtcsRegisteredPubkeySubscription) Pipe(
	sink chan *abi.TBTCSystemRegisteredPubkey,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(rps.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := rps.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - rps.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past RegisteredPubkey events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := rps.contract.PastRegisteredPubkeyEvents(
					fromBlock,
					nil,
					rps._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past RegisteredPubkey events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := rps.contract.watchRegisteredPubkey(
		sink,
		rps._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchRegisteredPubkey(
	sink chan *abi.TBTCSystemRegisteredPubkey,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchRegisteredPubkey(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event RegisteredPubkey had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event RegisteredPubkey failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) SetupFailed(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsSetupFailedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsSetupFailedSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsSetupFailedSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemSetupFailedFunc func(
	DepositContractAddress common.Address,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (sfs *TbtcsSetupFailedSubscription) OnEvent(
	handler tBTCSystemSetupFailedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemSetupFailed)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sfs.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sfs *TbtcsSetupFailedSubscription) Pipe(
	sink chan *abi.TBTCSystemSetupFailed,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sfs.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sfs.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sfs.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past SetupFailed events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sfs.contract.PastSetupFailedEvents(
					fromBlock,
					nil,
					sfs._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past SetupFailed events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sfs.contract.watchSetupFailed(
		sink,
		sfs._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchSetupFailed(
	sink chan *abi.TBTCSystemSetupFailed,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchSetupFailed(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event SetupFailed had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event SetupFailed failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) SignerFeeDivisorUpdateStarted(
	opts *ethlike.SubscribeOpts,
) *TbtcsSignerFeeDivisorUpdateStartedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsSignerFeeDivisorUpdateStartedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsSignerFeeDivisorUpdateStartedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemSignerFeeDivisorUpdateStartedFunc func(
	SignerFeeDivisor uint16,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (sfduss *TbtcsSignerFeeDivisorUpdateStartedSubscription) OnEvent(
	handler tBTCSystemSignerFeeDivisorUpdateStartedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemSignerFeeDivisorUpdateStarted)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.SignerFeeDivisor,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sfduss.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sfduss *TbtcsSignerFeeDivisorUpdateStartedSubscription) Pipe(
	sink chan *abi.TBTCSystemSignerFeeDivisorUpdateStarted,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sfduss.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sfduss.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sfduss.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past SignerFeeDivisorUpdateStarted events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sfduss.contract.PastSignerFeeDivisorUpdateStartedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past SignerFeeDivisorUpdateStarted events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sfduss.contract.watchSignerFeeDivisorUpdateStarted(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchSignerFeeDivisorUpdateStarted(
	sink chan *abi.TBTCSystemSignerFeeDivisorUpdateStarted,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchSignerFeeDivisorUpdateStarted(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event SignerFeeDivisorUpdateStarted had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event SignerFeeDivisorUpdateStarted failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) SignerFeeDivisorUpdated(
	opts *ethlike.SubscribeOpts,
) *TbtcsSignerFeeDivisorUpdatedSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsSignerFeeDivisorUpdatedSubscription{
		tbtcs,
		opts,
	}
}

type TbtcsSignerFeeDivisorUpdatedSubscription struct {
	contract *TBTCSystem
	opts     *ethlike.SubscribeOpts
}

type tBTCSystemSignerFeeDivisorUpdatedFunc func(
	SignerFeeDivisor uint16,
	blockNumber uint64,
)

func (sfdus *TbtcsSignerFeeDivisorUpdatedSubscription) OnEvent(
	handler tBTCSystemSignerFeeDivisorUpdatedFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemSignerFeeDivisorUpdated)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.SignerFeeDivisor,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sfdus.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sfdus *TbtcsSignerFeeDivisorUpdatedSubscription) Pipe(
	sink chan *abi.TBTCSystemSignerFeeDivisorUpdated,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sfdus.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sfdus.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sfdus.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past SignerFeeDivisorUpdated events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sfdus.contract.PastSignerFeeDivisorUpdatedEvents(
					fromBlock,
					nil,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past SignerFeeDivisorUpdated events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sfdus.contract.watchSignerFeeDivisorUpdated(
		sink,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchSignerFeeDivisorUpdated(
	sink chan *abi.TBTCSystemSignerFeeDivisorUpdated,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchSignerFeeDivisorUpdated(
			&bind.WatchOpts{Context: ctx},
			sink,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event SignerFeeDivisorUpdated had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event SignerFeeDivisorUpdated failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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

func (tbtcs *TBTCSystem) StartedLiquidation(
	opts *ethlike.SubscribeOpts,
	_depositContractAddressFilter []common.Address,
) *TbtcsStartedLiquidationSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &TbtcsStartedLiquidationSubscription{
		tbtcs,
		opts,
		_depositContractAddressFilter,
	}
}

type TbtcsStartedLiquidationSubscription struct {
	contract                      *TBTCSystem
	opts                          *ethlike.SubscribeOpts
	_depositContractAddressFilter []common.Address
}

type tBTCSystemStartedLiquidationFunc func(
	DepositContractAddress common.Address,
	WasFraud bool,
	Timestamp *big.Int,
	blockNumber uint64,
)

func (sls *TbtcsStartedLiquidationSubscription) OnEvent(
	handler tBTCSystemStartedLiquidationFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.TBTCSystemStartedLiquidation)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.DepositContractAddress,
					event.WasFraud,
					event.Timestamp,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := sls.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (sls *TbtcsStartedLiquidationSubscription) Pipe(
	sink chan *abi.TBTCSystemStartedLiquidation,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(sls.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := sls.contract.blockCounter.CurrentBlock()
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - sls.opts.PastBlocks

				tbtcsLogger.Infof(
					"subscription monitoring fetching past StartedLiquidation events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := sls.contract.PastStartedLiquidationEvents(
					fromBlock,
					nil,
					sls._depositContractAddressFilter,
				)
				if err != nil {
					tbtcsLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				tbtcsLogger.Infof(
					"subscription monitoring fetched [%v] past StartedLiquidation events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := sls.contract.watchStartedLiquidation(
		sink,
		sls._depositContractAddressFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (tbtcs *TBTCSystem) watchStartedLiquidation(
	sink chan *abi.TBTCSystemStartedLiquidation,
	_depositContractAddressFilter []common.Address,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return tbtcs.contract.WatchStartedLiquidation(
			&bind.WatchOpts{Context: ctx},
			sink,
			_depositContractAddressFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		tbtcsLogger.Errorf(
			"subscription to event StartedLiquidation had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		tbtcsLogger.Errorf(
			"subscription to event StartedLiquidation failed "+
				"with error: [%v]; resubscription attempt will be "+
				"performed",
			err,
		)
	}

	return chainutil.WithResubscription(
		chainutil.SubscriptionBackoffMax,
		subscribeFn,
		chainutil.SubscriptionAlertThreshold,
		thresholdViolatedFn,
		subscriptionFailedFn,
	)
}

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
