pragma solidity ^0.5.10;

/// @notice A medianizer price feed.
/// @dev Based off the MakerDAO medianizer (https://github.com/makerdao/median)
interface IMedianizer {
    function read() external view returns (uint256);
}