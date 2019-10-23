const TBTCToken = artifacts.require('TBTCToken')
const TestToken = artifacts.require('TestToken')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')
const TBTCSystem = artifacts.require('TBTCSystem')

const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')

const ECDSAKeepStub = artifacts.require('ECDSAKeepStub')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')
const TestDepositUtils = artifacts.require('TestDepositUtils')

const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name
  { name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'TestDeposit', contract: TestDeposit },
  { name: 'TestDepositUtils', contract: TestDepositUtils },
  { name: 'ECDSAKeepStub', contract: ECDSAKeepStub }]

import { UniswapFactoryAddress } from '../../migrations/externals'
import expectThrow from '../helpers/expectThrow'
import utils from '../utils'
import { integration } from './helpers/integration'
import { UniswapHelpers } from './helpers/uniswap'
import { AssertBalance } from '../helpers/assertBalance'

const BN = require('bn.js')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

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
  describe('deployment', async () => {
    it('created the TBTC UniswapExchange', async () => {
      const tbtcToken = await TBTCToken.deployed()
      expect(tbtcToken).to.not.be.empty

      const tbtcSystem = await TBTCSystem.deployed()
      const tbtcExchangeAddress = await tbtcSystem.getTBTCUniswapExchange()
      expect(tbtcExchangeAddress).to.not.be.empty

      const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)
      expect(await uniswapFactory.getExchange(tbtcToken.address)).to.equal(tbtcExchangeAddress)
    })
  })

  describe('end-to-end trade with TBTC', () => {
    let tbtcToken
    let tbtcExchange

    before(async () => {
      const tbtcSystem = await TBTCSystem.deployed()
      tbtcToken = await TestToken.new(tbtcSystem.address)

      // create a uniswap exchange for our TestToken
      const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)
      await uniswapFactory.createExchange(tbtcToken.address)
      const tbtcExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
      tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddress)
    })

    it('has no liquidity by default', async () => {
      await expectThrow(
        tbtcExchange.getEthToTokenOutputPrice(1)
      )
    })

    it('adds liquidity and trades ETH for TBTC', async () => {
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

      let deadline = await UniswapHelpers.getDeadline(web3)
      /* eslint-disable no-multi-spaces */
      await tbtcExchange.addLiquidity(
        '0',                           // min_liquidity
        TBTC_SUPPLY_AMOUNT,            // max_tokens
        deadline,                      // deadline
        { value: ETH_SUPPLY_AMOUNT, from: seller }
      )
      /* eslint-enable no-multi-spaces */

      // Rough price - we don't think about slippage
      // We are testing that Uniswap works, not testing the exact formulae of the price invariant.
      // When they come out with uniswap.js, this code could be made better
      /* eslint-disable no-irregular-whitespace */
      const priceEth = await tbtcExchange.getEthToTokenOutputPrice.call(TBTC_BUY_AMOUNT)

      // Uniswap uses constant product formula to price trades
      // Where (x,y) are the reserves of two tokens respectively
      //     k = xy  (1)
      // (1) is the price invariant. When liquidity is first added, k is set.
      // We added x=50 (TBTC) and y=1 (ETH) of liquidity.
      // k = xy = 50*1 = 50
      // Every trade includes a 0.3% Uniswap fee
      // x = 50
      // y = 1
      // The buy amount is 1 TBTC, we calculate the change in supply (cx)
      // b = 1
      // cx = (50 - b) = 49
      // We have to include the 0.3% fee into the calculation:
      // cx = (50 - b)*0.997
      //    = 48.853
      // Then we calculate the shift in the reserve of y, to maintain the invariant:
      // cy = y / cx
      //    = 1 / 48.853
      //    = 0.02047
      // And that's the price we must pay in ETH (y)
      /* eslint-enable no-irregular-whitespace */
      expect(priceEth.toString()).to.equal(web3.utils.toWei('0.020469571981249873'))

      deadline = await UniswapHelpers.getDeadline(web3)
      /* eslint-disable no-multi-spaces */
      await tbtcExchange.ethToTokenSwapOutput(
        TBTC_BUY_AMOUNT,                 // min_tokens
        deadline,                        // deadline
        { value: priceEth, from: buyer }
      )
      /* eslint-enable no-multi-spaces */

      const balance = await tbtcToken.balanceOf(buyer)
      expect(balance).to.eq.BN(TBTC_BUY_AMOUNT)
    })
  })

  describe('DepositLiquidation', async () => {
    let deployed
    let tbtcSystem

    before(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
    })

    describe('#attemptToLiquidateOnchain', async () => {
      let assertBalance
      let deposit
      let uniswapExchange
      let tbtcToken

      beforeEach(async () => {
        deposit = deployed.TestDeposit
        tbtcSystem = await TBTCSystemStub.new(utils.address0)
        tbtcToken = await TestToken.new(tbtcSystem.address)

        // create a uniswap exchange for our TestToken
        const uniswapFactory = await IUniswapFactory.at(UniswapFactoryAddress)
        await uniswapFactory.createExchange(tbtcToken.address)
        const uniswapExchangeAddress = await uniswapFactory.getExchange(tbtcToken.address)
        uniswapExchange = await IUniswapExchange.at(uniswapExchangeAddress)

        await tbtcSystem.reinitialize(uniswapExchangeAddress)

        deposit.setExteriorAddresses(tbtcSystem.address, tbtcToken.address)
        tbtcSystem.forceMint(accounts[0], web3.utils.toBN(deposit.address))

        // Helpers
        assertBalance = new AssertBalance(tbtcToken)
      })

      it('liquidates using Uniswap successfully', async () => {
        const ethAmount = web3.utils.toWei('0.2', 'ether')
        const tbtcAmount = '100000000'
        await UniswapHelpers.addLiquidity(accounts[0], uniswapExchange, tbtcToken, ethAmount, tbtcAmount)

        const minTbtcAmount = '100100000'
        // hard-coded from previous run
        const expectedPrice = new BN('223138580092984811')

        await assertBalance.eth(deposit.address, '0')
        await assertBalance.tbtc(deposit.address, '0')
        await deposit.send(expectedPrice, { from: accounts[0] })
        await assertBalance.eth(deposit.address, expectedPrice.toString())

        const price = await uniswapExchange.getEthToTokenOutputPrice.call(minTbtcAmount)
        expect(price).to.eq.BN(expectedPrice)

        const retval = await deposit.attemptToLiquidateOnchain.call()
        expect(retval).to.be.true
        await deposit.attemptToLiquidateOnchain()

        await assertBalance.tbtc(deposit.address, minTbtcAmount)
        await assertBalance.eth(deposit.address, '0')
      })

      it('returns false if cannot buy up enough TBTC', async () => {
        const ethAmount = web3.utils.toWei('0.2', 'ether')
        const tbtcAmount = '100000000'
        await UniswapHelpers.addLiquidity(accounts[0], uniswapExchange, tbtcToken, ethAmount, tbtcAmount)

        // hard-coded from previous run
        const expectedPrice = new BN('223138580092984811')
        const depositEthFunding = expectedPrice.sub(new BN(100))

        await assertBalance.eth(deposit.address, '0')
        await assertBalance.tbtc(deposit.address, '0')
        await deposit.send(depositEthFunding, { from: accounts[0] })
        await assertBalance.eth(deposit.address, depositEthFunding.toString())

        const retval = await deposit.attemptToLiquidateOnchain.call()
        expect(retval).to.be.false
      })
    })
  })
})
