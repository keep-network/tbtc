const TBTCToken = artifacts.require('TBTCToken')
const TestToken = artifacts.require('TestToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTCSystem = artifacts.require('TBTCSystem')

import utils from '../utils'
import { createSnapshot, restoreSnapshot } from '../helpers/snapshot'
import { UniswapHelpers } from './helpers/uniswap'
import { integration } from './helpers/integration'

const BN = require('bn.js')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCSystemStub', contract: TBTCSystemStub },
]

/**
 * Tests integration with the external Uniswap deployment.
 *
 * A briefer on Uniswap - it's programmed in Vyper. Vyper is an experimental alternative
 * language to Solidity.
 *
 * One problem with using anything built in Vyper, is that it doesn't error in
 * the same way a Solidity-compiled contract will. Where Solidity has `require` statements,
 * Vyper has `assert` statements. A failed `require` will `REVERT` the EVM, but a failed
 * `assert` in Vyper will throw cryptic errors like `invalid JUMP`
 *
 * In the test below, there are manual sanity checks for Vyper `assert`'s to avoid wasting programmer time.
 * One example is `assert balance(msg.sender) > 0`, which can easily fail due to the Ganache default account
 * running out of its 100ETH.
 */

integration('Uniswap', (accounts) => {
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
        // See file header for an explanation
        assert.include(e.message, 'invalid JUMP')
      }
    })

    it('adds liquidity and trades ETH for tBTC', async () => {
      const seller = accounts[0]
      const buyer = accounts[1]

      // Manual sanity checks for Vyper assertions. (see file header)
      expect(await web3.eth.getBalance(seller)).to.not.equal('0')
      expect(await web3.eth.getBalance(buyer)).to.not.equal('0')

      // Both tokens use 18 decimal places, so we can use toWei here.
      const TBTC_SUPPLY_AMOUNT = web3.utils.toWei('50', 'ether')
      const ETH_SUPPLY_AMOUNT = web3.utils.toWei('1', 'ether')
      const TBTC_BUY_AMOUNT = web3.utils.toWei('1', 'ether')

      // Mint TBTC
      await tbtcToken.forceMint(
        seller,
        TBTC_SUPPLY_AMOUNT
      )

      await tbtcToken.approve(tbtcExchange.address, TBTC_SUPPLY_AMOUNT, { from: seller })

      /* eslint-disable no-multi-spaces */
      await tbtcExchange.addLiquidity(
        '0',                           // min_liquidity
        TBTC_SUPPLY_AMOUNT,            // max_tokens
        UniswapHelpers.getDeadline(),  // deadline
        { value: ETH_SUPPLY_AMOUNT }
      )
      /* eslint-enable no-multi-spaces */

      // Rough price - we don't think about slippage
      // We are testing that Uniswap works, not testing the exact formulae of the price invariant.
      // When they come out with uniswap.js, this code could be made better
      const priceEth = await tbtcExchange.getEthToTokenOutputPrice.call(TBTC_BUY_AMOUNT)
      expect(priceEth.toString()).to.equal('20469571981249873')

      /* eslint-disable no-multi-spaces */
      await tbtcExchange.ethToTokenSwapOutput(
        TBTC_BUY_AMOUNT,                 // min_tokens
        UniswapHelpers.getDeadline(),    // deadline
        { value: priceEth, from: buyer }
      )
      /* eslint-enable no-multi-spaces */

      const balance = await tbtcToken.balanceOf(buyer)
      expect(balance).to.eq.BN(TBTC_BUY_AMOUNT)
    })
  })
})
