const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {increaseTime} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, expectRevert, expectEvent} = require("@openzeppelin/test-helpers")
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
    const tdtOwner = accounts[1]
    const keepOwner = accounts[2]

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

  describe("geRemainingSignerFeeDivisorUpdateTime", async () => {
    let totalDelay

    before(async () => {
      totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("reverts if update has not been initiated", async () => {
      await expectRevert(
        tbtcSystem.geRemainingSignerFeeDivisorUpdateTime.call(),
        "Update not initiated",
      )
    })

    it("returns total delay if no time has passed ", async () => {
      await tbtcSystem.beginSignerFeeDivisorUpdate(new BN("200"))
      const remaining = await tbtcSystem.geRemainingSignerFeeDivisorUpdateTime.call()
      expect(remaining).to.eq.BN(totalDelay)
    })

    it("returns the correct remaining time", async () => {
      await tbtcSystem.beginSignerFeeDivisorUpdate(new BN("200"))
      const expectedRemaining = 100
      await increaseTime(totalDelay.toNumber() - expectedRemaining)

      const remaining = await tbtcSystem.geRemainingSignerFeeDivisorUpdateTime.call()
      expect([expectedRemaining, expectedRemaining + 1]).to.include.toString(
        remaining.toNumber(),
      )
    })
  })

  describe("getRemainingLotSizesUpdateTime", async () => {
    let totalDelay
    const lotSizes = [new BN(10 ** 8), new BN(10 ** 6)]

    before(async () => {
      totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("reverts if update has not been initiated", async () => {
      await expectRevert(
        tbtcSystem.getRemainingLotSizesUpdateTime.call(),
        "Update not initiated",
      )
    })

    it("returns total delay if no time has passed ", async () => {
      await tbtcSystem.beginLotSizesUpdate(lotSizes)
      const remaining = await tbtcSystem.getRemainingLotSizesUpdateTime.call()
      expect(remaining).to.eq.BN(totalDelay)
    })

    it("returns the correct remaining time", async () => {
      await tbtcSystem.beginLotSizesUpdate(lotSizes)
      const expectedRemaining = 100
      await increaseTime(totalDelay.toNumber() - expectedRemaining)

      const remaining = await tbtcSystem.getRemainingLotSizesUpdateTime.call()
      expect([expectedRemaining, expectedRemaining + 1]).to.include.toString(
        remaining.toNumber(),
      )
    })
  })

  describe("getRemainingCollateralizationUpdateTime", async () => {
    let totalDelay

    before(async () => {
      totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("reverts if update has not been initiated", async () => {
      await expectRevert(
        tbtcSystem.getRemainingCollateralizationUpdateTime.call(),
        "Update not initiated",
      )
    })

    it("returns total delay if no time has passed ", async () => {
      await tbtcSystem.beginCollateralizationThresholdsUpdate(
        new BN("150"),
        new BN("130"),
        new BN("120"),
      )
      const remaining = await tbtcSystem.getRemainingCollateralizationUpdateTime.call()
      expect(remaining).to.eq.BN(totalDelay)
    })

    it("returns the correct remaining time", async () => {
      await tbtcSystem.beginCollateralizationThresholdsUpdate(
        new BN("150"),
        new BN("130"),
        new BN("120"),
      )
      const expectedRemaining = 100
      await increaseTime(totalDelay.toNumber() - expectedRemaining)

      const remaining = await tbtcSystem.getRemainingCollateralizationUpdateTime.call()
      expect([expectedRemaining, expectedRemaining + 1]).to.include.toString(
        remaining.toNumber(),
      )
    })
  })

  describe("update signer fee", async () => {
    describe("beginSignerFeeDivisorUpdate", async () => {
      const newFee = new BN("201")
      it("executes and fires SignerFeeDivisorUpdateStarted event", async () => {
        receipt = await tbtcSystem.beginSignerFeeDivisorUpdate(new BN("200"))

        expectEvent(receipt, "SignerFeeDivisorUpdateStarted", {
          _signerFeeDivisor: new BN("200"),
        })
      })

      it("overrides previous update and resets timer", async () => {
        receipt = await tbtcSystem.beginSignerFeeDivisorUpdate(newFee)
        const remainingTime = await tbtcSystem.geRemainingSignerFeeDivisorUpdateTime.call()
        const totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()

        expectEvent(receipt, "SignerFeeDivisorUpdateStarted", {
          _signerFeeDivisor: newFee,
        })
        expect([
          remainingTime.toString(),
          remainingTime.toString() - 1,
        ]).to.include(totalDelay.toString())
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
        const testSizes = [new BN(10 ** 8), new BN(10 ** 6)]
        const block = await web3.eth.getBlock("latest")
        const receipt = await tbtcSystem.beginLotSizesUpdate(testSizes)
        expectEvent(receipt, "LotSizesUpdateStarted", {})
        expect(receipt.logs[0].args[0][0]).to.eq.BN(testSizes[0])
        expect(receipt.logs[0].args[0][1]).to.eq.BN(testSizes[1])
        expect([
          receipt.logs[0].args[1].toString(),
          receipt.logs[0].args[1].toString() - 1,
        ]).to.include(block.timestamp.toString())
      })

      it("overrides previous update and resets timer", async () => {
        const block = await web3.eth.getBlock("latest")
        const receipt = await tbtcSystem.beginLotSizesUpdate(lotSizes)
        const remainingTime = await tbtcSystem.getRemainingLotSizesUpdateTime.call()
        const totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()

        expectEvent(receipt, "LotSizesUpdateStarted", {})
        expect(receipt.logs[0].args[0][0]).to.eq.BN(lotSizes[0])
        expect(receipt.logs[0].args[0][1]).to.eq.BN(lotSizes[1])
        expect([
          receipt.logs[0].args[1].toString(),
          receipt.logs[0].args[1].toString() - 1,
        ]).to.include(block.timestamp.toString())
        expect([
          remainingTime.toString(),
          remainingTime.toString() - 1,
        ]).to.include(totalDelay.toString())
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
          new BN("213"),
          new BN("156"),
          new BN("128"),
        )

        expectEvent(receipt, "CollateralizationThresholdsUpdateStarted", {
          _initialCollateralizedPercent: new BN("213"),
          _undercollateralizedThresholdPercent: new BN("156"),
          _severelyUndercollateralizedThresholdPercent: new BN("128"),
        })
      })

      it("overrides previous update and resets timer", async () => {
        receipt = await tbtcSystem.beginCollateralizationThresholdsUpdate(
          initialPercent,
          undercollateralizedPercent,
          severelyUndercollateralizedPercent,
        )

        const remainingTime = await tbtcSystem.getRemainingCollateralizationUpdateTime.call()
        const totalDelay = await tbtcSystem.getGovernanceTimeDelay.call()

        expectEvent(receipt, "CollateralizationThresholdsUpdateStarted", {
          _initialCollateralizedPercent: initialPercent,
          _undercollateralizedThresholdPercent: undercollateralizedPercent,
          _severelyUndercollateralizedThresholdPercent: severelyUndercollateralizedPercent,
        })
        expect([
          remainingTime.toString(),
          remainingTime.toString() - 1,
        ]).to.include(totalDelay.toString())
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
