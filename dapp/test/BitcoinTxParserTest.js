const BitcoinTxParser = require('../scripts/BitcoinTxParser.js')

const chai = require('chai')
const assert = chai.assert
const fs = require('fs')


describe('BitcoinTxParser', async () => {
  let tx

  before(async () => {
    const txData = fs.readFileSync('./test/tx.json', 'utf8')
    tx = JSON.parse(txData)
  })

  it('parses witness transaction details', async () => {
    const result = await BitcoinTxParser.parse(tx.hex)

    const expectedResult = {
      version: tx.version,
      txInVector: tx.txInVector,
      txOutVector: tx.txOutVector,
      locktime: tx.locktime,
    }

    assert.deepEqual(
      result,
      expectedResult,
    )
  })
})
