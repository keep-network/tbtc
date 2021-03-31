package header

import (
	"context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/btc"
	chainlocal "github.com/keep-network/tbtc/relay/pkg/chain/local"
)

func TestRelay_Integration(t *testing.T) {
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

	headers := []*btc.Header{
		{Hash: to32Bytes(0), Height: 0, PrevHash: to32Bytes(0), Raw: toBytes(0)},
		{Hash: to32Bytes(1), Height: 1, PrevHash: to32Bytes(0), Raw: toBytes(1)},
		{Hash: to32Bytes(2), Height: 2, PrevHash: to32Bytes(1), Raw: toBytes(2)},
		{Hash: to32Bytes(3), Height: 3, PrevHash: to32Bytes(2), Raw: toBytes(3)},
		{Hash: to32Bytes(4), Height: 4, PrevHash: to32Bytes(3), Raw: toBytes(4)},
		{Hash: to32Bytes(5), Height: 5, PrevHash: to32Bytes(4), Raw: toBytes(5)},
		{Hash: to32Bytes(6), Height: 6, PrevHash: to32Bytes(5), Raw: toBytes(6)},
		{Hash: to32Bytes(7), Height: 7, PrevHash: to32Bytes(6), Raw: toBytes(7)},
		// The first batch: difficulty change at the first header of the batch
		{Hash: to32Bytes(8), Height: 8, PrevHash: to32Bytes(7), Raw: toBytes(8)},
		{Hash: to32Bytes(9), Height: 9, PrevHash: to32Bytes(8), Raw: toBytes(9)},
		{Hash: to32Bytes(10), Height: 10, PrevHash: to32Bytes(9), Raw: toBytes(10)},
		{Hash: to32Bytes(11), Height: 11, PrevHash: to32Bytes(10), Raw: toBytes(11)},
		{Hash: to32Bytes(12), Height: 12, PrevHash: to32Bytes(11), Raw: toBytes(12)},
		// The second batch: headers span difficulty change (at header 16)
		{Hash: to32Bytes(13), Height: 13, PrevHash: to32Bytes(12), Raw: toBytes(13)},
		{Hash: to32Bytes(14), Height: 14, PrevHash: to32Bytes(13), Raw: toBytes(14)},
		{Hash: to32Bytes(15), Height: 15, PrevHash: to32Bytes(14), Raw: toBytes(15)},
		{Hash: to32Bytes(16), Height: 16, PrevHash: to32Bytes(15), Raw: toBytes(16)},
		{Hash: to32Bytes(17), Height: 17, PrevHash: to32Bytes(16), Raw: toBytes(17)},
		// The third batch: no difficulty change accross headers
		{Hash: to32Bytes(18), Height: 18, PrevHash: to32Bytes(17), Raw: toBytes(18)},
		{Hash: to32Bytes(19), Height: 19, PrevHash: to32Bytes(18), Raw: toBytes(19)},
		{Hash: to32Bytes(20), Height: 20, PrevHash: to32Bytes(19), Raw: toBytes(20)},
	}

	btcChain.SetHeaders(headers)

	// Set the best known digest to header 7, so the relay will be pulling headers
	// starting with the 8th header
	localChain.SetBestKnownDigest(to32Bytes(7))

	// Start the relay with difficulty changing at every eighth header and
	// a shortened pushing loop sleep time
	const (
		testDifficultyEpochDuration = 8
		testRelayPushingSleepTime   = 200 * time.Millisecond
	)
	startRelay(
		ctx,
		btcChain,
		localChain,
		testDifficultyEpochDuration,
		relayPullingSleepTime,
		testRelayPushingSleepTime,
		&mockObserver{},
	)
	time.Sleep(10 * time.Millisecond)

	//************ Verify processing of the first batch ************
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
		OldPeriodStartHeader: toBytes(0),
		OldPeriodEndHeader:   toBytes(7),
		Headers:              toBytes(8, 9, 10, 11, 12),
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
		AncestorDigest:    to32Bytes(7),
		CurrentBestHeader: toBytes(7),
		NewBestHeader:     toBytes(12),
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

	localChain.SetBestKnownDigest(to32Bytes(12))
	time.Sleep(testRelayPushingSleepTime)

	//************ Verify processing of the second batch ************
	addHeadersEvents := localChain.AddHeadersEvents()

	expectedEventsCount = 1
	actualEventsCount = len(addHeadersEvents)
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
		AnchorHeader: toBytes(12),
		Headers:      toBytes(13, 14, 15),
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

	addHeadersWithRetargetEvents = localChain.AddHeadersWithRetargetEvents()

	expectedEventsCount = 2
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

	expectedAddHeadersWithRetargetEvent = &chainlocal.AddHeadersWithRetargetEvent{
		OldPeriodStartHeader: toBytes(8),
		OldPeriodEndHeader:   toBytes(15),
		Headers:              toBytes(16, 17),
	}

	actualAddHeadersWithRetargetEvent = addHeadersWithRetargetEvents[1]

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

	markNewHeaviestEvents = localChain.MarkNewHeaviestEvents()

	expectedEventsCount = 2
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

	expectedMarkNewHeaviestEvent = &chainlocal.MarkNewHeaviestEvent{
		AncestorDigest:    to32Bytes(12),
		CurrentBestHeader: toBytes(12),
		NewBestHeader:     toBytes(17),
		Limit:             big.NewInt(6),
	}
	actualMarkNewHeaviestEvent = markNewHeaviestEvents[1]
	if !reflect.DeepEqual(expectedMarkNewHeaviestEvent, actualMarkNewHeaviestEvent) {
		t.Errorf(
			"unexpected mark new heaviest event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedMarkNewHeaviestEvent,
			actualMarkNewHeaviestEvent,
		)
	}

	localChain.SetBestKnownDigest(to32Bytes(17))
	time.Sleep(testRelayPushingSleepTime + headerTimeout)

	//************ Verify processing of the third batch ************
	addHeadersEvents = localChain.AddHeadersEvents()

	expectedEventsCount = 2
	actualEventsCount = len(addHeadersEvents)
	if expectedEventsCount != actualEventsCount {
		t.Fatalf(
			"unexpected number of add headers events:\n"+
				"expected: [%v]\n"+
				"actual:   [%v]\n",
			expectedEventsCount,
			actualEventsCount,
		)
	}

	expectedAddHeadersEvent = &chainlocal.AddHeadersEvent{
		AnchorHeader: toBytes(17),
		Headers:      toBytes(18, 19, 20),
	}
	actualAddHeadersEvent = addHeadersEvents[1]
	if !reflect.DeepEqual(expectedAddHeadersEvent, actualAddHeadersEvent) {
		t.Errorf(
			"unexpected add headers event:\n"+
				"expected: [%+v]\n"+
				"actual:   [%+v]\n",
			expectedAddHeadersEvent,
			actualAddHeadersEvent,
		)
	}

	// no change in MarkNewHeaviestEvents as there were only 3 headers in the batch
	markNewHeaviestEvents = localChain.MarkNewHeaviestEvents()
	expectedEventsCount = 2
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
