import expectThrow from './helpers/expectThrow'

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TBTCSystem = artifacts.require('TBTCSystem')

const KeepRegistryStub = artifacts.require('KeepRegistryStub')
const ECDSAKeepVendorStub = artifacts.require('ECDSAKeepVendorStub')

const DepositFunding = artifacts.require('DepositFunding')
const DepositLiquidation = artifacts.require('DepositLiquidation')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositUtils = artifacts.require('DepositUtils')
const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const DepositStates = artifacts.require('DepositStates')
const TBTCConstants = artifacts.require('TBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')
const TestDepositFactory = artifacts.require('TestDepositFactory')
const TestVendingMachine = artifacts.require('TestVendingMachine')

const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCSystem', contract: TBTCSystem, param: utils.address0 },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the names
  { name: 'DepositFactory', contract: TestDepositFactory, param: 'TBTCSystem' }, // we don't care about ACL param. Bypassed in test
  { name: 'TestVendingMachine', contract: TestVendingMachine, param: 'TBTCSystem' },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'TBTCConstants', contract: TBTCConstants },
  { name: 'TestDeposit', contract: TestDeposit, param: utils.address0 },
]

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

  describe('setAllowNewDeposits', async () => {
    it('sets allowNewDeposits', async () => {
      await tbtcSystem.setAllowNewDeposits(false)

      const allowNewDeposits = await tbtcSystem.getAllowNewDeposits()
      expect(allowNewDeposits).to.equal(false)
    })

    it('reverts if msg.sender != owner', async () => {
      await expectThrow(
        tbtcSystem.setAllowNewDeposits(false, { from: accounts[1] }),
        ''
      )
    })
  })
})
