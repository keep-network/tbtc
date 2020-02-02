pragma solidity ^0.5.10;

import {IECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeepVendor.sol";

/// @notice Implementation of ECDSAKeepVendor interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepVendorStub is IECDSAKeepVendor {
     address public keepOwner;
     address public keepAddress = address(888);

    function openKeep(
        uint256,
        uint256,
        address _owner
    ) external payable returns (address) {
        keepOwner = _owner;

        return keepAddress;
    }
}
