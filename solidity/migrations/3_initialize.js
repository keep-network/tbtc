const TBTCSystem = artifacts.require("TBTCSystem")

const BTCETHPriceFeed = artifacts.require("BTCETHPriceFeed")
const MockBTCETHPriceFeed = artifacts.require("BTCETHPriceFeedMock")

const DepositFactory = artifacts.require("DepositFactory")
const Deposit = artifacts.require("Deposit")
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")
const VendingMachine = artifacts.require("VendingMachine")

const {BondedECDSAKeepVendorAddress, BTCETHMedianizer} = require("./externals")

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
  const btcEthPriceFeed = await BTCETHPriceFeed.deployed()
  if (network === "mainnet") {
    // Inject mainnet price feeds.
    await btcEthPriceFeed.initialize(tbtcSystem.address, BTCETHMedianizer)
  } else {
    // Inject mock price feeds.
    const mockBtcEthPriceFeed = await MockBTCETHPriceFeed.deployed()
    await btcEthPriceFeed.initialize(
      tbtcSystem.address,
      mockBtcEthPriceFeed.address,
    )
  }
}
