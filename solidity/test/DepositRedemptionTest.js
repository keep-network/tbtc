const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {
  states,
  bytes32zero,
  increaseTime,
  fundingTx,
  expectEvent,
  resolveAllLogs,
  expectNoEvent,
} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const {
  BN,
  constants,
  expectRevert,
  time,
} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'
describe("DepositRedemption", async function() {
  const redeemerOutputScript =
    "0x16001486e7303082a6a21d5837176bc808bf4828371ab6"

  let tbtcConstants
  let mockRelay
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let feeRebateToken
  let testDeposit
  let ecdsaKeepStub

  let withdrawalRequestTime

  let depositValue
  let signerFee
  let tdtId

  // Default holders for various accounts.
  const [owner, redeemer, , frtHolder] = accounts
  const tdtHolder = owner

  before(async () => {
    let deployed
    ;({
      tbtcConstants,
      mockRelay,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      ecdsaKeepStub,
      deployed,
    } = await deployAndLinkAll())
    vendingMachine = deployed.VendingMachine

    await testDeposit.setSignerFeeDivisor(new BN("200"))

    tdtId = await web3.utils.toBN(testDeposit.address)

    depositValue = await testDeposit.lotSizeTbtc.call()
    signerFee = await testDeposit.signerFeeTbtc.call()
    depositTerm = await tbtcConstants.getDepositTerm.call()
  })

  beforeEach(async () => {
    await testDeposit.reset()
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  describe("requestRedemption", async () => {
    // the TX produced will be:
    // 01000000000101333333333333333333333333333333333333333333333333333333333333333333333333000000000001111111110000000016001433333333333333333333333333333333333333330000000000
    // the signer pkh will be:
    // 5eb9b5e445db673f0ed8935d18cd205b214e5187
    // == hash160(023333333333333333333333333333333333333333333333333333333333333333)
    // the sighash preimage will be:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788ac111111111111111100000000e4ca7a168bd64e3123edd7f39e1ab7d670b32311cac2dda8e083822139c7936c0000000001000000
    const sighash =
      "0xb08d3b935947dd03c2b485deecb3629bb9d7bc10c80e3cc6af43b8673e07d41c"
    const outpoint = "0x" + "33".repeat(36)
    const valueBytes = "0x1111111111111111"
    const keepPubkeyX = "0x" + "33".repeat(32)
    const keepPubkeyY = "0x" + "44".repeat(32)
    // Override redeemer output script for this test.
    // No real reason here, it's just how we derived the below values.
    const redeemerOutputScript = "0x160014" + "33".repeat(20)
    let requiredBalance

    before(async () => {
      requiredBalance = depositValue
    })

    beforeEach(async () => {
      await createSnapshot()

      await feeRebateToken.forceMint(
        frtHolder,
        web3.utils.toBN(testDeposit.address),
      )
      await tbtcDepositToken.forceMint(tdtHolder, tdtId)

      await testDeposit.setState(states.ACTIVE)
      await testDeposit.setFundingInfo(valueBytes, 0, outpoint)

      // make sure there is sufficient balance to request redemption. Then approve deposit
      await tbtcToken.resetBalance(requiredBalance, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, requiredBalance, {
        from: owner,
      })
    })

    afterEach(restoreSnapshot)

    it("updates state successfully and fires a RedemptionRequested event", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is 2.86331153 BTC
      await testDeposit.requestRedemption(
        "0x0000111111111111",
        redeemerOutputScript,
        {from: owner},
      )

      const requestInfo = await testDeposit.getRequestInfo()
      expect(requestInfo[1]).to.equal(redeemerOutputScript)
      expect(requestInfo[3]).to.not.equal(0) // withdrawalRequestTime is set
      expect(requestInfo[4]).to.equal(sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "RedemptionRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._digest).to.equal(sighash)
    })

    it("updates state successfully and fires a RedemptionRequested event from COURTESY_CALL state", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setState(states.COURTESY_CALL)

      // the fee is 2.86331153 BTC
      await testDeposit.requestRedemption(
        "0x0000111111111111",
        redeemerOutputScript,
        {from: owner},
      )

      const requestInfo = await testDeposit.getRequestInfo()
      expect(requestInfo[1]).to.equal(redeemerOutputScript)
      expect(requestInfo[3]).to.not.equal(0) // withdrawalRequestTime is set
      expect(requestInfo[4]).to.equal(sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "RedemptionRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._digest).to.equal(sighash)
    })

    it("reverts if not in Active or Courtesy", async () => {
      await testDeposit.setState(states.LIQUIDATED)

      await expectRevert(
        testDeposit.requestRedemption(
          "0x1111111100000000",
          "0x" + "33".repeat(20),
          {from: owner},
        ),
        "Redemption only available from Active or Courtesy state",
      )
    })

    it("reverts if the fee is too low", async () => {
      await expectRevert(
        testDeposit.requestRedemption(
          "0x0011111111111111",
          "0x1976a914" + "33".repeat(20) + "88ac",
          {from: owner},
        ),
        "Fee is too low",
      )
    })

    it("reverts if the fee is too high", async () => {
      await expectRevert(
        testDeposit.requestRedemption(
          "0x8888888888888808",
          "0x1976a914" + "33".repeat(20) + "88ac",
          {from: owner},
        ),
        "Initial fee cannot exceed half of the deposit's value",
      )
    })

    it("does not revert if the fee is just under the max threshold", async () => {
      await testDeposit.requestRedemption(
        "0x8a88888888888808",
        "0x1976a914" + "33".repeat(20) + "88ac",
        {from: owner},
      )
    })

    it("reverts if the output script is non-standard", async () => {
      await testDeposit.setFundingInfo(
        valueBytes,
        await time.latest(),
        outpoint,
      )

      await expectRevert(
        testDeposit.requestRedemption(
          "0x1111111100000000",
          "0x" + "33".repeat(20),
          {from: owner},
        ),
        "Output script must be a standard type",
      )
    })

    const badLengths = {
      // declare 20 bytes, include 21
      p2pkh: "0x1976a914" + "33".repeat(21) + "88ac",
      // declare 20 bytes, include 21
      p2sh: "0x17a914" + "33".repeat(21) + "87",
      // declare 20 bytes, include 21
      p2wpkh: "0x160014" + "33".repeat(21),
      // declare 32 bytes, include 33
      p2wsh: "0x220020" + "33".repeat(33),
    }
    for (const [type, script] of Object.entries(badLengths)) {
      it(`reverts if ${type} output script has standard type but bad length`, async () => {
        await testDeposit.setFundingInfo(
          valueBytes,
          await time.latest(),
          outpoint,
        )

        await expectRevert(
          testDeposit.requestRedemption("0x1111111100000000", script, {
            from: tdtHolder,
          }),
          "Output script must be a standard type",
        )
      })
    }

    it("reverts if the caller is not the deposit owner", async () => {
      await testDeposit.setFundingInfo(
        valueBytes,
        await time.latest(),
        outpoint,
      )

      await tbtcDepositToken.transferFrom(tdtHolder, frtHolder, tdtId, {
        from: owner,
      })

      await expectRevert(
        testDeposit.requestRedemption(
          "0x1111111100000000",
          "0x1976a914" + "33".repeat(20) + "88ac",
          {from: owner},
        ),
        "Only TDT holder can redeem unless deposit is at-term or in COURTESY_CALL",
      )
    })
  })

  describe("transferAndRequestRedemption", async () => {
    const sighash =
      "0xb08d3b935947dd03c2b485deecb3629bb9d7bc10c80e3cc6af43b8673e07d41c"
    const outpoint = "0x" + "33".repeat(36)
    const valueBytes = "0x1111111111111111"
    const keepPubkeyX = "0x" + "33".repeat(32)
    const keepPubkeyY = "0x" + "44".repeat(32)
    // Override redeemer output script for this test.
    // No real reason here, it's just how we derived the below values.
    const redeemerOutputScript = "0x160014" + "33".repeat(20)
    let requiredBalance

    before(async () => {
      requiredBalance = depositValue
    })

    beforeEach(async () => {
      await createSnapshot()
      await testDeposit.setState(states.ACTIVE)
      await testDeposit.setFundingInfo(valueBytes, 0, outpoint)

      // make sure there is sufficient balance to request redemption. Then approve deposit
      await tbtcToken.resetBalance(requiredBalance, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, requiredBalance, {
        from: owner,
      })

      await feeRebateToken.forceMint(
        frtHolder,
        web3.utils.toBN(testDeposit.address),
      )

      await tbtcDepositToken.forceMint(tdtHolder, tdtId)

      await tbtcDepositToken.approve(testDeposit.address, tdtId, {
        from: owner,
      })
    })

    afterEach(restoreSnapshot)

    it("updates state successfully and fires a RedemptionRequested event", async () => {
      await testDeposit.setVendingMachineAddress(owner)
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is 2.86331153 BTC
      await testDeposit.transferAndRequestRedemption(
        "0x0000111111111111",
        redeemerOutputScript,
        owner,
        {from: owner},
      )
      const tdtOwner = await tbtcDepositToken.ownerOf(tdtId)
      const requestInfo = await testDeposit.getRequestInfo()

      expect(requestInfo[1]).to.equal(redeemerOutputScript)
      expect(requestInfo[3]).to.not.equal(0) // withdrawalRequestTime is set
      expect(requestInfo[4]).to.equal(sighash)
      expect(tdtOwner).to.equal(owner)
      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "RedemptionRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._digest).to.equal(sighash)
    })
    it("fails if the caller is not the Vending Machine", async () => {
      await expectRevert(
        testDeposit.transferAndRequestRedemption(
          "0x1111111100000000",
          redeemerOutputScript,
          owner,
          {from: owner},
        ),
        "Only the vending machine can call transferAndRequestRedemption",
      )
    })
  })

  describe("approveDigest", async () => {
    beforeEach(async () => {
      await createSnapshot()
      await testDeposit.setSigningGroupPublicKey("0x00", "0x00")
    })

    afterEach(restoreSnapshot)

    it("calls keep for signing", async () => {
      const digest = "0x" + "08".repeat(32)

      await testDeposit.approveDigest(digest).catch(err => {
        assert.fail(`cannot approve digest: ${err}`)
      })

      const blockNumber = await web3.eth.getBlockNumber()

      // Check if ECDSAKeep has been called and event emitted.
      const eventList = await ecdsaKeepStub.getPastEvents(
        "SignatureRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(
        eventList[0].returnValues.digest,
        "incorrect digest in emitted event",
      ).to.equal(digest)
    })

    it("registers timestamp for digest approval", async () => {
      const digest = "0x" + "02".repeat(32)

      const approvalTx = await testDeposit.approveDigest(digest).catch(err => {
        assert.fail(`cannot approve digest: ${err}`)
      })

      const block = await web3.eth.getBlock(approvalTx.receipt.blockNumber)
      const expectedTimestamp = block.timestamp

      const timestamp = await testDeposit
        .wasDigestApprovedForSigning(digest)
        .catch(err => {
          assert.fail(`cannot check digest approval: ${err}`)
        })
      expect(timestamp, "incorrect registered timestamp").to.eq.BN(
        new BN(expectedTimestamp),
      )
    })
  })

  describe("provideRedemptionSignature", async () => {
    // signing the sha 256 of '11' * 32
    // signing with privkey '11' * 32
    // using RFC 6979 nonce (libsecp256k1)
    const pubkeyX =
      "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
    const pubkeyY =
      "0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"
    const digest =
      "0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc"
    const v = 27
    const r =
      "0xd7e83e8687ba8b555f553f22965c74e81fd08b619a7337c5c16e4b02873b537e"
    const s =
      "0x633bf745cdf7ae303ca8a6f41d71b2c3a21fcbd1aed9e7ffffa295c08918c1b3"

    beforeEach(async () => {
      await createSnapshot()
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_SIGNATURE)
      await tbtcDepositToken.forceMint(tdtHolder, tdtId)
    })

    afterEach(restoreSnapshot)

    it("updates the state and logs GotRedemptionSignature", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.setSigningGroupPublicKey(pubkeyX, pubkeyY)
      await testDeposit.setRequestInfo(
        "0x" + "11".repeat(20),
        "0x" + "11".repeat(20),
        0,
        0,
        digest,
      )

      await testDeposit.provideRedemptionSignature(v, r, s)

      const state = await testDeposit.getState.call()
      expect(state).to.eq.BN(states.AWAITING_WITHDRAWAL_PROOF)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "GotRedemptionSignature",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._r).to.equal(r)
      expect(eventList[0].returnValues._s).to.equal(s)
    })

    it("errors if not awaiting withdrawal signature", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideRedemptionSignature(v, r, s),
        "Not currently awaiting a signature",
      )
    })

    it("errors on invaid sig", async () => {
      await expectRevert(
        testDeposit.provideRedemptionSignature(28, r, s),
        "Invalid signature",
      )
    })

    it("reverts if S value is on the upper half of the secp256k1 curve's order", async () => {
      const s =
        "0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A1"
      await expectRevert(
        testDeposit.provideRedemptionSignature(0, "0x0", s),
        "Malleable signature - s should be in the low half of secp256k1 curve's order",
      )
    })
  })

  describe("increaseRedemptionFee", async () => {
    // the signer pkh will be:
    // 5eb9b5e445db673f0ed8935d18cd205b214e5187
    // == hash160(023333333333333333333333333333333333333333333333333333333333333333)
    // prevSighash preimage is:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788acffffffffffffffff0000000077feba7d287431f43e61c3e716e92d266658eb223ebad0e64d83b747652e2b9b0000000001000000
    // nextSighash preimage is:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788acffffffffffffffff0000000044bf045101f1d83d0f2e017a73bb85857f25137c31cc7382ef363c909659f55a0000000001000000

    // Override redeemer output script for this test.
    // No real reason here, it's just how we derived the below values.
    const redeemerOutputScript = "0x160014" + "33".repeat(20)
    const prevSighash =
      "0xd94b6f3bf19147cc3305ef202d6bd64f9b9a12d4d19cc2d8c7f93ef58fc8fffe"
    const nextSighash =
      "0xbb56d80cfd71e90215c6b5200c0605b7b80689d3479187bc2c232e756033a560"
    const keepPubkeyX = "0x" + "33".repeat(32)
    const keepPubkeyY = "0x" + "44".repeat(32)
    const prevoutValueBytes = "0xffffffffffffffff"
    const previousOutputBytes = "0x0000ffffffffffff"
    const newOutputBytes = "0x0100feffffffffff"
    const initialFee = 0xffff
    const outpoint = "0x" + "33".repeat(36)
    let feeIncreaseTimer

    before(async () => {
      feeIncreaseTimer = await tbtcConstants.getIncreaseFeeTimer.call()
    })

    beforeEach(async () => {
      await createSnapshot()
      const blockTimestamp = await time.latest()
      withdrawalRequestTime = blockTimestamp - feeIncreaseTimer.toNumber()
      await testDeposit.setDigestApprovedAtTime(
        prevSighash,
        withdrawalRequestTime,
      )
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setFundingInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        initialFee,
        withdrawalRequestTime,
        prevSighash,
      )

      await feeRebateToken.forceMint(
        frtHolder,
        web3.utils.toBN(testDeposit.address),
      )

      await tbtcDepositToken.forceMint(tdtHolder, tdtId)
    })

    afterEach(restoreSnapshot)

    it("correctly increases fee", async () => {
      // Currently little-endian = big-endian for this var.
      const startValue = new BN(prevoutValueBytes.slice(2), 16)
      // Fee increment that will reach the full UTXO value on the 4th bump.
      const feeIncrement = startValue.divn(4)

      const outputBytes = [1, 2, 3, 4].map(increment => {
        const hexIncrement = startValue
          .sub(feeIncrement.muln(increment))
          .toString(16)
          .padStart(16, "0")

        // Convert to little-endian for contract.
        return "0x" + [...hexIncrement.matchAll(/../g)].reverse().join("")
      })

      // Hardcoded sighashes for the first 3 output value bytes.
      const sigHashes = [
        "0xdcf7d96f5fbb1fb49d32e13a27d3a007cb91116170b504a3c3b87d07c8f6f3b3",
        "0x5261fa987e51cde646eb8c229c5ae39955ee1c995e8e5fb38bc5e6fc470d2276",
        "0xac51039a4029871ca5509ce9e17f11efd7821a461b1b203b7994cb2feefe49a1",
      ]

      const blockTimestamp = await time.latest()
      await testDeposit.setDigestApprovedAtTime(sigHashes[0], blockTimestamp)

      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        "0x3fffffffffffffff",
        blockTimestamp,
        sigHashes[0],
      )

      // Loop through first several fee bumps to check them normally; the last
      // bump requires a tweak.
      for (let i = 0; i < sigHashes.length - 1; i++) {
        // State resets every fee bump.
        await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
        await increaseTime(feeIncreaseTimer)

        await testDeposit.increaseRedemptionFee(
          outputBytes[i],
          outputBytes[i + 1],
        )

        const updatedFee = await testDeposit.getLatestRedemptionFee.call()
        expect(updatedFee).to.eq.BN(feeIncrement.muln(2 + i))
      }

      // The last fee increase pushes the deposit below the minimum UTXO value,
      // so it will in fact be clamped to (UTXO value) - 2000 satoshis (the
      // minimum UTXO value constant).
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
      await increaseTime(feeIncreaseTimer)

      await testDeposit.increaseRedemptionFee(
        outputBytes.slice(-2)[0],
        outputBytes.slice(-2)[1],
      )

      const updatedFee = await testDeposit.getLatestRedemptionFee.call()
      expect(updatedFee).to.eq.BN(startValue.subn(2000))
    })

    it("approves a new digest for signing, updates the state, and logs RedemptionRequested", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.increaseRedemptionFee(
        previousOutputBytes,
        newOutputBytes,
      )
      const requestInfo = await testDeposit.getRequestInfo.call()
      expect(requestInfo[4]).to.equal(nextSighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "RedemptionRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._digest).to.equal(nextSighash)
    })

    it("reverts if not awaiting withdrawal proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        "Fee increase only available after signature provided",
      )
    })

    it("reverts if the increase fee timer has not elapsed", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        initialFee,
        await time.latest(),
        prevSighash,
      )

      await expectRevert(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        "Fee increase not yet permitted",
      )
    })

    it("reverts if the fee step is not linear", async () => {
      await expectRevert(
        testDeposit.increaseRedemptionFee(
          previousOutputBytes,
          "0x1101010101102201",
        ),
        "Not an allowed fee step",
      )
    })

    it("reverts if the previous sighash was not the latest approved", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        initialFee,
        withdrawalRequestTime,
        keepPubkeyX,
      )

      // Previous sigHash is not approved for signing.
      await testDeposit.setDigestApprovedAtTime(prevSighash, 0)

      await expectRevert(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        "Provided previous value does not yield previous sighash",
      )
    })
  })

  describe("provideRedemptionProof", async () => {
    beforeEach(async () => {
      await createSnapshot()

      await tbtcDepositToken.forceMint(tdtHolder, tdtId)
      // Mint the signer fee so we don't try to transfer nonexistent tokens eh.
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await testDeposit.setFundingInfo(
        fundingTx.prevoutValueBytes,
        0,
        fundingTx.prevoutOutpoint,
      )
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setRequestInfo(
        "0x" + "11".repeat(20),
        redeemerOutputScript,
        14544,
        0,
        "0x" + "11" * 32,
      )
      await testDeposit.setLatestRedemptionFee(14544)
    })

    afterEach(restoreSnapshot)

    it("updates the state, clears struct info except for redeemer address, calls TBTC and Keep, and emits a Redeemed event", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.provideRedemptionProof(
        fundingTx.version,
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
        fundingTx.txLocktime,
        fundingTx.merkleProof,
        fundingTx.txIndexInBlock,
        fundingTx.bitcoinHeaders,
      )

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.REDEEMED)

      const requestInfo = await testDeposit.getRequestInfo.call()
      expect(requestInfo[0]).to.equal(
        "0x1111111111111111111111111111111111111111",
      ) // this value should not be cleared
      expect(requestInfo[1]).to.equal(null)
      expect(requestInfo[4]).to.equal(bytes32zero)

      const eventList = await tbtcSystemStub.getPastEvents("Redeemed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList[0].returnValues._txid).to.equal(fundingTx.txidLE)
    })

    it("reverts if not in the redemption flow", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideRedemptionProof(
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ),
        "Redemption proof only allowed from redemption flow",
      )
    })

    it("reverts if the merkle proof is not validated successfully", async () => {
      await expectRevert(
        testDeposit.provideRedemptionProof(
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.merkleProof,
          0,
          fundingTx.bitcoinHeaders,
        ),
        "Tx merkle proof is not valid for provided header",
      )
    })

    it("reverts if a higher fee is sent", async () => {
      const currentFee = await testDeposit.getLatestRedemptionFee.call()
      await testDeposit.setLatestRedemptionFee(currentFee.sub(new BN(1)))
      await expectRevert(
        testDeposit.provideRedemptionProof(
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ),
        "Incorrect fee amount",
      )
    })
  })

  describe("redemptionTransactionChecks", async () => {
    beforeEach(async () => {
      await testDeposit.setFundingInfo(
        fundingTx.prevoutValueBytes,
        0,
        fundingTx.prevoutOutpoint,
      )
      await testDeposit.setRequestInfo(
        "0x" + "11".repeat(20),
        redeemerOutputScript,
        14544,
        0,
        "0x" + "11" * 32,
      )
    })

    it("returns the output value", async () => {
      const redemptionChecks = await testDeposit.redemptionTransactionChecks.call(
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
      )
      expect(redemptionChecks).to.eq.BN(new BN(fundingTx.outputValue))
    })

    it("accepts all standard output types", async () => {
      const outputScripts = [
        "1976a914" + "00".repeat(20) + "88ac", // pkh
        "17a914" + "00".repeat(20) + "87", // sh
        "160014" + "00".repeat(20), // wpkh
        "220020" + "00".repeat(32), // wsh
      ]

      for (let i = 0; i < outputScripts.length; i++) {
        const script = outputScripts[i]
        const tempOutputVector = "0x012040351d00000000" + script
        await testDeposit.setRequestInfo(
          "0x" + "11".repeat(20),
          "0x" + script,
          14544,
          0,
          "0x" + "11" * 32,
        )
        const redemptionChecks = await testDeposit.redemptionTransactionChecks.call(
          fundingTx.txInputVector,
          tempOutputVector,
        )
        expect(redemptionChecks).to.eq.BN(new BN(fundingTx.outputValue))
      }
    })

    it("reverts if bad input vector is provided", async () => {
      await expectRevert(
        testDeposit.redemptionTransactionChecks(
          "0x00",
          fundingTx.txOutputVector,
        ),
        "invalid input vector provided",
      )
    })

    it("reverts if bad output vector is provided", async () => {
      await expectRevert(
        testDeposit.redemptionTransactionChecks(
          fundingTx.txInputVector,
          "0x00",
        ),
        "invalid output vector provided",
      )
    })

    it("reverts if the tx spends the wrong utxo", async () => {
      await testDeposit.setFundingInfo(
        fundingTx.prevoutValueBytes,
        0,
        "0x" + "33".repeat(36),
      )

      await expectRevert(
        testDeposit.redemptionTransactionChecks.call(
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
        ),
        "Tx spends the wrong UTXO",
      )
    })

    it("reverts if the tx sends value to the wrong output script", async () => {
      await testDeposit.setRequestInfo(
        "0x" + "11".repeat(20),
        "0x" + "11".repeat(20),
        14544,
        0,
        "0x" + "11" * 32,
      )

      await expectRevert(
        testDeposit.redemptionTransactionChecks.call(
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
        ),
        "Tx sends value to wrong output script",
      )
    })
  })

  const abortScenarios = {
    "signature timeout": {
      timeoutFn: constants => constants.getSignatureTimeout,
      timeoutError: "Signature timer has not elapsed",
      state: states.AWAITING_WITHDRAWAL_SIGNATURE,
      stateError: "Not currently awaiting a signature",
      notifyFn: deposit => deposit.notifyRedemptionSignatureTimedOut,
    },
    "proof timeout": {
      timeoutFn: constants => constants.getRedemptionProofTimeout,
      timeoutError: "Proof timer has not elapsed",
      state: states.AWAITING_WITHDRAWAL_PROOF,
      stateError: "Not currently awaiting a redemption proof",
      notifyFn: deposit => deposit.notifyRedemptionProofTimedOut,
    },
  }

  for (const [
    scenario,
    {timeoutFn, timeoutError, state, stateError, notifyFn},
  ] of Object.entries(abortScenarios)) {
    describe(`when reporting a signer abort due to ${scenario}`, async () => {
      let abortTimeout

      before(async () => {
        abortTimeout = await timeoutFn(tbtcConstants).call()
      })

      beforeEach(async () => {
        await createSnapshot()
        await tbtcDepositToken.forceMint(tdtHolder, tdtId)

        await ecdsaKeepStub.burnContractBalance()
        await testDeposit.setState(state)
        await testDeposit.setRequestInfo(
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          0,
          await time.latest(),
          bytes32zero,
        )
      })

      it("should revert if not in correct state", async () => {
        testDeposit.setState(states.START)

        await expectRevert(notifyFn(testDeposit)(), stateError)
      })

      it("reverts if the signature timeout has not elapsed", async () => {
        await expectRevert(notifyFn(testDeposit)(), timeoutError)
      })

      it("reverts if no funds received as signer bond", async () => {
        await time.increase(time.duration.seconds(abortTimeout + 1))

        const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
        expect(new BN(0), "no bond should be sent").to.eq.BN(bond)

        await expectRevert(
          notifyFn(testDeposit)(),
          "No funds received, unexpected",
        )
      })

      it("initiates a liquidation auction", async () => {
        await time.increase(time.duration.seconds(abortTimeout + 1))

        const signerBonds = new BN("1000000")
        await ecdsaKeepStub.send(signerBonds, {from: owner})
        await testDeposit.setRedeemerAddress(redeemer)

        const {receipt} = await notifyFn(testDeposit)()
        const notificationTime = (await web3.eth.getBlock(receipt.blockNumber))
          .timestamp
        const fullReceipt = resolveAllLogs(receipt, {tbtcSystemStub})

        expectNoEvent(fullReceipt, "Liquidated")

        expectEvent(fullReceipt, "StartedLiquidation", {
          _depositContractAddress: testDeposit.address,
          _wasFraud: false,
          _timestamp: new BN(notificationTime),
        })
      })
    })
  }
})
