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
	Raw        []byte
	Hash       [32]byte
	Height     uint64
	Prevhash   [32]byte
	MerkleRoot [32]byte
}

type Config struct {
	URL      string
	Password string
	Username string
}
