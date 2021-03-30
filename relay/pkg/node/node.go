package node

import (
	"context"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/block"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-node")

// Node represents a relay node.
type Node struct {
	stats *Stats
}

// Initialize initializes the relay node.
//
// TODO: This function will be probably the right place to handle relay auctions
//  which will require starting and stopping the block forwarder.
func Initialize(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) *Node {
	logger.Infof("initializing relay node")

	node := &Node{
		stats: newStats(),
	}

	go node.runForwarderControlLoop(ctx, btcChain, hostChain)

	return node
}

// runForwarderControlLoop runs a block forwarder control loop which is
// responsible for starting the forwarder and acting upon errors by restarting
// the forwarder instance. The lifecycle of the control loop itself can
// be managed using the passed context.
func (n *Node) runForwarderControlLoop(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) {
	logger.Infof("running block forwarding")
	n.stats.notifyBlockForwardingEnabled()

	defer func() {
		logger.Infof("stopping block forwarding")
		n.stats.notifyBlockForwardingDisabled()
	}()

	for {
		forwarder := block.RunForwarder(ctx, btcChain, hostChain, n.stats)

		select {
		case err := <-forwarder.ErrChan():
			logger.Errorf(
				"error occurred during block forwarding: [%v]",
				err,
			)

			n.stats.notifyBlockForwardingErrored()
		case <-ctx.Done():
			logger.Infof("stopping block forwarding")
			return
		}
	}
}

// Stats returns relay node statistics.
func (n *Node) Stats() *Stats {
	return n.stats
}
