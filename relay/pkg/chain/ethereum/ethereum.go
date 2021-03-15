package ethereum

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-chain-ethereum")

// ethereumChain is an implementation of the host chain interface for Ethereum.
type ethereumChain struct {
	// TODO: implementation
}

// Connect performs initialization for communication with Ethereum blockchain
// based on provided config.
func Connect(
	ctx context.Context,
	accountKey *keystore.Key,
	config *ethereum.Config,
) (chain.Handle, error) {
	// TODO: implementation
	logger.Infof("connecting Ethereum chain")

	return &ethereumChain{}, nil
}

// GetBestKnownDigest returns the best known digest.
func (ec *ethereumChain) GetBestKnownDigest() ([32]byte, error) {
	// TODO: implementation
	return [32]byte{}, nil
}
