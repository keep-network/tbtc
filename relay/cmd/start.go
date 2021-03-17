package cmd

import (
	"context"
	"fmt"

	commoneth "github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/tbtc/relay/pkg/chain/ethereum"

	"github.com/keep-network/tbtc/relay/pkg/btc/remote"
	"github.com/keep-network/tbtc/relay/pkg/node"

	"github.com/ipfs/go-log"
	"github.com/keep-network/tbtc/relay/config"
	"github.com/keep-network/tbtc/relay/pkg/chain"
	"github.com/urfave/cli"
)

var logger = log.Logger("relay-cmd")

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

	btcChain, err := remote.Connect()
	if err != nil {
		return fmt.Errorf("could not connect BTC chain: [%v]", err)
	}

	hostChain, err := connectHostChain(ctx, config)
	if err != nil {
		return fmt.Errorf("could not connect host chain: [%v]", err)
	}

	err = node.Initialize(btcChain, hostChain)
	if err != nil {
		return fmt.Errorf("could not initialize relay node: [%v]", err)
	}

	logger.Info("relay started")

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}

func connectHostChain(
	ctx context.Context,
	config *config.Config,
) (chain.Handle, error) {
	// TODO: add support for multiple host chains (like Celo).
	return connectEthereum(ctx, config.Ethereum)
}

func connectEthereum(
	ctx context.Context,
	config commoneth.Config,
) (chain.Handle, error) {
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

	return ethereum.Connect(ctx, key, &config)
}
