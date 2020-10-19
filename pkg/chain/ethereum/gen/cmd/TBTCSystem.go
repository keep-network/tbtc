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
	"github.com/keep-network/tbtc/go/contract"

	"github.com/urfave/cli"
)

var TBTCSystemCommand cli.Command

var tBTCSystemDescription = `The t-b-t-c-system command allows calling the TBTCSystem contract on an
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
		Name:        "t-b-t-c-system",
		Usage:       `Provides access to the TBTCSystem contract.`,
		Description: tBTCSystemDescription,
		Subcommands: []cli.Command{{
			Name:      "get-undercollateralized-threshold-percent",
			Usage:     "Calls the constant method getUndercollateralizedThresholdPercent on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetUndercollateralizedThresholdPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-eth-btc-price-feed-addition-time",
			Usage:     "Calls the constant method getRemainingEthBtcPriceFeedAdditionTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingEthBtcPriceFeedAdditionTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-initial-collateralized-percent",
			Usage:     "Calls the constant method getInitialCollateralizedPercent on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetInitialCollateralizedPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-allow-new-deposits",
			Usage:     "Calls the constant method getAllowNewDeposits on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetAllowNewDeposits,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-severely-undercollateralized-threshold-percent",
			Usage:     "Calls the constant method getSeverelyUndercollateralizedThresholdPercent on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetSeverelyUndercollateralizedThresholdPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-price-feed-governance-time-delay",
			Usage:     "Calls the constant method getPriceFeedGovernanceTimeDelay on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetPriceFeedGovernanceTimeDelay,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-keep-factories-update-time",
			Usage:     "Calls the constant method getRemainingKeepFactoriesUpdateTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingKeepFactoriesUpdateTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-pause-term",
			Usage:     "Calls the constant method getRemainingPauseTerm on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingPauseTerm,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "relay",
			Usage:     "Calls the constant method relay on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsRelay,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "approved-to-log",
			Usage:     "Calls the constant method approvedToLog on the TBTCSystem contract.",
			ArgsUsage: "[_caller] ",
			Action:    tbtcsApprovedToLog,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "fetch-relay-previous-difficulty",
			Usage:     "Calls the constant method fetchRelayPreviousDifficulty on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFetchRelayPreviousDifficulty,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-maximum-lot-size",
			Usage:     "Calls the constant method getMaximumLotSize on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetMaximumLotSize,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-signer-fee-divisor-update-time",
			Usage:     "Calls the constant method getRemainingSignerFeeDivisorUpdateTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingSignerFeeDivisorUpdateTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "fetch-bitcoin-price",
			Usage:     "Calls the constant method fetchBitcoinPrice on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFetchBitcoinPrice,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "owner",
			Usage:     "Calls the constant method owner on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsOwner,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-minimum-lot-size",
			Usage:     "Calls the constant method getMinimumLotSize on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetMinimumLotSize,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "fetch-relay-current-difficulty",
			Usage:     "Calls the constant method fetchRelayCurrentDifficulty on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFetchRelayCurrentDifficulty,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-keep-factories-upgradeability-period",
			Usage:     "Calls the constant method getKeepFactoriesUpgradeabilityPeriod on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetKeepFactoriesUpgradeabilityPeriod,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "price-feed",
			Usage:     "Calls the constant method priceFeed on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsPriceFeed,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-governance-time-delay",
			Usage:     "Calls the constant method getGovernanceTimeDelay on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetGovernanceTimeDelay,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-keep-factories-upgradeability-time",
			Usage:     "Calls the constant method getRemainingKeepFactoriesUpgradeabilityTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingKeepFactoriesUpgradeabilityTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-collateralization-thresholds-update-time",
			Usage:     "Calls the constant method getRemainingCollateralizationThresholdsUpdateTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingCollateralizationThresholdsUpdateTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-remaining-lot-sizes-update-time",
			Usage:     "Calls the constant method getRemainingLotSizesUpdateTime on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetRemainingLotSizesUpdateTime,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-signer-fee-divisor",
			Usage:     "Calls the constant method getSignerFeeDivisor on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetSignerFeeDivisor,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "keep-size",
			Usage:     "Calls the constant method keepSize on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsKeepSize,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-new-deposit-fee-estimate",
			Usage:     "Calls the constant method getNewDepositFeeEstimate on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetNewDepositFeeEstimate,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "is-owner",
			Usage:     "Calls the constant method isOwner on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsIsOwner,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "keep-threshold",
			Usage:     "Calls the constant method keepThreshold on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsKeepThreshold,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-allowed-lot-sizes",
			Usage:     "Calls the constant method getAllowedLotSizes on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsGetAllowedLotSizes,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "log-setup-failed",
			Usage:     "Calls the method logSetupFailed on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsLogSetupFailed,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "begin-eth-btc-price-feed-addition",
			Usage:     "Calls the method beginEthBtcPriceFeedAddition on the TBTCSystem contract.",
			ArgsUsage: "[_ethBtcPriceFeed] ",
			Action:    tbtcsBeginEthBtcPriceFeedAddition,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "begin-keep-factories-update",
			Usage:     "Calls the method beginKeepFactoriesUpdate on the TBTCSystem contract.",
			ArgsUsage: "[_keepStakedFactory] [_fullyBackedFactory] [_factorySelector] ",
			Action:    tbtcsBeginKeepFactoriesUpdate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(3))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "refresh-minimum-bondable-value",
			Usage:     "Calls the method refreshMinimumBondableValue on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsRefreshMinimumBondableValue,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "finalize-eth-btc-price-feed-addition",
			Usage:     "Calls the method finalizeEthBtcPriceFeedAddition on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFinalizeEthBtcPriceFeedAddition,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "finalize-signer-fee-divisor-update",
			Usage:     "Calls the method finalizeSignerFeeDivisorUpdate on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFinalizeSignerFeeDivisorUpdate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-courtesy-called",
			Usage:     "Calls the method logCourtesyCalled on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsLogCourtesyCalled,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "resume-new-deposits",
			Usage:     "Calls the method resumeNewDeposits on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsResumeNewDeposits,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "finalize-lot-sizes-update",
			Usage:     "Calls the method finalizeLotSizesUpdate on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFinalizeLotSizesUpdate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "transfer-ownership",
			Usage:     "Calls the method transferOwnership on the TBTCSystem contract.",
			ArgsUsage: "[newOwner] ",
			Action:    tbtcsTransferOwnership,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-fraud-during-setup",
			Usage:     "Calls the method logFraudDuringSetup on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsLogFraudDuringSetup,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "renounce-ownership",
			Usage:     "Calls the method renounceOwnership on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsRenounceOwnership,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-liquidated",
			Usage:     "Calls the method logLiquidated on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsLogLiquidated,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-exited-courtesy-call",
			Usage:     "Calls the method logExitedCourtesyCall on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsLogExitedCourtesyCall,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "emergency-pause-new-deposits",
			Usage:     "Calls the method emergencyPauseNewDeposits on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsEmergencyPauseNewDeposits,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-funder-requested-abort",
			Usage:     "Calls the method logFunderRequestedAbort on the TBTCSystem contract.",
			ArgsUsage: "[_abortOutputScript] ",
			Action:    tbtcsLogFunderRequestedAbort,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "finalize-keep-factories-update",
			Usage:     "Calls the method finalizeKeepFactoriesUpdate on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFinalizeKeepFactoriesUpdate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "log-created",
			Usage:     "Calls the method logCreated on the TBTCSystem contract.",
			ArgsUsage: "[_keepAddress] ",
			Action:    tbtcsLogCreated,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "finalize-collateralization-thresholds-update",
			Usage:     "Calls the method finalizeCollateralizationThresholdsUpdate on the TBTCSystem contract.",
			ArgsUsage: "",
			Action:    tbtcsFinalizeCollateralizationThresholdsUpdate,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func tbtcsGetUndercollateralizedThresholdPercent(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetUndercollateralizedThresholdPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingEthBtcPriceFeedAdditionTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingEthBtcPriceFeedAdditionTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetInitialCollateralizedPercent(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetInitialCollateralizedPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetAllowNewDeposits(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetAllowNewDepositsAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetSeverelyUndercollateralizedThresholdPercent(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetSeverelyUndercollateralizedThresholdPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetPriceFeedGovernanceTimeDelay(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetPriceFeedGovernanceTimeDelayAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingKeepFactoriesUpdateTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingKeepFactoriesUpdateTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingPauseTerm(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingPauseTermAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsRelay(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.RelayAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsApprovedToLog(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsFetchRelayPreviousDifficulty(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.FetchRelayPreviousDifficultyAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetMaximumLotSize(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetMaximumLotSizeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingSignerFeeDivisorUpdateTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingSignerFeeDivisorUpdateTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsFetchBitcoinPrice(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.FetchBitcoinPriceAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsOwner(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.OwnerAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetMinimumLotSize(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetMinimumLotSizeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsFetchRelayCurrentDifficulty(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.FetchRelayCurrentDifficultyAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetKeepFactoriesUpgradeabilityPeriod(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetKeepFactoriesUpgradeabilityPeriodAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsPriceFeed(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.PriceFeedAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetGovernanceTimeDelay(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetGovernanceTimeDelayAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingKeepFactoriesUpgradeabilityTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingKeepFactoriesUpgradeabilityTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingCollateralizationThresholdsUpdateTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingCollateralizationThresholdsUpdateTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetRemainingLotSizesUpdateTime(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetRemainingLotSizesUpdateTimeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetSignerFeeDivisor(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetSignerFeeDivisorAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsKeepSize(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.KeepSizeAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetNewDepositFeeEstimate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetNewDepositFeeEstimateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsIsOwner(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.IsOwnerAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsKeepThreshold(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.KeepThresholdAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func tbtcsGetAllowedLotSizes(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	result, err := contract.GetAllowedLotSizesAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func tbtcsLogSetupFailed(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsBeginEthBtcPriceFeedAddition(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	_ethBtcPriceFeed, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _ethBtcPriceFeed, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.BeginEthBtcPriceFeedAddition(
			_ethBtcPriceFeed,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallBeginEthBtcPriceFeedAddition(
			_ethBtcPriceFeed,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsBeginKeepFactoriesUpdate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	_keepStakedFactory, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _keepStakedFactory, a address, from passed value %v",
			c.Args()[0],
		)
	}

	_fullyBackedFactory, err := ethutil.AddressFromHex(c.Args()[1])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _fullyBackedFactory, a address, from passed value %v",
			c.Args()[1],
		)
	}

	_factorySelector, err := ethutil.AddressFromHex(c.Args()[2])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _factorySelector, a address, from passed value %v",
			c.Args()[2],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.BeginKeepFactoriesUpdate(
			_keepStakedFactory,
			_fullyBackedFactory,
			_factorySelector,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallBeginKeepFactoriesUpdate(
			_keepStakedFactory,
			_fullyBackedFactory,
			_factorySelector,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsRefreshMinimumBondableValue(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RefreshMinimumBondableValue()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRefreshMinimumBondableValue(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsFinalizeEthBtcPriceFeedAddition(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FinalizeEthBtcPriceFeedAddition()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFinalizeEthBtcPriceFeedAddition(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsFinalizeSignerFeeDivisorUpdate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FinalizeSignerFeeDivisorUpdate()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFinalizeSignerFeeDivisorUpdate(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsLogCourtesyCalled(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsResumeNewDeposits(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ResumeNewDeposits()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallResumeNewDeposits(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsFinalizeLotSizesUpdate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FinalizeLotSizesUpdate()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFinalizeLotSizesUpdate(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsTransferOwnership(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	newOwner, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter newOwner, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.TransferOwnership(
			newOwner,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallTransferOwnership(
			newOwner,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsLogFraudDuringSetup(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsRenounceOwnership(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RenounceOwnership()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRenounceOwnership(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsLogLiquidated(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsLogExitedCourtesyCall(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsEmergencyPauseNewDeposits(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.EmergencyPauseNewDeposits()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallEmergencyPauseNewDeposits(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsLogFunderRequestedAbort(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsFinalizeKeepFactoriesUpdate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FinalizeKeepFactoriesUpdate()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFinalizeKeepFactoriesUpdate(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func tbtcsLogCreated(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
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

func tbtcsFinalizeCollateralizationThresholdsUpdate(c *cli.Context) error {
	contract, err := initializeTBTCSystem(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.FinalizeCollateralizationThresholdsUpdate()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallFinalizeCollateralizationThresholdsUpdate(
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

func initializeTBTCSystem(c *cli.Context) (*contract.TBTCSystem, error) {
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

	address := common.HexToAddress(config.ContractAddresses["TBTCSystem"])

	return contract.NewTBTCSystem(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
