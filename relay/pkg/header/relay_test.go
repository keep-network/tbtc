package header

import (
	"context"
	"fmt"
	"github.com/keep-network/tbtc/relay/config"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestRelay_PullingLoop_ContextCancellationShutdown(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}
	btcChain := bc.(*btc.LocalChain)

	localChain, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// Run relay with an empty Bitcoin chain and wait for a moment so
	// the pulling loop goes to sleep
	relay := StartRelay(ctx, config.Relay{}, btcChain, localChain, &mockObserver{})
	time.Sleep(100 * time.Millisecond)

	// While the pulling loop is sleeping, add headers to Bitcoin chain and
	// cancel context
	btcChain.SetHeaders([]*btc.Header{
		{Height: 1, Hash: [32]byte{1}, PrevHash: [32]byte{0}},
		{Height: 2, Hash: [32]byte{2}, PrevHash: [32]byte{1}},
	})
	cancelCtx()
	time.Sleep(100 * time.Millisecond)

	// The relay's queue should be empty
	expectedQueueLength := 0
	actualQueueLength := len(relay.headersQueue)
	if expectedQueueLength != actualQueueLength {
		t.Errorf(
			"unexpected headers queue length:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedQueueLength,
			actualQueueLength,
		)
	}
}

func TestRelay_PullingLoop_ErrorShutdown(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

	lc, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	localChain := lc.(*chainlocal.Chain)

	// Add one header to Bitcoin chain and set a best known digest in host chain
	// that does not correspond to any header in Bitcoin chain
	btcChain.SetHeaders([]*btc.Header{
		{Hash: [32]byte{1}, Height: 1, PrevHash: [32]byte{0}},
	})

	localChain.SetBestKnownDigest([32]byte{2})

	relay := StartRelay(ctx, config.Relay{}, btcChain, localChain, &mockObserver{})

	select {
	case err = <-relay.ErrChan():
	case <-time.After(10 * time.Second):
		t.Fatal("test timeout has been exceeded")
	}

	// An error should appear as the host chain returns digest that the Bitcoin
	// chain does not recognize
	expectedError := fmt.Errorf(
		"could not find best header for pulling loop: " +
			"[no header with digest " +
			"[02000000000000000000000000000000000" +
			"00000000000000000000000000000]]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedError,
			err,
		)
	}

	// Because the pulling loop returns early, the header queue should be empty
	expectedQueueLength := 0
	actualQueueLength := len(relay.headersQueue)
	if expectedQueueLength != actualQueueLength {
		t.Errorf(
			"unexpected headers queue length:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedQueueLength,
			actualQueueLength,
		)
	}
}

func TestRelay_PushingLoop_ContextCancellationShutdown(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	btcChain, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	localChain, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	relay := StartRelay(ctx, config.Relay{}, btcChain, localChain, &mockObserver{})

	// Shutdown the pushing loop.
	cancelCtx()

	// Fill the queue with two headers batches.
	for i := 0; i < 10; i++ {
		relay.headersQueue <- &btc.Header{Height: int64(i)}
	}

	// Without the shutdown, the relay should pick at least a batch of
	// headers from the queue. Wait some time to make sure this won't happen.
	time.Sleep(1 * time.Second)

	// All headers should remain in the queue as the pushing loop has
	// been disabled.
	expectedQueueLength := 10
	actualQueueLength := len(relay.headersQueue)
	if expectedQueueLength != actualQueueLength {
		t.Errorf(
			"unexpected headers queue length:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedQueueLength,
			actualQueueLength,
		)
	}
}

func TestRelay_PushingLoop_ErrorShutdown(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

	btcChain.SetHeaders([]*btc.Header{
		{Hash: [32]byte{255}, Height: 255, PrevHash: [32]byte{254}},
	})

	localChain, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	// Setting this allows the pulling loop to start the work in a proper way
	// and prevents it to fail earlier than the planned failure of the
	// pushing loop.
	localChain.(*chainlocal.Chain).SetBestKnownDigest([32]byte{255})

	relay := StartRelay(ctx, config.Relay{}, btcChain, localChain, &mockObserver{})

	// Fill the queue with two headers batches.
	for i := 1; i <= 10; i++ {
		var hash [32]byte
		hash[31] = byte(i)

		var prevHash [32]byte
		prevHash[31] = byte(i - 1)

		header := &btc.Header{
			Hash:     hash,
			Height:   int64(i),
			PrevHash: prevHash,
		}
		relay.headersQueue <- header
	}

	select {
	case err = <-relay.ErrChan():
	case <-time.After(10 * time.Second):
		t.Fatal("test timeout has been exceeded")
	}

	// An error should appear as the BTC local chain doesn't contain a header
	// whose hash is a previous hash of one of the header passed to the queue.
	expectedError := fmt.Errorf(
		"could not push headers: " +
			"[could not add headers: " +
			"[could not get anchor header by digest: " +
			"[no header with digest [00000000000000000000000000000000000" +
			"00000000000000000000000000000]]]]",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedError,
			err,
		)
	}

	// First batch should be picked but the second should still remain as
	// the pushing loop has been disabled due to an error.
	expectedQueueLength := 5
	actualQueueLength := len(relay.headersQueue)
	if expectedQueueLength != actualQueueLength {
		t.Errorf(
			"unexpected headers queue length:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedQueueLength,
			actualQueueLength,
		)
	}
}

type mockObserver struct{}

func (mo *mockObserver) NotifyHeaderPulled(headerHeight int64) {
	// no-op
}

func (mo *mockObserver) NotifyHeadersPushed(headersHeights []int64) {
	// no-op
}
