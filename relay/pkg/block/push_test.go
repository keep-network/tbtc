package block

import (
	"context"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	btclocal "github.com/keep-network/tbtc/relay/pkg/btc/local"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestPullHeadersFromQueue(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	forwarder := &Forwarder{
		headersQueue: make(chan *btc.Header, headersQueueSize),
	}

	var wg sync.WaitGroup
	wg.Add(20) // number of headers which will be sent on queue

	headersBatches := make([][]*btc.Header, 0)

	go func() {
		for ctx.Err() == nil {
			headers := forwarder.pullHeadersFromQueue(ctx)

			if len(headers) == 0 {
				continue
			}

			headersBatches = append(headersBatches, headers)

			// Each header decrements wait group counter.
			for range headers {
				wg.Done()
			}
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			// Producer sends batches with following size: 14-3-2-1
			if i == 14 || i == 17 || i == 19 {
				time.Sleep(1100 * time.Millisecond)
			}

			forwarder.headersQueue <- &btc.Header{Height: int64(i)}
		}
	}()

	wg.Wait()
	cancelCtx() // stop consumer side

	// Consumer should receive batches with following size: 5-5-4-3-2-1
	expectedHeadersBatches := map[int][]int64{
		1: {0, 1, 2, 3, 4},
		2: {5, 6, 7, 8, 9},
		3: {10, 11, 12, 13},
		4: {14, 15, 16},
		5: {17, 18},
		6: {19},
	}

	for i, headersBatch := range headersBatches {
		actualHeadersNumbers := make([]int64, 0)
		for _, header := range headersBatch {
			actualHeadersNumbers = append(actualHeadersNumbers, header.Height)
		}

		expectedHeadersNumbers, ok := expectedHeadersBatches[i+1]
		if !ok {
			t.Errorf("batch number [%v] doesn't exist", i+1)
			continue
		}

		if !reflect.DeepEqual(expectedHeadersNumbers, actualHeadersNumbers) {
			t.Errorf(
				"unexpected batch number [%v]:\n"+
					"expected: [%v]\n"+
					"actual:   [%v]\n",
				i+1,
				expectedHeadersNumbers,
				actualHeadersNumbers,
			)
		}
	}
}

func TestPushHeadersToHostChain_AddHeaders(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btclocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btclocal.Chain)

	// When adding headers, the ancestor of the first header in batch
	// will be searched by its digest. Because of that, the local BTC
	// chain must be aware of it.
	btcChain.SetHeaders([]*btc.Header{
		{Hash: [32]byte{1}, Height: 1, Raw: []byte{1}},
	})

	lc, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	localChain := lc.(*chainlocal.Chain)

	forwarder := &Forwarder{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	headers := []*btc.Header{
		// Setting the first header's ancestor by setting the right value
		// of the PrevHash field.
		{Hash: [32]byte{2}, Height: 2, PrevHash: [32]byte{1}, Raw: []byte{2}},
		{Hash: [32]byte{3}, Height: 3, PrevHash: [32]byte{2}, Raw: []byte{3}},
		{Hash: [32]byte{4}, Height: 4, PrevHash: [32]byte{3}, Raw: []byte{4}},
		{Hash: [32]byte{5}, Height: 5, PrevHash: [32]byte{4}, Raw: []byte{5}},
	}

	err = forwarder.pushHeadersToHostChain(ctx, headers)
	if err != nil {
		t.Fatal(err)
	}

	addHeadersEvents := localChain.AddHeadersEvents()

	expectedEventsCount := 1
	actualEventsCount := len(addHeadersEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of add headers events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedAddHeadersEvent := &chainlocal.AddHeadersEvent{
		AnchorHeader: []byte{1},
		Headers:      []byte{2, 3, 4, 5},
	}
	actualAddHeadersEvent := addHeadersEvents[0]
	if !reflect.DeepEqual(expectedAddHeadersEvent, actualAddHeadersEvent) {
		t.Errorf(
			"unexpected add headers event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedAddHeadersEvent,
			actualAddHeadersEvent,
		)
	}

	markNewHeaviestEvents := localChain.MarkNewHeaviestEvents()

	// Four headers have been added. This number is not bigger or equal than
	// the batch size so update of best header should not be triggered.
	expectedEventsCount = 0
	actualEventsCount = len(markNewHeaviestEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of mark new heaviest events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}
}

func TestPushHeadersToHostChain_AddHeadersWithUpdateBestHeader(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btclocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btclocal.Chain)

	// When adding headers, the ancestor of the first header in batch
	// will be searched by its digest. Because of that, the local BTC
	// chain must be aware of it.
	btcChain.SetHeaders([]*btc.Header{
		{Hash: [32]byte{1}, Height: 1, Raw: []byte{1}},
	})

	lc, err := chainlocal.Connect()
	if err != nil {
		t.Fatal(err)
	}

	localChain := lc.(*chainlocal.Chain)

	// Simulate the situation when the best known header is the ancestor
	// of the first header in the batch.
	localChain.SetBestKnownDigest([32]byte{1})

	forwarder := &Forwarder{
		btcChain:  btcChain,
		hostChain: localChain,
	}

	headers := []*btc.Header{
		// Setting the first header's ancestor by setting the right value
		// of the PrevHash field.
		{Hash: [32]byte{2}, Height: 2, PrevHash: [32]byte{1}, Raw: []byte{2}},
		{Hash: [32]byte{3}, Height: 3, PrevHash: [32]byte{2}, Raw: []byte{3}},
		{Hash: [32]byte{4}, Height: 4, PrevHash: [32]byte{3}, Raw: []byte{4}},
		{Hash: [32]byte{5}, Height: 5, PrevHash: [32]byte{4}, Raw: []byte{5}},
		{Hash: [32]byte{6}, Height: 6, PrevHash: [32]byte{5}, Raw: []byte{6}},
	}

	err = forwarder.pushHeadersToHostChain(ctx, headers)
	if err != nil {
		t.Fatal(err)
	}

	addHeadersEvents := localChain.AddHeadersEvents()

	expectedEventsCount := 1
	actualEventsCount := len(addHeadersEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of add headers events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedAddHeadersEvent := &chainlocal.AddHeadersEvent{
		AnchorHeader: []byte{1},
		Headers:      []byte{2, 3, 4, 5, 6},
	}
	actualAddHeadersEvent := addHeadersEvents[0]
	if !reflect.DeepEqual(expectedAddHeadersEvent, actualAddHeadersEvent) {
		t.Errorf(
			"unexpected add headers event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedAddHeadersEvent,
			actualAddHeadersEvent,
		)
	}

	markNewHeaviestEvents := localChain.MarkNewHeaviestEvents()

	// Five headers have been added so update of best header should be
	// triggered with the last header in batch.
	expectedEventsCount = 1
	actualEventsCount = len(markNewHeaviestEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of mark new heaviest events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedMarkNewHeaviestEvent := &chainlocal.MarkNewHeaviestEvent{
		AncestorDigest:    [32]byte{1},
		CurrentBestHeader: []byte{1},
		NewBestHeader:     []byte{6},
		Limit:             big.NewInt(6),
	}
	actualMarkNewHeaviestEvent := markNewHeaviestEvents[0]
	if !reflect.DeepEqual(expectedMarkNewHeaviestEvent, actualMarkNewHeaviestEvent) {
		t.Errorf(
			"unexpected mark new heaviest event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedMarkNewHeaviestEvent,
			actualMarkNewHeaviestEvent,
		)
	}
}
