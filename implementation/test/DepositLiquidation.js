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

  beforeEach(async () => {
    await testInstance.reset()
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
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.purchaseSignerBondsAtAuction()
      } catch (e) {
        assert.include(e.message, 'No active auction')
      }
    })

    it('reverts if TBTC balance is insufficient', async () => {
      // mint 1 less than lot size
      const lotSize = await deployed.TBTCConstants.getLotSize.call()
      await deployed.TBTCTokenStub.mint(accounts[0], lotSize - 1)
      try {
        await testInstance.purchaseSignerBondsAtAuction()
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Not enough TBTC to cover outstanding debt')
      }
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
    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('sets courtesy call state, sets the timestamp, and logs CourtesyCalled', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.COURTESY_CALL)

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('CourtesyCalled', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in active state', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Can only courtesy call from active state')
      }
    })

    it('reverts if sufficiently collateralized', async () => {
      try {
        await deployed.KeepStub.setBondAmount(1000)
        await deployed.TBTCSystemStub.setOraclePrice(new BN(1))
        await testInstance.notifyCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Signers have sufficient collateral')
      }
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
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Not currently in courtesy call')
      }
    })

    it('reverts if the deposit term is expiring anyway', async () => {
      try {
        await testInstance.setUTXOInfo('0x' + '00'.repeat(8), 0, '0x' + '00'.repeat(36))
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is expiring')
      }
    })

    it('reverts if the deposit is still undercollateralized', async () => {
      try {
        await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))
        await deployed.KeepStub.setBondAmount(0)
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is still undercollateralized')
      }
    })
  })

  describe('notifyUndercollateralizedLiquidation', async () => {
    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.TBTCSystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('executes', async () => {
      await testInstance.notifyUndercollateralizedLiquidation()
    })

    it('reverts if not in active or courtesy call', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyUndercollateralizedLiquidation()
      } catch (e) {
        assert.include(e.message, 'Deposit not in active or courtesy call')
      }
    })

    it('reverts if the deposit is not severely undercollateralized', async () => {
      try {
        await deployed.KeepStub.setBondAmount(10000000)
        await deployed.TBTCSystemStub.setOraclePrice(new BN(1))
        await testInstance.notifyUndercollateralizedLiquidation()
      } catch (e) {
        assert.include(e.message, 'Deposit has sufficient collateral')
      }
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
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyCourtesyTimeout()
      } catch (e) {
        assert.include(e.message, 'Not in a courtesy call period')
      }
    })

    it('reverts if the period has not elapsed', async () => {
      try {
        await testInstance.setLiquidationAndCourtesyInitated(0, courtesyTime * 5)
        await testInstance.notifyCourtesyTimeout()
      } catch (e) {
        assert.include(e.message, 'Courtesy period has not elapsed')
      }
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
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyDepositExpiryCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is not active')
      }
    })

    it('reverts if deposit not yet expiring', async () => {
      try {
        await testInstance.setUTXOInfo('0x' + '00'.repeat(8), fundedTime * 5, '0x' + '00'.repeat(36))
        await testInstance.notifyDepositExpiryCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit term not elapsed')
      }
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

  describe('startLiquidation', async () => {
    // it('seizes signer bonds')

    describe('redemption', async () => {
      // _tbtc.transferFrom(address(this), _d.requesterAddress, TBTCConstants.getLotSize());

      // _d.distributeBeneficiaryReward();

      // setLiquidated
    })

    describe('non-redemption', async () => {
      // maintain supply peg
      // _tbtc.burnFrom(address(this), TBTCConstants.getLotSize());

      // _d.distributeBeneficiaryReward();

      // setLiquidated
    })

    describe('was not liquidated', async () => {
      // _d.liquidationInitiated = block.timestamp;
    })
  })

  describe('startSignerFraudLiquidation', async () => {
    it('was liquidated', async () => {
      // msg.sender.transfer(address(this).balance);
    })

    it('was not liquidated', async () => {
      // _d.setFraudLiquidationInProgress();
    })
  })

  describe('startSignerAbortLiquidation', async () => {
    it('was liquidated', async () => {
      // _d.pushFundsToKeepGroup(address(this).balance);
    })

    it('was not liquidated', async () => {
      // _d.setLiquidationInProgress();
    })
  })
})

// * tests redemption flows of startSignerAbortLiquidation/startSignerFraudLiquidation
// * tests non-redemption flows of startSignerAbortLiquidation/startSignerFraudLiquidation
// * changes to fraud/abort liquidation logic, auctionTBTCAmount as detailed in [this comment](https://github.com/keep-network/tbtc/issues/105#issuecomment-511513898)
