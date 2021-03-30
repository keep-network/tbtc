// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	hostchainabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ipfs/go-log"

	chainutil "github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/chain/ethlike"
	abi "github.com/keep-network/tbtc/pkg/chain/ethereum/gen/abi/deposit"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var dLogger = log.Logger("keep-contract-Deposit")

type Deposit struct {
	contract          *abi.Deposit
	contractAddress   common.Address
	contractABI       *hostchainabi.ABI
	caller            bind.ContractCaller
	transactor        bind.ContractTransactor
	callerOptions     *bind.CallOpts
	transactorOptions *bind.TransactOpts
	errorResolver     *chainutil.ErrorResolver
	nonceManager      *ethlike.NonceManager
	miningWaiter      *ethlike.MiningWaiter
	blockCounter      *ethlike.BlockCounter

	transactionMutex *sync.Mutex
}

func NewDeposit(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *ethlike.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*Deposit, error) {
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

	contract, err := abi.NewDeposit(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.DepositABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &Deposit{
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
func (d *Deposit) ExitCourtesyCall(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction exitCourtesyCall",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ExitCourtesyCall(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"exitCourtesyCall",
		)
	}

	dLogger.Infof(
		"submitted transaction exitCourtesyCall with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ExitCourtesyCall(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"exitCourtesyCall",
				)
			}

			dLogger.Infof(
				"submitted transaction exitCourtesyCall with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallExitCourtesyCall(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"exitCourtesyCall",
		&result,
	)

	return err
}

func (d *Deposit) ExitCourtesyCallGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"exitCourtesyCall",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) IncreaseRedemptionFee(
	_previousOutputValueBytes [8]uint8,
	_newOutputValueBytes [8]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction increaseRedemptionFee",
		"params: ",
		fmt.Sprint(
			_previousOutputValueBytes,
			_newOutputValueBytes,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.IncreaseRedemptionFee(
		transactorOptions,
		_previousOutputValueBytes,
		_newOutputValueBytes,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"increaseRedemptionFee",
			_previousOutputValueBytes,
			_newOutputValueBytes,
		)
	}

	dLogger.Infof(
		"submitted transaction increaseRedemptionFee with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.IncreaseRedemptionFee(
				transactorOptions,
				_previousOutputValueBytes,
				_newOutputValueBytes,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"increaseRedemptionFee",
					_previousOutputValueBytes,
					_newOutputValueBytes,
				)
			}

			dLogger.Infof(
				"submitted transaction increaseRedemptionFee with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallIncreaseRedemptionFee(
	_previousOutputValueBytes [8]uint8,
	_newOutputValueBytes [8]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"increaseRedemptionFee",
		&result,
		_previousOutputValueBytes,
		_newOutputValueBytes,
	)

	return err
}

func (d *Deposit) IncreaseRedemptionFeeGasEstimate(
	_previousOutputValueBytes [8]uint8,
	_newOutputValueBytes [8]uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"increaseRedemptionFee",
		d.contractABI,
		d.transactor,
		_previousOutputValueBytes,
		_newOutputValueBytes,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) Initialize(
	_factory common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction initialize",
		"params: ",
		fmt.Sprint(
			_factory,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.Initialize(
		transactorOptions,
		_factory,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"initialize",
			_factory,
		)
	}

	dLogger.Infof(
		"submitted transaction initialize with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.Initialize(
				transactorOptions,
				_factory,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"initialize",
					_factory,
				)
			}

			dLogger.Infof(
				"submitted transaction initialize with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallInitialize(
	_factory common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"initialize",
		&result,
		_factory,
	)

	return err
}

