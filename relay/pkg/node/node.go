package node

import (
	"context"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/block"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-node")

// Initialize initializes the relay node.
//
// TODO: This function will be probably the right place to handle new
//  requirements which will require starting and stopping the block forwarder.
func Initialize(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) {
	logger.Infof("initializing relay node")

	go runForwarderControlLoop(ctx, btcChain, hostChain)
}

// runForwarderControlLoop runs a block forwarder control loop which is
// responsible for starting the forwarder and acting upon errors by restarting
// the forwarder instance. The lifecycle of the control loop itself can
// be managed using the passed context.
func runForwarderControlLoop(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) {
	logger.Infof("running forwarder control loop")

	for {
		logger.Infof("running new forwarder instance")

		forwarder := block.RunForwarder(ctx, btcChain, hostChain)

		select {
		case err := <-forwarder.ErrChan():
			logger.Errorf(
				"forwarder terminated with error: [%v] "+
					"and will be restarted immediately",
				err,
			)
		case <-ctx.Done():
			logger.Infof("forwarder control loop context is done")
			return
		}
	}
}
