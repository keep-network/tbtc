const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, expectRevert} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

const TBTCSystem = contract.fromArtifact("TBTCSystem")
const SatWeiPriceFeed = contract.fromArtifact("SatWeiPriceFeed")
const MockMedianizer = contract.fromArtifact("MockMedianizer")

describe("TBTCSystem", async function() {
  let tbtcSystem
  let ecdsaKeepFactory
  let tdt

  const medianizerValue = 100000000000

  before(async () => {
    const {
      tbtcSystemStub,
      ecdsaKeepFactoryStub,
      tbtcDepositToken,
      mockSatWeiPriceFeed,
    } = await deployAndLinkAll(
      [],
      // Though deployTestDeposit deploys a TBTCSystemStub for us, we want to
      // test TBTCSystem itself.
      {
        TBTCSystemStub: TBTCSystem,
        MockSatWeiPriceFeed: SatWeiPriceFeed,
      },
    )
    // Refer to this correctly throughout the rest of the test.
    tbtcSystem = tbtcSystemStub
    ecdsaKeepFactory = ecdsaKeepFactoryStub
    tdt = tbtcDepositToken

    const ethBtcMedianizer = await MockMedianizer.new()
    await ethBtcMedianizer.setValue(medianizerValue)
    mockSatWeiPriceFeed.initialize(tbtcSystem.address, ethBtcMedianizer.address)
  })

  describe("requestNewKeep()", async () => {
    let openKeepFee
    const tdtOwner = accounts[1]
    const keepOwner = accounts[2]
    const maxSecuredLifetime = 123

    before(async () => {
      openKeepFee = await ecdsaKeepFactory.openKeepFeeEstimate.call()
      await tdt.forceMint(tdtOwner, web3.utils.toBN(keepOwner))
    })

    it("sends caller as owner to open new keep", async () => {
      await tbtcSystem.requestNewKeep(5, 10, 0, maxSecuredLifetime, {
        from: keepOwner,
        value: openKeepFee,
      })
      const actualKeepOwner = await ecdsaKeepFactory.keepOwner.call()

      expect(keepOwner, "incorrect keep owner address").to.equal(
        actualKeepOwner,
      )
    })

    it("returns keep address", async () => {
      const expectedKeepAddress = await ecdsaKeepFactory.keepAddress.call()

      const result = await tbtcSystem.requestNewKeep.call(
        5,
        10,
        0,
        maxSecuredLifetime,
        {
          value: openKeepFee,
          from: keepOwner,
        },
      )

      expect(expectedKeepAddress, "incorrect keep address").to.equal(result)
    })

    it("forwards value to keep factory", async () => {
      const initialBalance = await web3.eth.getBalance(ecdsaKeepFactory.address)

      await tbtcSystem.requestNewKeep(5, 10, 0, maxSecuredLifetime, {
        value: openKeepFee,
        from: keepOwner,
      })

      const finalBalance = await web3.eth.getBalance(ecdsaKeepFactory.address)
      const balanceCheck = new BN(finalBalance).sub(new BN(initialBalance))
      expect(
        balanceCheck,
        "TBTCSystem did not correctly forward value to keep factory",
      ).to.eq.BN(openKeepFee)
    })

    it("reverts if caller does not match a valid TDT", async () => {
      await expectRevert(
        tbtcSystem.requestNewKeep(5, 10, 0, maxSecuredLifetime, {
          value: openKeepFee,
          from: accounts[0],
        }),
        "Caller must be a Deposit contract",
      )
    })

    it("reverts if caller is the owner of a valid TDT", async () => {
      await expectRevert(
        tbtcSystem.requestNewKeep(5, 10, 0, maxSecuredLifetime, {
          value: openKeepFee,
          from: tdtOwner,
        }),
        "Caller must be a Deposit contract",
      )
    })
  })

  describe("when fetching Bitcoin price", async () => {
    const tdtOwner = accounts[1]
    const keepOwner = accounts[2]

    it("should revert if the caller does not have an associated TDT", async () => {
      await expectRevert(
        tbtcSystem.fetchBitcoinPrice(),
        "Caller must be a Deposit contract",
      )
    })

    it("should give the price if the caller does have an associated TDT", async () => {
      await tdt.forceMint(tdtOwner, web3.utils.toBN(keepOwner))
      const priceFeedValue = new BN(10)
        .pow(new BN(28))
        .div(new BN(medianizerValue))

      expect(
        await tbtcSystem.fetchBitcoinPrice.call({from: keepOwner}),
      ).to.eq.BN(priceFeedValue)
    })
  })
})
