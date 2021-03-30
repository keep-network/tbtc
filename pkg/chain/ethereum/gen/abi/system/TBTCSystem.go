// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TBTCSystemABI is the input ABI used to generate the binding from.
const TBTCSystemABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_priceFeed\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_relay\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_allowNewDeposits\",\"type\":\"bool\"}],\"name\":\"AllowNewDepositsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_initialCollateralizedPercent\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_undercollateralizedThresholdPercent\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_severelyUndercollateralizedThresholdPercent\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"CollateralizationThresholdsUpdateStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_initialCollateralizedPercent\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_undercollateralizedThresholdPercent\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_severelyUndercollateralizedThresholdPercent\",\"type\":\"uint16\"}],\"name\":\"CollateralizationThresholdsUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"CourtesyCalled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_keepAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Created\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_priceFeed\",\"type\":\"address\"}],\"name\":\"EthBtcPriceFeedAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_priceFeed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"EthBtcPriceFeedAdditionStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"ExitedCourtesyCall\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"FraudDuringSetup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Funded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_abortOutputScript\",\"type\":\"bytes\"}],\"name\":\"FunderAbortRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"GotRedemptionSignature\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_keepStakedFactory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_fullyBackedFactory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_factorySelector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"KeepFactoriesUpdateStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_keepStakedFactory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_fullyBackedFactory\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_factorySelector\",\"type\":\"address\"}],\"name\":\"KeepFactoriesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Liquidated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64[]\",\"name\":\"_lotSizes\",\"type\":\"uint64[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"LotSizesUpdateStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64[]\",\"name\":\"_lotSizes\",\"type\":\"uint64[]\"}],\"name\":\"LotSizesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_utxoValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_requestedFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_outpoint\",\"type\":\"bytes\"}],\"name\":\"RedemptionRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyX\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyY\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"RegisteredPubkey\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"SetupFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_signerFeeDivisor\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"SignerFeeDivisorUpdateStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"_signerFeeDivisor\",\"type\":\"uint16\"}],\"name\":\"SignerFeeDivisorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_depositContractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"_wasFraud\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"StartedLiquidation\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_caller\",\"type\":\"address\"}],\"name\":\"approvedToLog\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_initialCollateralizedPercent\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_undercollateralizedThresholdPercent\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_severelyUndercollateralizedThresholdPercent\",\"type\":\"uint16\"}],\"name\":\"beginCollateralizationThresholdsUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIMedianizer\",\"name\":\"_ethBtcPriceFeed\",\"type\":\"address\"}],\"name\":\"beginEthBtcPriceFeedAddition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_keepStakedFactory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_fullyBackedFactory\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_factorySelector\",\"type\":\"address\"}],\"name\":\"beginKeepFactoriesUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"_lotSizes\",\"type\":\"uint64[]\"}],\"name\":\"beginLotSizesUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_signerFeeDivisor\",\"type\":\"uint16\"}],\"name\":\"beginSignerFeeDivisorUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"emergencyPauseNewDeposits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fetchBitcoinPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fetchRelayCurrentDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"fetchRelayPreviousDifficulty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeCollateralizationThresholdsUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeEthBtcPriceFeedAddition\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeKeepFactoriesUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeLotSizesUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeSignerFeeDivisorUpdate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllowNewDeposits\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAllowedLotSizes\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"\",\"type\":\"uint64[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGovernanceTimeDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getInitialCollateralizedPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getKeepFactoriesUpgradeabilityPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMaximumLotSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinimumLotSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNewDepositFeeEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPriceFeedGovernanceTimeDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingCollateralizationThresholdsUpdateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingEthBtcPriceFeedAdditionTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingKeepFactoriesUpdateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingKeepFactoriesUpgradeabilityTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingLotSizesUpdateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingPauseTerm\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRemainingSignerFeeDivisorUpdateTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSeverelyUndercollateralizedThresholdPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSignerFeeDivisor\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getUndercollateralizedThresholdPercent\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"contractIBondedECDSAKeepFactory\",\"name\":\"_defaultKeepFactory\",\"type\":\"address\"},{\"internalType\":\"contractDepositFactory\",\"name\":\"_depositFactory\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_masterDepositAddress\",\"type\":\"address\"},{\"internalType\":\"contractTBTCToken\",\"name\":\"_tbtcToken\",\"type\":\"address\"},{\"internalType\":\"contractTBTCDepositToken\",\"name\":\"_tbtcDepositToken\",\"type\":\"address\"},{\"internalType\":\"contractFeeRebateToken\",\"name\":\"_feeRebateToken\",\"type\":\"address\"},{\"internalType\":\"contractVendingMachine\",\"name\":\"_vendingMachine\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_keepThreshold\",\"type\":\"uint16\"},{\"internalType\":\"uint16\",\"name\":\"_keepSize\",\"type\":\"uint16\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_requestedLotSizeSatoshis\",\"type\":\"uint64\"}],\"name\":\"isAllowedLotSize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keepSize\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keepThreshold\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logCourtesyCalled\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_keepAddress\",\"type\":\"address\"}],\"name\":\"logCreated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logExitedCourtesyCall\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logFraudDuringSetup\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"}],\"name\":\"logFunded\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_abortOutputScript\",\"type\":\"bytes\"}],\"name\":\"logFunderRequestedAbort\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"name\":\"logGotRedemptionSignature\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logLiquidated\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_txid\",\"type\":\"bytes32\"}],\"name\":\"logRedeemed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_requester\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_digest\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_utxoValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_redeemerOutputScript\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"_requestedFee\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_outpoint\",\"type\":\"bytes\"}],\"name\":\"logRedemptionRequested\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyX\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_signingGroupPubkeyY\",\"type\":\"bytes32\"}],\"name\":\"logRegisteredPubkey\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"logSetupFailed\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_wasFraud\",\"type\":\"bool\"}],\"name\":\"logStartedLiquidation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"priceFeed\",\"outputs\":[{\"internalType\":\"contractISatWeiPriceFeed\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"refreshMinimumBondableValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"relay\",\"outputs\":[{\"internalType\":\"contractIRelay\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"_requestedLotSizeSatoshis\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"_maxSecuredLifetime\",\"type\":\"uint256\"}],\"name\":\"requestNewKeep\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resumeNewDeposits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TBTCSystem is an auto generated Go binding around an Ethereum contract.
type TBTCSystem struct {
	TBTCSystemCaller     // Read-only binding to the contract
	TBTCSystemTransactor // Write-only binding to the contract
	TBTCSystemFilterer   // Log filterer for contract events
}

// TBTCSystemCaller is an auto generated read-only Go binding around an Ethereum contract.
type TBTCSystemCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TBTCSystemTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TBTCSystemTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TBTCSystemFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TBTCSystemFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TBTCSystemSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TBTCSystemSession struct {
	Contract     *TBTCSystem       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TBTCSystemCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TBTCSystemCallerSession struct {
	Contract *TBTCSystemCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// TBTCSystemTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TBTCSystemTransactorSession struct {
	Contract     *TBTCSystemTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// TBTCSystemRaw is an auto generated low-level Go binding around an Ethereum contract.
type TBTCSystemRaw struct {
	Contract *TBTCSystem // Generic contract binding to access the raw methods on
}

// TBTCSystemCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TBTCSystemCallerRaw struct {
	Contract *TBTCSystemCaller // Generic read-only contract binding to access the raw methods on
}

// TBTCSystemTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TBTCSystemTransactorRaw struct {
	Contract *TBTCSystemTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTBTCSystem creates a new instance of TBTCSystem, bound to a specific deployed contract.
func NewTBTCSystem(address common.Address, backend bind.ContractBackend) (*TBTCSystem, error) {
	contract, err := bindTBTCSystem(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TBTCSystem{TBTCSystemCaller: TBTCSystemCaller{contract: contract}, TBTCSystemTransactor: TBTCSystemTransactor{contract: contract}, TBTCSystemFilterer: TBTCSystemFilterer{contract: contract}}, nil
}

// NewTBTCSystemCaller creates a new read-only instance of TBTCSystem, bound to a specific deployed contract.
func NewTBTCSystemCaller(address common.Address, caller bind.ContractCaller) (*TBTCSystemCaller, error) {
	contract, err := bindTBTCSystem(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemCaller{contract: contract}, nil
}

// NewTBTCSystemTransactor creates a new write-only instance of TBTCSystem, bound to a specific deployed contract.
func NewTBTCSystemTransactor(address common.Address, transactor bind.ContractTransactor) (*TBTCSystemTransactor, error) {
	contract, err := bindTBTCSystem(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemTransactor{contract: contract}, nil
}

// NewTBTCSystemFilterer creates a new log filterer instance of TBTCSystem, bound to a specific deployed contract.
func NewTBTCSystemFilterer(address common.Address, filterer bind.ContractFilterer) (*TBTCSystemFilterer, error) {
	contract, err := bindTBTCSystem(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemFilterer{contract: contract}, nil
}

// bindTBTCSystem binds a generic wrapper to an already deployed contract.
func bindTBTCSystem(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TBTCSystemABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TBTCSystem *TBTCSystemRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TBTCSystem.Contract.TBTCSystemCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TBTCSystem *TBTCSystemRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.Contract.TBTCSystemTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TBTCSystem *TBTCSystemRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TBTCSystem.Contract.TBTCSystemTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TBTCSystem *TBTCSystemCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TBTCSystem.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TBTCSystem *TBTCSystemTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TBTCSystem *TBTCSystemTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TBTCSystem.Contract.contract.Transact(opts, method, params...)
}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) view returns(bool)
func (_TBTCSystem *TBTCSystemCaller) ApprovedToLog(opts *bind.CallOpts, _caller common.Address) (bool, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "approvedToLog", _caller)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) view returns(bool)
func (_TBTCSystem *TBTCSystemSession) ApprovedToLog(_caller common.Address) (bool, error) {
	return _TBTCSystem.Contract.ApprovedToLog(&_TBTCSystem.CallOpts, _caller)
}

// ApprovedToLog is a free data retrieval call binding the contract method 0x9ffb3862.
//
// Solidity: function approvedToLog(address _caller) view returns(bool)
func (_TBTCSystem *TBTCSystemCallerSession) ApprovedToLog(_caller common.Address) (bool, error) {
	return _TBTCSystem.Contract.ApprovedToLog(&_TBTCSystem.CallOpts, _caller)
}

// FetchBitcoinPrice is a free data retrieval call binding the contract method 0xa6c1691c.
//
// Solidity: function fetchBitcoinPrice() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) FetchBitcoinPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "fetchBitcoinPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FetchBitcoinPrice is a free data retrieval call binding the contract method 0xa6c1691c.
//
// Solidity: function fetchBitcoinPrice() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) FetchBitcoinPrice() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchBitcoinPrice(&_TBTCSystem.CallOpts)
}

// FetchBitcoinPrice is a free data retrieval call binding the contract method 0xa6c1691c.
//
// Solidity: function fetchBitcoinPrice() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) FetchBitcoinPrice() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchBitcoinPrice(&_TBTCSystem.CallOpts)
}

// FetchRelayCurrentDifficulty is a free data retrieval call binding the contract method 0xdab70cb1.
//
// Solidity: function fetchRelayCurrentDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) FetchRelayCurrentDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "fetchRelayCurrentDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FetchRelayCurrentDifficulty is a free data retrieval call binding the contract method 0xdab70cb1.
//
// Solidity: function fetchRelayCurrentDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) FetchRelayCurrentDifficulty() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchRelayCurrentDifficulty(&_TBTCSystem.CallOpts)
}

// FetchRelayCurrentDifficulty is a free data retrieval call binding the contract method 0xdab70cb1.
//
// Solidity: function fetchRelayCurrentDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) FetchRelayCurrentDifficulty() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchRelayCurrentDifficulty(&_TBTCSystem.CallOpts)
}

// FetchRelayPreviousDifficulty is a free data retrieval call binding the contract method 0x402b783d.
//
// Solidity: function fetchRelayPreviousDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) FetchRelayPreviousDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "fetchRelayPreviousDifficulty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FetchRelayPreviousDifficulty is a free data retrieval call binding the contract method 0x402b783d.
//
// Solidity: function fetchRelayPreviousDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) FetchRelayPreviousDifficulty() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchRelayPreviousDifficulty(&_TBTCSystem.CallOpts)
}

// FetchRelayPreviousDifficulty is a free data retrieval call binding the contract method 0x402b783d.
//
// Solidity: function fetchRelayPreviousDifficulty() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) FetchRelayPreviousDifficulty() (*big.Int, error) {
	return _TBTCSystem.Contract.FetchRelayPreviousDifficulty(&_TBTCSystem.CallOpts)
}

