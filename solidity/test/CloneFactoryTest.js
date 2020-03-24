const {contract, web3} = require("@openzeppelin/test-environment")
const {BN} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

const Dummy = contract.fromArtifact("Dummy")
const CloneFactoryStub = contract.fromArtifact("CloneFactoryStub")

describe("CloneFactory", async function() {
  let cloneFactory
  let masterContract
  before(async () => {
    cloneFactory = await CloneFactoryStub.new()
    masterContract = await Dummy.new()
  })

  describe("createClone()", async () => {
    it("creates new clone", async () => {
      await cloneFactory.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await cloneFactory.getPastEvents(
        "ContractCloneCreated",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList.length, "eventList length should be 1").to.equal(1)

      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      expect(
        web3.utils.isAddress(cloneAddress),
        "cloneAddress should be an address",
      ).to.be.true

      const cloneInstance = await Dummy.at(cloneAddress)
      const cloneInstanceState = await cloneInstance.getState()
      expect(cloneInstanceState, "State should be unset").to.be.bignumber.equal(
        new BN("0"),
      )
    })

    it("does not impact master contract state", async () => {
      await cloneFactory.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await cloneFactory.getPastEvents(
        "ContractCloneCreated",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      const cloneAddress = eventList[0].returnValues.contractCloneAddress
      const cloneInstance = await Dummy.at(cloneAddress)

      await cloneInstance.setState(1)

      const cloneInstanceState = await cloneInstance.getState()

      // master should still be in START
      const masterContractState = await masterContract.getState()
      expect(cloneInstanceState, "state should be 1").to.eq.BN(1)
      expect(masterContractState, "State should be unset").to.eq.BN(0)
    })
  })

  describe("isClone()", async () => {
    it("correctly checks if address is a clone", async () => {
      await cloneFactory.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await cloneFactory.getPastEvents(
        "ContractCloneCreated",
        {fromBlock: blockNumber, toBlock: "latest"},
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await cloneFactory.isClone_exposed.call(
        masterContract.address,
        cloneAddress,
      )
      expect(checkClone, "isClone() should return true").to.be.true
    })

    it("correctly checks if address is not a clone", async () => {
      const masterContract2 = await Dummy.new()

      await cloneFactory.createClone_exposed(masterContract.address)

      const blockNumber = await web3.eth.getBlockNumber()
      const eventList = await cloneFactory.getPastEvents(
        "ContractCloneCreated",
        {fromBlock: blockNumber, toBlock: "latest"},
      )

      const cloneAddress = eventList[0].returnValues.contractCloneAddress

      const checkClone = await cloneFactory.isClone_exposed.call(
        masterContract2.address,
        cloneAddress,
      )
      expect(checkClone, "isClone() should return false").not.to.be.true
    })
  })
})
