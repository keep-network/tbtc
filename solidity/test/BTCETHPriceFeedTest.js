const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {contract, accounts} = require("@openzeppelin/test-environment")
const {BN, expectRevert, constants} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

const BTCETHPriceFeed = contract.fromArtifact("BTCETHPriceFeed")
const MockMedianizer = contract.fromArtifact("MockMedianizer")

describe("BTCETHPriceFeed", async function() {
  let btcEthPriceFeed
  let btceth

  before(async () => {
    btcEthPriceFeed = await BTCETHPriceFeed.new()
    btceth = await MockMedianizer.new()

    await btcEthPriceFeed.initialize(accounts[0], btceth.address)
  })

  describe("#getPrice", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("returns BTCETH price feed value", async () => {
      const btcethPrice = "502709446162"

      await btceth.setValue(btcethPrice)

      const price = await btcEthPriceFeed.getPrice()

      expect(price).to.be.bignumber.equal(new BN(btcethPrice))
    })

    it("casts down each medianizer price to lower 128 bits", async () => {
      const btcethPrice =
        "0xdef00000000000000000000000000000000000000000000000000004000000"

      await btceth.setValue(btcethPrice)

      const price = await btcEthPriceFeed.getPrice()
      expect(price).to.be.bignumber.equal(new BN("67108864"))
    })

    it("Reverts if there are no active price feeds", async () => {
      await btceth.setValue(0)
      await expectRevert(btcEthPriceFeed.getPrice(), "Price feed offline")
    })

    it("Retrieve price through array of empty feeds", async () => {
      // set non-zero value temporarily to avoid addBtcEthFeed revert
      await btceth.setValue(1)
      await btcEthPriceFeed.addBtcEthFeed(btceth.address, {from: accounts[0]})
      await btceth.setValue(0)

      // Array should now contain 2 inactive price feeds, append a working one
      const workingBtcEthFeed = await MockMedianizer.new()
      const btcEth = "123"

      await workingBtcEthFeed.setValue(btcEth)
      await btcEthPriceFeed.addBtcEthFeed(workingBtcEthFeed.address, {
        from: accounts[0],
      })

      const activeBtcEthFeed = await btcEthPriceFeed.getWorkingBtcEthFeed()
      // ensure the feed we read from is indeed workingBtcFeed and workingEthFeed
      expect(activeBtcEthFeed).to.equal(workingBtcEthFeed.address)

      const price = await btcEthPriceFeed.getPrice()

      expect(price).to.be.bignumber.equal(btcEth)
    })
  })

  describe("Add price feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Only allowed address can add feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        btcEthPriceFeed.addBtcEthFeed(med.address),
        "Caller must be tbtcSystem contract",
      )
    })

    it("Cannot add dead feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        btcEthPriceFeed.addBtcEthFeed(med.address, {from: accounts[0]}),
        "Cannot add inactive feed",
      )
    })
  })

  describe("Get Working BTC/ETH feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Returns first active feed", async () => {
      const workingBtcEthFeed = await MockMedianizer.new()
      const btcEth = "123"

      await workingBtcEthFeed.setValue(btcEth)
      await btcEthPriceFeed.addBtcEthFeed(workingBtcEthFeed.address, {
        from: accounts[0],
      })

      const activeBtcEthFeed = await btcEthPriceFeed.getWorkingBtcEthFeed()

      expect(activeBtcEthFeed).to.equal(workingBtcEthFeed.address)
    })

    it("Returns address(0) if no active feed found ", async () => {
      const activeBtcEthFeed = await btcEthPriceFeed.getWorkingBtcEthFeed()

      expect(activeBtcEthFeed).to.equal(ZERO_ADDRESS)
    })
  })
})
