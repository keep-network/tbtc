const PriceOracleV1 = artifacts.require('PriceOracleV1')

import expectThrow from './helpers/expectThrow';
const BN = require('bn.js')

contract('PriceOracleV1', function () {

    // 1 wei = 1 satoshi
    const DEFAULT_PRICE = new BN(1);

    describe("#constructor", () => {
        it('deploys', async () => {
            let instance = await PriceOracleV1.deployed()
        })
    })

    describe("Methods", () => {

        describe("#getPrice", () => {
            it('returns a default price', async () => {
                let instance = await PriceOracleV1.new(
                    DEFAULT_PRICE
                )

                let res = await instance.getPrice.call()
                assert(res.eq(DEFAULT_PRICE))
            })
        })
    
        describe("#updatePrice", () => {
            it('sets new price', async () => {
                const NEW_PRICE = new BN(2)

                let instance = await PriceOracleV1.new(
                    DEFAULT_PRICE
                )

                await instance.updatePrice(2)

                let res = await instance.getPrice.call()
                assert(res.eq(NEW_PRICE))
            })

            
        })
    })
})