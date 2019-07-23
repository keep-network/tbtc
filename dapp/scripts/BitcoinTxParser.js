const ByteUtils = require('./utils/ByteUtils.js')

async function parse(tx) {
  console.log('Parse transaction...\nTX:', tx)

  tx = ByteUtils.fromHex(tx)
  if (tx.length == 0) {
    return Promise.reject(new Error('cannot decode the transaction'))
  }

  const txDetails = {
    version: ByteUtils.toHex(getTxVersion(tx)),
    txInVector: ByteUtils.toHex(getTxInputVector(tx)),
    txOutVector: ByteUtils.toHex(getTxOutputVector(tx)),
    locktime: ByteUtils.toHex(getTxLocktime(tx)),
  }

  return txDetails
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

  const outputsCount = tx.slice(outStartPosition, outStartPosition + 1).readIntBE(0, 1)

  if (outputsCount >= 253) {
    throw new Error('VarInts not supported')
  }

  return outputsCount
}

function getTxLocktime(tx) {
  return tx.slice(tx.length - 4)
}

module.exports.parse = parse