func (d *Deposit) InitializeGasEstimate(
	_factory common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"initialize",
		d.contractABI,
		d.transactor,
		_factory,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) InitializeDeposit(
	_tbtcSystem common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachineAddress common.Address,
	_lotSizeSatoshis uint64,
	value *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction initializeDeposit",
		"params: ",
		fmt.Sprint(
			_tbtcSystem,
			_tbtcToken,
			_tbtcDepositToken,
			_feeRebateToken,
			_vendingMachineAddress,
			_lotSizeSatoshis,
		),
		"value: ", value,
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	transactorOptions.Value = value

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.InitializeDeposit(
		transactorOptions,
		_tbtcSystem,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachineAddress,
		_lotSizeSatoshis,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			value,
			"initializeDeposit",
			_tbtcSystem,
			_tbtcToken,
			_tbtcDepositToken,
			_feeRebateToken,
			_vendingMachineAddress,
			_lotSizeSatoshis,
		)
	}

	dLogger.Infof(
		"submitted transaction initializeDeposit with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.InitializeDeposit(
				transactorOptions,
				_tbtcSystem,
				_tbtcToken,
				_tbtcDepositToken,
				_feeRebateToken,
				_vendingMachineAddress,
				_lotSizeSatoshis,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					value,
					"initializeDeposit",
					_tbtcSystem,
					_tbtcToken,
					_tbtcDepositToken,
					_feeRebateToken,
					_vendingMachineAddress,
					_lotSizeSatoshis,
				)
			}

			dLogger.Infof(
				"submitted transaction initializeDeposit with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallInitializeDeposit(
	_tbtcSystem common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachineAddress common.Address,
	_lotSizeSatoshis uint64,
	value *big.Int,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, value,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"initializeDeposit",
		&result,
		_tbtcSystem,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachineAddress,
		_lotSizeSatoshis,
	)

	return err
}

func (d *Deposit) InitializeDepositGasEstimate(
	_tbtcSystem common.Address,
	_tbtcToken common.Address,
	_tbtcDepositToken common.Address,
	_feeRebateToken common.Address,
	_vendingMachineAddress common.Address,
	_lotSizeSatoshis uint64,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"initializeDeposit",
		d.contractABI,
		d.transactor,
		_tbtcSystem,
		_tbtcToken,
		_tbtcDepositToken,
		_feeRebateToken,
		_vendingMachineAddress,
		_lotSizeSatoshis,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyCourtesyCall(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyCourtesyCall",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyCourtesyCall(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyCourtesyCall",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyCourtesyCall with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyCourtesyCall(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyCourtesyCall",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyCourtesyCall with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyCourtesyCall(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyCourtesyCall",
		&result,
	)

	return err
}

func (d *Deposit) NotifyCourtesyCallGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyCourtesyCall",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyCourtesyCallExpired(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyCourtesyCallExpired",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyCourtesyCallExpired(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyCourtesyCallExpired",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyCourtesyCallExpired with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyCourtesyCallExpired(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyCourtesyCallExpired",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyCourtesyCallExpired with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyCourtesyCallExpired(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyCourtesyCallExpired",
		&result,
	)

	return err
}

func (d *Deposit) NotifyCourtesyCallExpiredGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyCourtesyCallExpired",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyFundingTimedOut(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyFundingTimedOut",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyFundingTimedOut(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyFundingTimedOut",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyFundingTimedOut with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyFundingTimedOut(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyFundingTimedOut",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyFundingTimedOut with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyFundingTimedOut(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyFundingTimedOut",
		&result,
	)

	return err
}

func (d *Deposit) NotifyFundingTimedOutGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyFundingTimedOut",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyRedemptionProofTimedOut(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyRedemptionProofTimedOut",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyRedemptionProofTimedOut(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyRedemptionProofTimedOut",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyRedemptionProofTimedOut with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyRedemptionProofTimedOut(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyRedemptionProofTimedOut",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyRedemptionProofTimedOut with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyRedemptionProofTimedOut(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyRedemptionProofTimedOut",
		&result,
	)

	return err
}

func (d *Deposit) NotifyRedemptionProofTimedOutGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyRedemptionProofTimedOut",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyRedemptionSignatureTimedOut(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyRedemptionSignatureTimedOut",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyRedemptionSignatureTimedOut(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyRedemptionSignatureTimedOut",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyRedemptionSignatureTimedOut with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyRedemptionSignatureTimedOut(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyRedemptionSignatureTimedOut",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyRedemptionSignatureTimedOut with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyRedemptionSignatureTimedOut(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyRedemptionSignatureTimedOut",
		&result,
	)

	return err
}

func (d *Deposit) NotifyRedemptionSignatureTimedOutGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyRedemptionSignatureTimedOut",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifySignerSetupFailed(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifySignerSetupFailed",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifySignerSetupFailed(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifySignerSetupFailed",
		)
	}

	dLogger.Infof(
		"submitted transaction notifySignerSetupFailed with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifySignerSetupFailed(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifySignerSetupFailed",
				)
			}

			dLogger.Infof(
				"submitted transaction notifySignerSetupFailed with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifySignerSetupFailed(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifySignerSetupFailed",
		&result,
	)

	return err
}

