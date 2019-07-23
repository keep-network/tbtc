const BitcoinSPV = require('./BitcoinSPV.js')
const BitcoinTxParser = require('./BitcoinTxParser.js')

async function getTransactionProof(txID, confirmations) {
  console.log('Get transaction proof...')

  if (txID.length != 64) {
    throw new Error(`invalid transaction id length [${txID.length}], required: [64]`)
  }

  bitcoinSPV = new BitcoinSPV()
  await bitcoinSPV.initialize()
    .catch((err) => {
      return Promise.reject(new Error(`failed to initialize bitcoin spv: ${err}`))
    })

  const spvProof = await bitcoinSPV.getProof(txID, confirmations)
    .catch((err) => {
      return Promise.reject(new Error(`failed to get bitcoin spv proof: ${err}`))
    })

  const txDetails = await BitcoinTxParser.parse(spvProof.tx)
    .catch((err) => {
      return Promise.reject(new Error(`failed to parse spv proof: ${err}`))
    })

  return {
    version: txDetails.version,
    txInVector: txDetails.txInVector,
    txOutVector: txDetails.txOutVector,
    locktime: txDetails.locktime,
    merkleProof: spvProof.merkleProof,
    txInBlockIndex: spvProof.txInBlockIndex,
    chainHeaders: spvProof.chainHeaders,
  }
}

async function submitTransactionProof() {
  // TODO: Submit proof to a contract
}

module.exports.getTransactionProof = getTransactionProof
