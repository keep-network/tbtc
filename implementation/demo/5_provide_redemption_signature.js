// The script provides the redemption signature.
//
// Format:
// truffle exec 5_provide_redemption_signature.js <DEPOSIT_ADDRESS> <V> <R> <S>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// V, R, S - parameters of signature verification
// truffle exec 5_provide_redemption_signature.js 0xc536685ca46654f0e8e250382132b583d25e7fdd2e 27 0xd7e83e8687ba8b555f553f22965c74e81fd08b619a7337c5c16e4b02873b537e 0x633bf745cdf7ae303ca8a6f41d71b2c3a21fcbd1aed9e7ffffa295c08918c1b3

const Deposit = artifacts.require('./Deposit.sol')

module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const v = process.argv[5]
  const r = process.argv[6]
  const s = process.argv[7]

  let deposit

  try {
    deposit = await Deposit.at(depositAddress)
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  await deposit.provideRedemptionSignature(
    v,
    r,
    s
  ).catch((err) => {
    console.error(`failed in providing redemption signature: ${err}`)
    process.exit(1)
  })

  process.exit()
}
