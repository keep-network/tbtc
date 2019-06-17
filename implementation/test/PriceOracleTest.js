const PriceOracleV1 = artifacts.require('PriceOracleV1')

import expectThrow from './helpers/expectThrow';
const BN = require('bn.js')

contract('PriceOracleV1', function (accounts) {
    const DEFAULT_OPERATOR = accounts[0];

    // 1 wei = 1 satoshi
    const DEFAULT_PRICE = new BN(1);

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
        })
    
        describe("#updatePrice", () => {
            it('sets new price', async () => {
                const NEW_PRICE = new BN(2)

                let instance = await PriceOracleV1.new(
                    DEFAULT_OPERATOR,
                    DEFAULT_PRICE
                )

                await instance.updatePrice(2)

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
        })
    })
})