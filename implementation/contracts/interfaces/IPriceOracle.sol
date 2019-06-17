pragma solidity 0.4.25;

/**
 * @title Bitcoin-Ether price oracle interface.
 */

contract IPriceOracle {
    event PriceUpdated(uint128 price, uint256 zzz);
    
    function getPrice() external view returns (uint128);
    function updatePrice(uint128 price) external;
}