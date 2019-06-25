const TBTC = artifacts.require('TBTC')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const UniswapDeployment = artifacts.require('UniswapDeployment')
const IUniswapFactory = artifacts.require('IUniswapFactory')

const uniswap = require('../uniswap')

module.exports = async (deployer, network, accounts) => {
    const tbtcSystemStub = await TBTCSystemStub.deployed();
    const tbtc = await TBTC.deployed()
    const uniswapDeployment = await UniswapDeployment.deployed()
    const uniswapFactory = await IUniswapFactory.at(
        await uniswapDeployment.factory.call()
    )
    
    await uniswapFactory.createExchange(tbtc.address)

    await tbtcSystemStub.setExteroriorAddresses(
        uniswapFactory.address,
        tbtc.address
    );
}