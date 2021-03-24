package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// LocalChain represents a local Bitcoin chain.
type LocalChain struct {
	headers []*Header
}

// ConnectLocal connects to the local Bitcoin chain and returns a chain handle.
func ConnectLocal() (Handle, error) {
	logger.Infof("connecting local Bitcoin chain")

	return &LocalChain{}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (lc *LocalChain) GetHeaderByHeight(height int64) (*Header, error) {
	for _, header := range lc.headers {
		if header.Height == height {
			return header, nil
		}
	}

	return nil, fmt.Errorf("no header with height [%v]", height)
}

// GetHeaderByDigest returns the block header for given digest (hash).
func (lc *LocalChain) GetHeaderByDigest(
	digest Digest,
) (*Header, error) {
	for _, header := range lc.headers {
		if bytes.Equal(header.Hash[:], digest[:]) {
			return header, nil
		}
	}

	return nil, fmt.Errorf(
		"no header with digest [%v]",
		hex.EncodeToString(digest[:]),
	)
}

// GetBlockCount returns the number of blocks in the longest blockchain
func (lc *LocalChain) GetBlockCount() (int64, error) {
	return int64(len(lc.headers)), nil
}

// SetHeaders set internal headers for testing purposes.
func (lc *LocalChain) SetHeaders(headers []*Header) {
	lc.headers = headers
}
