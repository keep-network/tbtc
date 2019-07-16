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
const TestDepositUtils = artifacts.require('TestDepositUtils')

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
  { name: 'TestDepositUtils', contract: TestDepositUtils },
  { name: 'KeepStub', contract: KeepStub },
  { name: 'TBTCStub', contract: TBTCStub },
  { name: 'TBTCSystemStub', contract: TBTCSystemStub }]


contract('DepositLiquidation', (accounts) => {
  let deployed
  let testInstance


  before(async () => {
    deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
    testInstance = deployed.TestDeposit
    testInstance.setExteroriorAddresses(
      deployed.TBTCSystemStub.address, deployed.TBTCStub.address, deployed.KeepStub.address
    )
  })

  beforeEach(async () => {
    await testInstance.reset()
  })

  describe.only('liquidation flows', async () => {
    let uniswapExchange
    let tbtc
    let tbtcSystem
    let deposit
    let keep

    let assertBalance

    let lotSize
    let beneficiaryReward
    // let signerFee

    const keepBondAmount = 1000000

    beforeEach(async () => {
      tbtc = await TBTC.new()
      const uniswapFactory = await UniswapFactoryStub.new()
      uniswapExchange = await UniswapExchangeStub.new(tbtc.address)
      deposit = deployed.TestDepositUtils

      lotSize = await deployed.TBTCConstants.getLotSize()
      beneficiaryReward = await deposit.beneficiaryReward()
      // signerFee = await deposit.signerFee()

      tbtcSystem = deployed.TBTCSystemStub
      keep = deployed.KeepStub

      // Set exterior addresses
      await tbtcSystem.setExteroriorAddresses(
        uniswapFactory.address,
        tbtc.address
      )
      await uniswapFactory.setTbtcExchange(uniswapExchange.address)
      await deposit.setExteroriorAddresses(
        tbtcSystem.address,
        tbtc.address,
        keep.address,
      )

      // Setup mock Uniswap with liquidity, set prices
      const tbtcSupply = new BN(lotSize).mul(new BN(5))
      const tbtcPrice = new BN(lotSize).add(new BN(beneficiaryReward))
      await tbtc.mint(accounts[0], tbtcSupply, { from: accounts[0] })
      await tbtc.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })
      await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: keepBondAmount })
      await uniswapExchange.setPrices(keepBondAmount, tbtcPrice)

      assertBalance = new AssertBalanceHelpers(tbtc)
    })


    describe('#attemptToLiquidateOnchain', async () => {
      it('liquidates for 1.005 tbtc', async () => {
        await deposit.send(keepBondAmount, { from: accounts[0] })
        await assertBalance.tbtc(deposit.address, '0')
        await assertBalance.eth(deposit.address, ''+keepBondAmount)

        await deposit.attemptToLiquidateOnchain()
      })

      it('returns false if buy under MIN_TBTC threshold', async () => {})
    })

    describe('redemption', async () => {
      it('#startSignerFraudLiquidation', async () => {
        // eth bonds should be liquidated for tbtc
        // via uniswap
        //
        // if this is coming from redemption, then tbtc is refunded to redeemer
        // if this is not, tbtc is burnt to maintain supply peg
        //
        // in either case, the beneficiary reward (0.005 tbtc) should be sent back to the beneficiary
        // and the requestor should get whatever remaining eth
        //
        // we also need to check for the case in which onchain liquidation fails

        const beneficiary = accounts[1]
        const requestor = accounts[2]

        const requestorBalance1 = new BN(await web3.eth.getBalance(requestor))

        // set the owner/beneficiary of the deposit
        // TODO(liamz): set beneficiary more realistically
        await tbtcSystem.setOwner(beneficiary)

        // give the keep/deposit eth
        // TODO(liamz): use createNewDeposit when it works
        await assertBalance.eth(keep.address, '0')
        await keep.send(keepBondAmount, { from: accounts[0] })
        await assertBalance.eth(keep.address, '1000000') // TODO(liamz): replace w/ keep.checkBondAmount()
        await assertBalance.tbtc(deposit.address, '0')

        // set requestor
        // TODO(liamz): set requestorAddress more realistically
        const digest = '0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc'
        await deposit.setRequestInfo(
          requestor,
          '0x' + '11'.repeat(20),
          0, 0, digest
        )


        await deposit.startSignerFraudLiquidation()

        // truffleAssert.eventEmitted(tx, 'ECDSAKeepCreated', (ev) => {
        //   instanceAddress = ev.keepAddress;
        //   return true;
        // });

        const requestorBalance2 = new BN(await web3.eth.getBalance(requestor))


        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState).to.eq.BN(utils.states.LIQUIDATED)

        // tbtc distributed to requestor and beneficiary
        await assertBalance.tbtc(requestor, lotSize)
        await assertBalance.tbtc(beneficiary, beneficiaryReward)

        // requestor gets (any) remaining eth
        // TODO(liamz): no remaining eth in this test, but in another forsho
        expect(
          requestorBalance2.sub(requestorBalance1)
        ).to.eq.BN(new BN('0'))

        // deposit should have 0
        await assertBalance.tbtc(deposit.address, new BN('0'))
        await assertBalance.eth(deposit.address, '0')

        // check distributeEthToKeepGroup
        const keepGroupTotalEth = await keep.keepGroupTotalEth()
        expect(keepGroupTotalEth).to.eq.BN(new BN('0'))
      })

      it('#startSignerAbortLiquidation', async () => {
        // eth bonds should be liquidated for tbtc
        // via uniswap
        //
        // if this is coming from redemption, then tbtc is refunded to redeemer
        // if this is not, tbtc lotSize is burnt to maintain supply peg
        //
        // in either case, the beneficiary reward (0.005 tbtc) should be sent back to the beneficiary
        // and the keep group should get the remaining eth via pushFundsToKeepGroup
        //
        // we also need to check for the case in which onchain liquidation fails
      })
    })

    describe('non-redemption', async () => {

    })
  })
})