func (d *Deposit) NotifySignerSetupFailedGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifySignerSetupFailed",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) NotifyUndercollateralizedLiquidation(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction notifyUndercollateralizedLiquidation",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.NotifyUndercollateralizedLiquidation(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"notifyUndercollateralizedLiquidation",
		)
	}

	dLogger.Infof(
		"submitted transaction notifyUndercollateralizedLiquidation with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.NotifyUndercollateralizedLiquidation(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"notifyUndercollateralizedLiquidation",
				)
			}

			dLogger.Infof(
				"submitted transaction notifyUndercollateralizedLiquidation with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallNotifyUndercollateralizedLiquidation(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"notifyUndercollateralizedLiquidation",
		&result,
	)

	return err
}

func (d *Deposit) NotifyUndercollateralizedLiquidationGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"notifyUndercollateralizedLiquidation",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) ProvideBTCFundingProof(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_fundingOutputIndex uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction provideBTCFundingProof",
		"params: ",
		fmt.Sprint(
			_txVersion,
			_txInputVector,
			_txOutputVector,
			_txLocktime,
			_fundingOutputIndex,
			_merkleProof,
			_txIndexInBlock,
			_bitcoinHeaders,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ProvideBTCFundingProof(
		transactorOptions,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_fundingOutputIndex,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"provideBTCFundingProof",
			_txVersion,
			_txInputVector,
			_txOutputVector,
			_txLocktime,
			_fundingOutputIndex,
			_merkleProof,
			_txIndexInBlock,
			_bitcoinHeaders,
		)
	}

	dLogger.Infof(
		"submitted transaction provideBTCFundingProof with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ProvideBTCFundingProof(
				transactorOptions,
				_txVersion,
				_txInputVector,
				_txOutputVector,
				_txLocktime,
				_fundingOutputIndex,
				_merkleProof,
				_txIndexInBlock,
				_bitcoinHeaders,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"provideBTCFundingProof",
					_txVersion,
					_txInputVector,
					_txOutputVector,
					_txLocktime,
					_fundingOutputIndex,
					_merkleProof,
					_txIndexInBlock,
					_bitcoinHeaders,
				)
			}

			dLogger.Infof(
				"submitted transaction provideBTCFundingProof with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallProvideBTCFundingProof(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_fundingOutputIndex uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"provideBTCFundingProof",
		&result,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_fundingOutputIndex,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)

	return err
}

func (d *Deposit) ProvideBTCFundingProofGasEstimate(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_fundingOutputIndex uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"provideBTCFundingProof",
		d.contractABI,
		d.transactor,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_fundingOutputIndex,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) ProvideECDSAFraudProof(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction provideECDSAFraudProof",
		"params: ",
		fmt.Sprint(
			_v,
			_r,
			_s,
			_signedDigest,
			_preimage,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ProvideECDSAFraudProof(
		transactorOptions,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"provideECDSAFraudProof",
			_v,
			_r,
			_s,
			_signedDigest,
			_preimage,
		)
	}

	dLogger.Infof(
		"submitted transaction provideECDSAFraudProof with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ProvideECDSAFraudProof(
				transactorOptions,
				_v,
				_r,
				_s,
				_signedDigest,
				_preimage,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"provideECDSAFraudProof",
					_v,
					_r,
					_s,
					_signedDigest,
					_preimage,
				)
			}

			dLogger.Infof(
				"submitted transaction provideECDSAFraudProof with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallProvideECDSAFraudProof(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"provideECDSAFraudProof",
		&result,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)

	return err
}

