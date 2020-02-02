pragma solidity ^0.5.10;

import {Deposit} from '../../../contracts/deposit/Deposit.sol';

contract TestDeposit is Deposit {

    constructor(address _factoryAddress) 
        Deposit(_factoryAddress)
    public{}

    function createNewDeposit(
        address _TBTCSystem,
        address _TBTCToken,
        address _TBTCDepositToken,
        address _FeeRebateToken,
        address _VendingMachine,
        uint256 _m,
        uint256 _n,
        uint256 _lotSize
    ) public payable returns (bool) {
        self.TBTCSystem = _TBTCSystem;
        self.TBTCToken = _TBTCToken;
        self.TBTCDepositToken = _TBTCDepositToken;
        self.FeeRebateToken = _FeeRebateToken;
        self.VendingMachine = _VendingMachine;
        self.createNewDeposit(_m, _n, _lotSize);
        return true;
    }

    function setExteriorAddresses(
        address _sys,
        address _token,
        address _tbtcDepositToken,
        address _feeRebateToken,
        address _vendingMachine
    ) public {
        self.TBTCSystem = _sys;
        self.TBTCToken = _token;
        self.TBTCDepositToken = _tbtcDepositToken;
        self.FeeRebateToken = _feeRebateToken;
        self.VendingMachine = _vendingMachine;
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

    function setSignerFeeDivisor(uint256 _signerFeeDivisor) public {
        self.signerFeeDivisor = _signerFeeDivisor;
    }

    function getSignerFeeDivisor() public view returns (uint256) { return self.signerFeeDivisor; }
    
    function setLotSize(uint256 _lotSizeSatoshis) public {
        self.lotSizeSatoshis = _lotSizeSatoshis;
    }

    function setUndercollateralizedThresholdPercent(uint128 _undercollateralizedThresholdPercent) public {
        self.undercollateralizedThresholdPercent = _undercollateralizedThresholdPercent;
    }

    function getUndercollateralizedThresholdPercent() public view returns (uint128) { return self.undercollateralizedThresholdPercent; }

    function setSeverelyUndercollateralizedThresholdPercent(uint128 _severelyUndercollateralizedThresholdPercent) public {
        self.severelyUndercollateralizedThresholdPercent = _severelyUndercollateralizedThresholdPercent;
    }

    function getSeverelyUndercollateralizedThresholdPercent() public view returns (uint128) {
        return self.severelyUndercollateralizedThresholdPercent;
    }

    function setLiquidationAndCourtesyInitated(
        uint256 _liquidation,
        uint256 _courtesy
    ) public {
        self.liquidationInitiated = _liquidation;
        self.courtesyCallInitiated = _courtesy;
    }
    function setLiquidationInitiator(address payable _initiator) public {
        self.liquidationInitiator = _initiator;
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

    function getSigningGroupPublicKey() public view returns (bytes32, bytes32){
        return (self.signingGroupPubkeyX, self.signingGroupPubkeyY);
    }

    function setRequestInfo(
        address payable _redeemerAddress,
        bytes20 _redeemerPKH,
        uint256 _initialRedemptionFee,
        uint256 _withdrawalRequestTime,
        bytes32 _lastRequestedDigest
    ) public {
        self.redeemerAddress = _redeemerAddress;
        self.redeemerPKH = _redeemerPKH;
        self.initialRedemptionFee = _initialRedemptionFee;
        self.withdrawalRequestTime = _withdrawalRequestTime;
        self.lastRequestedDigest = _lastRequestedDigest;
    }
    function setRedeemerAddress(
        address payable _redeemerAddress
    ) public {
        self.redeemerAddress = _redeemerAddress;
    }
    function getRequestInfo() public view returns (address, bytes20, uint256, uint256, bytes32) {
        return (
            self.redeemerAddress,
            self.redeemerPKH,
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

    function getRedemptionTbtcRequirement(address _redeemer) public view returns (uint256) {
        return self.getRedemptionTbtcRequirement(_redeemer);
    }

   function getOwnerRedemptionTbtcRequirement(address _redeemer) public view returns (uint256) {
        return self.getOwnerRedemptionTbtcRequirement(_redeemer);
    }

     function performRedemptionTBTCTransfers() public {
        self.performRedemptionTBTCTransfers();
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
}
