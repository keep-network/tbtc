const { createInstance } = require('./liquidation-test-utils/createInstance.js')
const { states, increaseTime } = require('../helpers/utils.js')
const { accounts, web3 } = require('@openzeppelin/test-environment')
const { BN, expectRevert } = require('@openzeppelin/test-helpers')
const { expect } = require('chai')

describe('Integration -- Undercollateralized', async function() {
  const lotSize = '10000000'
  const lotSizeTbtc = new BN('10000000000').mul(new BN(lotSize))
  const satwei = new BN('466666666666')
  const auctionBuyer = accounts[2]
  const depositInitiator = accounts[1]
  let testDeposit

  before(async () => {
    ; ({
      mockSatWeiPriceFeed,
      tbtcConstants,
      tbtcToken,
      collateralAmount,
      deposits,
    } = await createInstance(
      { collateral: 110,
        depositOwner: depositInitiator,
        lotSize: lotSize }
    ))
    testDeposit = deposits
  })

  describe('undercollateralized', async () => {
    it('unable to start liquidation with sufficient collateral', async () => {
      await expectRevert(
        testDeposit.notifyUndercollateralizedLiquidation(),
        'Deposit has sufficient collateral',
      )
      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)
    })

    it('starts liquidation when bond value falls under collateralization requirement', async () => {
      // price was just on the limit before, since we used minUndercollateralized.
      // increaseing price of sat by just 1 wei, we're over the limit.
      await mockSatWeiPriceFeed.setPrice(satwei.add(new BN(1)))
      await testDeposit.notifyUndercollateralizedLiquidation()

      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it('reverts if no TBTC balance has been approved by the buyer to the Deposit', async () => {
      await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        'Not enough TBTC to cover outstanding debt',
      )
      const depositState = await testDeposit.currentState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it('completes liquidation correctly - end-time period', async () => {
      const duration = await tbtcConstants.getAuctionDuration.call()
      await increaseTime(duration.toNumber())

      await tbtcToken.resetBalance(lotSizeTbtc, { from: auctionBuyer })
      await tbtcToken.approve(testDeposit.address, lotSizeTbtc, { from: auctionBuyer })

      await testDeposit.purchaseSignerBondsAtAuction({ from: auctionBuyer })

      const allowance = await testDeposit.withdrawableAmount({ from: auctionBuyer })
      await testDeposit.withdrawFunds({ from: auctionBuyer })

      const depositState = await testDeposit.currentState.call()
      const endingBalance = await web3.eth.getBalance(testDeposit.address)
      const tbtcBalance = await tbtcToken.balanceOf(depositInitiator)

      expect(tbtcBalance).to.eq.BN(lotSizeTbtc)
      expect(allowance).to.eq.BN(collateralAmount)
      expect(new BN(endingBalance)).to.eq.BN(new BN(0))
      expect(depositState).to.eq.BN(states.LIQUIDATED)
    })
  })
})
