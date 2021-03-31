package node

import "sync"

// Stats exposes statistics of the relay node.
type Stats interface {
	// HeadersRelayActive returns whether the headers relay process is active.
	HeadersRelayActive() bool

	// HeadersRelayErrors returns the total number of headers relay errors.
	HeadersRelayErrors() int

	// UniqueHeadersPulled returns the number of unique headers pulled during
	// the relay node lifetime.
	UniqueHeadersPulled() int

	// UniqueHeadersPushed returns the number of unique headers pushed during
	// the relay node lifetime.
	UniqueHeadersPushed() int
}

// stats gathers and exposes statistics of the relay node.
type stats struct {
	mutex sync.RWMutex

	headersRelayActive  bool
	headersRelayErrors  int
	uniqueHeadersPulled map[int64]bool
	uniqueHeadersPushed map[int64]bool
}

func newStats() *stats {
	return &stats{
		uniqueHeadersPulled: make(map[int64]bool),
		uniqueHeadersPushed: make(map[int64]bool),
	}
}

func (s *stats) notifyHeadersRelayActive() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.headersRelayActive = true
}

func (s *stats) notifyHeadersRelayInactive() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.headersRelayActive = false
}

func (s *stats) notifyHeadersRelayErrored() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.headersRelayErrors++
}

// NotifyHeaderPulled notifies about new header pulled from the Bitcoin chain.
func (s *stats) NotifyHeaderPulled(headerHeight int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.uniqueHeadersPulled[headerHeight] = true
}

// NotifyHeadersPushed notifies about new headers pushed to the host chain.
func (s *stats) NotifyHeadersPushed(headersHeights []int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, headerHeight := range headersHeights {
		s.uniqueHeadersPushed[headerHeight] = true
	}
}

// HeadersRelayActive returns whether the headers relay process is active.
func (s *stats) HeadersRelayActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.headersRelayActive
}

// HeadersRelayErrors returns the total number of headers relay errors.
func (s *stats) HeadersRelayErrors() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.headersRelayErrors
}

// UniqueHeadersPulled returns the number of unique headers pulled during
// the relay node lifetime.
func (s *stats) UniqueHeadersPulled() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.uniqueHeadersPulled)
}

// UniqueHeadersPushed returns the number of unique headers pushed during
// the relay node lifetime.
func (s *stats) UniqueHeadersPushed() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return len(s.uniqueHeadersPushed)
}
