const { createInstance, toAwaitingWithdrawalSignature } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- Signature-timeout", async function () {
  const lotSize = "10000000";
  const lotSizeTbtc = new BN("10000000000").mul(new BN(lotSize));
  const depositInitiator = accounts[1]
  const auctionBuyer = accounts[3]
  let testDeposit

  describe("Signature-timeout no FRT minted", async () => {

    before(async () => {
      ;({
        tbtcConstants,
        tbtcToken,
        collateralAmount,
        deposits
      } = await createInstance(
        { collateral: 125,
          lotSize : lotSize,
        state: states.AWAITING_WITHDRAWAL_SIGNATURE.toNumber(),
        depositOwner: depositInitiator }
      ))
      testDeposit = deposits
    })

    it("unable to start liquidation when timer has not elapsed", async () => {
      await expectRevert(
        testDeposit.notifyRedemptionSignatureTimedOut(),
        "Signature timer has not elapsed",
      )
      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it("starts liquidation auction", async () => {
      await ecdsaKeepStub.setSuccess(true)

      const timer = await tbtcConstants.getSignatureTimeout.call()
      await increaseTime(timer.toNumber() + 1)

      await testDeposit.notifyRedemptionSignatureTimedOut()

      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("reverts if no TBTC balance has been approved by the auction buyer to the Deposit", async () => {
      await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        "Not enough TBTC to cover outstanding debt",
      )
      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("liquidates correctly", async () => {
      const duration = await tbtcConstants.getAuctionDuration.call()
      await increaseTime(duration.toNumber())

      await tbtcToken.approve(testDeposit.address, lotSizeTbtc, { from: auctionBuyer })
      await testDeposit.purchaseSignerBondsAtAuction({ from: auctionBuyer })

      const allowance = await testDeposit.withdrawableAmount({ from: auctionBuyer })
      await testDeposit.withdrawFunds({ from: auctionBuyer })

      const depositState = await testDeposit.currentState.call()
      const endingBalance = await web3.eth.getBalance(testDeposit.address)

      expect(allowance).to.eq.BN(collateralAmount)
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
    })
  })

  describe("Signature-timeout - FRT minted", async () => {

    before(async () => {
      ;({
        tbtcConstants,
        tbtcToken,
        collateralAmount,
        feeRebateToken,
        deposits
      } = await createInstance(
        { collateral: 125,
        state: states.ACTIVE.toNumber(),
        depositOwner: depositInitiator }
      ))
      testDeposit = deposits
    })

    it("trades TDT for TBTC and FRT via the VendingMachine", async () => {
      tdtId = await web3.utils.toBN(testDeposit.address)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: depositInitiator
      })
      await vendingMachine.tdtToTbtc(tdtId, {from: depositInitiator})
      const frtOwner = await feeRebateToken.ownerOf(tdtId)
      expect(frtOwner).to.equal(depositInitiator)
      
    })
    it("Faisl to redeem while sufficiently collateralized, within term, and TDT owner not caller", async () => {
      await expectRevert(
        toAwaitingWithdrawalSignature(testDeposit),
        "Only TDT holder can redeem unless deposit is at-term or in COURTESY_CALL",
      )
    })
    it("Purchases TDT and moves to AWAITING_WITHDRAWAL_SIGNATURE", async () => {
      await tbtcToken.resetBalance(lotSizeTbtc, { from: depositInitiator })
      await tbtcToken.approve(vendingMachine.address, lotSizeTbtc, { from: depositInitiator })
      await vendingMachine.tbtcToTdt(tdtId, {from: depositInitiator})

      const tdtOwner = await tbtcDepositToken.ownerOf(tdtId)
      await toAwaitingWithdrawalSignature(testDeposit)
      const depositState = await testDeposit.currentState.call()

      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
      expect(tdtOwner).to.equal(depositInitiator)

    })

    it("unable to start liquidation when timer has not elapsed", async () => {
      await expectRevert(
        testDeposit.notifyRedemptionSignatureTimedOut(),
        "Signature timer has not elapsed",
      )
    })

    it("starts liquidation auction", async () => {
      await ecdsaKeepStub.setSuccess(true)

      const timer = await tbtcConstants.getSignatureTimeout.call()
      await increaseTime(timer.toNumber() + 1)

      await testDeposit.notifyRedemptionSignatureTimedOut()

      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("reverts if no TBTC balance has been approved by the auction buyer to the Deposit", async () => {
      await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        "Not enough TBTC to cover outstanding debt",
      )
      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("liquidates correctly", async () => {
      const duration = await tbtcConstants.getAuctionDuration.call()
      await increaseTime(duration.toNumber())

      await tbtcToken.approve(testDeposit.address, lotSizeTbtc, { from: auctionBuyer })
      await testDeposit.purchaseSignerBondsAtAuction({ from: auctionBuyer })

      const allowance = await testDeposit.withdrawableAmount({ from: auctionBuyer })
      await testDeposit.withdrawFunds({ from: auctionBuyer })

      const depositState = await testDeposit.currentState.call()
      const endingBalance = await web3.eth.getBalance(testDeposit.address)

      expect(allowance).to.eq.BN(collateralAmount)
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
    })
  })
})