// GetAllowNewDeposits is a free data retrieval call binding the contract method 0x0d7eb1c4.
//
// Solidity: function getAllowNewDeposits() view returns(bool)
func (_TBTCSystem *TBTCSystemCaller) GetAllowNewDeposits(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getAllowNewDeposits")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetAllowNewDeposits is a free data retrieval call binding the contract method 0x0d7eb1c4.
//
// Solidity: function getAllowNewDeposits() view returns(bool)
func (_TBTCSystem *TBTCSystemSession) GetAllowNewDeposits() (bool, error) {
	return _TBTCSystem.Contract.GetAllowNewDeposits(&_TBTCSystem.CallOpts)
}

// GetAllowNewDeposits is a free data retrieval call binding the contract method 0x0d7eb1c4.
//
// Solidity: function getAllowNewDeposits() view returns(bool)
func (_TBTCSystem *TBTCSystemCallerSession) GetAllowNewDeposits() (bool, error) {
	return _TBTCSystem.Contract.GetAllowNewDeposits(&_TBTCSystem.CallOpts)
}

// GetAllowedLotSizes is a free data retrieval call binding the contract method 0x086c9edd.
//
// Solidity: function getAllowedLotSizes() view returns(uint64[])
func (_TBTCSystem *TBTCSystemCaller) GetAllowedLotSizes(opts *bind.CallOpts) ([]uint64, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getAllowedLotSizes")

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

// GetAllowedLotSizes is a free data retrieval call binding the contract method 0x086c9edd.
//
// Solidity: function getAllowedLotSizes() view returns(uint64[])
func (_TBTCSystem *TBTCSystemSession) GetAllowedLotSizes() ([]uint64, error) {
	return _TBTCSystem.Contract.GetAllowedLotSizes(&_TBTCSystem.CallOpts)
}

// GetAllowedLotSizes is a free data retrieval call binding the contract method 0x086c9edd.
//
// Solidity: function getAllowedLotSizes() view returns(uint64[])
func (_TBTCSystem *TBTCSystemCallerSession) GetAllowedLotSizes() ([]uint64, error) {
	return _TBTCSystem.Contract.GetAllowedLotSizes(&_TBTCSystem.CallOpts)
}

// GetGovernanceTimeDelay is a free data retrieval call binding the contract method 0xf2e72347.
//
// Solidity: function getGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetGovernanceTimeDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getGovernanceTimeDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGovernanceTimeDelay is a free data retrieval call binding the contract method 0xf2e72347.
//
// Solidity: function getGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetGovernanceTimeDelay() (*big.Int, error) {
	return _TBTCSystem.Contract.GetGovernanceTimeDelay(&_TBTCSystem.CallOpts)
}

// GetGovernanceTimeDelay is a free data retrieval call binding the contract method 0xf2e72347.
//
// Solidity: function getGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetGovernanceTimeDelay() (*big.Int, error) {
	return _TBTCSystem.Contract.GetGovernanceTimeDelay(&_TBTCSystem.CallOpts)
}

// GetInitialCollateralizedPercent is a free data retrieval call binding the contract method 0x987ecea7.
//
// Solidity: function getInitialCollateralizedPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) GetInitialCollateralizedPercent(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getInitialCollateralizedPercent")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetInitialCollateralizedPercent is a free data retrieval call binding the contract method 0x987ecea7.
//
// Solidity: function getInitialCollateralizedPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) GetInitialCollateralizedPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetInitialCollateralizedPercent(&_TBTCSystem.CallOpts)
}

// GetInitialCollateralizedPercent is a free data retrieval call binding the contract method 0x987ecea7.
//
// Solidity: function getInitialCollateralizedPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) GetInitialCollateralizedPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetInitialCollateralizedPercent(&_TBTCSystem.CallOpts)
}

// GetKeepFactoriesUpgradeabilityPeriod is a free data retrieval call binding the contract method 0xe5a6d77d.
//
// Solidity: function getKeepFactoriesUpgradeabilityPeriod() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetKeepFactoriesUpgradeabilityPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getKeepFactoriesUpgradeabilityPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetKeepFactoriesUpgradeabilityPeriod is a free data retrieval call binding the contract method 0xe5a6d77d.
//
// Solidity: function getKeepFactoriesUpgradeabilityPeriod() pure returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetKeepFactoriesUpgradeabilityPeriod() (*big.Int, error) {
	return _TBTCSystem.Contract.GetKeepFactoriesUpgradeabilityPeriod(&_TBTCSystem.CallOpts)
}

// GetKeepFactoriesUpgradeabilityPeriod is a free data retrieval call binding the contract method 0xe5a6d77d.
//
// Solidity: function getKeepFactoriesUpgradeabilityPeriod() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetKeepFactoriesUpgradeabilityPeriod() (*big.Int, error) {
	return _TBTCSystem.Contract.GetKeepFactoriesUpgradeabilityPeriod(&_TBTCSystem.CallOpts)
}

// GetMaximumLotSize is a free data retrieval call binding the contract method 0x2753d84b.
//
// Solidity: function getMaximumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetMaximumLotSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getMaximumLotSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaximumLotSize is a free data retrieval call binding the contract method 0x2753d84b.
//
// Solidity: function getMaximumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetMaximumLotSize() (*big.Int, error) {
	return _TBTCSystem.Contract.GetMaximumLotSize(&_TBTCSystem.CallOpts)
}

// GetMaximumLotSize is a free data retrieval call binding the contract method 0x2753d84b.
//
// Solidity: function getMaximumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetMaximumLotSize() (*big.Int, error) {
	return _TBTCSystem.Contract.GetMaximumLotSize(&_TBTCSystem.CallOpts)
}

// GetMinimumLotSize is a free data retrieval call binding the contract method 0x34d534a9.
//
// Solidity: function getMinimumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetMinimumLotSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getMinimumLotSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumLotSize is a free data retrieval call binding the contract method 0x34d534a9.
//
// Solidity: function getMinimumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetMinimumLotSize() (*big.Int, error) {
	return _TBTCSystem.Contract.GetMinimumLotSize(&_TBTCSystem.CallOpts)
}

// GetMinimumLotSize is a free data retrieval call binding the contract method 0x34d534a9.
//
// Solidity: function getMinimumLotSize() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetMinimumLotSize() (*big.Int, error) {
	return _TBTCSystem.Contract.GetMinimumLotSize(&_TBTCSystem.CallOpts)
}

// GetNewDepositFeeEstimate is a free data retrieval call binding the contract method 0x2d00f1ee.
//
// Solidity: function getNewDepositFeeEstimate() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetNewDepositFeeEstimate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getNewDepositFeeEstimate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNewDepositFeeEstimate is a free data retrieval call binding the contract method 0x2d00f1ee.
//
// Solidity: function getNewDepositFeeEstimate() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetNewDepositFeeEstimate() (*big.Int, error) {
	return _TBTCSystem.Contract.GetNewDepositFeeEstimate(&_TBTCSystem.CallOpts)
}

// GetNewDepositFeeEstimate is a free data retrieval call binding the contract method 0x2d00f1ee.
//
// Solidity: function getNewDepositFeeEstimate() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetNewDepositFeeEstimate() (*big.Int, error) {
	return _TBTCSystem.Contract.GetNewDepositFeeEstimate(&_TBTCSystem.CallOpts)
}

// GetPriceFeedGovernanceTimeDelay is a free data retrieval call binding the contract method 0xae7f4a5f.
//
// Solidity: function getPriceFeedGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetPriceFeedGovernanceTimeDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getPriceFeedGovernanceTimeDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPriceFeedGovernanceTimeDelay is a free data retrieval call binding the contract method 0xae7f4a5f.
//
// Solidity: function getPriceFeedGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetPriceFeedGovernanceTimeDelay() (*big.Int, error) {
	return _TBTCSystem.Contract.GetPriceFeedGovernanceTimeDelay(&_TBTCSystem.CallOpts)
}

// GetPriceFeedGovernanceTimeDelay is a free data retrieval call binding the contract method 0xae7f4a5f.
//
// Solidity: function getPriceFeedGovernanceTimeDelay() pure returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetPriceFeedGovernanceTimeDelay() (*big.Int, error) {
	return _TBTCSystem.Contract.GetPriceFeedGovernanceTimeDelay(&_TBTCSystem.CallOpts)
}

// GetRemainingCollateralizationThresholdsUpdateTime is a free data retrieval call binding the contract method 0xc074d550.
//
// Solidity: function getRemainingCollateralizationThresholdsUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingCollateralizationThresholdsUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingCollateralizationThresholdsUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingCollateralizationThresholdsUpdateTime is a free data retrieval call binding the contract method 0xc074d550.
//
// Solidity: function getRemainingCollateralizationThresholdsUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingCollateralizationThresholdsUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingCollateralizationThresholdsUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingCollateralizationThresholdsUpdateTime is a free data retrieval call binding the contract method 0xc074d550.
//
// Solidity: function getRemainingCollateralizationThresholdsUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingCollateralizationThresholdsUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingCollateralizationThresholdsUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingEthBtcPriceFeedAdditionTime is a free data retrieval call binding the contract method 0x0af488f9.
//
// Solidity: function getRemainingEthBtcPriceFeedAdditionTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingEthBtcPriceFeedAdditionTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingEthBtcPriceFeedAdditionTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingEthBtcPriceFeedAdditionTime is a free data retrieval call binding the contract method 0x0af488f9.
//
// Solidity: function getRemainingEthBtcPriceFeedAdditionTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingEthBtcPriceFeedAdditionTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingEthBtcPriceFeedAdditionTime(&_TBTCSystem.CallOpts)
}

// GetRemainingEthBtcPriceFeedAdditionTime is a free data retrieval call binding the contract method 0x0af488f9.
//
// Solidity: function getRemainingEthBtcPriceFeedAdditionTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingEthBtcPriceFeedAdditionTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingEthBtcPriceFeedAdditionTime(&_TBTCSystem.CallOpts)
}

// GetRemainingKeepFactoriesUpdateTime is a free data retrieval call binding the contract method 0x57535088.
//
// Solidity: function getRemainingKeepFactoriesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingKeepFactoriesUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingKeepFactoriesUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingKeepFactoriesUpdateTime is a free data retrieval call binding the contract method 0x57535088.
//
// Solidity: function getRemainingKeepFactoriesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingKeepFactoriesUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingKeepFactoriesUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingKeepFactoriesUpdateTime is a free data retrieval call binding the contract method 0x57535088.
//
// Solidity: function getRemainingKeepFactoriesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingKeepFactoriesUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingKeepFactoriesUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingKeepFactoriesUpgradeabilityTime is a free data retrieval call binding the contract method 0xb196b5a3.
//
// Solidity: function getRemainingKeepFactoriesUpgradeabilityTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingKeepFactoriesUpgradeabilityTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingKeepFactoriesUpgradeabilityTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingKeepFactoriesUpgradeabilityTime is a free data retrieval call binding the contract method 0xb196b5a3.
//
// Solidity: function getRemainingKeepFactoriesUpgradeabilityTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingKeepFactoriesUpgradeabilityTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingKeepFactoriesUpgradeabilityTime(&_TBTCSystem.CallOpts)
}

// GetRemainingKeepFactoriesUpgradeabilityTime is a free data retrieval call binding the contract method 0xb196b5a3.
//
// Solidity: function getRemainingKeepFactoriesUpgradeabilityTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingKeepFactoriesUpgradeabilityTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingKeepFactoriesUpgradeabilityTime(&_TBTCSystem.CallOpts)
}

// GetRemainingLotSizesUpdateTime is a free data retrieval call binding the contract method 0x3ee850bc.
//
// Solidity: function getRemainingLotSizesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingLotSizesUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingLotSizesUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingLotSizesUpdateTime is a free data retrieval call binding the contract method 0x3ee850bc.
//
// Solidity: function getRemainingLotSizesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingLotSizesUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingLotSizesUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingLotSizesUpdateTime is a free data retrieval call binding the contract method 0x3ee850bc.
//
// Solidity: function getRemainingLotSizesUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingLotSizesUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingLotSizesUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingPauseTerm is a free data retrieval call binding the contract method 0x013b0f30.
//
// Solidity: function getRemainingPauseTerm() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingPauseTerm(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingPauseTerm")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingPauseTerm is a free data retrieval call binding the contract method 0x013b0f30.
//
// Solidity: function getRemainingPauseTerm() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingPauseTerm() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingPauseTerm(&_TBTCSystem.CallOpts)
}

// GetRemainingPauseTerm is a free data retrieval call binding the contract method 0x013b0f30.
//
// Solidity: function getRemainingPauseTerm() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingPauseTerm() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingPauseTerm(&_TBTCSystem.CallOpts)
}

// GetRemainingSignerFeeDivisorUpdateTime is a free data retrieval call binding the contract method 0xb792a38e.
//
// Solidity: function getRemainingSignerFeeDivisorUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCaller) GetRemainingSignerFeeDivisorUpdateTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getRemainingSignerFeeDivisorUpdateTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRemainingSignerFeeDivisorUpdateTime is a free data retrieval call binding the contract method 0xb792a38e.
//
// Solidity: function getRemainingSignerFeeDivisorUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemSession) GetRemainingSignerFeeDivisorUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingSignerFeeDivisorUpdateTime(&_TBTCSystem.CallOpts)
}

// GetRemainingSignerFeeDivisorUpdateTime is a free data retrieval call binding the contract method 0xb792a38e.
//
// Solidity: function getRemainingSignerFeeDivisorUpdateTime() view returns(uint256)
func (_TBTCSystem *TBTCSystemCallerSession) GetRemainingSignerFeeDivisorUpdateTime() (*big.Int, error) {
	return _TBTCSystem.Contract.GetRemainingSignerFeeDivisorUpdateTime(&_TBTCSystem.CallOpts)
}

// GetSeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x6f4fef62.
//
// Solidity: function getSeverelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) GetSeverelyUndercollateralizedThresholdPercent(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getSeverelyUndercollateralizedThresholdPercent")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetSeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x6f4fef62.
//
// Solidity: function getSeverelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) GetSeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetSeverelyUndercollateralizedThresholdPercent(&_TBTCSystem.CallOpts)
}

// GetSeverelyUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0x6f4fef62.
//
// Solidity: function getSeverelyUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) GetSeverelyUndercollateralizedThresholdPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetSeverelyUndercollateralizedThresholdPercent(&_TBTCSystem.CallOpts)
}

// GetSignerFeeDivisor is a free data retrieval call binding the contract method 0x60e98d59.
//
// Solidity: function getSignerFeeDivisor() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) GetSignerFeeDivisor(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getSignerFeeDivisor")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetSignerFeeDivisor is a free data retrieval call binding the contract method 0x60e98d59.
//
// Solidity: function getSignerFeeDivisor() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) GetSignerFeeDivisor() (uint16, error) {
	return _TBTCSystem.Contract.GetSignerFeeDivisor(&_TBTCSystem.CallOpts)
}

// GetSignerFeeDivisor is a free data retrieval call binding the contract method 0x60e98d59.
//
// Solidity: function getSignerFeeDivisor() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) GetSignerFeeDivisor() (uint16, error) {
	return _TBTCSystem.Contract.GetSignerFeeDivisor(&_TBTCSystem.CallOpts)
}

// GetUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0xa2cd75da.
//
// Solidity: function getUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) GetUndercollateralizedThresholdPercent(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "getUndercollateralizedThresholdPercent")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// GetUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0xa2cd75da.
//
// Solidity: function getUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) GetUndercollateralizedThresholdPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetUndercollateralizedThresholdPercent(&_TBTCSystem.CallOpts)
}

// GetUndercollateralizedThresholdPercent is a free data retrieval call binding the contract method 0xa2cd75da.
//
// Solidity: function getUndercollateralizedThresholdPercent() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) GetUndercollateralizedThresholdPercent() (uint16, error) {
	return _TBTCSystem.Contract.GetUndercollateralizedThresholdPercent(&_TBTCSystem.CallOpts)
}

// IsAllowedLotSize is a free data retrieval call binding the contract method 0xa28b79f1.
//
// Solidity: function isAllowedLotSize(uint64 _requestedLotSizeSatoshis) view returns(bool)
func (_TBTCSystem *TBTCSystemCaller) IsAllowedLotSize(opts *bind.CallOpts, _requestedLotSizeSatoshis uint64) (bool, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "isAllowedLotSize", _requestedLotSizeSatoshis)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAllowedLotSize is a free data retrieval call binding the contract method 0xa28b79f1.
//
// Solidity: function isAllowedLotSize(uint64 _requestedLotSizeSatoshis) view returns(bool)
func (_TBTCSystem *TBTCSystemSession) IsAllowedLotSize(_requestedLotSizeSatoshis uint64) (bool, error) {
	return _TBTCSystem.Contract.IsAllowedLotSize(&_TBTCSystem.CallOpts, _requestedLotSizeSatoshis)
}

// IsAllowedLotSize is a free data retrieval call binding the contract method 0xa28b79f1.
//
// Solidity: function isAllowedLotSize(uint64 _requestedLotSizeSatoshis) view returns(bool)
func (_TBTCSystem *TBTCSystemCallerSession) IsAllowedLotSize(_requestedLotSizeSatoshis uint64) (bool, error) {
	return _TBTCSystem.Contract.IsAllowedLotSize(&_TBTCSystem.CallOpts, _requestedLotSizeSatoshis)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TBTCSystem *TBTCSystemCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TBTCSystem *TBTCSystemSession) IsOwner() (bool, error) {
	return _TBTCSystem.Contract.IsOwner(&_TBTCSystem.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TBTCSystem *TBTCSystemCallerSession) IsOwner() (bool, error) {
	return _TBTCSystem.Contract.IsOwner(&_TBTCSystem.CallOpts)
}

// KeepSize is a free data retrieval call binding the contract method 0x64bdb667.
//
// Solidity: function keepSize() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) KeepSize(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "keepSize")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// KeepSize is a free data retrieval call binding the contract method 0x64bdb667.
//
// Solidity: function keepSize() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) KeepSize() (uint16, error) {
	return _TBTCSystem.Contract.KeepSize(&_TBTCSystem.CallOpts)
}

// KeepSize is a free data retrieval call binding the contract method 0x64bdb667.
//
// Solidity: function keepSize() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) KeepSize() (uint16, error) {
	return _TBTCSystem.Contract.KeepSize(&_TBTCSystem.CallOpts)
}

// KeepThreshold is a free data retrieval call binding the contract method 0xe5426d2e.
//
// Solidity: function keepThreshold() view returns(uint16)
func (_TBTCSystem *TBTCSystemCaller) KeepThreshold(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "keepThreshold")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// KeepThreshold is a free data retrieval call binding the contract method 0xe5426d2e.
//
// Solidity: function keepThreshold() view returns(uint16)
func (_TBTCSystem *TBTCSystemSession) KeepThreshold() (uint16, error) {
	return _TBTCSystem.Contract.KeepThreshold(&_TBTCSystem.CallOpts)
}

// KeepThreshold is a free data retrieval call binding the contract method 0xe5426d2e.
//
// Solidity: function keepThreshold() view returns(uint16)
func (_TBTCSystem *TBTCSystemCallerSession) KeepThreshold() (uint16, error) {
	return _TBTCSystem.Contract.KeepThreshold(&_TBTCSystem.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TBTCSystem *TBTCSystemCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TBTCSystem *TBTCSystemSession) Owner() (common.Address, error) {
	return _TBTCSystem.Contract.Owner(&_TBTCSystem.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TBTCSystem *TBTCSystemCallerSession) Owner() (common.Address, error) {
	return _TBTCSystem.Contract.Owner(&_TBTCSystem.CallOpts)
}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_TBTCSystem *TBTCSystemCaller) PriceFeed(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "priceFeed")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_TBTCSystem *TBTCSystemSession) PriceFeed() (common.Address, error) {
	return _TBTCSystem.Contract.PriceFeed(&_TBTCSystem.CallOpts)
}

// PriceFeed is a free data retrieval call binding the contract method 0x741bef1a.
//
// Solidity: function priceFeed() view returns(address)
func (_TBTCSystem *TBTCSystemCallerSession) PriceFeed() (common.Address, error) {
	return _TBTCSystem.Contract.PriceFeed(&_TBTCSystem.CallOpts)
}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_TBTCSystem *TBTCSystemCaller) Relay(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TBTCSystem.contract.Call(opts, &out, "relay")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_TBTCSystem *TBTCSystemSession) Relay() (common.Address, error) {
	return _TBTCSystem.Contract.Relay(&_TBTCSystem.CallOpts)
}

// Relay is a free data retrieval call binding the contract method 0xb59589d1.
//
// Solidity: function relay() view returns(address)
func (_TBTCSystem *TBTCSystemCallerSession) Relay() (common.Address, error) {
	return _TBTCSystem.Contract.Relay(&_TBTCSystem.CallOpts)
}

// BeginCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xdde3fdd2.
//
// Solidity: function beginCollateralizationThresholdsUpdate(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent) returns()
func (_TBTCSystem *TBTCSystemTransactor) BeginCollateralizationThresholdsUpdate(opts *bind.TransactOpts, _initialCollateralizedPercent uint16, _undercollateralizedThresholdPercent uint16, _severelyUndercollateralizedThresholdPercent uint16) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "beginCollateralizationThresholdsUpdate", _initialCollateralizedPercent, _undercollateralizedThresholdPercent, _severelyUndercollateralizedThresholdPercent)
}

// BeginCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xdde3fdd2.
//
// Solidity: function beginCollateralizationThresholdsUpdate(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent) returns()
func (_TBTCSystem *TBTCSystemSession) BeginCollateralizationThresholdsUpdate(_initialCollateralizedPercent uint16, _undercollateralizedThresholdPercent uint16, _severelyUndercollateralizedThresholdPercent uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginCollateralizationThresholdsUpdate(&_TBTCSystem.TransactOpts, _initialCollateralizedPercent, _undercollateralizedThresholdPercent, _severelyUndercollateralizedThresholdPercent)
}

// BeginCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xdde3fdd2.
//
// Solidity: function beginCollateralizationThresholdsUpdate(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) BeginCollateralizationThresholdsUpdate(_initialCollateralizedPercent uint16, _undercollateralizedThresholdPercent uint16, _severelyUndercollateralizedThresholdPercent uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginCollateralizationThresholdsUpdate(&_TBTCSystem.TransactOpts, _initialCollateralizedPercent, _undercollateralizedThresholdPercent, _severelyUndercollateralizedThresholdPercent)
}

// BeginEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0x07a3d659.
//
// Solidity: function beginEthBtcPriceFeedAddition(address _ethBtcPriceFeed) returns()
func (_TBTCSystem *TBTCSystemTransactor) BeginEthBtcPriceFeedAddition(opts *bind.TransactOpts, _ethBtcPriceFeed common.Address) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "beginEthBtcPriceFeedAddition", _ethBtcPriceFeed)
}

// BeginEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0x07a3d659.
//
// Solidity: function beginEthBtcPriceFeedAddition(address _ethBtcPriceFeed) returns()
func (_TBTCSystem *TBTCSystemSession) BeginEthBtcPriceFeedAddition(_ethBtcPriceFeed common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginEthBtcPriceFeedAddition(&_TBTCSystem.TransactOpts, _ethBtcPriceFeed)
}

// BeginEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0x07a3d659.
//
// Solidity: function beginEthBtcPriceFeedAddition(address _ethBtcPriceFeed) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) BeginEthBtcPriceFeedAddition(_ethBtcPriceFeed common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginEthBtcPriceFeedAddition(&_TBTCSystem.TransactOpts, _ethBtcPriceFeed)
}

// BeginKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0xeae6191f.
//
// Solidity: function beginKeepFactoriesUpdate(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector) returns()
func (_TBTCSystem *TBTCSystemTransactor) BeginKeepFactoriesUpdate(opts *bind.TransactOpts, _keepStakedFactory common.Address, _fullyBackedFactory common.Address, _factorySelector common.Address) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "beginKeepFactoriesUpdate", _keepStakedFactory, _fullyBackedFactory, _factorySelector)
}

// BeginKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0xeae6191f.
//
// Solidity: function beginKeepFactoriesUpdate(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector) returns()
func (_TBTCSystem *TBTCSystemSession) BeginKeepFactoriesUpdate(_keepStakedFactory common.Address, _fullyBackedFactory common.Address, _factorySelector common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginKeepFactoriesUpdate(&_TBTCSystem.TransactOpts, _keepStakedFactory, _fullyBackedFactory, _factorySelector)
}

// BeginKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0xeae6191f.
//
// Solidity: function beginKeepFactoriesUpdate(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) BeginKeepFactoriesUpdate(_keepStakedFactory common.Address, _fullyBackedFactory common.Address, _factorySelector common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginKeepFactoriesUpdate(&_TBTCSystem.TransactOpts, _keepStakedFactory, _fullyBackedFactory, _factorySelector)
}

// BeginLotSizesUpdate is a paid mutator transaction binding the contract method 0x2b155e37.
//
// Solidity: function beginLotSizesUpdate(uint64[] _lotSizes) returns()
func (_TBTCSystem *TBTCSystemTransactor) BeginLotSizesUpdate(opts *bind.TransactOpts, _lotSizes []uint64) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "beginLotSizesUpdate", _lotSizes)
}

// BeginLotSizesUpdate is a paid mutator transaction binding the contract method 0x2b155e37.
//
// Solidity: function beginLotSizesUpdate(uint64[] _lotSizes) returns()
func (_TBTCSystem *TBTCSystemSession) BeginLotSizesUpdate(_lotSizes []uint64) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginLotSizesUpdate(&_TBTCSystem.TransactOpts, _lotSizes)
}

// BeginLotSizesUpdate is a paid mutator transaction binding the contract method 0x2b155e37.
//
// Solidity: function beginLotSizesUpdate(uint64[] _lotSizes) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) BeginLotSizesUpdate(_lotSizes []uint64) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginLotSizesUpdate(&_TBTCSystem.TransactOpts, _lotSizes)
}

// BeginSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x49b64730.
//
// Solidity: function beginSignerFeeDivisorUpdate(uint16 _signerFeeDivisor) returns()
func (_TBTCSystem *TBTCSystemTransactor) BeginSignerFeeDivisorUpdate(opts *bind.TransactOpts, _signerFeeDivisor uint16) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "beginSignerFeeDivisorUpdate", _signerFeeDivisor)
}

// BeginSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x49b64730.
//
// Solidity: function beginSignerFeeDivisorUpdate(uint16 _signerFeeDivisor) returns()
func (_TBTCSystem *TBTCSystemSession) BeginSignerFeeDivisorUpdate(_signerFeeDivisor uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginSignerFeeDivisorUpdate(&_TBTCSystem.TransactOpts, _signerFeeDivisor)
}

// BeginSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x49b64730.
//
// Solidity: function beginSignerFeeDivisorUpdate(uint16 _signerFeeDivisor) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) BeginSignerFeeDivisorUpdate(_signerFeeDivisor uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.BeginSignerFeeDivisorUpdate(&_TBTCSystem.TransactOpts, _signerFeeDivisor)
}

// EmergencyPauseNewDeposits is a paid mutator transaction binding the contract method 0x80f04b8c.
//
// Solidity: function emergencyPauseNewDeposits() returns()
func (_TBTCSystem *TBTCSystemTransactor) EmergencyPauseNewDeposits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "emergencyPauseNewDeposits")
}

// EmergencyPauseNewDeposits is a paid mutator transaction binding the contract method 0x80f04b8c.
//
// Solidity: function emergencyPauseNewDeposits() returns()
func (_TBTCSystem *TBTCSystemSession) EmergencyPauseNewDeposits() (*types.Transaction, error) {
	return _TBTCSystem.Contract.EmergencyPauseNewDeposits(&_TBTCSystem.TransactOpts)
}

// EmergencyPauseNewDeposits is a paid mutator transaction binding the contract method 0x80f04b8c.
//
// Solidity: function emergencyPauseNewDeposits() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) EmergencyPauseNewDeposits() (*types.Transaction, error) {
	return _TBTCSystem.Contract.EmergencyPauseNewDeposits(&_TBTCSystem.TransactOpts)
}

// FinalizeCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xde1e57d0.
//
// Solidity: function finalizeCollateralizationThresholdsUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactor) FinalizeCollateralizationThresholdsUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "finalizeCollateralizationThresholdsUpdate")
}

// FinalizeCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xde1e57d0.
//
// Solidity: function finalizeCollateralizationThresholdsUpdate() returns()
func (_TBTCSystem *TBTCSystemSession) FinalizeCollateralizationThresholdsUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeCollateralizationThresholdsUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeCollateralizationThresholdsUpdate is a paid mutator transaction binding the contract method 0xde1e57d0.
//
// Solidity: function finalizeCollateralizationThresholdsUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) FinalizeCollateralizationThresholdsUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeCollateralizationThresholdsUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0xadc3ef70.
//
// Solidity: function finalizeEthBtcPriceFeedAddition() returns()
func (_TBTCSystem *TBTCSystemTransactor) FinalizeEthBtcPriceFeedAddition(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "finalizeEthBtcPriceFeedAddition")
}

// FinalizeEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0xadc3ef70.
//
// Solidity: function finalizeEthBtcPriceFeedAddition() returns()
func (_TBTCSystem *TBTCSystemSession) FinalizeEthBtcPriceFeedAddition() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeEthBtcPriceFeedAddition(&_TBTCSystem.TransactOpts)
}

// FinalizeEthBtcPriceFeedAddition is a paid mutator transaction binding the contract method 0xadc3ef70.
//
// Solidity: function finalizeEthBtcPriceFeedAddition() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) FinalizeEthBtcPriceFeedAddition() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeEthBtcPriceFeedAddition(&_TBTCSystem.TransactOpts)
}

// FinalizeKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0x0d4455dc.
//
// Solidity: function finalizeKeepFactoriesUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactor) FinalizeKeepFactoriesUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "finalizeKeepFactoriesUpdate")
}

// FinalizeKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0x0d4455dc.
//
// Solidity: function finalizeKeepFactoriesUpdate() returns()
func (_TBTCSystem *TBTCSystemSession) FinalizeKeepFactoriesUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeKeepFactoriesUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeKeepFactoriesUpdate is a paid mutator transaction binding the contract method 0x0d4455dc.
//
// Solidity: function finalizeKeepFactoriesUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) FinalizeKeepFactoriesUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeKeepFactoriesUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeLotSizesUpdate is a paid mutator transaction binding the contract method 0xcd3a9490.
//
// Solidity: function finalizeLotSizesUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactor) FinalizeLotSizesUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "finalizeLotSizesUpdate")
}

// FinalizeLotSizesUpdate is a paid mutator transaction binding the contract method 0xcd3a9490.
//
// Solidity: function finalizeLotSizesUpdate() returns()
func (_TBTCSystem *TBTCSystemSession) FinalizeLotSizesUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeLotSizesUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeLotSizesUpdate is a paid mutator transaction binding the contract method 0xcd3a9490.
//
// Solidity: function finalizeLotSizesUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) FinalizeLotSizesUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeLotSizesUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x0ce0e700.
//
// Solidity: function finalizeSignerFeeDivisorUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactor) FinalizeSignerFeeDivisorUpdate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "finalizeSignerFeeDivisorUpdate")
}

// FinalizeSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x0ce0e700.
//
// Solidity: function finalizeSignerFeeDivisorUpdate() returns()
func (_TBTCSystem *TBTCSystemSession) FinalizeSignerFeeDivisorUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeSignerFeeDivisorUpdate(&_TBTCSystem.TransactOpts)
}

// FinalizeSignerFeeDivisorUpdate is a paid mutator transaction binding the contract method 0x0ce0e700.
//
// Solidity: function finalizeSignerFeeDivisorUpdate() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) FinalizeSignerFeeDivisorUpdate() (*types.Transaction, error) {
	return _TBTCSystem.Contract.FinalizeSignerFeeDivisorUpdate(&_TBTCSystem.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x1c52d90c.
//
// Solidity: function initialize(address _defaultKeepFactory, address _depositFactory, address _masterDepositAddress, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachine, uint16 _keepThreshold, uint16 _keepSize) returns()
func (_TBTCSystem *TBTCSystemTransactor) Initialize(opts *bind.TransactOpts, _defaultKeepFactory common.Address, _depositFactory common.Address, _masterDepositAddress common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachine common.Address, _keepThreshold uint16, _keepSize uint16) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "initialize", _defaultKeepFactory, _depositFactory, _masterDepositAddress, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachine, _keepThreshold, _keepSize)
}

// Initialize is a paid mutator transaction binding the contract method 0x1c52d90c.
//
// Solidity: function initialize(address _defaultKeepFactory, address _depositFactory, address _masterDepositAddress, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachine, uint16 _keepThreshold, uint16 _keepSize) returns()
func (_TBTCSystem *TBTCSystemSession) Initialize(_defaultKeepFactory common.Address, _depositFactory common.Address, _masterDepositAddress common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachine common.Address, _keepThreshold uint16, _keepSize uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.Initialize(&_TBTCSystem.TransactOpts, _defaultKeepFactory, _depositFactory, _masterDepositAddress, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachine, _keepThreshold, _keepSize)
}

// Initialize is a paid mutator transaction binding the contract method 0x1c52d90c.
//
// Solidity: function initialize(address _defaultKeepFactory, address _depositFactory, address _masterDepositAddress, address _tbtcToken, address _tbtcDepositToken, address _feeRebateToken, address _vendingMachine, uint16 _keepThreshold, uint16 _keepSize) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) Initialize(_defaultKeepFactory common.Address, _depositFactory common.Address, _masterDepositAddress common.Address, _tbtcToken common.Address, _tbtcDepositToken common.Address, _feeRebateToken common.Address, _vendingMachine common.Address, _keepThreshold uint16, _keepSize uint16) (*types.Transaction, error) {
	return _TBTCSystem.Contract.Initialize(&_TBTCSystem.TransactOpts, _defaultKeepFactory, _depositFactory, _masterDepositAddress, _tbtcToken, _tbtcDepositToken, _feeRebateToken, _vendingMachine, _keepThreshold, _keepSize)
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_TBTCSystem *TBTCSystemTransactor) LogCourtesyCalled(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logCourtesyCalled")
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_TBTCSystem *TBTCSystemSession) LogCourtesyCalled() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogCourtesyCalled(&_TBTCSystem.TransactOpts)
}

// LogCourtesyCalled is a paid mutator transaction binding the contract method 0x22a147e6.
//
// Solidity: function logCourtesyCalled() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogCourtesyCalled() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogCourtesyCalled(&_TBTCSystem.TransactOpts)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogCreated(opts *bind.TransactOpts, _keepAddress common.Address) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logCreated", _keepAddress)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_TBTCSystem *TBTCSystemSession) LogCreated(_keepAddress common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogCreated(&_TBTCSystem.TransactOpts, _keepAddress)
}

// LogCreated is a paid mutator transaction binding the contract method 0x282bfd38.
//
// Solidity: function logCreated(address _keepAddress) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogCreated(_keepAddress common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogCreated(&_TBTCSystem.TransactOpts, _keepAddress)
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_TBTCSystem *TBTCSystemTransactor) LogExitedCourtesyCall(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logExitedCourtesyCall")
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_TBTCSystem *TBTCSystemSession) LogExitedCourtesyCall() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogExitedCourtesyCall(&_TBTCSystem.TransactOpts)
}

// LogExitedCourtesyCall is a paid mutator transaction binding the contract method 0x22e5724c.
//
// Solidity: function logExitedCourtesyCall() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogExitedCourtesyCall() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogExitedCourtesyCall(&_TBTCSystem.TransactOpts)
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_TBTCSystem *TBTCSystemTransactor) LogFraudDuringSetup(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logFraudDuringSetup")
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_TBTCSystem *TBTCSystemSession) LogFraudDuringSetup() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFraudDuringSetup(&_TBTCSystem.TransactOpts)
}

// LogFraudDuringSetup is a paid mutator transaction binding the contract method 0xe2c50ad8.
//
// Solidity: function logFraudDuringSetup() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogFraudDuringSetup() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFraudDuringSetup(&_TBTCSystem.TransactOpts)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogFunded(opts *bind.TransactOpts, _txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logFunded", _txid)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemSession) LogFunded(_txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFunded(&_TBTCSystem.TransactOpts, _txid)
}

// LogFunded is a paid mutator transaction binding the contract method 0x7ed451a4.
//
// Solidity: function logFunded(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogFunded(_txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFunded(&_TBTCSystem.TransactOpts, _txid)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogFunderRequestedAbort(opts *bind.TransactOpts, _abortOutputScript []byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logFunderRequestedAbort", _abortOutputScript)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_TBTCSystem *TBTCSystemSession) LogFunderRequestedAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFunderRequestedAbort(&_TBTCSystem.TransactOpts, _abortOutputScript)
}

// LogFunderRequestedAbort is a paid mutator transaction binding the contract method 0xce2c07ce.
//
// Solidity: function logFunderRequestedAbort(bytes _abortOutputScript) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogFunderRequestedAbort(_abortOutputScript []byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogFunderRequestedAbort(&_TBTCSystem.TransactOpts, _abortOutputScript)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogGotRedemptionSignature(opts *bind.TransactOpts, _digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logGotRedemptionSignature", _digest, _r, _s)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_TBTCSystem *TBTCSystemSession) LogGotRedemptionSignature(_digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogGotRedemptionSignature(&_TBTCSystem.TransactOpts, _digest, _r, _s)
}

// LogGotRedemptionSignature is a paid mutator transaction binding the contract method 0xf760621e.
//
// Solidity: function logGotRedemptionSignature(bytes32 _digest, bytes32 _r, bytes32 _s) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogGotRedemptionSignature(_digest [32]byte, _r [32]byte, _s [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogGotRedemptionSignature(&_TBTCSystem.TransactOpts, _digest, _r, _s)
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_TBTCSystem *TBTCSystemTransactor) LogLiquidated(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logLiquidated")
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_TBTCSystem *TBTCSystemSession) LogLiquidated() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogLiquidated(&_TBTCSystem.TransactOpts)
}

// LogLiquidated is a paid mutator transaction binding the contract method 0xc8fba243.
//
// Solidity: function logLiquidated() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogLiquidated() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogLiquidated(&_TBTCSystem.TransactOpts)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogRedeemed(opts *bind.TransactOpts, _txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logRedeemed", _txid)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemSession) LogRedeemed(_txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRedeemed(&_TBTCSystem.TransactOpts, _txid)
}

// LogRedeemed is a paid mutator transaction binding the contract method 0x6e1ba283.
//
// Solidity: function logRedeemed(bytes32 _txid) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogRedeemed(_txid [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRedeemed(&_TBTCSystem.TransactOpts, _txid)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogRedemptionRequested(opts *bind.TransactOpts, _requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logRedemptionRequested", _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_TBTCSystem *TBTCSystemSession) LogRedemptionRequested(_requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRedemptionRequested(&_TBTCSystem.TransactOpts, _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRedemptionRequested is a paid mutator transaction binding the contract method 0x18e647dd.
//
// Solidity: function logRedemptionRequested(address _requester, bytes32 _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogRedemptionRequested(_requester common.Address, _digest [32]byte, _utxoValue *big.Int, _redeemerOutputScript []byte, _requestedFee *big.Int, _outpoint []byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRedemptionRequested(&_TBTCSystem.TransactOpts, _requester, _digest, _utxoValue, _redeemerOutputScript, _requestedFee, _outpoint)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogRegisteredPubkey(opts *bind.TransactOpts, _signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logRegisteredPubkey", _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_TBTCSystem *TBTCSystemSession) LogRegisteredPubkey(_signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRegisteredPubkey(&_TBTCSystem.TransactOpts, _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogRegisteredPubkey is a paid mutator transaction binding the contract method 0x869f9469.
//
// Solidity: function logRegisteredPubkey(bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogRegisteredPubkey(_signingGroupPubkeyX [32]byte, _signingGroupPubkeyY [32]byte) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogRegisteredPubkey(&_TBTCSystem.TransactOpts, _signingGroupPubkeyX, _signingGroupPubkeyY)
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_TBTCSystem *TBTCSystemTransactor) LogSetupFailed(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logSetupFailed")
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_TBTCSystem *TBTCSystemSession) LogSetupFailed() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogSetupFailed(&_TBTCSystem.TransactOpts)
}

// LogSetupFailed is a paid mutator transaction binding the contract method 0xa831c816.
//
// Solidity: function logSetupFailed() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogSetupFailed() (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogSetupFailed(&_TBTCSystem.TransactOpts)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_TBTCSystem *TBTCSystemTransactor) LogStartedLiquidation(opts *bind.TransactOpts, _wasFraud bool) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "logStartedLiquidation", _wasFraud)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_TBTCSystem *TBTCSystemSession) LogStartedLiquidation(_wasFraud bool) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogStartedLiquidation(&_TBTCSystem.TransactOpts, _wasFraud)
}

// LogStartedLiquidation is a paid mutator transaction binding the contract method 0x3aac3467.
//
// Solidity: function logStartedLiquidation(bool _wasFraud) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) LogStartedLiquidation(_wasFraud bool) (*types.Transaction, error) {
	return _TBTCSystem.Contract.LogStartedLiquidation(&_TBTCSystem.TransactOpts, _wasFraud)
}

// RefreshMinimumBondableValue is a paid mutator transaction binding the contract method 0x7c75b115.
//
// Solidity: function refreshMinimumBondableValue() returns()
func (_TBTCSystem *TBTCSystemTransactor) RefreshMinimumBondableValue(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "refreshMinimumBondableValue")
}

// RefreshMinimumBondableValue is a paid mutator transaction binding the contract method 0x7c75b115.
//
// Solidity: function refreshMinimumBondableValue() returns()
func (_TBTCSystem *TBTCSystemSession) RefreshMinimumBondableValue() (*types.Transaction, error) {
	return _TBTCSystem.Contract.RefreshMinimumBondableValue(&_TBTCSystem.TransactOpts)
}

// RefreshMinimumBondableValue is a paid mutator transaction binding the contract method 0x7c75b115.
//
// Solidity: function refreshMinimumBondableValue() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) RefreshMinimumBondableValue() (*types.Transaction, error) {
	return _TBTCSystem.Contract.RefreshMinimumBondableValue(&_TBTCSystem.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TBTCSystem *TBTCSystemTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TBTCSystem *TBTCSystemSession) RenounceOwnership() (*types.Transaction, error) {
	return _TBTCSystem.Contract.RenounceOwnership(&_TBTCSystem.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TBTCSystem.Contract.RenounceOwnership(&_TBTCSystem.TransactOpts)
}

// RequestNewKeep is a paid mutator transaction binding the contract method 0x82f91968.
//
// Solidity: function requestNewKeep(uint64 _requestedLotSizeSatoshis, uint256 _maxSecuredLifetime) payable returns(address)
func (_TBTCSystem *TBTCSystemTransactor) RequestNewKeep(opts *bind.TransactOpts, _requestedLotSizeSatoshis uint64, _maxSecuredLifetime *big.Int) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "requestNewKeep", _requestedLotSizeSatoshis, _maxSecuredLifetime)
}

