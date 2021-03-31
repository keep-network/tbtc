package local

import (
	"encoding/binary"
	"math/big"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("tbtc-relay-localchain")

// Chain is a local implementation of the host chain interface.
type Chain struct {
	bestKnownDigest [32]byte

	addHeadersEvents             []*AddHeadersEvent
	addHeadersWithRetargetEvents []*AddHeadersWithRetargetEvent
	markNewHeaviestEvent         []*MarkNewHeaviestEvent
}

// Connect performs initialization for communication with the local blockchain.
func Connect() (chain.Handle, error) {
	logger.Infof("connecting local host chain")

	return &Chain{
		addHeadersEvents:             make([]*AddHeadersEvent, 0),
		addHeadersWithRetargetEvents: make([]*AddHeadersWithRetargetEvent, 0),
		markNewHeaviestEvent:         make([]*MarkNewHeaviestEvent, 0),
	}, nil
}

// GetBestKnownDigest returns the best known digest.
func (c *Chain) GetBestKnownDigest() ([32]byte, error) {
	return c.bestKnownDigest, nil
}

// IsAncestor checks if ancestorDigest is an ancestor of the descendantDigest.
// The limit parameter determines the number of blocks to check.
func (c *Chain) IsAncestor(
	ancestorDigest [32]byte,
	descendantDigest [32]byte,
	limit *big.Int,
) (bool, error) {
	// Naive implementation for testing purposes. If the int representation
	// of the descendant digest is bigger than the ancestor's one, that
	// means the condition is true.
	descendant := binary.LittleEndian.Uint32(descendantDigest[:])
	ancestor := binary.LittleEndian.Uint32(ancestorDigest[:])
	return descendant > ancestor, nil
}

// FindHeight finds the height of a header by its digest.
func (c *Chain) FindHeight(digest [32]byte) (*big.Int, error) {
	panic("not implemented yet")
}

// AddHeaders adds headers to storage after validating. The anchorHeader
// parameter is the header immediately preceding the new chain. Headers
// parameter should be a tightly-packed list of 80-byte Bitcoin headers.
func (c *Chain) AddHeaders(anchorHeader []byte, headers []byte) error {
	c.addHeadersEvents = append(
		c.addHeadersEvents,
		&AddHeadersEvent{
			AnchorHeader: anchorHeader,
			Headers:      headers,
		},
	)

	return nil
}

// AddHeadersWithRetarget adds headers to storage, performs additional
// validation of retarget. The oldPeriodStartHeader is the first header in the
// difficulty period being closed while oldPeriodEndHeader is the last.
// Headers parameter should be a tightly-packed list of 80-byte Bitcoin headers.
func (c *Chain) AddHeadersWithRetarget(
	oldPeriodStartHeader []byte,
	oldPeriodEndHeader []byte,
	headers []byte,
) error {
	c.addHeadersWithRetargetEvents = append(
		c.addHeadersWithRetargetEvents,
		&AddHeadersWithRetargetEvent{
			OldPeriodStartHeader: oldPeriodStartHeader,
			OldPeriodEndHeader:   oldPeriodEndHeader,
			Headers:              headers,
		},
	)

	return nil
}

// MarkNewHeaviest gives a new starting point for the relay. The
// ancestorDigest param is the digest of the most recent common ancestor.
// The currentBestHeader is a 80-byte header referenced by bestKnownDigest
// while the newBestHeader param should be the header to mark as new best.
// Limit parameter limits the amount of traversal of the chain.
func (c *Chain) MarkNewHeaviest(
	ancestorDigest [32]byte,
	currentBestHeader []byte,
	newBestHeader []byte,
	limit *big.Int,
) error {
	c.markNewHeaviestEvent = append(
		c.markNewHeaviestEvent,
		&MarkNewHeaviestEvent{
			AncestorDigest:    ancestorDigest,
			CurrentBestHeader: currentBestHeader,
			NewBestHeader:     newBestHeader,
			Limit:             limit,
		},
	)

	return nil
}

// MarkNewHeaviestPreflight performs a preflight call of the
// MarkNewHeaviest method to check whether its execution will
// succeed.
func (c *Chain) MarkNewHeaviestPreflight(
	ancestorDigest [32]byte,
	currentBestHeader []byte,
	newBestHeader []byte,
	limit *big.Int,
) bool {
	return true
}

// AddHeadersEvents returns all invocations of the AddHeaders method for
// testing purposes.
func (c *Chain) AddHeadersEvents() []*AddHeadersEvent {
	return c.addHeadersEvents
}

// AddHeadersWithRetargetEvents returns all invocations of the
// AddHeadersWithRetarget method for testing purposes.
func (c *Chain) AddHeadersWithRetargetEvents() []*AddHeadersWithRetargetEvent {
	return c.addHeadersWithRetargetEvents
}

// MarkNewHeaviestEvents returns all invocations of the MarkNewHeaviest method
// for testing purposes.
func (c *Chain) MarkNewHeaviestEvents() []*MarkNewHeaviestEvent {
	return c.markNewHeaviestEvent
}

// SetBestKnownDigest sets the internal best known digest for testing purposes.
func (c *Chain) SetBestKnownDigest(bestKnownDigest [32]byte) {
	c.bestKnownDigest = bestKnownDigest
}

// AddHeadersEvent represents an invocation of the AddHeaders method.
type AddHeadersEvent struct {
	AnchorHeader []byte
	Headers      []byte
}

// AddHeadersWithRetargetEvent represents an invocation of the
// AddHeadersWithRetarget method.
type AddHeadersWithRetargetEvent struct {
	OldPeriodStartHeader []byte
	OldPeriodEndHeader   []byte
	Headers              []byte
}

// MarkNewHeaviestEvent represents an invocation of the MarkNewHeaviest method.
type MarkNewHeaviestEvent struct {
	AncestorDigest    [32]byte
	CurrentBestHeader []byte
	NewBestHeader     []byte
	Limit             *big.Int
}
