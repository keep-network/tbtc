const BitcoinSPV = require('../scripts/BitcoinSPV.js')
const fs = require('fs')
const chai = require('chai')
const assert = chai.assert

describe('BitcoinSPV', async () => {
  let tx
  let bitcoinSPV

  before(async () => {
    const txData = fs.readFileSync('./test/tx.json', 'utf8')
    tx = JSON.parse(txData)

    bitcoinSPV = new BitcoinSPV()
    await bitcoinSPV.initialize()
  })

  after(async () => {
    bitcoinSPV.close()
  })

  it('getProof', async () => {
    const expectedResult = {
      tx: tx.hex,
      merkleProof: tx.merkleProof,
      txInBlockIndex: tx.indexInBlock,
      chainHeaders: tx.chainHeaders,
    }

    const result = await bitcoinSPV.getProof(tx.hash, tx.chainHeadersNumber)

    assert.deepEqual(result, expectedResult)
  })

  it('verifyProof', async () => {
    const proofHex = tx.merkleProof
    const index = tx.indexInBlock
    const result = bitcoinSPV.verifyProof(proofHex, index)

    assert.isTrue(result)
  })
})
