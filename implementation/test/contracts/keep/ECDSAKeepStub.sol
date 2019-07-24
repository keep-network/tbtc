pragma solidity ^0.5.10;

import {ECDSAKeepContract} from '../../../contracts/interfaces/KeepBridge.sol';

/// @notice Implementation of ECDSAKeepContract interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepStub is ECDSAKeepContract {
    bytes publicKey;

    // Notification that the keep was requested to sign a digest.
    event SignatureRequested(
        bytes32 digest
    );

    function setPublicKey(bytes memory _publicKey) public {
        publicKey = _publicKey;
    }

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }

    function sign(bytes32 _digest) external {
          emit SignatureRequested(_digest);
    }
}
