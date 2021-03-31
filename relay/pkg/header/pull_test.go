package header

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestPutHeaderToQueue(t *testing.T) {
	relay := &Relay{
		headersQueue:         make(chan *btc.Header, headersQueueSize),
		nextPullHeaderHeight: 1,
	}

	headers := []*btc.Header{
		{Hash: [32]byte{1}, Height: 1, PrevHash: [32]byte{0}},
		{Hash: [32]byte{2}, Height: 2, PrevHash: [32]byte{1}},
	}

	for _, header := range headers {
		relay.putHeaderToQueue(header)
	}

	// Check nextPullHeaderHeight
	expectedNextPullHeight := int64(3)
	actualNextPullHeight := relay.nextPullHeaderHeight

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
	actualQueueLength := len(relay.headersQueue)
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
		actualHeader := <-relay.headersQueue
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

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

	btcChain.SetHeaders([]*btc.Header{
		{Height: 1, Hash: [32]byte{1}, Raw: []byte{1}},
		{Height: 2, Hash: [32]byte{2}, Raw: []byte{2}},
	})

	relay := &Relay{
		btcChain:             btcChain,
		pullingSleepTime:     300 * time.Millisecond,
		nextPullHeaderHeight: 1,
	}

	// Test that we can pull headers without waiting when we have not reached
	// Bitcoin chain tip
	type result struct {
		header *btc.Header
		err    error
	}

	resultChan := make(chan result, 2)

	go func() {
		for i := 0; i < 2; i++ {
			header, err := relay.pullHeaderFromBtcChain(ctx)
			resultChan <- result{header, err}
			relay.nextPullHeaderHeight++
		}
	}()

	var headers []btc.Header
	for i := 0; i < 2; i++ {
		select {
		case <-time.After(10 * time.Millisecond):
			t.Fatal("header timeout has been exceeded")
		case result := <-resultChan:
			if result.err != nil {
				t.Fatal(err)
			} else {
				headers = append(headers, *result.header)
			}
		}
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
	// and waits for more blocks to be appended.
	go func() {
		header, err := relay.pullHeaderFromBtcChain(ctx)
		resultChan <- result{header, err}
	}()

	select {
	case <-time.After(100 * time.Millisecond):
	case result := <-resultChan:
		if err != nil {
			t.Fatal(result.err)
		} else {
			t.Fatalf("unexpected header returned [%+v]", result.header)
		}
	}

	btcChain.AppendHeader(
		&btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}},
	)

	select {
	case <-time.After(300 * time.Millisecond):
		t.Fatal("header timeout has been exceeded")
	case result := <-resultChan:
		if result.err != nil {
			t.Fatal(err)
		}

		expectedHeader := &btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}}
		if !expectedHeader.Equals(result.header) {
			t.Errorf(
				"unexpected header:\n"+
					"expected: [%+v]\n"+
					"actual:   [%+v]\n",
				expectedHeader,
				result.header,
			)
		}
	}
}

func TestFindBestHeader_HostChainReturnsBestHeader(t *testing.T) {
	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

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

	relay := &Relay{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	header, err := relay.findBestHeader()
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

func TestFindBestHeader_HostChainDoesNotReturnBestHeader(t *testing.T) {
	bc, err := btc.ConnectLocal()
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

	btcChain := bc.(*btc.LocalChain)

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

	relay := &Relay{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	header, err := relay.findBestHeader()
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
