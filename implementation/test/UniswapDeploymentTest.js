const TBTCToken = artifacts.require('TBTCToken')
const TBTCSystem = artifacts.require('TBTCSystem')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

import { UniswapHelpers } from './helpers/uniswap'
import { createSnapshot, restoreSnapshot } from './helpers/snapshot'

// Tests the Uniswap deployment

contract.only('Uniswap', (accounts) => {
  let tbtcToken

  before(async () => {
    await createSnapshot()
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  after(async () => {
    await restoreSnapshot()
  })

  describe('deployment', async () => {
    it('deployed the uniswap factory and created the exchange', async () => {
      tbtcToken = await TBTCToken.deployed()
      expect(tbtcToken).to.not.be.empty

      const tbtcSystem = await TBTCSystem.deployed()
      const uniswapFactoryAddr = await tbtcSystem.getUniswapFactory()
      expect(uniswapFactoryAddr).to.not.be.empty

      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtcToken.address)
      expect(tbtcExchangeAddr).to.not.be.empty
    })
  })

  describe('tBTC Uniswap Exchange', () => {
    let tbtc
    let tbtcExchange

    beforeEach(async () => {
      tbtc = await TBTCToken.new()

      const tbtcSystem = await TBTCSystem.deployed()
      const uniswapFactoryAddr = await tbtcSystem.getUniswapFactory()
      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)

      await uniswapFactory.createExchange(tbtc.address)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtc.address)
      tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddr)
    })

    it('has no liquidity by default', async () => {
      try {
        await tbtcExchange.getTokenToEthInputPrice.call(1)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'invalid JUMP')
      }
    })

    describe('e2e testing of a trade', () => {
      it('adds liquidity and trades ETH for TBTC', async () => {
        // This avoids rabbit-hole debugging
        // stemming from the fact Vyper is new and they don't do REVERT's
        // so any failed assert's will throw an invalid JMP or something cryptic
        expect(
          await web3.eth.getBalance(accounts[0])
        ).to.not.eq('0')

        expect(
          await web3.eth.getBalance(accounts[1])
        ).to.not.eq('0')

        // Both tokens use 18 decimal places, so we can use toWei here.
        const TBTC_AMT = web3.utils.toWei('50', 'ether')
        const ETH_AMT = web3.utils.toWei('1', 'ether')

        // Mint TBTC
        await tbtc.mint(
          accounts[0],
          TBTC_AMT
        )
        await tbtc.mint(
          accounts[1],
          TBTC_AMT
        )

        await tbtc.approve(tbtcExchange.address, TBTC_AMT, { from: accounts[0] })
        await tbtc.approve(tbtcExchange.address, TBTC_AMT, { from: accounts[1] })

        // min_liquidity, max_tokens, deadline
        const TBTC_ADDED = web3.utils.toWei('10', 'ether')
        await tbtcExchange.addLiquidity(
          '0',
          TBTC_ADDED,
          UniswapHelpers.getDeadline(),
          { value: ETH_AMT }
        )

        // it will be at an exchange rate of
        // 10 TBTC : 1 ETH
        const TBTC_BUY_AMT = web3.utils.toWei('1', 'ether')

        // rough price - we don't think about slippage
        // we are testing that Uniswap works, not testing the exact
        // formulae of the price invariant
        // when they come out with uniswap.js, this code could be made better
        const priceEth = await tbtcExchange.getTokenToEthInputPrice.call(TBTC_BUY_AMT)
        expect(priceEth.toString()).to.eq('90661089388014913')

        const buyer = accounts[1]

        // def ethToTokenSwapInput(min_tokens: uint256, deadline: timestamp) -> uint256:
        await tbtcExchange.ethToTokenSwapInput(
          TBTC_BUY_AMT,
          UniswapHelpers.getDeadline(),
          { value: UniswapHelpers.calcWithFee(priceEth), from: buyer }
        )

        const balance = await tbtc.balanceOf(buyer)
        expect(balance.gt(TBTC_BUY_AMT))
      })
    })
  })
})
