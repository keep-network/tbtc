package chain

import (
	"math/big"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

// Handle represents a handle to a host chain.
type Handle interface {
	Relay
}

// Relay is an interface that provides ability to interact with Relay contract.
type Relay interface {
	// GetBestKnownDigest returns the best known digest.
	GetBestKnownDigest() (btc.Digest, error)

	// IsAncestor checks if ancestorDigest is an ancestor of the
	// descendantDigest. The limit parameter determines the number of blocks
	// to check.
	IsAncestor(
		ancestorDigest btc.Digest,
		descendantDigest btc.Digest,
		limit *big.Int,
	) (bool, error)

	// FindHeight finds the height of a header by its digest.
	FindHeight(digest btc.Digest) (*big.Int, error)

	// AddHeaders adds headers to storage after validating. The anchorHeader
	// parameter is the header immediately preceding the new chain. Headers
	// parameter should be a tightly-packed list of 80-byte Bitcoin headers.
	AddHeaders(anchorHeader []byte, headers []byte) error

	// AddHeadersWithRetarget adds headers to storage, performs additional
	// validation of retarget. The oldPeriodStartHeader is the first header in
	// the difficulty period being closed while oldPeriodEndHeader is the last.
	// Headers parameter should be a tightly-packed list of 80-byte
	// Bitcoin headers.
	AddHeadersWithRetarget(
		oldPeriodStartHeader []byte,
		oldPeriodEndHeader []byte,
		headers []byte,
	) error

	// MarkNewHeaviest gives a new starting point for the relay. The
	// ancestorDigest param is the digest of the most recent common ancestor.
	// The currentBestHeader is a 80-byte header referenced by bestKnownDigest
	// while the newBestHeader param should be the header to mark as new best.
	// Limit parameter limits the amount of traversal of the chain.
	MarkNewHeaviest(
		ancestorDigest btc.Digest,
		currentBestHeader []byte,
		newBestHeader []byte,
		limit *big.Int,
	) error

	// MarkNewHeaviestPreflight performs a preflight call of the
	// MarkNewHeaviest method to check whether its execution will
	// succeed. If the preflight call was successful, `true` is returned.
	// In case the preflight returns an error, `false` is returned.
	MarkNewHeaviestPreflight(
		ancestorDigest btc.Digest,
		currentBestHeader []byte,
		newBestHeader []byte,
		limit *big.Int,
	) bool
}
