const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {contract, web3, accounts} = require("@openzeppelin/test-environment")
const {BN, expectRevert, constants} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

const BTCETHPriceFeed = contract.fromArtifact("BTCETHPriceFeed")
const MockMedianizer = contract.fromArtifact("MockMedianizer")

describe("BTCETHPriceFeed", async function() {
  let btcEthPriceFeed
  let btc
  let eth

  before(async () => {
    btcEthPriceFeed = await BTCETHPriceFeed.new()
    btc = await MockMedianizer.new()
    eth = await MockMedianizer.new()

    await btcEthPriceFeed.initialize(accounts[0], btc.address, eth.address)
  })

  describe("#getPrice", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("computes a ratio of the two medianizers", async () => {
      const btcUsd = "7152.55"
      const ethUsd = "142.28"

      await btc.setValue(web3.utils.toWei(btcUsd))
      await eth.setValue(web3.utils.toWei(ethUsd))

      const price = await btcEthPriceFeed.getPrice()

      // 7152.55 / 142.28 = 50.2709446162
      // 50.2709446162 * 10^10
      // 502,709,446,162 wei
      expect(price).to.be.bignumber.equal(new BN("502709446162"))
    })

    it("casts down each medianizer price to lower 128 bits", async () => {
      const btcPrice =
        "0xdef00000000000000000000000000000000000000000000000000000000004"
      const ethPrice =
        "0xabc00000000000000000000000000000000000000000000000000000000002"

      await btc.setValue(btcPrice)
      await eth.setValue(ethPrice)

      const price = await btcEthPriceFeed.getPrice()
      expect(price).to.be.bignumber.equal(new BN("20000000000"))
    })

    it("Reverts if there are no active price feeds", async () => {
      await btc.setValue(0)
      await eth.setValue(0)
      await expectRevert(btcEthPriceFeed.getPrice(), "Price feed offline")
    })

    it("Retrieve price through array of empty feeds", async () => {
      // set non-zero value temporarily to avoid addBtcUsdFeed/addEthUsdFeed revert
      await btc.setValue(1)
      await eth.setValue(1)
      await btcEthPriceFeed.addBtcUsdFeed(btc.address, {from: accounts[0]})
      await btcEthPriceFeed.addEthUsdFeed(eth.address, {from: accounts[0]})
      await btc.setValue(0)
      await eth.setValue(0)

      // Arrays should now each contain 2 inactive price feeds, append a working one
      const workingBtcFeed = await MockMedianizer.new()
      const workingEthFeed = await MockMedianizer.new()

      const btcUsd = "7152.55"
      const ethUsd = "142.28"

      await workingBtcFeed.setValue(web3.utils.toWei(btcUsd))
      await workingEthFeed.setValue(web3.utils.toWei(ethUsd))

      await btcEthPriceFeed.addBtcUsdFeed(workingBtcFeed.address, {
        from: accounts[0],
      })
      await btcEthPriceFeed.addEthUsdFeed(workingEthFeed.address, {
        from: accounts[0],
      })

      const activeBtcFeed = await btcEthPriceFeed.getWorkingBtcUsdFeed()
      const activeEthFeed = await btcEthPriceFeed.getWorkingEthUsdFeed()
      // ensure the feed we read from is indeed workingBtcFeed and workingEthFeed
      expect(activeBtcFeed).to.equal(workingBtcFeed.address)
      expect(activeEthFeed).to.equal(workingEthFeed.address)

      const price = await btcEthPriceFeed.getPrice()

      // 7152.55 / 142.28 = 50.2709446162
      // 50.2709446162 * 10^10
      // 502,709,446,162 wei
      expect(price).to.be.bignumber.equal(new BN("502709446162"))
    })
  })

  describe("Add price feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Only allowed address can add feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        btcEthPriceFeed.addBtcUsdFeed(med.address),
        "Caller must be tbtcSystem contract",
      )

      await expectRevert(
        btcEthPriceFeed.addBtcUsdFeed(med.address),
        "Caller must be tbtcSystem contract",
      )
    })

    it("Cannot add dead feed", async () => {
      const med = await MockMedianizer.new()

      await expectRevert(
        btcEthPriceFeed.addBtcUsdFeed(med.address, {from: accounts[0]}),
        "Cannot add inactive feed",
      )

      await expectRevert(
        btcEthPriceFeed.addEthUsdFeed(med.address, {from: accounts[0]}),
        "Cannot add inactive feed",
      )
    })
  })

  describe("Get Working BTC/ETH feed", async () => {
    beforeEach(createSnapshot)
    afterEach(restoreSnapshot)

    it("Returns first active feed", async () => {
      const workingBtcFeed = await MockMedianizer.new()
      const workingEthFeed = await MockMedianizer.new()

      const btcUsd = "7152.55"
      const ethUsd = "142.28"

      await workingBtcFeed.setValue(web3.utils.toWei(btcUsd))
      await workingEthFeed.setValue(web3.utils.toWei(ethUsd))

      await btcEthPriceFeed.addBtcUsdFeed(workingBtcFeed.address, {
        from: accounts[0],
      })
      await btcEthPriceFeed.addEthUsdFeed(workingEthFeed.address, {
        from: accounts[0],
      })

      const activeBtcFeed = await btcEthPriceFeed.getWorkingBtcUsdFeed()
      const activeEthFeed = await btcEthPriceFeed.getWorkingEthUsdFeed()

      expect(activeBtcFeed).to.equal(workingBtcFeed.address)
      expect(activeEthFeed).to.equal(workingEthFeed.address)
    })

    it("Returns address(0) if no active feed found ", async () => {
      const activeBtcFeed = await btcEthPriceFeed.getWorkingBtcUsdFeed()
      const activeEthFeed = await btcEthPriceFeed.getWorkingEthUsdFeed()

      expect(activeBtcFeed).to.equal(ZERO_ADDRESS)
      expect(activeEthFeed).to.equal(ZERO_ADDRESS)
    })
  })
})
