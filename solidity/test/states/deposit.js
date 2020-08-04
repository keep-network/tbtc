// @ts-check
const {BN} = require("@openzeppelin/test-helpers")
const {accounts} = require("@openzeppelin/test-environment")
const {expect} = require("chai")

const {depositRoundTrip, bytes32zero, expectEvent} = require("../helpers/utils.js")

/** @typedef {any} BN */
/** @typedef { import("./run.js").StateDefinition<object> } StateDefinition */
/** @typedef { import("./run.js").TruffleReceipt } TruffleReceipt */
/** @typedef { import("./run.js").StateTransitionResolver<any,any> } StateTransitionResolver */
/** @typedef { import("./run.js").StateTransitionResult } StateTransitionResult */
/** @typedef { import("./run.js").StateTransitionResolvers } StateTransitionResolvers */

const opener = accounts[0]
// const redeemer = accounts[1]
const liquidator = accounts[2]

const System = {
    lotSizes: async ({ TBTCSystem }) => {
        return await TBTCSystem.getAllowedLotSizes()
    },
    feeEstimate: async ({ TBTCSystem }) => TBTCSystem.getNewDepositFeeEstimate(),
    setEcdsaKey: async ({ ECDSAKeepStub }) => {
        ECDSAKeepStub.setPublicKey(depositRoundTrip.signerPubkey.concatenated)
    },
    courtesyTimeout: async ({ TBTCConstants }) => {
        return await TBTCConstants.getCourtesyCallTimeout()
    },
    fundingTimeout: async ({ TBTCConstants }) => {
        return await TBTCConstants.getFundingTimeout()
    },
    signatureTimeout: async ({ TBTCConstants }) => {
        return await TBTCConstants.getSignatureTimeout()
    },
    setUpBond: async ({ ECDSAKeepStub, bondAmount }) => {
        ECDSAKeepStub.send(bondAmount)
        ECDSAKeepStub.setBondAmount(bondAmount)
    },
    drainBond: async ({ ECDSAKeepStub }) => {
        ECDSAKeepStub.drain()
        ECDSAKeepStub.setBondAmount(0)
    },
    signerSetupTimeout: async ({ TBTCConstants }) => {
        return await TBTCConstants.getSigningGroupFormationTimeout()
    },
    beforeSignerSetupTimeout: async (state) => {
        return (await System.signerSetupTimeout(state)).sub(new BN(10))
    },
    expectedBond: async ({ TBTCSystem, deposit }) => {
        const lotSize = await deposit.lotSizeSatoshis()
        const initial = await TBTCSystem.getInitialCollateralizedPercent()
        const price = await TBTCSystem.fetchBitcoinPrice({ from: deposit.address })
        return price.mul(lotSize).mul(initial).div(new BN(100))
    },
    fundingDifficulty: async () => {
        return depositRoundTrip.fundingTx.difficulty
    },
    setDifficulty: async ({ MockRelay, difficulty }) => {
        await MockRelay.setCurrentEpochDifficulty(difficulty)
    },
    depositLotSize: async ({ deposit }) => {
        return await deposit.lotSizeTbtc();
    },
    depositSignerFee: async ({ deposit }) => {
        return await deposit.signerFeeTbtc();
    },
    depositRedemptionRequirement: async ({ deposit }) => {
        return await deposit.getRedemptionTbtcRequirement(opener)
    },
    depositVendingMachineRedemptionRequirement: async ({ deposit, lotSize }) => {
        return lotSize.add(await deposit.getOwnerRedemptionTbtcRequirement(opener))
    },
    setAndApproveRedemptionBalance: async ({ TestTBTCToken, deposit, redemptionRequirement }) => {
        await TestTBTCToken.forceMint(
            opener,
            // Let's just play it safe.
            redemptionRequirement,
        )
        await TestTBTCToken.approve(
            deposit.address,
            redemptionRequirement,
            { from: opener }
        )
    },
    setAndApproveVendingMachineRedemptionBalance: async ({ VendingMachine, TestTBTCToken, redemptionRequirement }) => {
        await TestTBTCToken.forceMint(
            opener,
            // Let's just play it safe.
            redemptionRequirement,
        )
        await TestTBTCToken.approve(
            VendingMachine.address,
            redemptionRequirement,
            { from: opener }
        )
    },
    setAndApproveLiquidationBalance: async ({ TestTBTCToken, deposit }) => {
        const lotSize = await deposit.lotSizeSatoshis()
        const requirement = lotSize.mul(new BN(10).pow(new BN(10)))
        await TestTBTCToken.forceMint(liquidator, requirement)
        await TestTBTCToken.approve(deposit.address, requirement, { from: liquidator })
    },
    setWellCollateralized: async ({ deposit, ECDSAKeepStub, mockSatWeiPriceFeed }) => {
        const bondAmount = await ECDSAKeepStub.checkBondAmount()
        const lotSize = await deposit.lotSizeSatoshis()
        const initial = await deposit.initialCollateralizedPercent()
        const wellCollateralized = bondAmount.div(lotSize.mul(initial).div(new BN(100)))
        await mockSatWeiPriceFeed.setPrice(wellCollateralized)
    },
    setUndercollateralized: async ({ deposit, ECDSAKeepStub, mockSatWeiPriceFeed }) => {
        const bondAmount = await ECDSAKeepStub.checkBondAmount()
        const lotSize = await deposit.lotSizeSatoshis()
        const under = await deposit.undercollateralizedThresholdPercent()
        const undercollateralized = bondAmount.div(lotSize.mul(under.sub(new BN(1))).div(new BN(100)))
        await mockSatWeiPriceFeed.setPrice(undercollateralized)
    },
    setSeverelyUndercollateralized: async ({ deposit, ECDSAKeepStub, mockSatWeiPriceFeed }) => {
        const bondAmount = await ECDSAKeepStub.checkBondAmount()
        const lotSize = await deposit.lotSizeSatoshis()
        const severe = await deposit.severelyUndercollateralizedThresholdPercent()
        const severelyUndercollateralized = bondAmount.div(lotSize.mul(severe.sub(new BN(1))).div(new BN(100)))
        await mockSatWeiPriceFeed.setPrice(severelyUndercollateralized)
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
            availableLotSizes: System.lotSizes,
            feeEstimate: System.feeEstimate,
        },
        next: {
            awaitingSignerSetup: {
                transition: async ({ DepositFactory, availableLotSizes, feeEstimate }) => {
                    return {
                        state: "awaitingSignerSetup",
                        tx: DepositFactory.createDeposit(
                            availableLotSizes[0],
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
                    expect(await deposit.currentState()).to.eq.BN(
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
                    expect(await deposit.currentState()).to.eq.BN(
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
                        tx: deposit.notifySignerSetupFailed(),
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "SetupFailed")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FAILED_SETUP
                    )
                }
            }
        },
        failNext: {
            "signerSetupFailure too early": {
                dependencies: {
                    bondAmount: System.expectedBond,
                },
                after: System.beforeSignerSetupTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)

                    return {
                        state: "signerSetupFailure",
                        tx: deposit.notifySignerSetupFailed(),
                    }
                },
                expectError: async (_, error) => {
                    expect(error.message).to.match(
                        /Signing group formation timeout not yet elapsed/
                    )
                }
            },
            "signerSetupFailure when no bond is available": {
                dependencies: {
                    bondAmount: System.expectedBond,
                },
                after: System.signerSetupTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.drainBond(state)
                    return {
                        state: "signerSetupFailure",
                        tx: deposit.notifySignerSetupFailed(),
                    }
                },
                expectError: async (_, error) => {
                    expect(error.message).to.match(
                        /No funds received, unexpected/
                    )
                }
            },
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
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.ACTIVE
                    )
                },
            },
            failedSetup_timeout: {
                after: System.fundingTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setDifficulty(state)

                    return {
                        state: "failedSetup",
                        tx: deposit.notifyFundingTimedOut()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "SetupFailed")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FAILED_SETUP
                    )
                },
            },
            failedSetup_ECDSAFraud: {
                after: System.fundingTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "fraudLiquidation",
                        tx: deposit.provideFundingECDSAFraudProof(
                            0,
                            bytes32zero,
                            bytes32zero,
                            bytes32zero,
                            "0x00"
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "SetupFailed")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FAILED_SETUP
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {
            "signerSetupFailure after signer group formation": {
                transition: async ({ deposit }) => {
                    return {
                        state: "signerSetupFailure",
                        tx: deposit.notifySignerSetupFailed(),
                    }
                },
                expectError: async (_, error) => {
                    expect(error.message).to.match(/Not awaiting setup/)
                }
            },
        }
    },
    active: {
        name: "active",
        dependencies: {
            bondAmount: System.expectedBond,
            redemptionRequirement: System.depositRedemptionRequirement,
            lotSize: System.depositLotSize,
            signerFee: System.depositSignerFee,
        },
        next: {
            minted: {
                transition: async (state) => {
                    const { deposit, VendingMachine, TBTCDepositToken } = state
                    await TBTCDepositToken.approve(
                        VendingMachine.address,
                        deposit.address,
                        { from: opener },
                    )

                    return {
                        state: "minted",
                        tx: VendingMachine.tdtToTbtc(deposit.address, { from: opener })
                    }
                },
                expect: async (_, receipt, { deposit, VendingMachine, lotSize, signerFee }) => {
                    // Owner should have received the TBTC less signer fee.
                    expectEvent(
                        receipt,
                        "Transfer",
                        {
                            "from": "0x0000000000000000000000000000000000000000",
                            "to": opener,
                            "value": lotSize.sub(signerFee),
                        },
                    )
                    // Owner should have received the FRT.
                    expectEvent(
                        receipt,
                        "Transfer",
                        {
                            "from": "0x0000000000000000000000000000000000000000",
                            "to": opener,
                            "tokenId": new BN(deposit.address.substring(2), 16),
                        },
                    )
                    // Vending machine should have received the TDT.
                    expectEvent(
                        receipt,
                        "Transfer",
                        {
                            "from": opener,
                            "to": VendingMachine.address,
                            "tokenId": new BN(deposit.address.substring(2), 16),
                        },
                    )
                    // Deposit should have received the signer fee.
                    expectEvent(
                        receipt,
                        "Transfer",
                        {
                            "from": "0x0000000000000000000000000000000000000000",
                            "to": deposit.address,
                            "value": signerFee,
                        },
                    )
                },
            },
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
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.AWAITING_WITHDRAWAL_SIGNATURE
                    )
                },
            },
            courtesyCall: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    await System.setUndercollateralized(state)

                    return {
                        state: "courtesyCall",
                        tx: deposit.notifyCourtesyCall()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "CourtesyCalled")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.COURTESY_CALL
                    )
                },
            },
            liquidationInProgress: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    await System.setSeverelyUndercollateralized(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.notifyUndercollateralizedLiquidation()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATION_IN_PROGRESS
                    )
                },
            },
            liquidationInProgress_fraud: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.provideECDSAFraudProof(
                            0,
                            bytes32zero,
                            bytes32zero,
                            bytes32zero,
                            "0x00"
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FRAUD_LIQUIDATION_IN_PROGRESS
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {}
    },
    minted: {
        name: "minted",
        dependencies: {
            redemptionRequirement: System.depositVendingMachineRedemptionRequirement,
        },
        next: {
            awaitingWithdrawalSignature: {
                transition: async (state) => {
                    const { deposit, VendingMachine } = state
                    await System.setAndApproveVendingMachineRedemptionBalance(state)

                    return {
                        state: "active",
                        tx: VendingMachine.tbtcToBtc(
                            deposit.address,
                            depositRoundTrip.redemptionTx.outputValueBytes,
                            depositRoundTrip.redemptionTx.outputScript,
                            { from: opener },
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "RedemptionRequested")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.AWAITING_WITHDRAWAL_SIGNATURE
                    )
                },
            },
            courtesyCall: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    await System.setUndercollateralized(state)

                    return {
                        state: "courtesyCall",
                        tx: deposit.notifyCourtesyCall()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "CourtesyCalled")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.COURTESY_CALL
                    )
                },
            },
            liquidationInProgress: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    await System.setSeverelyUndercollateralized(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.notifyUndercollateralizedLiquidation()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATION_IN_PROGRESS
                    )
                },
            },
            liquidationInProgress_fraud: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.provideECDSAFraudProof(
                            0,
                            bytes32zero,
                            bytes32zero,
                            bytes32zero,
                            "0x00"
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FRAUD_LIQUIDATION_IN_PROGRESS
                    )
                },
            },
        }
    },
    awaitingWithdrawalSignature: {
        name: "awaitingWithdrawalSignature",
        dependencies: {
        },
        next: {
            redeemed: {
                transition: async (state) => {
                    const { deposit } = state

                    return {
                        state: "redeemed",
                        tx: deposit.provideRedemptionProof(
                            depositRoundTrip.redemptionTx.version,
                            depositRoundTrip.redemptionTx.txInputVector,
                            depositRoundTrip.redemptionTx.txOutputVector,
                            depositRoundTrip.redemptionTx.txLocktime,
                            depositRoundTrip.redemptionTx.merkleProof,
                            depositRoundTrip.redemptionTx.txIndexInBlock,
                            depositRoundTrip.redemptionTx.bitcoinHeaders,
                          )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "Redeemed")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.REDEEMED
                    )
                },
            },
            liquidated: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "liquidated",
                        tx: deposit.provideECDSAFraudProof(
                            0,
                            bytes32zero,
                            bytes32zero,
                            bytes32zero,
                            "0x00",
                        ),
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(
                        receipt,
                        "StartedLiquidation",
                        {
                            "_depositContractAddress": deposit.address,
                        },
                    )
                    expectEvent(
                        receipt,
                        "Liquidated",
                        {
                            "_depositContractAddress": deposit.address,
                        },
                    )
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATED
                    )
                },
            },
            liquidationInProgress: {
                after: System.signatureTimeout,
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.notifyRedemptionSignatureTimedOut()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATION_IN_PROGRESS
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {}
    },
    courtesyCall: {
        name: "courtesyCall",
        dependencies: {
            bondAmount: System.expectedBond,
            redemptionRequirement: System.depositRedemptionRequirement,
        },
        next: {
            active: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    await System.setWellCollateralized(state)
                    return {
                        state: "active",
                        tx: deposit.exitCourtesyCall()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "ExitedCourtesyCall")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.ACTIVE
                    )
                },
            },
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
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.AWAITING_WITHDRAWAL_SIGNATURE
                    )
                },
            },
            liquidationInProgress: {
                after: System.courtesyTimeout,
                transition: async (state) => {
                    const { deposit } = state

                    return {
                        state: "liquidationInProgress",
                        tx: deposit.notifyCourtesyCallExpired()
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATION_IN_PROGRESS
                    )
                },
            },
            liquidationInProgress_fraud: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setUpBond(state)
                    return {
                        state: "liquidatinInProgress",
                        tx: deposit.provideECDSAFraudProof(
                            0,
                            bytes32zero,
                            bytes32zero,
                            bytes32zero,
                            "0x00"
                        )
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "StartedLiquidation")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.FRAUD_LIQUIDATION_IN_PROGRESS
                    )
                },
            },
        },
        // TODO What can't happen here?
        failNext: {}
    },
    liquidationInProgress: {
        name: "liquidationInProgress",
        dependencies: {
            bondAmount: System.expectedBond,
        },
        next: {
            liquidated: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setAndApproveLiquidationBalance(state)
                    return {
                        state: "liquidated",
                        tx: deposit.purchaseSignerBondsAtAuction({from: liquidator})
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "Liquidated")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATED
                    )
                },
            },
        },
        failNext: {
            "active after liquidation started": {
                transition: async ({ deposit }) => {
                    return {
                        state: "active",
                        tx: deposit.exitCourtesyCall(),
                    }
                },
                expectError: async (_, error) => {
                    expect(error.message).to.match(/Not currently in courtesy call./)
                }
            },
        }
    },
    liquidationInProgress_fraud: {
        name: "liquidationInProgress_fraud",
        dependencies: {
            bondAmount: System.expectedBond,
        },
        next: {
            liquidated: {
                transition: async (state) => {
                    const { deposit } = state
                    await System.setAndApproveLiquidationBalance(state)
                    return {
                        state: "liquidated",
                        tx: deposit.purchaseSignerBondsAtAuction({from: liquidator})
                    }
                },
                expect: async (_, receipt, { deposit }) => {
                    expectEvent(receipt, "Liquidated")
                    expect(await deposit.currentState()).to.eq.BN(
                        System.States.LIQUIDATED
                    )
                },
            },
        },
        failNext: {
            "active after fraud liquidation started": {
                transition: async ({ deposit }) => {
                    return {
                        state: "active",
                        tx: deposit.exitCourtesyCall(),
                    }
                },
                expectError: async (_, error) => {
                    expect(error.message).to.match(/Not currently in courtesy call./)
                }
            },
        }
    }
}
