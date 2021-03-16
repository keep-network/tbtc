package btc

import (
	"math/big"
)

// Handle represents a handle to the Bitcoin chain.
type Handle interface {
	// GetHeaderByHeight returns the block header for the given block height.
	GetHeaderByHeight(height *big.Int) *Header
}

// Header represents a Bitcoin block header.
type Header struct {
	raw        []byte
	hash       [32]byte
	height     uint64
	prevhash   [32]byte
	merkleRoot [32]byte
}
