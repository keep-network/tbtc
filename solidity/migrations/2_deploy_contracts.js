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
const BTCETHPriceFeed = artifacts.require("BTCETHPriceFeed")
const BTCUSDPriceFeed = artifacts.require("BTCUSDPriceFeed")
const ETHUSDPriceFeed = artifacts.require("ETHUSDPriceFeed")
const prices = require("./prices")

// Bitcoin difficulty relays.
const Relay = artifacts.require("@summa-tx/relay-sol/contracts/Relay")
const MockRelay = artifacts.require("MockRelay")

// system
const TBTCConstants = artifacts.require("TBTCConstants")
const TBTCSystem = artifacts.require("TBTCSystem")

// tokens
const TBTCToken = artifacts.require("TBTCToken")
const TBTCDepositToken = artifacts.require("TBTCDepositToken")
const FeeRebateToken = artifacts.require("FeeRebateToken")

// deposit factory
const DepositFactory = artifacts.require("DepositFactory")

const all = [
  BytesLib,
  BTCUtils,
  ValidateSPV,
  TBTCConstants,
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
  BTCETHPriceFeed,
  VendingMachine,
  FeeRebateToken,
]

const bitcoinMain = {
  genesis:
    "0x00000020d208b5e50a8d3bd3a87f7a238e3f196621d0f9ffb5f302000000000000000000ee3af51ad3643a8a109935b45d9ca32b1003cda41df39dd75a17e13ba13aff4211aa585d39301c174a8ead73",
  height: 590588,
  epochStart:
    "0x704de08dc5329269011b878835be108a8202a93a0a2a1c000000000000000000",
}

const bitcoinTest = {
  genesis:
    "0x0000ff3ffc663e3a0b12b4cc2c05a425bdaf51922ce090acd8fa3a8a180300000000000084080b23fc40476d284da49fedaea9f7cee3aba33a8bad1347fa54740a29f02752b4c45dfcff031a279c2b3a",
  height: 1607272,
  epochStart:
    "0x84a9ec3b82556297ea36d1377901ecaef0bb5a5cf683f9f05103000000000000",
}

module.exports = (deployer, network, accounts) => {
  deployer.then(async () => {
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
    await deployer.deploy(TBTCConstants)
    await deployer.link(TBTCConstants, all)

    // logging
    await deployer.deploy(OutsourceDepositLogging)
    await deployer.link(OutsourceDepositLogging, all)

    let difficultyRelay
    // price feeds
    if (network !== "mainnet") {
      // On mainnet, we use the MakerDAO-deployed price feeds.
      // See: https://github.com/makerdao/oracles-v2#live-mainnet-oracles
      // Otherwise, we deploy our own mock price feeds, which are simpler
      // to maintain.
      await deployer.deploy(BTCUSDPriceFeed)
      await deployer.deploy(ETHUSDPriceFeed)

      const btcPriceFeed = await BTCUSDPriceFeed.deployed()
      const ethPriceFeed = await ETHUSDPriceFeed.deployed()

      await btcPriceFeed.setValue(web3.utils.toWei(prices.BTCUSD))
      await ethPriceFeed.setValue(web3.utils.toWei(prices.ETHUSD))
    }

    // On mainnet and Ropsten, we use the Summa-built, Keep-operated relay;
    // see https://github.com/summa-tx/relays . On testnet, we use a local
    // mock.
    if (network === "mainnet") {
      const {genesis, height, epochStart} = bitcoinMain

      await deployer.deploy(Relay, genesis, height, epochStart)
      difficultyRelay = await Relay.deployed()
    } else if (network == "ropsten") {
      const {genesis, height, epochStart} = bitcoinTest

      await deployer.deploy(Relay, genesis, height, epochStart)
      difficultyRelay = await Relay.deployed()
    } else if (network == "keep_dev") {
      const {genesis, height, epochStart} = bitcoinTest

      await deployer.deploy(Relay, genesis, height, epochStart)
      difficultyRelay = await Relay.deployed()
    } else {
      await deployer.deploy(MockRelay)
      difficultyRelay = await MockRelay.deployed()
    }

    // TODO This should be dropped soon.
    await deployer.deploy(BTCETHPriceFeed)

    if (!difficultyRelay) {
      throw new Error("Difficulty relay not found.")
    }

    // system
    await deployer.deploy(
      TBTCSystem,
      BTCETHPriceFeed.address,
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
  })
}
