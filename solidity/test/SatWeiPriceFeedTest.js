const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {contract, accounts} = require("@openzeppelin/test-environment")
const {BN, expectRevert, constants} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

const SatWeiPriceFeed = contract.fromArtifact("SatWeiPriceFeed")
const MockMedianizer = contract.fromArtifact("MockMedianizer")

describe("SatWeiPriceFeed", async function() {
  let satWeiPriceFeed
  let ethbtc

  before(async () => {
    satWeiPriceFeed = await SatWeiPriceFeed.new()
    ethbtc = await MockMedianizer.new()

    await satWeiPriceFeed.initialize(accounts[0], ethbtc.address)
  })

  describe("#getPrice", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("returns correct satwei price feed value", async () => {
      const ethusd = new BN(150)
      const btcusd = new BN(7000)
      // multiplication before division since BN does not store decimal points
      const ethbtcPrice = ethusd.mul(new BN(10).pow(new BN(18))).div(btcusd)

      await ethbtc.setValue(ethbtcPrice)

      const price = await satWeiPriceFeed.getPrice()
      const expectedSatWeiPrice = btcusd
        .mul(new BN(10).pow(new BN(10)))
        .div(ethusd)

      expect(new BN(price)).to.be.bignumber.equal(new BN(expectedSatWeiPrice))
    })

    it("casts down medianizer input price to lower 128 bits", async () => {
      const ethbtcPrice =
        "0xdef00000000000000000000000000000000000000000000000000004000000"

      await ethbtc.setValue(ethbtcPrice)

      const price = await satWeiPriceFeed.getPrice()
      const expectedSatWei = new BN(10).pow(new BN(28)).div(new BN("67108864"))

      expect(price).to.be.bignumber.equal(expectedSatWei)
    })

    it("Reverts if there are no active price feeds", async () => {
      await ethbtc.setValue(0)
      await expectRevert(satWeiPriceFeed.getPrice(), "Price feed offline")
    })

    it("Retrieve price through array of empty feeds", async () => {
      // set non-zero value temporarily to avoid addEthBtcFeed revert
      await ethbtc.setValue(1)
      await satWeiPriceFeed.addEthBtcFeed(ethbtc.address, {from: accounts[0]})
      await ethbtc.setValue(0)

      // Array should now contain 2 inactive price feeds, append a working one
      const workingEthBtcFeed = await MockMedianizer.new()
      const ethBtc = "123"

      await workingEthBtcFeed.setValue(ethBtc)
      await satWeiPriceFeed.addEthBtcFeed(workingEthBtcFeed.address, {
        from: accounts[0],
      })

      const activeEthBtcFeed = await satWeiPriceFeed.getWorkingEthBtcFeed()
      // ensure the feed we read from is indeed workingBtcFeed and workingEthFeed
      expect(activeEthBtcFeed).to.equal(workingEthBtcFeed.address)

      const price = await satWeiPriceFeed.getPrice()
      const expectedSatWei = new BN(10).pow(new BN(28)).div(new BN(ethBtc))

      expect(price).to.be.bignumber.equal(expectedSatWei)
    })
  })

  describe("Add price feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Only allowed address can add feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        satWeiPriceFeed.addEthBtcFeed(med.address),
        "Caller must be tbtcSystem contract",
      )
    })

    it("Cannot add dead feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        satWeiPriceFeed.addEthBtcFeed(med.address, {from: accounts[0]}),
        "Cannot add inactive feed",
      )
    })
  })

  describe("Get Working ETH/BTC feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Returns first active feed", async () => {
      const workingEthBtcFeed = await MockMedianizer.new()
      const ethBtc = "123"

      await workingEthBtcFeed.setValue(ethBtc)
      await satWeiPriceFeed.addEthBtcFeed(workingEthBtcFeed.address, {
        from: accounts[0],
      })

      const activeEthBtcFeed = await satWeiPriceFeed.getWorkingEthBtcFeed()

      expect(activeEthBtcFeed).to.equal(workingEthBtcFeed.address)
    })

    it("Returns address(0) if no active feed found ", async () => {
      const activeEthBtcFeed = await satWeiPriceFeed.getWorkingEthBtcFeed()

      expect(activeEthBtcFeed).to.equal(ZERO_ADDRESS)
    })
  })
})
