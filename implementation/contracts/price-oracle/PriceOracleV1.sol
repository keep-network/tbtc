pragma solidity 0.4.25;

import "../interfaces/IPriceOracle.sol";

/**
 * The price oracle implements a simple price feed, managed
 * by a trusted operator.
 */
contract PriceOracleV1 is IPriceOracle {
    constructor() public {
    }

    function getPrice() external view returns (uint128) {
        return 0;
    }

    function updatePrice(uint128 price) public {
        return;
    }
}