func (d *Deposit) ProvideECDSAFraudProofGasEstimate(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"provideECDSAFraudProof",
		d.contractABI,
		d.transactor,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) ProvideFundingECDSAFraudProof(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction provideFundingECDSAFraudProof",
		"params: ",
		fmt.Sprint(
			_v,
			_r,
			_s,
			_signedDigest,
			_preimage,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ProvideFundingECDSAFraudProof(
		transactorOptions,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"provideFundingECDSAFraudProof",
			_v,
			_r,
			_s,
			_signedDigest,
			_preimage,
		)
	}

	dLogger.Infof(
		"submitted transaction provideFundingECDSAFraudProof with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ProvideFundingECDSAFraudProof(
				transactorOptions,
				_v,
				_r,
				_s,
				_signedDigest,
				_preimage,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"provideFundingECDSAFraudProof",
					_v,
					_r,
					_s,
					_signedDigest,
					_preimage,
				)
			}

			dLogger.Infof(
				"submitted transaction provideFundingECDSAFraudProof with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallProvideFundingECDSAFraudProof(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"provideFundingECDSAFraudProof",
		&result,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)

	return err
}

func (d *Deposit) ProvideFundingECDSAFraudProofGasEstimate(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	_signedDigest [32]uint8,
	_preimage []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"provideFundingECDSAFraudProof",
		d.contractABI,
		d.transactor,
		_v,
		_r,
		_s,
		_signedDigest,
		_preimage,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) ProvideRedemptionProof(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction provideRedemptionProof",
		"params: ",
		fmt.Sprint(
			_txVersion,
			_txInputVector,
			_txOutputVector,
			_txLocktime,
			_merkleProof,
			_txIndexInBlock,
			_bitcoinHeaders,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ProvideRedemptionProof(
		transactorOptions,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"provideRedemptionProof",
			_txVersion,
			_txInputVector,
			_txOutputVector,
			_txLocktime,
			_merkleProof,
			_txIndexInBlock,
			_bitcoinHeaders,
		)
	}

	dLogger.Infof(
		"submitted transaction provideRedemptionProof with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ProvideRedemptionProof(
				transactorOptions,
				_txVersion,
				_txInputVector,
				_txOutputVector,
				_txLocktime,
				_merkleProof,
				_txIndexInBlock,
				_bitcoinHeaders,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"provideRedemptionProof",
					_txVersion,
					_txInputVector,
					_txOutputVector,
					_txLocktime,
					_merkleProof,
					_txIndexInBlock,
					_bitcoinHeaders,
				)
			}

			dLogger.Infof(
				"submitted transaction provideRedemptionProof with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallProvideRedemptionProof(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"provideRedemptionProof",
		&result,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)

	return err
}

func (d *Deposit) ProvideRedemptionProofGasEstimate(
	_txVersion [4]uint8,
	_txInputVector []uint8,
	_txOutputVector []uint8,
	_txLocktime [4]uint8,
	_merkleProof []uint8,
	_txIndexInBlock *big.Int,
	_bitcoinHeaders []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"provideRedemptionProof",
		d.contractABI,
		d.transactor,
		_txVersion,
		_txInputVector,
		_txOutputVector,
		_txLocktime,
		_merkleProof,
		_txIndexInBlock,
		_bitcoinHeaders,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) ProvideRedemptionSignature(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction provideRedemptionSignature",
		"params: ",
		fmt.Sprint(
			_v,
			_r,
			_s,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.ProvideRedemptionSignature(
		transactorOptions,
		_v,
		_r,
		_s,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"provideRedemptionSignature",
			_v,
			_r,
			_s,
		)
	}

	dLogger.Infof(
		"submitted transaction provideRedemptionSignature with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.ProvideRedemptionSignature(
				transactorOptions,
				_v,
				_r,
				_s,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"provideRedemptionSignature",
					_v,
					_r,
					_s,
				)
			}

			dLogger.Infof(
				"submitted transaction provideRedemptionSignature with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallProvideRedemptionSignature(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"provideRedemptionSignature",
		&result,
		_v,
		_r,
		_s,
	)

	return err
}

func (d *Deposit) ProvideRedemptionSignatureGasEstimate(
	_v uint8,
	_r [32]uint8,
	_s [32]uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"provideRedemptionSignature",
		d.contractABI,
		d.transactor,
		_v,
		_r,
		_s,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) PurchaseSignerBondsAtAuction(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction purchaseSignerBondsAtAuction",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.PurchaseSignerBondsAtAuction(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"purchaseSignerBondsAtAuction",
		)
	}

	dLogger.Infof(
		"submitted transaction purchaseSignerBondsAtAuction with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.PurchaseSignerBondsAtAuction(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"purchaseSignerBondsAtAuction",
				)
			}

			dLogger.Infof(
				"submitted transaction purchaseSignerBondsAtAuction with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallPurchaseSignerBondsAtAuction(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"purchaseSignerBondsAtAuction",
		&result,
	)

	return err
}

