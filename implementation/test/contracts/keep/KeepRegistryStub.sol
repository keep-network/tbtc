pragma solidity ^0.5.10;

import {IKeepRegistry} from "@keep-network/keep-tecdsa/contracts/api/IKeepRegistry.sol";

/// @notice Implementation of KeepRegistry interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract KeepRegistryStub is IKeepRegistry {
    address vendor;

    function setVendor(address _vendorAddress) public {
        vendor = _vendorAddress;
    }

    function getVendor(string calldata _keepType) external view returns (address){
        return vendor;
    }
}
