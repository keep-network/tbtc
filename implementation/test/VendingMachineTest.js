
import expectThrow from './helpers/expectThrow'
import { AssertBalance } from './helpers/assertBalance'

const DepositUtils = artifacts.require('DepositUtils')
const TestTBTCConstants = artifacts.require('TestTBTCConstants')

const TestVendingMachine = artifacts.require('TestVendingMachine')
const TestToken = artifacts.require('TestToken')
const TestDepositOwnerToken = artifacts.require('TestDepositOwnerToken')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
// eslint-disable-next-line
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name

  { name: 'TestVendingMachine', contract: TestVendingMachine },
]

contract('VendingMachine', (accounts) => {
  let vendingMachine
  let depositOwnerToken
  let tbtcToken

  let assertBalance

  const dotId = '1'

  before(async () => {
    // VendingMachine relies on linked libraries, hence we use deploySystem for consistency.
    const deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)

    vendingMachine = deployed.TestVendingMachine
    tbtcToken = await TestToken.new(vendingMachine.address)
    depositOwnerToken = await TestDepositOwnerToken.new()

    assertBalance = new AssertBalance(tbtcToken)

    await vendingMachine.setExteriorAddresses(tbtcToken.address, depositOwnerToken.address)
  })

  describe('#dotToTbtc', async () => {
    before(async () => {
      await depositOwnerToken.forceMint(accounts[0], dotId)
    })

    it('converts DOT to TBTC', async () => {
      await depositOwnerToken.approve(vendingMachine.address, dotId, { from: accounts[0] })

      await vendingMachine.dotToTbtc(dotId)

      await assertBalance.tbtc(accounts[0], '999999950000000000')
    })

    it.skip('fails if deposit not qualified', async () => {})

    it(`fails if DOT doesn't exist`, async () => {
      await expectThrow(
        vendingMachine.dotToTbtc(123345),
      )
    })

    it(`fails if DOT transfer not approved`, async () => {
      await expectThrow(
        vendingMachine.dotToTbtc(dotId),
      )
    })
  })
})