func (d *Deposit) PurchaseSignerBondsAtAuctionGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"purchaseSignerBondsAtAuction",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) RequestFunderAbort(
	_abortOutputScript []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction requestFunderAbort",
		"params: ",
		fmt.Sprint(
			_abortOutputScript,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.RequestFunderAbort(
		transactorOptions,
		_abortOutputScript,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"requestFunderAbort",
			_abortOutputScript,
		)
	}

	dLogger.Infof(
		"submitted transaction requestFunderAbort with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.RequestFunderAbort(
				transactorOptions,
				_abortOutputScript,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"requestFunderAbort",
					_abortOutputScript,
				)
			}

			dLogger.Infof(
				"submitted transaction requestFunderAbort with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallRequestFunderAbort(
	_abortOutputScript []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"requestFunderAbort",
		&result,
		_abortOutputScript,
	)

	return err
}

func (d *Deposit) RequestFunderAbortGasEstimate(
	_abortOutputScript []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"requestFunderAbort",
		d.contractABI,
		d.transactor,
		_abortOutputScript,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) RequestRedemption(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction requestRedemption",
		"params: ",
		fmt.Sprint(
			_outputValueBytes,
			_redeemerOutputScript,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.RequestRedemption(
		transactorOptions,
		_outputValueBytes,
		_redeemerOutputScript,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"requestRedemption",
			_outputValueBytes,
			_redeemerOutputScript,
		)
	}

	dLogger.Infof(
		"submitted transaction requestRedemption with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.RequestRedemption(
				transactorOptions,
				_outputValueBytes,
				_redeemerOutputScript,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"requestRedemption",
					_outputValueBytes,
					_redeemerOutputScript,
				)
			}

			dLogger.Infof(
				"submitted transaction requestRedemption with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallRequestRedemption(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"requestRedemption",
		&result,
		_outputValueBytes,
		_redeemerOutputScript,
	)

	return err
}

func (d *Deposit) RequestRedemptionGasEstimate(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"requestRedemption",
		d.contractABI,
		d.transactor,
		_outputValueBytes,
		_redeemerOutputScript,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) RetrieveSignerPubkey(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction retrieveSignerPubkey",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.RetrieveSignerPubkey(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"retrieveSignerPubkey",
		)
	}

	dLogger.Infof(
		"submitted transaction retrieveSignerPubkey with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.RetrieveSignerPubkey(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"retrieveSignerPubkey",
				)
			}

			dLogger.Infof(
				"submitted transaction retrieveSignerPubkey with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallRetrieveSignerPubkey(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"retrieveSignerPubkey",
		&result,
	)

	return err
}

func (d *Deposit) RetrieveSignerPubkeyGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"retrieveSignerPubkey",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) TransferAndRequestRedemption(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,
	_finalRecipient common.Address,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction transferAndRequestRedemption",
		"params: ",
		fmt.Sprint(
			_outputValueBytes,
			_redeemerOutputScript,
			_finalRecipient,
		),
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.TransferAndRequestRedemption(
		transactorOptions,
		_outputValueBytes,
		_redeemerOutputScript,
		_finalRecipient,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"transferAndRequestRedemption",
			_outputValueBytes,
			_redeemerOutputScript,
			_finalRecipient,
		)
	}

	dLogger.Infof(
		"submitted transaction transferAndRequestRedemption with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.TransferAndRequestRedemption(
				transactorOptions,
				_outputValueBytes,
				_redeemerOutputScript,
				_finalRecipient,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"transferAndRequestRedemption",
					_outputValueBytes,
					_redeemerOutputScript,
					_finalRecipient,
				)
			}

			dLogger.Infof(
				"submitted transaction transferAndRequestRedemption with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallTransferAndRequestRedemption(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,
	_finalRecipient common.Address,
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"transferAndRequestRedemption",
		&result,
		_outputValueBytes,
		_redeemerOutputScript,
		_finalRecipient,
	)

	return err
}

func (d *Deposit) TransferAndRequestRedemptionGasEstimate(
	_outputValueBytes [8]uint8,
	_redeemerOutputScript []uint8,
	_finalRecipient common.Address,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"transferAndRequestRedemption",
		d.contractABI,
		d.transactor,
		_outputValueBytes,
		_redeemerOutputScript,
		_finalRecipient,
	)

	return result, err
}

