package remote

import (
	"testing"

	"github.com/keep-network/tbtc/relay/pkg/btc"
)

func TestConnectionWrongURL(t *testing.T) {
	config := &btc.Config{
		URL:      "dummyURL",
		Password: "pass",
		Username: "user",
	}

	btcChain, err := Connect(config)
	if err == nil {
		t.Fatal("No error received")
	}

	if btcChain != nil {
		t.Errorf("Non nil btc chain received")
	}
}
