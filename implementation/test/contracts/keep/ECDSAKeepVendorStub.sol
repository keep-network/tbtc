pragma solidity ^0.5.10;

import {IECDSAKeepVendor} from "keep-tecdsa/solidity/contracts/api/IECDSAKeepVendor.sol";

/// @notice Implementation of ECDSAKeepVendor interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepVendorStub is IECDSAKeepVendor {
     address public keepOwner;

    function openKeep(
        uint256 _groupSize,
        uint256 _honestThreshold,
        address _owner
    ) external payable returns (address _keepAddress) {
        keepOwner = _owner;
        _keepAddress = address(888);
    }
}
