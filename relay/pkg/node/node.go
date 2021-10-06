package node

import (
	"context"
	"github.com/keep-network/tbtc/relay/config"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/header"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

// Back-off time which should be applied when the relay is restarted.
// It helps to avoid being flooded with error logs in case of a permanent error
// in the relay.
const restartBackoffTime = 10 * time.Second

var logger = log.Logger("tbtc-relay-node")

// Node represents a relay node.
type Node struct {
	stats *stats
}

// Initialize initializes the relay node.
//
// TODO: This function will be probably the right place to handle relay auctions
//  which will require starting and stopping the headers relay.
func Initialize(
	ctx context.Context,
	config *config.Config,
	btcChain btc.Handle,
	hostChain chain.Handle,
) *Node {
	logger.Infof("initializing relay node")

	node := &Node{
		stats: newStats(),
	}

	go node.startRelayControlLoop(ctx, config.Relay, btcChain, hostChain)

	return node
}

// startRelayControlLoop starts a headers relay control loop which is
// responsible for starting the relay and acting upon errors by restarting
// the relay instance. The lifecycle of the control loop itself can
// be managed using the passed context.
func (n *Node) startRelayControlLoop(
	ctx context.Context,
	config config.Relay,
	btcChain btc.Handle,
	hostChain chain.Handle,
) {
	logger.Infof("starting headers relay")
	n.stats.notifyHeadersRelayActive()

	defer func() {
		logger.Infof("stopping headers relay")
		n.stats.notifyHeadersRelayInactive()
	}()

	for {
		relay := header.StartRelay(ctx, config, btcChain, hostChain, n.stats)

		select {
		case err := <-relay.ErrChan():
			logger.Errorf(
				"headers relay raised an error: [%v]",
				err,
			)

			n.stats.notifyHeadersRelayErrored()
		case <-ctx.Done():
			return
		}

		time.Sleep(restartBackoffTime)
	}
}

// Stats returns relay node statistics.
func (n *Node) Stats() Stats {
	return n.stats
}
