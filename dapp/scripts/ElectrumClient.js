const ElectrumCli = require('electrum-client')
const ByteUtils = require('./utils/ByteUtils.js')

module.exports = class ElectrumClient {
  constructor(server, port, protocol) {
    this.client = new ElectrumCli(port, server, protocol)
  }

  async connect() {
    return new Promise((resolve, reject) => {
      this.client.connect()
        .then(() => {
          this.client.server_banner()
            .then((banner) => {
              console.log(banner)
            })
          resolve()
        }
        )
        .catch((err) => {
          reject(err)
        })
    })
  }

  close() {
    this.client.close()
  }

  latestBlockHeight() {
    return new Promise((resolve, reject) => {
      this.client.blockchainHeaders_subscribe()
        .then((header) => {
          resolve(header.height)
        })
        .catch((err) => {
          reject(new Error(JSON.stringify(err)))
        })
    })
  }

  getTransaction(txHash) {
    return new Promise((resolve, reject) => {
      this.client.blockchainTransaction_get(txHash, true)
        .then((tx) => {
          resolve(tx)
        })
        .catch((err) => {
          reject(new Error(JSON.stringify(err)))
        })
    })
  }

  getMerkleRoot(blockHeight) {
    return new Promise((resolve, reject) => {
      this.client.blockchainBlock_header(blockHeight)
        .then((header) => {
          const merkleRoot = ByteUtils.fromHex(header).slice(36, 68) // ??? WHY JAMES SLICE IT ????
          resolve(merkleRoot)
        })
        .catch((err) => {
          reject(new Error(JSON.stringify(err)))
        })
    })
  }

  getHeadersChain(blockHeight, confirmations) {
    return new Promise((resolve, reject) => {
      this.client.blockchainBlock_headers(blockHeight, confirmations + 1)
        .then((headersChain) => {
          resolve(headersChain.hex)
        })
        .catch((err) => {
          reject(new Error(JSON.stringify(err)))
        })
    })
  }

  getMerkleProof(txHash, blockHeight) {
    return new Promise((resolve, reject) => {
      this.client.blockchainTransaction_getMerkle(txHash, blockHeight)
        .then((merkle) => {
          const position = merkle.pos + 1 // add 1 because proof uses 1-indexed positions

          let proof = ByteUtils.fromHex(txHash).reverse()

          merkle.merkle.forEach(function(item) {
            proof = Buffer.concat([proof, ByteUtils.fromHex(item).reverse()])
          })

          // GET MERKLE ROOT
          this.getMerkleRoot(blockHeight)
            .then((merkleRoot) => {
              proof = Buffer.concat([proof, ByteUtils.fromHex(merkleRoot)])
              resolve({ 'proof': ByteUtils.toHex(proof), 'position': position })
            })
            .catch((err) => {
              reject(new Error(JSON.stringify(err)))
            })
        })
    })
  }

  findOutputForAddress(txHash, address) {
    return new Promise((resolve, reject) => {
      this.getTransaction(txHash)
        .then((tx) => {
          const outputs = tx.vout
          outputs.forEach((output, index) => {
            output.scriptPubKey.addresses.forEach((a) => {
              if (a == address) {
                resolve(index)
              }
            })
          })
          reject(new Error(`output for address ${address} not found`))
        })
        .catch((err) => {
          reject(new Error(JSON.stringify(err)))
        })
    })
  }

  watchHeaders() {
    // TODO: Leaving this snippet for later to find the example easily
    //
    // client.subscribe.on('blockchain.headers.subscribe', (v) => {
    //   console.log('emitted', v)
    // }) // subscribe message(EventEmitter)
  }
}
