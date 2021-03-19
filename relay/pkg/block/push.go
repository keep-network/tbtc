package block

import (
	"context"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

// pullHeadersFromQueue waits until we have `headersBatchSize` headers from
// the queue or until the queue fails to yield a header for
// `headerTimeout` duration.
func (f *Forwarder) pullHeadersFromQueue(ctx context.Context) []*btc.Header {
	headers := make([]*btc.Header, 0)

	headerTimer := time.NewTimer(headerTimeout)
	defer headerTimer.Stop()

	for len(headers) < headersBatchSize {
		logger.Debugf("waiting for new header appear on queue")

		select {
		case header := <-f.headersQueue:
			logger.Debugf("got header (%v) from queue", header.Height)

			headers = append(headers, header)

			// Stop the timer. In case it already expired, drain the channel
			// before performing reset.
			if !headerTimer.Stop() {
				<-headerTimer.C
			}
			headerTimer.Reset(headerTimeout)
		case <-headerTimer.C:
			if len(headers) > 0 {
				logger.Debugf(
					"new header did not appear in the given timeout; " +
						"returning headers pulled so far",
				)
				return headers
			}

			logger.Debugf(
				"new header did not appear in the given timeout; " +
					"resetting timer as no headers have been pulled so far",
			)

			// Timer expired and channel is drained so one can reset directly.
			headerTimer.Reset(headerTimeout)
		case <-ctx.Done():
			return headers
		}
	}

	return headers
}

func (f *Forwarder) pushHeadersToHostChain(headers []*btc.Header) error {
	if len(headers) == 0 {
		return nil
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
		logger.Infof(
			"performing simple headers adding as difficulty doesn't " +
				"change within headers batch",
		)

		if err := f.addHeaders(headers); err != nil {
			return err
		}
	}

	f.processedHeaders += len(headers)
	if f.processedHeaders >= headersBatchSize {
		newBestHeader := headers[len(headers)-1]

		if err := f.updateBestHeader(newBestHeader); err != nil {
			return err
		}

		f.processedHeaders = 0
	}

	return nil
}

func (f *Forwarder) addHeaders(headers []*btc.Header) error {
	// TODO: Get real anchor by calling f.btcChain.GetHeaderByHash with
	//  the `PrevHash` value of the first header from headers slice.
	var anchor []uint8

	return f.hostChain.AddHeaders(anchor, packHeaders(headers))
}

func (f *Forwarder) updateBestHeader(header *btc.Header) error {
	// TODO: implementation
	return nil
}

func packHeaders(headers []*btc.Header) []uint8 {
	packed := make([]uint8, 0)

	for _, header := range headers {
		packed = append(packed, header.Raw...)
	}

	return packed
}
