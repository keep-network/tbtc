// The script increases the redemption fee for a deposit redemption request.
//
// Format:
// truffle exec 6_increase_redemption_fee.js <DEPOSIT_ADDRESS>
//
// Arguments:
// DEPOSIT_ADDRESS - Address of Deposit contract instance
// truffle exec 6_increase_redemption_fee.js 0xc536685ca46654f0e8e250382132b583d25e7fdd2e

const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')

const BN = web3.utils.BN

async function run() {
  // Parse arguments
  const depositAddress = process.argv[4]

  let deposit
  let depositLog

  try {
    deposit = await Deposit.at(depositAddress)
    depositLog = await TBTCSystem.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  // Get the initialRedemptionFee from the first request made
  console.log(`Retrieving most recent redemption fee for deposit...`)
  const redemptionEvents = await depositLog.getPastEvents(
    'RedemptionRequested',
    {
      filter: { _depositContractAddress: depositAddress },
      fromBlock: 0,
      toBlock: 'latest',
    }
  ).catch((err) => {
    console.error(`failed to get past redemption requested events`)
    process.exit(1)
  })

  if (!redemptionEvents.length) {
    console.error(`no RedemptionRequested event found for deposit`)
    process.exit(1)
  }

  let latestRedemptionRequest
  if (redemptionEvents.length > 1) {
    latestRedemptionRequest = redemptionEvents[redemptionEvents.length - 1]
  } else {
    latestRedemptionRequest = redemptionEvents[0]
  }

  const fee = web3.utils.toBN(latestRedemptionRequest.returnValues._requestedFee)
  const uxtoSize = web3.utils.toBN(latestRedemptionRequest.returnValues._utxoSize)

  const previousOutputValueBytes = uxtoSize.sub(fee).toBuffer('le')
  const newFee = fee.mul(new BN(2))
  const newOutputValueBytes = uxtoSize.sub(newFee).toBuffer('le')

  console.log(`Current fee is ${fee.toString()} sats`)
  console.log(`Increasing fee to ${newFee.toString()} sats`)

  async function logEvents(startBlockNumber) {
    const eventList = await depositLog.getPastEvents('RedemptionRequested', {
      fromBlock: startBlockNumber,
      toBlock: 'latest',
    })

    const ev = eventList[eventList.length - 1]
    const {
      _digest,
    } = ev.returnValues

    console.log(`Digest approved for signing: ${_digest}`)
  }

  const startBlockNumber = await web3.eth.getBlock('latest').number

  await deposit.increaseRedemptionFee(previousOutputValueBytes, newOutputValueBytes)
    .catch((err) => {
      console.error(`increasing redemption fee failed: ${err}`)
      process.exit(1)
    })

  await logEvents(startBlockNumber)
    .catch((err) => {
      console.error('getting events log failed\n', err)
      process.exit(1)
    })

  process.exit()
}

module.exports = async function() {
  try {
    await run()
  } catch (ex) {
    console.error(ex)
    process.exit(1)
  }
}
