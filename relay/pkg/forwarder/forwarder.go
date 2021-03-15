package forwarder

import (
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

// Initialize initializes the header forwarder process.
func Initialize(
	btcChain *btc.Chain,
	hostChain chain.Handle,
) error {
	// TODO: implementation:
	//  - implement pull logic using `btcChain`
	//	- implement push logic using `hostChain`

	return nil
}
