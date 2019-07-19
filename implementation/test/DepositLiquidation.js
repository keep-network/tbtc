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
  { name: 'TBTCSystemStub', contract: TBTCSystemStub }]


const lotSize = '100000000'
const beneficiaryReward = '100000'
const tbtcPrice = new BN(lotSize).add(new BN(beneficiaryReward))
const keepBondAmount = 1000000

// used in a few tests
const remainingEthAfterLiquidation = new BN('1000')

contract('DepositLiquidation', (accounts) => {
  async function mockUniswapExchange(tbtcSystem, tbtc) {
    if (!tbtcSystem) throw new Error
    if (!tbtc) throw new Error

    const uniswapFactory = await UniswapFactoryStub.new()
    const uniswapExchange = await UniswapExchangeStub.new(tbtc.address)

    await uniswapFactory.setTbtcExchange(uniswapExchange.address)

    // Inject mock exchange
    await tbtcSystem.setExteroriorAddresses(
      uniswapFactory.address,
      tbtc.address
    )

    return uniswapExchange
  }

  async function mockUniswapLiquidity(tbtc, uniswapExchange) {
    const ethSupply = web3.utils.toWei('3', 'ether')
    const tbtcSupply = new BN(lotSize).mul(new BN(5))
    await tbtc.mint(accounts[0], tbtcSupply, { from: accounts[0] })
    await tbtc.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })
    await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: ethSupply })
  }

  describe.only('#attemptToLiquidateOnchain', () => {
    let deployed
    let deposit
    let tbtc
    let uniswapExchange

    let assertBalance

    beforeEach(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
      tbtc = await TBTC.new()
      deposit = await TestDeposit.new()

      const tbtcSystem = deployed.TBTCSystemStub
      const keep = deployed.KeepStub

      uniswapExchange = await mockUniswapExchange(tbtcSystem, tbtc)

      await deposit.setExteroriorAddresses(
        tbtcSystem.address,
        tbtc.address,
        keep.address,
      )

      // Test helpers
      assertBalance = new AssertBalanceHelpers(tbtc)


      // Mint tBTC, mock liquidity
      const tbtcSupply = new BN(lotSize).mul(new BN(5))
      await tbtc.mint(accounts[0], tbtcSupply, { from: accounts[0] })
      await tbtc.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })
      await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: keepBondAmount })
    })

    it('returns false when exchange = 0x0', async () => {

    })

    it('returns false if buy under MIN_TBTC threshold', async () => {

    })

    it('liquidates eth for 1.005 tbtc (MIN_TBTC)', async () => {
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

    it('only liquidates enough eth to buy MIN_TBTC', async () => {
      await uniswapExchange.setPrices(
        new BN(keepBondAmount).sub(remainingEthAfterLiquidation),
        tbtcPrice
      )

      await deposit.send(keepBondAmount, { from: accounts[0] })
      await assertBalance.tbtc(deposit.address, '0')
      await assertBalance.eth(deposit.address, ''+keepBondAmount)

      const retval = await deposit.attemptToLiquidateOnchain.call()
      expect(retval).to.be.true
      await deposit.attemptToLiquidateOnchain()

      await assertBalance.tbtc(deposit.address, '100100000')
      await assertBalance.eth(deposit.address, remainingEthAfterLiquidation.toString())
    })
  })

  describe('liquidation flows', async () => {
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

      const tbtcSystem = deployed.TBTCSystemStub
      deposit = deployed.TestDepositLiquidation
      keep = deployed.KeepStub

      uniswapExchange = await mockUniswapExchange(tbtcSystem, tbtc)

      await deposit.setExteroriorAddresses(
        tbtcSystem.address,
        tbtc.address,
        keep.address,
      )

      assertBalance = new AssertBalanceHelpers(tbtc)

      // ----------
      await mockUniswapLiquidity(tbtc, uniswapExchange)

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


    describe('#startLiquidation', async () => {
      async function assertCommon() {
        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

        // deposit should have 0
        await assertBalance.tbtc(deposit.address, new BN('0'))
        await assertBalance.eth(deposit.address, '0')

        // tbtc distributed to requestor and beneficiary
        await assertBalance.tbtc(beneficiary, beneficiaryReward)

        // check distributeEthToKeepGroup
        const keepGroupTotalEth = await keep.keepGroupTotalEth()
        expect(keepGroupTotalEth).to.eq.BN(new BN('0'))
      }

      it('passes (redemption)', async () => {
        await uniswapExchange.setPrices(keepBondAmount, tbtcPrice)

        const requestor = accounts[2]
        // const requestorBalance1 = new BN(await web3.eth.getBalance(requestor))

        // set requestor
        // TODO(liamz): set requestorAddress more realistically
        await deposit.setRequestInfo(
          requestor,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )

        await deposit.startLiquidation()
        await assertCommon()

        // expect tbtc to be refunded to requestor
        await assertBalance.tbtc(requestor, lotSize)

        // requestor gets (any) remaining eth
        // TODO(liamz): no remaining eth in this test, but in another forsho
        // const requestorBalance2 = new BN(await web3.eth.getBalance(requestor))
        // expect(
        //   requestorBalance2.sub(requestorBalance1)
        // ).to.eq.BN(new BN('0'))
      })


      it('passes (non-redemption)', async () => {
        await uniswapExchange.setPrices(keepBondAmount, tbtcPrice)

        const NON_REDEMPTION_ADDRESS = '0x' + '00'.repeat(20)
        // set requestor
        // TODO(liamz): set requestorAddress more realistically
        await deposit.setRequestInfo(
          NON_REDEMPTION_ADDRESS,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )

        const tx = await deposit.startLiquidation()
        await assertCommon()

        // expect TBTC to be burnt from deposit's account
        const evs = await tbtc.getPastEvents({ fromBlock: tx.receipt.blockNumber })
        const ev = evs[evs.length - 1]
        expect(ev.event).to.equal('Transfer')
        expect(ev.returnValues.from).to.equal(deposit.address)
        expect(ev.returnValues.to).to.equal('0x0000000000000000000000000000000000000000')
        expect(ev.returnValues.value).to.equal(lotSize.toString())
      })
    })

    it('#startSignerFraudLiquidation', async () => {
      await uniswapExchange.setPrices(
        new BN(keepBondAmount).sub(remainingEthAfterLiquidation),
        tbtcPrice
      )

      const requestor = accounts[2]
      await deposit.setRequestInfo(
        requestor,
        '0x' + '11'.repeat(20),
        0, 0, DIGEST
      )

      const maintainer = accounts[5]
      const maintainerBalance1 = new BN(await web3.eth.getBalance(maintainer))

      await deposit.startSignerFraudLiquidation({ from: maintainer })

      const depositState = await deposit.getState.call()
      expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

      const maintainerBalance2 = new BN(await web3.eth.getBalance(maintainer))

      expect(
        maintainerBalance2.sub(maintainerBalance1)
      ).to.eq.BN(remainingEthAfterLiquidation)
    })

    it('#startSignerAbortLiquidation', async () => {
      await uniswapExchange.setPrices(
        new BN(keepBondAmount).sub(remainingEthAfterLiquidation),
        tbtcPrice
      )

      const requestor = accounts[2]
      await deposit.setRequestInfo(
        requestor,
        '0x' + '11'.repeat(20),
        0, 0, DIGEST
      )

      await deposit.startSignerAbortLiquidation({ from: maintainer })

      const keepGroupTotalEth = await keep.keepGroupTotalEth.call()
      expect(keepGroupTotalEth).to.equal(remainingEthAfterLiquidation.toString())
    })
  })
})
