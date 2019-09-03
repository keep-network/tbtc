uniswap
=======

A vendored deployment of the Uniswap contracts, for private testnets.
 
 * `contracts-vyper` contains the pinned Git submodule of the **current** Uniswap mainnet deployment, as of 26 June 2019
 * `migrations/` deploys the Uniswap Factory/Exchange contracts, initialises the Factory, and creates the TBTC Exchange
 * see `UniswapDeploymentTest.js` for more