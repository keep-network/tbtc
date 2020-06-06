const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {increaseTime, states, fundingTx} = require("./helpers/utils.js")
const {AssertBalance} = require("./helpers/assertBalance.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")

function btcToTbtc(n) {
  return new BN(10).pow(new BN(18)).mul(new BN(n))
}

describe("VendingMachine", async function() {
  let vendingMachine
  let mockRelay
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let feeRebateToken
  let testDeposit

  let assertBalance
  let tdtId
  // this is the amount of TBTC exchanged for a TDT.
  let depositValue
  let signerFee

  const signerFeeDivisor = new BN("200")

  before(async () => {
    let deployed
    ;({
      mockRelay,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      deployed,
      redemptionScript,
      fundingScript,
    } = await deployAndLinkAll())
    vendingMachine = deployed.VendingMachine

    assertBalance = new AssertBalance(tbtcToken)

    await testDeposit.setSignerFeeDivisor(signerFeeDivisor)

    tdtId = await web3.utils.toBN(testDeposit.address)

    depositValue = await testDeposit.lotSizeTbtc()
    signerFee = depositValue.div(signerFeeDivisor)
  })

  describe("#isQualified", async () => {
    it("returns true if deposit is in ACTIVE State", async () => {
      await testDeposit.setState(states.ACTIVE)
      const qualified = await vendingMachine.isQualified.call(
        testDeposit.address,
      )
      expect(qualified).to.be.true
    })

    it("returns false if deposit is not in ACTIVE State", async () => {
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      const qualified = await vendingMachine.isQualified.call(
        testDeposit.address,
      )
      expect(qualified).to.be.false
    })
  })

  describe("#tdtToTbtc", async () => {
    before(async () => {
      await testDeposit.setState(states.ACTIVE)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("converts TDT to TBTC", async () => {
      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })

      await vendingMachine.tdtToTbtc(tdtId, {from: owner})

      await assertBalance.tbtc(owner, depositValue.sub(signerFee))
    })

    it("mints full lot size if backing deposit has signer fee escrowed", async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })

      await vendingMachine.tdtToTbtc(tdtId, {from: owner})

      await assertBalance.tbtc(owner, depositValue)
    })

    it("fails if deposit not qualified", async () => {
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectRevert(
        vendingMachine.tdtToTbtc(tdtId, {from: owner}),
        "Deposit must be qualified",
      )
    })

    it(`fails if TDT doesn't exist`, async () => {
      await expectRevert(
        vendingMachine.tdtToTbtc(new BN(123345), {from: owner}),
        "tBTC Deposit Token does not exist",
      )
    })

    it(`fails if TDT transfer not approved`, async () => {
      await tbtcDepositToken.forceMint(owner, tdtId)

      await expectRevert(
        vendingMachine.tdtToTbtc(tdtId, {from: owner}),
        "ERC721: transfer caller is not owner nor approved.",
      )
    })

    it("updates the minted supply", async () => {
      const mintedSupply = await vendingMachine.getMintedSupply()

      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })
      await vendingMachine.tdtToTbtc(tdtId, {from: owner})

      const newMintedSupply = await vendingMachine.getMintedSupply()

      expect(mintedSupply.add(depositValue)).to.eq.BN(newMintedSupply)
    })

    it("fails if the max supply has been reached", async () => {
      const maxSupply = await vendingMachine.getMaxSupply()
      const mintedSupply = await vendingMachine.getMintedSupply()

      // mint the difference, minus the deposit value, plus 1 weitoshi
      const toMint = new BN("1")
        .add(maxSupply)
        .sub(mintedSupply)
        .sub(depositValue)
      await tbtcToken.forceMint(owner, toMint)

      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })
      await expectRevert(
        vendingMachine.tdtToTbtc(tdtId, {from: owner}),
        "Can't mint more than the max supply cap",
      )

      await assertBalance.tbtc(owner, toMint)
    })
  })

  describe("#tbtcToTdt", async () => {
    before(async () => {
      await testDeposit.setState(states.ACTIVE)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("converts TBTC to TDT", async () => {
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
      await tbtcToken.forceMint(owner, depositValue)
      await tbtcToken.approve(vendingMachine.address, depositValue, {
        from: owner,
      })

      const fromBlock = await web3.eth.getBlockNumber()
      await vendingMachine.tbtcToTdt(tdtId, {from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock,
        toBlock: "latest",
      })
      const tbtcBurntEvent = events[0]
      expect(tbtcBurntEvent.returnValues.from).to.equal(owner)
      expect(tbtcBurntEvent.returnValues.to).to.equal(ZERO_ADDRESS)
      expect(tbtcBurntEvent.returnValues.value).to.equal(
        depositValue.toString(),
      )

      expect(await tbtcDepositToken.ownerOf(tdtId)).to.equal(owner)
    })

    it("fails if deposit not qualified", async () => {
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectRevert(
        vendingMachine.tbtcToTdt(tdtId, {from: owner}),
        "Deposit must be qualified",
      )
    })

    it(`fails if caller hasn't got enough TBTC`, async () => {
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectRevert(
        vendingMachine.tbtcToTdt(tdtId, {from: owner}),
        "Not enough TBTC for TDT exchange.",
      )
    })

    it(`fails if TDT doesn't exist`, async () => {
      await expectRevert(
        vendingMachine.tdtToTbtc(new BN(123345), {from: owner}),
        "tBTC Deposit Token does not exist",
      )
    })

    it(`fails if deposit is locked`, async () => {
      // Deposit is locked if the tBTC Deposit Token is not owned by the vending machine
      const depositOwner = accounts[1]
      await tbtcDepositToken.forceMint(depositOwner, tdtId)
      await tbtcToken.forceMint(owner, depositValue)
      await tbtcToken.approve(vendingMachine.address, depositValue, {
        from: owner,
      })

      await expectRevert(
        vendingMachine.tbtcToTdt(tdtId, {from: owner}),
        "Deposit is locked.",
      )
    })
  })

  describe("#unqualifiedDepositToTbtc", async () => {
    before(async () => {
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(
        fundingTx.signerPubkeyX,
        fundingTx.signerPubkeyY,
      )
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("qualifies a Deposit", async () => {
      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })
      const blockNumber = await web3.eth.getBlockNumber()

      await vendingMachine.unqualifiedDepositToTbtc(
        testDeposit.address,
        fundingTx.version,
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
        fundingTx.txLocktime,
        fundingTx.fundingOutputIndex,
        fundingTx.merkleProof,
        fundingTx.txIndexInBlock,
        fundingTx.bitcoinHeaders,
        {from: owner},
      )

      const UTXOInfo = await testDeposit.fundingInfo.call()
      expect(UTXOInfo[0]).to.equal(fundingTx.outValueBytes)
      expect(UTXOInfo[2]).to.equal(fundingTx.expectedUTXOOutpoint)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt not updated",
      ).to.not.equal(0)

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "undingProofTimerStart not updated",
      ).to.not.equal(0)

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.be.bignumber.equal(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents("Funded", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("mints TBTC to the TDT owner and siger fee to Deposit", async () => {
      await tbtcDepositToken.forceMint(owner, tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, {
        from: owner,
      })

      await vendingMachine.unqualifiedDepositToTbtc(
        testDeposit.address,
        fundingTx.version,
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
        fundingTx.txLocktime,
        fundingTx.fundingOutputIndex,
        fundingTx.merkleProof,
        fundingTx.txIndexInBlock,
        fundingTx.bitcoinHeaders,
        {from: owner},
      )
      await assertBalance.tbtc(owner, depositValue.sub(signerFee))
      await assertBalance.tbtc(testDeposit.address, signerFee)
    })
  })

  describe("#tbtcToBtc", async () => {
    const sighash =
      "0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90"
    const outpoint = "0x" + "33".repeat(36)
    const valueBytes = "0x1111111111111111"
    const keepPubkeyX = "0x" + "33".repeat(32)
    const keepPubkeyY = "0x" + "44".repeat(32)
    const redeemerOutputScript = "0x160014" + "33".repeat(20)
    let requiredBalance
    let block

    before(async () => {
      requiredBalance = await testDeposit.lotSizeTbtc.call()

      await tbtcToken.zeroBalance({from: owner})
      block = await web3.eth.getBlock("latest")
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
      await tbtcToken.resetBalance(requiredBalance, {from: owner})
      await tbtcToken.resetAllowance(vendingMachine.address, requiredBalance, {
        from: owner,
      })
      await testDeposit.setState(states.ACTIVE)
      await testDeposit.setFundingInfo(valueBytes, block.timestamp, outpoint)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })
    it("successfully redeems via wrapper", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      // the fee is 2.86331153 BTC
      await feeRebateToken.forceMint(owner, tdtId)
      await vendingMachine.tbtcToBtc(
        testDeposit.address,
        "0x1111111100000000",
        redeemerOutputScript,
        owner,
        {from: owner},
      )
      const requestInfo = await testDeposit.getRequestInfo()

      expect(requestInfo[1]).to.equal(redeemerOutputScript)
      expect(requestInfo[3]).to.not.equal(0)
      expect(requestInfo[4]).to.equal(sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents(
        "RedemptionRequested",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList[0].returnValues._digest).to.equal(sighash)
    })

    it("fails to redeem with insufficient balance", async () => {
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is 2.86331153 BTC
      // requester does not own the FRT, and therefore owes an additional SignerFee
      await feeRebateToken.forceMint(accounts[1], tdtId)

      await expectRevert(
        vendingMachine.tbtcToBtc(
          testDeposit.address,
          "0x1111111100000000",
          redeemerOutputScript,
          owner,
        ),
        "SafeMath: subtraction overflow.",
      )
    })

    describe("RedemptionScript", async () => {
      beforeEach(async () => {
        await createSnapshot()
      })

      afterEach(async () => {
        await restoreSnapshot()
      })

      it("successfully requests redemption", async () => {
        await testDeposit.setState(states.ACTIVE)
        await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
        await tbtcToken.forceMint(owner, depositValue.add(signerFee))
        await feeRebateToken.forceMint(owner, tdtId)

        const blockNumber = await web3.eth.getBlockNumber()

        await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

        const tbtcToBtc = vendingMachine.abi.filter(
          x => x.name == "tbtcToBtc",
        )[0]
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          testDeposit.address,
          "0x1111111100000000",
          redeemerOutputScript,
          owner,
        ])

        await tbtcToken.approveAndCall(
          redemptionScript.address,
          depositValue.add(signerFee),
          calldata,
          {from: owner},
        )

        const requestInfo = await testDeposit.getRequestInfo()
        expect(requestInfo[1]).to.equal(redeemerOutputScript)
        expect(requestInfo[3]).to.not.equal(0)
        expect(requestInfo[4]).to.equal(sighash)

        // fired an event
        const eventList = await tbtcSystemStub.getPastEvents(
          "RedemptionRequested",
          {fromBlock: blockNumber, toBlock: "latest"},
        )
        expect(eventList[0].returnValues._digest).to.equal(sighash)
      })

      it("returns true on success", async () => {
        await testDeposit.setState(states.ACTIVE)
        await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
        await tbtcToken.forceMint(owner, depositValue.add(signerFee))
        await feeRebateToken.forceMint(owner, tdtId)
        const tbtcToBtc = vendingMachine.abi.filter(
          x => x.name == "tbtcToBtc",
        )[0]
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          testDeposit.address,
          "0x1111111100000000",
          redeemerOutputScript,
          owner,
        ])

        const success = await tbtcToken.approveAndCall.call(
          redemptionScript.address,
          depositValue.add(signerFee),
          calldata,
          {from: owner},
        )

        expect(success).to.equal(true)
      })

      it("forwards nested revert error messages", async () => {
        await testDeposit.setState(states.ACTIVE)
        await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
        await tbtcToken.forceMint(owner, depositValue.add(signerFee))
        await feeRebateToken.forceMint(owner, tdtId)
        const tbtcToBtc = vendingMachine.abi.filter(
          x => x.name == "tbtcToBtc",
        )[0]
        const nonexistentDeposit = "0000000000000000000000000000000000000000"
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          nonexistentDeposit,
          "0x1111111100000000",
          redeemerOutputScript,
          owner,
        ])

        await expectRevert(
          tbtcToken.approveAndCall(
            redemptionScript.address,
            depositValue.add(signerFee),
            calldata,
            {from: owner},
          ),
          "tBTC Deposit Token does not exist",
        )
      })

      it("reverts for unknown function calls encoded in _extraData", async () => {
        const unknownFunctionSignature = "0xCAFEBABE"
        await tbtcToken.forceMint(owner, depositValue.add(signerFee))

        await expectRevert(
          tbtcToken.approveAndCall(
            redemptionScript.address,
            depositValue.add(signerFee),
            unknownFunctionSignature,
            {from: owner},
          ),
          "Bad _extraData signature. Call must be to tbtcToBtc.",
        )
      })
    })
  })

  describe("FundingScript", async () => {
    before(async () => {
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(
        fundingTx.signerPubkeyX,
        fundingTx.signerPubkeyY,
      )
      await tbtcToken.zeroBalance({from: owner})
      await tbtcDepositToken.forceMint(owner, tdtId)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("calls unqualifiedDepositToTbtcABI", async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.filter(
        x => x.name == "unqualifiedDepositToTbtc",
      )[0]
      const calldata = web3.eth.abi.encodeFunctionCall(
        unqualifiedDepositToTbtcABI,
        [
          testDeposit.address,
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.fundingOutputIndex,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ],
      )

      await tbtcDepositToken.approveAndCall(
        fundingScript.address,
        tdtId,
        calldata,
        {from: owner},
      )

      const UTXOInfo = await testDeposit.fundingInfo.call()
      expect(UTXOInfo[0]).to.equal(fundingTx.outValueBytes)
      expect(UTXOInfo[2]).to.equal(fundingTx.expectedUTXOOutpoint)

      await assertBalance.tbtc(owner, depositValue.sub(signerFee))
      expect(await tbtcDepositToken.ownerOf(tdtId)).to.equal(
        vendingMachine.address,
      )
      expect(await feeRebateToken.ownerOf(tdtId)).to.equal(owner)
    })

    it("reverts true on success", async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.filter(
        x => x.name == "unqualifiedDepositToTbtc",
      )[0]
      const calldata = web3.eth.abi.encodeFunctionCall(
        unqualifiedDepositToTbtcABI,
        [
          testDeposit.address,
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.fundingOutputIndex,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ],
      )

      const success = await tbtcDepositToken.approveAndCall.call(
        fundingScript.address,
        tdtId,
        calldata,
        {from: owner},
      )

      expect(success).to.be.true
    })

    it("forwards nested revert error messages", async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.filter(
        x => x.name == "unqualifiedDepositToTbtc",
      )[0]
      const calldata = web3.eth.abi.encodeFunctionCall(
        unqualifiedDepositToTbtcABI,
        [
          testDeposit.address,
          fundingTx.version,
          fundingTx.txInputVector,
          fundingTx.txOutputVector,
          fundingTx.txLocktime,
          fundingTx.fundingOutputIndex,
          fundingTx.merkleProof,
          fundingTx.txIndexInBlock,
          fundingTx.bitcoinHeaders,
        ],
      )

      // To make the funding call fail.
      await testDeposit.setState(states.ACTIVE)

      await expectRevert(
        tbtcDepositToken.approveAndCall.call(
          fundingScript.address,
          tdtId,
          calldata,
          {from: owner},
        ),
        "Not awaiting funding",
      )
    })

    it("reverts for unknown function calls encoded in _extraData", async () => {
      const unknownFunctionSignature = "0xCAFEBABE"

      await expectRevert(
        tbtcDepositToken.approveAndCall(
          fundingScript.address,
          tdtId,
          unknownFunctionSignature,
          {from: owner},
        ),
        "Bad _extraData signature. Call must be to unqualifiedDepositToTbtc.",
      )
    })
  })

  describe("getMaxSupply", async () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("has a max supply of 2 on the first day", async () => {
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(2))

      await increaseTime(23.5 * 60 * 60) // 23.5 hours

      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(2))

      await increaseTime(60 * 60) // 1 hour

      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.not.eq.BN(btcToTbtc(2))
    })

    it("has a max supply of 100 BTC between the first day and 30th day", async () => {
      await increaseTime(24 * 60 * 60 + 1) // 1 day and 1 second
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(100))

      await increaseTime(29 * 24 * 60 * 60 - 10 * 60) // 30 days minus 10 minutes
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(100))

      await increaseTime(1 * 60 * 60) // one hour
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.not.eq.BN(btcToTbtc(100))
    })

    it("has a max supply of 250 BTC between the 30th day and 60th day", async () => {
      await increaseTime(30 * 24 * 60 * 60 + 1) // 30 days and 1 second
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(250))

      await increaseTime(30 * 24 * 60 * 60 - 10 * 60) // 30 days minus 10 minutes
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(250))

      await increaseTime(60 * 60) // one hour
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.not.eq.BN(btcToTbtc(250))
    })

    it("has a max supply of 500 BTC between the 60th day and 90th day", async () => {
      await increaseTime(60 * 24 * 60 * 60 + 1) // 60 days and 1 second
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(500))

      await increaseTime(30 * 24 * 60 * 60 - 10 * 60) // 30 days minus 10 minutes
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(500))

      await increaseTime(60 * 60) // one hour
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.not.eq.BN(btcToTbtc(500))
    })

    it("has a max supply of 1000 BTC between the 90th day and 120th day", async () => {
      await increaseTime(90 * 24 * 60 * 60 + 1) // 90 days and 1 second
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(1000))

      await increaseTime(30 * 24 * 60 * 60 - 10 * 60) // 30 days minus 10 minutes
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(1000))

      await increaseTime(60 * 60) // one hour
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.not.eq.BN(btcToTbtc(1000))
    })

    it("has a max supply of 21M BTC after the 120th day", async () => {
      await increaseTime(120 * 24 * 60 * 60 + 1) // 120 days and 1 second
      let maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(21000000))

      await increaseTime(30 * 24 * 60 * 60) // 30 days
      maxSupply = await vendingMachine.getMaxSupply()
      expect(maxSupply).to.eq.BN(btcToTbtc(21000000))
    })
  })
})
