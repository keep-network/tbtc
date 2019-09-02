uniswap
=======

A vendored deployment of the Uniswap contracts, for private testnets.
 
 * `contracts-vyper` contains the pinned Git submodule of the **current** Uniswap mainnet deployment, as of January 25, 2019 ([c10c08d81d](https://github.com/Uniswap/contracts-vyper/tree/c10c08d81d6114f694baa8bd32f555a40f6264da))
 * `migrations/` deploys the Uniswap Factory/Exchange contracts, initialises the Factory, and creates the TBTC Exchange
 * see `UniswapDeploymentTest.js` for more
 * `example.sh` contains some [Seth](https://dev.to/liamzebedee/a-primer-on-seth-solidity-s-swiss-army-knife-mi3) code to get quickly started with playing with Uniswap methods (ie. adding liquidity)