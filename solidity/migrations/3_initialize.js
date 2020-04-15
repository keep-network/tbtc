const TBTCSystem = artifacts.require("TBTCSystem")

const SatWeiPriceFeed = artifacts.require("SatWeiPriceFeed")
const MockSatWeiPriceFeed = artifacts.require("ETHBTCPriceFeedMock")

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
    // Inject medianizer intermediary.
    await satWeiPriceFeed.initialize(tbtcSystem.address, RopstenETHBTCPriceFeed)
  } else {
    // Inject mock price feeds.
    const mockSatWeiPriceFeed = await MockSatWeiPriceFeed.deployed()
    await satWeiPriceFeed.initialize(
      tbtcSystem.address,
      mockSatWeiPriceFeed.address,
    )
  }
}
