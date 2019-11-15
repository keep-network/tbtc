// The script redeems TBTC for BTC.
//
// Format:
// truffle exec 4_redemption.js <DEPOSIT_ADDRESS> <OUTPUT_VALUE> <REQUESTER_PKH>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// OUTPUT_VALUE - value to be redeemed into BTC, in satoshis
// REQUESTER_PKH - public key hash of the user requesting redemption
// truffle exec demo/4_request_redemption.js 0x281447b37FFddEDE449B94edB212C49c9358D0AA 1000 0x3333333333333333333333333333333333333333

const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')
const TBTCToken = artifacts.require('./TBTCToken.sol')
const BN = web3.utils.BN

// We approve the Deposit contract to transfer the maximum number of tokens
// from the user's balance.
// Temporary solution until TBTC includes approveAndCall support.
// TODO: remove after https://github.com/keep-network/tbtc/issues/273 is merged.
const MAX_TOKEN_ALLOWANCE = (new BN(2)).pow(new BN(256)).sub(new BN(1))

module.exports = async function() {
  // Parse arguments
  const depositAddress = process.argv[4]
  const outputValue = process.argv[5]
  const outputValueBytes = new BN(outputValue).toBuffer('le', 8)
  const requesterPKH = process.argv[6]

  let deposit
  let depositLog
  let tbtcToken

  try {
    deposit = await Deposit.at(depositAddress)
    depositLog = await TBTCSystem.deployed()
    tbtcToken = await TBTCToken.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  await tbtcToken.approve(deposit.address, MAX_TOKEN_ALLOWANCE)
    .catch((err) => {
      console.error(`TBTC approval failed: ${err}`)
      process.exit(1)
    })

  async function logEvents(startBlockNumber) {
    const eventList = await depositLog.getPastEvents('RedemptionRequested', {
      fromBlock: startBlockNumber,
      toBlock: 'latest',
    })

    const {
      _depositContractAddress,
      _requester,
      _digest,
      _utxoSize,
      _requesterPKH,
      _requestedFee,
    } = eventList[0].returnValues

    console.debug(`Redemption requested for deposit: ${_depositContractAddress}`)
    console.debug(`Request details:`)
    console.debug(`\trequestor: ${_requester}`)
    console.debug(`\trequestor PKH: ${_requesterPKH}`)
    console.debug(`\tvalue: ${_utxoSize} satoshis`)
    console.debug(`\tfee: ${_requestedFee} satoshis`)
    console.debug(`Digest approved for signing: ${_digest}`)
  }

  const startBlockNumber = await web3.eth.getBlock('latest').number

  await deposit.requestRedemption(outputValueBytes, requesterPKH)
    .catch((err) => {
      console.error(`requesting redemption failed: ${err}`)
      process.exit(1)
    })

  await logEvents(startBlockNumber)
    .catch((err) => {
      console.error(`getting events log failed: ${err}`)
      process.exit(1)
    })

  process.exit()
}
