pragma solidity 0.5.17;

import {
    IBondedECDSAKeepVendor
} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";

contract ECDSAKeepVendorStub is IBondedECDSAKeepVendor {
    address payable factory;

    constructor(address payable _factory) public {
        factory = _factory;
    }

    function setFactory(address payable _factory) public {
        factory = _factory;
    }

    function selectFactory() public view returns (address payable) {
        return factory;
    }
}
