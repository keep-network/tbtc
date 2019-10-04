// The script redeems TBTC for BTC.
//
// Format:
// truffle exec 4_redemption.js <DEPOSIT_ADDRESS> <OUTPUT_VALUE> <REQUESTER_PKH>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// OUTPUT_VALUE - value to be redeemed into BTC
// REQUESTER_PKH - public key hash of the user requesting redemption
// truffle exec 4_redemption.js 0xc536685ca46654f0e8e250382132b583d25e7fdd2e 200 0x3333333333333333333333333333333333333333

const Deposit = artifacts.require('./Deposit.sol')

module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const outputValue = process.argv[5]
  const outputValueBytes = web3.utils.padLeft(web3.utils.numberToHex(outputValue), 16)
  const requesterPKH = process.argv[6]

  let deposit

  try {
    deposit = await Deposit.at(depositAddress)
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  await deposit.requestRedemption(
    outputValueBytes,
    requesterPKH,
  ).catch((err) => {
    console.error(`requesting redemption failed: ${err}`)
    process.exit(1)
  })

  process.exit()
}
