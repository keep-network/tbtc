const { contract } = require('@openzeppelin/test-environment')
const { BN, expectRevert, constants } = require('@openzeppelin/test-helpers')
const { createSnapshot, restoreSnapshot } = require('./helpers/snapshot.js')
const { expect } = require('chai')

const KeepFactorySelection = contract.fromArtifact('KeepFactorySelection')
const KeepFactorySelectionStub = contract.fromArtifact(
  'KeepFactorySelectionStub',
)
const KeepFactorySelectorStub = contract.fromArtifact('KeepFactorySelectorStub')
const ECDSAKeepFactoryStub = contract.fromArtifact('ECDSAKeepFactoryStub')

describe('KeepFactorySelection', async () => {
  let keepFactorySelection

  let keepStakedFactory
  let fullyBackedFactory
  let keepFactorySelector

  before(async () => {
    await KeepFactorySelection.detectNetwork()
    const library = await KeepFactorySelection.new()

    await KeepFactorySelectionStub.detectNetwork()
    await KeepFactorySelectionStub.link('KeepFactorySelection', library.address)
    keepFactorySelection = await KeepFactorySelectionStub.new()

    keepStakedFactory = await ECDSAKeepFactoryStub.new()
    fullyBackedFactory = await ECDSAKeepFactoryStub.new()
    keepFactorySelector = await KeepFactorySelectorStub.new()

    await keepFactorySelection.initialize(keepStakedFactory.address)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe('initialize', async () => {
    it('can be called only one time', async () => {
      // already initialized in before
      await expectRevert(
        keepFactorySelection.initialize(keepStakedFactory.address),
        'Already initialized',
      )
    })
  })

  describe('selectFactory', async () => {
    // No ETH-bond-only factory set,
    // No selection strategy set.
    it('returns KEEP stake factory by default', async () => {
      const selected = await keepFactorySelection.selectFactory()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    // No ETH-bond-only factory set,
    // No selection strategy set.
    it('returns selected KEEP-staked factory even after version update', async () => {
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFactories(
        newKeepFactory.address,
        constants.ZERO_ADDRESS,
        constants.ZERO_ADDRESS,
      )

      expect(
        await keepFactorySelection.selectFactory(),
        'unexpected factory selected',
      ).to.equal(keepStakedFactory.address)
    })

    // ETH-bond-only factory set.
    // Selection strategy set.
    it('returns selected ETH-only factory even after version update', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // update ETH-only factory version
      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        newEthFactory.address,
        keepFactorySelector.address,
      )

      expect(
        await keepFactorySelection.selectFactory(),
        'unexpected factory selected',
      ).to.equal(fullyBackedFactory.address)
    })

    // ETH-bond-only factory set.
    // Selection strategy set.
    it('returns the same factory until refreshed', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setDefaultMode()

      let selected = await keepFactorySelection.selectFactory()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // although the strategy wants to return fully-backed factory
      // selectFactory() is still returning the default option, until
      // the selection will be refreshed
      selected = await keepFactorySelection.selectFactory()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // the factory has changed because selectFactoryAndRefresh
      // refreshed the selection
      selected = await keepFactorySelection.selectFactory()
      expect(selected, 'unexpected factory selected').to.equal(
        fullyBackedFactory.address,
      )
    })
  })

  describe('selectFactoryAndRefresh', async () => {
    // No ETH-bond-only factory set.
    // No selection strategy set.
    it('returns KEEP stake factory by default', async () => {
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    // ETH-bond-only factory set.
    // No selection strategy set.
    it('returns KEEP stake factory if strategy is not set', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        fullyBackedFactory.address,
        constants.ZERO_ADDRESS,
      )

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    // No ETH-only vendor set.
    // Selection strategy set.
    it('returns KEEP stake factory if ETH-only factory is not set', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        constants.ZERO_ADDRESS,
        keepFactorySelector.address,
      )

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    // ETH-bond-only factory set.
    // Selection strategy set.
    it('returns fully-backed factory if selected by the strategy', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // returns the previously selected factory and performs a new
      // selection; we don't care about the previous result here
      await keepFactorySelection.selectFactoryAndRefresh()

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        fullyBackedFactory.address,
      )
    })

    // ETH-bond-only factory set.
    // Selection strategy set.
    it('returns default factory if selected by the strategy', async () => {
      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setDefaultMode()

      // returns the previously selected factory and performs a new
      // selection; we don't care about the previous result here
      await keepFactorySelection.selectFactoryAndRefresh()

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    // ETH-bond-only factory set.
    // Selection strategy set.
    it('returns the factory selected before refreshing', async () => {
      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setDefaultMode()

      // refresh the choice; it should be the default factory now
      await keepFactorySelection.selectFactoryAndRefresh()

      await keepFactorySelector.setFullyBackedMode()

      // although the strategy selects fully-backed factory, selectFactoryAndRefresh
      // should return the factory based on the previous choice (before it
      // refreshed)
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )
    })

    it('reverts if the returned factory is not one of the set factories', async () => {
      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setMaliciousMode()

      // refresh the choice; it should be the default factory now
      await expectRevert(
        keepFactorySelection.selectFactoryAndRefresh(),
        'Factory selector returned unknown factory',
      )
    })

    // No ETH-only vendor set.
    // No selection strategy set.
    it('gets new KEEP stake factory after update', async () => {
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFactories(
        newKeepFactory.address,
        constants.ZERO_ADDRESS,
        constants.ZERO_ADDRESS,
      )

      // refresh the choice; it should be the old factory now
      const selected1 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected1, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )

      // refresh the choice; it should be the new factory now
      const selected2 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected2, 'unexpected factory selected').to.equal(
        newKeepFactory.address,
      )
    })

    // ETH-only vendor set.
    // Selection strategy set.
    it('gets new ETH-only factory after update', async () => {
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice; it should be the default factory now
      const selected1 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected1, 'unexpected factory selected').to.equal(
        keepStakedFactory.address,
      )

      // refresh the choice; it should be the old fully backed factory now
      const selected2 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected2, 'unexpected factory selected').to.equal(
        fullyBackedFactory.address,
      )

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFactories(
        keepFactorySelector.address,
        newEthFactory.address,
        keepFactorySelector.address,
      )

      // refresh the choice; it should still be the old fully backed factory now
      const selected3 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected3, 'unexpected factory selected').to.equal(
        fullyBackedFactory.address,
      )

      // refresh the choice; it should be the new fully backed factory now
      const selected4 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected4, 'unexpected factory selected').to.equal(
        newEthFactory.address,
      )
    })
  })

  describe('setMinimumBondableValue', async () => {
    let defaultValue
    let newValue

    before(async () => {
      defaultValue = new BN(999)
      newValue = new BN(123987)
    })

    // No KEEP stake factory set.
    it('completes when KEEP stake factory is not set', async () => {
      const newKeepFactorySelection = await KeepFactorySelectionStub.new()

      await newKeepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakedFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
    })

    // KEEP stake factory set.
    // ETH-only factory not set.
    it('updates value in KEEP stake factory', async () => {
      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakedFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
    })

    // KEEP stake factory set.
    // ETH-only factory set.
    it('updates values in KEEP stake and ETH-only factories', async () => {
      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        fullyBackedFactory.address,
        constants.ZERO_ADDRESS,
      )

      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakedFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(newValue)
    })
  })

  describe('setFactories', async () => {
    it('can be called multiple times', async () => {
      const newKeepStakedFactory = await ECDSAKeepFactoryStub.new()
      const newFullyBackedFactory = await ECDSAKeepFactoryStub.new()
      const newSelector = await KeepFactorySelectorStub.new()

      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        fullyBackedFactory.address,
        keepFactorySelector.address,
      )

      expect(
        await keepFactorySelection.factories(),
        'invalid factories after first update',
      ).to.include({
        _keepStakedFactory: keepStakedFactory.address,
        _fullyBackedFactory: fullyBackedFactory.address,
        _factorySelector: keepFactorySelector.address,
      })

      await keepFactorySelection.setFactories(
        newKeepStakedFactory.address,
        newFullyBackedFactory.address,
        newSelector.address,
      )

      expect(
        await keepFactorySelection.factories(),
        'invalid factories after second update',
      ).to.include({
        _keepStakedFactory: newKeepStakedFactory.address,
        _fullyBackedFactory: newFullyBackedFactory.address,
        _factorySelector: newSelector.address,
      })
    })

    it('can be called for fully backed and factory selector zero addresses', async () => {
      await keepFactorySelection.setFactories(
        keepStakedFactory.address,
        constants.ZERO_ADDRESS,
        constants.ZERO_ADDRESS,
      )

      expect(
        await keepFactorySelection.factories(),
        'invalid factories',
      ).to.include({
        _keepStakedFactory: keepStakedFactory.address,
        _fullyBackedFactory: constants.ZERO_ADDRESS,
        _factorySelector: constants.ZERO_ADDRESS,
      })
    })

    it('reverts when KEEP-staked factory address is zero', async () => {
      await expectRevert(
        keepFactorySelection.setFactories(
          constants.ZERO_ADDRESS,
          fullyBackedFactory.address,
          keepFactorySelector.address,
        ),
        'Invalid KEEP-staked factory address',
      )
    })
  })
})
