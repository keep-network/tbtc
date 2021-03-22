package block

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

const (
	// Size of the headers queue.
	headersQueueSize = 50

	// Maximum size of processed headers batch.
	headersBatchSize = 5

	// Maximum time for which the pulling process will wait for a single header
	// to be delivered by the headers queue.
	headerTimeout = 1 * time.Second

	// Block duration of a Bitcoin difficulty epoch.
	difficultyEpochDuration = 2016

	// Duration for which the forwarder should rest after performing
	// a push action.
	forwarderPushingSleepTime = 45 * time.Second
)

var logger = log.Logger("relay-block-forwarder")

// Forwarder takes blocks from the Bitcoin chain and forwards them to the
// given host chain.
type Forwarder struct {
	btcChain  btc.Handle
	hostChain chain.Handle

	processedHeaders int

	headersQueue chan *btc.Header
	errChan      chan error
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
		btcChain:     btcChain,
		hostChain:    hostChain,
		headersQueue: make(chan *btc.Header, headersQueueSize),
		errChan:      make(chan error, 1),
	}

	go forwarder.pushingLoop(ctx)

	return forwarder
}

func (f *Forwarder) pushingLoop(ctx context.Context) {
	logger.Infof("running forwarder pushing loop")

	for {
		select {
		case <-ctx.Done():
			logger.Infof("forwarder context is done")
			return
		default:
			logger.Infof("pulling new headers from queue")

			headers := f.pullHeadersFromQueue(ctx)
			if len(headers) == 0 {
				continue
			}

			logger.Infof(
				"pushing [%v] header(s) to host chain",
				len(headers),
			)

			if err := f.pushHeadersToHostChain(headers); err != nil {
				f.errChan <- fmt.Errorf("could not push headers: [%v]", err)
				return
			}

			logger.Infof(
				"suspending forwarder pushing loop for [%v]",
				forwarderPushingSleepTime,
			)

			// Sleep for a while to achieve a limited rate.
			select {
			case <-time.After(forwarderPushingSleepTime):
			case <-ctx.Done():
			}
		}
	}
}

// ErrChan returns the error channel of the forwarder. Once an error
// appears here, the forwarder loop is immediately terminated.
func (f *Forwarder) ErrChan() <-chan error {
	return f.errChan
}
