const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const TBTCConstants = artifacts.require('TBTCConstants')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const Deposit = artifacts.require('Deposit')

module.exports = (deployer) => {
  deployer.then(async () => {
    await deployer.deploy(BytesLib)

    await deployer.link(BytesLib, [BTCUtils, ValidateSPV, Deposit, CheckBitcoinSigs])
    await deployer.deploy(BTCUtils)

    await deployer.link(BTCUtils, [ValidateSPV, Deposit, CheckBitcoinSigs])
    await deployer.deploy(ValidateSPV)
    await deployer.deploy(CheckBitcoinSigs)
    await deployer.deploy(TBTCConstants)

    await deployer.link(TBTCConstants, Deposit)
    await deployer.link(CheckBitcoinSigs, Deposit)
    await deployer.link(ValidateSPV, Deposit)
    await deployer.deploy(Deposit)
  })
}
