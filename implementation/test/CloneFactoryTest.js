const CloneFactoryTestDummy = artifacts.require('CloneFactoryTestDummy')
const CloneFactoryStub = artifacts.require('CloneFactoryStub')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const DEPLOY = [
  { name: 'CloneFactoryStub', contract: CloneFactoryStub },
  { name: 'CloneFactoryTestDummy', contract: CloneFactoryTestDummy }]

contract('CloneFactory', (accounts) => {
  let deployed
  let dummyContract

  before(async () => {
    deployed = await utils.deploySystem(DEPLOY)
    dummyContract = deployed.CloneFactoryTestDummy
  })

  describe('createClone()', async () => {
    it('creates new clone', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(dummyContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1, 'eventList length should be 1')

      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      assert(web3.utils.isAddress(cloneAddress), 'cloneAddress should be an address')

      const dummyContractInstance = await CloneFactoryTestDummy.at(cloneAddress)
      const dummyContractInstanceState = await dummyContractInstance.getState()
      expect(dummyContractInstanceState, 'State should be unset').to.eq.BN(0)
    })

    it('does not impact master contract state', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(dummyContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      const dummyContractInstance = await CloneFactoryTestDummy.at(cloneAddress)

      await dummyContractInstance.setState(1)

      const dummyContractInstanceState = await dummyContractInstance.getState()

      // master should still be in START
      const masterDummyState = await dummyContract.getState()
      expect(dummyContractInstanceState, 'state should be 1').to.eq.BN(1)
      expect(masterDummyState, 'State should be unset').to.eq.BN(utils.states.START)
    })
  })

  describe('isClone()', async () => {
    it('correctly checks if address is a clone', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(dummyContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(dummyContract.address, cloneAddress)
      assert(checkClone, 'isClone() should return true')
    })

    it('correctly checks if address is not a clone', async () => {
      const dummyContract2 = await CloneFactoryTestDummy.new()

      await deployed.CloneFactoryStub.createClone_exposed(dummyContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(
        dummyContract2.address,
        cloneAddress
      )
      assert(!checkClone, 'isClone() should return false')
    })
  })
})
