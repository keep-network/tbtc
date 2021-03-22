package block

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
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
