package block

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	btclocal "github.com/keep-network/tbtc/relay/pkg/btc/local"
)

func TestPushHeaderToQueue(t *testing.T) {
	// TODO
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
		forwarderPullingSleepTime: 3 * time.Second,
		nextPullHeaderHeight:      1,
	}

	// Test that we can pull headers without waiting when
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
	case <-time.After(time.Second):
	case header := <-headerChan:
		t.Fatalf("unexpected header returned[%v]", header)
	}

	btcChain.AppendHeader(
		&btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}},
	)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("failed to return header")
	case header := <-headerChan:
		expectedHeader := &btc.Header{Height: 3, Hash: [32]byte{3}, Raw: []byte{3}}

		if !reflect.DeepEqual(expectedHeader, header) {
			t.Errorf(
				"unexpected header:\n"+
					"expected: [%+v]\n"+
					"actual:   [%+v]\n",
				expectedHeaders,
				header,
			)
		}
	}
}

func TestFindBestHeader(t *testing.T) {
	//TODO
}
