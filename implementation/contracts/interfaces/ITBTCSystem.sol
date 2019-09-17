pragma solidity ^0.5.10;

/**
 * @title Keep interface
 */

interface ITBTCSystem {

    // expected behavior:
    // return the price of 1 sat in wei
    // these are the native units of the deposit contract
    function fetchOraclePrice() external view returns (uint256);

    // passthrough requests for the oracle
    function fetchRelayCurrentDifficulty() external view returns (uint256);
    function fetchRelayPreviousDifficulty() external view returns (uint256);

    function getTBTCUniswapExchange() external view returns (address);
}
