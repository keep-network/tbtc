package chain

// Handle represents a handle to a host chain.
type Handle interface {
	Relay
}

// Relay is an interface that provides ability to interact with Relay contract.
type Relay interface {
	// GetBestKnownDigest returns the best known digest.
	GetBestKnownDigest() ([32]byte, error)
}
