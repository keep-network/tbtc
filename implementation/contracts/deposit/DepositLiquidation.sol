pragma solidity 0.4.25;

import {SafeMath} from "../bitcoin-spv/SafeMath.sol";
import {BTCUtils} from "../bitcoin-spv/BTCUtils.sol";
import {DepositStates} from "./DepositStates.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {TBTCConstants} from "./TBTCConstants.sol";
import {IKeep} from "../interfaces/IKeep.sol";
import {OutsourceDepositLogging} from "./OutsourceDepositLogging.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {TBTCSystem} from "../system/TBTCSystem.sol";

library DepositLiquidation {

    using BTCUtils for bytes;
    using SafeMath for uint256;

    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;
    using OutsourceDepositLogging for DepositUtils.Deposit;

    /// @notice     Tries to liquidate the position on-chain using the signer bond
    /// @dev        Calls out to other contracts, watch for re-entrance
    /// @return     True if Liquidated, False otherwise
    function attemptToLiquidateOnchain() public pure returns (bool) {
        return false;
    }

    /// @notice                 Notifies the keep contract of fraud
    /// @dev                    Calls out to the keep contract. this could get expensive if preimage is large
    /// @param  _d               deposit storage pointer
    /// @param  _v              Signature recovery value
    /// @param  _r              Signature R value
    /// @param  _s              Signature S value
    /// @param _signedDigest    The digest signed by the signature vrs tuple
    /// @param _preimage        The sha256 preimage of the digest
    /// @return                 True if fraud, otherwise revert
    function submitSignatureFraud(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) public returns (bool _isFraud) {
        IKeep _keep = IKeep(_d.KeepBridge);
        return _keep.submitSignatureFraud(_d.keepID, _v, _r, _s, _signedDigest, _preimage);
    }

    /// @notice     Determines the collateralization percentage of the signing group
    /// @dev        Compares the bond value and lot value
    /// @param _d   deposit storage pointer
    /// @return     collateralization percentage as uint
    function getCollateralizationPercentage(DepositUtils.Deposit storage _d) public view returns (uint256) {

        // Determine value of the lot in wei
        uint256 _oraclePrice = _d.fetchOraclePrice();
        if (_oraclePrice == 0 || _oraclePrice > 10 ** 18) {
            /*
              This is if a sat is worth 0 wei, or is worth 1 ether
              TODO: what should this behavior be?
            */
            revert("Oracle returned a bad price");
        }

        uint256 _lotSize = TBTCConstants.getLotSize();
        uint256 _lotValue = _lotSize * _oraclePrice;

        // Amount of wei the signers have
        uint256 _bondValue = _d.fetchBondAmount();

        // This converts into a percentage
        return (_bondValue.mul(100).div(_lotValue));
    }

    /// @notice         Starts signer liquidation due to fraud
    /// @dev            We first attempt to liquidate on chain, then by auction
    /// @param  _d      deposit storage pointer
    function startSignerFraudLiquidation(DepositUtils.Deposit storage _d) public {
        _d.logStartedLiquidation(true);

        // Reclaim used state for gas savings
        _d.redemptionTeardown();
        uint256 _seized = _d.seizeSignerBonds();

        if (_d.auctionTBTCAmount() == 0) {
            // we came from the redemption flow
            _d.setLiquidated();
            _d.requesterAddress.transfer(_seized);
            _d.logLiquidated();
            return;
        }

        bool _liquidated = attemptToLiquidateOnchain();

        if (_liquidated) {
            _d.distributeBeneficiaryReward();
            _d.setLiquidated();
            _d.logLiquidated();
            address(0).transfer(address(this).balance);  // burn it down
        }
        if (!_liquidated) {
            _d.setFraudLiquidationInProgress();
            _d.liquidationInitiated = block.timestamp;  // Store the timestamp for auction
        }
    }

    /// @notice         Starts signer liquidation due to abort or undercollateralization
    /// @dev            We first attempt to liquidate on chain, then by auction
    /// @param  _d      deposit storage pointer
    function startSignerAbortLiquidation(DepositUtils.Deposit storage _d) public {
        _d.logStartedLiquidation(false);
        // Reclaim used state for gas savings
        _d.redemptionTeardown();
        _d.seizeSignerBonds();

        bool _liquidated = attemptToLiquidateOnchain();

        if (_liquidated) {
            _d.distributeBeneficiaryReward();
            _d.pushFundsToKeepGroup(address(this).balance);
            _d.setLiquidated();
            _d.logLiquidated();
        }
        if (!_liquidated) {
            _d.liquidationInitiated = block.timestamp;  // Store the timestamp for auction
            _d.setFraudLiquidationInProgress();
        }
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud
    /// @dev                    ECDSA is NOT SECURE unless you verify the digest
    /// @param  _d              deposit storage pointer
    /// @param  _v              Signature recovery value
    /// @param  _r              Signature R value
    /// @param  _s              Signature S value
    /// @param _signedDigest    The digest signed by the signature vrs tuple
    /// @param _preimage        The sha256 preimage of the digest
    function provideECDSAFraudProof(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
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

    /// @notice                 Anyone may notify the deposit of fraud via an SPV proof
    /// @dev                    We strong prefer ECDSA fraud proofs
    /// @param  _d              deposit storage pointer
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contains the funding output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    function provideSPVFraudProof(
        DepositUtils.Deposit storage _d,
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) public {
        bytes memory _input;
        bytes memory _output;
        bool _inputConsumed;
        uint8 i;

        require(
            !_d.inFunding() && !_d.inFundingFailure(),
            "SPV Fraud proofs not valid before Active state."
        );
        require(
            !_d.inSignerLiquidation(),
            "Signer liquidation already in progress"
        );
        require(!_d.inEndState(), "Contract has halted");

        _d.checkProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);
        for (i = 0; i < _bitcoinTx.extractNumInputs(); i++) {
            _input = _bitcoinTx.extractInputAtIndex(i);
            if (keccak256(_input.extractOutpoint()) == keccak256(_d.utxoOutpoint)) {
                _inputConsumed = true;
            }
        }
        require(_inputConsumed, "No input spending custodied UTXO found");

        uint256 _permittedFeeBumps = 5;  /* TODO: can we refactor withdrawal flow to improve this? */
        uint256 _requiredOutputSize = _d.utxoSize().sub((_d.initialRedemptionFee * (1 + _permittedFeeBumps)));
        for (i = 0; i < _bitcoinTx.extractNumOutputs(); i++) {
            _output = _bitcoinTx.extractOutputAtIndex(i);

            if (_output.extractValue() >= _requiredOutputSize
                && keccak256(_output.extractHash()) == keccak256(abi.encodePacked(_d.requesterPKH))) {
                revert("Found an output paying the redeemer as requested");
            }
        }

        startSignerFraudLiquidation(_d);
    }

    /// @notice     Closes an auction and purchases the signer bonds. Payout to buyer, funder, then signers if not fraud
    /// @dev        For interface, reading auctionValue will give a past value. the current is better
    /// @param  _d  deposit storage pointer
    function purchaseSignerBondsAtAuction(DepositUtils.Deposit storage _d) public {
        bool _wasFraud = _d.inFraudLiquidationInProgress();
        require(_d.inSignerLiquidation(), "No active auction");

        _d.setLiquidated();
        _d.logLiquidated();

        // Burn the outstanding TBTC
        TBTCToken _tbtc = TBTCToken(_d.TBTCToken);
        TBTCSystem _system = TBTCSystem(_d.TBTCSystem);
        require(_tbtc.balanceOf(msg.sender) >= TBTCConstants.getLotSize(), "Not enough TBTC to cover outstanding debt");
        _system.systemBurnFrom(msg.sender, TBTCConstants.getLotSize());  // burn minimal amount to cover size

        // Distribute funds to auction buyer
        uint256 _valueToDistribute = _d.auctionValue();
        msg.sender.transfer(_valueToDistribute);

        // Send any TBTC left to the beneficiary
        _d.distributeBeneficiaryReward();

        // then if there are funds left, and it wasn't fraud, pay out the signers
        if (address(this).balance > 0) {
            if (_wasFraud) {
                // Burn it
                address(0).transfer(address(this).balance);
            } else {
                // Send it back
                _d.pushFundsToKeepGroup(address(this).balance);
            }
        }
    }

    /// @notice     Notify the contract that the signers are undercollateralized
    /// @dev        Calls out to the system for oracle info
    /// @param  _d  deposit storage pointer
    function notifyCourtesyCall(DepositUtils.Deposit storage _d) public  {
        require(_d.inActive(), "Can only courtesy call from active state");
        require(getCollateralizationPercentage(_d) < TBTCConstants.getUndercollateralizedPercent(), "Signers have sufficient collateral");
        _d.courtesyCallInitiated = block.timestamp;
        _d.setCourtesyCall();
        _d.logCourtesyCalled();
    }

    /// @notice     Goes from courtesy call to active
    /// @dev        Only callable if collateral is sufficient and the deposit is not expiring
    /// @param  _d  deposit storage pointer
    function exitCourtesyCall(DepositUtils.Deposit storage _d) public {
        require(_d.inCourtesyCall(), "Not currently in courtesy call");
        require(block.timestamp <= _d.fundedAt + TBTCConstants.getDepositTerm(), "Deposit is expiring");
        require(getCollateralizationPercentage(_d) >= TBTCConstants.getUndercollateralizedPercent(), "Deposit is still undercollateralized");
        _d.setActive();
        _d.logExitedCourtesyCall();
    }

    /// @notice     Notify the contract that the signers are undercollateralized
    /// @dev        Calls out to the system for oracle info
    /// @param  _d  deposit storage pointer
    function notifyUndercollateralizedLiquidation(DepositUtils.Deposit storage _d) public {
        require(_d.inRedeemableState(), "Deposit not in active or courtesy call");
        require(getCollateralizationPercentage(_d) < TBTCConstants.getSeverelyUndercollateralizedPercent(), "Deposit has sufficient collateral");
        startSignerAbortLiquidation(_d);
    }

    /// @notice     Notifies the contract that the courtesy period has elapsed
    /// @dev        This is treated as an abort, rather than fraud
    /// @param  _d  deposit storage pointer
    function notifyCourtesyTimeout(DepositUtils.Deposit storage _d) public {
        require(_d.inCourtesyCall(), "Not in a courtesy call period");
        require(block.timestamp >= _d.courtesyCallInitiated + TBTCConstants.getCourtesyCallTimeout(), "Courtesy period has not elapsed");
        startSignerAbortLiquidation(_d);
    }

    /// @notice     Notifies the contract that its term limit has been reached
    /// @dev        This initiates a courtesy call
    /// @param  _d  deposit storage pointer
    function notifyDepositExpiryCourtesyCall(DepositUtils.Deposit storage _d) public {
        require(_d.inActive(), "Deposit is not active");
        require(block.timestamp >= _d.fundedAt + TBTCConstants.getDepositTerm(), "Deposit term not elapsed");
        _d.setCourtesyCall();
        _d.logCourtesyCalled();
        _d.courtesyCallInitiated = block.timestamp;
    }
}
