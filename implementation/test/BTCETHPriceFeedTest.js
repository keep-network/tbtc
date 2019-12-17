const BTCETHPriceFeed = artifacts.require('BTCETHPriceFeed')
const MockMedianizer = artifacts.require('MockMedianizer')

import { createSnapshot, restoreSnapshot } from './helpers/snapshot'
const BN = require('bn.js')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

contract('BTCETHPriceFeed', (accounts) => {
  let btcEthPriceFeed
  let btc
  let eth

  before(async () => {
    btcEthPriceFeed = await BTCETHPriceFeed.new()
    btc = await MockMedianizer.new()
    eth = await MockMedianizer.new()

    await btcEthPriceFeed.initialize(btc.address, eth.address)
  })

  describe('#getPrice', async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it('computes a ratio of the two medianizers', async () => {
      await btc.setValue('200')
      await eth.setValue('100')

      const price = await btcEthPriceFeed.getPrice()
      expect(price).to.eq.BN('2')
    })

    it('casts down each medianizer price to lower 128 bits', async () => {
      const btcPrice = '0xdef00000000000000000000000000000000000000000000000000000000004'
      const ethPrice = '0xabc00000000000000000000000000000000000000000000000000000000002'

      await btc.setValue(btcPrice)
      await eth.setValue(ethPrice)

      const price = await btcEthPriceFeed.getPrice()
      expect(price).to.eq.BN('2')
    })
  })
})
