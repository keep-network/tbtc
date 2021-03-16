package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/logging"
	"github.com/keep-network/tbtc/relay/cmd"
	"github.com/urfave/cli"
)

var logger = log.Logger("relay-main")

const (
	logLevelEnvVariable = "LOG_LEVEL"
	defaultConfigPath   = "./config/config.toml"
)

const appDescription = `
CLI for the relay maintainer.

It requires the config file path to be passed via the '--config' flag or looks
for the config under the default '` + defaultConfigPath + `' path if the flag 
is missing.

Log level can be customized via ` + logLevelEnvVariable + ` env variable.
`

var (
	version  string
	revision string

	configPath string
)

func main() {
	resolveVersionTags()
	configureLogging()

	app := configureCLIApp()

	app.Commands = []cli.Command{
		cmd.StartCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf("could not run relay maintainer: [%v]", err)
	}
}

func resolveVersionTags() {
	if version == "" {
		version = "unknown"
	}

	if revision == "" {
		revision = "unknown"
	}
}

func configureLogging() {
	err := logging.Configure(os.Getenv(logLevelEnvVariable))

	if err != nil {
		_, _ = fmt.Fprintf(
			os.Stderr,
			"failed to configure logging: [%v]\n",
			err,
		)
		os.Exit(1)
	}
}

func configureCLIApp() *cli.App {
	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "CLI for the relay maintainer"
	app.Description = appDescription
	app.Compiled = time.Now()
	app.Version = fmt.Sprintf("%s (revision %s)", version, revision)
	app.Authors = []cli.Author{
		{
			Name:  "Keep Network",
			Email: "info@keep.network",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Value:       defaultConfigPath,
			Destination: &configPath,
			Usage:       "full path to the configuration file",
		},
	}

	return app
}
