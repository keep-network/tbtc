pragma solidity 0.5.17;

import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";
import {BTCUtils} from "@summa-tx/bitcoin-spv-sol/contracts/BTCUtils.sol";
import {IBondedECDSAKeep} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeep.sol";
import {SafeMath} from "openzeppelin-solidity/contracts/math/SafeMath.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {DepositLiquidation} from "./DepositLiquidation.sol";
import {DepositStates} from "./DepositStates.sol";
import {OutsourceDepositLogging} from "./OutsourceDepositLogging.sol";
import {TBTCConstants} from "../system/TBTCConstants.sol";

library DepositFunding {

    using SafeMath for uint256;
    using SafeMath for uint64;
    using BTCUtils for bytes;
    using BytesLib for bytes;

    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;
    using DepositLiquidation for DepositUtils.Deposit;
    using OutsourceDepositLogging for DepositUtils.Deposit;

    /// @notice     Deletes state after funding.
    /// @dev        This is called when we go to ACTIVE or setup fails without fraud.
    function fundingTeardown(DepositUtils.Deposit storage _d) internal {
        _d.signingGroupRequestedAt = 0;
        _d.fundingProofTimerStart = 0;
    }

    /// @notice     Deletes state after the funding ECDSA fraud process.
    /// @dev        This is only called as we transition to setup failed.
    function fundingFraudTeardown(DepositUtils.Deposit storage _d) internal {
        _d.keepAddress = address(0);
        _d.signingGroupRequestedAt = 0;
        _d.fundingProofTimerStart = 0;
        _d.signingGroupPubkeyX = bytes32(0);
        _d.signingGroupPubkeyY = bytes32(0);
    }

    /// @notice         Internally called function to set up a newly created Deposit instance.
    ///                 This should not be called by developers, use `DepositFactory.createNewDeposit`
    ///                 to create a new deposit.
    /// @dev            If called directly, the transaction will revert since the call will be
    ///                 executed on an already set-up instance.
    /// @param _d       Deposit storage pointer.
    /// @param _m       Signing group honesty threshold.
    /// @param _n       Signing group size.
    function createNewDeposit(
        DepositUtils.Deposit storage _d,
        uint16 _m,
        uint16 _n,
        uint64 _lotSizeSatoshis
    ) public {
        require(_d.tbtcSystem.getAllowNewDeposits(), "New deposits aren't allowed.");
        require(_d.inStart(), "Deposit setup already requested");
        require(_d.tbtcSystem.isAllowedLotSize(_lotSizeSatoshis), "provided lot size not supported");

        _d.lotSizeSatoshis = _lotSizeSatoshis;
        uint256 _bondRequirementSatoshi = _lotSizeSatoshis.mul(_d.tbtcSystem.getInitialCollateralizedPercent()).div(100);
        uint256 _bondRequirementWei = _d.fetchBitcoinPrice().mul(_bondRequirementSatoshi);

        _d.keepSetupFee = _d.tbtcSystem.getNewDepositFeeEstimate();
        /* solium-disable-next-line value-in-payable */
        _d.keepAddress = _d.tbtcSystem.requestNewKeep.value(msg.value)(
            _m,
            _n,
            _bondRequirementWei,
            TBTCConstants.getDepositTerm()
        );

        _d.signerFeeDivisor = _d.tbtcSystem.getSignerFeeDivisor();
        _d.undercollateralizedThresholdPercent = _d.tbtcSystem.getUndercollateralizedThresholdPercent();
        _d.severelyUndercollateralizedThresholdPercent = _d.tbtcSystem.getSeverelyUndercollateralizedThresholdPercent();
        _d.initialCollateralizedPercent = _d.tbtcSystem.getInitialCollateralizedPercent();
        _d.signingGroupRequestedAt = block.timestamp;

        _d.setAwaitingSignerSetup();
        _d.logCreated(_d.keepAddress);
    }

    /// @notice     Anyone may notify the contract that signing group setup has timed out.
    /// @param  _d  Deposit storage pointer.
    function notifySignerSetupFailure(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingSignerSetup(), "Not awaiting setup");
        require(
            block.timestamp > _d.signingGroupRequestedAt.add(TBTCConstants.getSigningGroupFormationTimeout()),
            "Signing group formation timeout not yet elapsed"
        );

        // refund the deposit owner the cost to create a new Deposit at the time the Deposit was opened.
        uint256 _seized = _d.seizeSignerBonds();

        /* solium-disable-next-line security/no-send */
        _d.enableWithdrawal(_d.depositOwner(), _d.keepSetupFee);
        _d.pushFundsToKeepGroup(_seized.sub(_d.keepSetupFee));

        _d.setFailedSetup();
        _d.logSetupFailed();

        fundingTeardown(_d);
    }

    /// @notice             we poll the Keep contract to retrieve our pubkey.
    /// @dev                We store the pubkey as 2 bytestrings, X and Y.
    /// @param  _d          Deposit storage pointer.
    /// @return             True if successful, otherwise revert.
    function retrieveSignerPubkey(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingSignerSetup(), "Not currently awaiting signer setup");

        bytes memory _publicKey = IBondedECDSAKeep(_d.keepAddress).getPublicKey();
        require(_publicKey.length == 64, "public key not set or not 64-bytes long");

        _d.signingGroupPubkeyX = _publicKey.slice(0, 32).toBytes32();
        _d.signingGroupPubkeyY = _publicKey.slice(32, 32).toBytes32();
        require(_d.signingGroupPubkeyY != bytes32(0) && _d.signingGroupPubkeyX != bytes32(0), "Keep returned bad pubkey");
        _d.fundingProofTimerStart = block.timestamp;

        _d.setAwaitingBTCFundingProof();
        _d.logRegisteredPubkey(
            _d.signingGroupPubkeyX,
            _d.signingGroupPubkeyY);
    }

    /// @notice     Anyone may notify the contract that the funder has failed to send BTC.
    /// @dev        This is considered a funder fault, and the funder's payment
    ///             for opening the deposit is not refunded.
    /// @param  _d  Deposit storage pointer.
    function notifyFundingTimeout(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingBTCFundingProof(), "Funding timeout has not started");
        require(
            block.timestamp > _d.fundingProofTimerStart.add(TBTCConstants.getFundingTimeout()),
            "Funding timeout has not elapsed."
        );
        _d.setFailedSetup();
        _d.logSetupFailed();

        _d.closeKeep();
        fundingTeardown(_d);
    }

    /// @notice Requests a funder abort for a failed-funding deposit; that is,
    ///         requests return of a sent UTXO to `_abortOutputScript`. This can
    ///         be used for example when a UTXO is sent that is the wrong size
    ///         for the lot. Must be called after setup fails for any reason,
    ///         and imposes no requirement or incentive on the signing group to
    ///         return the UTXO.
    /// @dev This is a self-admitted funder fault, and should only be callable
    ///      by the TDT holder.
    /// @param _d Deposit storage pointer.
    /// @param _abortOutputScript The output script the funder wishes to request
    ///        a return of their UTXO to.
    function requestFunderAbort(
        DepositUtils.Deposit storage _d,
        bytes memory _abortOutputScript
    ) public {
        require(
            _d.inFailedSetup(),
            "The deposit has not failed funding"
        );

        _d.logFunderRequestedAbort(_abortOutputScript);
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud during funding.
    /// @dev                    Calls out to the keep to verify if there was fraud.
    /// @param  _d              Deposit storage pointer.
    /// @param  _v              Signature recovery value.
    /// @param  _r              Signature R value.
    /// @param  _s              Signature S value.
    /// @param _signedDigest    The digest signed by the signature vrs tuple.
    /// @param _preimage        The sha256 preimage of the digest.
    /// @return                 True if successful, otherwise revert.
    function provideFundingECDSAFraudProof(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes memory _preimage
    ) public {
        require(
            _d.inAwaitingBTCFundingProof(),
            "Signer fraud during funding flow only available while awaiting funding"
        );

        _d.submitSignatureFraud(_v, _r, _s, _signedDigest, _preimage);
        _d.logFraudDuringSetup();

        // Allow deposit owner to withdraw seized bonds after contract termination.
        uint256 _seized = _d.seizeSignerBonds();
        _d.enableWithdrawal(_d.depositOwner(), _seized);

        fundingFraudTeardown(_d);
        _d.setFailedSetup();
        _d.logSetupFailed();
    }

    /// @notice                     Anyone may notify the deposit of a funding proof to activate the deposit.
    ///                             This is the happy-path of the funding flow. It means that we have succeeded.
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding.
    /// @param  _d                  Deposit storage pointer.
    /// @param _txVersion           Transaction version number (4-byte LE).
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs.
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs.
    /// @param _txLocktime          Final 4 bytes of the transaction.
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed).
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block.
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed).
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first.
    function provideBTCFundingProof(
        DepositUtils.Deposit storage _d,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public {

        require(_d.inAwaitingBTCFundingProof(), "Not awaiting funding");

        bytes8 _valueBytes;
        bytes memory  _utxoOutpoint;

        (_valueBytes, _utxoOutpoint) = _d.validateAndParseFundingSPVProof(
            _txVersion,
            _txInputVector,
            _txOutputVector,
            _txLocktime,
            _fundingOutputIndex,
            _merkleProof,
            _txIndexInBlock,
            _bitcoinHeaders
        );

        // Write down the UTXO info and set to active. Congratulations :)
        _d.utxoSizeBytes = _valueBytes;
        _d.utxoOutpoint = _utxoOutpoint;
        _d.fundedAt = block.timestamp;

        bytes32 _txid = abi.encodePacked(_txVersion, _txInputVector, _txOutputVector, _txLocktime).hash256();

        fundingTeardown(_d);
        _d.setActive();
        _d.logFunded(_txid);
    }
}
