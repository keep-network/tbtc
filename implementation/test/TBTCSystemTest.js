const {deployAndLinkAll} = require("../testHelpers/testDeployer.js")
const {increaseTime} = require("../testHelpers/utils.js")
const {
  createSnapshot,
  restoreSnapshot,
} = require("../testHelpers/helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, expectRevert, expectEvent} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

const TBTCSystem = contract.fromArtifact("TBTCSystem")

// eslint-disable-next-line no-only-tests/no-only-tests
describe.only("TBTCSystem", async function() {
  let tbtcSystem
  let ecdsaKeepFactory

  before(async () => {
    const {tbtcSystemStub, ecdsaKeepFactoryStub} = await deployAndLinkAll(
      [],
      // Though deployTestDeposit deploys a TBTCSystemStub for us, we want to
      // test TBTCSystem itself.
      {TBTCSystemStub: TBTCSystem},
    )
    // Refer to this correctly throughout the rest of the test.
    tbtcSystem = tbtcSystemStub
    ecdsaKeepFactory = ecdsaKeepFactoryStub
  })

  describe("requestNewKeep()", async () => {
    let openKeepFee
    before(async () => {
      openKeepFee = await ecdsaKeepFactory.openKeepFeeEstimate.call()
    })
    it("sends caller as owner to open new keep", async () => {
      const expectedKeepOwner = accounts[2]

      await tbtcSystem.requestNewKeep(5, 10, 0, {
        from: expectedKeepOwner,
        value: openKeepFee,
      })
      const keepOwner = await ecdsaKeepFactory.keepOwner.call()

      expect(expectedKeepOwner, "incorrect keep owner address").to.equal(
        keepOwner,
      )
    })

    it("returns keep address", async () => {
      const expectedKeepAddress = await ecdsaKeepFactory.keepAddress.call()

      const result = await tbtcSystem.requestNewKeep.call(5, 10, 0, {
        value: openKeepFee,
      })

      expect(expectedKeepAddress, "incorrect keep address").to.equal(result)
    })

    it("forwards value to keep factory", async () => {
      const initialBalance = await web3.eth.getBalance(ecdsaKeepFactory.address)

      await tbtcSystem.requestNewKeep(5, 10, 0, {value: openKeepFee})

      const finalBalance = await web3.eth.getBalance(ecdsaKeepFactory.address)
      const balanceCheck = new BN(finalBalance).sub(new BN(initialBalance))
      expect(
        balanceCheck,
        "TBTCSystem did not correctly forward value to keep factory",
      ).to.eq.BN(openKeepFee)
    })
  })

  describe("update signer fee", async () => {
    describe("beginSignerFeeDivisorUpdate", async () => {
      const newFee = new BN("201")
      it("executes and fires SignerFeeDivisorUpdateStarted event", async () => {
        receipt = await tbtcSystem.beginSignerFeeDivisorUpdate(newFee)

        expectEvent(receipt, "SignerFeeDivisorUpdateStarted", {
          _signerFeeDivisor: newFee,
        })
      })

      it("reverts if msg.sender != owner", async () => {
        await expectRevert.unspecified(
          tbtcSystem.beginSignerFeeDivisorUpdate(newFee, {
            from: accounts[1],
          }),
          "",
        )
      })

      it("reverts if fee divisor is smaller than 10", async () => {
        await expectRevert(
          tbtcSystem.beginSignerFeeDivisorUpdate(new BN("9")),
          "Signer fee divisor must be greater than 9, for a signer fee that is <= 10%.",
        )
      })
    })

    describe("finalizeSignerFeeDivisorUpdate", async () => {
      it("reverts if the governance timer has not elapsed", async () => {
        await expectRevert(
          tbtcSystem.finalizeSignerFeeDivisorUpdate(),
          "Timer not elapsed",
        )
      })

      it("updates signer fee and fires SignerFeeDivisorUpdated event", async () => {
        const remainingTime = await tbtcSystem.geRemainingSignerFeeDivisorUpdateTime()

        await increaseTime(remainingTime.toNumber() + 1)

        receipt = await tbtcSystem.finalizeSignerFeeDivisorUpdate()

        const signerFeeDivisor = await tbtcSystem.getSignerFeeDivisor.call()

        expectEvent(receipt, "SignerFeeDivisorUpdated", {
          _signerFeeDivisor: new BN("201"),
        })
        expect(signerFeeDivisor).to.eq.BN(new BN("201"))
      })

      it("reverts if a change has not been initiated", async () => {
        await expectRevert(
          tbtcSystem.finalizeSignerFeeDivisorUpdate(),
          "Change not initiated",
        )
      })
    })
  })

  describe("update lot sizes", async () => {
    const lotSizes = [new BN(10 ** 8), new BN(10 ** 6)]
    describe("beginLotSizesUpdate", async () => {
      it("executes and emits a LotSizesUpdateStarted event", async () => {
        const block = await web3.eth.getBlock("latest")
        const receipt = await tbtcSystem.beginLotSizesUpdate(lotSizes)
        expectEvent(receipt, "LotSizesUpdateStarted", {})
        expect(receipt.logs[0].args[0][0]).to.eq.BN(lotSizes[0])
        expect(receipt.logs[0].args[0][1]).to.eq.BN(lotSizes[1])
        expect(receipt.logs[0].args[1]).to.eq.BN(block.timestamp)
      })

      it("reverts if lot size array is empty", async () => {
        const lotSizes = []
        await expectRevert(
          tbtcSystem.beginLotSizesUpdate(lotSizes),
          "Lot size array must always contain 1BTC",
        )
      })

      it("reverts if lot size array does not contain a 1BTC lot size", async () => {
        const lotSizes = [10 ** 7]
        await expectRevert(
          tbtcSystem.beginLotSizesUpdate(lotSizes),
          "Lot size array must always contain 1BTC",
        )
      })
    })

    describe("finalizeLotSizesUpdate", async () => {
      it("reverts if the governance timer has not elapsed", async () => {
        await expectRevert(
          tbtcSystem.finalizeLotSizesUpdate(),
          "Timer not elapsed",
        )
      })

      it("updates lot sizes and fires LotSizesUpdated event", async () => {
        const remainingTime = await tbtcSystem.getRemainingLotSizesUpdateTime()

        await increaseTime(remainingTime.toNumber() + 1)

        receipt = await tbtcSystem.finalizeLotSizesUpdate()

        const currentLotSizes = await tbtcSystem.getAllowedLotSizes.call()

        expectEvent(receipt, "LotSizesUpdated", {})
        expect(receipt.logs[0].args._lotSizes[0]).to.eq.BN(lotSizes[0])
        expect(receipt.logs[0].args._lotSizes[1]).to.eq.BN(lotSizes[1])
        expect(currentLotSizes[0]).to.eq.BN(lotSizes[0])
        expect(currentLotSizes[1]).to.eq.BN(lotSizes[1])
      })

      it("reverts if a change has not been initiated", async () => {
        await expectRevert(
          tbtcSystem.finalizeLotSizesUpdate(),
          "Change not initiated",
        )
      })
    })
  })

  describe("update collateralization thresholds", async () => {
    const initialPercent = new BN("150")
    const undercollateralizedPercent = new BN("130")
    const severelyUndercollateralizedPercent = new BN("120")
    describe("beginCollateralizationThresholdsUpdate", async () => {
      it("executes and fires CollateralizationThresholdsUpdateStarted event", async () => {
        receipt = await tbtcSystem.beginCollateralizationThresholdsUpdate(
          initialPercent,
          undercollateralizedPercent,
          severelyUndercollateralizedPercent,
        )

        expectEvent(receipt, "CollateralizationThresholdsUpdateStarted", {
          _initialCollateralizedPercent: initialPercent,
          _undercollateralizedThresholdPercent: undercollateralizedPercent,
          _severelyUndercollateralizedThresholdPercent: severelyUndercollateralizedPercent,
        })
      })

      it("reverts if Initial collateralized percent > 300", async () => {
        await expectRevert(
          tbtcSystem.beginCollateralizationThresholdsUpdate(
            new BN("301"),
            new BN("130"),
            new BN("120"),
          ),
          "Initial collateralized percent must be <= 300%",
        )
      })

      it("reverts if Undercollateralized threshold > initial collateralize percent", async () => {
        await expectRevert(
          tbtcSystem.beginCollateralizationThresholdsUpdate(
            new BN("150"),
            new BN("160"),
            new BN("120"),
          ),
          "Undercollateralized threshold must be < initial collateralized percent",
        )
      })

      it("reverts if Severe undercollateralized threshold > undercollateralized threshold", async () => {
        await expectRevert(
          tbtcSystem.beginCollateralizationThresholdsUpdate(
            new BN("150"),
            new BN("130"),
            new BN("131"),
          ),
          "Severe undercollateralized threshold must be < undercollateralized threshold",
        )
      })
    })

    describe("finalizeCollateralizationThresholdsUpdate", async () => {
      it("reverts if the governance timer has not elapsed", async () => {
        await expectRevert(
          tbtcSystem.finalizeCollateralizationThresholdsUpdate(),
          "Timer not elapsed",
        )
      })

      it("updates collateralization thresholds and fires CollateralizationThresholdsUpdated event", async () => {
        const remainingTime = await tbtcSystem.getRemainingCollateralizationUpdateTime()

        await increaseTime(remainingTime.toNumber() + 1)

        receipt = await tbtcSystem.finalizeCollateralizationThresholdsUpdate()

        const initial = await tbtcSystem.getInitialCollateralizedPercent.call()
        const undercollateralized = await tbtcSystem.getUndercollateralizedThresholdPercent.call()
        const severelyUndercollateralized = await tbtcSystem.getSeverelyUndercollateralizedThresholdPercent.call()

        expectEvent(receipt, "CollateralizationThresholdsUpdated", {
          _initialCollateralizedPercent: initialPercent,
          _undercollateralizedThresholdPercent: undercollateralizedPercent,
          _severelyUndercollateralizedThresholdPercent: severelyUndercollateralizedPercent,
        })

        expect(initialPercent).to.eq.BN(initial)
        expect(undercollateralizedPercent).to.eq.BN(undercollateralized)
        expect(severelyUndercollateralizedPercent).to.eq.BN(
          severelyUndercollateralized,
        )
      })

      it("reverts if a change has not been initiated", async () => {
        await expectRevert(
          tbtcSystem.finalizeCollateralizationThresholdsUpdate(),
          "Change not initiated",
        )
      })
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