// Transaction submission.
func (d *Deposit) WithdrawFunds(

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	dLogger.Debug(
		"submitting transaction withdrawFunds",
	)

	d.transactionMutex.Lock()
	defer d.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *d.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := d.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := d.contract.WithdrawFunds(
		transactorOptions,
	)
	if err != nil {
		return transaction, d.errorResolver.ResolveError(
			err,
			d.transactorOptions.From,
			nil,
			"withdrawFunds",
		)
	}

	dLogger.Infof(
		"submitted transaction withdrawFunds with id: [%v] and nonce [%v]",
		transaction.Hash().Hex(),
		transaction.Nonce(),
	)

	go d.miningWaiter.ForceMining(
		&ethlike.Transaction{
			Hash:     ethlike.Hash(transaction.Hash()),
			GasPrice: transaction.GasPrice(),
		},
		func(newGasPrice *big.Int) (*ethlike.Transaction, error) {
			transactorOptions.GasLimit = transaction.Gas()
			transactorOptions.GasPrice = newGasPrice

			transaction, err := d.contract.WithdrawFunds(
				transactorOptions,
			)
			if err != nil {
				return nil, d.errorResolver.ResolveError(
					err,
					d.transactorOptions.From,
					nil,
					"withdrawFunds",
				)
			}

			dLogger.Infof(
				"submitted transaction withdrawFunds with id: [%v] and nonce [%v]",
				transaction.Hash().Hex(),
				transaction.Nonce(),
			)

			return &ethlike.Transaction{
				Hash:     ethlike.Hash(transaction.Hash()),
				GasPrice: transaction.GasPrice(),
			}, nil
		},
	)

	d.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (d *Deposit) CallWithdrawFunds(
	blockNumber *big.Int,
) error {
	var result interface{} = nil

	err := chainutil.CallAtBlock(
		d.transactorOptions.From,
		blockNumber, nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"withdrawFunds",
		&result,
	)

	return err
}

func (d *Deposit) WithdrawFundsGasEstimate() (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		d.callerOptions.From,
		d.contractAddress,
		"withdrawFunds",
		d.contractABI,
		d.transactor,
	)

	return result, err
}

// ----- Const Methods ------

func (d *Deposit) AuctionValue() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.AuctionValue(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"auctionValue",
		)
	}

	return result, err
}

func (d *Deposit) AuctionValueAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"auctionValue",
		&result,
	)

	return result, err
}

func (d *Deposit) CollateralizationPercentage() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.CollateralizationPercentage(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"collateralizationPercentage",
		)
	}

	return result, err
}

func (d *Deposit) CollateralizationPercentageAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"collateralizationPercentage",
		&result,
	)

	return result, err
}

func (d *Deposit) CurrentState() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.CurrentState(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"currentState",
		)
	}

	return result, err
}

func (d *Deposit) CurrentStateAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"currentState",
		&result,
	)

	return result, err
}

type fundingInfo struct {
	UtxoValueBytes [8]uint8
	FundedAt       *big.Int
	UtxoOutpoint   []uint8
}

func (d *Deposit) FundingInfo() (fundingInfo, error) {
	var result fundingInfo
	result, err := d.contract.FundingInfo(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"fundingInfo",
		)
	}

	return result, err
}

func (d *Deposit) FundingInfoAtBlock(
	blockNumber *big.Int,
) (fundingInfo, error) {
	var result fundingInfo

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"fundingInfo",
		&result,
	)

	return result, err
}

func (d *Deposit) GetOwnerRedemptionTbtcRequirement(
	_redeemer common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.GetOwnerRedemptionTbtcRequirement(
		d.callerOptions,
		_redeemer,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"getOwnerRedemptionTbtcRequirement",
			_redeemer,
		)
	}

	return result, err
}

func (d *Deposit) GetOwnerRedemptionTbtcRequirementAtBlock(
	_redeemer common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"getOwnerRedemptionTbtcRequirement",
		&result,
		_redeemer,
	)

	return result, err
}

func (d *Deposit) GetRedemptionTbtcRequirement(
	_redeemer common.Address,
) (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.GetRedemptionTbtcRequirement(
		d.callerOptions,
		_redeemer,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"getRedemptionTbtcRequirement",
			_redeemer,
		)
	}

	return result, err
}

