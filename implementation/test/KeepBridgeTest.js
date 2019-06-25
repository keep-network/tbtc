const KeepBridge = artifacts.require('KeepBridge')
const ECDSAKeepStub = artifacts.require('ECDSAKeepStub')

contract('KeepBridge', (accounts) => {
  let keepBridge
  let ecdsaKeepStub

  before(async () => {
    keepBridge = await KeepBridge.deployed()
    ecdsaKeepStub = await ECDSAKeepStub.new()
  })

  describe('getKeepPubkey()', async () => {
    it('returns null when public key is not set', async () => {
      const publicKey = await keepBridge.getKeepPubkey.call(ecdsaKeepStub.address)
        .catch((err) => {
          assert.fail(`cannot get public key: ${err}`)
        })

      assert.equal(
        publicKey,
        null,
        'incorrect public key'
      )
    })

    it('returns the public key when it is set', async () => {
      const expectedPublicKey = web3.utils.hexToBytes('0xcf49e51388d87d9d878e4382880d4cbc20daf3865f499b182649755ad75fd81300b68eb03383826fc5bf67b489ab65224efbea76b81074ee52839986112e9e5e')

      await ecdsaKeepStub.setPublicKey(expectedPublicKey)
        .catch((err) => {
          assert.fail(`cannot set public key for keep: ${err}`)
        })

      const publicKey = await keepBridge.getKeepPubkey.call(ecdsaKeepStub.address)
        .catch((err) => {
          assert.fail(`cannot get public key: ${err}`)
        })

      assert.equal(
        publicKey,
        web3.utils.bytesToHex(expectedPublicKey),
        'incorrect public key'
      )
    })
  })
})
