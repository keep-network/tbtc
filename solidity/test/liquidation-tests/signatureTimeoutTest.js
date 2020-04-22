const { createInstance, toAwaitingWithdrawalSignature } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- Signature-timeout", async function () {
  const lotSize = "10000000";
  const lotSizeTbtc = new BN("10000000000").mul(new BN(lotSize));
  const fee = new BN("500000000000000")
  const depositInitiator = accounts[1]
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
        testDeposit.notifySignatureTimeout(),
        "Signature timer has not elapsed",
      )
      const depositState = await testDeposit.getCurrentState.call()
      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it("liquidates correctly", async () => {
      const requirement = await testDeposit.getOwnerRedemptionTbtcRequirement.call(depositInitiator)
      const timer = await tbtcConstants.getRedemptionProofTimeout.call()
      await increaseTime(timer.toNumber())

      await testDeposit.notifySignatureTimeout()

      const depositState = await testDeposit.getCurrentState.call()
      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: depositInitiator,
      })
      await testDeposit.withdrawFunds({from: depositInitiator})
      const endingBalance = await web3.eth.getBalance(testDeposit.address)
  
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
      expect(requirement).to.eq.BN(fee)
      expect(withdrawable).to.eq.BN(collateralAmount)
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
      const depositState = await testDeposit.getCurrentState.call()

      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
      expect(tdtOwner).to.equal(depositInitiator)

    })

    it("unable to start liquidation when timer has not elapsed", async () => {
      await expectRevert(
        testDeposit.notifySignatureTimeout(),
        "Signature timer has not elapsed",
      )
    })
    
    it("liquidates correctly", async () => {
      const requirement = await testDeposit.getOwnerRedemptionTbtcRequirement.call(depositInitiator)
      const timer = await tbtcConstants.getRedemptionProofTimeout.call()
      await increaseTime(timer.toNumber())

      await testDeposit.notifySignatureTimeout()

      const depositState = await testDeposit.getCurrentState.call()
      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: depositInitiator,
      })
      await testDeposit.withdrawFunds({from: depositInitiator})
      const endingBalance = await web3.eth.getBalance(testDeposit.address)
  
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
      expect(requirement).to.eq.BN(new BN(0))
      expect(withdrawable).to.eq.BN(collateralAmount)
    })
  })
})
