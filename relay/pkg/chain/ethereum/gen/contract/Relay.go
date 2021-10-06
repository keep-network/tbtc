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
	"github.com/keep-network/tbtc/relay/pkg/chain/ethereum/gen/abi"
)

// Create a package-level logger for this contract. The logger exists at
// package level so that the logger is registered at startup and can be
// included or excluded from logging at startup by name.
var rLogger = log.Logger("keep-contract-Relay")

type Relay struct {
	contract          *abi.Relay
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

func NewRelay(
	contractAddress common.Address,
	chainId *big.Int,
	accountKey *keystore.Key,
	backend bind.ContractBackend,
	nonceManager *ethlike.NonceManager,
	miningWaiter *chainutil.MiningWaiter,
	blockCounter *ethlike.BlockCounter,
	transactionMutex *sync.Mutex,
) (*Relay, error) {
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

	contract, err := abi.NewRelay(
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

	contractABI, err := hostchainabi.JSON(strings.NewReader(abi.RelayABI))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate ABI: [%v]", err)
	}

	return &Relay{
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
func (r *Relay) AddHeaders(
	_anchor []uint8,
	_headers []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	rLogger.Debug(
		"submitting transaction addHeaders",
		" params: ",
		fmt.Sprint(
			_anchor,
			_headers,
		),
	)

	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *r.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := r.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := r.contract.AddHeaders(
		transactorOptions,
		_anchor,
		_headers,
	)
	if err != nil {
		return transaction, r.errorResolver.ResolveError(
			err,
			r.transactorOptions.From,
			nil,
			"addHeaders",
			_anchor,
			_headers,
		)
	}

	rLogger.Infof(
		"submitted transaction addHeaders with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go r.miningWaiter.ForceMining(
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

			transaction, err := r.contract.AddHeaders(
				newTransactorOptions,
				_anchor,
				_headers,
			)
			if err != nil {
				return nil, r.errorResolver.ResolveError(
					err,
					r.transactorOptions.From,
					nil,
					"addHeaders",
					_anchor,
					_headers,
				)
			}

			rLogger.Infof(
				"submitted transaction addHeaders with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	r.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (r *Relay) CallAddHeaders(
	_anchor []uint8,
	_headers []uint8,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		r.transactorOptions.From,
		blockNumber, nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"addHeaders",
		&result,
		_anchor,
		_headers,
	)

	return result, err
}

func (r *Relay) AddHeadersGasEstimate(
	_anchor []uint8,
	_headers []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		r.callerOptions.From,
		r.contractAddress,
		"addHeaders",
		r.contractABI,
		r.transactor,
		_anchor,
		_headers,
	)

	return result, err
}

// Transaction submission.
func (r *Relay) AddHeadersWithRetarget(
	_oldPeriodStartHeader []uint8,
	_oldPeriodEndHeader []uint8,
	_headers []uint8,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	rLogger.Debug(
		"submitting transaction addHeadersWithRetarget",
		" params: ",
		fmt.Sprint(
			_oldPeriodStartHeader,
			_oldPeriodEndHeader,
			_headers,
		),
	)

	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *r.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := r.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := r.contract.AddHeadersWithRetarget(
		transactorOptions,
		_oldPeriodStartHeader,
		_oldPeriodEndHeader,
		_headers,
	)
	if err != nil {
		return transaction, r.errorResolver.ResolveError(
			err,
			r.transactorOptions.From,
			nil,
			"addHeadersWithRetarget",
			_oldPeriodStartHeader,
			_oldPeriodEndHeader,
			_headers,
		)
	}

	rLogger.Infof(
		"submitted transaction addHeadersWithRetarget with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go r.miningWaiter.ForceMining(
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

			transaction, err := r.contract.AddHeadersWithRetarget(
				newTransactorOptions,
				_oldPeriodStartHeader,
				_oldPeriodEndHeader,
				_headers,
			)
			if err != nil {
				return nil, r.errorResolver.ResolveError(
					err,
					r.transactorOptions.From,
					nil,
					"addHeadersWithRetarget",
					_oldPeriodStartHeader,
					_oldPeriodEndHeader,
					_headers,
				)
			}

			rLogger.Infof(
				"submitted transaction addHeadersWithRetarget with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	r.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (r *Relay) CallAddHeadersWithRetarget(
	_oldPeriodStartHeader []uint8,
	_oldPeriodEndHeader []uint8,
	_headers []uint8,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		r.transactorOptions.From,
		blockNumber, nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"addHeadersWithRetarget",
		&result,
		_oldPeriodStartHeader,
		_oldPeriodEndHeader,
		_headers,
	)

	return result, err
}

func (r *Relay) AddHeadersWithRetargetGasEstimate(
	_oldPeriodStartHeader []uint8,
	_oldPeriodEndHeader []uint8,
	_headers []uint8,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		r.callerOptions.From,
		r.contractAddress,
		"addHeadersWithRetarget",
		r.contractABI,
		r.transactor,
		_oldPeriodStartHeader,
		_oldPeriodEndHeader,
		_headers,
	)

	return result, err
}

// Transaction submission.
func (r *Relay) MarkNewHeaviest(
	_ancestor [32]uint8,
	_currentBest []uint8,
	_newBest []uint8,
	_limit *big.Int,

	transactionOptions ...chainutil.TransactionOptions,
) (*types.Transaction, error) {
	rLogger.Debug(
		"submitting transaction markNewHeaviest",
		" params: ",
		fmt.Sprint(
			_ancestor,
			_currentBest,
			_newBest,
			_limit,
		),
	)

	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()

	// create a copy
	transactorOptions := new(bind.TransactOpts)
	*transactorOptions = *r.transactorOptions

	if len(transactionOptions) > 1 {
		return nil, fmt.Errorf(
			"could not process multiple transaction options sets",
		)
	} else if len(transactionOptions) > 0 {
		transactionOptions[0].Apply(transactorOptions)
	}

	nonce, err := r.nonceManager.CurrentNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
	}

	transactorOptions.Nonce = new(big.Int).SetUint64(nonce)

	transaction, err := r.contract.MarkNewHeaviest(
		transactorOptions,
		_ancestor,
		_currentBest,
		_newBest,
		_limit,
	)
	if err != nil {
		return transaction, r.errorResolver.ResolveError(
			err,
			r.transactorOptions.From,
			nil,
			"markNewHeaviest",
			_ancestor,
			_currentBest,
			_newBest,
			_limit,
		)
	}

	rLogger.Infof(
		"submitted transaction markNewHeaviest with id: [%s] and nonce [%v]",
		transaction.Hash(),
		transaction.Nonce(),
	)

	go r.miningWaiter.ForceMining(
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

			transaction, err := r.contract.MarkNewHeaviest(
				newTransactorOptions,
				_ancestor,
				_currentBest,
				_newBest,
				_limit,
			)
			if err != nil {
				return nil, r.errorResolver.ResolveError(
					err,
					r.transactorOptions.From,
					nil,
					"markNewHeaviest",
					_ancestor,
					_currentBest,
					_newBest,
					_limit,
				)
			}

			rLogger.Infof(
				"submitted transaction markNewHeaviest with id: [%s] and nonce [%v]",
				transaction.Hash(),
				transaction.Nonce(),
			)

			return transaction, nil
		},
	)

	r.nonceManager.IncrementNonce()

	return transaction, err
}

// Non-mutating call, not a transaction submission.
func (r *Relay) CallMarkNewHeaviest(
	_ancestor [32]uint8,
	_currentBest []uint8,
	_newBest []uint8,
	_limit *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		r.transactorOptions.From,
		blockNumber, nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"markNewHeaviest",
		&result,
		_ancestor,
		_currentBest,
		_newBest,
		_limit,
	)

	return result, err
}

