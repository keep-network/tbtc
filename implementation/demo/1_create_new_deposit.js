const Deposit = artifacts.require('./Deposit.sol')
const KeepBridge = artifacts.require('./KeepBridge.sol')
const TBTCSystem = artifacts.require('./TBTCSystemStub.sol')

module.exports = async function() {
  let deposit
  let tbtcSystem
  let keepBridge

  try {
    deposit = await Deposit.deployed()
    keepBridge = await KeepBridge.deployed()
    tbtcSystem = await TBTCSystem.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  async function createNewDeposit() {
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
