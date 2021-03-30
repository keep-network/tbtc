package block

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
	"github.com/keep-network/tbtc/relay/pkg/chain"
)

// TODO: Make the following refactoring:
//  - rename `block` package to `header` (align all logging and stuff)
//  - rename `Forwarder` to `Relay`
//  - rename `RunForwarder` to `StartRelay`

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

	// Duration for which the forwarder should rest after reaching the
	// tip of Bitcoin blockchain. The forwarder waits for this time before
	// it tries to fetch a new tip from the Bitcoin blockchain, giving it
	// some time to mine new blocks.
	forwarderPullingSleepTime = 60 * time.Second

	// Maximum number of attempts which will be performed while trying
	// to update the best header.
	updateBestHeaderMaxAttempts = 30

	// Back-off time which should be applied between updating best header
	// attempts.
	updateBestHeaderBackoffTime = 10 * time.Second
)

var logger = log.Logger("relay-block-forwarder")

// Forwarder takes blocks from the Bitcoin chain and forwards them to the
// given host chain.
type Forwarder struct {
	btcChain  btc.Handle
	hostChain chain.Handle

	pullingSleepTime time.Duration

	processedHeaders     int
	nextPullHeaderHeight int64
	lastPulledHeader     *btc.Header

	headersQueue chan *btc.Header
	errChan      chan error
}

// RunForwarder creates an instance of the block forwarder and runs its
// processing loops. The lifecycle of the forwarder can be managed using the
// passed context.
func RunForwarder(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
) *Forwarder {
	loopCtx, cancelLoopCtx := context.WithCancel(ctx)

	forwarder := &Forwarder{
		btcChain:         btcChain,
		hostChain:        hostChain,
		pullingSleepTime: forwarderPullingSleepTime,
		headersQueue:     make(chan *btc.Header, headersQueueSize),
		errChan:          make(chan error, 1),
	}

	go func() {
		forwarder.pullingLoop(loopCtx)
		cancelLoopCtx() // loop exited, cancel the context
	}()

	go func() {
		forwarder.pushingLoop(loopCtx)
		cancelLoopCtx() // loop exited, cancel the context
	}()

	return forwarder
}

func (f *Forwarder) pullingLoop(ctx context.Context) {
	logger.Infof("running new block pulling loop")
	defer logger.Infof("stopping current block pulling loop")

	latestHeader, err := f.findBestHeader()
	if err != nil {
		f.errChan <- fmt.Errorf(
			"could not find best block for pulling loop: [%v]",
			err,
		)
		return
	}

	// Start pulling Bitcoin headers with the one above the latest header
	f.nextPullHeaderHeight = latestHeader.Height + 1
	logger.Infof("starting pulling from block: [%d]", latestHeader.Height+1)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Infof("pulling new header from BTC chain")

			header, err := f.pullHeaderFromBtcChain(ctx)
			if err != nil {
				f.errChan <- fmt.Errorf("could not pull header: [%v]", err)
				return
			}

			logger.Infof("pushing new header to the queue")

			f.pushHeaderToQueue(header)
		}
	}
}

func (f *Forwarder) pushingLoop(ctx context.Context) {
	logger.Infof("running new block pushing loop")
	defer logger.Infof("stopping current block pushing loop")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Infof("pulling new headers from queue")

			headers := f.pullHeadersFromQueue(ctx)
			if len(headers) == 0 {
				// Empty headers slice is returned only in case when context
				// has been cancelled.
				continue
			}

			logger.Infof("pushing %v to host chain", headersSummary(headers))

			if err := f.pushHeadersToHostChain(ctx, headers); err != nil {
				f.errChan <- fmt.Errorf("could not push headers: [%v]", err)
				// We exit on the first error letting the code controlling the
				// relay to restart it. The relay is stateful and it is easier
				// to fetch the most recent information from BTC after the
				// restart instead of trying to recover here.
				return
			}

			logger.Infof(
				"suspending block pushing loop for [%v]",
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

func headersSummary(headers []*btc.Header) string {
	if len(headers) == 0 {
		return "no headers"
	}

	firstHeaderHeight := headers[0].Height
	lastHeaderHeight := headers[len(headers)-1].Height

	if firstHeaderHeight == lastHeaderHeight {
		return fmt.Sprintf("[1] header (%v)", firstHeaderHeight)
	}

	return fmt.Sprintf(
		"[%v] headers (from %v to %v)",
		len(headers),
		firstHeaderHeight,
		lastHeaderHeight,
	)
}
