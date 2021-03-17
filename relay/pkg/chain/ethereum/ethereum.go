package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/chain/ethlike"
	"github.com/keep-network/tbtc/relay/pkg/chain/ethereum/gen/contract"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-chain-ethereum")

// RelayContractName defines the name of the Relay contract.
const RelayContractName = "Relay"

var (
	// DefaultMiningCheckInterval is the default interval in which transaction
	// mining status is checked. If the transaction is not mined within this
	// time, the gas price is increased and transaction is resubmitted.
	// This value can be overwritten in the configuration file.
	DefaultMiningCheckInterval = 60 * time.Second

	// DefaultMaxGasPrice specifies the default maximum gas price the client is
	// willing to pay for the transaction to be mined. The offered transaction
	// gas price can not be higher than the max gas price value. If the maximum
	// allowed gas price is reached, no further resubmission attempts are
	// performed. This value can be overwritten in the configuration file.
	DefaultMaxGasPrice = big.NewInt(1000000000000) // 1000 Gwei
)

// ethereumChain is an implementation of the host chain interface for Ethereum.
type ethereumChain struct {
	config        *ethereum.Config
	accountKey    *keystore.Key
	client        ethutil.EthereumClient
	relayContract *contract.Relay
	blockCounter  *ethlike.BlockCounter
	miningWaiter  *ethlike.MiningWaiter
	nonceManager  *ethlike.NonceManager

	// transactionMutex allows interested parties to forcibly serialize
	// transaction submission.
	//
	// When transactions are submitted, they require a valid nonce. The nonce is
	// equal to the count of transactions the account has submitted so far, and
	// for a transaction to be accepted it should be monotonically greater than
	// any previous submitted transaction. To do this, transaction submission
	// asks the Ethereum client it is connected to for the next pending nonce,
	// and uses that value for the transaction. Unfortunately, if multiple
	// transactions are submitted in short order, they may all get the same
	// nonce. Serializing submission ensures that each nonce is requested after
	// a previous transaction has been submitted.
	transactionMutex *sync.Mutex
}

// Connect performs initialization for communication with Ethereum blockchain
// based on provided config.
func Connect(
	_ context.Context,
	accountKey *keystore.Key,
	config *ethereum.Config,
) (chain.Handle, error) {
	logger.Infof("connecting Ethereum host chain")

	client, err := ethclient.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	wrappedClient := addClientWrappers(client)

	transactionMutex := &sync.Mutex{}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve Ethereum chain id: [%v]",
			err,
		)
	}

	nonceManager := ethutil.NewNonceManager(wrappedClient, accountKey.Address)

	miningWaiter := ethutil.NewMiningWaiter(
		wrappedClient,
		DefaultMiningCheckInterval,
		DefaultMaxGasPrice,
	)

	blockCounter, err := ethutil.NewBlockCounter(wrappedClient)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create Ethereum blockcounter: [%v]",
			err,
		)
	}

	relayContractAddress, err := config.ContractAddress(
		RelayContractName,
	)
	if err != nil {
		return nil, err
	}

	relayContract, err := contract.NewRelay(
		relayContractAddress,
		chainID,
		accountKey,
		wrappedClient,
		nonceManager,
		miningWaiter,
		blockCounter,
		transactionMutex,
	)
	if err != nil {
		return nil, err
	}

	return &ethereumChain{
		config:           config,
		accountKey:       accountKey,
		client:           wrappedClient,
		relayContract:    relayContract,
		blockCounter:     blockCounter,
		nonceManager:     nonceManager,
		miningWaiter:     miningWaiter,
		transactionMutex: transactionMutex,
	}, nil
}

func addClientWrappers(
	client ethutil.EthereumClient,
) ethutil.EthereumClient {
	loggingClient := ethutil.WrapCallLogging(logger, client)

	return loggingClient
}

// GetBestKnownDigest returns the best known digest.
func (ec *ethereumChain) GetBestKnownDigest() ([32]uint8, error) {
	return ec.relayContract.GetBestKnownDigest()
}
