pragma solidity ^0.5.10;

import {IKeep} from '../../../contracts/interfaces/IKeep.sol';

contract KeepStub {

    bool success = true;
    uint256 bondAmount = 10000;
    address keepAddress = address(7);
    bytes pubkey = hex"00";

    function () payable external {}

    function setPubkey(bytes memory _pubkey) public {pubkey = _pubkey;}
    function setSuccess(bool _success) public {success = _success;}
    function setBondAmount(uint256 _bondAmount) public {bondAmount = _bondAmount;}
    function setKeepAddress(address _id) public {keepAddress = _id;}

    function submitSignatureFraud(
        address _keepAddress,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes calldata _preimage
    ) external returns (bool _isFraud) {
        _keepAddress; _v; _r; _s; _signedDigest; _preimage; success = success;
        _isFraud = success;
    }

    function distributeEthToKeepGroup(address _keepAddress) external payable returns (bool) {
        _keepAddress;
        return success;
    }

    function distributeERC20ToKeepGroup(address _keepAddress, address _asset, uint256 _value) external returns (bool) {
        _keepAddress; _asset; _value; success = success;
        return success;
    }

    function checkBondAmount(address _keepAddress) external view returns (uint256) {
        _keepAddress;
        return bondAmount;
    }

    function seizeSignerBonds(address _keepAddress) external returns (bool) {
        _keepAddress;
        if (address(this).balance > 0) {
            msg.sender.transfer(address(this).balance);
        }
        return true;
    }

    function burnContractBalance() public {
        address(0).transfer(address(this).balance);
    }
}
