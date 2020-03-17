pragma solidity ^0.5.10;

import {
    IBondedECDSAKeep
} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeep.sol";

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

/// @notice Implementation of ECDSAKeep interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is IBondedECDSAKeep {
    using SafeMath for uint256;

    bytes publicKey;
    bool success;
    uint256 memberCount = 3;
    bool isActive = true;
    uint256 bondAmount = 10000;

    // Notification that the keep was requested to sign a digest.
    event SignatureRequested(bytes32 digest);

    // Define fallback function so the contract can accept ether.
    function() external payable {}

    // Functions implemented for IBondedECDSAKeep interface.

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }

    function checkBondAmount() external view returns (uint256) {
        return bondAmount;
    }

    function sign(bytes32 _digest) external onlyWhenActive {
        emit SignatureRequested(_digest);
    }

    function distributeETHReward() external payable {
        uint256 dividend = msg.value.div(memberCount);

        require(dividend > 0, "Dividend value must be non-zero");
    }

    function distributeERC20Reward(address _asset, uint256 _value) external {
        uint256 dividend = _value.div(memberCount);

        require(dividend > 0, "Dividend value must be non-zero");
    }

    function seizeSignerBonds() external onlyWhenActive {
        isActive = false;

        if (address(this).balance > 0) {
            msg.sender.transfer(address(this).balance);
        }
    }

    function returnPartialSignerBonds() external payable {
        uint256 bondPerMember = msg.value.div(memberCount);

        require(bondPerMember > 0, "Partial signer bond must be non-zero");
    }

    function submitSignatureFraud(
        uint8,
        bytes32,
        bytes32,
        bytes32,
        bytes calldata
    ) external returns (bool) {
        require(success, "Signature is not fraudulent");

        return success;
    }

    function closeKeep() external onlyWhenActive {
        isActive = false;
    }

    modifier onlyWhenActive() {
        require(isActive, "Keep is not active");
        _;
    }

    // Functions to set data for tests.

    function reset() public {
        success = true;
        isActive = true;
    }

    function setPublicKey(bytes memory _publicKey) public {
        publicKey = _publicKey;
    }

    function setSuccess(bool _success) public {
        success = _success;
    }

    function setMemberCount(uint256 _memberCount) public {
        memberCount = _memberCount;
    }

    function setBondAmount(uint256 _bondAmount) public {
        bondAmount = _bondAmount;
    }

    function burnContractBalance() public {
        address(0).transfer(address(this).balance);
    }

    function pushFundsFromKeep(address payable _depositAddress) public payable {
        require(msg.value > 0, "value must be greater than 0");
        _depositAddress.transfer(msg.value);
    }
}
