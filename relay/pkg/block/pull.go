package block

import (
	"fmt"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

func (f *Forwarder) pullNextHeader() error {
	nextHeader, err := f.pullHeaderFromBtcChain(f.nextPullHeaderHeight)
	if err != nil {
		return fmt.Errorf(
			"could not get header by height at [%d]: [%v]",
			f.nextPullHeaderHeight,
			err,
		)
	}

	// TODO: Consider just comparing hashes - should be enough
	// and check if nextHeader and lastPulledHeader header can ever be not equal
	if !nextHeader.Equals(f.lastPulledHeader) {
		f.pushHeaderToQueue(nextHeader)
		f.lastPulledHeader = nextHeader
		f.nextPullHeaderHeight++
	}

	return nil
}

func (f *Forwarder) pullHeaderFromBtcChain(height int64) (*btc.Header, error) {
	return f.btcChain.GetHeaderByHeight(height)
}

func (f *Forwarder) pushHeaderToQueue(header *btc.Header) {
	f.headersQueue <- header
}

func (f *Forwarder) setNextPullHeaderHeight(height int64) {
	f.nextPullHeaderHeight = height
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
