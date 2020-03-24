pragma solidity ^0.5.10;

/// @notice A medianizer price feed.
/// @dev Based off the MakerDAO medianizer (https://github.com/makerdao/median)
interface IMedianizer {
    /// @notice Get the current price.
    /// @dev May revert if caller not whitelisted.
    /// @return Price (USD) with 18 decimal places.
    function read() external view returns (uint256);
}