uniswap
=======

A vendored deployment of the Uniswap contracts, for private testnets.
 
 * `contracts-vyper` contains the pinned Git submodule of the **current** Uniswap mainnet deployment, as of 26 June 2019
 * `migrations/` deploys the Uniswap Factory/Exchange contracts, initialises the Factory, and creates the TBTC Exchange
 * see `UniswapDeploymentTest.js` for more
 * `example.sh` contains some [Seth](https://dev.to/liamzebedee/a-primer-on-seth-solidity-s-swiss-army-knife-mi3) code to get quickly started with playing with Uniswap methods (ie. adding liquidity)