package remote

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

var logger = log.Logger("relay-btc-remote")

// remoteChain represents a remote Bitcoin chain.
type remoteChain struct {
	rpcClient *rpcclient.Client
}

func serializeHeader(header *wire.BlockHeader) ([]byte, error) {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	err := header.Serialize(writer)

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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
			"failed to create rpc client at [%s]: [%v]", config.URL, err)
	}

	// TODO: Try to find a better way to test connection
	_, err = client.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf(
			"rpc client failed to connect at [%s]: [%v]", config.URL, err)
	}

	// TODO: Remember to shutdown client
	return &remoteChain{rpcClient: client}, nil
}

// GetHeaderByHeight returns the block header for the given block height.
func (rc *remoteChain) GetHeaderByHeight(height *big.Int) (*btc.Header, error) {
	blockHash, err := rc.rpcClient.GetBlockHash(height.Int64())

	if err != nil {
		return nil, fmt.Errorf(
			"getblockhash failed for height [%d]: [%v]", height.Int64(), err)
	}

	blockHeader, err := rc.rpcClient.GetBlockHeader(blockHash)

	if err != nil {
		return nil, fmt.Errorf(
			"getblockheader failed for block with hash [%s]: [%v]",
			blockHash.String(),
			err)
	}

	rawHeader, err := serializeHeader(blockHeader)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to serialize header for block with hash [%s]: [%v]",
			blockHash.String(),
			err)
	}

	relayHeader := btc.Header{
		Hash:       blockHeader.BlockHash(),
		Prevhash:   blockHeader.PrevBlock,
		MerkleRoot: blockHeader.MerkleRoot,
		Raw:        rawHeader,
		Height:     height.Int64(),
	}

	return &relayHeader, nil
}
