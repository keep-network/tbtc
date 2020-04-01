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
    address internal tbtcSystemAddress;

    IMedianizer[] private btcUsdFeeds;
    IMedianizer[] private ethUsdFeeds;

    constructor() public {
    // solium-disable-previous-line no-empty-blocks
    }

    /// @notice Initialises the addresses of the BTC/USD and ETH/USD price feeds.
    /// @param _BTCUSDPriceFeed The BTC/USD price feed address.
    /// @param _ETHUSDPriceFeed The ETH/USD price feed address.
    function initialize(
        address _tbtcSystemAddress,
        IMedianizer _BTCUSDPriceFeed,
        IMedianizer _ETHUSDPriceFeed
    )
        external onlyOwner
    {
        require(!_initialized, "Already initialized.");
        tbtcSystemAddress = _tbtcSystemAddress;
        btcUsdFeeds.push(_BTCUSDPriceFeed);
        ethUsdFeeds.push(_ETHUSDPriceFeed);
        // btcPriceFeed = IMedianizer(_BTCUSDPriceFeed);
        // ethPriceFeed = IMedianizer(_ETHUSDPriceFeed);
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
        (uint256 btcUsd, uint256 ethUsd) = _getComponentPrices();
        // The price is a ratio of bitcoin to ether is expressed as:
        //  x btc : y eth
        // Bitcoin has 10 decimal places, ether has 18. Normalising the units, we have:
        //  x * 10^8 : y * 10^18
        // Simplfying down, we can express it as:
        //  x : y * 10^10
        // Due to order-of-ops, we can move the multiplication to get some more precision.
        return btcUsd.mul(10**10).div(ethUsd);
    }

    /// @notice Get the first active Medianizer contract from the btcUsdFeeds array.
    /// @return The address of the first Active Medianizer. address(0) if none found
    function getWorkingBtcUsdFeed() external view returns (address){
        bool btcUsdActive;

        for(uint i = 0; i < btcUsdFeeds.length; i++){
            (, btcUsdActive) = btcUsdFeeds[i].peek();
            if(btcUsdActive) {
                return address(btcUsdFeeds[i]);
            }
        }
        return address(0);
    }

    /// @notice Get the first active Medianizer contract from the ethUsdFeeds array.
    /// @return The address of the first Active Medianizer. address(0) if none found
    function getWorkingEthUsdFeed() external view returns (address){
        bool ethUsdActive;

        for(uint i = 0; i < ethUsdFeeds.length; i++){
            (, ethUsdActive) = ethUsdFeeds[i].peek();
            if(ethUsdActive) {
                return address(ethUsdFeeds[i]);
            }
        }
        return address(0);
    }

    /// @notice Add _btcUsdFeed to internal btcUsdFeeds array.
    /// @dev IMedianizer must be active in order to add.
    function addBtcUsdFeed(IMedianizer _btcUsdFeed) external onlyTbtcSystem {
        bool btcUsdActive;
        (, btcUsdActive) = _btcUsdFeed.peek();
        require(btcUsdActive, "Cannot add inactive feed");
        btcUsdFeeds.push(_btcUsdFeed);
    }

    /// @notice Add _ethUsdFeed to internal btcUsdFeeds array.
    /// @dev IMedianizer must be active in order to add.
    function addEthUsdFeed(IMedianizer _ethUsdFeed) external onlyTbtcSystem {
        bool ethUsdActive;
        (, ethUsdActive) = _ethUsdFeed.peek();
        require(ethUsdActive, "Cannot add inactive feed");
        ethUsdFeeds.push(_ethUsdFeed);
    }

    /// @notice Get the current value of BTCUSD and ETHUSD price feeds.
    /// @dev sequentially traverses the `btcUsdFeeds` and `ethUsdFeeds` arrays
    /// and will revert if no price feed is active from either set.
    /// @return The BTCUSD and ETHUSD price feed values.
    function _getComponentPrices()
        internal view returns (uint256 btcUsd, uint256 ethUsd)
    {
        bool btcUsdActive;
        bool ethUsdActive;

        for(uint i = 0; i < btcUsdFeeds.length; i++){
            (btcUsd, btcUsdActive) = btcUsdFeeds[i].peek();
            if(btcUsdActive) {
                break;
            }
        }

        for(uint i = 0; i < ethUsdFeeds.length; i++){
            (ethUsd, ethUsdActive) = ethUsdFeeds[i].peek();
            if(ethUsdActive) {
                break;
            }
        }

        require(btcUsdActive && ethUsdActive, "Price feed offline");
    }

    /// @notice Function modifier ensures modified function is only called by tbtcSystemAddress.
    modifier onlyTbtcSystem(){
        require(msg.sender == tbtcSystemAddress, "Caller must be tbtcSystem contract");
        _;
    }
}