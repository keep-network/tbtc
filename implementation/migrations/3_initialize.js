const TBTCToken = artifacts.require('TBTCToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const KeepBridge = artifacts.require('KeepBridge')
const TBTCSystem = artifacts.require('TBTCSystem')

// Configuration for addresses of externally deployed smart contracts
const KeepRegistryAddress = '0x622B525445aD939f1b0fF1193E57DC0ED75dAb6e'
const UniswapFactoryAddress = '0xE6766544a13aD18C54B72875A62e61033ac81D47'

module.exports = async function(deployer) {
  // Don't enact this setup during testing.
  if (process.env.NODE_ENV == 'test') return

  // Keep
  const keepBridge = await KeepBridge.deployed()
  await keepBridge.initialize(KeepRegistryAddress)

  // Uniswap
  const tbtcToken = await TBTCToken.deployed()
  const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)
  if (await uniswapFactory.getExchange(tbtcToken.address) == '0x0000000000000000000000000000000000000000') {
    await uniswapFactory.createExchange(tbtcToken.address)
  }

  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(UniswapFactoryAddress)
}
