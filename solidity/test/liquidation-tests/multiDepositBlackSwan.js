const { createInstance } = require("./liquidation-test-utils/createInstance.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { accounts, web3 } = require("@openzeppelin/test-environment")
const { BN } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

describe("Integration -- black-swan", async function () {
  const lotSize = new BN("10000000");
  const lotSizeTbtc = new BN("10000000000").mul(lotSize);
  const depositInitiator = accounts[2]

    before(async () => {
      ; ({
        mockSatWeiPriceFeed,
        tbtcConstants,
        tbtcToken,
        collateralAmount,
        deposits,
        vendingMachine
      } = await createInstance(
        {collateral:125, numToCreate: 10, depositOwner: depositInitiator}
      ))
    })


  it("trade TDTs for TBTC", async () => {
      for(let i = 0; i < deposits.length; i++) {
        tdtId = await web3.utils.toBN(deposits[i].address)
        await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
          from: depositInitiator,
        })
        await vendingMachine.tdtToTbtc(tdtId, {from: depositInitiator})
      }
    });

  it("moves deposits to liquidation, large relative collateral value dip", async () => {
    for(let i = 0; i < deposits.length; i++) {
      await mockSatWeiPriceFeed.setPrice(new BN("4666666666666"))// price down 90%
      await deposits[i].notifyUndercollateralizedLiquidation()

      const depositState = await deposits[i].getCurrentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    }
  })

  it("Liquidate Deposits", async () => {
    await createSnapshot()
    let liquidatedIndex
    let totalGasCost
    const initialBalance = await web3.eth.getBalance(depositInitiator)
    const gasPrice = await web3.eth.getGasPrice()

    for(let i = 0; i < deposits.length; i++) {
      const requirement = await deposits[i].getOwnerRedemptionTbtcRequirement.call(depositInitiator)
      const totalRequirement = lotSizeTbtc.add(requirement)
      const available = await tbtcToken.balanceOf(depositInitiator)
      const duration = await tbtcConstants.getAuctionDuration.call()
      await increaseTime(duration.toNumber())

      if(available.gte(totalRequirement)){
        const tx0 = await tbtcToken.approve(deposits[i].address, totalRequirement, {from: depositInitiator})
        const tx1 = await deposits[i].purchaseSignerBondsAtAuction({ from: depositInitiator })
        const tx2 = await deposits[i].withdrawFunds({ from: depositInitiator })
        const gasUSed = new BN(tx0.receipt.cumulativeGasUsed).add(new BN(tx1.receipt.cumulativeGasUsed).add(new BN(tx2.receipt.cumulativeGasUsed)))
        const totalTxCost = new BN(gasPrice).mul(new BN(gasUSed))
        totalGasCost = new BN(totalGasCost).add(totalTxCost)

        const endingBalance = await web3.eth.getBalance(deposits[i].address)
        expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      }else{
        liquidatedIndex = i
        break
      }
    }

    const expectedReward = new BN(collateralAmount).mul(new BN(liquidatedIndex))
    const finalBalance = await web3.eth.getBalance(depositInitiator)
    expect(liquidatedIndex,"Should be able to liquidate all but 1 deposit").to.eq.BN(deposits.length - 1)
    expect(new BN(initialBalance).add(new BN(expectedReward))).to.eq.BN(new BN(finalBalance).add(totalGasCost))

    await restoreSnapshot()
  })
})
