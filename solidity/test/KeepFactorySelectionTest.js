const {accounts, contract} = require("@openzeppelin/test-environment")
const {BN, expectRevert, constants} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {expect} = require("chai")

const KeepFactorySelection = contract.fromArtifact("KeepFactorySelection")
const KeepFactorySelectionStub = contract.fromArtifact(
  "KeepFactorySelectionStub",
)
const KeepFactorySelectorStub = contract.fromArtifact("KeepFactorySelectorStub")
const ECDSAKeepFactoryStub = contract.fromArtifact("ECDSAKeepFactoryStub")
const ECDSAKeepVendorStub = contract.fromArtifact("ECDSAKeepVendorStub")

describe("KeepFactorySelection", async () => {
  let keepFactorySelection

  let keepStakeFactory
  let fullyBackedFactory
  let keepFactorySelector

  const thirdParty = accounts[3]

  before(async () => {
    await KeepFactorySelection.detectNetwork()
    const library = await KeepFactorySelection.new()

    await KeepFactorySelectionStub.detectNetwork()
    await KeepFactorySelectionStub.link("KeepFactorySelection", library.address)
    keepFactorySelection = await KeepFactorySelectionStub.new()

    keepStakeFactory = await ECDSAKeepFactoryStub.new()
    keepStakeVendor = await ECDSAKeepVendorStub.new(keepStakeFactory.address)

    fullyBackedFactory = await ECDSAKeepFactoryStub.new()
    fullyBackedVendor = await ECDSAKeepVendorStub.new(
      fullyBackedFactory.address,
    )

    keepFactorySelector = await KeepFactorySelectorStub.new()

    await keepFactorySelection.initialize(keepStakeVendor.address)
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
        keepFactorySelection.initialize(keepStakeVendor.address),
        "Already initialized",
      )
    })

    it("reverts when vendor returns factory with zero address", async () => {
      const keepFactorySelection = await KeepFactorySelectionStub.new()

      await keepStakeVendor.setFactory(constants.ZERO_ADDRESS)

      await expectRevert(
        keepFactorySelection.initialize(keepStakeVendor.address),
        "Vendor returned invalid factory address",
      )
    })
  })

  describe("selectFactory", async () => {
    // No ETH-only factory set,
    // No selection strategy set.
    it("returns KEEP stake factory by default", async () => {
      const selected = await keepFactorySelection.selectFactory()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH-only factory set.
    // Selection strategy set.
    it("returns the same factory until refreshed", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
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
        fullyBackedFactory.address,
      )
    })
  })

  describe("selectFactoryAndRefresh", async () => {
    // No ETH-only vendor set.
    // No selection strategy set.
    it("returns KEEP stake factory by default", async () => {
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH-only vendor set.
    // No selection strategy set.
    it("returns KEEP stake factory if strategy is not set", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // No ETH-only vendor set.
    // Selection strategy set.
    it("returns KEEP stake factory if ETH-only factory is not set", async () => {
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH-only vendor set.
    // Selection strategy set.
    it("returns fully-backed factory if selected by the strategy", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
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
        fullyBackedFactory.address,
      )
    })

    // ETH-only vendor set.
    // Selection strategy set.
    it("returns default factory if selected by the strategy", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
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

    // ETH-only vendor set.
    // Selection strategy set.
    it("returns the factory selected before refreshing", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
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

    // No ETH-only vendor set.
    // No selection strategy set.
    it("gets new KEEP stake factory after update in vendor", async () => {
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepStakeVendor.setFactory(newKeepFactory.address)

      // refresh the choice; it should be the old factory now
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
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
    it("gets new ETH-only factory after update in vendor", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice; it should be the default factory now
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the old factory now
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        fullyBackedFactory.address,
      )

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await fullyBackedVendor.setFactory(newEthFactory.address)

      // refresh the choice; it should be the old factory now
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the new factory now
      const selected2 = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected2, "unexpected factory selected").to.equal(
        newEthFactory.address,
      )
    })

    // No ETH-only vendor set.
    // No selection strategy set.
    it("returns locked KEEP stake factory", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.lockFactoriesVersions(
        keepStakeFactory.address,
        fullyBackedFactory.address,
      )

      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepStakeVendor.setFactory(newKeepFactory.address)

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the locked factory
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        keepStakeFactory.address,
      )
    })

    // ETH-only vendor set.
    // Selection strategy set.
    it("returns locked ETH-only factory", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      await keepFactorySelection.lockFactoriesVersions(
        keepStakeFactory.address,
        fullyBackedFactory.address,
      )

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await fullyBackedVendor.setFactory(newEthFactory.address)

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the locked factory
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        fullyBackedFactory.address,
      )
    })

    it("reverts if the returned factory is not one of the set factories", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
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

    it("returns zero address when KEEP stake vendor returns factory with zero address", async () => {
      await keepStakeVendor.setFactory(constants.ZERO_ADDRESS)

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the KEEP staked factory (zero address)
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        constants.ZERO_ADDRESS,
      )
    })

    it("returns zero address when ETH-only vendor returns factory with zero address", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await fullyBackedVendor.setFactory(constants.ZERO_ADDRESS)

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      // refresh the choice; it should be the ETH-only factory (zero address)
      const selected = await keepFactorySelection.selectFactoryAndRefresh.call()
      await keepFactorySelection.selectFactoryAndRefresh()
      expect(selected, "unexpected factory selected").to.equal(
        constants.ZERO_ADDRESS,
      )
    })
  })

  describe("setMinimumBondableValue", async () => {
    const defaultValue = new BN(999)
    const newValue = new BN(123987)

    // No KEEP stake vendor set.
    it("reverts when KEEP stake vendor not set", async () => {
      const newKeepFactorySelection = await KeepFactorySelectionStub.new()

      await expectRevert(
        newKeepFactorySelection.setMinimumBondableValue(newValue, 5, 3),
        "KEEP-staked vendor not set",
      )
    })

    // KEEP stake vendor set, factory zero.
    it("completes when KEEP stake factory address is zero", async () => {
      await keepStakeVendor.setFactory(constants.ZERO_ADDRESS)

      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakeFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
    })

    // KEEP stake vendor set, factory set.
    it("updates value in KEEP stake factory", async () => {
      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakeFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
    })

    // ETH-only vendor set, factory not set.
    it("updates value in ETH-only factory", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      await fullyBackedVendor.setFactory(constants.ZERO_ADDRESS)

      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakeFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(
        defaultValue,
      )
    })

    // ETH-only vendor set, factory set.
    it("updates value in ETH-only factory", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakeFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(newValue)
    })

    // KEEP stake vendor set, factory set.
    // ETH-only vendor set, factory not set.
    it("updates value in locked factories", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      // Lock factories
      await keepFactorySelection.lockFactoriesVersions(
        keepStakeFactory.address,
        fullyBackedFactory.address,
      )

      // Upgrade factories in vendors
      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepStakeVendor.setFactory(newKeepFactory.address)

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await fullyBackedVendor.setFactory(newEthFactory.address)

      // Set minimum bondable value
      await keepFactorySelection.setMinimumBondableValue(newValue, 5, 3)

      expect(await keepStakeFactory.minimumBondableValue()).to.eq.BN(newValue)
      expect(await fullyBackedFactory.minimumBondableValue()).to.eq.BN(newValue)

      expect(await newKeepFactory.minimumBondableValue()).to.eq.BN(defaultValue)
      expect(await newEthFactory.minimumBondableValue()).to.eq.BN(defaultValue)
    })
  })

  describe("setFullyBackedKeepVendor", async () => {
    it("can be called only one time", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )
      // ok, this was the first time

      await expectRevert(
        keepFactorySelection.setFullyBackedKeepVendor(
          fullyBackedVendor.address,
        ),
        "Fully backed vendor already set",
      )
    })

    it("can not be called for 0 address", async () => {
      await expectRevert(
        keepFactorySelection.setFullyBackedKeepVendor(constants.ZERO_ADDRESS),
        "Invalid address",
      )
    })

    it("obtains factory from vendor", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )

      await keepFactorySelector.setFullyBackedMode()

      // refresh the choice
      await keepFactorySelection.selectFactoryAndRefresh()

      expect(
        await keepFactorySelection.selectFactory(),
        "unexpected factory selected",
      ).to.equal(fullyBackedFactory.address)
    })

    it("reverts when vendor returns factory with zero address", async () => {
      await fullyBackedVendor.setFactory(constants.ZERO_ADDRESS)

      await expectRevert(
        keepFactorySelection.setFullyBackedKeepVendor(
          fullyBackedVendor.address,
        ),
        "Vendor returned invalid factory address",
      )
    })
  })

  describe("setKeepFactorySelector", async () => {
    it("can be called only one time", async () => {
      await keepFactorySelection.setKeepFactorySelector(
        keepFactorySelector.address,
      )
      // ok, this was the first time

      await expectRevert(
        keepFactorySelection.setKeepFactorySelector(
          keepFactorySelector.address,
        ),
        "Factory selector already set",
      )
    })

    it("can not be called for 0 address", async () => {
      await expectRevert(
        keepFactorySelection.setKeepFactorySelector(constants.ZERO_ADDRESS),
        "Invalid address",
      )
    })
  })

  describe("lockFactoriesVersions", async () => {
    it("sets factories version lock", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.lockFactoriesVersions(
        keepStakeFactory.address,
        fullyBackedFactory.address,
      )

      expect(await keepFactorySelection.factoriesVersionsLocked()).to.be.true
    })

    it("gets latest factories from vendor", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.selectFactoryAndRefresh()

      expect(await keepFactorySelection.keepStakeFactory()).to.equal(
        constants.ZERO_ADDRESS,
      )

      expect(await keepFactorySelection.fullyBackedFactory()).to.equal(
        constants.ZERO_ADDRESS,
      )

      const newKeepFactory = await ECDSAKeepFactoryStub.new()
      await keepStakeVendor.setFactory(newKeepFactory.address)

      const newEthFactory = await ECDSAKeepFactoryStub.new()
      await fullyBackedVendor.setFactory(newEthFactory.address)

      await keepFactorySelection.lockFactoriesVersions(
        newKeepFactory.address,
        newEthFactory.address,
      )

      expect(await keepFactorySelection.keepStakeFactory()).to.equal(
        newKeepFactory.address,
      )

      expect(await keepFactorySelection.fullyBackedFactory()).to.equal(
        newEthFactory.address,
      )
    })

    it("can be called only one time", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await keepFactorySelection.lockFactoriesVersions(
        keepStakeFactory.address,
        fullyBackedFactory.address,
      )
      // ok, this was the first time

      await expectRevert(
        keepFactorySelection.lockFactoriesVersions(
          keepStakeFactory.address,
          fullyBackedFactory.address,
        ),
        "Already locked",
      )
    })

    it("reverts when KEEP stake vendor not set", async () => {
      const newKeepFactorySelection = await KeepFactorySelectionStub.new()

      await expectRevert(
        newKeepFactorySelection.lockFactoriesVersions(
          constants.ZERO_ADDRESS,
          fullyBackedFactory.address,
        ),
        "KEEP-staked vendor not set",
      )
    })

    it("reverts when ETH-only vendor not set", async () => {
      await expectRevert(
        keepFactorySelection.lockFactoriesVersions(
          constants.ZERO_ADDRESS,
          fullyBackedFactory.address,
        ),
        "Fully backed vendor not set",
      )
    })

    it("reverts for unexpected KEEP stake factory address", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await expectRevert(
        keepFactorySelection.lockFactoriesVersions(
          thirdParty,
          fullyBackedFactory.address,
        ),
        "Unexpected KEEP-staked factory",
      )
    })

    it("reverts for unexpected ETH-only factory address", async () => {
      await keepFactorySelection.setFullyBackedKeepVendor(
        fullyBackedVendor.address,
      )

      await expectRevert(
        keepFactorySelection.lockFactoriesVersions(
          keepStakeFactory.address,
          thirdParty,
        ),
        "Unexpected fully backed factory",
      )
    })
  })
})
