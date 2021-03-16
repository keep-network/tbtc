package node

import (
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
	"github.com/keep-network/tbtc/relay/pkg/forwarder"
)

// Initialize initializes the relay node.
func Initialize(
	btcChain btc.Handle,
	hostChain chain.Handle,
) error {
	err := forwarder.Initialize(btcChain, hostChain)
	if err != nil {
		return err
	}

	return nil
}
