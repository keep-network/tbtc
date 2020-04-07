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

    IMedianizer[] private btcEthFeeds;

    constructor() public {
    // solium-disable-previous-line no-empty-blocks
    }

    /// @notice Initialises the addresses of the BTC/USD and ETH/USD price feeds.
    /// @param _BTCETHPriceFeed The BTC/USD price feed address.
    function initialize(
        address _tbtcSystemAddress,
        IMedianizer _BTCETHPriceFeed
    )
        external onlyOwner
    {
        require(!_initialized, "Already initialized.");
        tbtcSystemAddress = _tbtcSystemAddress;
        btcEthFeeds.push(_BTCETHPriceFeed);
        _initialized = true;
    }

    /// @notice Get the current price of bitcoin in ether.
    /// @dev This does not account for any 'Flippening' event.
    /// @return The price of one satoshi in wei.
    function getPrice()
        external view returns (uint256)
    {
        bool btcEthActive;
        uint256 btcEth;

        for(uint i = 0; i < btcEthFeeds.length; i++){
            (btcEth, btcEthActive) = btcEthFeeds[i].peek();
            if(btcEthActive) {
                break;
            }
        }

        require(btcEthActive, "Price feed offline");

        // We typecast down to uint128, because the first 128 bits of
        // the medianizer oracle value is unrelated to the price.
        return uint256(uint128(btcEth));
    }

    /// @notice Get the first active Medianizer contract from the btcEthFeeds array.
    /// @return The address of the first Active Medianizer. address(0) if none found
    function getWorkingBtcEthFeed() external view returns (address){
        bool btcEthActive;

        for(uint i = 0; i < btcEthFeeds.length; i++){
            (, btcEthActive) = btcEthFeeds[i].peek();
            if(btcEthActive) {
                return address(btcEthFeeds[i]);
            }
        }
        return address(0);
    }

    /// @notice Add _btcEthFeed to internal btcEthFeeds array.
    /// @dev IMedianizer must be active in order to add.
    function addBtcEthFeed(IMedianizer _btcEthFeed) external onlyTbtcSystem {
        bool btcEthActive;
        (, btcEthActive) = _btcEthFeed.peek();
        require(btcEthActive, "Cannot add inactive feed");
        btcEthFeeds.push(_btcEthFeed);
    }

    /// @notice Function modifier ensures modified function is only called by tbtcSystemAddress.
    modifier onlyTbtcSystem(){
        require(msg.sender == tbtcSystemAddress, "Caller must be tbtcSystem contract");
        _;
    }
}