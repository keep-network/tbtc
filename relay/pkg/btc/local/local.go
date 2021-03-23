package local

import (
	"math/big"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-local")

// localChain represents a local Bitcoin chain.
type localChain struct{}

// Connect connects to the local Bitcoin chain and returns a chain handle.
func Connect() (btc.Handle, error) {
	logger.Infof("connecting local Bitcoin chain")

	return &localChain{}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (lc *localChain) GetHeaderByHeight(height *big.Int) (*btc.Header, error) {
	panic("not implemented yet")
}

// GetHeaderByDigest returns the block header for given digest (hash).
func (lc *localChain) GetHeaderByDigest(
	digest btc.Digest,
) (*btc.Header, error) {
	panic("not implemented yet")
}

// GetBlockCount returns the number of blocks in the longest blockchain
func (lc *localChain) GetBlockCount() (int64, error) {
	panic("not implemented yet")
}
