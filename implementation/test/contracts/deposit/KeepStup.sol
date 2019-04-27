pragma solidity 0.4.25;

import {IKeep} from '../../../contracts/interfaces/IKeep.sol';

contract KeepStub {

    mapping (bytes32 => uint256) approved;
    bool success = true;
    uint256 bondAmount = 10000;

    function () payable {}

    function setIsFraud(bool _success) public {
        success = _success;
    }

    function setBondAmount(uint256 _bondAmount) public {
        bondAmount = _bondAmount;
    }

    function wasDigestApprovedForSigning(uint256 _keepID, bytes32 _digest) external view returns (uint256) {
        _keepID;
        return approved[_digest];
    }

    function approveDigest(uint256 _keepID, bytes32 _digest) external returns (bool _success) {
        _keepID;
        approved[_digest] = 100;
        _success = success;
    }

    function submitSignatureFraud(
        uint256 _keepID,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) external returns (bool _isFraud) {
        _keepID; _v; _r; _s; _signedDigest; _preimage; success = success;
        _isFraud = success;
    }


    function distributeEthToKeepGroup(uint256 _keepID) external payable returns (bool) {
        _keepID;
        return success;
    }


    function distributeERC20ToKeepGroup(uint256 _keepID, address _asset, uint256 _value) external returns (bool) {
        _keepID; _asset; _value; success = success;
        return success;
    }


    function requestKeepGroup(uint256 _m, uint256 _n) external payable returns (uint256 _keepID) {
        _m; _n;
        return 7;
    }


    function getKeepPubkey(uint256 _keepID) external view returns (bytes) {
        _keepID; success;
        // this is the pubkey coresponding to 32 '11' bytes
        return hex"4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1";
    }


    function checkBondAmount(uint256 _keepID) external view returns (uint256) {
        _keepID;
        return bondAmount;
    }


    function seizeSignerBonds(uint256 _keepID) external returns (bool) {
        _keepID; success = success;
        msg.sender.transfer(address(this).balance);
        return true;
    }
}
