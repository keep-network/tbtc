// The script provides the redemption proof to a Deposit
//
// Format:
// truffle exec 7_provide_redemption_proof.js <DEPOSIT_ADDRESS> <TX> <PROOF> <INDEX> <HEADER_CHAIN>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// TX - the bitcoin tx that purportedly contain the redemption output
// PROOF - the merkle proof of inclusion of the tx in the bitcoin block
// INDEX - the index of the tx in the Bitcoin block (1-indexed)
// HEADER_CHAIN - an array of tightly-packed bitcoin headers
//
// Usage:
// truffle exec 7_provide_redemption_proof.js

const Deposit = artifacts.require('./Deposit.sol')

module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const tx = process.argv[5]
  const proof = process.argv[6]
  const index = process.argv[7]
  const headerChain = process.argv[8]

  let deposit

  try {
    deposit = await Deposit.at(depositAddress)
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  await deposit.provideRedemptionProof(
    tx,
    proof,
    index,
    headerChain
  ).catch((err) => {
    console.error(`failed to provide redemption proof: ${err}`)
    process.exit(1)
  })

  process.exit()
}
