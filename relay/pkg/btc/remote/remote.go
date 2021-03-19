package remote

import (
	"math/big"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-remote")

// remoteChain represents a remote Bitcoin chain.
type remoteChain struct {
	// TODO: implementation
}

// Connect connects to the Bitcoin chain and returns a chain handle.
func Connect() (btc.Handle, error) {
	// TODO: implementation
	logger.Infof("connecting remote Bitcoin chain")

	return &remoteChain{}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (rc *remoteChain) GetHeaderByHeight(height *big.Int) *btc.Header {
	// TODO: implementation
	return nil
}

// GetHeaderByDigest returns the block header for given digest (hash).
// The digest should be passed in little-endian system.
func (rc *remoteChain) GetHeaderByDigest(
	digest [32]uint8,
) (*btc.Header, error) {
	// TODO: implementation
	return nil, nil
}
