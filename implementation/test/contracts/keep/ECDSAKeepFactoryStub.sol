pragma solidity ^0.5.10;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

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