// @ts-check
require("mocha")
const {expect} = require("chai")

const {accounts, web3} = require("@openzeppelin/test-environment")
const {BN, expectRevert, time} = require("@openzeppelin/test-helpers")

const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {redemptionPaymentFixturesFor} = require("./helpers/spec-fixtures.js")

const {states, resolveAllLogs, expectEvent} = require("./helpers/utils.js")

/** @typedef { import('bn.js') } BN */
/** @typedef {import("./helpers/spec-fixtures.js").RedemptionPaymentFixture} RedemptionPaymentFixture */

describe("Deposit redemption disbursal tester", async () => {
  // Anchor to run the before hook that actually sets up tests.
  it("should create disbursal tests", async () => {
    expect(true).to.be.true
  })

  before(async () => {
    // System setup.
    const {
      tbtcConstants,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
    } = await deployAndLinkAll()

    const nftId = web3.utils.toBN(testDeposit.address)
    const depositValue = await testDeposit.lotSizeTbtc.call()
    const depositTerm = await tbtcConstants.getDepositTerm.call()

    await testDeposit.setSignerFeeDivisor(new BN("200"))
    /** @type {BN} */
    const signerFee = await testDeposit.signerFeeTbtc.call()
    /** @type {BN} */
    const lotSize = await testDeposit.lotSizeTbtc.call()

    // Test setup.
    describe("DepositRedemption when getting the redemption TBTC requirement", () => {
      before(async () => {
        // Set reference time for pre- vs at-term deposits.
        await testDeposit.setFundingInfo("0x0", await time.latest(), "0x0")
      })

      const fixtures = redemptionPaymentFixturesFor(depositValue, signerFee)
      fixtures.forEach(fixture => {
        const hasFrt = typeof fixture.frtHolder !== "undefined"

        // Below, we define two types of tests:
        // - Baseline tests check situations where the deposit balance has not
        //   been affected outside of the system expectations.
        // - Remaining tests provide a starting balance to the deposit that is
        //   outside of normal system functions. These are divided into three
        //   groups: existing balances that are less than the signer fee,
        //   existing balances that are greater than the signer fee but less
        //   than the lot size, and existing balances that are greater than
        //   the lot size. Actual existing balance values are randomized.
        //
        // In each case, the tests verify both that the internal function that
        // reports redemption TBTC amounts returns the expected values, and
        // that the performRedemptionTBTCTransfers function performs the
        // expected transfers and leaves the deposit with exactly enough
        // escrowed to pay the signers at the end of the redemption flow.
        //
        // The existing-balance tests require certain adjustments to the fixture
        // expectations, which are set up in their respective `describe` blocks.
        describe(`for spec row ${fixture.row}:\t${fixture.description}`, async () => {
          // Set shared state for all tests before each test.
          beforeEach(async () => {
            await createSnapshot()

            // Mint the TDT to the TDT holder's account.
            await tbtcDepositToken.forceMint(accounts[fixture.tdtHolder], nftId)

            // If an amount will be owed by the redeemer, set the redeemer's
            // account up to complete that transfer.
            if (typeof fixture.repaymentAmount !== "undefined") {
              await tbtcToken.forceMint(
                accounts[fixture.redeemer],
                fixture.repaymentAmount,
              )
              await tbtcToken.approve(
                testDeposit.address,
                fixture.repaymentAmount,
                {from: accounts[fixture.redeemer]},
              )
            }

            // If there is an FRT holder, set up the FRT and escrow state.
            if (hasFrt) {
              await feeRebateToken.forceMint(accounts[fixture.frtHolder], nftId)
              await tbtcToken.forceMint(testDeposit.address, signerFee)
            }

            // Set up the at-term or courtesy call state.
            if (fixture.courtesyCall) {
              await testDeposit.setState(states.COURTESY_CALL)
            } else if (!fixture.preTerm) {
              await time.increase(depositTerm)
            }

            await testDeposit.setRedeemerAddress(accounts[fixture.redeemer])
          })

          afterEach(restoreSnapshot)

          it(`should return correct baseline requirements`, async () => {
            await checkTbtcRequirements(
              fixture,
              fixture.repaymentAmount,
              fixture.disbursalAmounts[fixture.tdtHolder],
            )
          })

          it(`should perform the correct TBTC transfers`, async () => {
            await checkTbtcTransfers(
              fixture,
              fixture.repaymentAmount,
              fixture.disbursalAmounts[fixture.tdtHolder],
              fixture.disbursalAmounts[fixture.frtHolder] || new BN(0),
            )
          })

          describe("with additional escrow less than signer fee", async () => {
            const additionalEscrow = randomBnBelow(signerFee)

            let expectedEscrowValue = fixture.repaymentAmount
            if (fixture.repaymentAmount.gt(new BN(0))) {
              // If the redeemer should have _owed_, additional escrow is
              // subtracted from what they should owe, so adjust our
              // expectation accordingly.
              expectedEscrowValue = expectedEscrowValue.sub(additionalEscrow)
            }

            let expectedTdtValue = fixture.disbursalAmounts[fixture.tdtHolder]
            if (fixture.repaymentAmount.lt(additionalEscrow)) {
              // If the redeemer should owe 0 and there is excess escrow,
              // that excess goes to the TDT holder. Similarly, in cases
              // where the redeemer should owe the full lot size, any excess
              // escrow will ultimately go to the TDT holder's balance
              // (since it reduces the fee amount owed directly by the TDT
              // holder).
              expectedTdtValue = additionalEscrow.sub(fixture.repaymentAmount)
            }

            beforeEach(async function() {
              await tbtcToken.forceMint(testDeposit.address, additionalEscrow)
            })

            it(`should return correct requirements`, async () => {
              await checkTbtcRequirements(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
              )
            })

            it(`should perform the correct TBTC transfers`, async () => {
              await checkTbtcTransfers(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
                fixture.disbursalAmounts[fixture.frtHolder] || new BN(0),
              )
            })
          })

          describe("with additional escrow > signer fee, < lot size", async () => {
            const additionalEscrow = signerFee.add(
              randomBnBelow(lotSize.sub(signerFee)),
            )

            let expectedEscrowValue = fixture.repaymentAmount
            if (fixture.repaymentAmount.gt(new BN(0))) {
              // If the redeemer should have _owed_, additional escrow is
              // subtracted from what they should owe, so adjust our expectation
              // accordingly; however, expected escrow cannot dip below 0.
              expectedEscrowValue = expectedEscrowValue.sub(additionalEscrow)
              if (expectedEscrowValue.isNeg()) {
                expectedEscrowValue = new BN(0)
              }
            }

            let expectedTdtValue = fixture.disbursalAmounts[fixture.tdtHolder]
            if (fixture.repaymentAmount.lt(additionalEscrow)) {
              // If the redeemer should owe 0 and there is excess escrow,
              // that excess goes to the TDT holder. Similarly, in cases
              // where the redeemer should owe the full lot size, any excess
              // escrow will ultimately go to the TDT holder's balance
              // (since it reduces the fee amount owed directly by the TDT
              // holder).
              expectedTdtValue = additionalEscrow.sub(fixture.repaymentAmount)
            }

            beforeEach(async function() {
              await tbtcToken.forceMint(testDeposit.address, additionalEscrow)
            })

            it(`should return correct requirements`, async () => {
              await checkTbtcRequirements(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
              )
            })

            it(`should perform the correct TBTC transfers`, async () => {
              await checkTbtcTransfers(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
                fixture.disbursalAmounts[fixture.frtHolder] || new BN(0),
              )
            })
          })

          describe("with additional escrow > lot size", async () => {
            const additionalEscrow = lotSize.add(randomBnBelow(lotSize.muln(2)))

            let expectedEscrowValue = fixture.repaymentAmount
            if (fixture.repaymentAmount.gt(new BN(0))) {
              // If the redeemer should have _owed_, additional escrow is
              // subtracted from what they should owe, so adjust our expectation
              // accordingly; however, expected escrow cannot dip below 0.
              expectedEscrowValue = expectedEscrowValue.sub(additionalEscrow)
              if (expectedEscrowValue.isNeg()) {
                expectedEscrowValue = new BN(0)
              }
            }

            // Baseline expected TDT value. Adjustments are applied below based
            // on the fixture.
            let expectedTdtValue = additionalEscrow

            // If the fixture expects a revert, it generally doesn't handle
            // payouts, so we guard for that. Missing payouts when the
            // redemption shouldn't revert will fail here, which is fine since
            // that's unexpected fixture data and is worth investigating.
            if (!fixture.redemptionShouldRevert) {
              expectedTdtValue = expectedTdtValue.sub(
                fixture.disbursalAmounts["signers"],
              )

              // If there's an FRT holder, earlier setup will escrow that in
              // addition to the additional escrow from these tests, so adjust
              // the expected value accordingly.
              if (typeof fixture.frtHolder !== "undefined") {
                expectedTdtValue = expectedTdtValue.add(signerFee)

                if (fixture.frtHolder != fixture.tdtHolder) {
                  // If the FRT holder and TDT holder are not the same, then the
                  // TDT holder's expected value will not include the FRT
                  // payment.
                  expectedTdtValue = expectedTdtValue.sub(
                    fixture.disbursalAmounts[fixture.frtHolder],
                  )
                } else if (
                  fixture.preTerm &&
                  fixture.frtHolder != fixture.redeemer
                ) {
                  // If the FRT holder and TDT holder *are* the same, the FRT
                  // holder payout in the fixture includes the full TDT holder
                  // payout for redemption. These are handled separately by the
                  // redemption payout flow, however, so the expected TDT value
                  // has to be adjusted for the FRT payout.
                  expectedTdtValue = expectedTdtValue.sub(signerFee)
                }
              }
            }

            beforeEach(async function() {
              await tbtcToken.forceMint(testDeposit.address, additionalEscrow)
            })

            it(`should return correct requirements`, async () => {
              await checkTbtcRequirements(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
              )
            })

            it(`should perform the correct TBTC transfers`, async () => {
              await checkTbtcTransfers(
                fixture,
                expectedEscrowValue,
                expectedTdtValue,
                fixture.disbursalAmounts[fixture.frtHolder] || new BN(0),
              )
            })
          })
        })
      })
    })

    /**
     * Helper function that runs all checks related to
     * calculateRedemptionTbtcAmounts, given the expected redeemer payment.
     *
     * @param {RedemptionPaymentFixture} fixture
     * @param {BN} expectedRedeemerPayment
     * @param {BN} expectedTdtHolderValue
     */
    async function checkTbtcRequirements(
      fixture,
      expectedRedeemerPayment,
      expectedTdtHolderValue,
    ) {
      if (fixture.redemptionShouldRevert) {
        expectRevert.unspecified(
          testDeposit.getRedemptionTbtcRequirement(accounts[fixture.redeemer]),
        )
      } else {
        // Unpack return values accordingly.
        const {
          "0": escrowValue,
          "1": tdtHolderValue,
          "2": frtHolderValue,
        } = await testDeposit.calculateRedemptionTbtcAmounts(
          accounts[fixture.redeemer],
        )

        expect(escrowValue, "Unexpected redeemer payment requirement").to.eq.BN(
          expectedRedeemerPayment,
        )

        expect(tdtHolderValue, "Unexpected amount owed to TDT holder").to.eq.BN(
          expectedTdtHolderValue,
        )

        if (fixture.frtHolder != fixture.tdtHolder) {
          expect(
            frtHolderValue,
            "Unexpected amount owed to FRT holder",
          ).to.eq.BN(fixture.disbursalAmounts[fixture.frtHolder])
        } else {
          expect(
            frtHolderValue,
            "Unexpected amount owed to FRT holder",
          ).to.eq.BN(0)
        }
      }
    }

    /**
     * Helper function that runs all checks related to
     * performRedemptionTBTCTransfers, given the expected redeemer payment.
     *
     * @param {RedemptionPaymentFixture} fixture
     * @param {BN} expectedRedeemerPayment
     * @param {BN} expectedTdtHolderValue
     * @param {BN} expectedFrtHolderValue
     */
    async function checkTbtcTransfers(
      fixture,
      expectedRedeemerPayment,
      expectedTdtHolderValue,
      expectedFrtHolderValue,
    ) {
      if (fixture.redemptionShouldRevert) {
        expectRevert.unspecified(
          testDeposit.performRedemptionTBTCTransfers({
            from: accounts[fixture.redeemer],
          }),
        )
      } else {
        const receipt = resolveAllLogs(
          (
            await testDeposit.performRedemptionTBTCTransfers({
              from: accounts[fixture.redeemer],
            })
          ).receipt,
          {tbtcToken},
        )

        // Escrow should only hold the signer fee during redemption.
        expect(
          await tbtcToken.balanceOf(testDeposit.address),
          "Deposit does not have enough escrowed to pay signers",
        ).to.eq.BN(signerFee)

        if (expectedRedeemerPayment.gt(new BN(0))) {
          expectEvent(receipt, "Transfer", {
            from: accounts[fixture.redeemer],
            to: testDeposit.address,
            value: expectedRedeemerPayment,
          })
        }
        // else check there is no transfer event

        if (
          expectedFrtHolderValue.gt(new BN(0)) &&
          // If the TDT and FRT holder are the same, the TDT holder check
          // should win.
          fixture.tdtHolder != fixture.frtHolder
        ) {
          expectEvent(receipt, "Transfer", {
            from: testDeposit.address,
            to: accounts[fixture.frtHolder],
            value: expectedFrtHolderValue,
          })
        }

        if (expectedTdtHolderValue.gt(new BN(0))) {
          expectEvent(receipt, "Transfer", {
            from: testDeposit.address,
            to: accounts[fixture.tdtHolder],
            value: expectedTdtHolderValue,
          })
        }
        // else check there is no transfer event
      }
    }
  })
})

/**
 * @param {BN} maxValue The max value to generate.
 * @return {BN} A random BN between 0 and the given max value. Randomness is
 *         limited to the first 6 significant digits.
 */
function randomBnBelow(maxValue) {
  // Scale the random number up to something integer-y, then
  // divide the scale factor away to get the right range.
  return new BN(Math.random() * 1000000).mul(maxValue).divn(1000000)
}
