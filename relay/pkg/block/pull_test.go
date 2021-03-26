package block

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	btclocal "github.com/keep-network/tbtc/relay/pkg/btc/local"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestPushHeaderToQueue(t *testing.T) {
	forwarder := &Forwarder{
		headersQueue:         make(chan *btc.Header, headersQueueSize),
		nextPullHeaderHeight: 1,
	}

	headers := []*btc.Header{
		{Hash: [32]byte{1}, Height: 1, PrevHash: [32]byte{0}},
		{Hash: [32]byte{2}, Height: 2, PrevHash: [32]byte{1}},
	}

	for _, header := range headers {
		forwarder.pushHeaderToQueue(header)
	}

	// Check nextPullHeaderHeight
	expectedNextPullHeight := int64(3)
	actualNextPullHeight := forwarder.nextPullHeaderHeight

	if expectedNextPullHeight != actualNextPullHeight {
		t.Errorf(
			"unexpected add headers event:\n"+
				"expected: [%d]\n"+
				"actual:   [%d]\n",
			expectedNextPullHeight,
			actualNextPullHeight,
		)
	}

	// Check headers on the queue
	expectedQueueLength := 2
	actualQueueLength := len(forwarder.headersQueue)
	if expectedQueueLength != actualQueueLength {
		t.Errorf(
			"unexpected header queue length :\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedQueueLength,
			actualQueueLength,
		)
	}

	for _, expectedHeader := range headers {
		actualHeader := <-forwarder.headersQueue
		if !expectedHeader.Equals(actualHeader) {
			t.Errorf(
				"unexpected header in queue:\n"+
					"expected: [%+v]\n"+
					"actual:   [%+v]\n",
				expectedHeader,
				actualHeader,
			)
		}
	}
}

func TestPullHeaderFromBtcChain(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btclocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btclocal.Chain)

	btcChain.SetHeaders([]*btc.Header{
		{Height: 1, Hash: [32]byte{1}, Raw: []byte{1}},
		{Height: 2, Hash: [32]byte{2}, Raw: []byte{2}},
	})

	forwarder := &Forwarder{
		btcChain:                  btcChain,
		forwarderPullingSleepTime: 300 * time.Millisecond,
		nextPullHeaderHeight:      1,
	}

	// Test that we can pull headers without waiting when we have not reached
	// Bitcoin chain tip
	var headers []btc.Header
	for i := 0; i < 2; i++ {
		header, err := forwarder.pullHeaderFromBtcChain(ctx)
		if err != nil {
			t.Fatal(err)
		}
		headers = append(headers, *header)
		forwarder.nextPullHeaderHeight++
	}

	expectedHeaders := []btc.Header{
		{Height: 1, Hash: [32]byte{1}, Raw: []byte{1}},
		{Height: 2, Hash: [32]byte{2}, Raw: []byte{2}},
	}

	if !reflect.DeepEqual(expectedHeaders, headers) {
		t.Errorf(
			"unexpected headers:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedHeaders,
			headers,
		)
	}

	// Test that function hangs after we have reached the tip of Bitcoin chain
	// and waits for more blocks to be appended
	headerChan := make(chan *btc.Header, 1)

	go func() {
		header, err := forwarder.pullHeaderFromBtcChain(ctx)
		if err != nil {
			t.Fatal(err)
		}

		headerChan <- header
	}()

	select {
	case <-time.After(100 * time.Millisecond):
	case header := <-headerChan:
		t.Fatalf("unexpected header returned [%+v]", header)
	}

	btcChain.AppendHeader(
		&btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}},
	)

	select {
	case <-time.After(300 * time.Millisecond):
		t.Fatal("header timeout has been exceeded")
	case header := <-headerChan:
		expectedHeader := &btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}}
		if !expectedHeader.Equals(header) {
			t.Errorf(
				"unexpected header:\n"+
					"expected: [%+v]\n"+
					"actual:   [%+v]\n",
				expectedHeader,
				header,
			)
		}
	}
}

func TestFindBestHeader_HostChainReturnsBestBlock(t *testing.T) {
	bc, err := btclocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btclocal.Chain)

	btcChain.SetHeaders([]*btc.Header{
		{Height: 1, Hash: [32]byte{1}, Raw: []byte{1}},
		{Height: 2, Hash: [32]byte{2}, Raw: []byte{2}},
	})

	lc, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	localChain := lc.(*chainlocal.Chain)
	localChain.SetBestKnownDigest([32]byte{2})

	forwarder := &Forwarder{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	header, err := forwarder.findBestHeader()
	if err != nil {
		t.Fatal(err)
	}

	expectedHeader := &btc.Header{Height: 2, Hash: [32]byte{2}, Raw: []byte{2}}
	if !expectedHeader.Equals(header) {
		t.Errorf(
			"unexpected header:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedHeader,
			header,
		)
	}
}

func TestFindBestHeader_HostChainDoesNotReturnBestBlock(t *testing.T) {
	bc, err := btclocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// Create a split in chain
	//
	// 	       6 - 7  orphaned blocks
	// 	     /
	// 1 - 2
	// 	     \
	// 	       3 - 4 - 5  longest chain

	btcChain := bc.(*btclocal.Chain)

	btcChain.SetHeaders([]*btc.Header{
		{Height: 1, Hash: [32]byte{1}, PrevHash: [32]byte{0}},
		{Height: 2, Hash: [32]byte{2}, PrevHash: [32]byte{1}},
		{Height: 3, Hash: [32]byte{3}, PrevHash: [32]byte{2}},
		{Height: 4, Hash: [32]byte{4}, PrevHash: [32]byte{3}},
		{Height: 5, Hash: [32]byte{5}, PrevHash: [32]byte{4}},
	})

	btcChain.SetOrphanedHeaders([]*btc.Header{
		{Height: 3, Hash: [32]byte{6}, PrevHash: [32]byte{2}},
		{Height: 4, Hash: [32]byte{7}, PrevHash: [32]byte{6}},
	})

	lc, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	localChain := lc.(*chainlocal.Chain)
	localChain.SetBestKnownDigest([32]byte{7})

	forwarder := &Forwarder{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	header, err := forwarder.findBestHeader()
	if err != nil {
		t.Fatal(err)
	}

	// Should return header of the last common block
	expectedHeader := &btc.Header{
		Height: 2, Hash: [32]byte{2}, PrevHash: [32]byte{1},
	}
	if !expectedHeader.Equals(header) {
		t.Errorf(
			"unexpected header:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedHeader,
			header,
		)
	}
}
