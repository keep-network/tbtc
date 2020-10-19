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

var DepositCommand cli.Command

var depositDescription = `The deposit command allows calling the Deposit contract on an
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
		Name:        "deposit",
		Usage:       `Provides access to the Deposit contract.`,
		Description: depositDescription,
		Subcommands: []cli.Command{{
			Name:      "severely-undercollateralized-threshold-percent",
			Usage:     "Calls the constant method severelyUndercollateralizedThresholdPercent on the Deposit contract.",
			ArgsUsage: "",
			Action:    dSeverelyUndercollateralizedThresholdPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "lot-size-satoshis",
			Usage:     "Calls the constant method lotSizeSatoshis on the Deposit contract.",
			ArgsUsage: "",
			Action:    dLotSizeSatoshis,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "current-state",
			Usage:     "Calls the constant method currentState on the Deposit contract.",
			ArgsUsage: "",
			Action:    dCurrentState,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "keep-address",
			Usage:     "Calls the constant method keepAddress on the Deposit contract.",
			ArgsUsage: "",
			Action:    dKeepAddress,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "undercollateralized-threshold-percent",
			Usage:     "Calls the constant method undercollateralizedThresholdPercent on the Deposit contract.",
			ArgsUsage: "",
			Action:    dUndercollateralizedThresholdPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "withdrawable-amount",
			Usage:     "Calls the constant method withdrawableAmount on the Deposit contract.",
			ArgsUsage: "",
			Action:    dWithdrawableAmount,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "collateralization-percentage",
			Usage:     "Calls the constant method collateralizationPercentage on the Deposit contract.",
			ArgsUsage: "",
			Action:    dCollateralizationPercentage,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "funding-info",
			Usage:     "Calls the constant method fundingInfo on the Deposit contract.",
			ArgsUsage: "",
			Action:    dFundingInfo,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "auction-value",
			Usage:     "Calls the constant method auctionValue on the Deposit contract.",
			ArgsUsage: "",
			Action:    dAuctionValue,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "signer-fee-tbtc",
			Usage:     "Calls the constant method signerFeeTbtc on the Deposit contract.",
			ArgsUsage: "",
			Action:    dSignerFeeTbtc,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-owner-redemption-tbtc-requirement",
			Usage:     "Calls the constant method getOwnerRedemptionTbtcRequirement on the Deposit contract.",
			ArgsUsage: "[_redeemer] ",
			Action:    dGetOwnerRedemptionTbtcRequirement,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "initial-collateralized-percent",
			Usage:     "Calls the constant method initialCollateralizedPercent on the Deposit contract.",
			ArgsUsage: "",
			Action:    dInitialCollateralizedPercent,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "remaining-term",
			Usage:     "Calls the constant method remainingTerm on the Deposit contract.",
			ArgsUsage: "",
			Action:    dRemainingTerm,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "in-active",
			Usage:     "Calls the constant method inActive on the Deposit contract.",
			ArgsUsage: "",
			Action:    dInActive,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "lot-size-tbtc",
			Usage:     "Calls the constant method lotSizeTbtc on the Deposit contract.",
			ArgsUsage: "",
			Action:    dLotSizeTbtc,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "get-redemption-tbtc-requirement",
			Usage:     "Calls the constant method getRedemptionTbtcRequirement on the Deposit contract.",
			ArgsUsage: "[_redeemer] ",
			Action:    dGetRedemptionTbtcRequirement,
			Before:    cmd.ArgCountChecker(1),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "utxo-value",
			Usage:     "Calls the constant method utxoValue on the Deposit contract.",
			ArgsUsage: "",
			Action:    dUtxoValue,
			Before:    cmd.ArgCountChecker(0),
			Flags:     cmd.ConstFlags,
		}, {
			Name:      "exit-courtesy-call",
			Usage:     "Calls the method exitCourtesyCall on the Deposit contract.",
			ArgsUsage: "",
			Action:    dExitCourtesyCall,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-funding-timed-out",
			Usage:     "Calls the method notifyFundingTimedOut on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyFundingTimedOut,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-undercollateralized-liquidation",
			Usage:     "Calls the method notifyUndercollateralizedLiquidation on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyUndercollateralizedLiquidation,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-redemption-proof-timed-out",
			Usage:     "Calls the method notifyRedemptionProofTimedOut on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyRedemptionProofTimedOut,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-signer-setup-failed",
			Usage:     "Calls the method notifySignerSetupFailed on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifySignerSetupFailed,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "request-funder-abort",
			Usage:     "Calls the method requestFunderAbort on the Deposit contract.",
			ArgsUsage: "[_abortOutputScript] ",
			Action:    dRequestFunderAbort,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-courtesy-call-expired",
			Usage:     "Calls the method notifyCourtesyCallExpired on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyCourtesyCallExpired,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "purchase-signer-bonds-at-auction",
			Usage:     "Calls the method purchaseSignerBondsAtAuction on the Deposit contract.",
			ArgsUsage: "",
			Action:    dPurchaseSignerBondsAtAuction,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "retrieve-signer-pubkey",
			Usage:     "Calls the method retrieveSignerPubkey on the Deposit contract.",
			ArgsUsage: "",
			Action:    dRetrieveSignerPubkey,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "initialize",
			Usage:     "Calls the method initialize on the Deposit contract.",
			ArgsUsage: "[_factory] ",
			Action:    dInitialize,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(1))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-redemption-signature-timed-out",
			Usage:     "Calls the method notifyRedemptionSignatureTimedOut on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyRedemptionSignatureTimedOut,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "notify-courtesy-call",
			Usage:     "Calls the method notifyCourtesyCall on the Deposit contract.",
			ArgsUsage: "",
			Action:    dNotifyCourtesyCall,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}, {
			Name:      "withdraw-funds",
			Usage:     "Calls the method withdrawFunds on the Deposit contract.",
			ArgsUsage: "",
			Action:    dWithdrawFunds,
			Before:    cli.BeforeFunc(cmd.NonConstArgsChecker.AndThen(cmd.ArgCountChecker(0))),
			Flags:     cmd.NonConstFlags,
		}},
	})
}

/// ------------------- Const methods -------------------

func dSeverelyUndercollateralizedThresholdPercent(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.SeverelyUndercollateralizedThresholdPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dLotSizeSatoshis(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.LotSizeSatoshisAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dCurrentState(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.CurrentStateAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dKeepAddress(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.KeepAddressAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dUndercollateralizedThresholdPercent(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.UndercollateralizedThresholdPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dWithdrawableAmount(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.WithdrawableAmountAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dCollateralizationPercentage(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.CollateralizationPercentageAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dFundingInfo(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.FundingInfoAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dAuctionValue(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.AuctionValueAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dSignerFeeTbtc(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.SignerFeeTbtcAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dGetOwnerRedemptionTbtcRequirement(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}
	_redeemer, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _redeemer, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetOwnerRedemptionTbtcRequirementAtBlock(
		_redeemer,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dInitialCollateralizedPercent(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.InitialCollateralizedPercentAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dRemainingTerm(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.RemainingTermAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dInActive(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.InActiveAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dLotSizeTbtc(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.LotSizeTbtcAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dGetRedemptionTbtcRequirement(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}
	_redeemer, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _redeemer, a address, from passed value %v",
			c.Args()[0],
		)
	}

	result, err := contract.GetRedemptionTbtcRequirementAtBlock(
		_redeemer,

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

func dUtxoValue(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	result, err := contract.UtxoValueAtBlock(

		cmd.BlockFlagValue.Uint,
	)

	if err != nil {
		return err
	}

	cmd.PrintOutput(result)

	return nil
}

/// ------------------- Non-const methods -------------------

func dExitCourtesyCall(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.ExitCourtesyCall()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallExitCourtesyCall(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifyFundingTimedOut(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyFundingTimedOut()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyFundingTimedOut(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifyUndercollateralizedLiquidation(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyUndercollateralizedLiquidation()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyUndercollateralizedLiquidation(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifyRedemptionProofTimedOut(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyRedemptionProofTimedOut()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyRedemptionProofTimedOut(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifySignerSetupFailed(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifySignerSetupFailed()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifySignerSetupFailed(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dRequestFunderAbort(c *cli.Context) error {
	contract, err := initializeDeposit(c)
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
		transaction, err = contract.RequestFunderAbort(
			_abortOutputScript,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRequestFunderAbort(
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

func dNotifyCourtesyCallExpired(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyCourtesyCallExpired()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyCourtesyCallExpired(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dPurchaseSignerBondsAtAuction(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.PurchaseSignerBondsAtAuction()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallPurchaseSignerBondsAtAuction(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dRetrieveSignerPubkey(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.RetrieveSignerPubkey()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallRetrieveSignerPubkey(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dInitialize(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	_factory, err := ethutil.AddressFromHex(c.Args()[0])
	if err != nil {
		return fmt.Errorf(
			"couldn't parse parameter _factory, a address, from passed value %v",
			c.Args()[0],
		)
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.Initialize(
			_factory,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallInitialize(
			_factory,
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifyRedemptionSignatureTimedOut(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyRedemptionSignatureTimedOut()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyRedemptionSignatureTimedOut(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dNotifyCourtesyCall(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.NotifyCourtesyCall()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallNotifyCourtesyCall(
			cmd.BlockFlagValue.Uint,
		)
		if err != nil {
			return err
		}

		cmd.PrintOutput(nil)
	}

	return nil
}

func dWithdrawFunds(c *cli.Context) error {
	contract, err := initializeDeposit(c)
	if err != nil {
		return err
	}

	var (
		transaction *types.Transaction
	)

	if c.Bool(cmd.SubmitFlag) {
		// Do a regular submission. Take payable into account.
		transaction, err = contract.WithdrawFunds()
		if err != nil {
			return err
		}

		cmd.PrintOutput(transaction.Hash)
	} else {
		// Do a call.
		err = contract.CallWithdrawFunds(
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

func initializeDeposit(c *cli.Context) (*contract.Deposit, error) {
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

	address := common.HexToAddress(config.ContractAddresses["Deposit"])

	return contract.NewDeposit(
		address,
		key,
		client,
		ethutil.NewNonceManager(key.Address, client),
		miningWaiter,
		&sync.Mutex{},
	)
}
