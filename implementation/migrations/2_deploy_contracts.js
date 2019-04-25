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

module.exports = (deployer) => {
  deployer.then(async () => {
    await deployer.deploy(BytesLib)

    await deployer.link(BytesLib, [BTCUtils, ValidateSPV, Deposit, CheckBitcoinSigs, DepositUtils, DepositLiquidation, DepositFunding, DepositRedemption])
    await deployer.deploy(BTCUtils)

    await deployer.link(BTCUtils, [ValidateSPV, Deposit, CheckBitcoinSigs, DepositUtils, DepositLiquidation, DepositFunding, DepositRedemption])
    await deployer.deploy(ValidateSPV)
    await deployer.deploy(CheckBitcoinSigs)
    await deployer.deploy(TBTCConstants)

    await deployer.link(TBTCConstants, [OutsourceDepositLogging, Deposit, DepositUtils, DepositLiquidation, DepositFunding, DepositRedemption])
    await deployer.link(CheckBitcoinSigs, [Deposit, DepositUtils, DepositLiquidation, DepositFunding, DepositRedemption])
    await deployer.link(ValidateSPV, [Deposit, DepositUtils, DepositLiquidation, DepositFunding, DepositRedemption])
    await deployer.deploy(DepositStates)
    await deployer.deploy(OutsourceDepositLogging)

    await deployer.link(OutsourceDepositLogging, [Deposit, DepositRedemption, DepositFunding, DepositLiquidation])
    await deployer.link(DepositStates, [DepositUtils, DepositRedemption, DepositFunding, DepositLiquidation])
    await deployer.deploy(DepositUtils)

    await deployer.link(DepositUtils, [Deposit, DepositRedemption, DepositFunding, DepositLiquidation])
    await deployer.deploy(DepositLiquidation)

    await deployer.link(DepositLiquidation, [Deposit, DepositRedemption, DepositFunding])
    await deployer.deploy(DepositRedemption)
    await deployer.deploy(DepositFunding)

    await deployer.link(DepositFunding, [Deposit])
    await deployer.link(DepositRedemption, [Deposit])
    await deployer.deploy(Deposit)
  })
}
