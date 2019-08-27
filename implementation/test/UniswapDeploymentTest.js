const TBTCToken = artifacts.require('TBTCToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

const UniswapDeployment = artifacts.require('UniswapDeployment')

import { UniswapHelpers } from './helpers/uniswap'

// Tests the Uniswap deployment

contract('Uniswap', (accounts) => {
  let tbtcToken

  describe('deployment', async () => {
    it('deployed the uniswap factory and exchange', async () => {
      tbtcToken = await TBTCToken.deployed()
      expect(tbtcToken).to.not.be.empty

      const uniswapDeployment = await UniswapDeployment.deployed()
      const uniswapFactoryAddr = await uniswapDeployment.factory()
      expect(uniswapFactoryAddr).to.not.be.empty

      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtcToken.address)
      expect(tbtcExchangeAddr).to.not.be.empty
    })
  })

  describe('TBTC Uniswap Exchange', () => {
    let tbtc
    let tbtcExchange

    beforeEach(async () => {
      /* eslint-disable no-unused-vars */
      tbtc = await TBTCToken.new()

      // We rely on the already pre-deployed Uniswap factory here.
      const uniswapDeployment = await UniswapDeployment.deployed()
      const uniswapFactoryAddr = await uniswapDeployment.factory()

      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)

      const res = await uniswapFactory.createExchange(tbtc.address)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtc.address)


      tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddr)
      /* eslint-enable no-unused-vars */
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
