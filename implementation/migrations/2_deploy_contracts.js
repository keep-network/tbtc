// bitcoin-spv
const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

// logging
const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositLog = artifacts.require('DepositLog')

// deposit
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')
const Deposit = artifacts.require('Deposit')
const VendingMachine = artifacts.require('VendingMachine')

// price feed
const BTCETHPriceFeed = artifacts.require('BTCETHPriceFeed')
const BTCUSDPriceFeed = artifacts.require('BTCUSDPriceFeed')
const ETHUSDPriceFeed = artifacts.require('ETHUSDPriceFeed')
const prices = require('./prices')

// system
const TBTCConstants = artifacts.require('TBTCConstants')
const TBTCSystem = artifacts.require('TBTCSystem')

// tokens
const TBTCToken = artifacts.require('TBTCToken')
const TBTCDepositToken = artifacts.require('TBTCDepositToken')
const FeeRebateToken = artifacts.require('FeeRebateToken')

// deposit factory
const DepositFactory = artifacts.require('DepositFactory')

const all = [BytesLib, BTCUtils, ValidateSPV, TBTCConstants, CheckBitcoinSigs,
  OutsourceDepositLogging, DepositLog, DepositStates, DepositUtils,
  DepositFunding, DepositRedemption, DepositLiquidation, Deposit, TBTCSystem,
  BTCETHPriceFeed, VendingMachine, FeeRebateToken]

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

    await deployer.deploy(BTCETHPriceFeed)

    await deployer.deploy(TBTCSystem, BTCETHPriceFeed.address)

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

    await deployer.deploy(Deposit, DepositFactory.address)

    // price feeds
    if (network !== 'mainnet') {
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

    // token
    await deployer.deploy(TBTCToken, VendingMachine.address)
    await deployer.deploy(TBTCDepositToken, DepositFactory.address)
    await deployer.deploy(FeeRebateToken, VendingMachine.address)
  })
}
