package local

import (
	"math/big"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-chain-local")

// localChain is a local implementation of the host chain interface.
type localChain struct{}

// Connect performs initialization for communication with the local blockchain.
func Connect() (chain.Handle, error) {
	logger.Infof("connecting local host chain")

	return &localChain{}, nil
}

// GetBestKnownDigest returns the best known digest. Returned digest is
// presented in little-endian system.
func (lc *localChain) GetBestKnownDigest() ([32]uint8, error) {
	panic("not implemented yet")
}

// IsAncestor checks if a digest is an ancestor of the given descendant. The
// limit parameter determines the number of blocks to check.
func (lc *localChain) IsAncestor(
	ancestor [32]uint8,
	descendant [32]uint8,
	limit *big.Int,
) (bool, error) {
	panic("not implemented yet")
}

// FindHeight finds the height of a header by its digest.
func (lc *localChain) FindHeight(digest [32]uint8) (*big.Int, error) {
	panic("not implemented yet")
}

// AddHeaders adds headers to storage after validating. The anchor parameter is
// the header immediately preceding the new chain. Headers parameter should be
// a tightly-packed list of 80-byte Bitcoin headers.
func (lc *localChain) AddHeaders(anchor []uint8, headers []uint8) error {
	panic("not implemented yet")
}

// AddHeadersWithRetarget adds headers to storage, performs additional
// validation of retarget. The oldPeriodStartHeader is the first header in the
// difficulty period being closed while oldPeriodEndHeader is the last.
// Headers parameter should be a tightly-packed list of 80-byte Bitcoin headers.
func (lc *localChain) AddHeadersWithRetarget(
	oldPeriodStartHeader []uint8,
	oldPeriodEndHeader []uint8,
	headers []uint8,
) error {
	panic("not implemented yet")
}

// MarkNewHeaviest gives a new starting point for the relay. The ancestor param
// is the digest of the most recent common ancestor. The currentBest is a
// 80-byte header referenced by bestKnownDigest while the newBast param should
// be the header to mark as new best. Limit parameter limits the amount of
// traversal of the chain.
func (lc *localChain) MarkNewHeaviest(
	ancestor [32]uint8,
	currentBest []uint8,
	newBest []uint8,
	limit *big.Int,
) error {
	panic("not implemented yet")
}
