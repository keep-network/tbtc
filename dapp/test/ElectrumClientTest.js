const ElectrumClient = require('../scripts/ElectrumClient')
const Config = require('../scripts/utils/Config')

const fs = require('fs')
const chai = require('chai')
const assert = chai.assert

const config = Config.readFile(process.env.CONFIG_FILE)

describe('ElectrumClient', async () => {
  const client = new ElectrumClient(
    config.electrum.server,
    config.electrum.port,
    config.electrum.protocol
  )
  let tx

  before(async () => {
    const txData = fs.readFileSync('./test/tx.json', 'utf8')
    tx = JSON.parse(txData)

    await client.connect()
  })

  after(async () => {
    client.close()
  })

  describe('findOutputForAddress', async () => {
    let result
    let expectedResult

    it('getTransaction', async () => {
      const expectedTx = tx.hex
      const result = await client.getTransaction(tx.hash)

      assert.equal(
        result.hex,
        expectedTx,
        'unexpected result',
      )
    })

    it('getMerkleProof', async () => {
      const expectedResult = tx.merkleProof
      const expectedPosition = tx.indexInBlock
      const result = await client.getMerkleProof(tx.hash, tx.blockHeight)

      assert.equal(
        result.proof,
        expectedResult,
        'unexpected result',
      )

      assert.equal(
        result.position,
        expectedPosition,
        'unexpected result',
      )
    })

    it('getHeadersChain', async () => {
      const confirmations = tx.chainHeadersNumber
      const expectedResult = tx.chainHeaders
      const result = await client.getHeadersChain(tx.blockHeight, confirmations)

      assert.equal(
        result,
        expectedResult,
        'unexpected result',
      )
    })

    describe('findOutputForAddress', async () => {
      afterEach(() => {
        assert.equal(
          result,
          expectedResult
        )
      })

      it('finds first element', async () => {
        const address = 'tb1qfdru0xx39mw30ha5a2vw23reymmxgucujfnc7l'
        expectedResult = 0

        result = await client.findOutputForAddress(tx.hash, address)
      })

      it('finds second element', async () => {
        const address = 'tb1q78ezl08lyhuazzfz592sstenmegdns7durc4cl'
        expectedResult = 1

        result = await client.findOutputForAddress(tx.hash, address)
      })
    })
  })
})
