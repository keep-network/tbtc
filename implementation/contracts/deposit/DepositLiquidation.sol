pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import {BTCUtils} from "@summa-tx/bitcoin-spv-sol/contracts/BTCUtils.sol";
import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";
import {DepositStates} from "./DepositStates.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {TBTCConstants} from "./TBTCConstants.sol";
import {IBondedECDSAKeep} from "../external/IBondedECDSAKeep.sol";
import {OutsourceDepositLogging} from "./OutsourceDepositLogging.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {IUniswapExchange} from "../external/IUniswapExchange.sol";
import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";

library DepositLiquidation {

    using BTCUtils for bytes;
    using BytesLib for bytes;
    using SafeMath for uint256;

    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;
    using OutsourceDepositLogging for DepositUtils.Deposit;

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
        bytes memory _preimage
    ) public returns (bool _isFraud) {
        IBondedECDSAKeep _keep = IBondedECDSAKeep(_d.keepAddress);
        return _keep.submitSignatureFraud(_d.keepAddress, _v, _r, _s, _signedDigest, _preimage);
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

        bool _liquidated = attemptToLiquidateOnchain(_d);

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

        bool _liquidated = attemptToLiquidateOnchain(_d);

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

    /// @notice                   Anyone may notify the deposit of fraud via an SPV proof
    /// @dev                      We strong prefer ECDSA fraud proofs
    /// @param  _d                deposit storage pointer
    /// @param  _txVersion        Transaction version number (4-byte LE)
    /// @param  _txInputVector    All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param  _txOutputVector   All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param  _txLocktime       Final 4 bytes of the transaction
    /// @param  _merkleProof      The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _txIndexInBlock   Transaction index in the block (0-indexed)
    /// @param  _targetInputIndex Index of the input that spends the custodied UTXO
    /// @param  _bitcoinHeaders   An array of tightly-packed bitcoin headers
    function provideSPVFraudProof(
        DepositUtils.Deposit storage _d,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        uint8 _targetInputIndex,
        bytes memory _bitcoinHeaders
    ) public {
        bytes32 _txId;
        require(
            !_d.inFunding() && !_d.inFundingFailure(),
            "SPV Fraud proofs not valid before Active state."
        );
        require(
            !_d.inSignerLiquidation(),
            "Signer liquidation already in progress"
        );
        require(!_d.inEndState(), "Contract has halted");
        require(_txInputVector.validateVin(), "invalid input vector provided");
        require(_txOutputVector.validateVout(), "invalid output vector provided");

        // DRAFT comments:
        // we can assume this TX is witness?

        _txId = abi.encodePacked(_txVersion, _txInputVector, _txOutputVector, _txLocktime).hash256();

        _d.checkProofFromTxId(_txId, _merkleProof, _txIndexInBlock, _bitcoinHeaders);

        bytes memory _targetOutpoint = _txInputVector.extractInputAtIndex(_targetInputIndex).extractOutpoint();
        require(
            keccak256(_targetOutpoint) == keccak256(_d.utxoOutpoint),
            "No input spending custodied UTXO found at given index"
        );

        if (_d.requesterPKH != bytes20(0)) {
            require(
                validateRedeemerNotPaid(_d, _txOutputVector),
                "Found an output paying the redeemer as requested"
            );
        }

        startSignerFraudLiquidation(_d);
    }

    /// @notice                 Search _txOutputVector for output paying the requestor
    /// @dev                    Require that outputs checked are witness
    /// @param  _d              Deposit storage pointer
    /// @param _txOutputVector  All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @return                 False if output paying redeemer was found, true otherwise
    function validateRedeemerNotPaid(
        DepositUtils.Deposit storage _d,
        bytes memory _txOutputVector
    ) internal view returns (bool){
        bytes memory _output;
        uint256 _offset = 1;
        uint256 _permittedFeeBumps = TBTCConstants.getPermittedFeeBumps();  /* TODO: can we refactor withdrawal flow to improve this? */
        uint256 _requiredOutputValue = _d.utxoSize().sub((_d.initialRedemptionFee * (1 + _permittedFeeBumps)));

        uint8 _numOuts = uint8(_txOutputVector.slice(0, 1)[0]);
        for (uint8 i = 0; i < _numOuts; i++) {
            _output = _txOutputVector.slice(_offset, _txOutputVector.length - _offset);
            _offset += _output.determineOutputLength();

            if (_output.extractValue() >= _requiredOutputValue
                // extract the output flag and check that it is witness
                && keccak256(_output.slice(8, 3)) == keccak256(hex"160014")
                && keccak256(_output.extractHash()) == keccak256(abi.encodePacked(_d.requesterPKH))) {
                return false;
            }
        }
        return true;
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
        TBTCToken _tbtcToken = TBTCToken(_d.TBTCToken);
        require(_tbtcToken.balanceOf(msg.sender) >= TBTCConstants.getLotSize(), "Not enough TBTC to cover outstanding debt");
        _tbtcToken.burnFrom(msg.sender, TBTCConstants.getLotSize());  // burn minimal amount to cover size

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

    /// @notice     Tries to liquidate the position on-chain using the signer bond
    /// @dev        Calls out to other contracts, watch for re-entrance
    /// @return     True if Liquidated, False otherwise
    // TODO(liamz): check for re-entry
    function attemptToLiquidateOnchain(
        DepositUtils.Deposit storage _d
    ) internal returns (bool) {
        // Return early if there is no Uniswap TBTC Exchange.
        IUniswapExchange exchange = IUniswapExchange(ITBTCSystem(_d.TBTCSystem).getTBTCUniswapExchange());
        if(address(exchange) == address(0x0)) {
            return false;
        }

        // Only liquidate if we can buy up enough TBTC to burn,
        // otherwise go 100% for the falling-price auction
        uint tbtcAmount = _d.liquidationTBTCAmount();
        uint ethAmount = exchange.getEthToTokenOutputPrice(tbtcAmount);

        if(address(this).balance < ethAmount) {
            return false;
        }

        // Leverage uniswap’s frontrunning mitigation functionality.
        uint deadline = block.timestamp;
        exchange.ethToTokenSwapOutput.value(ethAmount)(tbtcAmount, deadline);

        return true;
    }
}
