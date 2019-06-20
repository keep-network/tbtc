pragma solidity 0.4.25;

import {ECDSAKeepContract} from '../../../contracts/interfaces/KeepBridge.sol';

contract ECDSAKeepStub is ECDSAKeepContract {
    bytes publicKey;

    function setPublicKey(bytes _publicKey) public {
        publicKey = _publicKey;
    }

    function getPublicKey() external view returns (bytes memory) {
        return publicKey;
    }
}
