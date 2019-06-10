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
        let tbtcExchange = await (await tbtcSystem.uniswapFactory()).getExchange(tbtc)

        console.log(tbtcSystem)
        console.log(tbtc)
        console.log(tbtcExchange)
    })

    describe('UniswapExchange TBTC', () => {
        it('has no liquidity by default', async () => {
            // getTokenToEthInputPrice
        })
        it('trades with tbtc', async () => {

        })
    })

})