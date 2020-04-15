const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {states, bytes32zero} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {AssertBalance} = require("./helpers/assertBalance.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
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
      const blockNumber = await web3.eth.getBlock("latest").number

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

  describe("startSignerFraudLiquidation", async () => {
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
      const block = await web3.eth.getBlock("latest")

      await testDeposit.startSignerFraudLiquidation({from: owner})

      const events = await tbtcSystemStub.getPastEvents("StartedLiquidation", {
        fromBlock: block.number,
        toBlock: "latest",
      })

      const initiator = await testDeposit.getLiquidationInitiator()
      const initiated = await testDeposit.getLiquidationTimestamp()

      expect(events[0].returnValues[0]).to.equal(testDeposit.address)
      expect(events[0].returnValues[1]).to.be.true
      expect(events[0].returnValues[2]).to.eq.BN(block.timestamp)

      expect(initiator).to.equal(owner)
      expect(initiated).to.eq.BN(block.timestamp)
    })

    it("liquidates immediately with bonds going to the redeemer if we came from the redemption flow", async () => {
      // setting redeemer address suggests we are coming from redemption flow
      testDeposit.setRedeemerAddress(owner)

      const currentBond = await web3.eth.getBalance(ecdsaKeepStub.address)
      const block = await web3.eth.getBlock("latest")
      await testDeposit.startSignerFraudLiquidation()

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

  describe("startSignerAbortLiquidation", async () => {
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
      const block = await web3.eth.getBlock("latest")

      await testDeposit.startSignerAbortLiquidation({from: owner})

      const events = await tbtcSystemStub.getPastEvents("StartedLiquidation", {
        fromBlock: block.number,
        toBlock: "latest",
      })
      const initiator = await testDeposit.getLiquidationInitiator()
      const initiated = await testDeposit.getLiquidationTimestamp()
      expect(events[0].returnValues[0]).to.equal(testDeposit.address)
      expect(events[0].returnValues[1]).to.be.false
      expect(events[0].returnValues[2]).to.eq.BN(block.timestamp)

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

  describe("validateRedeemerNotPaid", async () => {
    const _txOutputVector =
      "0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6"
    const redeemerOutputScript =
      "0x16001486e7303082a6a21d5837176bc808bf4828371ab6"
    const prevoutValueBytes = "0xf078351d00000000"
    const outpoint =
      "0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000"
    const _longTxOutputVector = `0x034897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b078952112040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6`

    beforeEach(async () => {
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        2424,
        0,
        bytes32zero,
      )
    })

    it("returns false if redeemer is paid and value is sufficient", async () => {
      const success = await testDeposit.validateRedeemerNotPaid(_txOutputVector)
      expect(success).to.be.false
    })

    it("returns false if redeemer is paid, value is sufficient and output is at 3rd position", async () => {
      const success = await testDeposit.validateRedeemerNotPaid(
        _longTxOutputVector,
      )
      expect(success).to.be.false
    })

    it("returns true if redeemer is not paid", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        "0x" + "0".repeat(20),
        2424,
        0,
        bytes32zero,
      )

      const success = await testDeposit.validateRedeemerNotPaid(_txOutputVector)
      expect(success).to.be.true
    })

    it("returns true if value is not sufficient", async () => {
      await testDeposit.setUTXOInfo("0xf078351d00000001", 0, outpoint)

      const success = await testDeposit.validateRedeemerNotPaid(_txOutputVector)
      expect(success).to.be.true
    })

    it("returns true if there is no witness flag", async () => {
      const _txOutputVectorNoWitness =
        "0x024897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211"
      const newPKH =
        "0xa4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c18" // note length > 20

      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        newPKH,
        2424,
        0,
        bytes32zero,
      )
      await testDeposit.setUTXOInfo("0xffff", 0, outpoint)

      const success = await testDeposit.validateRedeemerNotPaid(
        _txOutputVectorNoWitness,
      )
      expect(success).to.be.true
    })
  })
})
