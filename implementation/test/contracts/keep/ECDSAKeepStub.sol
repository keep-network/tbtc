pragma solidity ^0.5.10;

import {ECDSAKeep} from '../../../contracts/interfaces/KeepBridge.sol';

/// @notice Implementation of ECDSAKeep interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is ECDSAKeep {
    bytes publicKey;

    function setPublicKey(bytes memory _publicKey) public {
        publicKey = _publicKey;
    }

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }
}
