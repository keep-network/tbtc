package block

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestForwarder_PushingLoop_ContextCancellationShutdown(t *testing.T) {
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

	forwarder := RunForwarder(ctx, btcChain, localChain)

	// Shutdown the pushing loop.
	cancelCtx()

	// Fill the queue with two headers batches.
	for i := 0; i < 10; i++ {
		forwarder.headersQueue <- &btc.Header{Height: int64(i)}
	}

	// Without the shutdown, the forwarder should pick at least a batch of
	// headers from the queue. Wait some time to make sure this won't happen.
	time.Sleep(1 * time.Second)

	// All headers should remain in the queue as the pushing loop has
	// been disabled.
	expectedQueueLength := 10
	actualQueueLength := len(forwarder.headersQueue)
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

func TestForwarder_PushingLoop_ErrorShutdown(t *testing.T) {
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

	forwarder := RunForwarder(ctx, btcChain, localChain)

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
		forwarder.headersQueue <- header
	}

	select {
	case err = <-forwarder.ErrChan():
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
	actualQueueLength := len(forwarder.headersQueue)
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