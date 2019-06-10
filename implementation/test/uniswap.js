const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTC = artifacts.require('TBTC')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

// Tests the Uniswap deployment

contract('Uniswap deployment', () => {
    let tbtcSystem;
    let tbtc;

    it('deployed uniswap', async () => {
        let tbtcSystem = await TBTCSystemStub.deployed();
        let tbtc = await tbtcSystem.tbtc()
        expect(tbtc).to.not.be.empty;

        let uniswapFactoryAddr = await tbtcSystem.uniswapFactory()
        expect(uniswapFactoryAddr).to.not.be.empty;

        let uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr);
        let tbtcExchangeAddr = await uniswapFactory.getExchange(tbtc)
        expect(tbtcExchangeAddr).to.not.be.empty;
    })

    describe('UniswapExchange TBTC', () => {
        it('has no liquidity by default', async () => {
            // getTokenToEthInputPrice
        })

        it('successfully trades with tbtc', async () => {

        })
    })

})