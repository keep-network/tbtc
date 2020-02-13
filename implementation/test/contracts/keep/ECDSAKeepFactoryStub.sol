pragma solidity ^0.5.10;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

contract ECDSAKeepFactoryStub is IBondedECDSAKeepFactory {
     address public keepOwner;
     address public keepAddress = address(888);
     uint256 feeEstimate = 123456;

    function openKeep(
        uint256,
        uint256,
        address _owner,
        uint256
    ) external payable returns (address) {
        require(msg.value >= feeEstimate, "Insufficient value for new keep creation");
        keepOwner = _owner;
        return keepAddress;
    }

    function openKeepFeeEstimate() external returns (uint256){
        return feeEstimate;
    }

    function setOpenKeepFeeEstimate(uint256 _fee) external {
        feeEstimate = _fee;
    }

}