func (r *Relay) MarkNewHeaviestGasEstimate(
	_ancestor [32]uint8,
	_currentBest []uint8,
	_newBest []uint8,
	_limit *big.Int,
) (uint64, error) {
	var result uint64

	result, err := chainutil.EstimateGas(
		r.callerOptions.From,
		r.contractAddress,
		"markNewHeaviest",
		r.contractABI,
		r.transactor,
		_ancestor,
		_currentBest,
		_newBest,
		_limit,
	)

	return result, err
}

// ----- Const Methods ------

func (r *Relay) FindAncestor(
	_digest [32]uint8,
	_offset *big.Int,
) ([32]uint8, error) {
	var result [32]uint8
	result, err := r.contract.FindAncestor(
		r.callerOptions,
		_digest,
		_offset,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"findAncestor",
			_digest,
			_offset,
		)
	}

	return result, err
}

func (r *Relay) FindAncestorAtBlock(
	_digest [32]uint8,
	_offset *big.Int,
	blockNumber *big.Int,
) ([32]uint8, error) {
	var result [32]uint8

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"findAncestor",
		&result,
		_digest,
		_offset,
	)

	return result, err
}

func (r *Relay) FindHeight(
	_digest [32]uint8,
) (*big.Int, error) {
	var result *big.Int
	result, err := r.contract.FindHeight(
		r.callerOptions,
		_digest,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"findHeight",
			_digest,
		)
	}

	return result, err
}

