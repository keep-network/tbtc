pragma solidity ^0.5.10;

import {ECDSAKeepVendor} from '../../../contracts/interfaces/KeepBridge.sol';
import {ECDSAKeepStub} from './ECDSAKeepStub.sol';

/// @notice Implementation of ECDSAKeepVendor interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepVendorStub is ECDSAKeepVendor {
     address public keepOwner;

    function openKeep(
        uint256 _groupSize,
        uint256 _honestThreshold,
        address _owner
    ) external payable returns (address _keepAddress) {
        keepOwner = _owner;
    }
}
