const PriceOracleV1 = artifacts.require('PriceOracleV1')

import expectThrow from './helpers/expectThrow';
import increaseTime, { duration, increaseTimeTo } from './helpers/increaseTime';
const BN = require('bn.js')


contract('PriceOracleV1', function (accounts) {
    const DEFAULT_OPERATOR = accounts[0];

    // 1 wei = 1 satoshi
    const DEFAULT_PRICE = new BN(1);
    
    const PRICE_EXPIRY_S = 21600;

    describe("#constructor", () => {
        it('deploys', async () => {
            let instance = await PriceOracleV1.deployed()
            assert(instance.address.length == 42)
        })
    })

    describe("Methods", () => {

        describe("#getPrice", () => {
            it('returns a default price', async () => {
                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    DEFAULT_PRICE
                )

                let res = await instance.getPrice.call()
                assert(res.eq(DEFAULT_PRICE))
            })

            it('fails when the price is expired', async () => {
                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    DEFAULT_PRICE
                )

                let res = await instance.getPrice.call()
                assert(res.eq(DEFAULT_PRICE))

                await increaseTime(PRICE_EXPIRY_S + 1)

                await expectThrow(instance.getPrice.call())
            })
        })
    
        describe("#updatePrice", () => {
            it('sets new price', async () => {
                const NEW_PRICE = new BN(2)
                
                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    DEFAULT_PRICE
                )
                    
                await instance.updatePrice(2)

                const block = await web3.eth.getBlock('latest')
                let zzz = await instance.zzz.call()

                assert(
                    zzz.toNumber() == (block.timestamp + PRICE_EXPIRY_S)
                );
                
                let res = await instance.getPrice.call()
                assert(res.eq(NEW_PRICE))
            })

            it('fails when price delta < 1%', async () => {
                const price1 = new BN('1000000')
                const price2 = new BN('1000001')

                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    price1
                )

                await expectThrow(instance.updatePrice(price2))
            })

            it('fails when msg.sender != operator', async () => {
                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    DEFAULT_PRICE
                )

                await expectThrow(
                    instance.updatePrice(
                        new BN('1000000'),
                        { from: accounts[1] }
                    )
                );
            })
        })
    })
})