const DepositFactory = artifacts.require('DepositFactory')
const Deposit = artifacts.require('Deposit')

module.exports = (deployer) => {
  deployer.then(async () => {
    // proxy
    const deposit = await Deposit.deployed()
    await deployer.deploy(DepositFactory, deposit.address)
  })
}
