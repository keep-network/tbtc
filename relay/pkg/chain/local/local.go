package local

import (
	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-chain-local")

// localChain is a local implementation of the host chain interface.
type localChain struct {
	// TODO: implementation
}

// Connect performs initialization for communication with the local blockchain.
func Connect() (chain.Handle, error) {
	// TODO: implementation
	logger.Infof("connecting local chain")

	return &localChain{}, nil
}

// GetBestKnownDigest returns the best known digest.
func (lc *localChain) GetBestKnownDigest() ([32]byte, error) {
	// TODO: implementation
	return [32]byte{}, nil
}
