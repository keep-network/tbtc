// Code generated - DO NOT EDIT.
// This file is a generated command and any manual changes will be lost.

package cmd

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/tbtc/bindings/go/contract"

	"github.com/urfave/cli"
)

var DepositLogCommand cli.Command

var depositLogDescription = `The deposit-log command allows calling the DepositLog contract on an
	Ethereum network. It has subcommands corresponding to each contract method,
	which respectively each take parameters based on the contract method's
	parameters.

	Subcommands will submit a non-mutating call to the network and output the
	result.

	All subcommands can be called against a specific block by passing the
	-b/--block flag.

	All subcommands can be used to investigate the result of a previous
	transaction that called that same method by passing the -t/--transaction
	flag with the transaction hash.

	Subcommands for mutating methods may be submitted as a mutating transaction
	by passing the -s/--submit flag. In this mode, this command will terminate
	successfully once the transaction has been submitted, but will not wait for
	the transaction to be included in a block. They return the transaction hash.

	Calls that require ether to be paid will get 0 ether by default, which can
	be changed by passing the -v/--value flag.`

func init() {
	AvailableCommands = append(AvailableCommands, cli.Command{
		Name:        "deposit-log",
		Usage:       `Provides access to the DepositLog contract.`,
		Description: depositLogDescription,
		Subcommands: []cli.Command{{
			Name:      "approved-to-log",
			Usage:     "Calls the constant method approvedToLog on the DepositLog contract.",
			ArgsUsage: "[_caller] ",
			Action:    dlApprovedToLog,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "log-funder-requested-abort",
			Usage:     "Calls the method logFunderRequestedAbort on the DepositLog contract.",
			ArgsUsage: "[_abortOutputScript] ",
			Action:    dlLogFunderRequestedAbort,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-liquidated",
			Usage:     "Calls the method logLiquidated on the DepositLog contract.",
			ArgsUsage: "",
			Action:    dlLogLiquidated,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-exited-courtesy-call",
			Usage:     "Calls the method logExitedCourtesyCall on the DepositLog contract.",
			ArgsUsage: "",
			Action:    dlLogExitedCourtesyCall,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-created",
			Usage:     "Calls the method logCreated on the DepositLog contract.",
			ArgsUsage: "[_keepAddress] ",
			Action:    dlLogCreated,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-setup-failed",
			Usage:     "Calls the method logSetupFailed on the DepositLog contract.",
			ArgsUsage: "",
			Action:    dlLogSetupFailed,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-fraud-during-setup",
			Usage:     "Calls the method logFraudDuringSetup on the DepositLog contract.",
			ArgsUsage: "",
			Action:    dlLogFraudDuringSetup,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-courtesy-called",
			Usage:     "Calls the method logCourtesyCalled on the DepositLog contract.",
			ArgsUsage: "",
			Action:    dlLogCourtesyCalled,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func dlApprovedToLog(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}
	_caller, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _caller, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.ApprovedToLogAtBlock(
		_caller,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func dlLogFunderRequestedAbort(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	_abortOutputScript, err := hexutil.Decode(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _abortOutputScript, a bytes, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogFunderRequestedAbort(
			_abortOutputScript,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogFunderRequestedAbort(
			_abortOutputScript,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogLiquidated(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogLiquidated()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogLiquidated(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogExitedCourtesyCall(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogExitedCourtesyCall()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogExitedCourtesyCall(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogCreated(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	_keepAddress, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _keepAddress, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogCreated(
			_keepAddress,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogCreated(
			_keepAddress,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogSetupFailed(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogSetupFailed()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogSetupFailed(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogFraudDuringSetup(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogFraudDuringSetup()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogFraudDuringSetup(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dlLogCourtesyCalled(c *cli.Context) error {
	contract, err := initializeDepositLog(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.LogCourtesyCalled()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallLogCourtesyCalled(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

/// ------------------- Initialization -------------------

func initializeDepositLog(c *cli.Context) (*contract.DepositLog, error) {
	config, err := config.ReadEthereumConfig(c.GlobalString("config"))
	if err != nil {
		return nil, fmt.Errorf("error reading Ethereum config from file: [%v]", err)
	}

	client, _, _, err := ethutil.ConnectClients(config.URL, config.URLRPC)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	key, err := ethutil.DecryptKeyFile(
		config.Account.KeyFile,
		config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KeyFile: %s: [%v]",
			config.Account.KeyFile,
			err,
		)
	}

	checkInterval := cmd.DefaultMiningCheckInterval
	maxGasPrice := cmd.DefaultMaxGasPrice
	if config.MiningCheckInterval != 0 {
		checkInterval = time.Duration(config.MiningCheckInterval) * time.Second
	}
	if config.MaxGasPrice != 0 {
		maxGasPrice = new(big.Int).SetUint64(config.MaxGasPrice)
	}

	miningWaiter := ethutil.NewMiningWaiter(client, checkInterval, maxGasPrice)

	address := common.HexToAddress(config.ContractAddresses["DepositLog"])

	return contract.NewDepositLog(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
