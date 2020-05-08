pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

contract FullyBackedECDSAKeepFactoryStub is IBondedECDSAKeepFactory {
    address public keepOwner;
    address public keepAddress = address(444);
    uint256 public stakeLockDuration;
    uint256 feeEstimate = 789789;

    function openKeep(
        uint256,
        uint256,
        address _owner,
        uint256,
        uint256 _stakeLockDuration
    ) external payable returns (address) {
        require(msg.value >= feeEstimate, "Insufficient value for new keep creation");
        keepOwner = _owner;
        stakeLockDuration = _stakeLockDuration;
        return keepAddress;
    }

    function openKeepFeeEstimate() external view returns (uint256) {
        return feeEstimate;
    }

    function setOpenKeepFeeEstimate(uint256 _fee) external {
        feeEstimate = _fee;
    }

    function getSortitionPoolWeight(
        address _application
    ) external view returns (uint256 poolWeight) {
        return 1000;
    }
}