func (r *Relay) FindHeightAtBlock(
	_digest [32]uint8,
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"findHeight",
		&result,
		_digest,
	)

	return result, err
}

func (r *Relay) GetBestKnownDigest() ([32]uint8, error) {
	var result [32]uint8
	result, err := r.contract.GetBestKnownDigest(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"getBestKnownDigest",
		)
	}

	return result, err
}

func (r *Relay) GetBestKnownDigestAtBlock(
	blockNumber *big.Int,
) ([32]uint8, error) {
	var result [32]uint8

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"getBestKnownDigest",
		&result,
	)

	return result, err
}

func (r *Relay) GetCurrentEpochDifficulty() (*big.Int, error) {
	var result *big.Int
	result, err := r.contract.GetCurrentEpochDifficulty(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"getCurrentEpochDifficulty",
		)
	}

	return result, err
}

func (r *Relay) GetCurrentEpochDifficultyAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"getCurrentEpochDifficulty",
		&result,
	)

	return result, err
}

func (r *Relay) GetLastReorgCommonAncestor() ([32]uint8, error) {
	var result [32]uint8
	result, err := r.contract.GetLastReorgCommonAncestor(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"getLastReorgCommonAncestor",
		)
	}

	return result, err
}

func (r *Relay) GetLastReorgCommonAncestorAtBlock(
	blockNumber *big.Int,
) ([32]uint8, error) {
	var result [32]uint8

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"getLastReorgCommonAncestor",
		&result,
	)

	return result, err
}

func (r *Relay) GetPrevEpochDifficulty() (*big.Int, error) {
	var result *big.Int
	result, err := r.contract.GetPrevEpochDifficulty(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"getPrevEpochDifficulty",
		)
	}

	return result, err
}

func (r *Relay) GetPrevEpochDifficultyAtBlock(
	blockNumber *big.Int,
) (*big.Int, error) {
	var result *big.Int

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"getPrevEpochDifficulty",
		&result,
	)

	return result, err
}

func (r *Relay) GetRelayGenesis() ([32]uint8, error) {
	var result [32]uint8
	result, err := r.contract.GetRelayGenesis(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"getRelayGenesis",
		)
	}

	return result, err
}

func (r *Relay) GetRelayGenesisAtBlock(
	blockNumber *big.Int,
) ([32]uint8, error) {
	var result [32]uint8

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"getRelayGenesis",
		&result,
	)

	return result, err
}

