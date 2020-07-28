const { createInstance, toAwaitingWithdrawalProof } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime, expectEvent, resolveAllLogs } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- Redemption-proof timeout", async function () {
  const lotSize = new BN("10000000");
  const lotSizeTbtc = new BN("10000000000").mul(lotSize);
  const [depositInitiator, redeemer, auctionBuyer] = accounts

  let tbtcConstants
  let tbtcToken
  let deposits

  let testDeposit
  let ecdsaKeepStub
  let collateralAmount

  before(async () => {
    ; ({
      ecdsaKeepStub,
      tbtcConstants,
      tbtcToken,
      collateralAmount,
      deposits,
    } = await createInstance(
      { collateral: 125,
      state: states.AWAITING_WITHDRAWAL_SIGNATURE.toNumber(),
      depositOwner: depositInitiator }
    ))
    testDeposit = deposits

    await testDeposit.setRedeemerAddress(redeemer)
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

    it("starts liquidation auction", async () => {
      await ecdsaKeepStub.setSuccess(true)

     await toAwaitingWithdrawalProof(testDeposit)
      //  AWAITING_WITHDRAWAL_SIGNATURE -> AWAITING_WITHDRAWAL_PROOF
      const timer = await tbtcConstants.getRedemptionProofTimeout.call()
      await increaseTime(timer.toNumber())

      await testDeposit.notifyRedemptionProofTimeout()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("reverts if no TBTC balance has been approved by the auction buyer to the Deposit", async () => {
      await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        "Not enough TBTC to cover outstanding debt",
      )
      const depositState = await testDeposit.getCurrentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("liquidates correctly", async () => {
      const duration = await tbtcConstants.getAuctionDuration.call()
      await increaseTime(duration.toNumber())

      await tbtcToken.approve(testDeposit.address, lotSizeTbtc, { from: auctionBuyer })
      const { receipt } = await testDeposit.purchaseSignerBondsAtAuction({ from: auctionBuyer })

      expectEvent(
        resolveAllLogs(receipt, { tbtcToken }),
        "Transfer",
        {
          "from": auctionBuyer,
          "to": redeemer,
          "value": lotSizeTbtc,
        }
      )

      const allowance = await testDeposit.getWithdrawAllowance({ from: auctionBuyer })
      await testDeposit.withdrawFunds({ from: auctionBuyer })

      const depositState = await testDeposit.getCurrentState.call()
      const endingBalance = await web3.eth.getBalance(testDeposit.address)

      expect(allowance).to.eq.BN(collateralAmount)
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
    })
  })
})
