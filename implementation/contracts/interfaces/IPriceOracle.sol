pragma solidity 0.4.25;

/// @title Bitcoin-Ether price oracle interface.
contract IPriceOracle {
    event PriceUpdated(uint128 price, uint256 priceExpiryTime);

    /// @notice Get the current price of bitcoin in ether
    /// @dev reverts if the price is invalid due to expiry
    /// @return the price in terms of x wei : 1 satoshi
    function getPrice() external view returns (uint128);

    /// @notice Updates the price
    /// @param newPrice the new price in terms of x wei : 1 satoshi
    function updatePrice(uint128 newPrice) external;
}