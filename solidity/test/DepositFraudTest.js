const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {states, bytes32zero} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {AssertBalance} = require("./helpers/assertBalance.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, expectRevert} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

describe("DepositFraud", async function() {
  let tbtcConstants
  let tbtcSystemStub
  let tbtcDepositToken
  let testDeposit
  let ecdsaKeepStub
  let beneficiary
  let assertBalance

  before(async () => {
    ;({
      tbtcConstants,
      tbtcToken,
      tbtcSystemStub,
      tbtcDepositToken,
      testDeposit,
      ecdsaKeepStub,
    } = await deployAndLinkAll())

    beneficiary = accounts[4]
    assertBalance = new AssertBalance(tbtcToken)
    tbtcDepositToken.forceMint(
      beneficiary,
      web3.utils.toBN(testDeposit.address),
    )
  })

  beforeEach(async () => {
    await testDeposit.reset()
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  describe("provideFundingECDSAFraudProof", async () => {
    const bond = 1000

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      await ecdsaKeepStub.setSuccess(true)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await ecdsaKeepStub.send(bond, {from: owner})
    })

    it("updates to awaiting fraud funding proof, distributes signer bond to funder, and logs FraudDuringSetup", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.provideFundingECDSAFraudProof(
        0,
        bytes32zero,
        bytes32zero,
        bytes32zero,
        "0x00",
      )

      await assertBalance.eth(ecdsaKeepStub.address, new BN(0))
      await assertBalance.eth(testDeposit.address, new BN(bond))

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: beneficiary,
      })
      expect(withdrawable).to.eq.BN(new BN(bond))

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.FAILED_SETUP)

      const eventList = await tbtcSystemStub.getPastEvents("FraudDuringSetup", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length, "bad event length").to.equal(1)
    })

    it("reverts if not awaiting funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideFundingECDSAFraudProof(
          0,
          bytes32zero,
          bytes32zero,
          bytes32zero,
          "0x00",
        ),
        "Signer fraud during funding flow only available while awaiting funding",
      )
    })

    it("reverts if the signature is not fraud", async () => {
      await ecdsaKeepStub.setSuccess(false)

      await expectRevert(
        testDeposit.provideFundingECDSAFraudProof(
          0,
          bytes32zero,
          bytes32zero,
          bytes32zero,
          "0x00",
        ),
        "Signature is not fraudulent",
      )
    })
  })

  describe("startLiquidation - fraud", async () => {
    let signerBond
    before(async () => {
      signerBond = 10000000
      await ecdsaKeepStub.send(signerBond, {from: accounts[1]})
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("executes and emits StartedLiquidation event", async () => {
      const {
        receipt: {blockNumber: liquidationBlock},
      } = await testDeposit.startLiquidation(true, {from: owner})
      const block = await web3.eth.getBlock(liquidationBlock)

      const events = await tbtcSystemStub.getPastEvents("StartedLiquidation", {
        fromBlock: block.number,
        toBlock: "latest",
      })

      const initiator = await testDeposit.getLiquidationInitiator()
      const initiated = await testDeposit.getLiquidationTimestamp()

      expect(events[0].returnValues[0]).to.equal(testDeposit.address)
      expect(events[0].returnValues[1]).to.be.true

      expect(initiator).to.equal(owner)
      expect(initiated).to.eq.BN(block.timestamp)
    })

    it("liquidates immediately with bonds going to the redeemer if we came from the redemption flow", async () => {
      // setting redeemer address suggests we are coming from redemption flow
      testDeposit.setRedeemerAddress(owner)
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_SIGNATURE)

      const currentBond = await web3.eth.getBalance(ecdsaKeepStub.address)
      const {
        receipt: {blockNumber: liquidationBlock},
      } = await testDeposit.startLiquidation(true)
      const block = await web3.eth.getBlock(liquidationBlock)

      const events = await tbtcSystemStub.getPastEvents("Liquidated", {
        fromBlock: block.number,
        toBlock: "latest",
      })

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: owner,
      })
      expect(withdrawable).to.eq.BN(new BN(currentBond))

      expect(events[0].returnValues[0]).to.equal(testDeposit.address)

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATED)
    })
  })

  describe("startLiquidation - not fraud", async () => {
    let signerBond
    before(async () => {
      signerBond = 10000000
      await ecdsaKeepStub.send(signerBond, {from: accounts[1]})
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("executes and emits StartedLiquidation event", async () => {
      const {
        receipt: {blockNumber: liquidationBlock},
      } = await testDeposit.startLiquidation(false, {from: owner})
      const block = await web3.eth.getBlock(liquidationBlock)

      const events = await tbtcSystemStub.getPastEvents("StartedLiquidation", {
        fromBlock: block.number,
        toBlock: "latest",
      })
      const initiator = await testDeposit.getLiquidationInitiator()
      const initiated = await testDeposit.getLiquidationTimestamp()
      expect(events[0].returnValues[0]).to.equal(testDeposit.address)
      expect(events[0].returnValues[1]).to.be.false

      expect(initiator).to.equal(owner)
      expect(initiated).to.eq.BN(block.timestamp)
    })
  })

  describe("provideECDSAFraudProof", async () => {
    before(async () => {
      await testDeposit.setState(states.ACTIVE)
      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("executes and moves state to FRAUD_LIQUIDATION_IN_PROGRESS", async () => {
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

    it("reverts if in the funding flow", async () => {
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)

      await expectRevert(
        testDeposit.provideECDSAFraudProof(
          0,
          bytes32zero,
          bytes32zero,
          bytes32zero,
          "0x00",
        ),
        "Use provideFundingECDSAFraudProof instead",
      )
    })

    it("reverts if already in signer liquidation", async () => {
      await testDeposit.setState(states.LIQUIDATION_IN_PROGRESS)

      await expectRevert(
        testDeposit.provideECDSAFraudProof(
          0,
          bytes32zero,
          bytes32zero,
          bytes32zero,
          "0x00",
        ),
        "Signer liquidation already in progress",
      )
    })

    it("reverts if the contract has halted", async () => {
      await testDeposit.setState(states.REDEEMED)

      await expectRevert(
        testDeposit.provideECDSAFraudProof(
          0,
          bytes32zero,
          bytes32zero,
          bytes32zero,
          "0x00",
        ),
        "Contract has halted",
      )
    })

    it("reverts if signature is not fraud according to Keep", async () => {
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
    })
  })
})
