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

    assert.deepEqual(result, expectedResult)
  })

  it('parses legacy transaction details', async () => {
    const hex = '0100000001aea2e43133a533669b942d335ef6ebef7528f01c3ed1d43b4ccff1e9590d44c9010000006a4730440220785a31ce8bf2c63c5fbda079dea98f2740eaa81dfd09d6987b7ba9a4d2a5ccb702204ef4ff2f852a25fb4f75c2b16d61c09fcd411eea962f51e2ceec630de6e3cc8f0121028896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2effffffff022823000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4f0840d00000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac00000000'

    const result = await BitcoinTxParser.parse(hex)

    const expectedResult = {
      version: '01000000',
      txInVector: '01aea2e43133a533669b942d335ef6ebef7528f01c3ed1d43b4ccff1e9590d44c9010000006a4730440220785a31ce8bf2c63c5fbda079dea98f2740eaa81dfd09d6987b7ba9a4d2a5ccb702204ef4ff2f852a25fb4f75c2b16d61c09fcd411eea962f51e2ceec630de6e3cc8f0121028896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2effffffff',
      txOutVector: '022823000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4f0840d00000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac',
      locktime: '00000000',
    }

    assert.deepEqual(result, expectedResult)
  })
})
