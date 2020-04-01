pragma solidity ^0.5.10;
import "../external/IMedianizer.sol";

interface IBTCETHPriceFeed {
    /// @notice Get the current price of bitcoin in ether.
    /// @dev This does not account for any 'Flippening' event.
    /// @return The price of one satoshi in wei.
    function getPrice() external view returns (uint256);

    /// @notice add a new BTC/USD meidanizer to the internal btcUsdFeeds array
    function addBtcUsdFeed(IMedianizer _btcUsdFeed) external;

    /// @notice add a new ETH/USD meidanizer to the internal ethUsdFeeds array
    function addEthUsdFeed(IMedianizer _ethUsdFeed) external;
}
