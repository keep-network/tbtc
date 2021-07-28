const truffleContract = require("@truffle/contract")
const Kit = require("@celo/contractkit")

const TBTCSystem = artifacts.require("TBTCSystem")

const SatWeiPriceFeed = artifacts.require("SatWeiPriceFeed")
const ETHBTCPriceFeedMock = artifacts.require("ETHBTCPriceFeedMock")

const DepositFactory = artifacts.require("DepositFactory")
const Deposit = artifacts.require("Deposit")
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")
const VendingMachine = artifacts.require("VendingMachine")

const { contracts } = require("@keep-network/common.js")
const { readExternalContractAddress } = contracts

// Used for creating sortition pool.
const BondedECDSAKeepFactoryJson = require("@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json")

const { ETHBTCMedianizer } = require("./externals.js")

module.exports = async function (deployer, network, accounts) {
  const BondedECDSAKeepFactoryAddress = readExternalContractAddress(
    "@keep-network/keep-ecdsa",
    "BondedECDSAKeepFactory",
    deployer
  )

  // Don't enact this setup during unit testing.
  if (process.env.NODE_ENV == "test" && !process.env.INTEGRATION_TEST) return

  const keepThreshold = 3
  const keepGroupSize = 3

  console.debug(
    `Initializing TBTCSystem [${TBTCSystem.address}] with:\n` +
      `  keepFactory: ${BondedECDSAKeepFactoryAddress}\n` +
      `  depositFactory: ${DepositFactory.address}\n` +
      `  masterDepositAddress: ${Deposit.address}\n` +
      `  tbtcToken: ${TBTCToken.address}\n` +
      `  tbtcDepositToken: ${TBTCDepositToken.address}\n` +
      `  feeRebateToken: ${FeeRebateToken.address}\n` +
      `  vendingMachine: ${VendingMachine.address}\n` +
      `  keepThreshold: ${keepThreshold}\n` +
      `  keepSize: ${keepGroupSize}`
  )

  // System.
  const tbtcSystem = await TBTCSystem.deployed()
  await tbtcSystem.initialize(
    BondedECDSAKeepFactoryAddress,
    DepositFactory.address,
    Deposit.address,
    TBTCToken.address,
    TBTCDepositToken.address,
    FeeRebateToken.address,
    VendingMachine.address,
    keepThreshold,
    keepGroupSize
  )

  console.log("TBTCSystem initialized!")

  // Price feed.
  const satWeiPriceFeed = await SatWeiPriceFeed.deployed()
  if (network === "mainnet") {
    // Inject mainnet price feeds.
    await satWeiPriceFeed.initialize(tbtcSystem.address, ETHBTCMedianizer)
  } else {
    // Inject mock price feed as base.
    const ethBtcPriceFeedMock = await ETHBTCPriceFeedMock.deployed()
    await satWeiPriceFeed.initialize(
      tbtcSystem.address,
      ethBtcPriceFeedMock.address
    )
  }

  // Create sortition pool for new TBTCSystem.
  console.log(`Creating sortition pool for TBTCSystem: [${TBTCSystem.address}]`)
  const BondedECDSAKeepFactoryContract = truffleContract(
    BondedECDSAKeepFactoryJson
  )
  BondedECDSAKeepFactoryContract.setProvider(deployer.provider)

  const BondedECDSAKeepFactory = await BondedECDSAKeepFactoryContract.at(
    BondedECDSAKeepFactoryAddress
  )

  if (network === "alfajores") {
    await createSortitionPoolCelo(accounts[0], BondedECDSAKeepFactoryAddress)
  } else {
    await BondedECDSAKeepFactory.createSortitionPool(TBTCSystem.address, {
      from: accounts[0],
    })
  }

  const sortitionPoolContractAddress = await BondedECDSAKeepFactory.getSortitionPool.call(
    TBTCSystem.address
  )
  console.log(`Sortition pool address: [${sortitionPoolContractAddress}]`)

  // Initialize minimum bondable value in newly created sortition pool.
  console.log("Refreshing minimum bondable value")
  await tbtcSystem.refreshMinimumBondableValue()
}

async function createSortitionPoolCelo(account, BondedECDSAKeepFactoryAddress) {
  const celoKit = Kit.newKitFromWeb3(web3)

  const BondedECDSAKeepFactory = new web3.eth.Contract(
    BondedECDSAKeepFactoryJson.abi,
    BondedECDSAKeepFactoryAddress
  )

  const txObject = await BondedECDSAKeepFactory.methods.createSortitionPool(
    TBTCSystem.address
  )

  const tx = await celoKit.sendTransactionObject(txObject, { from: account })
  await tx.waitReceipt()
}
