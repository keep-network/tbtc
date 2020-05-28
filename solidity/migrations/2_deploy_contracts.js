// bitcoin-spv
const BytesLib = artifacts.require("BytesLib")
const BTCUtils = artifacts.require("BTCUtils")
const ValidateSPV = artifacts.require("ValidateSPV")
const CheckBitcoinSigs = artifacts.require("CheckBitcoinSigs")

// logging
const OutsourceDepositLogging = artifacts.require("OutsourceDepositLogging")
const DepositLog = artifacts.require("DepositLog")

// deposit
const DepositStates = artifacts.require("DepositStates")
const DepositUtils = artifacts.require("DepositUtils")
const DepositFunding = artifacts.require("DepositFunding")
const DepositRedemption = artifacts.require("DepositRedemption")
const DepositLiquidation = artifacts.require("DepositLiquidation")
const Deposit = artifacts.require("Deposit")
const VendingMachine = artifacts.require("VendingMachine")

// price feed
const ETHBTCPriceFeedMock = artifacts.require("ETHBTCPriceFeedMock")
const prices = require("./prices")
const SatWeiPriceFeed = artifacts.require("SatWeiPriceFeed")

// Bitcoin difficulty relays.
const relayConfig = require("./relay-config.json")

const OnDemandSPV = artifacts.require(
  "@summa-tx/relay-sol/contracts/OnDemandSPV",
)
const TestnetRelay = artifacts.require(
  "@summa-tx/relay-sol/contracts/TestnetRelay",
)
const MockRelay = artifacts.require("MockRelay")

// system
const TBTCConstants = artifacts.require("TBTCConstants")
const TBTCDevelopmentConstants = artifacts.require("TBTCDevelopmentConstants")
const KeepFactorySelection = artifacts.require("KeepFactorySelection")
const TBTCSystem = artifacts.require("TBTCSystem")

// tokens
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")

// deposit factory
const DepositFactory = artifacts.require("DepositFactory")

// scripts
const FundingScript = artifacts.require("FundingScript")
const RedemptionScript = artifacts.require("RedemptionScript")

const all = [
  BytesLib,
  BTCUtils,
  ValidateSPV,
  CheckBitcoinSigs,
  OutsourceDepositLogging,
  DepositLog,
  DepositStates,
  DepositUtils,
  DepositFunding,
  DepositRedemption,
  DepositLiquidation,
  Deposit,
  TBTCSystem,
  SatWeiPriceFeed,
  VendingMachine,
  FeeRebateToken,
]

module.exports = (deployer, network, accounts) => {
  deployer.then(async () => {
    let constantsContract = TBTCConstants
    if (network == "keep_dev" || network == "development") {
      // For keep_dev and development, replace constants with testnet constants.
      // Masquerade as TBTCConstants like a sinister fellow.
      TBTCDevelopmentConstants._json.contractName = "TBTCConstants"
      constantsContract = TBTCDevelopmentConstants
    }
    all.push(constantsContract)

    // bitcoin-spv
    await deployer.deploy(BytesLib)
    await deployer.link(BytesLib, all)

    await deployer.deploy(BTCUtils)
    await deployer.link(BTCUtils, all)

    await deployer.deploy(ValidateSPV)
    await deployer.link(ValidateSPV, all)

    await deployer.deploy(CheckBitcoinSigs)
    await deployer.link(CheckBitcoinSigs, all)

    // constants
    await deployer.deploy(constantsContract)
    await deployer.link(constantsContract, all)

    // logging
    await deployer.deploy(OutsourceDepositLogging)
    await deployer.link(OutsourceDepositLogging, all)

    let difficultyRelay
    // price feeds
    if (network !== "mainnet") {
      // On mainnet, we use the MakerDAO-deployed price feed.
      // See: https://github.com/makerdao/oracles-v2#live-mainnet-oracles
      // Otherwise, we deploy our own mock price feeds, which are simpler
      // to maintain.
      await deployer.deploy(ETHBTCPriceFeedMock)
      const ethBtcPriceFeedMock = await ETHBTCPriceFeedMock.deployed()
      await ethBtcPriceFeedMock.setValue(prices.satwei)
    }

    // On mainnet and Ropsten, we use the Summa-built, Keep-operated relay;
    // see https://github.com/summa-tx/relays . On testnet, we use a local
    // mock.
    if (network === "mainnet" || relayConfig.forceOnDemandSPV === true) {
      const {genesis, height, epochStart} = relayConfig.init.bitcoinMain

      await deployer.deploy(OnDemandSPV, genesis, height, epochStart, 0)
      difficultyRelay = await OnDemandSPV.deployed()
    } else if (
      network === "keep_dev" ||
      network === "ropsten" ||
      relayConfig.forceTestnetRelay === true
    ) {
      const {genesis, height, epochStart} = relayConfig.init.bitcoinTest

      await deployer.deploy(TestnetRelay, genesis, height, epochStart, 0)
      difficultyRelay = await TestnetRelay.deployed()
    } else {
      await deployer.deploy(MockRelay)
      difficultyRelay = await MockRelay.deployed()
    }

    await deployer.deploy(SatWeiPriceFeed)

    if (!difficultyRelay) {
      throw new Error("Difficulty relay not found.")
    }

    await deployer.deploy(KeepFactorySelection)
    await deployer.link(KeepFactorySelection, TBTCSystem)

    // system
    await deployer.deploy(
      TBTCSystem,
      SatWeiPriceFeed.address,
      difficultyRelay.address,
    )

    await deployer.deploy(DepositFactory, TBTCSystem.address)

    await deployer.deploy(VendingMachine, TBTCSystem.address)

    // deposit
    await deployer.deploy(DepositStates)
    await deployer.link(DepositStates, all)

    await deployer.deploy(DepositUtils)
    await deployer.link(DepositUtils, all)

    await deployer.deploy(DepositLiquidation)
    await deployer.link(DepositLiquidation, all)

    await deployer.deploy(DepositRedemption)
    await deployer.link(DepositRedemption, all)

    await deployer.deploy(DepositFunding)
    await deployer.link(DepositFunding, all)

    await deployer.deploy(Deposit)

    // token
    await deployer.deploy(TBTCToken, VendingMachine.address)
    await deployer.deploy(TBTCDepositToken, DepositFactory.address)
    await deployer.deploy(FeeRebateToken, VendingMachine.address)

    // scripts
    await deployer.deploy(
      FundingScript,
      VendingMachine.address,
      TBTCToken.address,
      TBTCDepositToken.address,
      FeeRebateToken.address,
    )
    await deployer.deploy(
      RedemptionScript,
      VendingMachine.address,
      TBTCToken.address,
      FeeRebateToken.address,
    )
  })
}
