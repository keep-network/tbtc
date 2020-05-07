const TBTCSystem = artifacts.require("TBTCSystem")

const SatWeiPriceFeed = artifacts.require("SatWeiPriceFeed")
const ETHBTCPriceFeedMock = artifacts.require("ETHBTCPriceFeedMock")

const DepositFactory = artifacts.require("DepositFactory")
const Deposit = artifacts.require("Deposit")
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")
const VendingMachine = artifacts.require("VendingMachine")

const {
  BondedECDSAKeepVendorAddress,
  ETHBTCMedianizer,
  RopstenETHBTCPriceFeed,
} = require("./externals")

module.exports = async function(deployer, network) {
  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == "test" && !process.env.INTEGRATION_TEST) return

  const keepThreshold = 3
  const keepGroupSize = 3

  console.debug(
    `Initializing TBTCSystem [${TBTCSystem.address}] with:\n` +
      `  keepVendor: ${BondedECDSAKeepVendorAddress}\n` +
      `  depositFactory: ${DepositFactory.address}\n` +
      `  masterDepositAddress: ${Deposit.address}\n` +
      `  tbtcToken: ${TBTCToken.address}\n` +
      `  tbtcDepositToken: ${TBTCDepositToken.address}\n` +
      `  feeRebateToken: ${FeeRebateToken.address}\n` +
      `  vendingMachine: ${VendingMachine.address}\n` +
      `  keepThreshold: ${keepThreshold}\n` +
      `  keepSize: ${keepGroupSize}`,
  )

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
    keepThreshold,
    keepGroupSize,
  )

  console.log("TBTCSystem initialized!")

  // Price feed.
  const satWeiPriceFeed = await SatWeiPriceFeed.deployed()
  if (network === "mainnet") {
    // Inject mainnet price feeds.
    await satWeiPriceFeed.initialize(tbtcSystem.address, ETHBTCMedianizer)
  } else if (network === "ropsten") {
    // Inject mock price feed as base.
    const ethBtcPriceFeedMock = await ethBtcPriceFeedMock.deployed()
    await satWeiPriceFeed.initialize(
      tbtcSystem.address,
      ethBtcPriceFeedMock.address,
    )

    // Add medianizer intermediary.
    await satWeiPriceFeed.addEthBtcFeed(RopstenETHBTCPriceFeed)
    // Disable mock feed so medianizer intermediary is active feed until we
    // choose to muck with the price.
    await ethBtcPriceFeedMock.setValue(0)
  } else {
    // Inject mock price feeds.
    const ethBtcPriceFeedMock = await ETHBTCPriceFeedMock.deployed()
    await satWeiPriceFeed.initialize(
      tbtcSystem.address,
      ethBtcPriceFeedMock.address,
    )
  }
}
