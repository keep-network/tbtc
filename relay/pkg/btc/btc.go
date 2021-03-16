package btc

import (
	"math/big"

	"github.com/ipfs/go-log"
)

var logger = log.Logger("relay-btc")

// Header represents a Bitcoin block header.
type Header struct {
	raw        []byte
	hash       [32]byte
	height     uint64
	prevhash   [32]byte
	merkleRoot [32]byte
}

// Client exposes methods needed to interact with the BTC blockchain.
// It's defined as an interface in order to allow swapping the actual
// implementation. For example, during the integration tests, it can
// be more convenient to use a mock instead of the real network.
type Client interface {
	// TODO: implementation
}

// Chain represents a Bitcoin chain handle and exposes methods needed to
// interact with the chain.
type Chain struct {
	client Client
}

// Connect connects to the Bitcoin chain and returns a chain handle.
func Connect() (*Chain, error) {
	// TODO: implementation
	logger.Infof("connecting Bitcoin chain")

	return &Chain{nil}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (c *Chain) GetHeaderByHeight(height *big.Int) *Header {
	// TODO: implementation
	return nil
}
