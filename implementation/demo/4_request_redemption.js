// The script redeems TBTC for BTC.
//
// Format:
// truffle exec 4_redemption.js <DEPOSIT_ADDRESS> <OUTPUT_VALUE> <REQUESTER_PKH>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// OUTPUT_VALUE - value to be redeemed into BTC, in satoshis
// REQUESTER_PKH - public key hash of the user requesting redemption
// truffle exec demo/4_request_redemption.js 0x281447b37FFddEDE449B94edB212C49c9358D0AA 2 0x3333333333333333333333333333333333333333

const Deposit = artifacts.require('./Deposit.sol')
const TBTCToken = artifacts.require('./TBTCToken.sol')
const BN = web3.utils.BN


module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const outputValue = process.argv[5]
  const outputValueBytes = web3.utils.bytesToHex(new BN(outputValue).toArray('le', 8))
  const requesterPKH = process.argv[6]

  console.log(outputValueBytes)

  let deposit
  let tbtcToken

  try {
    deposit = await Deposit.at(depositAddress)
    tbtcToken = await TBTCToken.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  await tbtcToken.approve(deposit.address, new BN('20000000000000'))
    .catch((err) => {
      console.error(`TBTC approval failed: ${err}`)
    })

  await deposit.requestRedemption(
    outputValueBytes,
    requesterPKH
  ).catch((err) => {
    console.error(`requesting redemption failed: ${err}`)
    process.exit(1)
  })

  console.log('Redemption requested!')
}
