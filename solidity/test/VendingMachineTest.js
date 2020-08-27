const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {increaseTime, states, fundingTx} = require("./helpers/utils.js")
const {AssertBalance} = require("./helpers/assertBalance.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")
const moment = require("moment")

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
      "0xb08d3b935947dd03c2b485deecb3629bb9d7bc10c80e3cc6af43b8673e07d41c"
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
        "0x0000111111111111",
        redeemerOutputScript,
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
        await tbtcToken.forceMint(testDeposit.address, signerFee)
        await feeRebateToken.forceMint(owner, tdtId)
        const blockNumber = await web3.eth.getBlockNumber()

        await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

        const tbtcToBtc = vendingMachine.abi.find(x => x.name == "tbtcToBtc")
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          testDeposit.address,
          "0x0000111111111111",
          redeemerOutputScript,
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
        await tbtcToken.forceMint(testDeposit.address, signerFee)
        await feeRebateToken.forceMint(owner, tdtId)
        const tbtcToBtc = vendingMachine.abi.find(x => x.name == "tbtcToBtc")
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          testDeposit.address,
          "0x0000111111111111",
          redeemerOutputScript,
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
        const tbtcToBtc = vendingMachine.abi.find(x => x.name == "tbtcToBtc")
        const nonexistentDeposit = "0000000000000000000000000000000000000000"
        const calldata = web3.eth.abi.encodeFunctionCall(tbtcToBtc, [
          nonexistentDeposit,
          "0x1111111100000000",
          redeemerOutputScript,
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

      it("protects receiveApproval from external callers", async () => {
        await expectRevert(
          redemptionScript.receiveApproval(
            redemptionScript.address,
            tdtId,
            owner,
            "0x00",
            {from: owner},
          ),
          "Only token contract can call receiveApproval",
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
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.find(
        x => x.name == "unqualifiedDepositToTbtc",
      )
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

    it("returns true on success", async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.find(
        x => x.name == "unqualifiedDepositToTbtc",
      )
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
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.find(
        x => x.name == "unqualifiedDepositToTbtc",
      )
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

    it("protects receiveApproval from external callers", async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.find(
        x => x.name == "unqualifiedDepositToTbtc",
      )
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

      await expectRevert(
        fundingScript.receiveApproval(
          redemptionScript.address,
          depositValue.add(signerFee),
          accounts[0],
          calldata,
          {from: owner},
        ),
        "Only token contract can call receiveApproval",
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

    const MINUTE = 60
    const HOUR = 60 * MINUTE
    const DAY = 24 * HOUR

    const SUPPLY_SCHEDULE = [
      [0, 2 * DAY, 2],
      [2 * DAY, 7 * DAY, 100],
      [7 * DAY, 14 * DAY, 250],
      [14 * DAY, 21 * DAY, 500],
      [21 * DAY, 28 * DAY, 750],
      [28 * DAY, 35 * DAY, 1000],
      [35 * DAY, 42 * DAY, 1500],
      [42 * DAY, 49 * DAY, 2000],
      [49 * DAY, 56 * DAY, 2500],
      [56 * DAY, 63 * DAY, 3000],
      [63 * DAY, Infinity, 21000000],
    ]

    SUPPLY_SCHEDULE.forEach(([start, end, expectedMax]) => {
      const humanize = function(interval) {
        if (interval === 0) {
          return "instantiation"
        } else if (interval === Infinity) {
          return "the end of time"
        }
        return moment.duration(interval, 'seconds').humanize({d: 7, w: 40})
      }

      it(`has a max supply of ${expectedMax} TBTC between ${humanize(start)}
          and ${humanize(end)}`, async () => {
        await increaseTime(start)
        let maxSupply = await vendingMachine.getMaxSupply()
        expect(maxSupply).to.eq.BN(btcToTbtc(expectedMax))

        await increaseTime(MINUTE)
        maxSupply = await vendingMachine.getMaxSupply()
        expect(maxSupply).to.eq.BN(btcToTbtc(expectedMax))

        if (end !== Infinity) {
          await increaseTime(end - start - 2 * MINUTE)

          maxSupply = await vendingMachine.getMaxSupply()
          expect(maxSupply).to.eq.BN(btcToTbtc(expectedMax))

          await increaseTime(2 * MINUTE)

          maxSupply = await vendingMachine.getMaxSupply()
          expect(maxSupply).to.not.eq.BN(btcToTbtc(expectedMax))
        }
        else {
          // if we've reached infinity, fast-forward 0-1000 days in advance
          // 10 times and confir
          for(i = 0; i<10; i++) {
            let rand = Math.floor(Math.random() * 1000)
            await increaseTime(rand * DAY)
            maxSupply = await vendingMachine.getMaxSupply()
            expect(maxSupply).to.eq.BN(btcToTbtc(expectedMax))
          }
        }
      })
    })
  })
})
