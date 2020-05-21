const { createInstance, toAwaitingWithdrawalProof } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- Redemption-proof timeout", async function () {
  const fee = new BN("50000000000000")
  const depositInitiator = accounts[1]
  let testDeposit

  before(async () => {
    ; ({
      tbtcConstants,
      collateralAmount,
      deposits
    } = await createInstance(
      { collateral: 125,
      state: states.AWAITING_WITHDRAWAL_SIGNATURE.toNumber(),
      depositOwner: depositInitiator }
    ))
    testDeposit = deposits
  })

  describe("Redemption-proof timeout", async () => {
    it("unable to start liquidation if timer not elapsed", async () => {
      await expectRevert(
        testDeposit.notifyRedemptionProofTimeout(),
        "Not currently awaiting a redemption proof",
      )
      const depositState = await testDeposit.getCurrentState.call()
      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it("liquidates correctly", async () => {
     await toAwaitingWithdrawalProof(testDeposit)
      //  AWAITING_WITHDRAWAL_SIGNATURE -> AWAITING_WITHDRAWAL_PROOF
      const requirement = await testDeposit.getOwnerRedemptionTbtcRequirement.call(depositInitiator)
      const timer = await tbtcConstants.getRedemptionProofTimeout.call()
      await increaseTime(timer.toNumber())

      await testDeposit.notifyRedemptionProofTimeout()

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
})
