const TBTCToken = artifacts.require('TBTCToken')
const TestToken = artifacts.require('TestToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTCSystem = artifacts.require('TBTCSystem')

import utils from '../utils'
import { createSnapshot, restoreSnapshot } from '../helpers/snapshot'
import { UniswapHelpers } from '../helpers/uniswap'
import { e2e } from './helper'

const BN = require('bn.js')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCSystemStub', contract: TBTCSystemStub },
]

e2e('Uniswap', (accounts) => {
  let deployed
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
    it('created the tBTC UniswapExchange', async () => {
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

  describe('end-to-end trade with tBTC', () => {
    let tbtcExchange

    before(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
      tbtcToken = await TestToken.new(deployed.TBTCSystemStub.address)

      const tbtcSystem = await TBTCSystem.deployed()
      const uniswapFactoryAddr = await tbtcSystem.getUniswapFactory()
      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)

      await uniswapFactory.createExchange(tbtcToken.address)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtcToken.address)
      tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddr)
    })

    it('has no liquidity by default', async () => {
      try {
        await tbtcExchange.getEthToTokenOutputPrice(1)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'invalid JUMP')
      }
    })

    it('adds liquidity and trades ETH for tBTC', async () => {
      // This avoids rabbit-hole debugging
      // stemming from the fact Vyper is new and they don't do REVERT's
      // so any failed assertions in Vyper, will throw cryptic errors like "invalid JMP"
      // for something simple like `assert balance(msg.sender) > 0`
      expect(
        await web3.eth.getBalance(accounts[0])
      ).to.not.equal('0')

      expect(
        await web3.eth.getBalance(accounts[1])
      ).to.not.equal('0')

      // Both tokens use 18 decimal places, so we can use toWei here.
      const TBTC_AMT = web3.utils.toWei('50', 'ether')
      const ETH_AMT = web3.utils.toWei('1', 'ether')

      // Mint TBTC
      await tbtcToken.forceMint(
        accounts[0],
        TBTC_AMT
      )

      await tbtcToken.approve(tbtcExchange.address, TBTC_AMT, { from: accounts[0] })
      const TBTC_ADDED = web3.utils.toWei('10', 'ether')

      // min_liquidity, max_tokens, deadline
      await tbtcExchange.addLiquidity(
        '0',
        TBTC_ADDED,
        UniswapHelpers.getDeadline(),
        { value: ETH_AMT }
      )

      // it will be at an exchange rate of
      // 10 TBTC : 1 ETH
      const TBTC_BUY_AMT = web3.utils.toWei('1', 'ether')

      // Rough price - we don't think about slippage
      // We are testing that Uniswap works, not testing the exact formulae of the price invariant.
      // When they come out with uniswap.js, this code could be made better
      const priceEth = await tbtcExchange.getEthToTokenOutputPrice.call(TBTC_BUY_AMT)
      expect(priceEth.toString()).to.equal('111445447453471526')

      // def ethToTokenSwapInput(min_tokens: uint256, deadline: timestamp) -> uint256:
      const buyer = accounts[1]
      await tbtcExchange.ethToTokenSwapOutput(
        TBTC_BUY_AMT,
        UniswapHelpers.getDeadline(),
        { value: priceEth, from: buyer }
      )

      const balance = await tbtcToken.balanceOf(buyer)
      expect(balance).to.eq.BN(TBTC_BUY_AMT)
    })
  })
})
