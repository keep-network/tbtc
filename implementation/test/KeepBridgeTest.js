const KeepBridge = artifacts.require('KeepBridge')
const KeepRegistryStub = artifacts.require('KeepRegistryStub')
const ECDSAKeepVendorStub = artifacts.require('ECDSAKeepVendorStub')
const ECDSAKeepStub = artifacts.require('ECDSAKeepStub')

contract('KeepBridge', (accounts) => {
  let keepBridge
  let ecdsaKeepStub
  let ecdsaKeepVendor

  before(async () => {
    ecdsaKeepVendor = await ECDSAKeepVendorStub.new()

    const keepRegistry = await KeepRegistryStub.new()
    await keepRegistry.setVendor(ecdsaKeepVendor.address)

    keepBridge = await KeepBridge.deployed()
    await keepBridge.initialize(keepRegistry.address)

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

  describe('approveDigest()', async () => {
    it('calls ECDSA keep for signing', async () => {
      const digest = '0x' + '08'.repeat(32)

      await keepBridge.approveDigest(ecdsaKeepStub.address, digest)
        .catch((err) => {
          assert.fail(`cannot approve digest: ${err}`)
        })

      // Check if ECDSAKeep has been called and event emitted.
      const eventList = await ecdsaKeepStub.getPastEvents(
        'SignatureRequested',
        { fromBlock: 0, toBlock: 'latest' },
      )

      assert.equal(
        eventList[0].returnValues.digest,
        digest,
        'incorrect digest in emitted event',
      )
    })
  })

  describe('wasDigestApprovedForSigning()', async () => {
    it('returns 0 when digest has not been registered', async () => {
      const digest = '0x' + '01'.repeat(32)

      const result = await keepBridge.wasDigestApprovedForSigning(ecdsaKeepStub.address, digest)
        .catch((err) => {
          assert.fail(`cannot check digest approval: ${err}`)
        })

      assert.equal(
        result,
        0,
        'incorrect result',
      )
    })

    it('returns timestamp registered on digest approval', async () => {
      const digest = '0x' + '02'.repeat(32)

      const approvalTx = await keepBridge.approveDigest(ecdsaKeepStub.address, digest)
        .catch((err) => {
          assert.fail(`cannot approve digest: ${err}`)
        })

      const block = await web3.eth.getBlock(approvalTx.receipt.blockNumber)
      const expectedTimestamp = block.timestamp

      const timestamp = await keepBridge.wasDigestApprovedForSigning(ecdsaKeepStub.address, digest)
        .catch((err) => {
          assert.fail(`cannot check digest approval: ${err}`)
        })

      assert.equal(
        timestamp,
        expectedTimestamp,
        'incorrect registered timestamp',
      )
    })

    it('returns 0 when digest has been registered for another keep', async () => {
      const digest = '0x' + '03'.repeat(32)
      const keep2address = '0x' + '04'.repeat(20)

      await keepBridge.approveDigest(ecdsaKeepStub.address, digest)
        .catch((err) => {
          assert.fail(`cannot approve digest: ${err}`)
        })

      const timestamp = await keepBridge.wasDigestApprovedForSigning(keep2address, digest)
        .catch((err) => {
          assert.fail(`cannot check digest approval: ${err}`)
        })

      assert.equal(
        timestamp,
        0,
        'incorrect registered timestamp',
      )
    })
  })

  describe('requestNewKeep()', async () => {
    it('sends caller as owner to open new keep', async () => {
      const expectedKeepOwner = accounts[2]

      await keepBridge.requestNewKeep(5, 10, { from: expectedKeepOwner })
      const keepOwner = await ecdsaKeepVendor.keepOwner.call()

      assert.equal(expectedKeepOwner, keepOwner, 'incorrect keep owner address')
    })
  })
})
