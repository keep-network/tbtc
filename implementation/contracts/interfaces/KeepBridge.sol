pragma solidity 0.4.25;

import {IKeep} from "./IKeep.sol";

contract KeepBridge is IKeep {
    address keepRegistry;

    function wasDigestApprovedForSigning(uint256 _keepID, bytes32 _digest) external view returns (uint256){
        //TODO: Implement
        return 0;
    }


    function approveDigest(uint256 _keepID, bytes32 _digest) external returns (bool _success){
        //TODO: Implement
        return _success;
    }

    function submitSignatureFraud(
        uint256 _keepID,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) external returns (bool _isFraud){
        //TODO: Implement
        return _isFraud;
    }

    function distributeEthToKeepGroup(uint256 _keepID) external payable returns (bool){
        //TODO: Implement
        return false;
    }

    function distributeERC20ToKeepGroup(uint256 _keepID, address _asset, uint256 _value) external returns (bool){
        //TODO: Implement
        return false;
    }

    function requestKeepGroup(uint256 _m, uint256 _n) external payable returns (uint256 _keepID){
        //TODO: Implement

        address keepAddress = KeepRegistryContract(keepRegistry).createECDSAKeep(_n,_m);
        // TODO: keepID type should be changed from uint256 to addrress
        _keepID = uint256(keepAddress);

        return _keepID;
    }

    // get the result of a keep formation
    // should return a 64 byte packed pubkey (x and y)
    // error if not ready yet
    function getKeepPubkey(uint256 _keepID) external view returns (bytes){
        // TODO: keepID type should be changed from uint256 to addrress
        address _keepAddress = address(_keepID);

        return ECDSAKeepContract(_keepAddress).getPublicKey();
    }


    // returns the amount of the keep's ETH bond in wei
    function checkBondAmount(uint256 _keepID) external view returns (uint256){
        //TODO: Implement
        return 0;
    }

    function seizeSignerBonds(uint256 _keepID) external returns (bool){
        //TODO: Implement
        return false;
    }

    //TODO: add: onlyOwner
    function initialize(address _keepRegistry) public {
        keepRegistry = _keepRegistry;
    }
}

/// @notice Interface for communication with `KeepRegistry` contract
/// @dev It allows to call a function without the need of low-level call
interface KeepRegistryContract {

    /// @notice Create a new ECDSA keep
    /// @param _groupSize Number of members in the keep
    /// @param _honestThreshold Minimum number of honest keep members
    /// @return Created keep address
    function createECDSAKeep(
        uint256 _groupSize,
        uint256 _honestThreshold
    ) external payable returns (address keep);
}

/// @notice Interface for communication with `ECDSAKeep` contract
/// @dev It allows to call a function without the need of low-level call
interface ECDSAKeepContract {

    /// @notice Returns the keep signer's public key.
    /// @return Signer's public key.
    function getPublicKey() external view returns (bytes memory);
}
