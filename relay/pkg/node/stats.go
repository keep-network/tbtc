package node

import "sync"

// Stats exposes statistics of the relay node.
type Stats interface {
	// BlockForwardingEnabled returns whether the block forwarding process
	// is enabled.
	BlockForwardingEnabled() bool

	// BlockForwardingErrors returns the total number of block
	// forwarding errors.
	BlockForwardingErrors() int

	// UniqueBlocksPulled returns the number of unique blocks pulled during the
	// relay node lifetime.
	UniqueBlocksPulled() int

	// UniqueBlocksPushed returns the number of unique blocks pushed during the
	// relay node lifetime.
	UniqueBlocksPushed() int
}

// stats gathers and exposes statistics of the relay node.
type stats struct {
	mutex sync.RWMutex

	blockForwardingEnabled bool
	blockForwardingErrors  int
	uniqueBlocksPulled     map[int64]bool
	uniqueBlocksPushed     map[int64]bool
}

func newStats() *stats {
	return &stats{
		uniqueBlocksPulled: make(map[int64]bool),
		uniqueBlocksPushed: make(map[int64]bool),
	}
}

func (s *stats) notifyBlockForwardingEnabled() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.blockForwardingEnabled = true
}

func (s *stats) notifyBlockForwardingDisabled() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.blockForwardingEnabled = false
}

func (s *stats) notifyBlockForwardingErrored() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.blockForwardingErrors++
}

// NotifyBlockPulled notifies about new block pulled from the Bitcoin chain.
func (s *stats) NotifyBlockPulled(blockNumber int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.uniqueBlocksPulled[blockNumber] = true
}

// NotifyBlocksPushed notifies about new block pushed to the host chain.
func (s *stats) NotifyBlocksPushed(blockNumbers []int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, blockNumber := range blockNumbers {
		s.uniqueBlocksPushed[blockNumber] = true
	}
}

// BlockForwardingEnabled returns whether the block forwarding process
// is enabled.
func (s *stats) BlockForwardingEnabled() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.blockForwardingEnabled
}

// BlockForwardingErrors returns the total number of block forwarding errors.
func (s *stats) BlockForwardingErrors() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.blockForwardingErrors
}

// UniqueBlocksPulled returns the number of unique blocks pulled during the
// relay node lifetime.
func (s *stats) UniqueBlocksPulled() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.uniqueBlocksPulled)
}

// UniqueBlocksPushed returns the number of unique blocks pushed during the
// relay node lifetime.
func (s *stats) UniqueBlocksPushed() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.uniqueBlocksPushed)
}
