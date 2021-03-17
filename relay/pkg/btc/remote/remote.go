package remote

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-remote")

// remoteChain represents a remote Bitcoin chain.
type remoteChain struct {
	rpcClient *rpcclient.Client
}

// Connect connects to the Bitcoin chain and returns a chain handle.
func Connect(config *btc.Config) (btc.Handle, error) {
	connCfg := &rpcclient.ConnConfig{
		User:         config.Username,
		Pass:         config.Password,
		Host:         config.URL,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	client, err := rpcclient.New(connCfg, nil)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to Bitcoin node at [%s]: [%v]",
			config.URL,
			err,
		)
	}

	// TODO: Remember to shutdown client
	return &remoteChain{rpcClient: client}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (rc *remoteChain) GetHeaderByHeight(height *big.Int) *btc.Header {
	return nil
}
