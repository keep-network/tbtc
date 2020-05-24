const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {
  states,
  fundingTx,
  increaseTime,
  legacyFundingTx,
} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")
const ECDSAKeepStub = contract.fromArtifact("ECDSAKeepStub")
const TBTCSystem = contract.fromArtifact("TBTCSystem")

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

describe("DepositFunding", async function() {
  let tbtcConstants
  let mockRelay
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let testDeposit
  let ecdsaKeepStub
  let ecdsaKeepFactory

  let fundingProofTimerStart
  let beneficiary

  const funderBondAmount = new BN("10").pow(new BN("5"))
  const fullBtc = 100000000

  before(async () => {
    ;({
      tbtcConstants,
      mockRelay,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      testDeposit,
      ecdsaKeepStub,
      ecdsaKeepFactoryStub,
    } = await deployAndLinkAll())

    ecdsaKeepFactory = ecdsaKeepFactoryStub
    beneficiary = accounts[4]
    await tbtcDepositToken.forceMint(
      beneficiary,
      web3.utils.toBN(testDeposit.address),
    )

    await testDeposit.reset()
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("createNewDeposit", async () => {
    it("runs and updates state and fires a created event", async () => {
      const expectedKeepAddress = "0x0000000000000000000000000000000000000007"

      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.createNewDeposit(
        tbtcSystemStub.address,
        tbtcToken.address,
        tbtcDepositToken.address,
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        1, // m
        1,
        fullBtc,
        {value: funderBondAmount},
      )

      // state updates
      const depositState = await testDeposit.getState.call()
      expect(depositState, "state not as expected").to.eq.BN(
        states.AWAITING_SIGNER_SETUP,
      )

      const systemSignerFeeDivisor = await tbtcSystemStub.getSignerFeeDivisor()
      const signerFeeDivisor = await testDeposit.getSignerFeeDivisor.call()
      expect(signerFeeDivisor).to.eq.BN(systemSignerFeeDivisor)

      const keepAddress = await testDeposit.getKeepAddress.call()
      expect(keepAddress, "keepAddress not as expected").to.equal(
        expectedKeepAddress,
      )

      const stakeLockDuration = await ecdsaKeepFactory.stakeLockDuration.call()
      expect(
        stakeLockDuration,
        "stake lock duration not as expected",
      ).not.to.eq.BN(await tbtcConstants.getDepositTerm.call())

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signing group timestamp not as expected",
      ).not.to.eq.BN(0)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents("Created", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList[0].returnValues._keepAddress).to.equal(
        expectedKeepAddress,
      )
    })

    it("reverts if not in the start state", async () => {
      await testDeposit.setState(states.REDEEMED)

      await expectRevert(
        testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          fullBtc,
        ),
        "Deposit setup already requested",
      )
    })

    it("fails if new deposits are disabled", async () => {
      await tbtcSystemStub.emergencyPauseNewDeposits()

      await expectRevert(
        testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          fullBtc,
        ),
        "Opening new deposits is currently disabled.",
      )
    })

    it("respects the supply cap schedule", async () => {
      const bn = web3.utils.toBN

      const mint = amountInSats =>
        tbtcToken.forceMint(
          testDeposit.address,
          bn(amountInSats).mul(bn(10).pow(bn(10))),
        )

      const createNewDeposit = amountInSats => {
        return testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          amountInSats,
        )
      }

      await mint(fullBtc - 1000)
      await createNewDeposit(fullBtc)
      await mint(fullBtc) // should now be at 2 BTC - 1000 sats

      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )

      await increaseTime(15 * 24 * 60 * 60) // 15 days, into the 1st month
      await createNewDeposit(fullBtc)
      await mint(98 * fullBtc) // should new be at 100 BTC - 1000 sats
      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )

      await increaseTime(30 * 24 * 60 * 60) // 30 days, into the 2nd month
      await createNewDeposit(fullBtc)
      await mint(150 * fullBtc) // should new be at 250 BTC - 1000 sats
      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )

      await increaseTime(30 * 24 * 60 * 60) // 30 days, into the 3rd month
      await createNewDeposit(fullBtc)
      await mint(250 * fullBtc) // should new be at 500 BTC - 1000 sats
      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )

      await increaseTime(30 * 24 * 60 * 60) // 30 days, into the 4th month
      await createNewDeposit(fullBtc)
      await mint(500 * fullBtc) // should new be at 1000 BTC - 1000 sats
      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )

      await increaseTime(30 * 24 * 60 * 60) // 30 days, into the 5th month
      await createNewDeposit(fullBtc)
      await mint("2099900000000000") // should new be at 21M BTC - 1000 sats
      await expectRevert(
        createNewDeposit(fullBtc),
        "Opening new deposits is currently disabled.",
      )
    })
  })
  describe("notifySignerSetupFailure", async () => {
    let timer
    let owner
    let openKeepFee

    before(async () => {
      ;({
        tbtcConstants,
        mockRelay,
        tbtcSystemStub,
        tbtcToken,
        tbtcDepositToken,
        testDeposit,
        ecdsaKeepStub,
      } = await deployAndLinkAll([], {TBTCSystemStub: TBTCSystem}))

      openKeepFee = await ecdsaKeepFactory.openKeepFeeEstimate.call()
      await testDeposit.setKeepSetupFee(openKeepFee)
      owner = accounts[1]
      await tbtcDepositToken.forceMint(
        owner,
        web3.utils.toBN(testDeposit.address),
      )
      timer = await tbtcConstants.getSigningGroupFormationTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      const value = openKeepFee + 100

      await ecdsaKeepStub.send(value)

      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_SIGNER_SETUP)

      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    it("updates state to setup failed, deconstes state, logs SetupFailed, and refunds TDT owner", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.notifySignerSetupFailure({from: owner})

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: owner,
      })

      const depositBalance = await web3.eth.getBalance(testDeposit.address)
      expect(withdrawable).to.eq.BN(new BN(openKeepFee))
      expect(withdrawable).to.eq.BN(new BN(depositBalance))

      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt should be 0",
      ).to.eq.BN(0)

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "fundingProofTimerStart should be 0",
      ).to.eq.BN(0)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length, "Event list is the wrong length").to.equal(1)
    })

    it("reverts if not awaiting signer setup", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifySignerSetupFailure(),
        "Not awaiting setup",
      )
    })

    it("reverts if the timer has not yet elapsed", async () => {
      await testDeposit.setSigningGroupRequestedAt(fundingProofTimerStart * 5)

      await expectRevert(
        testDeposit.notifySignerSetupFailure(),
        "Signing group formation timeout not yet elapsed",
      )
    })
  })

  describe("retrieveSignerPubkey", async () => {
    const publicKey =
      "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"
    const pubkeyX =
      "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
    const pubkeyY =
      "0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"

    let ecdsaKeepStub

    before(async () => {
      ecdsaKeepStub = await ECDSAKeepStub.new()
    })

    beforeEach(async () => {
      await testDeposit.setState(states.AWAITING_SIGNER_SETUP)
      await testDeposit.setKeepAddress(ecdsaKeepStub.address)
      await ecdsaKeepStub.setPublicKey(publicKey)
    })

    it("updates the pubkey X and Y, changes state, and logs RegisteredPubkey", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.retrieveSignerPubkey()

      const signingGroupPublicKey = await testDeposit.getSigningGroupPublicKey.call()
      expect(signingGroupPublicKey[0]).to.equal(pubkeyX)
      expect(signingGroupPublicKey[1]).to.equal(pubkeyY)

      const eventList = await tbtcSystemStub.getPastEvents("RegisteredPubkey", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(
        eventList[0].returnValues._signingGroupPubkeyX,
        "Logged X is wrong",
      ).to.equal(pubkeyX)
      expect(
        eventList[0].returnValues._signingGroupPubkeyY,
        "Logged Y is wrong",
      ).to.equal(pubkeyY)
    })

    it("reverts if not awaiting signer setup", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "Not currently awaiting signer setup",
      )
    })

    it("reverts when public key is not 64-bytes long", async () => {
      await ecdsaKeepStub.setPublicKey("0x" + "00".repeat(63))

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "public key not set or not 64-bytes long",
      )
    })

    it("reverts if either half of the pubkey is 0", async () => {
      await ecdsaKeepStub.setPublicKey("0x" + "00".repeat(64))

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "Keep returned bad pubkey",
      )
    })
  })

  describe("notifyFundingTimeout", async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    it("updates the state to failed setup, deconsts funding info, and logs SetupFailed", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.notifyFundingTimeout()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.FAILED_SETUP)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting a funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyFundingTimeout(),
        "Funding timeout has not started",
      )
    })

    it("reverts if the timeout has not elapsed", async () => {
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart * 5)

      await expectRevert(
        testDeposit.notifyFundingTimeout(),
        "Funding timeout has not elapsed",
      )
    })
  })

  describe("requestFunderAbort", async () => {
    let timer
    let owner

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
      owner = accounts[1]
      await tbtcDepositToken.forceMint(
        owner,
        web3.utils.toBN(testDeposit.address),
      )
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    it("fails if the deposit has not failed setup", () => {
      expectRevert(
        testDeposit.requestFunderAbort("0x1234", {from: owner}),
        "The deposit has not failed funding",
      )
    })

    it("emits a FunderAbortRequested event", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.notifyFundingTimeout()

      const outputScript = "0x012345"
      await testDeposit.requestFunderAbort(outputScript, {from: owner})

      const eventList = await tbtcSystemStub.getPastEvents(
        "FunderAbortRequested",
        {
          fromBlock: blockNumber,
          toBlock: "latest",
        },
      )
      expect(eventList.length).to.equal(1)
      expect(eventList[0].name == "FunderAbortRequested")
      expect(eventList[0].returnValues).to.contain({
        _depositContractAddress: testDeposit.address,
        _abortOutputScript: outputScript,
      })
    })
  })

  describe("provideBTCFundingProof", async () => {
    beforeEach(async () => {
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(
        fundingTx.signerPubkeyX,
        fundingTx.signerPubkeyY,
      )
      await ecdsaKeepStub.send(1000000, {from: accounts[0]})
    })

    it("updates to active, stores UTXO info, deconstes funding info, logs Funded", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      const {
        receipt: {blockNumber: proofBlock},
      } = await testDeposit.provideBTCFundingProof(
        fundingTx.version,
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
        fundingTx.txLocktime,
        fundingTx.fundingOutputIndex,
        fundingTx.merkleProof,
        fundingTx.txIndexInBlock,
        fundingTx.bitcoinHeaders,
      )
      const expectedFundedAt = (await web3.eth.getBlock(proofBlock)).timestamp

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      expect(UTXOInfo[0]).to.equal(fundingTx.outValueBytes)
      expect(UTXOInfo[1]).to.eq.BN(new BN(expectedFundedAt))
      expect(UTXOInfo[2]).to.equal(fundingTx.expectedUTXOOutpoint)

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

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents("Funded", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideBTCFundingProof(
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.fundingOutputIndex,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ),
        "Not awaiting funding",
      )
    })
  })

  describe("provideBTCFundingProof Legacy", async () => {
    beforeEach(async () => {
      await mockRelay.setCurrentEpochDifficulty(legacyFundingTx.difficulty)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(
        legacyFundingTx.signerPubkeyX,
        legacyFundingTx.signerPubkeyY,
      )
      await ecdsaKeepStub.send(1000000, {from: accounts[0]})
    })

    it("updates to active, stores UTXO info, deconsts funding info, logs Funded", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.provideBTCFundingProof(
        legacyFundingTx.version,
        legacyFundingTx.txInputVector,
        legacyFundingTx.txOutputVector,
        legacyFundingTx.txLocktime,
        legacyFundingTx.fundingOutputIndex,
        legacyFundingTx.merkleProof,
        legacyFundingTx.txIndexInBlock,
        legacyFundingTx.bitcoinHeaders,
      )

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      expect(UTXOInfo[0]).to.eql(legacyFundingTx.outValueBytes)
      expect(UTXOInfo[2]).to.eql(legacyFundingTx.expectedUTXOOutpoint)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(signingGroupRequestedAt).to.eq.BN(
        0,
        "signingGroupRequestedAt not deconsted",
      )

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(fundingProofTimerStart).to.eq.BN(
        0,
        "fundingProofTimerStart not deconsted",
      )

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents("Funded", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.eql(1)
    })
  })
})
