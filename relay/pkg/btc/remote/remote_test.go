package remote

import (
	"context"
	"testing"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

func TestConnectionWrongURL(t *testing.T) {
	config := &btc.Config{
		URL:      "dummyURL",
		Password: "pass",
		Username: "user",
	}

	ctx := context.Background()

	btcChain, err := Connect(ctx, config)
	if err == nil {
		t.Fatal("No error received")
	}

	if btcChain != nil {
		t.Errorf("Non nil btc chain received")
	}
}