// RequestNewKeep is a paid mutator transaction binding the contract method 0x82f91968.
//
// Solidity: function requestNewKeep(uint64 _requestedLotSizeSatoshis, uint256 _maxSecuredLifetime) payable returns(address)
func (_TBTCSystem *TBTCSystemSession) RequestNewKeep(_requestedLotSizeSatoshis uint64, _maxSecuredLifetime *big.Int) (*types.Transaction, error) {
	return _TBTCSystem.Contract.RequestNewKeep(&_TBTCSystem.TransactOpts, _requestedLotSizeSatoshis, _maxSecuredLifetime)
}

// RequestNewKeep is a paid mutator transaction binding the contract method 0x82f91968.
//
// Solidity: function requestNewKeep(uint64 _requestedLotSizeSatoshis, uint256 _maxSecuredLifetime) payable returns(address)
func (_TBTCSystem *TBTCSystemTransactorSession) RequestNewKeep(_requestedLotSizeSatoshis uint64, _maxSecuredLifetime *big.Int) (*types.Transaction, error) {
	return _TBTCSystem.Contract.RequestNewKeep(&_TBTCSystem.TransactOpts, _requestedLotSizeSatoshis, _maxSecuredLifetime)
}

// ResumeNewDeposits is a paid mutator transaction binding the contract method 0x7c33fc05.
//
// Solidity: function resumeNewDeposits() returns()
func (_TBTCSystem *TBTCSystemTransactor) ResumeNewDeposits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "resumeNewDeposits")
}

// ResumeNewDeposits is a paid mutator transaction binding the contract method 0x7c33fc05.
//
// Solidity: function resumeNewDeposits() returns()
func (_TBTCSystem *TBTCSystemSession) ResumeNewDeposits() (*types.Transaction, error) {
	return _TBTCSystem.Contract.ResumeNewDeposits(&_TBTCSystem.TransactOpts)
}

