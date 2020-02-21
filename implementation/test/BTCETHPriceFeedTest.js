const {
  createSnapshot,
  restoreSnapshot,
} = require("../testHelpers/helpers/snapshot.js")
const {contract, web3} = require("@openzeppelin/test-environment")
const {BN} = require("@openzeppelin/test-helpers")
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

    await btcEthPriceFeed.initialize(btc.address, eth.address)
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
  })
})
