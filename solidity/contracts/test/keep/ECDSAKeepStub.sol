pragma solidity 0.5.17;

import {
    IBondedECDSAKeep
} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeep.sol";
import {TestTBTCToken} from "../system/TestTBTCToken.sol";

/// @notice Implementation of ECDSAKeep interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is IBondedECDSAKeep {
    bytes publicKey;
    bool success = true;
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

    function sign(bytes32 _digest) external {
        emit SignatureRequested(_digest);
    }

    function distributeETHReward() external payable {
        // solium-disable-previous-line no-empty-blocks
    }

    function distributeERC20Reward(address _asset, uint256 _value) external {
        TestTBTCToken(_asset).transferFrom(msg.sender, address(this), _value);
    }

    function seizeSignerBonds() external {
        if (address(this).balance > 0) {
            msg.sender.transfer(address(this).balance);
        }
    }

    function returnPartialSignerBonds() external payable {
        // solium-disable-previous-line no-empty-blocks
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

    function closeKeep() external {
        // solium-disable-previous-line no-empty-blocks
    }

    // Functions to set data for tests.

    function reset() public {
        success = true;
    }

    function setPublicKey(bytes memory _publicKey) public {
        publicKey = _publicKey;
    }

    function setSuccess(bool _success) public {
        success = _success;
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

    function drain() public {
        address(0).transfer(address(this).balance);
    }
}
