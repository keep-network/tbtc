
import expectThrow from './helpers/expectThrow'
import { AssertBalance } from './helpers/assertBalance'
import { createSnapshot, restoreSnapshot } from './helpers/snapshot'

const DepositUtils = artifacts.require('DepositUtils')
const TestTBTCConstants = artifacts.require('TestTBTCConstants')

const TestVendingMachine = artifacts.require('TestVendingMachine')
const TestToken = artifacts.require('TestToken')
const TestDepositOwnerToken = artifacts.require('TestDepositOwnerToken')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
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

  // For lack of a better design, this is the amount of TBTC exchanged for DOT's.
  const depositValueLessSignerFee = '995000000000000000'

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
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('converts DOT to TBTC', async () => {
      await depositOwnerToken.forceMint(accounts[0], dotId)
      await depositOwnerToken.approve(vendingMachine.address, dotId, { from: accounts[0] })

      await vendingMachine.dotToTbtc(dotId)

      await assertBalance.tbtc(accounts[0], depositValueLessSignerFee)
    })

    it.skip('fails if deposit not qualified', async () => { })

    it(`fails if DOT doesn't exist`, async () => {
      await expectThrow(
        vendingMachine.dotToTbtc(123345),
        'ERC721: operator query for nonexistent token.'
      )
    })

    it(`fails if DOT transfer not approved`, async () => {
      await depositOwnerToken.forceMint(accounts[0], dotId)

      await expectThrow(
        vendingMachine.dotToTbtc(dotId),
        'ERC721: transfer caller is not owner nor approved.'
      )
    })
  })

  describe('#tbtcToDot', async () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('converts TBTC to DOT', async () => {
      await depositOwnerToken.forceMint(vendingMachine.address, dotId)
      await tbtcToken.forceMint(accounts[0], depositValueLessSignerFee)
      await tbtcToken.approve(vendingMachine.address, depositValueLessSignerFee, { from: accounts[0] })

      const fromBlock = await web3.eth.getBlockNumber()
      await vendingMachine.tbtcToDot(dotId)

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock, toBlock: 'latest' })
      const tbtcBurntEvent = events[0]
      expect(tbtcBurntEvent.returnValues.from).to.equal(accounts[0])
      expect(tbtcBurntEvent.returnValues.to).to.equal(utils.address0)
      expect(tbtcBurntEvent.returnValues.value).to.equal(depositValueLessSignerFee)

      expect(
        await depositOwnerToken.ownerOf(dotId)
      ).to.equal(accounts[0])
    })

    it.skip('fails if deposit not qualified', async () => { })

    it(`fails if caller hasn't got enough TBTC`, async () => {
      await depositOwnerToken.forceMint(vendingMachine.address, dotId)

      await expectThrow(
        vendingMachine.tbtcToDot(dotId),
        'Not enough TBTC for DOT exchange.'
      )
    })

    it(`fails if deposit is locked`, async () => {
      // Deposit is locked if the Deposit Owner Token is not owned by the vending machine
      const depositOwner = accounts[1]
      await depositOwnerToken.forceMint(depositOwner, dotId)
      await tbtcToken.forceMint(accounts[0], depositValueLessSignerFee)
      await tbtcToken.approve(vendingMachine.address, depositValueLessSignerFee, { from: accounts[0] })

      await expectThrow(
        vendingMachine.tbtcToDot(dotId),
        'Deposit is locked.'
      )
    })
  })
})
