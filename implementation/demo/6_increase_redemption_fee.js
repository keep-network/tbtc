// The script increases the redemption fee for a deposit redemption request.
//
// Format:
// truffle exec 6_increase_redemption_fee.js <DEPOSIT_ADDRESS> <PREV> <NEW>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// PREV - previous fee in satoshis
// NEW - new fee in satoshis
// truffle exec 6_increase_redemption_fee.js 0xc536685ca46654f0e8e250382132b583d25e7fdd2e

const Deposit = artifacts.require('./Deposit.sol')

module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const previousOutputValue = process.argv[5]
  const previousOutputValueBytes = web3.utils.padLeft(web3.utils.numberToHex(previousOutputValue), 16)
  const newOutputValue = process.argv[6]
  const newOutputValueBytes = web3.utils.padLeft(web3.utils.numberToHex(newOutputValue), 16)

  let deposit

  try {
    deposit = await Deposit.at(depositAddress)
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  console.log(`Increasing redemption fee from ${previousOutputValue} sats to ${newOutputValue} sats`)

  await deposit.increaseRedemptionFee(
    previousOutputValueBytes,
    newOutputValueBytes
  ).catch((err) => {
    console.error(`increasing redemption fee failed: ${err}`)
    process.exit(1)
  })

  process.exit()
}
