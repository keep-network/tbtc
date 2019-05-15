pragma solidity 0.4.25;

import {Deposit} from '../../../contracts/deposit/Deposit.sol';

contract TestDeposit is Deposit {

    function setExteroriorAddresses(
        address _sys,
        address _token,
        address _k
    ) public {
        self.TBTCSystem = _sys;
        self.KeepSystem = _k;
        self.TBTCToken = _token;
    }

    function reset() public {
        setState(0);
        setLiquidationAndCourtesyInitated(0, 0);
        setKeepInfo(0, 0, 0, bytes32(0), bytes32(0));
        setRequestInfo(address(0), bytes20(0), 0, 0, bytes32(0));
        setUTXOInfo(bytes8(0), 0, '');
    }

    function setState(uint8 _state) public {
        self.currentState = _state;
    }

    function getState() public view returns (uint8) { return self.currentState; }

    function setLiquidationAndCourtesyInitated(
        uint256 _liquidation,
        uint256 _courtesy
    ) public {
        self.liquidationInitiated = _liquidation;
        self.courtesyCallInitiated = _courtesy;
    }

    function getLiqudationAndCoutesyInitiated() public view returns (uint256, uint256) {
        return (self.liquidationInitiated, self.courtesyCallInitiated);
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

    function getKeepInfo() public view returns (uint256, uint256, uint256, bytes32, bytes32) {
        return (self.keepID, self.signingGroupRequestedAt, self.fundingProofTimerStart, self.signingGroupPubkeyX, self.signingGroupPubkeyY);
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

    function getRequestInfo() public view returns (address, bytes20, uint256, uint256, bytes32) {
        return (
            self.requesterAddress,
            self.requesterPKH,
            self.initialRedemptionFee,
            self.withdrawalRequestTime,
            self.lastRequestedDigest);
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

    function getUTXOInfo() public view returns (bytes8, uint256, bytes) {
        return (self.utxoSizeBytes, self.fundedAt, self.utxoOutpoint);
    }

    // passthrough for direct testing
    function redemptionTransactionChecks(bytes _bitcoinTx) public view returns (bytes32, uint256) {
        return self.redemptionTransactionChecks(_bitcoinTx);
    }

    function getWithdrawalRequestTime() public view returns(uint256){
        return self.withdrawalRequestTime;
    }
}
