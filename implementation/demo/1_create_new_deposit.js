const Deposit = artifacts.require('./Deposit.sol')
const KeepBridge = artifacts.require('./KeepBridge.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')

module.exports = async function() {
  let deposit
  let tbtcSystem
  let keepBridge

  try {
    keepBridge = await KeepBridge.deployed()
    tbtcSystem = await TBTCSystem.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  async function createNewDeposit() {
    // This is a temporary solution to deploy a new contract. It will be replaced
    // by a deposit factory.
    deposit = await Deposit.new()
    console.log('new deposit deployed: ', deposit.address)

    const result = await deposit.createNewDeposit(
      tbtcSystem.address, // address _TBTCSystem,
      '0x0000000000000000000000000000000000000000', // address _TBTCToken,
      keepBridge.address, // address _KeepBridge,
      5, // uint256 _m,
      10 // uint256 _n
    ).catch((err) => {
      console.error(`createNewDeposit failed: ${err}`)
      process.exit(1)
    })

    console.log('createNewDeposit transaction: ', result.tx)
  }

  await createNewDeposit()

  process.exit()
}
