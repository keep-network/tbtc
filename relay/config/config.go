package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/tbtc/relay/pkg/btc"
)

// PasswordEnvVariable environment variable name for operator key file password.
//
// #nosec G101 (look for hardcoded credentials)
// This line doesn't contain any credentials.
// It's just the name of the environment variable.
const PasswordEnvVariable = "OPERATOR_KEY_FILE_PASSWORD"

// Config is the top level config structure.
type Config struct {
	Ethereum ethereum.Config
	Bitcoin  btc.Config
	Metrics  Metrics
}

// Metrics stores meta-info about metrics.
type Metrics struct {
	Port             int
	ChainMetricsTick int
	NodeMetricsTick  int
}

// ReadConfig reads in the configuration file in .toml format. Chain key file
// password is expected to be provided as environment variable.
func ReadConfig(filePath string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, fmt.Errorf(
			"failed to decode file [%s]: [%v]",
			filePath,
			err,
		)
	}

	password := os.Getenv(PasswordEnvVariable)

	config.Ethereum.Account.KeyFilePassword = password

	return config, nil
}
