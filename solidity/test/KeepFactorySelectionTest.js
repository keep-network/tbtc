const {contract} = require("@openzeppelin/test-environment")
const {expectRevert, constants} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {expect} = require("chai")

const KeepFactorySelection = contract.fromArtifact("KeepFactorySelection")
const KeepFactorySelectionStub = contract.fromArtifact(
  "KeepFactorySelectionStub",
)
const KeepFactorySelectorStub = contract.fromArtifact("KeepFactorySelectorStub")
const ECDSAKeepFactoryStub = contract.fromArtifact("ECDSAKeepFactoryStub")

describe("KeepFactorySelection", async () => {
  let keepFactorySelection

  let keepStakeFactory
  let ethStakeFactory
  let keepFactorySelector

  before(async () => {
    await KeepFactorySelection.detectNetwork()
    const library = await KeepFactorySelection.new()

    await KeepFactorySelectionStub.detectNetwork()
    await KeepFactorySelectionStub.link("KeepFactorySelection", library.address)
    keepFactorySelection = await KeepFactorySelectionStub.new()

    keepStakeFactory = await ECDSAKeepFactoryStub.new()
    ethStakeFactory = await ECDSAKeepFactoryStub.new()
    keepFactorySelector = await KeepFactorySelectorStub.new()

    await keepFactorySelection.initialize(keepStakeFactory.address)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("initialize", async () => {
    it("can be called only one time", async () => {
      // already initialized in before
      await expectRevert(
        keepFactorySelection.initialize(keepStakeFactory.address),
        "Already initialized",
      )
    })
  })

  describe("selectFactory", async () => {
    // No ETH stake factory set,
    // No selection strategy set.
    it("returns KEEP stake factory by default", async () => {
      const selected = await keepFactorySelection.selectFactory()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // No ETH stake factory set,
    // No selection strategy set.
    it("returns selected KEEP-staked factory even after version update", async () => {
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setKeepStakedKeepFactory(
        newKeepFactory.address,
      )

      expect(
        await keepFactorySelection.selectFactory(),
        "unexpected factory selected",
      ).to.equal(keepStakeFactory.address)
    })

    // ETH stake factory set.
    // Selection strategy set.
    it("returns selected ETH-only factory even after version update", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // update ETH-only factory version
      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFullyBackedKeepFactory(
        newEthFactory.address,
      )

      expect(
        await keepFactorySelection.selectFactory(),
        "unexpected factory selected",
      ).to.equal(ethStakeFactory.address)
    })

    // ETH stake factory set.
    // Selection strategy set.
    it("returns the same factory until refreshed", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setDefaultMode()

      let selected = await keepFactorySelection.selectFactory()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // although the strategy wants to return fully-backed factory
      // selectFactory() is still returning the default option, until
      // the selection will be refreshed
      selected = await keepFactorySelection.selectFactory()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // the factory has changed because selectFactoryAndRefresh
      // refreshed the selection
      selected = await keepFactorySelection.selectFactory()
      expect(selected, "unexpected factory selected").to.equal(
        ethStakeFactory.address,
      )
    })
  })

  describe("selectFactoryAndRefresh", async () => {
    // No ETH stake factory set.
    // No selection strategy set.
    it("returns KEEP stake factory by default", async () => {
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH stake factory set.
    // No selection strategy set.
    it("returns KEEP stake factory if strategy is not set", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH stake factory set.
    // Selection strategy set.
    it("returns fully-backed factory if selected by the strategy", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // returns the previously selected factory and performs a new
      // selection; we don't care about the previous result here
      await keepFactorySelection.selectFactoryAndRefresh()

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        ethStakeFactory.address,
      )
    })

    // ETH stake factory set.
    // Selection strategy set.
    it("returns default factory if selected by the strategy", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setDefaultMode()

      // returns the previously selected factory and performs a new
      // selection; we don't care about the previous result here
      await keepFactorySelection.selectFactoryAndRefresh()

      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH stake factory set.
    // Selection strategy set.
    it("returns the factory selected before refreshing", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
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
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    it("reverts if the returned factory is not one of the set factories", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setMaliciousMode()

      // refresh the choice; it should be the default factory now
      await expectRevert(
        keepFactorySelection.selectFactoryAndRefresh(),
        "Factory selector returned unknown factory",
      )
    })

    // No ETH-only vendor set.
    // No selection strategy set.
    it("gets new KEEP stake factory after update", async () => {
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setKeepStakedKeepFactory(
        newKeepFactory.address,
      )

      // refresh the choice; it should be the old factory now
      const selected1 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected1, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )

      // refresh the choice; it should be the new factory now
      const selected2 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected2, "unexpected factory selected").to.equal(
        newKeepFactory.address,
      )
    })

    // ETH-only vendor set.
    // Selection strategy set.
    it("gets new ETH-only factory after update", async () => {
      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice; it should be the default factory now
      const selected1 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected1, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )

      // refresh the choice; it should be the old fully backed factory now
      const selected2 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected2, "unexpected factory selected").to.equal(
        ethStakeFactory.address,
      )

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await keepFactorySelection.setFullyBackedKeepFactory(
        newEthFactory.address,
      )

      // refresh the choice; it should still be the old fully backed factory now
      const selected3 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected3, "unexpected factory selected").to.equal(
        ethStakeFactory.address,
      )

      // refresh the choice; it should be the new fully backed factory now
      const selected4 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected4, "unexpected factory selected").to.equal(
        newEthFactory.address,
      )
    })
  })

  describe("setKeepStakedKeepFactory", async () => {
    it("can be called multiple times", async () => {
      const newKeepStakeFactory = await ECDSAKeepFactoryStub.new()

      await keepFactorySelection.setKeepStakedKeepFactory(
        ethStakeFactory.address,
      )

      expect(await keepFactorySelection.keepStakedFactory()).to.equal(
        ethStakeFactory.address,
      )

      await keepFactorySelection.setKeepStakedKeepFactory(
        newKeepStakeFactory.address,
      )

      expect(await keepFactorySelection.keepStakedFactory()).to.equal(
        newKeepStakeFactory.address,
      )
    })

    it("can not be called for 0 address", async () => {
      await expectRevert(
        keepFactorySelection.setKeepStakedKeepFactory(constants.ZERO_ADDRESS),
        "Invalid address",
      )
    })
  })

  describe("setFullyBackedKeepFactory", async () => {
    it("can be called multiple times", async () => {
      const newFullyBackedFactory = await ECDSAKeepFactoryStub.new()

      await keepFactorySelection.setFullyBackedKeepFactory(
        ethStakeFactory.address,
      )

      expect(await keepFactorySelection.fullyBackedFactory()).to.equal(
        ethStakeFactory.address,
      )

      await keepFactorySelection.setFullyBackedKeepFactory(
        newFullyBackedFactory.address,
      )

      expect(await keepFactorySelection.fullyBackedFactory()).to.equal(
        newFullyBackedFactory.address,
      )
    })

    it("can not be called for 0 address", async () => {
      await expectRevert(
        keepFactorySelection.setFullyBackedKeepFactory(constants.ZERO_ADDRESS),
        "Invalid address",
      )
    })
  })

  describe("setKeepFactorySelector", async () => {
    it("can be called multiple times", async () => {
      const newSelector = await KeepFactorySelectorStub.new()

      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )
      expect(await keepFactorySelection.factorySelector()).to.equal(
        keepFactorySelector.address,
      )

      await keepFactorySelection.setKeepFactorySelector(newSelector.address)

      expect(await keepFactorySelection.factorySelector()).to.equal(
        newSelector.address,
      )
    })

    it("can not be called for 0 address", async () => {
      await expectRevert(
        keepFactorySelection.setKeepFactorySelector(constants.ZERO_ADDRESS),
        "Invalid address",
      )
    })
  })
})
