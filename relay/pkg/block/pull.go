package block

import "github.com/keep-network/tbtc/relay/pkg/btc"

func (f *Forwarder) pullHeaderFromBtcNetwork(height int64) (*btc.Header, error) {
	return f.btcChain.GetHeaderByHeight(height)
}

func (f *Forwarder) pushHeaderToQueue(header *btc.Header) {
	f.headersQueue <- header
}

func (f *Forwarder) findBestBlock() (*btc.Header, error) {
	currentBestDigest, err := f.hostChain.GetBestKnownDigest()
	if err != nil {
		return nil, err
	}

	bestHeader, err := f.btcChain.GetHeaderByDigest(currentBestDigest)
	if err != nil {
		return nil, err
	}

	betterOrSameHeader, err := f.btcChain.GetHeaderByHeight(bestHeader.Height)
	if err != nil {
		return nil, err
	}

	// TODO: Is it ever possible that bestHeader and betterOrSameHeader are not
	// equal?
	// TODO: Consider just comparing hashes - it should be enough

	// see if there's a better block at that height
	// if so, crawl backwards
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
