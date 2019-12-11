pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../interfaces/IMedianizer.sol";

contract BTCETHPriceFeed is Ownable {
    using SafeMath for uint256;

    bool private _initialized = false;

    IMedianizer btcPriceFeed;
    IMedianizer ethPriceFeed;

    constructor() public {
    }

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