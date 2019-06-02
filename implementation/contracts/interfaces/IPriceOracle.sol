pragma solidity 0.4.25;

/**
 * @title Bitcoin-Ether price oracle interface.
 */

interface IPriceOracle {
    function getPrice() external view returns (uint128);
    function updatePrice(uint128 price) external;
}