const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {states, bytes32zero} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {AssertBalance} = require("./helpers/assertBalance.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

const currentDifficulty = 6353030562983
const _version = "0x01000000"
const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
const _txOutputVector =
  "0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6"
const _fundingOutputIndex = 0
const _txLocktime = "0x4ec10800"
const _txIndexInBlock = 129
const _bitcoinHeaders =
  "0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e"
const _signerPubkeyX =
  "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e"
const _signerPubkeyY =
  "0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1"
const _merkleProof =
  "0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049"

describe("DepositFraud", async function() {
  let tbtcConstants
  let mockRelay
  let tbtcSystemStub
  let tbtcDepositToken
  let testDeposit
  let ecdsaKeepStub
  let beneficiary
  let fundingProofTimerStart
  let assertBalance

  before(async () => {
    ;({
      tbtcConstants,
      tbtcToken,
      mockRelay,
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
    let timer

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1 // has elapsed
      await ecdsaKeepStub.setSuccess(true)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)

      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    it("updates to awaiting fraud funding proof and logs FraudDuringSetup if the timer has not elapsed", async () => {
      const blockNumber = await web3.eth.getBlock("latest").number

      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart * 5) // timer has not elapsed

      await testDeposit.provideFundingECDSAFraudProof(
        0,
        bytes32zero,
        bytes32zero,
        bytes32zero,
        "0x00",
      )

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.FRAUD_AWAITING_BTC_FUNDING_PROOF)

      const actualFundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        actualFundingProofTimerStart,
        "fundingProofTimerStart did not increase",
      ).to.be.bignumber.greaterThan(new BN(fundingProofTimerStart))

      const eventList = await tbtcSystemStub.getPastEvents("FraudDuringSetup", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length, "bad event length").to.equal(1)
    })

    it("updates to failed, logs FraudDuringSetup, and burns value if the timer has elapsed", async () => {
      const blockNumber = await web3.eth.getBlock("latest").number
      await testDeposit.provideFundingECDSAFraudProof(
        0,
        bytes32zero,
        bytes32zero,
        bytes32zero,
        "0x00",
      )

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

  describe("notifyFraudFundingTimeout", async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getFraudFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1 // timer has elapsed

      await testDeposit.setState(states.FRAUD_AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)

      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    it("updates state to setup failed, logs SetupFailed, and deconstes state", async () => {
      const blockNumber = await web3.eth.getBlock("latest").number

      await testDeposit.notifyFraudFundingTimeout()

      const keepAddress = await testDeposit.getKeepAddress.call()
      expect(keepAddress, "Keep address not deconsted").to.equal(ZERO_ADDRESS)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt not deconsted",
      ).to.be.bignumber.equal(new BN(0))

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "fundingProofTimerStart not deconsted",
      ).to.be.bignumber.equal(new BN(0))

      const signingGroupPublicKey = await testDeposit.getSigningGroupPublicKey.call()
      expect(signingGroupPublicKey[0]).to.equal(bytes32zero) // pubkey X
      expect(signingGroupPublicKey[1]).to.equal(bytes32zero) // pubkey Y

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.FAILED_SETUP)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting fraud funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyFraudFundingTimeout(),
        "Not currently awaiting fraud-related funding proof",
      )
    })

    it("reverts if the timer has not elapsed", async () => {
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart * 5)

      await expectRevert(
        testDeposit.notifyFraudFundingTimeout(),
        "Fraud funding proof timeout has not elapsed",
      )
    })

    it("asserts that it partially slashes signers", async () => {
      const initialBalance = await web3.eth.getBalance(beneficiary)
      const toSeize = await web3.eth.getBalance(ecdsaKeepStub.address)

      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)

      await testDeposit.notifyFraudFundingTimeout()

      const divisor = await tbtcConstants.getFundingFraudPartialSlashDivisor.call()
      const slash = new BN(toSeize).div(new BN(divisor))
      const balanceAfter = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(slash)

      expect(
        balanceCheck,
        "partial slash not correctly awarded to funder",
      ).to.eq.BN(balanceAfter)
    })
  })

  describe("provideFraudBTCFundingProof", async () => {
    beforeEach(async () => {
      await testDeposit.setSigningGroupPublicKey(_signerPubkeyX, _signerPubkeyY)
      await mockRelay.setCurrentEpochDifficulty(currentDifficulty)
      await testDeposit.setState(states.FRAUD_AWAITING_BTC_FUNDING_PROOF)
      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    it("updates to setup failed, logs SetupFailed, ", async () => {
      const blockNumber = await web3.eth.getBlock("latest").number

      await testDeposit.provideFraudBTCFundingProof(
        _version,
        _txInputVector,
        _txOutputVector,
        _txLocktime,
        _fundingOutputIndex,
        _merkleProof,
        _txIndexInBlock,
        _bitcoinHeaders,
      )

      const keepAddress = await testDeposit.getKeepAddress.call()
      expect(keepAddress, "Keep address not deconsted").to.equal(ZERO_ADDRESS)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt not deconsted",
      ).to.not.equal(0)

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "fundingProofTimerStart not deconsted",
      ).to.not.equal(0)

      const signingGroupPublicKey = await testDeposit.getSigningGroupPublicKey.call()
      expect(signingGroupPublicKey[0]).to.equal(bytes32zero) // pubkey X
      expect(signingGroupPublicKey[1]).to.equal(bytes32zero) // pubkey Y

      const depostState = await testDeposit.getState.call()
      expect(depostState).to.eq.BN(states.FAILED_SETUP)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting a funding proof during setup fraud", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideFraudBTCFundingProof(
          _version,
          _txInputVector,
          _txOutputVector,
          _txLocktime,
          _fundingOutputIndex,
          _merkleProof,
          _txIndexInBlock,
          _bitcoinHeaders,
        ),
        "Not awaiting a funding proof during setup fraud",
      )
    })

    it("assert distribute signer bonds to funder", async () => {
      const initialBalance = await web3.eth.getBalance(beneficiary)
      const signerBond = await web3.eth.getBalance(ecdsaKeepStub.address)

      await testDeposit.provideFraudBTCFundingProof(
        _version,
        _txInputVector,
        _txOutputVector,
        _txLocktime,
        _fundingOutputIndex,
        _merkleProof,
        _txIndexInBlock,
        _bitcoinHeaders,
      )

      const balanceAfter = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(new BN(signerBond))

      expect(
        balanceCheck,
        "partial slash not correctly awarded to funder",
      ).to.eq.BN(balanceAfter)
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

      const block = await web3.eth.getBlock("latest")
      const initialBalance = await web3.eth.getBalance(owner)
      await testDeposit.startSignerFraudLiquidation()

      const events = await tbtcSystemStub.getPastEvents("Liquidated", {
        fromBlock: block.number,
        toBlock: "latest",
      })

      await assertBalance.eth(
        owner,
        new BN(initialBalance).add(new BN(signerBond)),
      )

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
      await ecdsaKeepStub.setSuccess(true)
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