// ResumeNewDeposits is a paid mutator transaction binding the contract method 0x7c33fc05.
//
// Solidity: function resumeNewDeposits() returns()
func (_TBTCSystem *TBTCSystemTransactorSession) ResumeNewDeposits() (*types.Transaction, error) {
	return _TBTCSystem.Contract.ResumeNewDeposits(&_TBTCSystem.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TBTCSystem *TBTCSystemTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TBTCSystem.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TBTCSystem *TBTCSystemSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.TransferOwnership(&_TBTCSystem.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TBTCSystem *TBTCSystemTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TBTCSystem.Contract.TransferOwnership(&_TBTCSystem.TransactOpts, newOwner)
}

// TBTCSystemAllowNewDepositsUpdatedIterator is returned from FilterAllowNewDepositsUpdated and is used to iterate over the raw logs and unpacked data for AllowNewDepositsUpdated events raised by the TBTCSystem contract.
type TBTCSystemAllowNewDepositsUpdatedIterator struct {
	Event *TBTCSystemAllowNewDepositsUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemAllowNewDepositsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemAllowNewDepositsUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemAllowNewDepositsUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemAllowNewDepositsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemAllowNewDepositsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemAllowNewDepositsUpdated represents a AllowNewDepositsUpdated event raised by the TBTCSystem contract.
type TBTCSystemAllowNewDepositsUpdated struct {
	AllowNewDeposits bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterAllowNewDepositsUpdated is a free log retrieval operation binding the contract event 0x3a854be74be62dd3ba5f0fdb7aa5b535683f999e90cda09ba75a2d99b2722523.
//
// Solidity: event AllowNewDepositsUpdated(bool _allowNewDeposits)
func (_TBTCSystem *TBTCSystemFilterer) FilterAllowNewDepositsUpdated(opts *bind.FilterOpts) (*TBTCSystemAllowNewDepositsUpdatedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "AllowNewDepositsUpdated")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemAllowNewDepositsUpdatedIterator{contract: _TBTCSystem.contract, event: "AllowNewDepositsUpdated", logs: logs, sub: sub}, nil
}

// WatchAllowNewDepositsUpdated is a free log subscription operation binding the contract event 0x3a854be74be62dd3ba5f0fdb7aa5b535683f999e90cda09ba75a2d99b2722523.
//
// Solidity: event AllowNewDepositsUpdated(bool _allowNewDeposits)
func (_TBTCSystem *TBTCSystemFilterer) WatchAllowNewDepositsUpdated(opts *bind.WatchOpts, sink chan<- *TBTCSystemAllowNewDepositsUpdated) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "AllowNewDepositsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemAllowNewDepositsUpdated)
				if err := _TBTCSystem.contract.UnpackLog(event, "AllowNewDepositsUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAllowNewDepositsUpdated is a log parse operation binding the contract event 0x3a854be74be62dd3ba5f0fdb7aa5b535683f999e90cda09ba75a2d99b2722523.
//
// Solidity: event AllowNewDepositsUpdated(bool _allowNewDeposits)
func (_TBTCSystem *TBTCSystemFilterer) ParseAllowNewDepositsUpdated(log types.Log) (*TBTCSystemAllowNewDepositsUpdated, error) {
	event := new(TBTCSystemAllowNewDepositsUpdated)
	if err := _TBTCSystem.contract.UnpackLog(event, "AllowNewDepositsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemCollateralizationThresholdsUpdateStartedIterator is returned from FilterCollateralizationThresholdsUpdateStarted and is used to iterate over the raw logs and unpacked data for CollateralizationThresholdsUpdateStarted events raised by the TBTCSystem contract.
type TBTCSystemCollateralizationThresholdsUpdateStartedIterator struct {
	Event *TBTCSystemCollateralizationThresholdsUpdateStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemCollateralizationThresholdsUpdateStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemCollateralizationThresholdsUpdateStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemCollateralizationThresholdsUpdateStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemCollateralizationThresholdsUpdateStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemCollateralizationThresholdsUpdateStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemCollateralizationThresholdsUpdateStarted represents a CollateralizationThresholdsUpdateStarted event raised by the TBTCSystem contract.
type TBTCSystemCollateralizationThresholdsUpdateStarted struct {
	InitialCollateralizedPercent                uint16
	UndercollateralizedThresholdPercent         uint16
	SeverelyUndercollateralizedThresholdPercent uint16
	Timestamp                                   *big.Int
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterCollateralizationThresholdsUpdateStarted is a free log retrieval operation binding the contract event 0xc9e225c2db3e9f70966e6c0403de785bcda0172ec4e41111a6e8f4b85b1f30fb.
//
// Solidity: event CollateralizationThresholdsUpdateStarted(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterCollateralizationThresholdsUpdateStarted(opts *bind.FilterOpts) (*TBTCSystemCollateralizationThresholdsUpdateStartedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "CollateralizationThresholdsUpdateStarted")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemCollateralizationThresholdsUpdateStartedIterator{contract: _TBTCSystem.contract, event: "CollateralizationThresholdsUpdateStarted", logs: logs, sub: sub}, nil
}

// WatchCollateralizationThresholdsUpdateStarted is a free log subscription operation binding the contract event 0xc9e225c2db3e9f70966e6c0403de785bcda0172ec4e41111a6e8f4b85b1f30fb.
//
// Solidity: event CollateralizationThresholdsUpdateStarted(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchCollateralizationThresholdsUpdateStarted(opts *bind.WatchOpts, sink chan<- *TBTCSystemCollateralizationThresholdsUpdateStarted) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "CollateralizationThresholdsUpdateStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemCollateralizationThresholdsUpdateStarted)
				if err := _TBTCSystem.contract.UnpackLog(event, "CollateralizationThresholdsUpdateStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralizationThresholdsUpdateStarted is a log parse operation binding the contract event 0xc9e225c2db3e9f70966e6c0403de785bcda0172ec4e41111a6e8f4b85b1f30fb.
//
// Solidity: event CollateralizationThresholdsUpdateStarted(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseCollateralizationThresholdsUpdateStarted(log types.Log) (*TBTCSystemCollateralizationThresholdsUpdateStarted, error) {
	event := new(TBTCSystemCollateralizationThresholdsUpdateStarted)
	if err := _TBTCSystem.contract.UnpackLog(event, "CollateralizationThresholdsUpdateStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemCollateralizationThresholdsUpdatedIterator is returned from FilterCollateralizationThresholdsUpdated and is used to iterate over the raw logs and unpacked data for CollateralizationThresholdsUpdated events raised by the TBTCSystem contract.
type TBTCSystemCollateralizationThresholdsUpdatedIterator struct {
	Event *TBTCSystemCollateralizationThresholdsUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemCollateralizationThresholdsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemCollateralizationThresholdsUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemCollateralizationThresholdsUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemCollateralizationThresholdsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemCollateralizationThresholdsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemCollateralizationThresholdsUpdated represents a CollateralizationThresholdsUpdated event raised by the TBTCSystem contract.
type TBTCSystemCollateralizationThresholdsUpdated struct {
	InitialCollateralizedPercent                uint16
	UndercollateralizedThresholdPercent         uint16
	SeverelyUndercollateralizedThresholdPercent uint16
	Raw                                         types.Log // Blockchain specific contextual infos
}

// FilterCollateralizationThresholdsUpdated is a free log retrieval operation binding the contract event 0x07ac9ce7dc4b2edb6435fb2255e9e867f357ef2052b982ce468442aa9d6c1d50.
//
// Solidity: event CollateralizationThresholdsUpdated(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent)
func (_TBTCSystem *TBTCSystemFilterer) FilterCollateralizationThresholdsUpdated(opts *bind.FilterOpts) (*TBTCSystemCollateralizationThresholdsUpdatedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "CollateralizationThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemCollateralizationThresholdsUpdatedIterator{contract: _TBTCSystem.contract, event: "CollateralizationThresholdsUpdated", logs: logs, sub: sub}, nil
}

// WatchCollateralizationThresholdsUpdated is a free log subscription operation binding the contract event 0x07ac9ce7dc4b2edb6435fb2255e9e867f357ef2052b982ce468442aa9d6c1d50.
//
// Solidity: event CollateralizationThresholdsUpdated(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent)
func (_TBTCSystem *TBTCSystemFilterer) WatchCollateralizationThresholdsUpdated(opts *bind.WatchOpts, sink chan<- *TBTCSystemCollateralizationThresholdsUpdated) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "CollateralizationThresholdsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemCollateralizationThresholdsUpdated)
				if err := _TBTCSystem.contract.UnpackLog(event, "CollateralizationThresholdsUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCollateralizationThresholdsUpdated is a log parse operation binding the contract event 0x07ac9ce7dc4b2edb6435fb2255e9e867f357ef2052b982ce468442aa9d6c1d50.
//
// Solidity: event CollateralizationThresholdsUpdated(uint16 _initialCollateralizedPercent, uint16 _undercollateralizedThresholdPercent, uint16 _severelyUndercollateralizedThresholdPercent)
func (_TBTCSystem *TBTCSystemFilterer) ParseCollateralizationThresholdsUpdated(log types.Log) (*TBTCSystemCollateralizationThresholdsUpdated, error) {
	event := new(TBTCSystemCollateralizationThresholdsUpdated)
	if err := _TBTCSystem.contract.UnpackLog(event, "CollateralizationThresholdsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemCourtesyCalledIterator is returned from FilterCourtesyCalled and is used to iterate over the raw logs and unpacked data for CourtesyCalled events raised by the TBTCSystem contract.
type TBTCSystemCourtesyCalledIterator struct {
	Event *TBTCSystemCourtesyCalled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemCourtesyCalledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemCourtesyCalled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemCourtesyCalled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemCourtesyCalledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemCourtesyCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemCourtesyCalled represents a CourtesyCalled event raised by the TBTCSystem contract.
type TBTCSystemCourtesyCalled struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterCourtesyCalled is a free log retrieval operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterCourtesyCalled(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemCourtesyCalledIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "CourtesyCalled", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemCourtesyCalledIterator{contract: _TBTCSystem.contract, event: "CourtesyCalled", logs: logs, sub: sub}, nil
}

// WatchCourtesyCalled is a free log subscription operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchCourtesyCalled(opts *bind.WatchOpts, sink chan<- *TBTCSystemCourtesyCalled, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "CourtesyCalled", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemCourtesyCalled)
				if err := _TBTCSystem.contract.UnpackLog(event, "CourtesyCalled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCourtesyCalled is a log parse operation binding the contract event 0x6e7b45210b79c12cd1332babd8d86c0bbb9ca898a89ce0404f17064dbfba18c0.
//
// Solidity: event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseCourtesyCalled(log types.Log) (*TBTCSystemCourtesyCalled, error) {
	event := new(TBTCSystemCourtesyCalled)
	if err := _TBTCSystem.contract.UnpackLog(event, "CourtesyCalled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemCreatedIterator is returned from FilterCreated and is used to iterate over the raw logs and unpacked data for Created events raised by the TBTCSystem contract.
type TBTCSystemCreatedIterator struct {
	Event *TBTCSystemCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemCreated represents a Created event raised by the TBTCSystem contract.
type TBTCSystemCreated struct {
	DepositContractAddress common.Address
	KeepAddress            common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterCreated is a free log retrieval operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterCreated(opts *bind.FilterOpts, _depositContractAddress []common.Address, _keepAddress []common.Address) (*TBTCSystemCreatedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _keepAddressRule []interface{}
	for _, _keepAddressItem := range _keepAddress {
		_keepAddressRule = append(_keepAddressRule, _keepAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "Created", _depositContractAddressRule, _keepAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemCreatedIterator{contract: _TBTCSystem.contract, event: "Created", logs: logs, sub: sub}, nil
}

// WatchCreated is a free log subscription operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchCreated(opts *bind.WatchOpts, sink chan<- *TBTCSystemCreated, _depositContractAddress []common.Address, _keepAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _keepAddressRule []interface{}
	for _, _keepAddressItem := range _keepAddress {
		_keepAddressRule = append(_keepAddressRule, _keepAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "Created", _depositContractAddressRule, _keepAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemCreated)
				if err := _TBTCSystem.contract.UnpackLog(event, "Created", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreated is a log parse operation binding the contract event 0x822b3073be62c5c7f143c2dcd71ee266434ee935d90a1eec3be34710ac8ec1a2.
//
// Solidity: event Created(address indexed _depositContractAddress, address indexed _keepAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseCreated(log types.Log) (*TBTCSystemCreated, error) {
	event := new(TBTCSystemCreated)
	if err := _TBTCSystem.contract.UnpackLog(event, "Created", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemEthBtcPriceFeedAddedIterator is returned from FilterEthBtcPriceFeedAdded and is used to iterate over the raw logs and unpacked data for EthBtcPriceFeedAdded events raised by the TBTCSystem contract.
type TBTCSystemEthBtcPriceFeedAddedIterator struct {
	Event *TBTCSystemEthBtcPriceFeedAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemEthBtcPriceFeedAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemEthBtcPriceFeedAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemEthBtcPriceFeedAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemEthBtcPriceFeedAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemEthBtcPriceFeedAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemEthBtcPriceFeedAdded represents a EthBtcPriceFeedAdded event raised by the TBTCSystem contract.
type TBTCSystemEthBtcPriceFeedAdded struct {
	PriceFeed common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEthBtcPriceFeedAdded is a free log retrieval operation binding the contract event 0x5e4bd1f6e413d39e172d96a88ee6b1b2ba9e1a6207e2ca34fa8c3ccd152ff21a.
//
// Solidity: event EthBtcPriceFeedAdded(address _priceFeed)
func (_TBTCSystem *TBTCSystemFilterer) FilterEthBtcPriceFeedAdded(opts *bind.FilterOpts) (*TBTCSystemEthBtcPriceFeedAddedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "EthBtcPriceFeedAdded")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemEthBtcPriceFeedAddedIterator{contract: _TBTCSystem.contract, event: "EthBtcPriceFeedAdded", logs: logs, sub: sub}, nil
}

// WatchEthBtcPriceFeedAdded is a free log subscription operation binding the contract event 0x5e4bd1f6e413d39e172d96a88ee6b1b2ba9e1a6207e2ca34fa8c3ccd152ff21a.
//
// Solidity: event EthBtcPriceFeedAdded(address _priceFeed)
func (_TBTCSystem *TBTCSystemFilterer) WatchEthBtcPriceFeedAdded(opts *bind.WatchOpts, sink chan<- *TBTCSystemEthBtcPriceFeedAdded) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "EthBtcPriceFeedAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemEthBtcPriceFeedAdded)
				if err := _TBTCSystem.contract.UnpackLog(event, "EthBtcPriceFeedAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEthBtcPriceFeedAdded is a log parse operation binding the contract event 0x5e4bd1f6e413d39e172d96a88ee6b1b2ba9e1a6207e2ca34fa8c3ccd152ff21a.
//
// Solidity: event EthBtcPriceFeedAdded(address _priceFeed)
func (_TBTCSystem *TBTCSystemFilterer) ParseEthBtcPriceFeedAdded(log types.Log) (*TBTCSystemEthBtcPriceFeedAdded, error) {
	event := new(TBTCSystemEthBtcPriceFeedAdded)
	if err := _TBTCSystem.contract.UnpackLog(event, "EthBtcPriceFeedAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemEthBtcPriceFeedAdditionStartedIterator is returned from FilterEthBtcPriceFeedAdditionStarted and is used to iterate over the raw logs and unpacked data for EthBtcPriceFeedAdditionStarted events raised by the TBTCSystem contract.
type TBTCSystemEthBtcPriceFeedAdditionStartedIterator struct {
	Event *TBTCSystemEthBtcPriceFeedAdditionStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemEthBtcPriceFeedAdditionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemEthBtcPriceFeedAdditionStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemEthBtcPriceFeedAdditionStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemEthBtcPriceFeedAdditionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemEthBtcPriceFeedAdditionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemEthBtcPriceFeedAdditionStarted represents a EthBtcPriceFeedAdditionStarted event raised by the TBTCSystem contract.
type TBTCSystemEthBtcPriceFeedAdditionStarted struct {
	PriceFeed common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEthBtcPriceFeedAdditionStarted is a free log retrieval operation binding the contract event 0x5a3d3d9197f5c60c16de28887dccf83284a4fd034b930272637c83307b4fffe7.
//
// Solidity: event EthBtcPriceFeedAdditionStarted(address _priceFeed, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterEthBtcPriceFeedAdditionStarted(opts *bind.FilterOpts) (*TBTCSystemEthBtcPriceFeedAdditionStartedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "EthBtcPriceFeedAdditionStarted")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemEthBtcPriceFeedAdditionStartedIterator{contract: _TBTCSystem.contract, event: "EthBtcPriceFeedAdditionStarted", logs: logs, sub: sub}, nil
}

// WatchEthBtcPriceFeedAdditionStarted is a free log subscription operation binding the contract event 0x5a3d3d9197f5c60c16de28887dccf83284a4fd034b930272637c83307b4fffe7.
//
// Solidity: event EthBtcPriceFeedAdditionStarted(address _priceFeed, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchEthBtcPriceFeedAdditionStarted(opts *bind.WatchOpts, sink chan<- *TBTCSystemEthBtcPriceFeedAdditionStarted) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "EthBtcPriceFeedAdditionStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemEthBtcPriceFeedAdditionStarted)
				if err := _TBTCSystem.contract.UnpackLog(event, "EthBtcPriceFeedAdditionStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEthBtcPriceFeedAdditionStarted is a log parse operation binding the contract event 0x5a3d3d9197f5c60c16de28887dccf83284a4fd034b930272637c83307b4fffe7.
//
// Solidity: event EthBtcPriceFeedAdditionStarted(address _priceFeed, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseEthBtcPriceFeedAdditionStarted(log types.Log) (*TBTCSystemEthBtcPriceFeedAdditionStarted, error) {
	event := new(TBTCSystemEthBtcPriceFeedAdditionStarted)
	if err := _TBTCSystem.contract.UnpackLog(event, "EthBtcPriceFeedAdditionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemExitedCourtesyCallIterator is returned from FilterExitedCourtesyCall and is used to iterate over the raw logs and unpacked data for ExitedCourtesyCall events raised by the TBTCSystem contract.
type TBTCSystemExitedCourtesyCallIterator struct {
	Event *TBTCSystemExitedCourtesyCall // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemExitedCourtesyCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemExitedCourtesyCall)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemExitedCourtesyCall)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemExitedCourtesyCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemExitedCourtesyCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemExitedCourtesyCall represents a ExitedCourtesyCall event raised by the TBTCSystem contract.
type TBTCSystemExitedCourtesyCall struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterExitedCourtesyCall is a free log retrieval operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterExitedCourtesyCall(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemExitedCourtesyCallIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "ExitedCourtesyCall", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemExitedCourtesyCallIterator{contract: _TBTCSystem.contract, event: "ExitedCourtesyCall", logs: logs, sub: sub}, nil
}

// WatchExitedCourtesyCall is a free log subscription operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchExitedCourtesyCall(opts *bind.WatchOpts, sink chan<- *TBTCSystemExitedCourtesyCall, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "ExitedCourtesyCall", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemExitedCourtesyCall)
				if err := _TBTCSystem.contract.UnpackLog(event, "ExitedCourtesyCall", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExitedCourtesyCall is a log parse operation binding the contract event 0x07f0eaafadb9abb1d28da85d4b4c74f1939fd61b535c7f5ab501f618f07e76ee.
//
// Solidity: event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseExitedCourtesyCall(log types.Log) (*TBTCSystemExitedCourtesyCall, error) {
	event := new(TBTCSystemExitedCourtesyCall)
	if err := _TBTCSystem.contract.UnpackLog(event, "ExitedCourtesyCall", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemFraudDuringSetupIterator is returned from FilterFraudDuringSetup and is used to iterate over the raw logs and unpacked data for FraudDuringSetup events raised by the TBTCSystem contract.
type TBTCSystemFraudDuringSetupIterator struct {
	Event *TBTCSystemFraudDuringSetup // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemFraudDuringSetupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemFraudDuringSetup)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemFraudDuringSetup)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemFraudDuringSetupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemFraudDuringSetupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemFraudDuringSetup represents a FraudDuringSetup event raised by the TBTCSystem contract.
type TBTCSystemFraudDuringSetup struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFraudDuringSetup is a free log retrieval operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterFraudDuringSetup(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemFraudDuringSetupIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "FraudDuringSetup", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemFraudDuringSetupIterator{contract: _TBTCSystem.contract, event: "FraudDuringSetup", logs: logs, sub: sub}, nil
}

// WatchFraudDuringSetup is a free log subscription operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchFraudDuringSetup(opts *bind.WatchOpts, sink chan<- *TBTCSystemFraudDuringSetup, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "FraudDuringSetup", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemFraudDuringSetup)
				if err := _TBTCSystem.contract.UnpackLog(event, "FraudDuringSetup", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFraudDuringSetup is a log parse operation binding the contract event 0x1e61af503f1d7de21d5300094c18bf8700f82b2951a4d54dd2adda13f6b3da30.
//
// Solidity: event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseFraudDuringSetup(log types.Log) (*TBTCSystemFraudDuringSetup, error) {
	event := new(TBTCSystemFraudDuringSetup)
	if err := _TBTCSystem.contract.UnpackLog(event, "FraudDuringSetup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemFundedIterator is returned from FilterFunded and is used to iterate over the raw logs and unpacked data for Funded events raised by the TBTCSystem contract.
type TBTCSystemFundedIterator struct {
	Event *TBTCSystemFunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemFunded represents a Funded event raised by the TBTCSystem contract.
type TBTCSystemFunded struct {
	DepositContractAddress common.Address
	Txid                   [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFunded is a free log retrieval operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterFunded(opts *bind.FilterOpts, _depositContractAddress []common.Address, _txid [][32]byte) (*TBTCSystemFundedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "Funded", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemFundedIterator{contract: _TBTCSystem.contract, event: "Funded", logs: logs, sub: sub}, nil
}

// WatchFunded is a free log subscription operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *TBTCSystemFunded, _depositContractAddress []common.Address, _txid [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "Funded", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemFunded)
				if err := _TBTCSystem.contract.UnpackLog(event, "Funded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFunded is a log parse operation binding the contract event 0xe34c70bd3e03956978a5c76d2ea5f3a60819171afea6dee4fc12b2e45f72d43d.
//
// Solidity: event Funded(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseFunded(log types.Log) (*TBTCSystemFunded, error) {
	event := new(TBTCSystemFunded)
	if err := _TBTCSystem.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemFunderAbortRequestedIterator is returned from FilterFunderAbortRequested and is used to iterate over the raw logs and unpacked data for FunderAbortRequested events raised by the TBTCSystem contract.
type TBTCSystemFunderAbortRequestedIterator struct {
	Event *TBTCSystemFunderAbortRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemFunderAbortRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemFunderAbortRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemFunderAbortRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemFunderAbortRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemFunderAbortRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemFunderAbortRequested represents a FunderAbortRequested event raised by the TBTCSystem contract.
type TBTCSystemFunderAbortRequested struct {
	DepositContractAddress common.Address
	AbortOutputScript      []byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFunderAbortRequested is a free log retrieval operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_TBTCSystem *TBTCSystemFilterer) FilterFunderAbortRequested(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemFunderAbortRequestedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "FunderAbortRequested", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemFunderAbortRequestedIterator{contract: _TBTCSystem.contract, event: "FunderAbortRequested", logs: logs, sub: sub}, nil
}

// WatchFunderAbortRequested is a free log subscription operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_TBTCSystem *TBTCSystemFilterer) WatchFunderAbortRequested(opts *bind.WatchOpts, sink chan<- *TBTCSystemFunderAbortRequested, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "FunderAbortRequested", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemFunderAbortRequested)
				if err := _TBTCSystem.contract.UnpackLog(event, "FunderAbortRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFunderAbortRequested is a log parse operation binding the contract event 0xa6e9673b5d53b3fe3c62b6459720f9c2a1b129d4f69acb771404ba8681b6a930.
//
// Solidity: event FunderAbortRequested(address indexed _depositContractAddress, bytes _abortOutputScript)
func (_TBTCSystem *TBTCSystemFilterer) ParseFunderAbortRequested(log types.Log) (*TBTCSystemFunderAbortRequested, error) {
	event := new(TBTCSystemFunderAbortRequested)
	if err := _TBTCSystem.contract.UnpackLog(event, "FunderAbortRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemGotRedemptionSignatureIterator is returned from FilterGotRedemptionSignature and is used to iterate over the raw logs and unpacked data for GotRedemptionSignature events raised by the TBTCSystem contract.
type TBTCSystemGotRedemptionSignatureIterator struct {
	Event *TBTCSystemGotRedemptionSignature // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemGotRedemptionSignatureIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemGotRedemptionSignature)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemGotRedemptionSignature)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemGotRedemptionSignatureIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemGotRedemptionSignatureIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemGotRedemptionSignature represents a GotRedemptionSignature event raised by the TBTCSystem contract.
type TBTCSystemGotRedemptionSignature struct {
	DepositContractAddress common.Address
	Digest                 [32]byte
	R                      [32]byte
	S                      [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterGotRedemptionSignature is a free log retrieval operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterGotRedemptionSignature(opts *bind.FilterOpts, _depositContractAddress []common.Address, _digest [][32]byte) (*TBTCSystemGotRedemptionSignatureIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "GotRedemptionSignature", _depositContractAddressRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemGotRedemptionSignatureIterator{contract: _TBTCSystem.contract, event: "GotRedemptionSignature", logs: logs, sub: sub}, nil
}

// WatchGotRedemptionSignature is a free log subscription operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchGotRedemptionSignature(opts *bind.WatchOpts, sink chan<- *TBTCSystemGotRedemptionSignature, _depositContractAddress []common.Address, _digest [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "GotRedemptionSignature", _depositContractAddressRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemGotRedemptionSignature)
				if err := _TBTCSystem.contract.UnpackLog(event, "GotRedemptionSignature", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGotRedemptionSignature is a log parse operation binding the contract event 0x7f7d7327762d01d2c4a552ea0be2bc5a76264574a80aa78083e691a840e509f2.
//
// Solidity: event GotRedemptionSignature(address indexed _depositContractAddress, bytes32 indexed _digest, bytes32 _r, bytes32 _s, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseGotRedemptionSignature(log types.Log) (*TBTCSystemGotRedemptionSignature, error) {
	event := new(TBTCSystemGotRedemptionSignature)
	if err := _TBTCSystem.contract.UnpackLog(event, "GotRedemptionSignature", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemKeepFactoriesUpdateStartedIterator is returned from FilterKeepFactoriesUpdateStarted and is used to iterate over the raw logs and unpacked data for KeepFactoriesUpdateStarted events raised by the TBTCSystem contract.
type TBTCSystemKeepFactoriesUpdateStartedIterator struct {
	Event *TBTCSystemKeepFactoriesUpdateStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemKeepFactoriesUpdateStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemKeepFactoriesUpdateStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemKeepFactoriesUpdateStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemKeepFactoriesUpdateStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemKeepFactoriesUpdateStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemKeepFactoriesUpdateStarted represents a KeepFactoriesUpdateStarted event raised by the TBTCSystem contract.
type TBTCSystemKeepFactoriesUpdateStarted struct {
	KeepStakedFactory  common.Address
	FullyBackedFactory common.Address
	FactorySelector    common.Address
	Timestamp          *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterKeepFactoriesUpdateStarted is a free log retrieval operation binding the contract event 0x1608ec8025d64cdb0ed78e62a67d271a33b9d738842a6eb6e6449bc3afab6dca.
//
// Solidity: event KeepFactoriesUpdateStarted(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterKeepFactoriesUpdateStarted(opts *bind.FilterOpts) (*TBTCSystemKeepFactoriesUpdateStartedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "KeepFactoriesUpdateStarted")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemKeepFactoriesUpdateStartedIterator{contract: _TBTCSystem.contract, event: "KeepFactoriesUpdateStarted", logs: logs, sub: sub}, nil
}

// WatchKeepFactoriesUpdateStarted is a free log subscription operation binding the contract event 0x1608ec8025d64cdb0ed78e62a67d271a33b9d738842a6eb6e6449bc3afab6dca.
//
// Solidity: event KeepFactoriesUpdateStarted(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchKeepFactoriesUpdateStarted(opts *bind.WatchOpts, sink chan<- *TBTCSystemKeepFactoriesUpdateStarted) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "KeepFactoriesUpdateStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemKeepFactoriesUpdateStarted)
				if err := _TBTCSystem.contract.UnpackLog(event, "KeepFactoriesUpdateStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseKeepFactoriesUpdateStarted is a log parse operation binding the contract event 0x1608ec8025d64cdb0ed78e62a67d271a33b9d738842a6eb6e6449bc3afab6dca.
//
// Solidity: event KeepFactoriesUpdateStarted(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseKeepFactoriesUpdateStarted(log types.Log) (*TBTCSystemKeepFactoriesUpdateStarted, error) {
	event := new(TBTCSystemKeepFactoriesUpdateStarted)
	if err := _TBTCSystem.contract.UnpackLog(event, "KeepFactoriesUpdateStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemKeepFactoriesUpdatedIterator is returned from FilterKeepFactoriesUpdated and is used to iterate over the raw logs and unpacked data for KeepFactoriesUpdated events raised by the TBTCSystem contract.
type TBTCSystemKeepFactoriesUpdatedIterator struct {
	Event *TBTCSystemKeepFactoriesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemKeepFactoriesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemKeepFactoriesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemKeepFactoriesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemKeepFactoriesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemKeepFactoriesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemKeepFactoriesUpdated represents a KeepFactoriesUpdated event raised by the TBTCSystem contract.
type TBTCSystemKeepFactoriesUpdated struct {
	KeepStakedFactory  common.Address
	FullyBackedFactory common.Address
	FactorySelector    common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterKeepFactoriesUpdated is a free log retrieval operation binding the contract event 0x75cd06e2a95bd62ad447184bf0950b3af3aabd0960994d09da9724686d0c1720.
//
// Solidity: event KeepFactoriesUpdated(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector)
func (_TBTCSystem *TBTCSystemFilterer) FilterKeepFactoriesUpdated(opts *bind.FilterOpts) (*TBTCSystemKeepFactoriesUpdatedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "KeepFactoriesUpdated")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemKeepFactoriesUpdatedIterator{contract: _TBTCSystem.contract, event: "KeepFactoriesUpdated", logs: logs, sub: sub}, nil
}

// WatchKeepFactoriesUpdated is a free log subscription operation binding the contract event 0x75cd06e2a95bd62ad447184bf0950b3af3aabd0960994d09da9724686d0c1720.
//
// Solidity: event KeepFactoriesUpdated(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector)
func (_TBTCSystem *TBTCSystemFilterer) WatchKeepFactoriesUpdated(opts *bind.WatchOpts, sink chan<- *TBTCSystemKeepFactoriesUpdated) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "KeepFactoriesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemKeepFactoriesUpdated)
				if err := _TBTCSystem.contract.UnpackLog(event, "KeepFactoriesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseKeepFactoriesUpdated is a log parse operation binding the contract event 0x75cd06e2a95bd62ad447184bf0950b3af3aabd0960994d09da9724686d0c1720.
//
// Solidity: event KeepFactoriesUpdated(address _keepStakedFactory, address _fullyBackedFactory, address _factorySelector)
func (_TBTCSystem *TBTCSystemFilterer) ParseKeepFactoriesUpdated(log types.Log) (*TBTCSystemKeepFactoriesUpdated, error) {
	event := new(TBTCSystemKeepFactoriesUpdated)
	if err := _TBTCSystem.contract.UnpackLog(event, "KeepFactoriesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemLiquidatedIterator is returned from FilterLiquidated and is used to iterate over the raw logs and unpacked data for Liquidated events raised by the TBTCSystem contract.
type TBTCSystemLiquidatedIterator struct {
	Event *TBTCSystemLiquidated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemLiquidatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemLiquidated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemLiquidated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemLiquidatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemLiquidatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemLiquidated represents a Liquidated event raised by the TBTCSystem contract.
type TBTCSystemLiquidated struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLiquidated is a free log retrieval operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterLiquidated(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemLiquidatedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "Liquidated", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemLiquidatedIterator{contract: _TBTCSystem.contract, event: "Liquidated", logs: logs, sub: sub}, nil
}

// WatchLiquidated is a free log subscription operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchLiquidated(opts *bind.WatchOpts, sink chan<- *TBTCSystemLiquidated, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "Liquidated", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemLiquidated)
				if err := _TBTCSystem.contract.UnpackLog(event, "Liquidated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLiquidated is a log parse operation binding the contract event 0xa5ee7a2b0254fce91deed604506790ed7fa072d0b14cba4859c3bc8955b9caac.
//
// Solidity: event Liquidated(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseLiquidated(log types.Log) (*TBTCSystemLiquidated, error) {
	event := new(TBTCSystemLiquidated)
	if err := _TBTCSystem.contract.UnpackLog(event, "Liquidated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemLotSizesUpdateStartedIterator is returned from FilterLotSizesUpdateStarted and is used to iterate over the raw logs and unpacked data for LotSizesUpdateStarted events raised by the TBTCSystem contract.
type TBTCSystemLotSizesUpdateStartedIterator struct {
	Event *TBTCSystemLotSizesUpdateStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemLotSizesUpdateStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemLotSizesUpdateStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemLotSizesUpdateStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemLotSizesUpdateStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemLotSizesUpdateStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemLotSizesUpdateStarted represents a LotSizesUpdateStarted event raised by the TBTCSystem contract.
type TBTCSystemLotSizesUpdateStarted struct {
	LotSizes  []uint64
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLotSizesUpdateStarted is a free log retrieval operation binding the contract event 0xffb1e2bce3c7a63d0cac5492540bd370e1156621adde36cd481262c2846a2b7b.
//
// Solidity: event LotSizesUpdateStarted(uint64[] _lotSizes, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterLotSizesUpdateStarted(opts *bind.FilterOpts) (*TBTCSystemLotSizesUpdateStartedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "LotSizesUpdateStarted")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemLotSizesUpdateStartedIterator{contract: _TBTCSystem.contract, event: "LotSizesUpdateStarted", logs: logs, sub: sub}, nil
}

// WatchLotSizesUpdateStarted is a free log subscription operation binding the contract event 0xffb1e2bce3c7a63d0cac5492540bd370e1156621adde36cd481262c2846a2b7b.
//
// Solidity: event LotSizesUpdateStarted(uint64[] _lotSizes, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchLotSizesUpdateStarted(opts *bind.WatchOpts, sink chan<- *TBTCSystemLotSizesUpdateStarted) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "LotSizesUpdateStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemLotSizesUpdateStarted)
				if err := _TBTCSystem.contract.UnpackLog(event, "LotSizesUpdateStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLotSizesUpdateStarted is a log parse operation binding the contract event 0xffb1e2bce3c7a63d0cac5492540bd370e1156621adde36cd481262c2846a2b7b.
//
// Solidity: event LotSizesUpdateStarted(uint64[] _lotSizes, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseLotSizesUpdateStarted(log types.Log) (*TBTCSystemLotSizesUpdateStarted, error) {
	event := new(TBTCSystemLotSizesUpdateStarted)
	if err := _TBTCSystem.contract.UnpackLog(event, "LotSizesUpdateStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemLotSizesUpdatedIterator is returned from FilterLotSizesUpdated and is used to iterate over the raw logs and unpacked data for LotSizesUpdated events raised by the TBTCSystem contract.
type TBTCSystemLotSizesUpdatedIterator struct {
	Event *TBTCSystemLotSizesUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemLotSizesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemLotSizesUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemLotSizesUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemLotSizesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemLotSizesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemLotSizesUpdated represents a LotSizesUpdated event raised by the TBTCSystem contract.
type TBTCSystemLotSizesUpdated struct {
	LotSizes []uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterLotSizesUpdated is a free log retrieval operation binding the contract event 0xa801e49d33e856d89d06e647753d9d9dda3d0b0520c4346290ada455a00cafcc.
//
// Solidity: event LotSizesUpdated(uint64[] _lotSizes)
func (_TBTCSystem *TBTCSystemFilterer) FilterLotSizesUpdated(opts *bind.FilterOpts) (*TBTCSystemLotSizesUpdatedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "LotSizesUpdated")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemLotSizesUpdatedIterator{contract: _TBTCSystem.contract, event: "LotSizesUpdated", logs: logs, sub: sub}, nil
}

// WatchLotSizesUpdated is a free log subscription operation binding the contract event 0xa801e49d33e856d89d06e647753d9d9dda3d0b0520c4346290ada455a00cafcc.
//
// Solidity: event LotSizesUpdated(uint64[] _lotSizes)
func (_TBTCSystem *TBTCSystemFilterer) WatchLotSizesUpdated(opts *bind.WatchOpts, sink chan<- *TBTCSystemLotSizesUpdated) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "LotSizesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemLotSizesUpdated)
				if err := _TBTCSystem.contract.UnpackLog(event, "LotSizesUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLotSizesUpdated is a log parse operation binding the contract event 0xa801e49d33e856d89d06e647753d9d9dda3d0b0520c4346290ada455a00cafcc.
//
// Solidity: event LotSizesUpdated(uint64[] _lotSizes)
func (_TBTCSystem *TBTCSystemFilterer) ParseLotSizesUpdated(log types.Log) (*TBTCSystemLotSizesUpdated, error) {
	event := new(TBTCSystemLotSizesUpdated)
	if err := _TBTCSystem.contract.UnpackLog(event, "LotSizesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TBTCSystem contract.
type TBTCSystemOwnershipTransferredIterator struct {
	Event *TBTCSystemOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemOwnershipTransferred represents a OwnershipTransferred event raised by the TBTCSystem contract.
type TBTCSystemOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TBTCSystem *TBTCSystemFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TBTCSystemOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemOwnershipTransferredIterator{contract: _TBTCSystem.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TBTCSystem *TBTCSystemFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TBTCSystemOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemOwnershipTransferred)
				if err := _TBTCSystem.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TBTCSystem *TBTCSystemFilterer) ParseOwnershipTransferred(log types.Log) (*TBTCSystemOwnershipTransferred, error) {
	event := new(TBTCSystemOwnershipTransferred)
	if err := _TBTCSystem.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the TBTCSystem contract.
type TBTCSystemRedeemedIterator struct {
	Event *TBTCSystemRedeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemRedeemed represents a Redeemed event raised by the TBTCSystem contract.
type TBTCSystemRedeemed struct {
	DepositContractAddress common.Address
	Txid                   [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterRedeemed(opts *bind.FilterOpts, _depositContractAddress []common.Address, _txid [][32]byte) (*TBTCSystemRedeemedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "Redeemed", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemRedeemedIterator{contract: _TBTCSystem.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *TBTCSystemRedeemed, _depositContractAddress []common.Address, _txid [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _txidRule []interface{}
	for _, _txidItem := range _txid {
		_txidRule = append(_txidRule, _txidItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "Redeemed", _depositContractAddressRule, _txidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemRedeemed)
				if err := _TBTCSystem.contract.UnpackLog(event, "Redeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemed is a log parse operation binding the contract event 0x44b7f176bcc739b54bd0800fe491cbdea19df7d4d6b19c281462e6b4fc504344.
//
// Solidity: event Redeemed(address indexed _depositContractAddress, bytes32 indexed _txid, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseRedeemed(log types.Log) (*TBTCSystemRedeemed, error) {
	event := new(TBTCSystemRedeemed)
	if err := _TBTCSystem.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemRedemptionRequestedIterator is returned from FilterRedemptionRequested and is used to iterate over the raw logs and unpacked data for RedemptionRequested events raised by the TBTCSystem contract.
type TBTCSystemRedemptionRequestedIterator struct {
	Event *TBTCSystemRedemptionRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemRedemptionRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemRedemptionRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemRedemptionRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemRedemptionRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemRedemptionRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemRedemptionRequested represents a RedemptionRequested event raised by the TBTCSystem contract.
type TBTCSystemRedemptionRequested struct {
	DepositContractAddress common.Address
	Requester              common.Address
	Digest                 [32]byte
	UtxoValue              *big.Int
	RedeemerOutputScript   []byte
	RequestedFee           *big.Int
	Outpoint               []byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRedemptionRequested is a free log retrieval operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_TBTCSystem *TBTCSystemFilterer) FilterRedemptionRequested(opts *bind.FilterOpts, _depositContractAddress []common.Address, _requester []common.Address, _digest [][32]byte) (*TBTCSystemRedemptionRequestedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "RedemptionRequested", _depositContractAddressRule, _requesterRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemRedemptionRequestedIterator{contract: _TBTCSystem.contract, event: "RedemptionRequested", logs: logs, sub: sub}, nil
}

// WatchRedemptionRequested is a free log subscription operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_TBTCSystem *TBTCSystemFilterer) WatchRedemptionRequested(opts *bind.WatchOpts, sink chan<- *TBTCSystemRedemptionRequested, _depositContractAddress []common.Address, _requester []common.Address, _digest [][32]byte) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}
	var _requesterRule []interface{}
	for _, _requesterItem := range _requester {
		_requesterRule = append(_requesterRule, _requesterItem)
	}
	var _digestRule []interface{}
	for _, _digestItem := range _digest {
		_digestRule = append(_digestRule, _digestItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "RedemptionRequested", _depositContractAddressRule, _requesterRule, _digestRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemRedemptionRequested)
				if err := _TBTCSystem.contract.UnpackLog(event, "RedemptionRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedemptionRequested is a log parse operation binding the contract event 0x7959c380174061a21a3ba80243a032ba9cd10dc8bd1736d7e835c94e97a35a98.
//
// Solidity: event RedemptionRequested(address indexed _depositContractAddress, address indexed _requester, bytes32 indexed _digest, uint256 _utxoValue, bytes _redeemerOutputScript, uint256 _requestedFee, bytes _outpoint)
func (_TBTCSystem *TBTCSystemFilterer) ParseRedemptionRequested(log types.Log) (*TBTCSystemRedemptionRequested, error) {
	event := new(TBTCSystemRedemptionRequested)
	if err := _TBTCSystem.contract.UnpackLog(event, "RedemptionRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemRegisteredPubkeyIterator is returned from FilterRegisteredPubkey and is used to iterate over the raw logs and unpacked data for RegisteredPubkey events raised by the TBTCSystem contract.
type TBTCSystemRegisteredPubkeyIterator struct {
	Event *TBTCSystemRegisteredPubkey // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemRegisteredPubkeyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemRegisteredPubkey)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemRegisteredPubkey)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemRegisteredPubkeyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemRegisteredPubkeyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemRegisteredPubkey represents a RegisteredPubkey event raised by the TBTCSystem contract.
type TBTCSystemRegisteredPubkey struct {
	DepositContractAddress common.Address
	SigningGroupPubkeyX    [32]byte
	SigningGroupPubkeyY    [32]byte
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRegisteredPubkey is a free log retrieval operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterRegisteredPubkey(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemRegisteredPubkeyIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "RegisteredPubkey", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemRegisteredPubkeyIterator{contract: _TBTCSystem.contract, event: "RegisteredPubkey", logs: logs, sub: sub}, nil
}

// WatchRegisteredPubkey is a free log subscription operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchRegisteredPubkey(opts *bind.WatchOpts, sink chan<- *TBTCSystemRegisteredPubkey, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "RegisteredPubkey", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemRegisteredPubkey)
				if err := _TBTCSystem.contract.UnpackLog(event, "RegisteredPubkey", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisteredPubkey is a log parse operation binding the contract event 0x8ee737ab16909c4e9d1b750814a4393c9f84ab5d3a29c08c313b783fc846ae33.
//
// Solidity: event RegisteredPubkey(address indexed _depositContractAddress, bytes32 _signingGroupPubkeyX, bytes32 _signingGroupPubkeyY, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseRegisteredPubkey(log types.Log) (*TBTCSystemRegisteredPubkey, error) {
	event := new(TBTCSystemRegisteredPubkey)
	if err := _TBTCSystem.contract.UnpackLog(event, "RegisteredPubkey", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemSetupFailedIterator is returned from FilterSetupFailed and is used to iterate over the raw logs and unpacked data for SetupFailed events raised by the TBTCSystem contract.
type TBTCSystemSetupFailedIterator struct {
	Event *TBTCSystemSetupFailed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemSetupFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemSetupFailed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemSetupFailed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemSetupFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemSetupFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemSetupFailed represents a SetupFailed event raised by the TBTCSystem contract.
type TBTCSystemSetupFailed struct {
	DepositContractAddress common.Address
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetupFailed is a free log retrieval operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterSetupFailed(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemSetupFailedIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "SetupFailed", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemSetupFailedIterator{contract: _TBTCSystem.contract, event: "SetupFailed", logs: logs, sub: sub}, nil
}

// WatchSetupFailed is a free log subscription operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchSetupFailed(opts *bind.WatchOpts, sink chan<- *TBTCSystemSetupFailed, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "SetupFailed", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemSetupFailed)
				if err := _TBTCSystem.contract.UnpackLog(event, "SetupFailed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetupFailed is a log parse operation binding the contract event 0x8fd2cfb62a35fccc1ecef829f83a6c2f840b73dad49d3eaaa402909752086d4b.
//
// Solidity: event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseSetupFailed(log types.Log) (*TBTCSystemSetupFailed, error) {
	event := new(TBTCSystemSetupFailed)
	if err := _TBTCSystem.contract.UnpackLog(event, "SetupFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemSignerFeeDivisorUpdateStartedIterator is returned from FilterSignerFeeDivisorUpdateStarted and is used to iterate over the raw logs and unpacked data for SignerFeeDivisorUpdateStarted events raised by the TBTCSystem contract.
type TBTCSystemSignerFeeDivisorUpdateStartedIterator struct {
	Event *TBTCSystemSignerFeeDivisorUpdateStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemSignerFeeDivisorUpdateStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemSignerFeeDivisorUpdateStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemSignerFeeDivisorUpdateStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemSignerFeeDivisorUpdateStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemSignerFeeDivisorUpdateStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemSignerFeeDivisorUpdateStarted represents a SignerFeeDivisorUpdateStarted event raised by the TBTCSystem contract.
type TBTCSystemSignerFeeDivisorUpdateStarted struct {
	SignerFeeDivisor uint16
	Timestamp        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignerFeeDivisorUpdateStarted is a free log retrieval operation binding the contract event 0x38cb7049f0daf658ca989e9ef6b850ef11e3740ff07a0c16706042c39adf48fc.
//
// Solidity: event SignerFeeDivisorUpdateStarted(uint16 _signerFeeDivisor, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterSignerFeeDivisorUpdateStarted(opts *bind.FilterOpts) (*TBTCSystemSignerFeeDivisorUpdateStartedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "SignerFeeDivisorUpdateStarted")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemSignerFeeDivisorUpdateStartedIterator{contract: _TBTCSystem.contract, event: "SignerFeeDivisorUpdateStarted", logs: logs, sub: sub}, nil
}

// WatchSignerFeeDivisorUpdateStarted is a free log subscription operation binding the contract event 0x38cb7049f0daf658ca989e9ef6b850ef11e3740ff07a0c16706042c39adf48fc.
//
// Solidity: event SignerFeeDivisorUpdateStarted(uint16 _signerFeeDivisor, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchSignerFeeDivisorUpdateStarted(opts *bind.WatchOpts, sink chan<- *TBTCSystemSignerFeeDivisorUpdateStarted) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "SignerFeeDivisorUpdateStarted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemSignerFeeDivisorUpdateStarted)
				if err := _TBTCSystem.contract.UnpackLog(event, "SignerFeeDivisorUpdateStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignerFeeDivisorUpdateStarted is a log parse operation binding the contract event 0x38cb7049f0daf658ca989e9ef6b850ef11e3740ff07a0c16706042c39adf48fc.
//
// Solidity: event SignerFeeDivisorUpdateStarted(uint16 _signerFeeDivisor, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseSignerFeeDivisorUpdateStarted(log types.Log) (*TBTCSystemSignerFeeDivisorUpdateStarted, error) {
	event := new(TBTCSystemSignerFeeDivisorUpdateStarted)
	if err := _TBTCSystem.contract.UnpackLog(event, "SignerFeeDivisorUpdateStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemSignerFeeDivisorUpdatedIterator is returned from FilterSignerFeeDivisorUpdated and is used to iterate over the raw logs and unpacked data for SignerFeeDivisorUpdated events raised by the TBTCSystem contract.
type TBTCSystemSignerFeeDivisorUpdatedIterator struct {
	Event *TBTCSystemSignerFeeDivisorUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemSignerFeeDivisorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemSignerFeeDivisorUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemSignerFeeDivisorUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemSignerFeeDivisorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemSignerFeeDivisorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemSignerFeeDivisorUpdated represents a SignerFeeDivisorUpdated event raised by the TBTCSystem contract.
type TBTCSystemSignerFeeDivisorUpdated struct {
	SignerFeeDivisor uint16
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignerFeeDivisorUpdated is a free log retrieval operation binding the contract event 0x236dec26d53c6a51390e98ed703106e132fd062b0e38b8a9cf8b4d13f47952c8.
//
// Solidity: event SignerFeeDivisorUpdated(uint16 _signerFeeDivisor)
func (_TBTCSystem *TBTCSystemFilterer) FilterSignerFeeDivisorUpdated(opts *bind.FilterOpts) (*TBTCSystemSignerFeeDivisorUpdatedIterator, error) {

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "SignerFeeDivisorUpdated")
	if err != nil {
		return nil, err
	}
	return &TBTCSystemSignerFeeDivisorUpdatedIterator{contract: _TBTCSystem.contract, event: "SignerFeeDivisorUpdated", logs: logs, sub: sub}, nil
}

// WatchSignerFeeDivisorUpdated is a free log subscription operation binding the contract event 0x236dec26d53c6a51390e98ed703106e132fd062b0e38b8a9cf8b4d13f47952c8.
//
// Solidity: event SignerFeeDivisorUpdated(uint16 _signerFeeDivisor)
func (_TBTCSystem *TBTCSystemFilterer) WatchSignerFeeDivisorUpdated(opts *bind.WatchOpts, sink chan<- *TBTCSystemSignerFeeDivisorUpdated) (event.Subscription, error) {

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "SignerFeeDivisorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemSignerFeeDivisorUpdated)
				if err := _TBTCSystem.contract.UnpackLog(event, "SignerFeeDivisorUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignerFeeDivisorUpdated is a log parse operation binding the contract event 0x236dec26d53c6a51390e98ed703106e132fd062b0e38b8a9cf8b4d13f47952c8.
//
// Solidity: event SignerFeeDivisorUpdated(uint16 _signerFeeDivisor)
func (_TBTCSystem *TBTCSystemFilterer) ParseSignerFeeDivisorUpdated(log types.Log) (*TBTCSystemSignerFeeDivisorUpdated, error) {
	event := new(TBTCSystemSignerFeeDivisorUpdated)
	if err := _TBTCSystem.contract.UnpackLog(event, "SignerFeeDivisorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TBTCSystemStartedLiquidationIterator is returned from FilterStartedLiquidation and is used to iterate over the raw logs and unpacked data for StartedLiquidation events raised by the TBTCSystem contract.
type TBTCSystemStartedLiquidationIterator struct {
	Event *TBTCSystemStartedLiquidation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TBTCSystemStartedLiquidationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TBTCSystemStartedLiquidation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TBTCSystemStartedLiquidation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TBTCSystemStartedLiquidationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TBTCSystemStartedLiquidationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TBTCSystemStartedLiquidation represents a StartedLiquidation event raised by the TBTCSystem contract.
type TBTCSystemStartedLiquidation struct {
	DepositContractAddress common.Address
	WasFraud               bool
	Timestamp              *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterStartedLiquidation is a free log retrieval operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) FilterStartedLiquidation(opts *bind.FilterOpts, _depositContractAddress []common.Address) (*TBTCSystemStartedLiquidationIterator, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.FilterLogs(opts, "StartedLiquidation", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return &TBTCSystemStartedLiquidationIterator{contract: _TBTCSystem.contract, event: "StartedLiquidation", logs: logs, sub: sub}, nil
}

// WatchStartedLiquidation is a free log subscription operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) WatchStartedLiquidation(opts *bind.WatchOpts, sink chan<- *TBTCSystemStartedLiquidation, _depositContractAddress []common.Address) (event.Subscription, error) {

	var _depositContractAddressRule []interface{}
	for _, _depositContractAddressItem := range _depositContractAddress {
		_depositContractAddressRule = append(_depositContractAddressRule, _depositContractAddressItem)
	}

	logs, sub, err := _TBTCSystem.contract.WatchLogs(opts, "StartedLiquidation", _depositContractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TBTCSystemStartedLiquidation)
				if err := _TBTCSystem.contract.UnpackLog(event, "StartedLiquidation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStartedLiquidation is a log parse operation binding the contract event 0xbef11c059eefba82a15aea8a3a89c86fd08d7711c88fa7daea2632a55488510c.
//
// Solidity: event StartedLiquidation(address indexed _depositContractAddress, bool _wasFraud, uint256 _timestamp)
func (_TBTCSystem *TBTCSystemFilterer) ParseStartedLiquidation(log types.Log) (*TBTCSystemStartedLiquidation, error) {
	event := new(TBTCSystemStartedLiquidation)
	if err := _TBTCSystem.contract.UnpackLog(event, "StartedLiquidation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