func (r *Relay) HEIGHTINTERVAL() (uint32, error) {
	var result uint32
	result, err := r.contract.HEIGHTINTERVAL(
		r.callerOptions,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"hEIGHTINTERVAL",
		)
	}

	return result, err
}

func (r *Relay) HEIGHTINTERVALAtBlock(
	blockNumber *big.Int,
) (uint32, error) {
	var result uint32

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"hEIGHTINTERVAL",
		&result,
	)

	return result, err
}

func (r *Relay) IsAncestor(
	_ancestor [32]uint8,
	_descendant [32]uint8,
	_limit *big.Int,
) (bool, error) {
	var result bool
	result, err := r.contract.IsAncestor(
		r.callerOptions,
		_ancestor,
		_descendant,
		_limit,
	)

	if err != nil {
		return result, r.errorResolver.ResolveError(
			err,
			r.callerOptions.From,
			nil,
			"isAncestor",
			_ancestor,
			_descendant,
			_limit,
		)
	}

	return result, err
}

func (r *Relay) IsAncestorAtBlock(
	_ancestor [32]uint8,
	_descendant [32]uint8,
	_limit *big.Int,
	blockNumber *big.Int,
) (bool, error) {
	var result bool

	err := chainutil.CallAtBlock(
		r.callerOptions.From,
		blockNumber,
		nil,
		r.contractABI,
		r.caller,
		r.errorResolver,
		r.contractAddress,
		"isAncestor",
		&result,
		_ancestor,
		_descendant,
		_limit,
	)

	return result, err
}

// ------ Events -------

func (r *Relay) Extension(
	opts *ethlike.SubscribeOpts,
	_firstFilter [][32]uint8,
	_lastFilter [][32]uint8,
) *RExtensionSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &RExtensionSubscription{
		r,
		opts,
		_firstFilter,
		_lastFilter,
	}
}

type RExtensionSubscription struct {
	contract     *Relay
	opts         *ethlike.SubscribeOpts
	_firstFilter [][32]uint8
	_lastFilter  [][32]uint8
}

type relayExtensionFunc func(
	First [32]uint8,
	Last [32]uint8,
	blockNumber uint64,
)

