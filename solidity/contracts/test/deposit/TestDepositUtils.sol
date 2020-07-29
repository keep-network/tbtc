pragma solidity 0.5.17;

import {DepositUtils} from "../../../contracts/deposit/DepositUtils.sol";
import {TestDeposit} from "./TestDeposit.sol";

contract TestDepositUtils is TestDeposit {
    constructor() public {
    // solium-disable-previous-line no-empty-blocks
    }

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

    function signerFee() public view returns (uint256) {
        return self.signerFeeTbtc();
    }

    function signerPubkey() public view returns (bytes memory) {
        return self.signerPubkey();
    }

    function signerPKH() public view returns (bytes20) {
        return self.signerPKH();
    }

    function fetchBitcoinPrice() public view returns (uint256) {
        return self.fetchBitcoinPrice();
    }

    function fetchBondAmount() public view returns (uint256) {
        return self.fetchBondAmount();
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

    function pushFundsToKeepGroup(uint256 _ethValue) public {
        self.pushFundsToKeepGroup(_ethValue);
    }

    function enableWithdrawal(address _withdrawer, uint256 _amount) public {
        self.enableWithdrawal(_withdrawer, _amount);
    }
}

// Separate contract for testing SPV proofs, as putting this in the main
// TestDepositUtils contract causes it to run out of gas before finishing its
// deploy.
contract TestDepositUtilsSPV is TestDeposit {
    constructor() public {
        // solium-disable-previous-line no-empty-blocks
    }

    function setPubKey(
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) public {
        self.signingGroupPubkeyX = _signingGroupPubkeyX;
        self.signingGroupPubkeyY = _signingGroupPubkeyY;
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

    function checkProofFromTxId(
        bytes32 _bitcoinTxId,
        bytes memory _merkleProof,
        uint256 _index,
        bytes memory _bitcoinHeaders
    ) public view returns (bytes32) {
        self.checkProofFromTxId(_bitcoinTxId, _merkleProof, _index, _bitcoinHeaders);
    }
}
