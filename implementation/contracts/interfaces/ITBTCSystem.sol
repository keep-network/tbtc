pragma solidity 0.4.25;

/**
 * @title Keep interface
 */

interface ITBTCSystem {

    // Sets up the TBTC system
    // - uniswap factory initialisation
    // - TBTC uniswap exchange creation
    function setup(
        address _uniswapFactory,
        address _tbtc
    ) external;

    // expected behavior:
    // return the price of 1 sat in wei
    // these are the native units of the deposit contract
    function fetchOraclePrice() external view returns (uint256);

    // passthrough requests for the oracle
    function fetchRelayCurrentDifficulty() external view returns (uint256);
    function fetchRelayPreviousDifficulty() external view returns (uint256);

    function getTBTCUniswapExchange() external view returns (address);
}
