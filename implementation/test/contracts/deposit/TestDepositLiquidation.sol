pragma solidity 0.4.25;

import "./TestDeposit.sol";

contract TestDepositLiquidation is TestDeposit {
    function startLiquidation() public {
        return self.startLiquidation();
    }

    function startSignerFraudLiquidation() public {
        return self.startSignerFraudLiquidation();
    }

    function startSignerAbortLiquidation() public {
        return self.startSignerAbortLiquidation();
    }
}