func (es *RExtensionSubscription) OnEvent(
	handler relayExtensionFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.RelayExtension)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.First,
					event.Last,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := es.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (es *RExtensionSubscription) Pipe(
	sink chan *abi.RelayExtension,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(es.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := es.contract.blockCounter.CurrentBlock()
				if err != nil {
					rLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - es.opts.PastBlocks

				rLogger.Infof(
					"subscription monitoring fetching past Extension events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := es.contract.PastExtensionEvents(
					fromBlock,
					nil,
					es._firstFilter,
					es._lastFilter,
				)
				if err != nil {
					rLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				rLogger.Infof(
					"subscription monitoring fetched [%v] past Extension events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := es.contract.watchExtension(
		sink,
		es._firstFilter,
		es._lastFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (r *Relay) watchExtension(
	sink chan *abi.RelayExtension,
	_firstFilter [][32]uint8,
	_lastFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return r.contract.WatchExtension(
			&bind.WatchOpts{Context: ctx},
			sink,
			_firstFilter,
			_lastFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		rLogger.Errorf(
			"subscription to event Extension had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		rLogger.Errorf(
			"subscription to event Extension failed "+
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

func (r *Relay) PastExtensionEvents(
	startBlock uint64,
	endBlock *uint64,
	_firstFilter [][32]uint8,
	_lastFilter [][32]uint8,
) ([]*abi.RelayExtension, error) {
	iterator, err := r.contract.FilterExtension(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_firstFilter,
		_lastFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past Extension events: [%v]",
			err,
		)
	}

	events := make([]*abi.RelayExtension, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}

func (r *Relay) NewTip(
	opts *ethlike.SubscribeOpts,
	_fromFilter [][32]uint8,
	_toFilter [][32]uint8,
	_gcdFilter [][32]uint8,
) *RNewTipSubscription {
	if opts == nil {
		opts = new(ethlike.SubscribeOpts)
	}
	if opts.Tick == 0 {
		opts.Tick = chainutil.DefaultSubscribeOptsTick
	}
	if opts.PastBlocks == 0 {
		opts.PastBlocks = chainutil.DefaultSubscribeOptsPastBlocks
	}

	return &RNewTipSubscription{
		r,
		opts,
		_fromFilter,
		_toFilter,
		_gcdFilter,
	}
}

type RNewTipSubscription struct {
	contract    *Relay
	opts        *ethlike.SubscribeOpts
	_fromFilter [][32]uint8
	_toFilter   [][32]uint8
	_gcdFilter  [][32]uint8
}

type relayNewTipFunc func(
	From [32]uint8,
	To [32]uint8,
	Gcd [32]uint8,
	blockNumber uint64,
)

func (nts *RNewTipSubscription) OnEvent(
	handler relayNewTipFunc,
) subscription.EventSubscription {
	eventChan := make(chan *abi.RelayNewTip)
	ctx, cancelCtx := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-eventChan:
				handler(
					event.From,
					event.To,
					event.Gcd,
					event.Raw.BlockNumber,
				)
			}
		}
	}()

	sub := nts.Pipe(eventChan)
	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (nts *RNewTipSubscription) Pipe(
	sink chan *abi.RelayNewTip,
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		ticker := time.NewTicker(nts.opts.Tick)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastBlock, err := nts.contract.blockCounter.CurrentBlock()
				if err != nil {
					rLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
				}
				fromBlock := lastBlock - nts.opts.PastBlocks

				rLogger.Infof(
					"subscription monitoring fetching past NewTip events "+
						"starting from block [%v]",
					fromBlock,
				)
				events, err := nts.contract.PastNewTipEvents(
					fromBlock,
					nil,
					nts._fromFilter,
					nts._toFilter,
					nts._gcdFilter,
				)
				if err != nil {
					rLogger.Errorf(
						"subscription failed to pull events: [%v]",
						err,
					)
					continue
				}
				rLogger.Infof(
					"subscription monitoring fetched [%v] past NewTip events",
					len(events),
				)

				for _, event := range events {
					sink <- event
				}
			}
		}
	}()

	sub := nts.contract.watchNewTip(
		sink,
		nts._fromFilter,
		nts._toFilter,
		nts._gcdFilter,
	)

	return subscription.NewEventSubscription(func() {
		sub.Unsubscribe()
		cancelCtx()
	})
}

func (r *Relay) watchNewTip(
	sink chan *abi.RelayNewTip,
	_fromFilter [][32]uint8,
	_toFilter [][32]uint8,
	_gcdFilter [][32]uint8,
) event.Subscription {
	subscribeFn := func(ctx context.Context) (event.Subscription, error) {
		return r.contract.WatchNewTip(
			&bind.WatchOpts{Context: ctx},
			sink,
			_fromFilter,
			_toFilter,
			_gcdFilter,
		)
	}

	thresholdViolatedFn := func(elapsed time.Duration) {
		rLogger.Errorf(
			"subscription to event NewTip had to be "+
				"retried [%s] since the last attempt; please inspect "+
				"host chain connectivity",
			elapsed,
		)
	}

	subscriptionFailedFn := func(err error) {
		rLogger.Errorf(
			"subscription to event NewTip failed "+
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

func (r *Relay) PastNewTipEvents(
	startBlock uint64,
	endBlock *uint64,
	_fromFilter [][32]uint8,
	_toFilter [][32]uint8,
	_gcdFilter [][32]uint8,
) ([]*abi.RelayNewTip, error) {
	iterator, err := r.contract.FilterNewTip(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		_fromFilter,
		_toFilter,
		_gcdFilter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving past NewTip events: [%v]",
			err,
		)
	}

	events := make([]*abi.RelayNewTip, 0)

	for iterator.Next() {
		event := iterator.Event
		events = append(events, event)
	}

	return events, nil
}
