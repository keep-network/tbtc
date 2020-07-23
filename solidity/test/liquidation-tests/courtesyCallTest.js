const { createInstance } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- courtesy_call", async function () {
  const lotSize = new BN("10000000");
  const lotSizeTbtc = new BN("10000000000").mul(lotSize);
  const satwei = new BN("466666666666")
  const auctionBuyer = accounts[3]
  const liqInitiator = accounts[4]
  const depositOwner = accounts[1]
  let testDeposit

  before(async () => {
    ; ({
      mockSatWeiPriceFeed,
      ecdsaKeepStub,
      tbtcConstants,
      tbtcToken,
      collateralAmount,
      deposits
    } = await createInstance(
      { collateral: 125, depositOwner: accounts[1] }
    ))
    testDeposit = deposits
  })

  it("unable start courtesy-call with sufficient collateral", async () => {
    await expectRevert(
      testDeposit.notifyCourtesyCall(),
      "Signers have sufficient collateral",
    )
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.ACTIVE)
  })

  it("starts courtesy-call when bond value falls under collateralization requirement", async () => {
    // increaseing price of sat by just 1 wei, we're over the limit.
    await mockSatWeiPriceFeed.setPrice(satwei.add(new BN(1)))
    await testDeposit.notifyCourtesyCall()

    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.COURTESY_CALL)
  })

  it("unabse to liquidate if timer has not elapsed", async () => {
    await expectRevert(
      testDeposit.notifyCourtesyCallExpired(),
      "Courtesy period has not elapsed",
    )
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.COURTESY_CALL)
  })

  it("correclty exits courtesy_call if price recovers", async () => {
    await mockSatWeiPriceFeed.setPrice(satwei)
    await testDeposit.exitCourtesyCall()

    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.ACTIVE)
  })

  it("cannot exit courtesy_call from ACTIVE state", async () => {
    await expectRevert(
      testDeposit.exitCourtesyCall(),
      "Not currently in courtesy call",
    )

    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.ACTIVE)
  })

  it("moves back to courtesy_call if price falls below collateralization requirement again", async () => {
    await mockSatWeiPriceFeed.setPrice(satwei.add(new BN(1)))
    await testDeposit.notifyCourtesyCall()

    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.COURTESY_CALL)
  })

  it("Correctly starts liquidation if courtesy timeout elapses", async () => {
    timer = await tbtcConstants.getCourtesyCallTimeout.call()
    await increaseTime(timer.toNumber())

    await testDeposit.notifyCourtesyCallExpired({ from: liqInitiator })
    // not fraud and we did not come from redemption.
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
  })

  it("Cannot exit liquidation once auction starts", async () => {
    await mockSatWeiPriceFeed.setPrice(satwei)
    await expectRevert(
      testDeposit.exitCourtesyCall(),
      "Not currently in courtesy call",
    )

    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
  })

  it("reverts if no TBTC balance has been approved by the buyer to the Deposit", async () => {
    await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
    await expectRevert(
      testDeposit.purchaseSignerBondsAtAuction(),
      "Not enough TBTC to cover outstanding debt",
    )
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
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

  it("starts liquidation correctly - end-time period - FRT minted", async () => {
    ;({ deposits, feeRebateToken, vendingMachine} = await createInstance(
      { collateral: 125, depositOwner: accounts[1] }
    ))
    testDeposit = deposits
    tdtId = await web3.utils.toBN(testDeposit.address)
    await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
      from: depositOwner
    })

    await vendingMachine.tdtToTbtc(tdtId, {from: depositOwner})
    const tdtOwner = await feeRebateToken.ownerOf(tdtId)
    expect(tdtOwner).to.equal(depositOwner)

    await mockSatWeiPriceFeed.setPrice(satwei.add(new BN(1)))
    await testDeposit.notifyCourtesyCall()

    timer = await tbtcConstants.getCourtesyCallTimeout.call()
    await increaseTime(timer.toNumber())
    await testDeposit.notifyCourtesyCallExpired({ from: liqInitiator })
    // not fraud and we did not come from redemption.
    const depositState = await testDeposit.currentState.call()
    expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
  })

  it("completes liquidation correctly - end-time period - FRT minted", async () => {
    const duration = await tbtcConstants.getAuctionDuration.call()
    await increaseTime(duration.toNumber())

    await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
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
