const {deployAndLinkAll} = require("../testHelpers/testDeployer.js")
const {increaseTime} = require("../testHelpers/utils.js")
const {
  createSnapshot,
  restoreSnapshot,
} = require("../testHelpers/helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, expectRevert} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

const TBTCSystem = contract.fromArtifact("TBTCSystem")

describe("TBTCSystem", async function() {
  let tbtcSystem
  let ecdsaKeepFactory
  let tdt

  before(async () => {
    const {
      tbtcSystemStub,
      ecdsaKeepFactoryStub,
      tbtcDepositToken,
    } = await deployAndLinkAll(
      [],
      // Though deployTestDeposit deploys a TBTCSystemStub for us, we want to
      // test TBTCSystem itself.
      {TBTCSystemStub: TBTCSystem},
    )
    // Refer to this correctly throughout the rest of the test.
    tbtcSystem = tbtcSystemStub
    ecdsaKeepFactory = ecdsaKeepFactoryStub
    tdt = tbtcDepositToken
  })

  describe("requestNewKeep()", async () => {
    let openKeepFee
    const keepOwner = accounts[2]
    const tdtOwner = accounts[1]

    before(async () => {
      openKeepFee = await ecdsaKeepFactory.openKeepFeeEstimate.call()
      await tdt.forceMint(tdtOwner, web3.utils.toBN(keepOwner))
    })

    it("sends caller as owner to open new keep", async () => {
      await tbtcSystem.requestNewKeep(5, 10, 0, {
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

      const result = await tbtcSystem.requestNewKeep.call(5, 10, 0, {
        value: openKeepFee,
        from: keepOwner,
      })

      expect(expectedKeepAddress, "incorrect keep address").to.equal(result)
    })

    it("forwards value to keep factory", async () => {
      const initialBalance = await web3.eth.getBalance(ecdsaKeepFactory.address)

      await tbtcSystem.requestNewKeep(5, 10, 0, {
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
        tbtcSystem.requestNewKeep(5, 10, 0, {
          value: openKeepFee,
          from: accounts[0],
        }),
        "Caller must be a Deposit contract",
      )
    })

    it("reverts if caller is the owner of a valid TDT", async () => {
      await expectRevert(
        tbtcSystem.requestNewKeep(5, 10, 0, {
          value: openKeepFee,
          from: tdtOwner,
        }),
        "Caller must be a Deposit contract",
      )
    })
  })

  describe("setSignerFeeDivisor", async () => {
    it("sets the signer fee", async () => {
      await tbtcSystem.setSignerFeeDivisor(new BN("201"))

      const signerFeeDivisor = await tbtcSystem.getSignerFeeDivisor()
      expect(signerFeeDivisor).to.eq.BN(new BN("201"))
    })

    it("reverts if msg.sender != owner", async () => {
      await expectRevert.unspecified(
        tbtcSystem.setSignerFeeDivisor(new BN("201"), {
          from: accounts[1],
        }),
        "",
      )
    })
  })

  describe("setLotSizes", async () => {
    it("sets a different lot size array", async () => {
      const blockNumber = await web3.eth.getBlock("latest").number
      const lotSizes = [10 ** 8, 10 ** 6]
      await tbtcSystem.setLotSizes(lotSizes)

      const eventList = await tbtcSystem.getPastEvents("LotSizesUpdated", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
      expect(eventList[0].returnValues._lotSizes).to.eql([
        "100000000",
        "1000000",
      ]) // deep equality check
    })

    it("reverts if lot size array is empty", async () => {
      const lotSizes = []
      await expectRevert(
        tbtcSystem.setLotSizes(lotSizes),
        "Lot size array must always contain 1BTC",
      )
    })

    it("reverts if lot size array does not contain a 1BTC lot size", async () => {
      const lotSizes = [10 ** 7]
      await expectRevert(
        tbtcSystem.setLotSizes(lotSizes),
        "Lot size array must always contain 1BTC",
      )
    })
  })

  describe("emergencyPauseNewDeposits", async () => {
    let term

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("pauses new deposit creation", async () => {
      await tbtcSystem.emergencyPauseNewDeposits()

      const allowNewDeposits = await tbtcSystem.getAllowNewDeposits()
      expect(allowNewDeposits).to.equal(false)
    })

    it("reverts if msg.sender is not owner", async () => {
      await expectRevert(
        tbtcSystem.emergencyPauseNewDeposits({from: accounts[1]}),
        "Ownable: caller is not the owner",
      )
    })

    it("does not allows new deposit re-activation before 10 days", async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber() - 10) // T-10 seconds. toNumber because increaseTime doesn't support BN

      await expectRevert(
        tbtcSystem.resumeNewDeposits(),
        "Deposits are still paused",
      )
    })

    it("allows new deposit creation after 10 days", async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber()) // 10 days
      await tbtcSystem.resumeNewDeposits()
      const allowNewDeposits = await tbtcSystem.getAllowNewDeposits()
      expect(allowNewDeposits).to.equal(true)
    })

    it("reverts if emergencyPauseNewDeposits has already been called", async () => {
      await tbtcSystem.emergencyPauseNewDeposits()
      term = await tbtcSystem.getRemainingPauseTerm()

      await increaseTime(term.toNumber()) // 10 days
      tbtcSystem.resumeNewDeposits()

      await expectRevert(
        tbtcSystem.emergencyPauseNewDeposits(),
        "emergencyPauseNewDeposits can only be called once",
      )
    })
  })
})
