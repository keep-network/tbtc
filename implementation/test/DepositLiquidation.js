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
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')
const UniswapDeployment = artifacts.require('UniswapDeployment')
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
import { getTxCost } from './helpers/tx_cost'
import { UniswapHelpers } from './helpers/uniswap'

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

  async function realUniswapExchange(tbtcSystem, tbtc) {
    const uniswapDeployment = await UniswapDeployment.deployed()
    const uniswapFactoryAddr = await uniswapDeployment.factory()
    const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)

    // Create exchange
    await uniswapFactory.createExchange(tbtc.address)
    const tbtcExchangeAddr = await uniswapFactory.getExchange.call(tbtc.address)
    const uniswapExchange = await IUniswapExchange.at(tbtcExchangeAddr)

    // Inject
    await tbtcSystem.setExteroriorAddresses(
      uniswapFactory.address,
      tbtc.address
    )

    return uniswapExchange
  }

  async function realUniswapLiquidity(uniswapExchange, tbtcSupply, ethSupply) {
    await uniswapExchange.addLiquidity(
      '0',
      tbtcSupply,
      UniswapHelpers.getDeadline(),
      { from: accounts[0], value: ethSupply }
    )
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

    describe('using UniswapExchangeStub', async () => {
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

        // Uniswap requires a minimum of 1000000000 wei for the initial addLiquidity call
        // So we add an extreme amount of ETH here.
        // We will use different values for keepBondAmount in future,
        // so multiplying by 1000000000x, while excessive, is robust.
        const MINIMUM_UNISWAP_ETH_LIQUIDITY = '1000000000'
        const ethSupply = new BN(keepBondAmount).mul(new BN(MINIMUM_UNISWAP_ETH_LIQUIDITY))

        await uniswapExchange.mockLiquidity(tbtcSupply, { from: accounts[0], value: ethSupply })
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

    describe.only('using real Uniswap deployment', async () => {
      beforeEach(async () => {
        deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
        tbtc = await TBTC.new()
        deposit = await TestDeposit.new()

        const tbtcSystem = deployed.TBTCSystemStub
        const keep = deployed.KeepStub

        uniswapExchange = await realUniswapExchange(tbtcSystem, tbtc)

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

        // Uniswap requires a minimum of 1000000000 wei for the initial addLiquidity call
        // So we add an extreme amount of ETH here.
        // We will use different values for keepBondAmount in future,
        // so multiplying by 1000000000x, while excessive, is robust.
        const MINIMUM_UNISWAP_ETH_LIQUIDITY = '1000000000'
        const ethSupply = new BN(keepBondAmount).mul(new BN(MINIMUM_UNISWAP_ETH_LIQUIDITY))

        await realUniswapLiquidity(uniswapExchange, tbtcSupply, ethSupply)
      })

      it('returns false when exchange = 0x0 (todo)', async () => {
      })

      it('returns false if buy under MIN_TBTC threshold (todo)', async () => {
      })

      it.only('liquidates eth for 1.005 tbtc (MIN_TBTC)', async () => {
        const expectedPriceWei = new BN('166388379320629')
        // 0.3%
        const priceWithFee = expectedPriceWei.add(expectedPriceWei.mul(new BN(1003)).div(new BN(1000)))

        await deposit.send(priceWithFee, { from: accounts[0] })
        await assertBalance.tbtc(deposit.address, '0')
        await assertBalance.eth(deposit.address, priceWithFee.toString())

        const priceEth = await uniswapExchange.getTokenToEthInputPrice.call('100100000')
        expect(priceEth).to.eq.BN(expectedPriceWei)

        const retval = await deposit.attemptToLiquidateOnchain.call()
        expect(retval).to.be.true
        await deposit.attemptToLiquidateOnchain()

        await assertBalance.tbtc(deposit.address, '100100000')
        await assertBalance.eth(deposit.address, '82210148308270')
      })

      it('only liquidates enough eth to buy MIN_TBTC', async () => {
        // todo assert prices

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
  })

  describe('liquidation flows', async () => {
    let deployed
    let deposit
    let tbtc
    let uniswapExchange
    let keep

    let assertBalance

    const beneficiary = accounts[1]
    const maintainer = accounts[5]
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
      const maintainerBalance = new BN(await web3.eth.getBalance(maintainer))

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


      const tx = await deposit.startSignerFraudLiquidation({ from: maintainer })
      const txCostEth = await getTxCost(tx)

      const depositState = await deposit.getState.call()
      expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

      // check remaining eth sent to maintainer
      await assertBalance.eth(deposit.address, '0')
      const maintainerBalanceExpected = maintainerBalance.sub(txCostEth).add(remainingEthAfterLiquidation)
      const maintainerBalance2 = await web3.eth.getBalance(maintainer)
      expect(maintainerBalance2).to.eq.BN(maintainerBalanceExpected)
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

      // check remaining eth sent to keep group
      await assertBalance.eth(deposit.address, '0')
      const keepGroupTotalEth = await keep.keepGroupTotalEth.call()
      expect(keepGroupTotalEth).to.eq.BN(remainingEthAfterLiquidation)
    })
  })
})
