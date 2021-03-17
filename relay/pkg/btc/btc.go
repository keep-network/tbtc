package btc

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// Handle represents a handle to the Bitcoin chain.
type Handle interface {
	// GetHeaderByHeight returns the block header for the given block height.
	GetHeaderByHeight(height *big.Int) (*Header, error)
}

// Header represents a Bitcoin block header.
type Header struct {
	Hash       [32]byte
	Height     int64
	Prevhash   [32]byte
	MerkleRoot [32]byte
	Raw        []byte
}

func (h *Header) String() string {
	return fmt.Sprintf(
		"Hash: %s, Height: %d, PrevHash: %s, MerkleRoot: %s, Raw: %s",
		hex.EncodeToString(h.Hash[:]),
		h.Height,
		hex.EncodeToString(h.Prevhash[:]),
		hex.EncodeToString(h.MerkleRoot[:]),
		hex.EncodeToString(h.Raw),
	)
}

type Config struct {
	URL      string
	Password string
	Username string
}
