pragma solidity ^0.5.10;

import {IBondedECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";

/// @notice Implementation of ECDSAKeepVendor interface used in tests only
/// @dev This is a stub used in tests, so we don't have to call actual ECDSAKeep
contract ECDSAKeepVendorStub is IBondedECDSAKeepVendor {
     address payable factory;

     constructor(address payable _factory) public {
       factory = _factory;
     }

     function selectFactory() public view returns (address payable) {
       return factory;
     }
}
