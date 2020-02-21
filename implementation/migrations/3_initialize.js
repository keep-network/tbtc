const TBTCSystem = artifacts.require("TBTCSystem")

const BTCETHPriceFeed = artifacts.require("BTCETHPriceFeed")
const MockBTCUSDPriceFeed = artifacts.require("BTCUSDPriceFeed")
const MockETHUSDPriceFeed = artifacts.require("ETHUSDPriceFeed")

const DepositFactory = artifacts.require("DepositFactory")
const Deposit = artifacts.require("Deposit")
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")
const VendingMachine = artifacts.require("VendingMachine")

const {
  BondedECDSAKeepVendorAddress,
  BTCUSDPriceFeed,
  ETHUSDPriceFeed,
} = require("./externals")

module.exports = async function(deployer, network) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == "test" && !process.env.INTEGRATION_TEST) return

  // System.
  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(
    BondedECDSAKeepVendorAddress,
    DepositFactory.address,
    Deposit.address,
    TBTCToken.address,
    TBTCDepositToken.address,
    FeeRebateToken.address,
    VendingMachine.address,
    3,
    3,
  )

  // Price feed.
  const btcEthPriceFeed = await BTCETHPriceFeed.deployed()
  if (network === "mainnet") {
    // Inject mainnet price feeds.
    await btcEthPriceFeed.initialize(BTCUSDPriceFeed, ETHUSDPriceFeed)
  } else {
    // Inject mock price feeds.
    const mockBtcPriceFeed = await MockBTCUSDPriceFeed.deployed()
    const mockEthPriceFeed = await MockETHUSDPriceFeed.deployed()
    await btcEthPriceFeed.initialize(
      mockBtcPriceFeed.address,
      mockEthPriceFeed.address,
    )
  }
}
