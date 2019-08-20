const Dummy = artifacts.require('Dummy')
const CloneFactoryStub = artifacts.require('CloneFactoryStub')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const DEPLOY = [
  { name: 'CloneFactoryStub', contract: CloneFactoryStub },
  { name: 'Dummy', contract: Dummy }]

contract('CloneFactory', (accounts) => {
  let deployed
  let masterContract

  before(async () => {
    deployed = await utils.deploySystem(DEPLOY)
    masterContract = deployed.Dummy
  })

  describe('createClone()', async () => {
    it('creates new clone', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1, 'eventList length should be 1')

      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      assert(web3.utils.isAddress(cloneAddress), 'cloneAddress should be an address')

      const cloneInstance = await Dummy.at(cloneAddress)
      const cloneInstanceState = await cloneInstance.getState()
      expect(cloneInstanceState, 'State should be unset').to.eq.BN(0)
    })

    it('does not impact master contract state', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents('ContractCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      const cloneInstance = await Dummy.at(cloneAddress)

      await cloneInstance.setState(1)

      const cloneInstanceState = await cloneInstance.getState()

      // master should still be in START
      const masterContractState = await masterContract.getState()
      expect(cloneInstanceState, 'state should be 1').to.eq.BN(1)
      expect(masterContractState, 'State should be unset').to.eq.BN(0)
    })
  })

  describe('isClone()', async () => {
    it('correctly checks if address is a clone', async () => {
      await deployed.CloneFactoryStub.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(masterContract.address, cloneAddress)
      expect(checkClone, 'isClone() should return true').to.be.true
    })

    it('correctly checks if address is not a clone', async () => {
      const masterContract2 = await Dummy.new()

      await deployed.CloneFactoryStub.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await deployed.CloneFactoryStub.getPastEvents(
        'ContractCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' }
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await deployed.CloneFactoryStub.isClone_exposed.call(
        masterContract2.address,
        cloneAddress
      )
      expect(checkClone, 'isClone() should return false').not.to.be.true
    })
  })
})
