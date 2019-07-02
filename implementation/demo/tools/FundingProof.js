// Configure path to bitcoin-spv directory as `BITCOIN_SPV_DIR` environment
// variable, e.g:
// export BITCOIN_SPV_DIR="/Users/jakub/workspace/bitcoin-spv"

export async function getTransactionProof(txID, headersCount) {
  console.log('Get transaction proof...')

  if (txID == undefined || txID.length < 64) {
    throw new Error('missing or malformed argument')
  }

  console.log(`Transaction ID: ${txID}`)

  const spvProof = await getBitcoinSPVproof(txID, headersCount)
  const txDetails = parseTransaction(spvProof.tx)

  return {
    version: txDetails.version,
    txInVector: txDetails.txInVector,
    txOutVector: txDetails.txOutVector,
    locktime: txDetails.locktime,
    fundingOutputIndex: txDetails.fundingOutputIndex,
    merkleProof: spvProof.merkleProof,
    txInBlockIndex: spvProof.txInBlockIndex,
    chainHeaders: spvProof.chainHeaders,
  }
}

async function getBitcoinSPVproof(txID, headersCount) {
  console.log('Get bitcoin-spv proof...')

  const { spawn } = require('child_process')

  const spawnProcess = spawn(
    'pipenv',
    ['run', 'python', 'scripts/merkle.py', txID, headersCount],
    { cwd: process.env.BITCOIN_SPV_DIR }
  )

  return new Promise((resolve, reject) => {
    spawnProcess.stdout.on('data', (data) => {
      console.log(`Received data from bitcoin-spv`)

      const spvProof = parseBitcoinSPVOutput(data.toString())

      resolve(spvProof)
    })

    spawnProcess.stderr.on('data', (data) => {
      reject(new Error(`failure:\n${data}`))
    })
  })
}

function parseBitcoinSPVOutput(output) {
  console.log('Parse bitcoin-spv output...\n')

  const tx = hexToBytes(output.match(/(^-* TX -*$\n)^(.*)$/m)[2])
  const merkleProof = hexToBytes(output.match(/(^-* PROOF -*$\n)^(.*)$/m)[2])
  const txInBlockIndex = parseInt(output.match(/(^-* INDEX -*$\n)^(.*)$/m)[2])
  const chainHeaders = hexToBytes(output.match(/(^-* CHAIN -*$\n)^(.*)$/m)[2])

  return {
    tx: tx,
    merkleProof: merkleProof,
    txInBlockIndex: txInBlockIndex,
    chainHeaders: chainHeaders,
  }
}

export function parseTransaction(tx) {
  console.log('Parse transaction...\nTX:', bytesToHex(tx))

  if (tx.length == 0) {
    throw new Error('cannot decode the transaction')
  }

  const txDetails = {
    version: getTxVersion(tx),
    txInVector: getTxInputVector(tx),
    txOutVector: getTxOutputVector(tx),
    locktime: getTxLocktime(tx),
    fundingOutputIndex: getFundingOutputIndex(tx), // TODO: Find index in transaction based on deposit's public key
  }

  return txDetails
}

export function serialize(fundingProof) {
  return {
    version: bytesToHex(fundingProof.version),
    txInVector: bytesToHex(fundingProof.txInVector),
    txOutVector: bytesToHex(fundingProof.txOutVector),
    locktime: bytesToHex(fundingProof.locktime),
    fundingOutputIndex: fundingProof.fundingOutputIndex,
    merkleProof: bytesToHex(fundingProof.merkleProof),
    chainHeaders: bytesToHex(fundingProof.chainHeaders),
    txInBlockIndex: fundingProof.txInBlockIndex,
  }
}

// getTxPrefix returns a prefix of the transaction. For legacy transaction it is
// first 4 bytes (version) and for witness transaction it is first 6 bytes
// (version and flag).
function getTxPrefix(tx) {
  if (isWitnessFlagPresent(tx)) {
    return tx.slice(0, 6)
  }
  return tx.slice(0, 4)
}

function getTxVersion(tx) {
  return tx.slice(0, 4)
}

// isWitnessFlagPresent checks if witness flag is present on the transaction.
// Witness flag is indicated by 2-byte value `0001` after the transaction version.
function isWitnessFlagPresent(tx) {
  if (tx.slice(4, 6).equals(Buffer.from('0001', 'hex'))) {
    return true
  }

  return false
}

// getTxInputVector returns vector of inputs for the transaction. It consists of
// number of inputs and list of the inputs.
function getTxInputVector(tx) {
  const txInVectorStartPosition = getTxPrefix(tx).length
  let txInVectorEndPosition

  if (isWitnessFlagPresent(tx)) {
    // TODO: Implement for witness transaction
    console.log('witness is not fully supported')
    txInVectorEndPosition = txInVectorStartPosition + 42
  } else {
    const inputCount = tx.slice(txInVectorStartPosition, txInVectorStartPosition + 1).readIntBE(0, 1)

    if (inputCount != 1) {
      // TODO: Support multiple inputs
      throw new Error(`exactly one input is required, got [${inputCount}]`)
    } else {
      const startPos = txInVectorStartPosition + 1

      const scriptLength = tx.slice(startPos + 36, startPos + 37).readIntBE(0, 1)
      if (scriptLength >= 253) {
        throw new Error('VarInts not supported')
      }

      txInVectorEndPosition = startPos + 37 + scriptLength + 4
    }
  }

  return tx.slice(txInVectorStartPosition, txInVectorEndPosition)
}

// getTxOutputVector returns vector of outputs for the transaction. It consists
// of number of outputs and list of the outputs.
function getTxOutputVector(tx) {
  const outStartPosition = getTxOutputVectorPosition(tx)
  const outputsCount = getNumberOfOutputs(tx)

  let startPosition = outStartPosition + 1
  let outEndPosition

  for (let i = 0; i < outputsCount; i++) {
    const scriptLength = tx.slice(startPosition + 8, startPosition + 8 + 1).readIntBE(0, 1)

    if (scriptLength >= 253) {
      throw new Error('VarInts not supported')
    }

    outEndPosition = startPosition + 8 + 1 + scriptLength
    startPosition = outEndPosition
  }

  return tx.slice(outStartPosition, outEndPosition)
}

// getTxOutputVector returns position on which outputs vector starts.
function getTxOutputVectorPosition(tx) {
  const txPrefix = getTxPrefix(tx)
  const txInput = getTxInputVector(tx)

  return txPrefix.length + txInput.length
}

function getNumberOfOutputs(tx) {
  const outStartPosition = getTxOutputVectorPosition(tx)

  return tx.slice(outStartPosition, outStartPosition + 1).readIntBE(0, 1)
}

function getFundingOutputIndex(tx) {
  // TODO: Implement
  console.log('get funding output index not implemented - always expect 0')
  return 0
}

function getTxLocktime(tx) {
  return tx.slice(tx.length - 4)
}

export function hexToBytes(hex) {
  return Buffer.from(hex, 'hex')
}

export function bytesToHex(bytes) {
  const buffer = Buffer.from(bytes)
  return buffer.toString('hex')
}
