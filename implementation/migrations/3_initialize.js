const TBTCToken = artifacts.require('TBTCToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const TBTCSystem = artifacts.require('TBTCSystem')

const {
  KeepRegistryAddress,
  UniswapFactoryAddress,
} = require('./externals')

module.exports = async function(deployer) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == 'test' && !process.env.INTEGRATION_TEST) return

  // Uniswap
  const tbtcToken = await TBTCToken.deployed()
  const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)

  let tbtcExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
  if (tbtcExchangeAddress == '0x0000000000000000000000000000000000000000') {
    await uniswapFactory.createExchange(tbtcToken.address)
    tbtcExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
  }

  // system
  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(KeepRegistryAddress, tbtcExchangeAddress)
}
