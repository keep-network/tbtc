import expectThrow from './helpers/expectThrow'

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
const TBTCTokenStub = artifacts.require('TBTCTokenStub')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')
const TestDepositUtils = artifacts.require('TestDepositUtils')

const UniswapFactoryStub = artifacts.require('UniswapFactoryStub')
const UniswapExchangeStub = artifacts.require('UniswapExchangeStub')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

import { AssertBalanceHelpers } from './helpers/assertBalance'
import { UniswapHelpers } from './helpers/uniswap'
import { createSnapshot, restoreSnapshot } from './helpers/snapshot'
import { getTxCost } from './helpers/txCost'

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
  { name: 'TBTCTokenStub', contract: TBTCTokenStub },
  { name: 'TBTCSystemStub', contract: TBTCSystemStub },
  { name: 'UniswapFactoryStub', contract: UniswapFactoryStub }]

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

contract('DepositLiquidation', (accounts) => {
  let deployed
  let testInstance
  let beneficiary


  let snapshotId

  before(async () => {
    snapshotId = await createSnapshot()
  })

  beforeEach(async () => {
    snapshotId = await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot(snapshotId)
  })

  after(async () => {
    await restoreSnapshot(snapshotId)
  })

  before(async () => {
    deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
    testInstance = deployed.TestDeposit
    testInstance.setExteroriorAddresses(deployed.TBTCSystemStub.address, deployed.TBTCTokenStub.address, deployed.KeepStub.address)

    const uniswapFactory = deployed.UniswapFactoryStub
    const uniswapExchange = await UniswapExchangeStub.new(deployed.TBTCTokenStub.address)
    await uniswapFactory.setExchange(uniswapExchange.address)
    await deployed.TBTCSystemStub.setExternalAddresses(uniswapFactory.address)

    deployed.TBTCSystemStub.mint(accounts[4], web3.utils.toBN(deployed.TestDeposit.address))
    beneficiary = accounts[4]
  })

  describe('purchaseSignerBondsAtAuction', async () => {
    let requiredBalance

    before(async () => {
      requiredBalance = await deployed.TestDepositUtils.redemptionTBTCAmount.call()
    })

    beforeEach(async () => {
      await testInstance.setState(utils.states.LIQUIDATION_IN_PROGRESS)
      for (let i = 0; i < 4; i++) {
        await deployed.TBTCTokenStub.clearBalance(accounts[i])
      }
    })

    it('sets state to liquidated, logs Liquidated, ', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await deployed.TBTCTokenStub.mint(accounts[0], requiredBalance)

      await testInstance.purchaseSignerBondsAtAuction()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.LIQUIDATED)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('Liquidated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in a liquidation auction', async () => {
      await deployed.TBTCTokenStub.mint(accounts[0], requiredBalance)
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.purchaseSignerBondsAtAuction(),
        'No active auction'
      )
    })

    it('reverts if TBTC balance is insufficient', async () => {
      // mint 1 less than lot size
      const lotSize = await deployed.TBTCConstants.getLotSize.call()
      await deployed.TBTCTokenStub.mint(accounts[0], lotSize - 1)

      await expectThrow(
        testInstance.purchaseSignerBondsAtAuction(),
        'Not enough TBTC to cover outstanding debt'
      )
    })

    it(`burns msg.sender's tokens`, async () => {
      const caller = accounts[2]

      await deployed.TBTCTokenStub.mint(caller, requiredBalance)

      const lotSize = await deployed.TBTCConstants.getLotSize.call()
      const initialTokenBalance = await deployed.TBTCTokenStub.balanceOf(caller)

      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await deployed.TBTCTokenStub.balanceOf(caller)
      const tokenCheck = new BN(finalTokenBalance).add(new BN(lotSize))
      expect(tokenCheck, 'tokens not burned correctly').to.eq.BN(initialTokenBalance)
    })

    it('distributes beneficiary reward', async () => {
      const caller = accounts[2]
      const initialTokenBalance = await deployed.TBTCTokenStub.balanceOf(beneficiary)
      const returned = await deployed.TBTCTokenStub.balanceOf.call(caller)

      await deployed.TBTCTokenStub.mint(caller, requiredBalance)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await deployed.TBTCTokenStub.balanceOf(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add(new BN(returned))

      expect(finalTokenBalance, 'tokens not returned to beneficiary correctly').to.eq.BN(tokenCheck)
    })

    it('distributes value to the caller', async () => {
      const value = 1000000000000
      const caller = accounts[2]
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const initialBalance = await web3.eth.getBalance(caller)

      await testInstance.send(value, { from: accounts[0] })
      await deployed.TBTCTokenStub.mint(caller, requiredBalance)
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(caller)

      expect(new BN(finalBalance), 'caller balance should increase').to.be.gte.BN(initialBalance)
    })

    it('returns keep funds if not fraud', async () => {
      const value = 1000000000000
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const caller = accounts[2]
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })
      await deployed.TBTCTokenStub.mint(caller, requiredBalance)
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      assert(new BN(finalBalance).gtn(new BN(initialBalance)), 'caller balance should increase')
    })

    it('burns if fraud', async () => {
      const value = 1000000000000
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const caller = accounts[2]
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })
      await deployed.TBTCTokenStub.mint(caller, requiredBalance)
      await testInstance.setState(utils.states.FRAUD_LIQUIDATION_IN_PROGRESS)
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      expect(new BN(finalBalance)).to.eq.BN(initialBalance)
    })
  })

  describe('notifyCourtesyCall', async () => {
    let oraclePrice
    let lotSize
    let lotValue
    let undercollateralizedPercent

    before(async () => {
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))

      oraclePrice = await deployed.TBTCSystemStub.fetchOraclePrice.call()
      lotSize = await deployed.TBTCConstants.getLotSize.call()
      lotValue = lotSize.mul(oraclePrice)

      undercollateralizedPercent = await deployed.TBTCConstants.getUndercollateralizedPercent.call()
    })

    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
    })

    it('sets courtesy call state, sets the timestamp, and logs CourtesyCalled', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we subtract `1` to test collateralization less than undercollateralized
      // threshold (140%).
      const bondValue = undercollateralizedPercent.mul(lotValue).div(new BN(100)).sub(new BN(1))
      await deployed.KeepStub.setBondAmount(bondValue)

      await testInstance.notifyCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.COURTESY_CALL)

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('CourtesyCalled', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in active state', async () => {
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.notifyCourtesyCall(),
        'Can only courtesy call from active state'
      )
    })

    it('reverts if sufficiently collateralized', async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization equal undercollateralized threshold (140%).
      const bondValue = undercollateralizedPercent.mul(lotValue).div(new BN(100))
      await deployed.KeepStub.setBondAmount(bondValue)

      await expectThrow(
        testInstance.notifyCourtesyCall(),
        'Signers have sufficient collateral'
      )
    })
  })

  describe('exitCourtesyCall', async () => {
    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      const notifiedTime = blockTimestamp // not expired
      const fundedTime = blockTimestamp // not expired
      await deployed.KeepStub.setBondAmount(new BN('1000000000000000000000000', 10))
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1', 10))
      await testInstance.setState(utils.states.COURTESY_CALL)
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), fundedTime, '0x' + '00'.repeat(36))
      await testInstance.setLiquidationAndCourtesyInitated(0, notifiedTime)
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('transitions to active, and logs ExitedCourtesyCall', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.exitCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.ACTIVE)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('ExitedCourtesyCall', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in courtesy call state', async () => {
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.exitCourtesyCall(),
        'Not currently in courtesy call'
      )
    })

    it('reverts if the deposit term is expiring anyway', async () => {
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), 0, '0x' + '00'.repeat(36))

      await expectThrow(
        testInstance.exitCourtesyCall(),
        'Deposit is expiring'
      )
    })

    it('reverts if the deposit is still undercollateralized', async () => {
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))
      await deployed.KeepStub.setBondAmount(0)

      await expectThrow(
        testInstance.exitCourtesyCall(),
        'Deposit is still undercollateralized'
      )
    })
  })

  describe('notifyUndercollateralizedLiquidation', async () => {
    let oraclePrice
    let lotSize
    let lotValue
    let severelyUndercollateralizedPercent

    before(async () => {
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))

      oraclePrice = await deployed.TBTCSystemStub.fetchOraclePrice.call()
      lotSize = await deployed.TBTCConstants.getLotSize.call()
      lotValue = lotSize.mul(oraclePrice)

      severelyUndercollateralizedPercent = await deployed.TBTCConstants.getSeverelyUndercollateralizedPercent.call()
    })

    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('executes', async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization less than severely undercollateralized
      // threshold (120%).
      const bondValue = severelyUndercollateralizedPercent.mul(lotValue).div(new BN(100)).sub(new BN(1))
      await deployed.KeepStub.setBondAmount(bondValue)

      await testInstance.notifyUndercollateralizedLiquidation()
      // TODO: Add validations or cover with `reverts if the deposit is not
      // severely undercollateralized` test case.
    })

    it('reverts if not in active or courtesy call', async () => {
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.notifyUndercollateralizedLiquidation(),
        'Deposit not in active or courtesy call'
      )
    })

    it('reverts if the deposit is not severely undercollateralized', async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization equal severely undercollateralized threshold (120%).
      const bondValue = severelyUndercollateralizedPercent.mul(lotValue).div(new BN(100))
      await deployed.KeepStub.setBondAmount(bondValue)

      await expectThrow(
        testInstance.notifyUndercollateralizedLiquidation(),
        'Deposit has sufficient collateral'
      )
    })

    it('assert starts signer abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifyUndercollateralizedLiquidation()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyCourtesyTimeout', async () => {
    let courtesyTime
    let timer
    before(async () => {
      timer = await deployed.TBTCConstants.getCourtesyCallTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      courtesyTime = blockTimestamp - timer.toNumber() // has not expired
      await testInstance.setState(utils.states.COURTESY_CALL)
      await testInstance.setLiquidationAndCourtesyInitated(0, courtesyTime)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('executes', async () => {
      await testInstance.notifyCourtesyTimeout()
    })

    it('reverts if not in a courtesy call period', async () => {
      await testInstance.setState(utils.states.START)
      await expectThrow(
        testInstance.notifyCourtesyTimeout(),
        'Not in a courtesy call period'
      )
    })

    it('reverts if the period has not elapsed', async () => {
      await testInstance.setLiquidationAndCourtesyInitated(0, courtesyTime * 5)
      await expectThrow(
        testInstance.notifyCourtesyTimeout(),
        'Courtesy period has not elapsed'
      )
    })

    it('assert starts signer abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifyCourtesyTimeout()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyDepositExpiryCourtesyCall', async () => {
    let timer
    let fundedTime

    before(async () => {
      timer = await deployed.TBTCConstants.getCourtesyCallTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundedTime = blockTimestamp - timer.toNumber() - 1 // has expired
      await testInstance.setState(utils.states.ACTIVE)
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), 0, '0x' + '00'.repeat(36))
    })

    it('sets courtesy call state, stores the time, and logs CourtesyCalled', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyDepositExpiryCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.COURTESY_CALL)

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('CourtesyCalled', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in active', async () => {
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.notifyDepositExpiryCourtesyCall(),
        'Deposit is not active'
      )
    })

    it('reverts if deposit not yet expiring', async () => {
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), fundedTime * 5, '0x' + '00'.repeat(36))

      await expectThrow(
        testInstance.notifyDepositExpiryCourtesyCall(),
        'Deposit term not elapsed'
      )
    })
  })

  describe('#attemptToLiquidateOnchain', async () => {
    let assertBalance
    let deposit
    let uniswapExchange
    let tbtcToken

    async function addLiquidity(ethAmount, tbtcAmount) {
      // Mint tBTC, mock liquidity
      // supply the equivalent of 10 actors posting liquidity
      const supplyFactor = new BN(10)

      const tbtcSupply = new BN(tbtcAmount).mul(supplyFactor)
      await tbtcToken.mint(accounts[0], tbtcSupply, { from: accounts[0] })
      await tbtcToken.approve(uniswapExchange.address, tbtcSupply, { from: accounts[0] })

      // Uniswap requires a minimum of 1000000000 wei for the initial addLiquidity call
      // So we add an extreme amount of ETH here.
      // We will use different values for keepBondAmount in future,
      // so multiplying by 1000000000x, while excessive, is robust.
      const UNISWAP_MINIMUM_INITIAL_LIQUIDITY_WEI = new BN('1000000000')
      const ethSupply = new BN(ethAmount).add(UNISWAP_MINIMUM_INITIAL_LIQUIDITY_WEI).mul(supplyFactor)

      // def addLiquidity(min_liquidity: uint256, max_tokens: uint256, deadline: timestamp) -> uint256:
      await uniswapExchange.addLiquidity(
        '0',
        tbtcSupply,
        UniswapHelpers.getDeadline(),
        { from: accounts[0], value: ethSupply }
      )
    }

    beforeEach(async () => {
      deposit = testInstance
      tbtcToken = deployed.TBTCTokenStub

      // Deploy Uniswap
      const uniswap = require('../uniswap')

      async function deployFromBytecode(abi, bytecode) {
        const Contract = new web3.eth.Contract(abi)
        const instance = await Contract
          .deploy({
            data: `0x`+bytecode,
          })
          .send({
            from: accounts[0],
            gas: 4712388,
          })
        return instance
      }

      // eslint-disable-next-line camelcase
      const uniswapExchange_web3 = await deployFromBytecode(uniswap.abis.exchange, uniswap.bytecode.exchange)
      // eslint-disable-next-line camelcase
      const uniswapFactory_web3 = await deployFromBytecode(uniswap.abis.factory, uniswap.bytecode.factory)

      // Required for Uniswap to clone and create factories
      await uniswapFactory_web3.methods.initializeFactory(uniswapExchange_web3.options.address).send({ from: accounts[0] })

      // Create tBTC exchange
      const uniswapFactory = await IUniswapFactory.at(uniswapFactory_web3.options.address)
      await uniswapFactory.createExchange(tbtcToken.address)
      const tbtcExchangeAddr = await uniswapFactory.getExchange.call(tbtcToken.address)
      uniswapExchange = await IUniswapExchange.at(tbtcExchangeAddr)

      await deployed.TBTCSystemStub.setExternalAddresses(uniswapFactory.address)

      // Helpers
      assertBalance = new AssertBalanceHelpers(tbtcToken)


      await tbtcToken.clearBalance(deposit.address)
    })

    it('returns false if address(exchange) = 0x0', async () => {
      // Override and use mock
      const uniswapFactory = deployed.UniswapFactoryStub
      await uniswapFactory.setExchange('0x0000000000000000000000000000000000000000')
      await deployed.TBTCSystemStub.setExternalAddresses(uniswapFactory.address)

      const retval = await deposit.attemptToLiquidateOnchain.call()
      expect(retval).to.be.false
    })

    it('liquidates using Uniswap successfully', async () => {
      const ethAmount = web3.utils.toWei('0.2', 'ether') // 0.2 eth : 1 tBTC
      const tbtcAmount = '100000000'
      await addLiquidity(ethAmount, tbtcAmount)

      const minTbtcAmount = '100100000'

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

    it('returns false if cannot buy up enough tBTC', async () => {
      const ethAmount = web3.utils.toWei('0.2', 'ether') // 0.2 eth : 1 tBTC
      const tbtcAmount = '100000000'
      await addLiquidity(ethAmount, tbtcAmount)

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

  describe('liquidation flows', async () => {
    let assertBalance
    let uniswapExchange
    let deposit
    let keep
    let tbtcToken

    const DIGEST = '0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc'
    const beneficiaryReward = '100000'
    const keepBondAmount = 1000000
    const lotSize = '100000000'

    beforeEach(async () => {
      deposit = testInstance
      keep = deployed.KeepStub
      tbtcToken = deployed.TBTCTokenStub

      const uniswapFactory = deployed.UniswapFactoryStub
      uniswapExchange = await UniswapExchangeStub.new(deployed.TBTCTokenStub.address)
      await uniswapFactory.setExchange(uniswapExchange.address)
      await deployed.TBTCSystemStub.setExternalAddresses(uniswapFactory.address)

      await keep.send(keepBondAmount, { from: accounts[0] })

      // Helpers
      assertBalance = new AssertBalanceHelpers(tbtcToken)
    })

    describe('startLiquidation', async () => {
      it('redemption', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount)

        // give the keep/deposit eth
        // TODO(liamz): use createNewDeposit when it works
        await assertBalance.eth(keep.address, '1000000') // TODO(liamz): replace w/ keep.checkBondAmount()
        await assertBalance.tbtc(deposit.address, '0')

        const requestor = accounts[2]
        await deposit.setRequestInfo(
          requestor,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )

        await deposit.startLiquidation()

        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

        // deposit should have 0 eth
        await assertBalance.eth(deposit.address, '0')

        // tbtc distributions
        await assertBalance.tbtc(deposit.address, '0')
        await assertBalance.tbtc(requestor, lotSize)
        await assertBalance.tbtc(beneficiary, beneficiaryReward)
      })

      it('non-redemption', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount)

        const NON_REDEMPTION_ADDRESS = '0x' + '00'.repeat(20)
        await deposit.setRequestInfo(
          NON_REDEMPTION_ADDRESS,
          '0x' + '11'.repeat(20),
          0, 0, DIGEST
        )

        const tx = await deposit.startLiquidation()

        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

        // deposit should have 0 eth
        await assertBalance.eth(deposit.address, '0')

        // tbtc distributions
        await assertBalance.tbtc(deposit.address, '0')
        await assertBalance.tbtc(beneficiary, beneficiaryReward)

        // expect TBTC to be burnt from deposit's account
        const evs = await tbtcToken.getPastEvents({ fromBlock: tx.receipt.blockNumber })
        const ev = evs[1]
        expect(ev.event).to.equal('Transfer')
        expect(ev.returnValues.from).to.equal(deposit.address)
        expect(ev.returnValues.to).to.equal('0x0000000000000000000000000000000000000000')
        expect(ev.returnValues.value).to.equal(lotSize.toString())
      })

      it('not liquidated', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount + 1)

        const res = await deposit.startLiquidation()

        // assert not liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should not be LIQUIDATED').to.not.eq.BN(utils.states.LIQUIDATED)

        const liquidationTime = await deposit.getLiquidationAndCourtesyInitiated.call()
        const block = await web3.eth.getBlock(res.receipt.blockNumber)
        expect(liquidationTime[0]).to.eq.BN(block.timestamp)
      })
    })

    describe('startSignerFraudLiquidation', async () => {
      const maintainer = accounts[5]

      it('was liquidated', async () => {
        const remainingEthAfterLiquidation = new BN(10)
        const maintainerBalance1 = new BN(await web3.eth.getBalance(maintainer))
        await deployed.KeepStub.send(remainingEthAfterLiquidation, { from: accounts[0] })
        await uniswapExchange.setEthPrice(keepBondAmount)

        const res = await deposit.startSignerFraudLiquidation({ from: maintainer })
        const txCostEth = await getTxCost(res)

        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

        // assert msg.sender.transfer(address(this).balance)
        await assertBalance.eth(deposit.address, '0')
        const maintainerBalanceExpected = maintainerBalance1.sub(txCostEth).add(remainingEthAfterLiquidation)
        const maintainerBalance2 = await web3.eth.getBalance(maintainer)
        expect(maintainerBalance2).to.eq.BN(maintainerBalanceExpected)
      })

      it('was not liquidated', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount + 1)
        await deposit.startSignerFraudLiquidation()

        // assert liquidation in process
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be FRAUD_LIQUIDATION_IN_PROGRESS').to.eq.BN(utils.states.FRAUD_LIQUIDATION_IN_PROGRESS)
      })
    })

    describe('startSignerAbortLiquidation', async () => {
      it('was liquidated', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount - 10)

        await deposit.startSignerAbortLiquidation()

        // assert liquidated
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be LIQUIDATED').to.eq.BN(utils.states.LIQUIDATED)

        // check distributeEthToKeepGroup
        const keepGroupEth = await keep.keepGroupEth()
        expect(keepGroupEth).to.eq.BN(new BN('10'))
      })

      it('was not liquidated', async () => {
        await uniswapExchange.setEthPrice(keepBondAmount + 1)

        await deposit.startSignerAbortLiquidation()

        // assert liquidation in process
        const depositState = await deposit.getState.call()
        expect(depositState, 'Deposit state should be FRAUD_LIQUIDATION_IN_PROGRESS').to.eq.BN(utils.states.LIQUIDATION_IN_PROGRESS)
      })
    })
  })
})
