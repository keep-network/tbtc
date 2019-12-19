pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../external/IMedianizer.sol";
import "../interfaces/IBTCETHPriceFeed.sol";

/// @notice Bitcoin-Ether price feed.
/// @dev Based on the ratio of two medianizer price feeds, BTC/USD and ETH/USD.
contract BTCETHPriceFeed is Ownable, IBTCETHPriceFeed {
    using SafeMath for uint256;

    bool private _initialized = false;

    IMedianizer btcPriceFeed;
    IMedianizer ethPriceFeed;

    constructor() public {
    }

    /// @notice Initialises the addresses of the BTC/USD and ETH/USD price feeds.
    /// @param _BTCUSDPriceFeed The BTC/USD price feed address.
    /// @param _ETHUSDPriceFeed The ETH/USD price feed address.
    function initialize(
        address _BTCUSDPriceFeed,
        address _ETHUSDPriceFeed
    )
        external onlyOwner
    {
        require(!_initialized, "Already initialized.");

        btcPriceFeed = IMedianizer(_BTCUSDPriceFeed);
        ethPriceFeed = IMedianizer(_ETHUSDPriceFeed);
        _initialized = true;
    }

    /// @notice Get the current price of bitcoin in ether.
    /// @dev This does not account for any 'Flippening' event.
    /// @return The price of one satoshi in wei.
    function getPrice()
        external view returns (uint256)
    {
        // We typecast down to uint128, because the first 128 bits of
        // the medianizer oracle value is unrelated to the price.
        uint256 btcUsd = uint256(uint128(btcPriceFeed.read()));
        uint256 ethUsd = uint256(uint128(ethPriceFeed.read()));
        return btcUsd.div(ethUsd);
    }
}