func (d *Deposit) GetRedemptionTbtcRequirementAtBlock(
	_redeemer common.Address,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"getRedemptionTbtcRequirement",
		&result,
		_redeemer,
	)

	return result, err
}

func (d *Deposit) InActive() (bool, error) {
	var result bool
	result, err := d.contract.InActive(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"inActive",
		)
	}

	return result, err
}

func (d *Deposit) InActiveAtBlock(
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"inActive",
		&result,
	)

	return result, err
}

func (d *Deposit) InitialCollateralizedPercent() (uint16, error) {
	var result uint16
	result, err := d.contract.InitialCollateralizedPercent(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"initialCollateralizedPercent",
		)
	}

	return result, err
}

func (d *Deposit) InitialCollateralizedPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"initialCollateralizedPercent",
		&result,
	)

	return result, err
}

func (d *Deposit) KeepAddress() (common.Address, error) {
	var result common.Address
	result, err := d.contract.KeepAddress(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"keepAddress",
		)
	}

	return result, err
}

func (d *Deposit) KeepAddressAtBlock(
	blockNumber *big.Int,
) (common.Address, error) {
	var result common.Address

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"keepAddress",
		&result,
	)

	return result, err
}

func (d *Deposit) LotSizeSatoshis() (uint64, error) {
	var result uint64
	result, err := d.contract.LotSizeSatoshis(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"lotSizeSatoshis",
		)
	}

	return result, err
}

func (d *Deposit) LotSizeSatoshisAtBlock(
	blockNumber *big.Int,
) (uint64, error) {
	var result uint64

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"lotSizeSatoshis",
		&result,
	)

	return result, err
}

func (d *Deposit) LotSizeTbtc() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.LotSizeTbtc(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"lotSizeTbtc",
		)
	}

	return result, err
}

func (d *Deposit) LotSizeTbtcAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"lotSizeTbtc",
		&result,
	)

	return result, err
}

func (d *Deposit) RemainingTerm() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.RemainingTerm(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"remainingTerm",
		)
	}

	return result, err
}

func (d *Deposit) RemainingTermAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"remainingTerm",
		&result,
	)

	return result, err
}

func (d *Deposit) SeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	var result uint16
	result, err := d.contract.SeverelyUndercollateralizedThresholdPercent(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"severelyUndercollateralizedThresholdPercent",
		)
	}

	return result, err
}

func (d *Deposit) SeverelyUndercollateralizedThresholdPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"severelyUndercollateralizedThresholdPercent",
		&result,
	)

	return result, err
}

func (d *Deposit) SignerFeeTbtc() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.SignerFeeTbtc(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"signerFeeTbtc",
		)
	}

	return result, err
}

func (d *Deposit) SignerFeeTbtcAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"signerFeeTbtc",
		&result,
	)

	return result, err
}

func (d *Deposit) UndercollateralizedThresholdPercent() (uint16, error) {
	var result uint16
	result, err := d.contract.UndercollateralizedThresholdPercent(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"undercollateralizedThresholdPercent",
		)
	}

	return result, err
}

func (d *Deposit) UndercollateralizedThresholdPercentAtBlock(
	blockNumber *big.Int,
) (uint16, error) {
	var result uint16

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"undercollateralizedThresholdPercent",
		&result,
	)

	return result, err
}

func (d *Deposit) UtxoValue() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.UtxoValue(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"utxoValue",
		)
	}

	return result, err
}

func (d *Deposit) UtxoValueAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"utxoValue",
		&result,
	)

	return result, err
}

func (d *Deposit) WithdrawableAmount() (*big.Int, error) {
	var result *big.Int
	result, err := d.contract.WithdrawableAmount(
		d.callerOptions,
	)

	if err != nil {
		return result, d.errorResolver.ResolveError(
			err,
			d.callerOptions.From,
			nil,
			"withdrawableAmount",
		)
	}

	return result, err
}

func (d *Deposit) WithdrawableAmountAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		d.callerOptions.From,
		blockNumber,
		nil,
		d.contractABI,
		d.caller,
		d.errorResolver,
		d.contractAddress,
		"withdrawableAmount",
		&result,
	)

	return result, err
}

// ------ Events -------
