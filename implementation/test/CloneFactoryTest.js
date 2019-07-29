const KeepStub = artifacts.require('KeepStub')
const TBTCTokenStub = artifacts.require('TBTCTokenStub')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const Deposit = artifacts.require('Deposit')
const CloneFactoryStub = artifacts.require('CloneFactoryStub')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCTokenStub', contract: TBTCTokenStub },
  { name: 'CloneFactoryStub', contract: CloneFactoryStub },
  { name: 'TBTCSystemStub', contract: TBTCSystemStub }]

contract('CloneFactory', (accounts) => {
  let deployed
  let depositContract

  before(async () => {
    depositContract = await Deposit.deployed()
    deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
  })

  describe('createClone()', async () => {
    it('creates new clone', async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await deployed.CloneFactoryStub.createClone_exposed(depositContract.address)

      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1, 'eventList length should be 1')

      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      assert(web3.utils.isAddress(cloneAddress), 'cloneAddress should be an address')

      const depositInstance = await Deposit.at(cloneAddress)
      const depositInstanceState = await depositInstance.getCurrentState()
      expect(depositInstanceState, 'Deposit instance should be in START').to.eq.BN(utils.states.START)
    })
    it('does not impact master contract state', async () => {
      const keep = await KeepStub.new()
      const blockNumber = await web3.eth.getBlockNumber()

      await deployed.CloneFactoryStub.createClone_exposed(depositContract.address)

      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      const depositInstance = await Deposit.at(cloneAddress)

      await depositInstance.createNewDeposit(
        deployed.TBTCSystemStub.address,
        deployed.TBTCTokenStub.address,
        keep.address,
        1,
        1)

      // moves instance state to AWAITING_BTC_FUNDING_PROOF
      await depositInstance.retrieveSignerPubkey()

      const depositInstanceState = await depositInstance.getCurrentState()

      // master should still be in START
      const masterDepositState = await depositContract.getCurrentState()
      expect(depositInstanceState, 'Deposit 1 should be in AWAITING_BTC_FUNDING_PROOF').to.eq.BN(utils.states.AWAITING_BTC_FUNDING_PROOF)
      expect(masterDepositState, 'Deposit instance should be in START').to.eq.BN(utils.states.START)
    })
  })

  describe('isClone()', async () => {
    it('correctly checks if address is a clone', async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await deployed.CloneFactoryStub.createClone_exposed(depositContract.address)
      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(depositContract.address, cloneAddress)
      assert(checkClone, 'isClone() should return true')
    })
    it('correctly checks if address is not a clone', async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      const depositContract2 = await Deposit.new()

      await deployed.CloneFactoryStub.createClone_exposed(depositContract.address)

      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(
        depositContract2.address,
        cloneAddress
      )
      assert(!checkClone, 'isClone() should return false')
    })
  })
})
