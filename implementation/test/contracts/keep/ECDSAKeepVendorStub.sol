pragma solidity ^0.5.10;

import {IBondedECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";
import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

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

contract ECDSAKeepFactoryStub is IBondedECDSAKeepFactory {
     address public keepOwner;
     address public keepAddress = address(888);

    function openKeep(
        uint256,
        uint256,
        address _owner,
        uint256
    ) external payable returns (address) {
        keepOwner = _owner;
        return keepAddress;
    }
}
