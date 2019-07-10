const TBTC = artifacts.require('TBTC')
const TBTCSystem = artifacts.require('TBTCSystem')
const UniswapDeployment = artifacts.require('UniswapDeployment')
const IUniswapFactory = artifacts.require('IUniswapFactory')

module.exports = async (deployer, network, accounts) => {
  const tbtcSystem = await TBTCSystem.deployed()
  const tbtc = await TBTC.deployed()
  const uniswapDeployment = await UniswapDeployment.deployed()
  const uniswapFactory = await IUniswapFactory.at(
    await uniswapDeployment.factory.call()
  )

  await uniswapFactory.createExchange(tbtc.address)

  await tbtcSystem.setExteroriorAddresses(
    uniswapFactory.address,
    tbtc.address
  )
}
