package block

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

// TODO: Make a small refactor:
//  - On top of push.go and pull.go put a comment explaining the flow:
//    eg.: f.headersQueue -> pullHeadersFromQueue -> pushHeadersToHostChain
//  - Rename `pullHeadersFromQueue` to `getHeadersFromQueue`
//  - Rename `pushHeadersToQueue` to `putHeadersToQueue`

// pullHeadersFromQueue waits until we have `headersBatchSize` headers from
// the queue or until the queue fails to yield a header for
// `headerTimeout` duration and returns them from the function.
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

func (f *Forwarder) pushHeadersToHostChain(
	ctx context.Context,
	headers []*btc.Header,
) error {
	if len(headers) == 0 {
		return nil
	}

	startMod := headers[0].Height % difficultyEpochDuration
	endMod := headers[len(headers)-1].Height % difficultyEpochDuration

	if startMod == 0 {
		// we have a difficulty change first
		// TODO: implementation
	} else if startMod > endMod {
		// we span a difficulty change
		// TODO: implementation
	} else {
		// no difficulty change
		logger.Infof(
			"performing simple headers adding as difficulty doesn't " +
				"change within headers batch",
		)

		if err := f.addHeaders(headers); err != nil {
			return fmt.Errorf("could not add headers: [%v]", err)
		}
	}

	f.processedHeaders += len(headers)
	if f.processedHeaders >= headersBatchSize {
		newBestHeader := headers[len(headers)-1]

		if err := f.updateBestHeader(ctx, newBestHeader); err != nil {
			return fmt.Errorf("could not update best header: [%v]", err)
		}

		f.processedHeaders = 0
	}

	return nil
}

func (f *Forwarder) addHeaders(headers []*btc.Header) error {
	anchorDigest := headers[0].PrevHash

	anchorHeader, err := f.btcChain.GetHeaderByDigest(anchorDigest)
	if err != nil {
		return fmt.Errorf(
			"could not get anchor header by digest: [%v]",
			err,
		)
	}

	return f.hostChain.AddHeaders(anchorHeader.Raw, packHeaders(headers))
}

func (f *Forwarder) updateBestHeader(
	ctx context.Context,
	newBestHeader *btc.Header,
) error {
	currentBestDigest, err := f.hostChain.GetBestKnownDigest()
	if err != nil {
		return fmt.Errorf("could not get best known digest: [%v]", err)
	}

	currentBestHeader, err := f.btcChain.GetHeaderByDigest(currentBestDigest)
	if err != nil {
		return fmt.Errorf(
			"could not get current best header by digest: [%v]",
			err,
		)
	}

	lastCommonAncestor, err := f.findLastCommonAncestor(
		ctx,
		newBestHeader,
		currentBestHeader,
	)
	if err != nil {
		return fmt.Errorf("could not find last common ancestor: [%v]", err)
	}

	limit := newBestHeader.Height - lastCommonAncestor.Height + 1

	return f.hostChain.MarkNewHeaviest(
		lastCommonAncestor.Hash,
		currentBestHeader.Raw,
		newBestHeader.Raw,
		big.NewInt(limit),
	)
}

func (f *Forwarder) findLastCommonAncestor(
	ctx context.Context,
	newBestHeader *btc.Header,
	currentBestHeader *btc.Header,
) (*btc.Header, error) {
	totalAttempts := 5

	for attempt := 1; attempt <= totalAttempts; attempt++ {
		logger.Infof(
			"attempt [%v] to find LCA in in previous 20 blocks",
			attempt,
		)

		ancestorHeader := currentBestHeader

		for i := 0; i < 20; i++ {
			isAncestor, err := f.hostChain.IsAncestor(
				ancestorHeader.Hash,
				newBestHeader.Hash,
				big.NewInt(240), // default value used in legacy relay
			)
			if err != nil {
				return nil, fmt.Errorf("could not check ancestry: [%v]", err)
			}

			if isAncestor {
				return ancestorHeader, nil
			}

			ancestorHeader, err = f.btcChain.GetHeaderByDigest(
				ancestorHeader.PrevHash,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"could not get header by digest: [%v]",
					err,
				)
			}
		}

		// wait a constant back-off time
		select {
		case <-time.After(15 * time.Second):
		case <-ctx.Done():
		}
	}

	return nil, fmt.Errorf(
		"could not find LCA after [%v] attempts",
		totalAttempts,
	)
}

func packHeaders(headers []*btc.Header) []byte {
	packed := make([]byte, 0)

	for _, header := range headers {
		packed = append(packed, header.Raw...)
	}

	return packed
}
