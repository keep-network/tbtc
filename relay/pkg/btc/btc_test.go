package btc

import "testing"

func TestConnect(t *testing.T) {
	btcChain, err := Connect()
	if err != nil {
		t.Fatal(err)
	}

	if btcChain == nil {
		t.Errorf("btc chain handle is null")
	}
}
