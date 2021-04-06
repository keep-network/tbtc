package header

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

// push.go file contains the logic which performs the following flow:
// headersQueue -> getHeadersFromQueue -> pushHeadersToHostChain

// getHeadersFromQueue blocks until there is `headersBatchSize` headers in
// the queue or until `headerTimeout` is hit and no more headers are available
// in the queue. Normally, this function returns headers from the queue but not
// less than one and no more than `headersBatchSize` headers. Empty headers
// slice can be returned only in case the provided context is cancelled.
func (r *Relay) getHeadersFromQueue(ctx context.Context) []*btc.Header {
	headers := make([]*btc.Header, 0)

	headerTimer := time.NewTimer(headerTimeout)
	defer headerTimer.Stop()

	for len(headers) < headersBatchSize {
		logger.Debugf("waiting for new header appear on queue")

		select {
		case header := <-r.headersQueue:
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

func (r *Relay) pushHeadersToHostChain(
	ctx context.Context,
	headers []*btc.Header,
) error {
	if len(headers) == 0 {
		return nil
	}

	startMod := headers[0].Height % r.difficultyEpochDuration
	endMod := headers[len(headers)-1].Height % r.difficultyEpochDuration

	if startMod == 0 {
		// we have a difficulty change first
		logger.Info(
			"adding all headers with retarget as there is a difficulty " +
				"change at the beginning of headers batch",
		)

		if err := r.addHeadersWithRetarget(headers); err != nil {
			return fmt.Errorf("could not add headers with retarget: [%v]", err)
		}
	} else if startMod > endMod {
		// we span a difficulty change
		logger.Info(
			"adding some headers with retarget as there is a difficulty " +
				"change in the middle of headers batch",
		)

		preChangeHeaders, postChangeHeaders := r.splitBatch(headers, startMod)

		if len(preChangeHeaders) > 0 {
			if err := r.addHeaders(preChangeHeaders); err != nil {
				return fmt.Errorf("could not add headers: [%v]", err)
			}
		}

		if len(postChangeHeaders) > 0 {
			if err := r.addHeadersWithRetarget(postChangeHeaders); err != nil {
				return fmt.Errorf(
					"could not add headers with retarget: [%v]",
					err,
				)
			}
		}
	} else {
		// no difficulty change
		logger.Info(
			"adding all headers without retarget as there is no " +
				"difficulty change within headers batch",
		)

		if err := r.addHeaders(headers); err != nil {
			return fmt.Errorf("could not add headers: [%v]", err)
		}
	}

	r.processedHeaders += len(headers)
	if r.processedHeaders >= headersBatchSize {
		newBestHeader := headers[len(headers)-1]

		if err := r.updateBestHeader(ctx, newBestHeader); err != nil {
			return fmt.Errorf("could not update best header: [%v]", err)
		}

		r.processedHeaders = 0
	}

	return nil
}

func (r *Relay) addHeaders(headers []*btc.Header) error {
	anchorDigest := headers[0].PrevHash

	anchorHeader, err := r.btcChain.GetHeaderByDigest(anchorDigest)
	if err != nil {
		return fmt.Errorf(
			"could not get anchor header by digest: [%v]",
			err,
		)
	}

	return r.hostChain.AddHeaders(anchorHeader.Raw, packHeaders(headers))
}

func (r *Relay) addHeadersWithRetarget(headers []*btc.Header) error {
	epochStart := headers[0].Height - r.difficultyEpochDuration
	epochEnd := epochStart + r.difficultyEpochDuration - 1

	oldPeriodStartHeader, err := r.btcChain.GetHeaderByHeight(epochStart)
	if err != nil {
		return fmt.Errorf(
			"could not get header by height [%v]: [%v]",
			epochStart,
			err,
		)
	}

	oldPeriodEndHeader, err := r.btcChain.GetHeaderByHeight(epochEnd)
	if err != nil {
		return fmt.Errorf(
			"could not get header by height [%v]: [%v]",
			epochEnd,
			err,
		)
	}

	return r.hostChain.AddHeadersWithRetarget(
		oldPeriodStartHeader.Raw,
		oldPeriodEndHeader.Raw,
		packHeaders(headers),
	)
}

func (r *Relay) updateBestHeader(
	ctx context.Context,
	newBestHeader *btc.Header,
) error {
	for attempt := 1; attempt <= updateBestHeaderMaxAttempts; attempt++ {
		logger.Infof(
			"attempt [%v] to set header [%v] as new best",
			attempt,
			newBestHeader.Height,
		)

		currentBestDigest, err := r.hostChain.GetBestKnownDigest()
		if err != nil {
			return fmt.Errorf("could not get best known digest: [%v]", err)
		}

		currentBestHeader, err := r.btcChain.GetHeaderByDigest(
			currentBestDigest,
		)
		if err != nil {
			return fmt.Errorf(
				"could not get current best header by digest: [%v]",
				err,
			)
		}

		lastCommonAncestor, err := r.findLastCommonAncestor(
			ctx,
			newBestHeader,
			currentBestHeader,
		)
		if err != nil {
			return fmt.Errorf("could not find last common ancestor: [%v]", err)
		}

		limit := big.NewInt(
			newBestHeader.Height - lastCommonAncestor.Height + 1,
		)

		if willSucceed := r.hostChain.MarkNewHeaviestPreflight(
			lastCommonAncestor.Hash,
			currentBestHeader.Raw,
			newBestHeader.Raw,
			limit,
		); willSucceed {
			return r.hostChain.MarkNewHeaviest(
				lastCommonAncestor.Hash,
				currentBestHeader.Raw,
				newBestHeader.Raw,
				limit,
			)
		}

		// wait a constant back-off time
		select {
		case <-time.After(updateBestHeaderBackoffTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf(
		"could not set header [%v] as new best after [%v] attempts",
		newBestHeader.Height,
		updateBestHeaderMaxAttempts,
	)
}

func (r *Relay) findLastCommonAncestor(
	ctx context.Context,
	newBestHeader *btc.Header,
	currentBestHeader *btc.Header,
) (*btc.Header, error) {
	totalAttempts := 5

	for attempt := 1; attempt <= totalAttempts; attempt++ {
		logger.Infof(
			"attempt [%v] to find LCA in in previous 20 headers",
			attempt,
		)

		ancestorHeader := currentBestHeader

		for i := 0; i < 20; i++ {
			// This loop can be long-running so check the context before
			// each iteration.
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			isAncestor, err := r.hostChain.IsAncestor(
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

			ancestorHeader, err = r.btcChain.GetHeaderByDigest(
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
			return nil, ctx.Err()
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

func (r *Relay) splitBatch(headers []*btc.Header, startMod int64) (
	preChangeHeaders,
	postChangeHeaders []*btc.Header,
) {
	for _, header := range headers {
		if header.Height%r.difficultyEpochDuration >= startMod {
			preChangeHeaders = append(preChangeHeaders, header)
		} else if header.Height%r.difficultyEpochDuration < startMod {
			postChangeHeaders = append(postChangeHeaders, header)
		} else {
			logger.Errorf(
				"could not assign header [%v] to pre/post-change "+
					"part where start mod is [%v]",
				header,
				startMod,
			)
		}
	}

	return
}
