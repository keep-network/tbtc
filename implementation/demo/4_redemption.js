// The script redeems TBTC for BTC.
//
// Format:
// truffle exec 4_redemption.js <DEPOSIT_ADDRESS> <OUTPUT_VALUE> <REQUESTER_PKH>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// OUTPUT_VALUE - value to be redeemed into BTC
// REQUESTER_PKH - public key hash of the user requesting redemption

// truffle exec 0xc536685ca46654f0e8e250382132b583d25e7fdd2e 0x1111111100000000 0x3333333333333333333333333333333333333333

const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')
const TBTCToken = artifacts.require('./TBTCToken.sol')

module.exports = async function() {
  const outputValue = process.argv[3]
  const requesterPKH = process.argv[4]
  const depositAddress = process.argv[5]

  /* eslint-disable no-unused-vars */
  let deposit
  let tbtcSystem
  let tbtcToken
  /* eslint-enable no-unused-vars */

  try {
    deposit = await Deposit.at(depositAddress)
    depositLog = await TBTCSystem.deployed()
    tbtcSystem = await TBTCSystem.deployed()
    tbtcToken = await TBTCToken.deployed()
  } catch (err) {
    throw new Error('contracts initialization failed', err)
  }

  await deposit.requestRedemption(
    outputValue, // _outputValueBytes
    requesterPKH, // _requesterPKH
  )

  process.exit()
}
