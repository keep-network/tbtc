package node

// Stats gathers and exposes statistics of the relay node.
type Stats struct {
	blockForwardingEnabled bool
	blockForwardingErrors  int
	uniqueBlocksPulled     map[int64]bool
	uniqueBlocksPushed     map[int64]bool
}

func newStats() *Stats {
	return &Stats{
		uniqueBlocksPulled: make(map[int64]bool),
		uniqueBlocksPushed: make(map[int64]bool),
	}
}

func (s *Stats) notifyBlockForwardingEnabled() {
	s.blockForwardingEnabled = true
}

func (s *Stats) notifyBlockForwardingDisabled() {
	s.blockForwardingEnabled = false
}

func (s *Stats) notifyBlockForwardingErrored() {
	s.blockForwardingErrors++
}

// NotifyBlockPulled notifies about new block pulled from the Bitcoin chain.
func (s *Stats) NotifyBlockPulled(blockNumber int64) {
	s.uniqueBlocksPulled[blockNumber] = true
}

// NotifyBlocksPushed notifies about new block pushed to the host chain.
func (s *Stats) NotifyBlocksPushed(blockNumbers []int64) {
	for _, blockNumber := range blockNumbers {
		s.uniqueBlocksPushed[blockNumber] = true
	}
}

// BlockForwardingEnabled returns whether the block forwarding process
// is enabled.
func (s *Stats) BlockForwardingEnabled() bool {
	return s.blockForwardingEnabled
}

// BlockForwardingErrors returns the total number of block forwarding errors.
func (s *Stats) BlockForwardingErrors() int {
	return s.blockForwardingErrors
}

// UniqueBlocksPulled returns the number of unique blocks pulled during the
// relay node lifetime.
func (s *Stats) UniqueBlocksPulled() int {
	return len(s.uniqueBlocksPulled)
}

// UniqueBlocksPushed returns the number of unique blocks pushed during the
// relay node lifetime.
func (s *Stats) UniqueBlocksPushed() int {
	return len(s.uniqueBlocksPushed)
}
