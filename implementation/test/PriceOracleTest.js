import increaseTime from './helpers/increaseTime'
import BN from 'bn.js'

const PriceOracleV1 = artifacts.require('PriceOracleV1')

contract('PriceOracleV1', function(accounts) {
  const DEFAULT_OPERATOR = accounts[0]

  // 1 satoshi = 1 wei
  const DEFAULT_PRICE = new BN('323200000000')

  const PRICE_EXPIRY_SECONDS = 21600 // 6 hours = 21600 seconds

  describe('#constructor', () => {
    it('deploys', async () => {
      const instance = await PriceOracleV1.deployed()
      assert(instance.address.length == 42)

      const price = await instance.getPrice()
      assert(price.eq(DEFAULT_PRICE))

      const operator = await instance.operator()
      assert(operator == accounts[0])

      const expiry = new BN(await instance.expiry())
      const timestamp = new BN((await web3.eth.getBlock('latest')).timestamp)
      // using an error margin of seconds here, since idk how to get the Truffle deployment time
      const expiry_errorMargin = timestamp.add(new BN(PRICE_EXPIRY_SECONDS)).sub(expiry).abs() // eslint-disable-line
      assert(expiry_errorMargin.lt(new BN('20')), 'unexpected expiry value')
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

        await increaseTime(PRICE_EXPIRY_SECONDS + 1)

        try {
          await instance.getPrice.call()

          assert(false, 'Test call did not error as expected')
        } catch (e) {
          assert.include(e.message, 'Price expired')
        }
      })
    })

    describe('#updatePrice', () => {
      it('sets new price', async () => {
        const newPrice = DEFAULT_PRICE.mul(new BN('0.02'))

        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        await instance.updatePrice(newPrice)

        const block = await web3.eth.getBlock('latest')
        const expiry = await instance.expiry.call()

        assert(
          expiry.toNumber() == (block.timestamp + PRICE_EXPIRY_SECONDS)
        )

        const res = await instance.getPrice.call()
        assert(res.eq(newPrice))
      })

      it('fails when price delta < 1%', async () => {
        const delta = new BN('10')
        const price1 = new BN('323200000001')
        const price2 = new BN('323200000010').add(delta)
        const price3 = new BN('323200000000').sub(delta)

        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          price1
        )

        try {
          await instance.updatePrice(price2)
          await instance.updatePrice(price3)

          assert(false, 'Test call did not error as expected')
        } catch (e) {
          assert.include(e.message, 'Price change is negligible (<1%)')
        }
      })

      it('fails when msg.sender != operator', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          DEFAULT_PRICE
        )

        try {
          await instance.updatePrice(
            new BN('323200000001'),
            { from: accounts[1] }
          )

          assert(false, 'Test call did not error as expected')
        } catch (e) {
          assert.include(e.message, 'Unauthorised')
        }
      })

      it('ignores 1% threshold for update, when the price is close to expiry', async () => {
        const instance = await PriceOracleV1.new(
          DEFAULT_OPERATOR,
          new BN('1')
        )
        const price1 = new BN('323200000001')
        const price2 = new BN('323200000010')

        await instance.updatePrice(price1)

        try {
          await instance.updatePrice(price2)
          assert(false, 'Test call did not error as expected')
        } catch (e) {
          assert.include(e.message, 'Price change is negligible (<1%)')
        }

        const ONE_HOUR = 3600
        await increaseTime(PRICE_EXPIRY_SECONDS - ONE_HOUR)

        await instance.updatePrice(price2)
        const res = await instance.getPrice.call()
        assert(res.eq(price2))
      })
    })
  })
})
