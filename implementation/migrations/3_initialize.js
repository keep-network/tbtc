const TBTCSystem = artifacts.require('TBTCSystem')
const BTCETHPriceFeed = artifacts.require('BTCETHPriceFeed')

const {
  KeepRegistryAddress,
  BTCUSDPriceFeed,
  ETHUSDPriceFeed,
} = require('./externals')

module.exports = async function(deployer, network) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == 'test' && !process.env.INTEGRATION_TEST) return

  // system
  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(KeepRegistryAddress)

  // price feed
  const btcEthPriceFeed = await BTCETHPriceFeed.deployed()
  if (network === 'mainnet') {
    await btcEthPriceFeed.initialize(BTCUSDPriceFeed, ETHUSDPriceFeed)
  } else {
    // TODO: initialize with our mock price feeds.
  }
}
