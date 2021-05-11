package header

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

	// Maximum time for which the pushing process will wait for a single header
	// to be delivered by the headers queue.
	headerTimeout = 1 * time.Second

	// Block duration of a Bitcoin difficulty epoch.
	btcDifficultyEpochDuration = 2016

	// Duration for which the relay should rest after performing a push action.
	relayPushingSleepTime = 60 * time.Second

	// Duration for which the relay should rest after reaching the
	// tip of Bitcoin blockchain. The relay waits for this time before
	// it tries to fetch a new tip from the Bitcoin blockchain, giving it
	// some time to mine new blocks.
	relayPullingSleepTime = 60 * time.Second

	// Maximum number of attempts which will be performed while trying
	// to update the best header.
	updateBestHeaderMaxAttempts = 10

	// Back-off time which should be applied between updating best header
	// attempts.
	updateBestHeaderBackoffTime = 30 * time.Second
)

var logger = log.Logger("tbtc-relay-header")

// RelayObserver represents an observer of headers relay events.
type RelayObserver interface {
	// NotifyHeaderPulled notifies about new header pulled from the
	// Bitcoin chain.
	NotifyHeaderPulled(headerHeight int64)

	// NotifyHeadersPushed notifies about new headers pushed to the host chain.
	NotifyHeadersPushed(headersHeights []int64)
}

// Relay takes headers from the Bitcoin chain and relays them to the
// given host chain.
type Relay struct {
	btcChain  btc.Handle
	hostChain chain.Handle

	difficultyEpochDuration int64

	pullingSleepTime time.Duration
	pushingSleepTime time.Duration

	processedHeaders     int
	nextPullHeaderHeight int64
	lastPulledHeader     *btc.Header

	headersQueue chan *btc.Header
	errChan      chan error

	observer RelayObserver
}

// StartRelay creates an instance of the headers relay and runs its
// processing loops. The lifecycle of the relay can be managed using the
// passed context. The relay exits automatically once an error occurs.
func StartRelay(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
	observer RelayObserver,
) *Relay {
	return startRelay(
		ctx,
		btcChain,
		hostChain,
		btcDifficultyEpochDuration,
		relayPullingSleepTime,
		relayPushingSleepTime,
		observer,
	)
}

func startRelay(
	ctx context.Context,
	btcChain btc.Handle,
	hostChain chain.Handle,
	difficultyEpochDuration int64,
	pullingSleepTime time.Duration,
	pushingSleepTime time.Duration,
	observer RelayObserver,
) *Relay {
	loopCtx, cancelLoopCtx := context.WithCancel(ctx)

	relay := &Relay{
		btcChain:                btcChain,
		hostChain:               hostChain,
		difficultyEpochDuration: difficultyEpochDuration,
		pullingSleepTime:        pullingSleepTime,
		pushingSleepTime:        pushingSleepTime,
		headersQueue:            make(chan *btc.Header, headersQueueSize),
		errChan:                 make(chan error, 1),
		observer:                observer,
	}

	go func() {
		relay.pullingLoop(loopCtx)
		cancelLoopCtx() // loop exited, cancel the context
	}()

	go func() {
		relay.pushingLoop(loopCtx)
		cancelLoopCtx() // loop exited, cancel the context
	}()

	return relay
}

func (r *Relay) pullingLoop(ctx context.Context) {
	logger.Infof("starting new headers pulling loop")
	defer logger.Infof("stopping current headers pulling loop")

	latestHeader, err := r.findBestHeader()
	if err != nil {
		r.errChan <- fmt.Errorf(
			"could not find best header for pulling loop: [%v]",
			err,
		)
		return
	}

	// Start pulling Bitcoin headers with the one above the latest header
	r.nextPullHeaderHeight = latestHeader.Height + 1

	logger.Infof(
		"starting pulling from header: [%d]",
		latestHeader.Height+1,
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Infof("starting pulling header from BTC chain")

			header, err := r.pullHeaderFromBtcChain(ctx)
			if err != nil {
				r.errChan <- fmt.Errorf("could not pull header: [%v]", err)
				return
			}

			logger.Infof("pulled header [%v] from BTC chain", header.Height)

			r.putHeaderToQueue(header)

			r.observer.NotifyHeaderPulled(header.Height)
		}
	}
}

func (r *Relay) pushingLoop(ctx context.Context) {
	logger.Infof("starting new headers pushing loop")
	defer logger.Infof("stopping current headers pushing loop")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			headers := r.getHeadersFromQueue(ctx)
			if len(headers) == 0 {
				// Empty headers slice is returned only in case when context
				// has been cancelled.
				continue
			}

			logger.Infof(
				"starting pushing %v to host chain",
				headersSummary(headers),
			)

			if err := r.pushHeadersToHostChain(ctx, headers); err != nil {
				r.errChan <- fmt.Errorf("could not push headers: [%v]", err)
				// We exit on the first error letting the code controlling the
				// relay to restart it. The relay is stateful and it is easier
				// to fetch the most recent information from BTC after the
				// restart instead of trying to recover here.
				return
			}

			logger.Infof(
				"pushed %v to host chain",
				headersSummary(headers),
			)

			r.observer.NotifyHeadersPushed(headersHeights(headers))

			logger.Infof(
				"suspending headers pushing loop for [%v]",
				r.pushingSleepTime,
			)

			// Sleep for a while to achieve a limited rate.
			select {
			case <-time.After(r.pushingSleepTime):
			case <-ctx.Done():
			}
		}
	}
}

// ErrChan returns the error channel of the relay. Once an error
// appears here, all relay loops are immediately terminated.
func (r *Relay) ErrChan() <-chan error {
	return r.errChan
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

func headersHeights(headers []*btc.Header) []int64 {
	result := make([]int64, len(headers))

	for i, header := range headers {
		result[i] = header.Height
	}

	return result
}
