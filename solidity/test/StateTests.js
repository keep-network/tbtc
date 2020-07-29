const states = require("./states/deposit.js")
const {deployAndLinkAll} = require("./helpers/fullDeployer.js")
const {runStatePath} = require("./states/run.js")

describe("tBTC states", () => {
  describe("simple redeem", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "awaitingWithdrawalSignature",
    )
  })
  describe("Courtesy -> Active & liquidation", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "courtesyCall",
      "liquidationInProgress",
    )
  })
  describe("Courtesy -> Active & liquidation_fraud", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "courtesyCall",
      "liquidationInProgress_fraud",
    )
  })
  describe("Undercollateralized Liquidation", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "liquidationInProgress",
    )
  })
  describe("aborted redemption", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "awaitingWithdrawalSignature",
      "liquidationInProgress",
    )
  })
})
