package btc

import (
	"math/big"
)

// Handle represents a handle to the Bitcoin chain.
type Handle interface {
	// GetHeaderByHeight returns the block header for the given block height.
	GetHeaderByHeight(height *big.Int) *Header

	// GetHeaderByDigest returns the block header for given digest (hash).
	// The digest should be passed in little-endian system.
	GetHeaderByDigest(digest [32]uint8) (*Header, error)
}

// Header represents a Bitcoin block header.
type Header struct {
	// TODO: implementation
	Hash   [32]byte
	Height int64
	Raw    []byte
}
