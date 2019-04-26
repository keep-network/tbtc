pragma solidity 0.4.25;

import {DepositUtils} from '../../contracts/deposit/DepositUtils.sol';
import {Deposit} from '../../contracts/deposit/Deposit.sol';

contract TestDeposit is Deposit {

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
