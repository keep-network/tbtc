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

const bitcoinMain = {
  genesis:
    "0x00006020dd02d03c03dbc1f41312a6940e89919ce67fbf99a20307000000000000000000260d70e7ae07c80db07fbf29d09ec1a86d4f788e58098189a6f9021236572a7dd99eb15e397a11178294a823",
  height: 629070,
  epochStart:
    "0x459ec50d4ea62a89da04eb1ef3e352ec740bca50e8a808000000000000000000",
}

const bitcoinTest = {
  genesis:
    "0x0000c0205d1103efc13e6647977e3d65f253c3e762451e9ca9b920517d000000000000008442a07bcde3292a888277ea6337ba5bbdfa808ae01535846f19d843144c8f60478bb15e7b41011a88be5d36",
  height: 1723030,
  epochStart:
    "0xe2657f702faa9470815005305c45b4be2271c22ade1348e6fe00000000000000",
}

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
    if (network === "mainnet") {
      const {genesis, height, epochStart} = bitcoinMain

      await deployer.deploy(OnDemandSPV, genesis, height, epochStart, 0)
      difficultyRelay = await OnDemandSPV.deployed()
    } else if (network == "keep_dev" || "ropsten") {
      const {genesis, height, epochStart} = bitcoinTest

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
