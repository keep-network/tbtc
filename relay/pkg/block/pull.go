package block

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

func (f *Forwarder) pullHeaderFromBtcChain(
	ctx context.Context,
) (*btc.Header, error) {
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		chainHeight, err := f.btcChain.GetBlockCount()
		if err != nil {
			return nil, fmt.Errorf("could not get block count [%v]", err)
		}

		// Check if there are more headers to pull or we are above the chain's
		// tip and need to sleep until the chain adds more headers
		if f.nextPullHeaderHeight <= chainHeight {
			nextHeader, err := f.btcChain.GetHeaderByHeight(f.nextPullHeaderHeight)
			if err != nil {
				return nil, fmt.Errorf(
					"could not get header by height at [%d]: [%v]",
					f.nextPullHeaderHeight,
					err,
				)
			}

			// TODO: Consider just comparing hashes - should be enough
			//  and check if nextHeader and lastPulledHeader header can ever
			//  be not equal
			if !nextHeader.Equals(f.lastPulledHeader) {
				return nextHeader, nil
			}
		}

		select {
		case <-time.After(f.forwarderPullingSleepTime):
		case <-ctx.Done():
		}
	}
}

func (f *Forwarder) pushHeaderToQueue(header *btc.Header) {
	f.headersQueue <- header
	f.lastPulledHeader = header
	f.nextPullHeaderHeight++
}

func (f *Forwarder) findBestHeader() (*btc.Header, error) {
	currentBestDigest, err := f.hostChain.GetBestKnownDigest()
	if err != nil {
		return nil, err
	}

	bestHeader, err := f.btcChain.GetHeaderByDigest(currentBestDigest)
	if err != nil {
		return nil, err
	}

	// It may happen that the best header returned from the host chain is no
	// longer part of the longest bitcoin block chain (perhaps we registered
	// a header on the host chain and crashed and reorg happened on the Bitcoin
	// chain before we recovered from the crash).
	betterOrSameHeader, err := f.btcChain.GetHeaderByHeight(bestHeader.Height)
	if err != nil {
		return nil, err
	}

	// See if there's a better block at that height. If so, crawl backwards.
	for !bestHeader.Equals(betterOrSameHeader) {
		bestHeader, err = f.btcChain.GetHeaderByDigest(bestHeader.PrevHash)
		if err != nil {
			return nil, err
		}

		betterOrSameHeader, err = f.btcChain.GetHeaderByHeight(bestHeader.Height)
		if err != nil {
			return nil, err
		}
	}

	return bestHeader, nil
}
