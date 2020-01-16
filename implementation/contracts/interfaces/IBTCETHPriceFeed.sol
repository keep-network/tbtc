pragma solidity ^0.5.10;

interface IBTCETHPriceFeed {
    /// @notice Get the current price of bitcoin in ether.
    /// @dev This does not account for any 'Flippening' event.
    /// @return The price of one satoshi in wei.
    function getPrice() external view returns (uint256);
}
