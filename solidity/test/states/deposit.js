// @ts-check
const {BN, expectEvent} = require("@openzeppelin/test-helpers")
const {accounts} = require("@openzeppelin/test-environment")
const {expect} = require("chai")

/** @typedef {any} BN */
/** @typedef { import("./run.js").StateDefinition<object> } StateDefinition */
/** @typedef { import("./run.js").TruffleReceipt } TruffleReceipt */
/** @typedef { import("./run.js").StateTransitionResolver<any,any> } StateTransitionResolver */
/** @typedef { import("./run.js").StateTransitionResult } StateTransitionResult */
/** @typedef { import("./run.js").StateTransitionResolvers } StateTransitionResolvers */

const publicKey =
    "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"
const publicKeyX =
    "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
const publicKeyY =
    "0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"

const System = {
    lotSizes: async ({ TBTCSystem }) => {
        return (await TBTCSystem.getAllowedLotSizes())[0]
    },
    feeEstimate: async ({ TBTCSystem }) => TBTCSystem.getNewDepositFeeEstimate(),
    setEcdsaKey: async ({ ECDSAKeepStub }) => {
        ECDSAKeepStub.setPublicKey(publicKey)
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
    States: {
        START: new BN(0),
        AWAITING_SIGNER_SETUP: new BN(1),
        AWAITING_BTC_FUNDING_PROOF: new BN(2),
    }
}

const opener = accounts[0],
      redeemer = accounts[1],
      liquidator = accounts[2]

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
                        /** @type TruffleReceipt */
                        tx: DepositFactory.createDeposit(lotSize, { value: feeEstimate }),
                        /** @type StateTransitionResolver */
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
                expect: async (_, receipt, { ECDSAKeepStub, deposit }) => {
                    expectEvent(receipt, "DepositCloneCreated", {
                        depositCloneAddress: deposit.address,
                    })
                    expectEvent(receipt, "Created", {
                        _depositContractAddress: deposit.address,
                        _keepAddress: ECDSAKeepStub.address,
                    })
                    expect(await deposit.getCurrentState()).to.eq.BN(System.States.AWAITING_SIGNER_SETUP)
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
                        _signingGroupPubkeyX: publicKeyX,
                        _signingGroupPubkeyY: publicKeyY,
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
    }
}
