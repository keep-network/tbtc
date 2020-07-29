const { createInstance } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime, bytes32zero } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- fraud proof", async function () {
  const lotSize = new BN("10000000");
  const lotSizeTbtc = new BN("10000000000").mul(lotSize);
  const auctionBuyer = accounts[3]
  let testDeposit

  before(async () => {
    ; ({
      ecdsaKeepStub,
      tbtcConstants,
      tbtcToken,
      collateralAmount,
      deposits
    } = await createInstance({collateral: 125,  depositOwner: accounts[1]}))
    testDeposit = deposits
  })

  it("reverts if not fraud according to keep", async () => {
    await ecdsaKeepStub.setSuccess(false)

    await expectRevert(
      testDeposit.provideECDSAFraudProof(
        0,
        bytes32zero,
        bytes32zero,
        bytes32zero,
        "0x00",
      ),
      "Signature is not fraud",
    )

    const depositState = await testDeposit.getState.call()
    expect(depositState).to.eq.BN(states.ACTIVE)
  })

  it("starts liquidation if keep confirms fraud", async () => {
    await ecdsaKeepStub.setSuccess(true)

    await testDeposit.provideECDSAFraudProof(
      0,
      bytes32zero,
      bytes32zero,
      bytes32zero,
      "0x00",
    )

    const depositState = await testDeposit.getState.call()
    expect(depositState).to.eq.BN(states.FRAUD_LIQUIDATION_IN_PROGRESS)
  })

  it("reverts if no TBTC balance has been approved by the buyer to the Deposit", async () => {
    await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
    await expectRevert(
      testDeposit.purchaseSignerBondsAtAuction(),
      "Not enough TBTC to cover outstanding debt",
    )
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.FRAUD_LIQUIDATION_IN_PROGRESS)
  })

  it("completes liquidation correctly - end-time period", async () => {
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
