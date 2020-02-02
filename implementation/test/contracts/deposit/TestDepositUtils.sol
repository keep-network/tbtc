pragma solidity ^0.5.10;

import {DepositUtils} from '../../../contracts/deposit/DepositUtils.sol';
import {TestDeposit} from './TestDeposit.sol';

contract TestDepositUtils is TestDeposit {
    constructor(address _factoryAddress) 
        TestDeposit(_factoryAddress)
    public{}

    // Passthroughs to test view and pure functions

    function currentBlockDifficulty() public view returns (uint256) {
        return self.currentBlockDifficulty();
    }

    function previousBlockDifficulty() public view returns (uint256) {
        return self.previousBlockDifficulty();
    }

    function evaluateProofDifficulty(bytes memory _bitcoinHeaders) public view {
        return self.evaluateProofDifficulty(_bitcoinHeaders);
    }

    function checkProofFromTxId(
        bytes32 _bitcoinTxId,
        bytes memory _merkleProof,
        uint256 _index,
        bytes memory _bitcoinHeaders
    ) public view returns (bytes32) {
        self.checkProofFromTxId(_bitcoinTxId, _merkleProof, _index, _bitcoinHeaders);
    }
    
    function setPubKey(
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) public {
        self.signingGroupPubkeyX = _signingGroupPubkeyX;
        self.signingGroupPubkeyY = _signingGroupPubkeyY;
    }

    function findAndParseFundingOutput(
        bytes memory _txOutputVector,
        uint8 _fundingOutputIndex
    ) public view returns (bytes8) {
        return self.findAndParseFundingOutput(_txOutputVector, _fundingOutputIndex);
    }

    function validateAndParseFundingSPVProof(
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public view returns (bytes8 _valueBytes, bytes memory _utxoOutpoint){
      return self.validateAndParseFundingSPVProof(
        _txVersion,
        _txInputVector,
        _txOutputVector,
        _txLocktime,
        _fundingOutputIndex,
        _merkleProof,
        _txIndexInBlock,
        _bitcoinHeaders
        );
    }

    function auctionValue() public view returns (uint256) {
        return self.auctionValue();
    }

    function signerFee() public view returns (uint256) {
        return self.signerFee();
    }

    function auctionTBTCAmount() public view returns (uint256) {
        return self.auctionTBTCAmount();
    }

    function determineCompressionPrefix(bytes32 _pubkeyY) public pure returns (bytes memory) {
        return DepositUtils.determineCompressionPrefix(_pubkeyY);
    }

    function compressPubkey(bytes32 _pubkeyX, bytes32 _pubkeyY) public pure returns (bytes memory) {
        return DepositUtils.compressPubkey(_pubkeyX, _pubkeyY);
    }

    function signerPubkey() public view returns (bytes memory) {
        return self.signerPubkey();
    }

    function signerPKH() public view returns (bytes20) {
        return self.signerPKH();
    }

    function utxoSize() public view returns (uint256) {
        return self.utxoSize();
    }

    function fetchBitcoinPrice() public view returns (uint256) {
        return self.fetchBitcoinPrice();
    }

    function fetchBondAmount() public view returns (uint256) {
        return self.fetchBondAmount();
    }

    function bytes8LEToUint(bytes8 _b) public pure returns (uint256) {
        return DepositUtils.bytes8LEToUint(_b);
    }

    function feeRebateTokenHolder() public view returns (address payable) {
        return self.feeRebateTokenHolder();
    }

    function redemptionTeardown() public {
        return self.redemptionTeardown();
    }

    function seizeSignerBonds() public returns (uint256) {
        return self.seizeSignerBonds();
    }

    function distributeFeeRebate() public {
        return self.distributeFeeRebate();
    }

    function pushFundsToKeepGroup(uint256 _ethValue) public returns (bool) {
        return self.pushFundsToKeepGroup(_ethValue);
    }
}
