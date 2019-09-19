const TBTCSystem = artifacts.require('TBTCSystem')
const KeepRegistryStub = artifacts.require('KeepRegistryStub')
const ECDSAKeepVendorStub = artifacts.require('ECDSAKeepVendorStub')

contract('TBTCSystem', (accounts) => {
  let tbtcSystem
  let ecdsaKeepVendor

  before(async () => {
    ecdsaKeepVendor = await ECDSAKeepVendorStub.new()

    const keepRegistry = await KeepRegistryStub.new()
    await keepRegistry.setVendor(ecdsaKeepVendor.address)

    tbtcSystem = await TBTCSystem.deployed()
    await tbtcSystem.initialize(
      keepRegistry.address,
      '0x0000000000000000000000000000000000000000' // TBTC Uniswap Exchange
    )
  })

  describe('requestNewKeep()', async () => {
    it('sends caller as owner to open new keep', async () => {
      const expectedKeepOwner = accounts[2]

      await tbtcSystem.requestNewKeep(5, 10, { from: expectedKeepOwner })
      const keepOwner = await ecdsaKeepVendor.keepOwner.call()

      assert.equal(expectedKeepOwner, keepOwner, 'incorrect keep owner address')
    })

    it('returns keep address', async () => {
      const expectedKeepAddress = await ecdsaKeepVendor.keepAddress.call()

      const result = await tbtcSystem.requestNewKeep.call(5, 10)

      assert.equal(expectedKeepAddress, result, 'incorrect keep address')
    })
  })
})
