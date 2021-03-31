package header

import (
	"context"
	"encoding/binary"
	"math/big"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
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

func TestPushHeadersToHostChain_NoDifficultyChange(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

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

func TestPushHeadersToHostChain_NoDifficultyChange_WithUpdateBestHeader(
	t *testing.T,
) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

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

func TestPushHeadersToHostChain_DifficultyChangeAtBeginning(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

	// When adding headers with retarget, previous epoch boundary headers will
	// be searched by heights. Because of that, the local BTC chain must be
	// aware of it.
	btcChain.SetHeaders([]*btc.Header{
		{Hash: to32Bytes(2016), Height: 2016, Raw: toBytes(2016)},
		{Hash: to32Bytes(4031), Height: 4031, Raw: toBytes(4031)},
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
		// First header's height modulo epoch duration (2016) must be zero.
		{Hash: to32Bytes(4032), Height: 4032, Raw: toBytes(4032)},
		{Hash: to32Bytes(4033), Height: 4033, Raw: toBytes(4033)},
		{Hash: to32Bytes(4034), Height: 4034, Raw: toBytes(4034)},
		{Hash: to32Bytes(4035), Height: 4035, Raw: toBytes(4035)},
	}

	err = forwarder.pushHeadersToHostChain(ctx, headers)
	if err != nil {
		t.Fatal(err)
	}

	addHeadersWithRetargetEvents := localChain.AddHeadersWithRetargetEvents()

	expectedEventsCount := 1
	actualEventsCount := len(addHeadersWithRetargetEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of add headers with retarget events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedAddHeadersWithRetargetEvent := &chainlocal.AddHeadersWithRetargetEvent{
		OldPeriodStartHeader: toBytes(2016),
		OldPeriodEndHeader:   toBytes(4031),
		Headers:              toBytes(4032, 4033, 4034, 4035),
	}
	actualAddHeadersWithRetargetEvent := addHeadersWithRetargetEvents[0]
	if !reflect.DeepEqual(
		expectedAddHeadersWithRetargetEvent,
		actualAddHeadersWithRetargetEvent,
	) {
		t.Errorf(
			"unexpected add headers with retarget event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedAddHeadersWithRetargetEvent,
			actualAddHeadersWithRetargetEvent,
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

func TestPushHeadersToHostChain_DifficultyChangeInMiddle(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	bc, err := btc.ConnectLocal()
	if err != nil {
		t.Fatal(err)
	}

	btcChain := bc.(*btc.LocalChain)

	// Headers will be added with and without retarget. Because of that, BTC
	// chain must be aware about the ancestor of the first header in the
	// pre-change batch. It must be also aware of previous epoch boundary
	// headers for the post-change batch.
	btcChain.SetHeaders([]*btc.Header{
		{Hash: to32Bytes(2016), Height: 2016, Raw: toBytes(2016)},
		{Hash: to32Bytes(4029), Height: 4029, Raw: toBytes(4029)},
		{Hash: to32Bytes(4031), Height: 4031, Raw: toBytes(4031)},
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
		{Hash: to32Bytes(4030), Height: 4030, PrevHash: to32Bytes(4029), Raw: toBytes(4030)},
		{Hash: to32Bytes(4031), Height: 4031, PrevHash: to32Bytes(4030), Raw: toBytes(4031)},
		// Difficulty change occurs on the third header as their height
		// modulo 2016 is zero.
		{Hash: to32Bytes(4032), Height: 4032, PrevHash: to32Bytes(4031), Raw: toBytes(4032)},
		{Hash: to32Bytes(4033), Height: 4033, PrevHash: to32Bytes(4032), Raw: toBytes(4033)},
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
		AnchorHeader: toBytes(4029),
		Headers:      toBytes(4030, 4031),
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

	addHeadersWithRetargetEvents := localChain.AddHeadersWithRetargetEvents()

	expectedEventsCount = 1
	actualEventsCount = len(addHeadersWithRetargetEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of add headers with retarget events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedAddHeadersWithRetargetEvent := &chainlocal.AddHeadersWithRetargetEvent{
		OldPeriodStartHeader: toBytes(2016),
		OldPeriodEndHeader:   toBytes(4031),
		Headers:              toBytes(4032, 4033),
	}
	actualAddHeadersWithRetargetEvent := addHeadersWithRetargetEvents[0]
	if !reflect.DeepEqual(
		expectedAddHeadersWithRetargetEvent,
		actualAddHeadersWithRetargetEvent,
	) {
		t.Errorf(
			"unexpected add headers with retarget event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedAddHeadersWithRetargetEvent,
			actualAddHeadersWithRetargetEvent,
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

func toBytes(values ...int) []byte {
	result := make([]byte, 0)

	for _, value := range values {
		valueResult := make([]byte, 4)
		binary.LittleEndian.PutUint32(valueResult, uint32(value))
		result = append(result, valueResult...)
	}

	return result
}

func to32Bytes(value int) [32]byte {
	result := [32]byte{}
	copy(result[:], toBytes(value))
	return result
}
