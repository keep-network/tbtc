pragma solidity ^0.5.10;

import {Deposit} from '../../../contracts/deposit/Deposit.sol';

contract TestDeposit is Deposit {

    function setExteriorAddresses(
        address _sys,
        address _token
    ) public {
        self.TBTCSystem = _sys;
        self.TBTCToken = _token;
    }

    function reset() public {
        setState(0);
        setLiquidationAndCourtesyInitated(0, 0);
        setRequestInfo(address(0), bytes20(0), 0, 0, bytes32(0));
        setUTXOInfo(bytes8(0), 0, '');

        setKeepAddress(address(0));
        setSigningGroupRequestedAt(0);
        setFundingProofTimerStart(0);
        setSigningGroupPublicKey(bytes32(0), bytes32(0));
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

    function getLiquidationAndCourtesyInitiated() public view returns (uint256, uint256) {
        return (self.liquidationInitiated, self.courtesyCallInitiated);
    }

    function setKeepAddress(address _keepAddress) public {
        self.keepAddress = _keepAddress;
    }

    function getKeepAddress() public view returns (address) {
        return self.keepAddress;
    }

    function setSigningGroupRequestedAt(uint256 _signingGroupRequestedAt) public {
        self.signingGroupRequestedAt = _signingGroupRequestedAt;
    }

    function getSigningGroupRequestedAt() public view returns (uint256) {
        return self.signingGroupRequestedAt;
    }

    function setFundingProofTimerStart(uint256 _fundingProofTimerStart) public {
        self.fundingProofTimerStart = _fundingProofTimerStart;
    }

    function getFundingProofTimerStart() public view returns (uint256) {
        return self.fundingProofTimerStart;
    }

    function setSigningGroupPublicKey(bytes32 _x,bytes32 _y) public {
        self.signingGroupPubkeyX = _x;
        self.signingGroupPubkeyY = _y;
    }

    function getSigningGroupPublicKey() public returns (bytes32, bytes32){
        return (self.signingGroupPubkeyX, self.signingGroupPubkeyY);
    }

    function setRequestInfo(
        address payable _requesterAddress,
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
        bytes memory _utxoOutpoint
    ) public {
        self.utxoSizeBytes = _utxoSizeBytes;
        self.fundedAt = _fundedAt;
        self.utxoOutpoint = _utxoOutpoint;
    }

    function getUTXOInfo() public view returns (bytes8, uint256, bytes memory) {
        return (self.utxoSizeBytes, self.fundedAt, self.utxoOutpoint);
    }

    function setDigestApprovedAtTime(bytes32 _digest, uint256 _timestamp) public {
        self.approvedDigests[_digest] = _timestamp;
    }

    // passthrough for direct testing
    function approveDigest(bytes32 _digest) public {
        return self.approveDigest(_digest);
    }

    function wasDigestApprovedForSigning(bytes32 _digest) public view returns (uint256) {
        return self.wasDigestApprovedForSigning(_digest);
    }

     function redemptionTransactionChecks(
        bytes memory _txInputVector,
        bytes memory _txOutputVector
    ) public view returns (uint256) {
        return self.redemptionTransactionChecks(_txInputVector, _txOutputVector);
    }

    function validateRedeemerNotPaid(bytes memory _txOutputVector) public view returns (bool){
        return self.validateRedeemerNotPaid(_txOutputVector);
    }

    function getWithdrawalRequestTime() public view returns(uint256){
        return self.withdrawalRequestTime;
    }

    function pushFundsToKeepGroup(uint256 _ethValue) public returns (bool) {
        return self.pushFundsToKeepGroup(_ethValue);
    }

    function attemptToLiquidateOnchain() public returns (bool) {
        return self.attemptToLiquidateOnchain();
    }
}
