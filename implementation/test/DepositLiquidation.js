const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')

const KeepStub = artifacts.require('KeepStub')
const TBTCStub = artifacts.require('TBTCStub')
const TBTC = artifacts.require('TBTC')
const UniswapFactoryStub = artifacts.require('UniswapFactoryStub')
const UniswapExchangeStub = artifacts.require('UniswapExchangeStub')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')
const TestDepositLiquidation = artifacts.require('TestDepositLiquidation')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

import { AssertBalanceHelpers } from './helpers/assert_balance'

const TEST_DEPOSIT_DEPLOY = [
  { name: 'BytesLib', contract: BytesLib },
  { name: 'BTCUtils', contract: BTCUtils },
  { name: 'ValidateSPV', contract: ValidateSPV },
  { name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name
  { name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'TestDeposit', contract: TestDeposit },
  { name: 'TestDepositLiquidation', contract: TestDepositLiquidation },
  { name: 'KeepStub', contract: KeepStub },
  { name: 'TBTCStub', contract: TBTCStub },
  { name: 'TBTCSystemStub', contract: TBTCSystemStub },
  { name: 'UniswapFactoryStub', contract: UniswapFactoryStub }]


const lotSize = '100000000'
const beneficiaryReward = '100000'
const keepBondAmount = 1000000

contract('DepositLiquidation', (accounts) => {
  describe('#attemptToLiquidateOnchain', () => {
    let deployed
    let deposit
    let tbtc
    let uniswapExchange

    let assertBalance

    beforeEach(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
      tbtc = await TBTC.new()
      deposit = await TestDeposit.new()
      uniswapExchange = await UniswapExchangeStub.new(tbtc.address)

      const uniswapFactory = deployed.UniswapFactoryStub
      const tbtcSystem = deployed.TBTCSystemStub
      const keep = deployed.KeepStub

      await deposit.setExteroriorAddresses(
        tbtcSystem.address,
        tbtc.address,
        keep.address,
      )

      // Inject mock exchange
      await tbtcSystem.setExteroriorAddresses(
        uniswapFactory.address,
        tbtc.address
      )
      await uniswapFactory.setTbtcExchange(uniswapExchange.address)

      // Test helpers
      assertBalance = new AssertBalanceHelpers(tbtc)
    })

    it('liquidates eth for 1.005 tbtc (MIN_TBTC)', async () => {
      const tbtcSupply = new BN(lotSize).mul(new BN(5))
      const tbtcPrice = new BN(lotSize).add(new BN(beneficiaryReward))
      await tbtc.mint(accounts[0], tbtcSupply, { from: accounts[0] })
      await tbtc.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })
      await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: keepBondAmount })
      await uniswapExchange.setPrices(keepBondAmount, tbtcPrice)

      await deposit.send(keepBondAmount, { from: accounts[0] })
      await assertBalance.tbtc(deposit.address, '0')
      await assertBalance.eth(deposit.address, ''+keepBondAmount)

      const retval = await deposit.attemptToLiquidateOnchain.call()
      expect(retval).to.be.true
      await deposit.attemptToLiquidateOnchain()

      await assertBalance.tbtc(deposit.address, '100100000')
      await assertBalance.eth(deposit.address, '0')
    })

    // it('returns false if buy under MIN_TBTC threshold', async () => {})
  })

  describe.only('liquidation flows', async () => {
    let deployed
    let deposit
    let tbtc
    let uniswapExchange
    let keep

    let assertBalance

    const beneficiary = accounts[1]
    const DIGEST = '0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc'

    beforeEach(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
      tbtc = await TBTC.new()
      uniswapExchange = await UniswapExchangeStub.new(tbtc.address)

      deposit = deployed.TestDepositLiquidation
      keep = deployed.KeepStub

      const uniswapFactory = deployed.UniswapFactoryStub
      const tbtcSystem = deployed.TBTCSystemStub

      await tbtcSystem.setExteroriorAddresses(
        uniswapFactory.address,
        tbtc.address
      )
      await deposit.setExteroriorAddresses(
        tbtcSystem.address,
        tbtc.address,
        keep.address,
      )
      await uniswapFactory.setTbtcExchange(uniswapExchange.address)

      assertBalance = new AssertBalanceHelpers(tbtc)


      // ----------

      // setup stub uniswap exchange
      const ethSupply = web3.utils.toWei('3', 'ether')
      const tbtcSupply = new BN(lotSize).mul(new BN(5))
      await tbtc.mint(accounts[0], tbtcSupply, { from: accounts[0] })
      await tbtc.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })
      await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: ethSupply })

      // set the owner/beneficiary of the deposit
      // TODO(liamz): set beneficiary more realistically
      await tbtcSystem.setOwner(beneficiary)

      // give the keep/deposit eth
      // TODO(liamz): use createNewDeposit when it works
      await assertBalance.eth(keep.address, '0')
      await keep.send(keepBondAmount, { from: accounts[0] })
      await assertBalance.eth(keep.address, '1000000') // TODO(liamz): replace w/ keep.checkBondAmount()
      await assertBalance.tbtc(deposit.address, '0')
    })


    describe('redemption', async () => {
      const requestor = accounts[2]
      let requestorBalance1

      beforeEach(async () => {
        requestorBalance1 = new BN(await web3.eth.getBalance(requestor))

        // set requestor
        // TODO(liamz): set requestorAddress more realistically
        await deposit.setRequestInfo(
          requestor,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )
      })

      it('#startSignerFraudLiquidation', async () => {
        const tbtcPrice = new BN(lotSize).add(new BN(beneficiaryReward))
        await uniswapExchange.setPrices(keepBondAmount, tbtcPrice)

        await deposit.startSignerFraudLiquidation()

        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState).to.eq.BN(utils.states.LIQUIDATED)

        // deposit should have 0
        await assertBalance.tbtc(deposit.address, new BN('0'))
        await assertBalance.eth(deposit.address, '0')

        // tbtc distributed to requestor and beneficiary
        await assertBalance.tbtc(requestor, lotSize)
        await assertBalance.tbtc(beneficiary, beneficiaryReward)

        // requestor gets (any) remaining eth
        // TODO(liamz): no remaining eth in this test, but in another forsho
        const requestorBalance2 = new BN(await web3.eth.getBalance(requestor))
        expect(
          requestorBalance2.sub(requestorBalance1)
        ).to.eq.BN(new BN('0'))

        // check distributeEthToKeepGroup
        const keepGroupTotalEth = await keep.keepGroupTotalEth()
        expect(keepGroupTotalEth).to.eq.BN(new BN('0'))
      })

      it('#startSignerAbortLiquidation', async () => {

      })
    })

    describe('non-redemption', async () => {
      const NON_REDEMPTION_ADDRESS = '0x' + '00'.repeat(20)

      beforeEach(async () => {
        // set requestor
        // TODO(liamz): set requestorAddress more realistically
        await deposit.setRequestInfo(
          NON_REDEMPTION_ADDRESS,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )
      })

      it('#startSignerFraudLiquidation', async () => {})
      it('#startSignerAbortLiquidation', async () => {})
    })
  })
})
