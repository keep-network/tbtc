// The script uses merkle.py script from bitcoin-spv to get part of the data for
// the proof.
//
// Format:
// truffle exec provide_funding_proof.js <TX_ID> <HEADERS_COUNT>
//
// Arguments:
// TX_ID - id of the funding transaction
// HEADERS_COUNT - number of block headers required for the proof, it's a number
//                 of confirmations for the transaction

const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystemStub.sol')

const FundingProof = require('./tools/FundingProof')

module.exports = async function() {
  const txID = process.argv[4]
  const headersCount = process.argv[5]

  let deposit
  let depositLog

  try {
    deposit = await Deposit.deployed()
    depositLog = await TBTCSystem.deployed()
  } catch (err) {
    throw new Error('contracts initialization failed', err)
  }

  async function provideFundingProof(fundingProof) {
    console.log('Submit funding proof...')

    const blockNumber = await web3.eth.getBlock('latest').number

    const result = await deposit.provideBTCFundingProof(
      fundingProof.version,
      fundingProof.txInVector,
      fundingProof.txOutVector,
      fundingProof.locktime,
      fundingProof.fundingOutputIndex,
      fundingProof.merkleProof,
      fundingProof.txInBlockIndex,
      fundingProof.chainHeaders
    ).catch((err) => {
      throw new Error(`provideBTCFundingProof failed: ${err}`)
    })

    console.log('provideBTCFundingProof transaction: ', result.tx)

    const eventList = await depositLog.getPastEvents('Funded', {
      fromBlock: blockNumber,
      toBlock: 'latest',
    })

    console.log('Funding proof accepted for the deposit: ', eventList[0].returnValues._depositContractAddress)
  }

  const fundingProof = await FundingProof.getTransactionProof(txID, headersCount)
    .catch((err) => {
      console.error('getting transaction proof failed\n', err)
      process.exit(1)
    })

  console.log('Funding proof:', FundingProof.serialize(fundingProof))

  await provideFundingProof(fundingProof)
    .catch((err) => {
      console.error('funding proof submission failed\n', err)
      process.exit(1)
    })
}
