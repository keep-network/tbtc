package remote

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-remote")

const connectionTimeout = 3 * time.Second

// remoteChain represents a remote Bitcoin chain.
type remoteChain struct {
	client *rpcclient.Client
}

// Connect connects to the Bitcoin chain and returns a chain handle.
func Connect(
	ctx context.Context,
	config *btc.Config,
) (btc.Handle, error) {
	logger.Infof("connecting remote Bitcoin chain")

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
			"failed to create rpc client at [%s]: [%v]",
			config.URL,
			err,
		)
	}

	// when we are exiting the program, cancel all requests from the rpc client
	// and disconnect it
	go shutdownClient(ctx, client)

	err = testConnection(client, connectionTimeout)
	if err != nil {
		return nil, fmt.Errorf(
			"error while connecting to [%s]: [%v]; check if the Bitcoin node "+
				"is running and you provided correct credentials and url",
			config.URL,
			err,
		)
	}

	return &remoteChain{client: client}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (rc *remoteChain) GetHeaderByHeight(height *big.Int) (*btc.Header, error) {
	blockHash, err := rc.client.GetBlockHash(height.Int64())
	if err != nil {
		return nil, fmt.Errorf(
			"could not get block hash for height [%d]: [%v]",
			height.Int64(),
			err,
		)
	}

	blockHeader, err := rc.client.GetBlockHeader(blockHash)
	if err != nil {
		return nil, fmt.Errorf(
			"could not get block header for hash [%s]: [%v]",
			blockHash.String(),
			err,
		)
	}

	rawHeader, err := serializeHeader(blockHeader)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to serialize header for block with hash [%s]: [%v]",
			blockHash.String(),
			err,
		)
	}

	relayHeader := &btc.Header{
		Hash:       blockHeader.BlockHash(),
		PrevHash:   blockHeader.PrevBlock,
		MerkleRoot: blockHeader.MerkleRoot,
		Raw:        rawHeader,
		Height:     height.Int64(),
	}

	return relayHeader, nil
}

func testConnection(client *rpcclient.Client, timeout time.Duration) error {
	errChan := make(chan error, 1)

	go func() {
		_, err := client.GetBlockCount()
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(timeout):
		return fmt.Errorf(
			"connection timed out after [%f] seconds",
			timeout.Seconds(),
		)
	}
}

func shutdownClient(ctx context.Context, client *rpcclient.Client) {
	<-ctx.Done()
	logger.Info("Shutting down Bitcoin rpc client")
	client.Shutdown()
}

func serializeHeader(header *wire.BlockHeader) ([]byte, error) {
	var buffer bytes.Buffer

	err := header.Serialize(&buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
