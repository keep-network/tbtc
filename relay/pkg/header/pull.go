package header

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

func (r *Relay) pullHeaderFromBtcChain(
	ctx context.Context,
) (*btc.Header, error) {
	for {
		chainHeight, err := r.btcChain.GetBlockCount()
		if err != nil {
			return nil, fmt.Errorf("could not get block count [%v]", err)
		}

		// Check if there are more headers to pull or we are above the chain's
		// tip and need to sleep until the Bitcoin chain adds more blocks.
		if r.nextPullHeaderHeight <= chainHeight {
			nextHeader, err := r.btcChain.GetHeaderByHeight(r.nextPullHeaderHeight)
			if err != nil {
				return nil, fmt.Errorf(
					"could not get header by height at [%d]: [%v]",
					r.nextPullHeaderHeight,
					err,
				)
			}

			// TODO: Consider just comparing hashes (it should be enough)
			//  and check if nextHeader and lastPulledHeader header can ever
			//  be equal.
			if !nextHeader.Equals(r.lastPulledHeader) {
				return nextHeader, nil
			}
		}

		select {
		case <-time.After(r.pullingSleepTime):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (r *Relay) putHeaderToQueue(header *btc.Header) {
	r.headersQueue <- header
	r.lastPulledHeader = header
	r.nextPullHeaderHeight++
}

func (r *Relay) findBestHeader() (*btc.Header, error) {
	currentBestDigest, err := r.hostChain.GetBestKnownDigest()
	if err != nil {
		return nil, err
	}

	bestHeader, err := r.btcChain.GetHeaderByDigest(currentBestDigest)
	if err != nil {
		return nil, err
	}

	// It may happen that the best header returned from the host chain is no
	// longer part of the longest Bitcoin blockchain (perhaps we registered
	// a header on the host chain and crashed and reorg happened on the Bitcoin
	// chain before we recovered from the crash).
	betterOrSameHeader, err := r.btcChain.GetHeaderByHeight(bestHeader.Height)
	if err != nil {
		return nil, err
	}

	// See if there's a better header at that height. If so, crawl backwards.
	for !bestHeader.Equals(betterOrSameHeader) {
		bestHeader, err = r.btcChain.GetHeaderByDigest(bestHeader.PrevHash)
		if err != nil {
			return nil, err
		}

		betterOrSameHeader, err = r.btcChain.GetHeaderByHeight(bestHeader.Height)
		if err != nil {
			return nil, err
		}
	}

	return bestHeader, nil
}
