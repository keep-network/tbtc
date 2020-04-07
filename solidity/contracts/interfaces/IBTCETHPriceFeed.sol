pragma solidity ^0.5.10;
import "../external/IMedianizer.sol";

interface IBTCETHPriceFeed {
    /// @notice Get the current price of bitcoin in ether.
    /// @dev This does not account for any 'Flippening' event.
    /// @return The price of one satoshi in wei.
    function getPrice() external view returns (uint256);

    /// @notice add a new BTC/ETH meidanizer to the internal btcEthFeeds array
    function addBtcEthFeed(IMedianizer _btcEthFeed) external;
}
