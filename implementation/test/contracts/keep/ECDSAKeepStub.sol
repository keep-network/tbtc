pragma solidity ^0.5.10;

import {IECDSAKeep} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeep.sol";
import {IBondedECDSAKeep} from "../../../contracts/external/IBondedECDSAKeep.sol";

/// @notice Implementation of ECDSAKeep interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is IECDSAKeep, IBondedECDSAKeep {
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

    function distributeETHToMembers() external payable {

    }

    function distributeERC20ToMembers(address _asset, uint256 _value) external {

    }


    // Functions implemented for IBondedECDSAKeep interface.

    function submitSignatureFraud(
        address,
        uint8,
        bytes32,
        bytes32,
        bytes32,
        bytes calldata
    ) external returns (bool){
       return success;
    }

    function checkBondAmount(address) external view returns (uint256){
        return bondAmount;
    }

    function seizeSignerBonds(address) external returns (bool){
        if (address(this).balance > 0) {
            msg.sender.transfer(address(this).balance);
        }
        return true;
    }
}
