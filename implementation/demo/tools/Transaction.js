const bcoin = require('bcoin')
const secp256k1 = require('bcrypto').secp256k1
const Signature = require('bcrypto/lib/internal/signature')

export function oneInputOneOutputWitnessTX(
  inputPreviousOutpoint, // 36 byte UTXO id
  inputSequence,
  outputValue, // BN
  outputPKH,
) {
  // Input
  const prevOutpoint = bcoin.Outpoint.fromRaw(inputPreviousOutpoint)

  const input = bcoin.Input.fromOptions({
    prevout: prevOutpoint,
    sequence: inputSequence,
  })

  // Output
  const outputScript = bcoin.Script.fromProgram(
    0, // Witness program version
    outputPKH
  )

  const output = bcoin.Output.fromOptions({
    value: outputValue.toNumber(),
    script: outputScript,
  })

  // Transaction
  const preimageTX = bcoin.TX.fromOptions({
    inputs: [input],
    outputs: [output],
  })

  return preimageTX.toRaw().toString('hex')
}

export function addWitnessSignature(unsignedTransaction, inputIndex, r, s, publicKey) {
  const signedTransaction = bcoin.TX.fromRaw(unsignedTransaction, 'hex').clone()

  // Signature
  const size = 32
  const signature = new Signature(size, r, s)
  const signatureDER = signature.toDER(size)

  const hashType = Buffer.from([bcoin.Script.hashType.ALL])
  const sig = Buffer.concat([signatureDER, hashType])

  // Public Key
  const compressedPublicKey = secp256k1.publicKeyImport(publicKey, true)


  signedTransaction.inputs[inputIndex].witness.fromItems([sig, compressedPublicKey])

  return signedTransaction.toRaw().toString('hex')
}
