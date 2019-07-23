// JS implementation of merkle.py script from [summa-tx/bitcoin-spv] repository.
//
// [summa-tx/bitcoin-spv]: https://github.com/summa-tx/bitcoin-spv/

const ElectrumClient = require('./ElectrumClient')
const ByteUtils = require('./utils/ByteUtils')
const Hash = require('./utils/Hash')
const Config = require('./utils/Config')

const config = Config.readFile(process.env.CONFIG_FILE)

module.exports = class BitcoinSPV {
  async initialize() {
    this.electrumClient = new ElectrumClient(
      config.electrum.server,
      config.electrum.port,
      config.electrum.protocol
    )

    await this.electrumClient.connect()
  }

  close() {
    this.electrumClient.close()
  }

  async getProof(txHash, confirmations) {
    // GET TRANSACTION
    const tx = await this.electrumClient.getTransaction(txHash)
      .catch((err) => {
        return Promise.reject(new Error(`failed to get transaction: [${err}]`))
      })

    if (tx.confirmations < confirmations) {
      return Promise.reject(new Error(`transaction confirmations number [${tx.confirmations}] is not enough, required [${confirmations}]`))
    }

    const latestBlockHeight = await this.electrumClient.latestBlockHeight()
      .catch((err) => {
        return Promise.reject(new Error(`failed to get latest block height: [${err}]`))
      })

    const txBlockHeight = latestBlockHeight - tx.confirmations + 1

    // GET HEADER CHAIN
    const headersChain = await this.electrumClient.getHeadersChain(txBlockHeight, confirmations)
      .catch((err) => {
        return Promise.reject(new Error(`failed to get headers chain: [${err}]`))
      })

    // GET MERKLE PROOF
    const merkleProof = await this.electrumClient.getMerkleProof(txHash, txBlockHeight)
      .catch((err) => {
        return Promise.reject(new Error(`failed to get merkle proof: [${err}]`))
      })

    // VERIFY PROOF
    if (!this.verifyProof(merkleProof.proof, merkleProof.position)) {
      return Promise.reject(new Error('invalid merkle proof'))
    }

    this.close()

    return Promise.resolve({
      tx: tx.hex,
      merkleProof: merkleProof.proof,
      txInBlockIndex: merkleProof.position,
      chainHeaders: headersChain,
    })
  }

  verifyProof(proofHex, index) {
    const proof = ByteUtils.fromHex(proofHex)

    const root = proof.slice(proof.length - 32, proof.length) // expectedRoot

    let currentHash = proof.slice(0, 32)

    // For all hashes between first and last
    for (let i = 1; i < (Math.floor(proof.length / 32) - 1); i++) {
      // If the current index is even,
      // The next hash goes before the current one
      if ((index % 2) == 0) {
        const children = Buffer.concat([proof.slice(i * 32, (i + 1) * 32), currentHash])
        currentHash = Hash.hash256(children)

        // Halve and floor the index
        index = Math.floor(index / 2)
      } else {
        // The next hash goes after the current one
        const children = Buffer.concat([currentHash, proof.slice(i * 32, (i + 1) * 32)])
        currentHash = Hash.hash256(children)

        // Halve and ceil the index
        index = Math.ceil(index / 2)
      }
    }

    // At the end we should have made the root
    if (currentHash.equals(root)) {
      return true
    } else {
      return false
    }
  }
}
