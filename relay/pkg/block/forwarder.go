package block

import (
	"context"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

var logger = log.Logger("relay-block-forwarder")

// Forwarder takes blocks from the Bitcoin chain and forwards them to the
// given host chain.
type Forwarder struct {
	errChan chan error
}

// RunForwarder creates an instance of the block forwarder and runs its
// processing loop. The lifecycle of the forwarder loop can be managed
// using the passed context.
func RunForwarder(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) *Forwarder {
	forwarder := &Forwarder{
		errChan: make(chan error, 1),
	}

	go forwarder.loop(ctx)

	return forwarder
}

func (f *Forwarder) loop(ctx context.Context) {
	logger.Infof("running forwarder loop")

	for {
		select {
		// TODO: implementation
		case <-ctx.Done():
			logger.Infof("forwarder context is done")
			return
		}
	}
}

// ErrChan returns the error channel of the forwarder. Once an error
// appears here, the forwarder loop is immediately terminated.
func (f *Forwarder) ErrChan() <-chan error {
	return f.errChan
}
