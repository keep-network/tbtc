import expectThrow from './helpers/expectThrow'

const utils = require('./utils')

const TBTCSystem = artifacts.require('TBTCSystem')

const KeepRegistryStub = artifacts.require('KeepRegistryStub')
const ECDSAKeepVendorStub = artifacts.require('ECDSAKeepVendorStub')

const DepositFunding = artifacts.require('DepositFunding')
const DepositLiquidation = artifacts.require('DepositLiquidation')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositUtils = artifacts.require('DepositUtils')
const DepositStates = artifacts.require('DepositStates')
const TBTCConstants = artifacts.require('TBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')
const DepositFactory = artifacts.require('DepositFactory')

const TEST_DEPOSIT_DEPLOY = [
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'TBTCConstants', contract: TBTCConstants },
  { name: 'TestDeposit', contract: TestDeposit },
]

contract('TBTCSystem', (accounts) => {
  let tbtcSystem
  let ecdsaKeepVendor

  describe('requestNewKeep()', async () => {
    before(async () => {
      const deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)

      ecdsaKeepVendor = await ECDSAKeepVendorStub.new()

      const keepRegistry = await KeepRegistryStub.new()
      await keepRegistry.setVendor(ecdsaKeepVendor.address)

      const depositFactory = await DepositFactory.new(deployed.TestDeposit.address)
      tbtcSystem = await TBTCSystem.new()

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
})
