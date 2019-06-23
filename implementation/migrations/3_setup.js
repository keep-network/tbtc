const TBTC = artifacts.require('TBTC')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const UniswapDeployment = artifacts.require('UniswapDeployment')
const IUniswapFactory = artifacts.require('IUniswapFactory')

module.exports = async (deployer) => {
    const tbtcSystem = await TBTCSystemStub.deployed();
    const tbtc = await TBTC.deployed()
    const uniswapDeployment = await UniswapDeployment.deployed()
    const uniswapFactory = await IUniswapFactory.at(
        await uniswapDeployment.factory()
    )

    await uniswapFactory.createExchange(tbtc.address)

    await tbtcSystem.setExteroriorAddresses(
        uniswapFactory.address,
        tbtc.address
    );
}