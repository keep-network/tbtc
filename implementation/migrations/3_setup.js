const TBTCToken = artifacts.require('TBTCToken')
const UniswapDeployment = artifacts.require('UniswapDeployment')
const IUniswapFactory = artifacts.require('IUniswapFactory')

module.exports = async (deployer, network, accounts) => {
  const tbtcToken = await TBTCToken.deployed()
  const uniswapDeployment = await UniswapDeployment.deployed()
  const uniswapFactory = await IUniswapFactory.at(
    await uniswapDeployment.factory.call()
  )

  await uniswapFactory.createExchange(tbtcToken.address)
}
