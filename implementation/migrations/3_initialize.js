const TBTCSystem = artifacts.require('TBTCSystem')

const {
  KeepRegistryAddress,
} = require('./externals')

module.exports = async function(deployer) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == 'test' && !process.env.INTEGRATION_TEST) return

  // system
  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(KeepRegistryAddress)
}
