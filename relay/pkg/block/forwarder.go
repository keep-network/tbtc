package block

import (
	"context"
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
	forwarderSleepTime = 45 * time.Second
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

	go forwarder.loop(ctx)

	return forwarder
}

func (f *Forwarder) loop(ctx context.Context) {
	logger.Infof("running forwarder loop")

	// Init the forwarder timer with a short value for the first time.
	timer := time.After(10 * time.Second)

	for {
		select {
		case <-timer:
			logger.Debugf("running forwarder iteration")

			headers := f.pullHeaders()
			if len(headers) == 0 {
				continue
			}

			logger.Infof("pushing [%v] headers", len(headers))

			f.pushHeaders(headers)

			timer = time.After(forwarderSleepTime)
		case <-ctx.Done():
			logger.Infof("forwarder loop context is done")
			return
		}
	}
}

// pullHeaders wait until we have `headersBatchSize` headers from the queue or
// until the queue fails to yield a header for `headerTimeout`.
func (f *Forwarder) pullHeaders() []*btc.Header {
	headers := make([]*btc.Header, 0)

	for len(headers) < headersBatchSize {
		select {
		case header := <-f.headersQueue:
			headers = append(headers, header)
		case <-time.After(headerTimeout):
			if len(headers) > 0 {
				break
			}
		}
	}

	return headers
}

func (f *Forwarder) pushHeaders(headers []*btc.Header) {
	if len(headers) == 0 {
		return
	}

	startDifficulty := headers[0].Height % difficultyEpochDuration
	endDifficulty := headers[len(headers)-1].Height % difficultyEpochDuration

	if startDifficulty == 0 {
		// we have a difficulty change first
		// TODO: implementation
	} else if startDifficulty > endDifficulty {
		// we span a difficulty change
		// TODO: implementation
	} else {
		// no difficulty change
		// TODO: implementation
	}

	f.processedHeaders += len(headers)
	if f.processedHeaders >= headersBatchSize {
		newBestHeader := headers[len(headers)-1]
		f.updateBestHeader(newBestHeader)
		f.processedHeaders = 0
	}
}

func (f *Forwarder) updateBestHeader(header *btc.Header) {
	// TODO: implementation
}

// ErrChan returns the error channel of the forwarder. Once an error
// appears here, the forwarder loop is immediately terminated.
func (f *Forwarder) ErrChan() <-chan error {
	return f.errChan
}
