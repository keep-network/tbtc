const TBTCToken = artifacts.require('TBTCToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const KeepBridge = artifacts.require('KeepBridge')
const TBTCSystem = artifacts.require('TBTCSystem')

const {
  KeepRegistryAddress,
  UniswapFactoryAddress,
} = require('./externals')

module.exports = async function(deployer) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == 'test' && !process.env.INTEGRATION_TEST) return

  // Keep
  const keepBridge = await KeepBridge.deployed()
  await keepBridge.initialize(KeepRegistryAddress)

  // Uniswap
  const tbtcToken = await TBTCToken.deployed()
  const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)

  let tbtcExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
  if (tbtcExchangeAddress == '0x0000000000000000000000000000000000000000') {
    await uniswapFactory.createExchange(tbtcToken.address)
    tbtcExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
  }

  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(tbtcExchangeAddress)
}
