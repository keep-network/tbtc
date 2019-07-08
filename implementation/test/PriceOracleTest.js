const PriceOracleV1 = artifacts.require('PriceOracleV1')

import expectThrow from './helpers/expectThrow'
import increaseTime from './helpers/increaseTime'
const BN = require('bn.js')


contract('PriceOracleV1', function(accounts) {
  const DEFAULT_OPERATOR = accounts[0]

  // 1 wei = 1 satoshi
  const DEFAULT_PRICE = new BN(1)

  const PRICE_EXPIRY_S = 21600

  describe('#constructor', () => {
    it('deploys', async () => {
      const instance = await PriceOracleV1.deployed()
      assert(instance.address.length == 42)
    })
  })

  describe('Methods', () => {
    describe('#getPrice', () => {
      it('returns a default price', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        const res = await instance.getPrice.call()
        assert(res.eq(DEFAULT_PRICE))
      })

      it('fails when the price is expired', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        const res = await instance.getPrice.call()
        assert(res.eq(DEFAULT_PRICE))

        await increaseTime(PRICE_EXPIRY_S + 1)

        await expectThrow(instance.getPrice.call())
      })
    })

    describe('#updatePrice', () => {
      it('sets new price', async () => {
        const NEW_PRICE = new BN(2)

        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        await instance.updatePrice(2)

        const block = await web3.eth.getBlock('latest')
        const expiry = await instance.expiry.call()

        assert(
          expiry.toNumber() == (block.timestamp + PRICE_EXPIRY_S)
        )

        const res = await instance.getPrice.call()
        assert(res.eq(NEW_PRICE))
      })

      it('fails when price delta < 1%', async () => {
        const price1 = new BN('1000000')
        const price2 = new BN('1000001')

        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          price1
        )

        await expectThrow(instance.updatePrice(price2))
      })

      it('fails when msg.sender != operator', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        await expectThrow(
          instance.updatePrice(
            new BN('1000000'),
            { from: accounts[1] }
          )
        )
      })

      it('ignores 1% threshold for update, when the price is close to expiry', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )
        const price1 = new BN('1000000')
        const price2 = new BN('1000001')

        await instance.updatePrice(price1)
        await expectThrow(
          instance.updatePrice(price2)
        )

        const ONE_HOUR = 3600
        await increaseTime(PRICE_EXPIRY_S - ONE_HOUR + 1)
        
        await instance.updatePrice(price2)
        const res = await instance.getPrice.call()
        assert(res.eq(price2))
      })
    })
  })
})
