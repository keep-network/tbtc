pragma solidity 0.4.25;

import {SafeMath} from "../bitcoin-spv/SafeMath.sol";
import {BytesLib} from "../bitcoin-spv/BytesLib.sol";
import {BTCUtils} from "../bitcoin-spv/BTCUtils.sol";
import {IBurnableERC20} from "../interfaces/IBurnableERC20.sol";
import {IKeep} from "../interfaces/IKeep.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {DepositLiquidation} from "./DepositLiquidation.sol";
import {DepositStates} from "./DepositStates.sol";
import {OutsourceDepositLogging} from "./OutsourceDepositLogging.sol";
import {TBTCConstants} from "./TBTCConstants.sol";

library DepositFunding {

    using SafeMath for uint256;
    using BTCUtils for bytes;
    using BytesLib for bytes;

    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;
    using DepositLiquidation for DepositUtils.Deposit;
    using OutsourceDepositLogging for DepositUtils.Deposit;

    /// @notice     Deletes state after funding
    /// @dev        This is called when we go to ACTIVE or setup fails without fraud
    function fundingTeardown(DepositUtils.Deposit storage _d) public {
        _d.signingGroupRequestedAt = 0;
        _d.fundingProofTimerStart = 0;
    }

    /// @notice     Deletes state after the funding ECDSA fraud process
    /// @dev        This is only called as we transition to setup failed
    function fundingFraudTeardown(DepositUtils.Deposit storage _d) public {
        _d.keepID = 0;
        _d.signingGroupRequestedAt = 0;
        _d.fundingProofTimerStart = 0;
        _d.signingGroupPubkeyX = bytes32(0);
        _d.signingGroupPubkeyY = bytes32(0);
    }

    /// @notice         get the signer pubkey for our keep
    /// @dev            calls out to the keep contract, should get 64 bytes back
    /// @return         the 64 byte pubkey
    function getKeepPubkeyResult(DepositUtils.Deposit storage _d) public view returns (bytes) {
        IKeep _keep = IKeep(_d.KeepBridge);
        bytes memory _pubkey = _keep.getKeepPubkey(_d.keepID);
        /* solium-disable-next-line */
        require(_pubkey.length == 64, "public key not set or not 64-bytes long");
        return _pubkey;
    }

    /// @notice         The system can spin up a new deposit
    /// @dev            This should be called by an approved contract, not a developer
    /// @param _d       deposit storage pointer
    /// @param _m       m for m-of-n
    /// @param _m       n for m-of-n
    /// @return         True if successful, otherwise revert
    function createNewDeposit(
        DepositUtils.Deposit storage _d,
        uint256 _m,
        uint256 _n
    ) public returns (bool) {
        require(_d.inStart(), "Deposit setup already requested");

        IKeep _keep = IKeep(_d.KeepBridge);

        /* solium-disable-next-line value-in-payable */
        _d.keepID = _keep.requestKeepGroup.value(msg.value)(_m, _n);  // kinda gross but
        _d.signingGroupRequestedAt = block.timestamp;

        _d.setAwaitingSignerSetup();
        _d.logCreated(_d.keepID);

        return true;
    }

    /// @notice     Transfers the funders bond to the signers if the funder never funds
    /// @dev        Called only by notifyFundingTimeout
    function revokeFunderBond(DepositUtils.Deposit storage _d) public {
        if (address(this).balance >= TBTCConstants.getFunderBondAmount()) {
            _d.pushFundsToKeepGroup(TBTCConstants.getFunderBondAmount());
        } else if (address(this).balance > 0) {
            _d.pushFundsToKeepGroup(address(this).balance);
        }
    }

    /// @notice     Returns the funder's bond plus a payment at contract teardown
    /// @dev        Returns the balance if insufficient. Always call this before distributing signer payments
    function returnFunderBond(DepositUtils.Deposit storage _d) public {
        if (address(this).balance >= TBTCConstants.getFunderBondAmount()) {
            _d.depositBeneficiary().transfer(TBTCConstants.getFunderBondAmount());
        } else if (address(this).balance > 0) {
            _d.depositBeneficiary().transfer(address(this).balance);
        }
    }

    /// @notice     slashes the signers partially for committing fraud before funding occurs
    /// @dev        called only by notifyFraudFundingTimeout
    function partiallySlashForFraudInFunding(DepositUtils.Deposit storage _d) public {
        uint256 _seized = _d.seizeSignerBonds();
        uint256 _slash = _seized.div(TBTCConstants.getFundingFraudPartialSlashDivisor());
        _d.pushFundsToKeepGroup(_seized.sub(_slash));
        _d.depositBeneficiary().transfer(_slash);
    }

    /// @notice     Seizes signer bonds and distributes them to the funder
    /// @dev        This is only called as part of funding fraud flow
    function distributeSignerBondsToFunder(DepositUtils.Deposit storage _d) public {
        uint256 _seized = _d.seizeSignerBonds();
        _d.depositBeneficiary().transfer(_seized);  // Transfer whole amount
    }

    /// @notice     Anyone may notify the contract that signing group setup has timed out
    /// @dev        We rely on the keep system punishes the signers in this case
    /// @param  _d  deposit storage pointer
    function notifySignerSetupFailure(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingSignerSetup(), "Not awaiting setup");
        require(
            block.timestamp > _d.signingGroupRequestedAt + TBTCConstants.getSigningGroupFormationTimeout(),
            "Signing group formation timeout not yet elapsed"
        );
        _d.setFailedSetup();
        _d.logSetupFailed();

        returnFunderBond(_d);
        fundingTeardown(_d);
    }

    /// @notice             we poll the Keep contract to retrieve our pubkey
    /// @dev                We store the pubkey as 2 bytestrings, X and Y.
    /// @param  _d          deposit storage pointer
    /// @return             True if successful, otherwise revert
    function retrieveSignerPubkey(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingSignerSetup(), "Not currently awaiting signer setup");
        bytes memory _keepResult = getKeepPubkeyResult(_d);

        _d.signingGroupPubkeyX = _keepResult.slice(0, 32).toBytes32();
        _d.signingGroupPubkeyY = _keepResult.slice(32, 32).toBytes32();
        require(_d.signingGroupPubkeyY != bytes32(0) && _d.signingGroupPubkeyX != bytes32(0), "Keep returned bad pubkey");
        _d.fundingProofTimerStart = block.timestamp;

        _d.setAwaitingBTCFundingProof();
        _d.logRegisteredPubkey(
            _d.signingGroupPubkeyX,
            _d.signingGroupPubkeyY);
    }

    /// @notice     Anyone may notify the contract that the funder has failed to send BTC
    /// @dev        This is considered a funder fault, and we revoke their bond
    /// @param  _d  deposit storage pointer
    function notifyFundingTimeout(DepositUtils.Deposit storage _d) public {
        require(_d.inAwaitingBTCFundingProof(), "Funding timeout has not started");
        require(
            block.timestamp > _d.fundingProofTimerStart + TBTCConstants.getFundingTimeout(),
            "Funding timeout has not elapsed."
        );
        _d.setFailedSetup();
        _d.logSetupFailed();

        revokeFunderBond(_d);
        fundingTeardown(_d);
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud during funding
    /// @dev                    ECDSA is NOT SECURE unless you verify the digest
    /// @param  _d              deposit storage pointer
    /// @param  _v              Signature recovery value
    /// @param  _r              Signature R value
    /// @param  _s              Signature S value
    /// @param _signedDigest    The digest signed by the signature vrs tuple
    /// @param _preimage        The sha256 preimage of the digest
    /// @return                 True if successful, otherwise revert
    function provideFundingECDSAFraudProof(
        DepositUtils.Deposit storage _d,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) public {
        require(
            _d.inAwaitingBTCFundingProof(),
            "Signer fraud during funding flow only available while awaiting funding"
        );

        bool _isFraud = _d.submitSignatureFraud(_v, _r, _s, _signedDigest, _preimage);
        require(_isFraud, "Signature is not fraudulent");
        _d.seizeSignerBonds();
        _d.logFraudDuringSetup();

        // If the funding timeout has elapsed, punish the funder too!
        if (block.timestamp > _d.fundingProofTimerStart + TBTCConstants.getFundingTimeout()) {
            address(0).transfer(address(this).balance);  // Burn it all down (fire emoji)
            _d.setFailedSetup();
        } else {
            /* NB: This is reuse of the variable */
            _d.fundingProofTimerStart = block.timestamp;
            _d.setFraudAwaitingBTCFundingProof();
            returnFunderBond(_d);
        }
    }

    /// @notice     Anyone may notify the contract no funding proof was submitted during funding fraud
    /// @dev        This is not a funder fault. The signers have faulted, so the funder shouldn't fund
    /// @param  _d  deposit storage pointer
    function notifyFraudFundingTimeout(DepositUtils.Deposit storage _d) public {
        require(
            _d.inFraudAwaitingBTCFundingProof(),
            "Not currently awaiting fraud-related funding proof"
        );
        require(
            block.timestamp > _d.fundingProofTimerStart + TBTCConstants.getFraudFundingTimeout(),
            "Fraud funding proof timeout has not elapsed"
        );
        _d.setFailedSetup();
        _d.logSetupFailed();

        partiallySlashForFraudInFunding(_d);
        fundingFraudTeardown(_d);
    }

    /// @notice                     Anyone may notify the deposit of a funding proof during funding fraud
    //                              We reward the funder the entire bond if this occurs
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding
    /// @param  _d                  Deposit storage pointer
    /// @param _txVersion           Transaction version number (4-byte LE)
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param _txLocktime          Final 4 bytes of the transaction
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block
    /// @param _txIndexInBlock      Transaction index in the block (1-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first
    /// @return                     True if no errors are thrown
    function provideFraudBTCFundingProof(
        DepositUtils.Deposit storage _d,
        bytes _txVersion,
        bytes _txInputVector,
        bytes _txOutputVector,
        bytes _txLocktime,
        uint8 _fundingOutputIndex,
        bytes _merkleProof,
        uint256 _txIndexInBlock,
        bytes _bitcoinHeaders
    ) public returns (bool) {
        require(_d.inFraudAwaitingBTCFundingProof(), "Not awaiting a funding proof during setup fraud");

        bytes8 _valueBytes;
        bytes memory _utxoOutpoint;

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

        _d.setFailedSetup();
        _d.logSetupFailed();

        // If the proof is accepted, update to failed, and distribute signer bonds
        distributeSignerBondsToFunder(_d);
        fundingFraudTeardown(_d);

        return true;
    }

    /// @notice                     Anyone may notify the deposit of a funding proof to activate the deposit
    ///                             This is the happy-path of the funding flow. It means that we have succeeded
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding
    /// @param  _d                  Deposit storage pointer
    /// @param _txVersion           Transaction version number (4-byte LE)
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param _txLocktime          Final 4 bytes of the transaction
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block
    /// @param _txIndexInBlock      Transaction index in the block (1-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first
    /// @return                     True if no errors are thrown
    function provideBTCFundingProof(
        DepositUtils.Deposit storage _d,
        bytes _txVersion,
        bytes _txInputVector,
        bytes _txOutputVector,
        bytes _txLocktime,
        uint8 _fundingOutputIndex,
        bytes _merkleProof,
        uint256 _txIndexInBlock,
        bytes _bitcoinHeaders
    ) public returns (bool) {

        require(_d.inAwaitingBTCFundingProof(), "Not awaiting funding");

        // Design decision:
        // We COULD revoke the funder bond here if the funding proof timeout has elapsed
        // HOWEVER, that would only create a situation where the funder loses eerything
        // It would be a large punishment for a small crime (being slightly late)
        // So if the funder manages to call this before anyone notifies of timeout
        // We let them have a freebie

        bytes8 _valueBytes;
        bytes memory _utxoOutpoint;

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

        fundingTeardown(_d);
        _d.setActive();
        _d.logFunded();

        returnFunderBond(_d);

        // Mint 95% of the deposit size
        IBurnableERC20 _tbtc = IBurnableERC20(_d.TBTCToken);
        uint256 _value = TBTCConstants.getLotSize();
        _tbtc.mint(_d.depositBeneficiary(), _value.mul(95).div(100));

        return true;
    }
}
