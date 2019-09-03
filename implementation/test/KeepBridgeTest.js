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
    it('returns 0 when digest has not been approved', async () => {
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

    it('returns 0 when digest has been approved to sign by another keep', async () => {
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
})
