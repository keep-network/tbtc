pragma solidity 0.4.25;

import "./TestDeposit.sol";

contract TestDepositLiquidation is TestDeposit {
    bool attemptToLiquidateOnchainSuccess = false;

    function setAttemptToLiquidateOnchain(bool _success) public returns (bool) {
        attemptToLiquidateOnchainSuccess = _success;
    }

    function attemptToLiquidateOnchain() public returns (bool) {
        return attemptToLiquidateOnchainSuccess;
    }

    function startSignerFraudLiquidation() public {
        return self.startSignerFraudLiquidation();
    }

    function startSignerAbortLiquidation() public {
        return self.startSignerAbortLiquidation();
    }
}