pragma solidity 0.4.25;

import {ECDSAKeepContract} from '../../../contracts/interfaces/KeepBridge.sol';

/// @notice Implementation of ECDSAKeepContract interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is ECDSAKeepContract {
    bytes publicKey;

    function setPublicKey(bytes _publicKey) public {
        publicKey = _publicKey;
    }

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }
}
