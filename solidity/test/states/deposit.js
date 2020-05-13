// @ts-check
const {BN, expectEvent} = require("@openzeppelin/test-helpers")
const {accounts} = require("@openzeppelin/test-environment")
const {expect} = require("chai")

const {depositRoundTrip} = require("../helpers/utils.js")

/** @typedef {any} BN */
/** @typedef { import("./run.js").StateDefinition<object> } StateDefinition */
/** @typedef { import("./run.js").TruffleReceipt } TruffleReceipt */
/** @typedef { import("./run.js").StateTransitionResolver<any,any> } StateTransitionResolver */
/** @typedef { import("./run.js").StateTransitionResult } StateTransitionResult */
/** @typedef { import("./run.js").StateTransitionResolvers } StateTransitionResolvers */

const opener = accounts[0]
// const redeemer = accounts[1]
// const liquidator = accounts[2]

const System = {
    lotSizes: async ({ TBTCSystem }) => {
        return (await TBTCSystem.getAllowedLotSizes())[0]
    },
    feeEstimate: async ({ TBTCSystem }) => TBTCSystem.getNewDepositFeeEstimate(),
    setEcdsaKey: async ({ ECDSAKeepStub }) => {
        ECDSAKeepStub.setPublicKey(depositRoundTrip.signerPubkey.concatenated)
    },
    setUpBond: async ({ ECDSAKeepStub, bondAmount }) => {
        ECDSAKeepStub.send(bondAmount)
        ECDSAKeepStub.setBondAmount(bondAmount)
    },
    signerSetupTimeout: async ({ TBTCConstants }) => {
        return await TBTCConstants.getSigningGroupFormationTimeout()
    },
    expectedBond: async ({ TBTCSystem, deposit }) => {
        const lotSize = await deposit.lotSizeSatoshis()
        const initial = await TBTCSystem.getInitialCollateralizedPercent()
        return lotSize.mul(await TBTCSystem.fetchBitcoinPrice()).mul(
            initial.div(new BN(100))
        )
    },
    fundingDifficulty: async () => {
        return depositRoundTrip.fundingTx.difficulty
    },
    setDifficulty: async ({ MockRelay, difficulty }) => {
        await MockRelay.setCurrentEpochDifficulty(difficulty)
    },
    setAndApproveRedemptionBalance: async ({ TestTBTCToken, deposit }) => {
        const redemptionRequirement = 
            await deposit.getRedemptionTbtcRequirement(opener)
        await TestTBTCToken.forceMint(
            opener,
            // Let's just play it safe.
            redemptionRequirement,
        )
        await TestTBTCToken.approve(deposit.address, redemptionRequirement, { from: opener })
    },
    States: {
        START: new BN(0),
        AWAITING_SIGNER_SETUP: new BN(1),
        AWAITING_BTC_FUNDING_PROOF: new BN(2),
        FAILED_SETUP: new BN(3),
        ACTIVE: new BN(4),
        AWAITING_WITHDRAWAL_SIGNATURE: new BN(5),
        AWAITING_WITHDRAWAL_PROOF: new BN(6),
        REDEEMED: new BN(7),
        COURTESY_CALL: new BN(8),
        FRAUD_LIQUIDATION_IN_PROGRESS: new BN(9),
        LIQUIDATION_IN_PROGRESS: new BN(10),
        LIQUIDATED: new BN(11),
    }
}

/** @type Object.<string,StateDefinition> */
module.exports = {
    start: {
        name: "start",
        dependencies: {
            lotSize: System.lotSizes,
            feeEstimate: System.feeEstimate,
        },
        next: {
            awaitingSignerSetup: {
                transition: async ({ DepositFactory, lotSize, feeEstimate }) => {
                    return {
                        state: "awaitingSignerSetup",
                        tx: DepositFactory.createDeposit(
                            lotSize,
                            {
                                value: feeEstimate,
                                from: opener,
                            }
                        ),
                        resolveDeposit: ({ Deposit }, receipt) => {
                            const depositAddress =
                                receipt.logs
                                    .find(_ => _.event == "DepositCloneCreated")
                                    .args
                                    .depositCloneAddress

                            return Deposit.at(depositAddress)
                        }
                    }
                },
                expect: async (_, receipt, { TBTCDepositToken, ECDSAKeepStub, deposit }) => {
                    expectEvent(receipt, "DepositCloneCreated", {
                        depositCloneAddress: deposit.address,
                    })
                    expectEvent(receipt, "Created", {
                        _depositContractAddress: deposit.address,
                        _keepAddress: ECDSAKeepStub.address,
                    })
                    expect(await deposit.getCurrentState()).to.eq.BN(
                        System.States.AWAITING_SIGNER_SETUP
                    )
                    expect(await TBTCDepositToken.ownerOf(deposit.address)).to.equal(
                        opener,
                    )
                }
            }
        },
        // TODO What can't happen here?
        failNext: {}
    },
    awaitingSignerSetup: {
        name: "awaitingSignerSetup",
        dependencies: {
            bondAmount: System.expectedBond,
        },
        next: {
            awaitingFundingProof: {
                after: System.generateEcdsaKey,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setEcdsaKey(state)

                    return {
                        state: "awaitingFundingProof",
                        tx: deposit.retrieveSignerPubkey()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "RegisteredPubkey", {
                        _depositContractAddress: deposit.address,
                        _signingGroupPubkeyX: depositRoundTrip.signerPubkey.x,
                        _signingGroupPubkeyY: depositRoundTrip.signerPubkey.y,
                    })
                    expect(await deposit.getCurrentState()).to.eq.BN(
                        System.States.AWAITING_BTC_FUNDING_PROOF
                    )
                },
            },
            signerSetupFailure: {
                after: System.signerSetupTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)

                    return {
                        state: "signerSetupFailure",
                        tx: deposit.notifySignerSetupFailure(),
                    }
                }
            }
        },
        // TODO What can't happen here?
        failNext: {
            "signerSetupFailure too early": {
                transition: async ({ deposit }) => {
                    return {
                        state: "signerSetupFailure",
                        tx: deposit.notifySignerSetupFailure(),
                    }
                }
            }
        },
    },
    awaitingFundingProof: {
        name: "awaitingFundingProof",
        dependencies: {
            bondAmount: System.expectedBond,
            difficulty: System.fundingDifficulty,
        },
        next: {
            active: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setDifficulty(state)

                    return {
                        state: "active",
                        tx: deposit.provideBTCFundingProof(
                            depositRoundTrip.fundingTx.version,
                            depositRoundTrip.fundingTx.txInputVector,
                            depositRoundTrip.fundingTx.txOutputVector,
                            depositRoundTrip.fundingTx.txLocktime,
                            depositRoundTrip.fundingTx.fundingOutputIndex,
                            depositRoundTrip.fundingTx.merkleProof,
                            depositRoundTrip.fundingTx.txIndexInBlock,
                            depositRoundTrip.fundingTx.bitcoinHeaders,
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "Funded")
                    expect(await deposit.getCurrentState()).to.eq.BN(
                        System.States.ACTIVE
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {}
    },
    active: {
        name: "active",
        dependencies: {
        },
        next: {
            awaitingWithdrawalSignature: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setAndApproveRedemptionBalance(state)

                    return {
                        state: "active",
                        tx: deposit.requestRedemption(
                            depositRoundTrip.redemptionTx.outputValueBytes,
                            depositRoundTrip.redemptionTx.outputScript,
                            { from: opener },
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "RedemptionRequested")
                    expect(await deposit.getCurrentState()).to.eq.BN(
                        System.States.AWAITING_WITHDRAWAL_SIGNATURE
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {}
    },
}
