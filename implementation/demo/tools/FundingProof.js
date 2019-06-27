// Configure path to bitcoin-spv merkle proof script.
const merkleScript = '/Users/jakub/workspace/bitcoin-spv/scripts/merkle.py'

const fundingProof = {}

export function initialize() {
  console.log('Install python environment...')

  const { spawn } = require('child_process')

  const spawnProcess = spawn('pipenv', ['install'])

  spawnProcess.stdout.on('data', (data) => {
    console.log(`${data}`)
  })

  spawnProcess.stderr.on('data', (data) => {
    console.error(`Failure:\n${data}`)
    process.exit(1)
  })

  spawnProcess.on('close', (code) => {
    console.log(`child process exited with code ${code}`)
  })
}

export async function getTransactionProof(txID, headerLen, callback) {
  console.log('Get transaction proof...')

  if (txID == undefined || txID.length < 64) {
    console.error('missing txID argument')
    process.exit(1)
  }

  console.log(`Transaction ID: ${txID}`)

  await getBitcoinSPVproof(txID, headerLen, callback)
}

async function getBitcoinSPVproof(txID, headerLen, callback) {
  console.log('Get bitcoin-spv proof...')

  const { spawn } = require('child_process')

  const spawnProcess = spawn('pipenv', ['run', 'python', merkleScript, txID, headerLen])

  spawnProcess.stdout.on('data', (data) => {
    console.log(`Received data from bitcoin-spv`)
    const spvProof = parseBitcoinSPVOutput(data.toString())

    fundingProof.merkleProof = spvProof.merkleProof
    fundingProof.txInBlockIndex = spvProof.txInBlockIndex
    fundingProof.chainHeaders = spvProof.chainHeaders

    parseTransaction(spvProof.tx, callback)
  })

  spawnProcess.stderr.on('data', (data) => {
    console.error(`Failure:\n${data}`)
    process.exit(1)
  })

  spawnProcess.on('close', (code) => {
    console.log(`child process exited with code ${code}`)
    return
  })
}

function parseBitcoinSPVOutput(output) {
  console.log('Parse bitcoin-spv output...\n')

  const tx = hexToBytes(output.match(/(^-* TX -*$\n)^(.*)$/m)[2])
  const merkleProof = hexToBytes(output.match(/(^-* PROOF -*$\n)^(.*)$/m)[2])
  const txInBlockIndex = output.match(/(^-* INDEX -*$\n)^(.*)$/m)[2]
  const chainHeaders = hexToBytes(output.match(/(^-* CHAIN -*$\n)^(.*)$/m)[2])

  return {
    tx: tx,
    merkleProof: merkleProof,
    txInBlockIndex: txInBlockIndex,
    chainHeaders: chainHeaders,
  }
}

async function parseTransaction(tx, callback) {
  console.log('Parse transaction...\nTX:', bytesToHex(tx))

  if (tx.length == 0) {
    console.error('cannot decode the transaction', bytesToHex(tx))
    process.exit(1)
  }

  fundingProof.version = getVersion(tx)
  fundingProof.txInVector = getTxInputVector(tx)
  fundingProof.txOutVector = getTxOutputVector(tx)
  fundingProof.locktime = getLocktime(tx)
  fundingProof.fundingOutputIndex = getFundingOutputIndex(tx)// TODO: Find index in transaction based on deposit's public key

  const serializedProof = serializeProof(fundingProof)

  console.log('Funding Proof:', serializedProof)

  callback(fundingProof)
}

function serializeProof(fundingProof) {
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

function getPrefix(tx) {
  if (isFlagPresent(tx)) {
    return tx.slice(0, 6)
  }
  return tx.slice(0, 4)
}

function getVersion(tx) {
  return tx.slice(0, 4)
}

function isFlagPresent(tx) {
  if (tx.slice(4, 6).equals(Buffer.from('0001', 'hex'))) {
    return true
  }

  return false
}

function getTxInputVector(tx) {
  const txInVectorStartPosition = getPrefix(tx).length
  let txInVectorEndPosition

  if (isFlagPresent(tx)) {
    // TODO: Implement for witness transaction
    console.log('witness is not fully supported')
    txInVectorEndPosition = txInVectorStartPosition + 42
  } else {
    const inputCount = tx.slice(txInVectorStartPosition, txInVectorStartPosition + 1).readIntBE(0, 1)

    if (inputCount != 1) {
      // TODO: Support multiple inputs
      console.error(`exactly one input is required, got [${inputCount}]`)
      process.exit(1)
    } else {
      const startPos = txInVectorStartPosition + 1

      // const previousHash = tx.slice(startPos, startPos + 32).reverse()

      // const previousOutIndex = tx.slice(startPos + 32, startPos + 36).readIntBE(0, 4)

      const scriptLength = tx.slice(startPos + 36, startPos + 37).readIntBE(0, 1)
      if (scriptLength >= 253) {
        console.error(`VarInts not supported`)
        process.exit(1)
      }

      // const script = tx.slice(startPos + 37, startPos + 37 + scriptLength)

      // const sequenceNumber = tx.slice(startPos + 37 + scriptLength, startPos + 37 + scriptLength + 4)

      txInVectorEndPosition = startPos + 37 + scriptLength + 4
    }
  }

  return tx.slice(txInVectorStartPosition, txInVectorEndPosition)
}

function getTxOutputVector(tx) {
  const outStartPosition = getTxOutputVectorPosition(tx)
  const outputsCount = getNumberOfOutputs(tx)

  let startPosition = outStartPosition + 1
  let outEndPosition

  for (let i = 0; i < outputsCount; i++) {
    // const value = tx.slice(startPosition, startPosition + 8)

    const scriptLength = tx.slice(startPosition + 8, startPosition + 8 + 1).readIntBE(0, 1)

    if (scriptLength >= 253) {
      console.error(`VarInts not supported`)
      process.exit(1)
    }

    // const script = tx.slice(startPosition + 8 + 1, startPosition + 8 + 1 + scriptLength)

    outEndPosition = startPosition + 8 + 1 + scriptLength
    startPosition = outEndPosition
  }

  return tx.slice(outStartPosition, outEndPosition)
}

function getTxOutputVectorPosition(tx) {
  const txPrefix = getPrefix(tx)
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

function getLocktime(tx) {
  return tx.slice(tx.length - 4)
}

export function hexToBytes(hex) {
  return Buffer.from(hex, 'hex')
}

export function bytesToHex(bytes) {
  const buffer = Buffer.from(bytes)
  return buffer.toString('hex')
}
