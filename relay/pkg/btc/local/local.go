package local

import (
	"math/big"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-local")

// localChain represents a local Bitcoin chain.
type localChain struct {
	// TODO: implementation
}

// Connect connects to the local Bitcoin chain and returns a chain handle.
func Connect() (btc.Handle, error) {
	// TODO: implementation
	logger.Infof("connecting local Bitcoin chain")

	return &localChain{}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (lc *localChain) GetHeaderByHeight(height *big.Int) *btc.Header {
	// TODO: implementation
	return nil
}
