package chain

import "math/big"

// Handle represents a handle to a host chain.
type Handle interface {
	Relay
}

// Relay is an interface that provides ability to interact with Relay contract.
type Relay interface {
	// GetBestKnownDigest returns the best known digest. Returned digest is
	// presented in little-endian system.
	GetBestKnownDigest() ([32]uint8, error)

	// IsAncestor checks if a digest is an ancestor of the given descendant.
	// The limit parameter determines the number of blocks to check.
	IsAncestor(
		ancestor [32]uint8,
		descendant [32]uint8,
		limit *big.Int,
	) (bool, error)

	// FindHeight finds the height of a header by its digest.
	FindHeight(digest [32]uint8) (*big.Int, error)

	// AddHeaders adds headers to storage after validating. The anchor
	// parameter is the header immediately preceding the new chain. Headers
	// parameter should be a tightly-packed list of 80-byte Bitcoin headers.
	AddHeaders(anchor []uint8, headers []uint8) error

	// AddHeadersWithRetarget adds headers to storage, performs additional
	// validation of retarget. The oldPeriodStartHeader is the first header in
	// the difficulty period being closed while oldPeriodEndHeader is the last.
	// Headers parameter should be a tightly-packed list of 80-byte
	// Bitcoin headers.
	AddHeadersWithRetarget(
		oldPeriodStartHeader []uint8,
		oldPeriodEndHeader []uint8,
		headers []uint8,
	) error

	// MarkNewHeaviest gives a new starting point for the relay. The ancestor
	// param is the digest of the most recent common ancestor. The currentBest
	// is a 80-byte header referenced by bestKnownDigest while the newBast
	// param should be the header to mark as new best. Limit parameter limits
	// the amount of traversal of the chain.
	MarkNewHeaviest(
		ancestor [32]uint8,
		currentBest []uint8,
		newBest []uint8,
		limit *big.Int,
	) error
}
