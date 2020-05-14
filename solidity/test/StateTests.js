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
  describe("Courtesy-Active", () => {
    runStatePath(
      states,
      deployAndLinkAll(),
      "start",
      "awaitingSignerSetup",
      "awaitingFundingProof",
      "active",
      "courtesyCall",
    )
  })
})
