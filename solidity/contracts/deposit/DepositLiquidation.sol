pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import {BTCUtils} from "@summa-tx/bitcoin-spv-sol/contracts/BTCUtils.sol";
import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";
import {DepositStates} from "./DepositStates.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {TBTCConstants} from "./TBTCConstants.sol";
import {IBondedECDSAKeep} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeep.sol";
import {OutsourceDepositLogging} from "./OutsourceDepositLogging.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";

library DepositLiquidation {

    using BTCUtils for bytes;
    using BytesLib for bytes;
    using SafeMath for uint256;

    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;
    using OutsourceDepositLogging for DepositUtils.Deposit;

    /// @notice                 Notifies the keep contract of fraud.
    /// @dev                    Calls out to the keep contract. this could get expensive if preimage is large.
    /// @param  _d              Deposit storage pointer.
    /// @param  _v              Signature recovery value.
    /// @param  _r              Signature R value.
    /// @param  _s              Signature S value.
    /// @param _signedDigest    The digest signed by the signature vrs tuple.
    /// @param _preimage        The sha256 preimage of the digest.
    /// @return                 True if fraud, otherwise revert.
    function submitSignatureFraud(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes memory _preimage
    ) public returns (bool _isFraud) {
        IBondedECDSAKeep _keep = IBondedECDSAKeep(_d.keepAddress);
        return _keep.submitSignatureFraud(_v, _r, _s, _signedDigest, _preimage);
    }

    /// @notice     Determines the collateralization percentage of the signing group.
    /// @dev        Compares the bond value and lot value.
    /// @param _d   Deposit storage pointer.
    /// @return     Collateralization percentage as uint.
    function getCollateralizationPercentage(DepositUtils.Deposit storage _d) public view returns (uint256) {

        // Determine value of the lot in wei
        uint256 _satoshiPrice = _d.fetchBitcoinPrice();
        uint256 _lotSizeSatoshis = _d.lotSizeSatoshis;
        uint256 _lotValue = _lotSizeSatoshis.mul(_satoshiPrice);

        // Amount of wei the signers have
        uint256 _bondValue = _d.fetchBondAmount();

        // This converts into a percentage
        return (_bondValue.mul(100).div(_lotValue));
    }

    /// @notice         Starts signer liquidation due to fraud.
    /// @dev            Liquidation is done by auction.
    /// @param  _d      Deposit storage pointer.
    function startSignerFraudLiquidation(DepositUtils.Deposit storage _d) internal {
        _d.logStartedLiquidation(true);

        // Reclaim used state for gas savings
        _d.redemptionTeardown();
        uint256 _seized = _d.seizeSignerBonds();

        if (_d.auctionTBTCAmount() == 0) {
            // we came from the redemption flow
            _d.setLiquidated();
            _d.redeemerAddress.transfer(_seized);
            _d.logLiquidated();
            return;
        }

        _d.liquidationInitiator = msg.sender;
        _d.liquidationInitiated = block.timestamp;  // Store the timestamp for auction

        _d.setFraudLiquidationInProgress();
    }

    /// @notice         Starts signer liquidation due to abort or undercollateralization.
    /// @dev            Liquidation is done by auction.
    /// @param  _d      Deposit storage pointer.
    function startSignerAbortLiquidation(DepositUtils.Deposit storage _d) internal {
        _d.logStartedLiquidation(false);
        // Reclaim used state for gas savings
        _d.redemptionTeardown();
        _d.seizeSignerBonds();

        _d.liquidationInitiated = block.timestamp;  // Store the timestamp for auction
        _d.liquidationInitiator = msg.sender;
        _d.setLiquidationInProgress();
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud.
    /// @dev                    Calls out to the keep to verify if there was fraud.
    /// @param  _d              Deposit storage pointer.
    /// @param  _v              Signature recovery value.
    /// @param  _r              Signature R value.
    /// @param  _s              Signature S value.
    /// @param _signedDigest    The digest signed by the signature vrs tuple.
    /// @param _preimage        The sha256 preimage of the digest.
    function provideECDSAFraudProof(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes memory _preimage
    ) public {
        require(
            !_d.inFunding() && !_d.inFundingFailure(),
            "Use provideFundingECDSAFraudProof instead"
        );
        require(
            !_d.inSignerLiquidation(),
            "Signer liquidation already in progress"
        );
        require(!_d.inEndState(), "Contract has halted");
        require(submitSignatureFraud(_d, _v, _r, _s, _signedDigest, _preimage), "Signature is not fraud");
        startSignerFraudLiquidation(_d);
    }

    /// @notice                 Search _txOutputVector for output paying the redeemer.
    /// @dev                    Require that outputs checked are witness.
    /// @param  _d              Deposit storage pointer.
    /// @param _txOutputVector  All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs.
    /// @return                 False if output paying redeemer was found, true otherwise.
    function validateRedeemerNotPaid(
        DepositUtils.Deposit storage _d,
        bytes memory _txOutputVector
    ) internal view returns (bool){
        bytes memory _output;
        uint256 _offset = 1;
        uint256 _permittedFeeBumps = TBTCConstants.getPermittedFeeBumps();  /* TODO: can we refactor withdrawal flow to improve this? */
        uint256 _requiredOutputValue = _d.utxoSize().sub((_d.initialRedemptionFee.mul((_permittedFeeBumps.add(1)))));

        uint8 _numOuts = uint8(_txOutputVector.slice(0, 1)[0]);
        for (uint8 i = 0; i < _numOuts; i++) {
            _output = _txOutputVector.slice(_offset, _txOutputVector.length.sub(_offset));
            _offset = _offset.add(_output.determineOutputLength());

            if (_output.extractValue() >= _requiredOutputValue &&
                keccak256(_output.slice(8, 3).concat(_output.extractHash())) == keccak256(abi.encodePacked(_d.redeemerOutputScript))) {
                return false;
            }
        }
        return true;
    }

    /// @notice     Closes an auction and purchases the signer bonds. Payout to buyer, funder, then signers if not fraud.
    /// @dev        For interface, reading auctionValue will give a past value. the current is better.
    /// @param  _d  Deposit storage pointer.
    function purchaseSignerBondsAtAuction(DepositUtils.Deposit storage _d) public {
        bool _wasFraud = _d.inFraudLiquidationInProgress();
        require(_d.inSignerLiquidation(), "No active auction");

        _d.setLiquidated();
        _d.logLiquidated();

        // send the TBTC to the TDT holder. If the TDT holder is the Vending Machine, burn it to maintain the peg.
        address tdtHolder = _d.depositOwner();

        uint256 lotSizeTbtc = _d.lotSizeTbtc();
        require(_d.tbtcToken.balanceOf(msg.sender) >= lotSizeTbtc, "Not enough TBTC to cover outstanding debt");

        if(tdtHolder == _d.VendingMachineAddress){
            _d.tbtcToken.burnFrom(msg.sender, lotSizeTbtc);  // burn minimal amount to cover size
        }
        else{
            _d.tbtcToken.transferFrom(msg.sender, tdtHolder, lotSizeTbtc);
        }

        // Distribute funds to auction buyer
        uint256 _valueToDistribute = _d.auctionValue();
        msg.sender.transfer(_valueToDistribute);

        // Send any TBTC left to the Fee Rebate Token holder
        _d.distributeFeeRebate();

        // For fraud, pay remainder to the liquidation initiator.
        // For non-fraud, split 50-50 between initiator and signers. if the transfer amount is 1,
        // division will yield a 0 value which causes a revert; instead,
        // we simply ignore such a tiny amount and leave some wei dust in escrow
        uint256 contractEthBalance = address(this).balance;
        address payable initiator = _d.liquidationInitiator;

        if (initiator == address(0)){
            initiator = address(0xdead);
        }
        if (contractEthBalance > 1) {
            if (_wasFraud) {
                /* solium-disable-next-line security/no-send */
                initiator.send(contractEthBalance);
            } else {
                // There will always be a liquidation initiator.
                uint256 split = contractEthBalance.div(2);
                _d.pushFundsToKeepGroup(split);
                /* solium-disable-next-line security/no-send */
                initiator.send(address(this).balance);
            }
        }
    }

    /// @notice     Notify the contract that the signers are undercollateralized.
    /// @dev        Calls out to the system for oracle info.
    /// @param  _d  Deposit storage pointer.
    function notifyCourtesyCall(DepositUtils.Deposit storage _d) public  {
        require(_d.inActive(), "Can only courtesy call from active state");
        require(getCollateralizationPercentage(_d) < _d.undercollateralizedThresholdPercent, "Signers have sufficient collateral");
        _d.courtesyCallInitiated = block.timestamp;
        _d.setCourtesyCall();
        _d.logCourtesyCalled();
    }

    /// @notice     Goes from courtesy call to active.
    /// @dev        Only callable if collateral is sufficient and the deposit is not expiring.
    /// @param  _d  Deposit storage pointer.
    function exitCourtesyCall(DepositUtils.Deposit storage _d) public {
        require(_d.inCourtesyCall(), "Not currently in courtesy call");
        require(getCollateralizationPercentage(_d) >= _d.undercollateralizedThresholdPercent, "Deposit is still undercollateralized");
        _d.setActive();
        _d.logExitedCourtesyCall();
    }

    /// @notice     Notify the contract that the signers are undercollateralized.
    /// @dev        Calls out to the system for oracle info.
    /// @param  _d  Deposit storage pointer.
    function notifyUndercollateralizedLiquidation(DepositUtils.Deposit storage _d) public {
        require(_d.inRedeemableState(), "Deposit not in active or courtesy call");
        require(getCollateralizationPercentage(_d) < _d.severelyUndercollateralizedThresholdPercent, "Deposit has sufficient collateral");
        startSignerAbortLiquidation(_d);
    }

    /// @notice     Notifies the contract that the courtesy period has elapsed.
    /// @dev        This is treated as an abort, rather than fraud.
    /// @param  _d  Deposit storage pointer.
    function notifyCourtesyTimeout(DepositUtils.Deposit storage _d) public {
        require(_d.inCourtesyCall(), "Not in a courtesy call period");
        require(block.timestamp >= _d.courtesyCallInitiated.add(TBTCConstants.getCourtesyCallTimeout()), "Courtesy period has not elapsed");
        startSignerAbortLiquidation(_d);
    }
}
