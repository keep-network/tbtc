package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/tbtc/relay/pkg/metrics"

	"github.com/keep-network/tbtc/relay/pkg/btc"

	commoneth "github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/tbtc/relay/pkg/chain/ethereum"

	"github.com/keep-network/tbtc/relay/pkg/node"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/config"
	"github.com/keep-network/tbtc/relay/pkg/chain"
	"github.com/urfave/cli"
)

var logger = log.Logger("tbtc-relay-cmd")

const startDescription = `
Starts the relay maintainer in the foreground.

It requires the password of the operator host chain key file to be provided
as ` + config.PasswordEnvVariable + ` environment variable.
`

// StartCommand contains the definition of the start command-line sub-command.
var StartCommand = cli.Command{
	Name:        "start",
	Usage:       `Starts the relay maintainer in the foreground`,
	Description: startDescription,
	Action:      Start,
}

// Start starts the relay maintainer.
func Start(c *cli.Context) error {
	ctx := context.Background()

	config, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("could not read config file: [%v]", err)
	}

	btcChain, err := btc.Connect(ctx, &config.Bitcoin)
	if err != nil {
		return fmt.Errorf("could not connect BTC chain: [%v]", err)
	}

	hostChain, err := connectHostChain(config)
	if err != nil {
		return fmt.Errorf("could not connect host chain: [%v]", err)
	}

	node := node.Initialize(ctx, config, btcChain, hostChain)

	initializeMetrics(ctx, config, btcChain, hostChain, node.Stats())

	logger.Info("relay started")

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}

func connectHostChain(config *config.Config) (chain.Handle, error) {
	// TODO: add support for multiple host chains (like Celo).
	return connectEthereum(config.Ethereum)
}

func connectEthereum(config commoneth.Config) (chain.Handle, error) {
	key, err := ethutil.DecryptKeyFile(
		config.Account.KeyFile,
		config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read key file [%s]: [%v]",
			config.Account.KeyFile,
			err,
		)
	}

	return ethereum.Connect(key, &config)
}

func initializeMetrics(
	ctx context.Context,
	config *config.Config,
	btcChain btc.Handle,
	hostChain chain.Handle,
	nodeStats node.Stats,
) {
	registry, isConfigured := metrics.Initialize(
		config.Metrics.Port,
	)
	if !isConfigured {
		logger.Infof("metrics are not configured")
		return
	}

	logger.Infof(
		"enabled metrics on port [%v]",
		config.Metrics.Port,
	)

	metrics.ObserveBtcChainConnectivity(
		ctx,
		registry,
		btcChain,
		time.Duration(config.Metrics.ChainMetricsTick)*time.Second,
	)

	metrics.ObserveHostChainConnectivity(
		ctx,
		registry,
		hostChain,
		time.Duration(config.Metrics.ChainMetricsTick)*time.Second,
	)

	metrics.ObserveHeadersRelayActive(
		ctx,
		registry,
		nodeStats,
		time.Duration(config.Metrics.NodeMetricsTick)*time.Second,
	)

	metrics.ObserveHeadersRelayErrors(
		ctx,
		registry,
		nodeStats,
		time.Duration(config.Metrics.NodeMetricsTick)*time.Second,
	)

	metrics.ObserveHeadersPulled(
		ctx,
		registry,
		nodeStats,
		time.Duration(config.Metrics.NodeMetricsTick)*time.Second,
	)

	metrics.ObserveHeadersPushed(
		ctx,
		registry,
		nodeStats,
		time.Duration(config.Metrics.NodeMetricsTick)*time.Second,
	)
}
