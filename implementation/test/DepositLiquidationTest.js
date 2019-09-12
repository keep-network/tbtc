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
const TestToken = artifacts.require('TestToken')
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
  { name: 'TBTCSystemStub', contract: TBTCSystemStub }]

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

contract('Deposit', (accounts) => {
  let deployed
  let testInstance
  let beneficiary
  let tbtcToken

  before(async () => {
    beneficiary = accounts[2]
    deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
    tbtcToken = await TestToken.new(deployed.TBTCSystemStub.address)
    testInstance = deployed.TestDeposit
    testInstance.setExteroriorAddresses(deployed.TBTCSystemStub.address, tbtcToken.address, deployed.KeepStub.address)
    deployed.TBTCSystemStub.mint(beneficiary, web3.utils.toBN(deployed.TestDeposit.address))
  })

  beforeEach(async () => {
    await testInstance.reset()
  })

  describe('purchaseSignerBondsAtAuction', async () => {
    let lotSize

    before(async () => {
      lotSize = await deployed.TBTCConstants.getLotSize.call()
    })

    beforeEach(async () => {
      await testInstance.setState(utils.states.LIQUIDATION_IN_PROGRESS)
      for (let i = 0; i < 2; i++) {
        await tbtcToken.resetBalance(lotSize, { from: accounts[i] } )
        await tbtcToken.resetAllowance(testInstance.address, lotSize, { from: accounts[i] })
      }
    })

    it('sets state to liquidated, logs Liquidated, ', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.purchaseSignerBondsAtAuction()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.LIQUIDATED)

      const eventList = await deployed.TBTCSystemStub.getPastEvents('Liquidated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in a liquidation auction', async () => {
      await testInstance.setState(utils.states.START)

      await expectThrow(
        testInstance.purchaseSignerBondsAtAuction(),
        'No active auction'
      )
    })

    it('reverts if TBTC balance is insufficient', async () => {
      // burn 1 from caller to make balance insufficient
      await tbtcToken.forceBurn(accounts[0], 1)

      await expectThrow(
        testInstance.purchaseSignerBondsAtAuction(),
        'Not enough TBTC to cover outstanding debt'
      )
    })

    it(`burns msg.sender's tokens`, async () => {
      const caller = accounts[1]
      const initialTokenBalance = await tbtcToken.balanceOf(caller)

      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await tbtcToken.balanceOf(caller)
      const tokenCheck = new BN(finalTokenBalance).add(new BN(lotSize))
      expect(tokenCheck, 'tokens not burned correctly').to.eq.BN(initialTokenBalance)
    })

    it('distributes beneficiary reward', async () => {
      const caller = accounts[1]

      // Make sure Deposit has enough to cover beneficiary reward
      const beneficiaryReward = await deployed.TestDepositUtils.beneficiaryReward.call()
      await tbtcToken.forceMint(testInstance.address, beneficiaryReward)

      const initialTokenBalance = await tbtcToken.balanceOf(beneficiary)

      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await tbtcToken.balanceOf(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add(new BN(beneficiaryReward))

      expect(finalTokenBalance, 'tokens not returned to beneficiary correctly').to.eq.BN(tokenCheck)
    })

    it('distributes value to the caller', async () => {
      const value = 1000000000000
      const caller = accounts[1]
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const initialBalance = await web3.eth.getBalance(caller)

      await testInstance.send(value, { from: accounts[0] })

      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(caller)

      expect(new BN(finalBalance), 'caller balance should increase').to.be.gte.BN(initialBalance)
    })

    it('returns keep funds if not fraud', async () => {
      const value = 1000000000000
      const caller = accounts[1]
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })

      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      assert(new BN(finalBalance).gtn(new BN(initialBalance)), 'caller balance should increase')
    })

    it('burns if fraud', async () => {
      const value = 1000000000000
      const caller = accounts[1]
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })

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
})
