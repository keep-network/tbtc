import expectThrow from './helpers/expectThrow'
import increaseTime from './helpers/increaseTime'
import {
  createSnapshot,
  restoreSnapshot,
} from './helpers/snapshot'

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TBTCSystem = artifacts.require('TBTCSystem')

const KeepRegistryStub = artifacts.require('KeepRegistryStub')
const ECDSAKeepVendorStub = artifacts.require('ECDSAKeepVendorStub')

contract('TBTCSystem', (accounts) => {
  let tbtcSystem
  let ecdsaKeepVendor
  let factory
  let vendingMachine

  before(async () => {
    const deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)

    ecdsaKeepVendor = await ECDSAKeepVendorStub.new()

    const keepRegistry = await KeepRegistryStub.new()
    await keepRegistry.setVendor(ecdsaKeepVendor.address)

    tbtcSystem = deployed.TBTCSystem

    factory = deployed.DepositFactory
    vendingMachine = deployed.TestVendingMachine

    await tbtcSystem.initialize(
      keepRegistry.address,
      factory.address,
      utils.address0,
      utils.address0,
      utils.address0,
      utils.address0,
      utils.address0,
      vendingMachine.address,
      1,
      1
    )
  })

  describe('requestNewKeep()', async () => {
    before(async () => {
      ecdsaKeepVendor = await ECDSAKeepVendorStub.new()

      const keepRegistry = await KeepRegistryStub.new()
      await keepRegistry.setVendor(ecdsaKeepVendor.address)

      tbtcSystem = await TBTCSystem.new(utils.address0)

      await tbtcSystem.initialize(
        keepRegistry.address
      )
    })

    it('sends caller as owner to open new keep', async () => {
      const expectedKeepOwner = accounts[2]

      await tbtcSystem.requestNewKeep(5, 10, { from: expectedKeepOwner })
      const keepOwner = await ecdsaKeepVendor.keepOwner.call()

      assert.equal(expectedKeepOwner, keepOwner, 'incorrect keep owner address')
    })

    it('returns keep address', async () => {
      const expectedKeepAddress = await ecdsaKeepVendor.keepAddress.call()

      const result = await tbtcSystem.requestNewKeep.call(5, 10)

      assert.equal(expectedKeepAddress, result, 'incorrect keep address')
    })
  })

  describe('setSignerFeeDivisor', async () => {
    it('sets the signer fee', async () => {
      await tbtcSystem.setSignerFeeDivisor(new BN('201'))

      const signerFeeDivisor = await tbtcSystem.getSignerFeeDivisor()
      expect(signerFeeDivisor).to.eq.BN(new BN('201'))
    })

    it('reverts if msg.sender != owner', async () => {
      await expectThrow(
        tbtcSystem.setSignerFeeDivisor(new BN('201'), { from: accounts[1] }),
        ''
      )
    })
  })

  describe('setLotSizes', async () => {
    it('sets a different lot size array', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      const lotSizes = [10**8, 10**6]
      await tbtcSystem.setLotSizes(lotSizes)

      const eventList = await tbtcSystem.getPastEvents('LotSizesUpdated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
      expect(eventList[0].returnValues._lotSizes).to.eql(['100000000', '1000000']) // deep equality check
    })

    it('reverts if lot size array is empty', async () => {
      const lotSizes = []
      await expectThrow(
        tbtcSystem.setLotSizes(lotSizes),
        'Lot size array must always contain 1BTC'
      )
    })

    it('reverts if lot size array does not contain a 1BTC lot size', async () => {
      const lotSizes = [10**7]
      await expectThrow(
        tbtcSystem.setLotSizes(lotSizes),
        'Lot size array must always contain 1BTC'
      )
    })
  })

  describe('emergencyPauseNewDeposits', async () => {
    let term

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('pauses new deposit creation', async () => {
      await tbtcSystem.emergencyPauseNewDeposits()

      const allowNewDeposits = await tbtcSystem.getAllowNewDeposits()
      expect(allowNewDeposits).to.equal(false)
    })

    it('reverts if msg.sender is not owner', async () => {
      await expectThrow(
        tbtcSystem.emergencyPauseNewDeposits({ from: accounts[1] }),
        'Ownable: caller is not the owner'
      )
    })

    it('does not allows new deposit re-activation before 10 days', async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber() - 10) // T-10 seconds. toNumber because increaseTime doesn't support BN

      await expectThrow(
        tbtcSystem.resumeNewDeposits(),
        'Deposits are still paused'
      )
    })

    it('allows new deposit creation after 10 days', async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber()) // 10 days
      await tbtcSystem.resumeNewDeposits()
      const allowNewDeposits = await tbtcSystem.getAllowNewDeposits()
      expect(allowNewDeposits).to.equal(true)
    })

    it('reverts if emergencyPauseNewDeposits has already been called', async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber()) // 10 days
      tbtcSystem.resumeNewDeposits()

      await expectThrow(
        tbtcSystem.emergencyPauseNewDeposits(),
        'emergencyPauseNewDeposits can only be called once'
      )
    })
  })
})
