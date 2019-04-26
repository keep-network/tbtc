pragma solidity 0.4.25;

import {DepositUtils} from '../../contracts/deposit/DepositUtils.sol';
import {Deposit} from '../../contracts/deposit/Deposit.sol';

contract TestDeposit is Deposit {

    struct Deposit {

        // SET DURING CONSTRUCTION
        uint8 currentState;

        // SET ON FRAUD
        uint256 liquidationInitiated;  // Timestamp of when liquidation starts
        uint256 courtesyCallInitiated; // When the courtesy call is issued

        // written when we request a keep
        uint256 keepID;  // The ID of our keep group
        uint256 signingGroupRequestedAt;  // timestamp of signing group request

        // written when we get a keep result
        uint256 fundingProofTimerStart;  // start of the funding proof period. reused for funding fraud proof period
        bytes32 signingGroupPubkeyX;  // The X coordinate of the signing group's pubkey
        bytes32 signingGroupPubkeyY;  // The Y coordinate of the signing group's pubkey

        // INITIALLY WRITTEN BY REDEMPTION FLOW
        address requesterAddress;  // The requester's addr4ess, used as fallback for fraud in redemption
        bytes20 requesterPKH;  // The 20-byte requeser PKH
        uint256 initialRedemptionFee;  // the initial fee as requested
        uint256 withdrawalRequestTime;  // the most recent withdrawal request timestamp
        bytes32 lastRequestedDigest;  // the digest most recently requested for signing

        // written when we get funded
        bytes8 utxoSizeBytes;  // LE uint. the size of the deposit UTXO in satoshis
        uint256 fundedAt; // timestamp when funding proof was received
        bytes utxoOutpoint;  // the 36-byte outpoint of the custodied UTXO
    }

    function setState(uint8 _state) public {
        self.currentState = _state;
    }

    function setLiquidationOrCourtesyInitated(
        uint256 _liquidation,
        uint256 _courtesy
    ) public {
        self.liquidationInitiated = _liquidation;
        self.courtesyCallInitiated = _courtesy;
    }

    function setKeepInfo(
        uint256 _keepID,
        uint256 _signingGroupRequestedAt,
        uint256 _fundingProofTimerStart,
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) public {
        self.keepID = _keepID;
        self.signingGroupRequestedAt = _signingGroupRequestedAt;
        self.fundingProofTimerStart = _fundingProofTimerStart;
        self.signingGroupPubkeyX = _signingGroupPubkeyX;
        self.signingGroupPubkeyY = _signingGroupPubkeyY;
    }

    function setRequestInfo(
        address _requesterAddress,
        bytes20 _requesterPKH,
        uint256 _initialRedemptionFee,
        uint256 _withdrawalRequestTime,
        bytes32 _lastRequestedDigest
    ) public {
        self.requesterAddress = _requesterAddress;
        self.requesterPKH = _requesterPKH;
        self.initialRedemptionFee = _initialRedemptionFee;
        self.withdrawalRequestTime = _withdrawalRequestTime;
        self.lastRequestedDigest = _lastRequestedDigest;
    }

    function setUTXOInfo(
        bytes8 _utxoSizeBytes,
        uint256 _fundedAt,
        bytes _utxoOutpoint
    ) public {
        self.utxoSizeBytes = _utxoSizeBytes;
        self.fundedAt = _fundedAt;
        self.utxoOutpoint = _utxoOutpoint;
    }

}
