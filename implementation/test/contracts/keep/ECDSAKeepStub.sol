pragma solidity ^0.5.10;

import {IECDSAKeep} from "keep-tecdsa/solidity/contracts/api/IECDSAKeep.sol";
import {IKeep} from "../../../contracts/external/IKeep.sol";

/// @notice Implementation of ECDSAKeep interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is IECDSAKeep, IKeep {
    bytes publicKey;
    bool success;
    uint256 bondAmount = 10000;

    // Notification that the keep was requested to sign a digest.
    event SignatureRequested(
        bytes32 digest
    );

    // Define fallback function so the contract can accept ether.
    function() external payable {}

    // Functions to set data for tests.

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

    // Functions implemented for IECDSAKeep interface.

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }

    function sign(bytes32 _digest) external {
          emit SignatureRequested(_digest);
    }

    // Functions implemented for IKeep interface.

    function submitSignatureFraud(
        address _keepAddress,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes calldata _preimage
    ) external returns (bool){
       return success;
    }

    function distributeEthToKeepGroup(address _keepAddress) external payable returns (bool){
        return success;
    }

    function distributeERC20ToKeepGroup(address _keepAddress, address _asset, uint256 _value) external returns (bool){
        return success;
    }

    function checkBondAmount(address _keepAddress) external view returns (uint256){
        return bondAmount;
    }

    function seizeSignerBonds(address _keepAddress) external returns (bool){
        if (address(this).balance > 0) {
            msg.sender.transfer(address(this).balance);
        }
        return true;
    }
}
