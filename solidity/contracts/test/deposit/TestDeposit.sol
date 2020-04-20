pragma solidity 0.5.17;

import {Deposit} from "../../../contracts/deposit/Deposit.sol";
import {ITBTCSystem} from "../../interfaces/ITBTCSystem.sol";
import {IERC721} from "openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import {TBTCToken} from "../../system/TBTCToken.sol";
import {FeeRebateToken} from "../../system/FeeRebateToken.sol";

contract TestDeposit is Deposit {

    constructor() public {
    // solium-disable-previous-line no-empty-blocks
    }

    function createNewDeposit(
        ITBTCSystem _tbtcSystem,
        TBTCToken _tbtcToken,
        IERC721 _tbtcDepositToken,
        FeeRebateToken _feeRebateToken,
        address _vendingMachineAddress,
        uint16 _m,
        uint16 _n,
        uint64 _lotSizeSatoshis
    ) public payable returns (bool) {
        self.tbtcSystem = _tbtcSystem;
        self.tbtcToken = _tbtcToken;
        self.tbtcDepositToken = _tbtcDepositToken;
        self.feeRebateToken = _feeRebateToken;
        self.vendingMachineAddress = _vendingMachineAddress;
        self.createNewDeposit(_m, _n, _lotSizeSatoshis);
        return true;
    }

    function setExteriorAddresses(
        address _sys,
        address _token,
        address _tbtcDepositToken,
        address _feeRebateToken,
        address _vendingMachineAddress
    ) public {
        self.tbtcSystem = ITBTCSystem(_sys);
        self.tbtcToken = TBTCToken(_token);
        self.tbtcDepositToken = IERC721(_tbtcDepositToken);
        self.feeRebateToken = FeeRebateToken(_feeRebateToken);
        self.vendingMachineAddress = _vendingMachineAddress;
    }

    function reset() public {
        setState(0);
        setLiquidationAndCourtesyInitated(0, 0);
        setRequestInfo(address(0), "", 0, 0, bytes32(0));
        setUTXOInfo(bytes8(0), 0, "");

        setKeepAddress(address(0));
        setSigningGroupRequestedAt(0);
        setFundingProofTimerStart(0);
        setSigningGroupPublicKey(bytes32(0), bytes32(0));
    }

    function setState(uint8 _state) public {
        self.currentState = _state;
    }

    function getState() public view returns (uint8) { return self.currentState; }

    function setSignerFeeDivisor(uint16 _signerFeeDivisor) public {
        self.signerFeeDivisor = _signerFeeDivisor;
    }

    function setKeepSetupFee(uint256 _fee) public {
        self.keepSetupFee = _fee;
    }

    function getSignerFeeDivisor() public view returns (uint256) { return self.signerFeeDivisor; }

    function setLotSize(uint64 _lotSizeSatoshis) public {
        self.lotSizeSatoshis = _lotSizeSatoshis;
    }

    function setUndercollateralizedThresholdPercent(uint16 _undercollateralizedThresholdPercent) public {
        self.undercollateralizedThresholdPercent = _undercollateralizedThresholdPercent;
    }

    function setInitialCollateralizedPercent(uint16 _initialCollateralizedPercent) public {
        self.initialCollateralizedPercent = _initialCollateralizedPercent;
    }

    function getUndercollateralizedThresholdPercent() public view returns (uint16) { return self.undercollateralizedThresholdPercent; }

    function setSeverelyUndercollateralizedThresholdPercent(uint16 _severelyUndercollateralizedThresholdPercent) public {
        self.severelyUndercollateralizedThresholdPercent = _severelyUndercollateralizedThresholdPercent;
    }

    function getSeverelyUndercollateralizedThresholdPercent() public view returns (uint16) {
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

    function getLiquidationInitiator() public view returns (address) {
        return self.liquidationInitiator;
    }

    function getLiquidationTimestamp() public view returns (uint256) {
        return self.liquidationInitiated;
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
        bytes memory _redeemerOutputScript,
        uint256 _initialRedemptionFee,
        uint256 _withdrawalRequestTime,
        bytes32 _lastRequestedDigest
    ) public {
        self.redeemerAddress = _redeemerAddress;
        self.redeemerOutputScript = _redeemerOutputScript;
        self.initialRedemptionFee = _initialRedemptionFee;
        self.withdrawalRequestTime = _withdrawalRequestTime;
        self.lastRequestedDigest = _lastRequestedDigest;
    }

    function setLatestRedemptionFee(uint256 _latestRedemptionFee) public {
        self.latestRedemptionFee = _latestRedemptionFee;
    }

    function setRedeemerAddress(
        address payable _redeemerAddress
    ) public {
        self.redeemerAddress = _redeemerAddress;
    }
    function getRequestInfo() public view returns (address, bytes memory, uint256, uint256, bytes32) {
        return (
            self.redeemerAddress,
            self.redeemerOutputScript,
            self.initialRedemptionFee,
            self.withdrawalRequestTime,
            self.lastRequestedDigest
        );
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

    function startLiquidation(bool _wasFraud) public {
        self.startLiquidation(_wasFraud);
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

    function getCollateralizationPercentage() public view returns (uint256) {
        return self.getCollateralizationPercentage();
    }

    function pushFundsToKeepGroup(uint256 _ethValue) public returns (bool) {
        return self.pushFundsToKeepGroup(_ethValue);
    }

    function getAuctionBasePercentage() public view returns (uint256) {
        return self.getAuctionBasePercentage();
    }

    function auctionValue() public view returns (uint256) {
        return self.auctionValue();
    }

}
