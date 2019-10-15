const bcoin = require('bcoin')
const BN = require('bcrypto/lib/BN')
const Signature = require('bcrypto/lib/internal/signature')
const secp256k1 = require('bcrypto/lib/js/secp256k1')

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
  // TODO: When we want to give user a possibility to provide an address instead
  // of a public key hash we need to change it to `fromAddress`.
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

function bitcoinSignatureDER(r, s) {
  const size = secp256k1.size
  const signature = new Signature(size, r, s)

  // Check if `s` is a high value. As per BIP-0062 signature's `s` value should
  // be in a low half of curve's order. If it's a high value we convert it to `-s`.
  // Reference: https://en.bitcoin.it/wiki/BIP_0062#Low_S_values_in_signatures
  if (!secp256k1.isLowS(signature.encode(size))) {
    const newS = BN.fromBuffer(signature.s, 'be')

    newS.ineg().imod(secp256k1.curve.n)

    signature.s = secp256k1.curve.encodeScalar(newS)
    signature.param ^= 1
  }

  return signature.toDER(size)
}

export function addWitnessSignature(unsignedTransaction, inputIndex, r, s, publicKey) {
  // Signature
  const signatureDER = bitcoinSignatureDER(r, s)

  const hashType = Buffer.from([bcoin.Script.hashType.ALL])
  const sig = Buffer.concat([signatureDER, hashType])

  // Public Key
  const compressedPublicKey = secp256k1.publicKeyImport(publicKey, true)

  // Combine witness
  const signedTransaction = bcoin.TX.fromRaw(unsignedTransaction, 'hex').clone()
  signedTransaction.inputs[inputIndex].witness.fromItems([sig, compressedPublicKey])

  return signedTransaction.toRaw().toString('hex')
}
