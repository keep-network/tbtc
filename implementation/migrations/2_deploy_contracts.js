const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const TBTCConstants = artifacts.require('TBTCConstants')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositLog = artifacts.require('DepositLog')
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')

const Deposit = artifacts.require('Deposit')

const KeepBridge = artifacts.require('KeepBridge')

const all = [BytesLib, BTCUtils, ValidateSPV, TBTCConstants, CheckBitcoinSigs,
  OutsourceDepositLogging, DepositLog, DepositStates, DepositUtils,
  DepositFunding, DepositRedemption, DepositLiquidation, Deposit, KeepBridge]

module.exports = (deployer) => {
  deployer.then(async () => {
    await deployer.deploy(BytesLib)

    await deployer.link(BytesLib, all)
    await deployer.deploy(BTCUtils)

    await deployer.link(BTCUtils, all)
    await deployer.deploy(ValidateSPV)
    await deployer.deploy(CheckBitcoinSigs)
    await deployer.deploy(TBTCConstants)
    await deployer.link(TBTCConstants, all)

    await deployer.link(CheckBitcoinSigs, all)
    await deployer.link(ValidateSPV, all)
    await deployer.deploy(DepositStates)
    await deployer.deploy(OutsourceDepositLogging)

    await deployer.link(OutsourceDepositLogging, all)
    await deployer.link(DepositStates, all)
    await deployer.deploy(DepositUtils)

    await deployer.link(DepositUtils, all)
    await deployer.deploy(DepositLiquidation)

    await deployer.link(DepositLiquidation, all)
    await deployer.deploy(DepositRedemption)
    await deployer.deploy(DepositFunding)

    await deployer.link(DepositFunding, all)
    await deployer.link(DepositRedemption, all)

    await deployer.deploy(Deposit)

    await deployer.deploy(KeepBridge)
  })
}
