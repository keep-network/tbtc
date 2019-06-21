const TBTC = artifacts.require('TBTC')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const UniswapDeployment = artifacts.require('UniswapDeployment')

module.exports = async (deployer) => {
    const tbtcSystem = await TBTCSystemStub.deployed();
    const tbtc = await TBTC.deployed()
    const uniswapDeployment = await UniswapDeployment.deployed()

    await tbtcSystem.setup(
        await uniswapDeployment.factory(),
        tbtc.address
    );
}