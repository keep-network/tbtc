const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {
  states,
  bytes32zero,
  increaseTime,
  fundingTx,
} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
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
  let depositTerm
  let tdtId
  // TODO tdtHolder, frtHolder
  let vendingMachine

  const tdtHolder = owner
  const frtHolder = accounts[4]

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

    await feeRebateToken.forceMint(
      frtHolder,
      web3.utils.toBN(testDeposit.address),
    )

    tdtId = await web3.utils.toBN(testDeposit.address)
    await tbtcDepositToken.forceMint(tdtHolder, tdtId)

    depositValue = await testDeposit.lotSizeTbtc.call()
    signerFee = await testDeposit.signerFee.call()
    depositTerm = await tbtcConstants.getDepositTerm.call()
  })

  beforeEach(async () => {
    await testDeposit.reset()
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  describe("getOwnerRedemptionTbtcRequirement", async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = "0x" + "33".repeat(36)
      valueBytes = "0x1111111111111111"
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock("latest")
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("returns signerFee if we are pre-term and owner is not FRT holder", async () => {
      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(
        owner,
      )
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it("returns zero if deposit is pre-term, owner is FRT holder and signer fee is escrowed", async () => {
      await feeRebateToken.transferFrom(accounts[4], owner, tdtId, {
        from: accounts[4],
      })
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(
        owner,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns signer fee if deposit is pre-term and signer fee is not escrowed", async () => {
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(
        owner,
      )
      expect(tbtcOwed).to.eq.BN(signerFee)
    })
  })

  describe("getRedemptionTbtcRequirement", async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = "0x" + "33".repeat(36)
      valueBytes = "0x1111111111111111"
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock("latest")
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("returns signerFee if we are pre term and redeemer is not FRT holder", async () => {
      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it("returns zero if deposit is pre-term, redeemer is FRT holder and signer fee is escrowed", async () => {
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, {
        from: frtHolder,
      })
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns full TBTC if we are pre term and we are at COURTESY_CALL - not TDT owner", async () => {
      await testDeposit.setState(states.COURTESY_CALL)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        accounts[9],
      )
      expect(tbtcOwed).to.eq.BN(depositValue)
    })

    it("returns zero if we are pre term and we are at COURTESY_CALL - tdt owner", async () => {
      await testDeposit.setState(states.COURTESY_CALL)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        owner,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns zero if deposit is pre-term and redeemer is FRT holder", async () => {
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, {
        from: frtHolder,
      })

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns SignerFee if deposit is pre-term, FRT does not exist and signer fee is partially escrowed", async () => {
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })

      await tbtcToken.forceMint(testDeposit.address, signerFee)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it("reverts if deposit is pre-term and redeemer is not Deposit owner", async () => {
      await expectRevert(
        testDeposit.getRedemptionTbtcRequirement.call(accounts[1]),
        "Only TDT holder can redeem unless deposit is at-term or in COURTESY_CALL",
      )
    })

    it("returns full TBTC if we are at-term and caller is not TDT holder", async () => {
      await increaseTime(depositTerm)
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId, {
        from: owner,
      })

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(depositValue)
    })

    it("returns zero if we are at-term, caller is TDT holder, and fee is not escrowed", async () => {
      await increaseTime(depositTerm)
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })
      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns signerFee if we are pre-term, caller is TDT holder and fee is not escrowed", async () => {
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })
      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it("returns zero if we are in courtesy_call, caller is TDT holder and fee is not escrowed", async () => {
      await testDeposit.setState(states.COURTESY_CALL)
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })
      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it("returns zero if we are at-term, caller is TDT holder and signer fee is escrowed", async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await increaseTime(depositTerm)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(
        tdtHolder,
      )
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })
  })

  describe("performRedemptionTBTCTransfers", async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = "0x" + "33".repeat(36)
      valueBytes = "0x1111111111111111"
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock("latest")
      await testDeposit.setRedeemerAddress(owner)
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
      await tbtcToken.resetBalance(depositValue, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, depositValue, {
        from: owner,
      })
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("does nothing if deposit is pre-term, redeemer is FRT holder and signerFee is escrowed", async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, {
        from: frtHolder,
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events.length).to.equal(0)
    })

    it("escrows signerFee if deposit is pre-term, and signerFee is not escrowed", async () => {
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })
      await tbtcToken.resetBalance(signerFee, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, signerFee, {
        from: owner,
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})
      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })

      expect(events[0].returnValues.from).to.equal(owner)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
    })

    it("burns 1 TBTC if deposit is in COURTESY_CALL and TDT holder is the Vending Machine", async () => {
      const {
        receipt: {blockNumber: transferBlock},
      } = await tbtcDepositToken.transferFrom(
        owner,
        vendingMachine.address,
        tdtId,
        {from: owner},
      )
      await testDeposit.setState(states.COURTESY_CALL)

      await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })

      expect(events[0].returnValues.from).to.equal(owner)
      expect(events[0].returnValues.to).to.equal(ZERO_ADDRESS)
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it("escrows signerFee if deposit is pre-term, owner is TDT holder and signer fee is partially escrowed", async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })

      await tbtcToken.resetBalance(signerFee, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, signerFee, {
        from: owner,
      })

      block = await web3.eth.getBlock("latest")

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events[0].returnValues.from).to.equal(owner)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
    })

    it("does nothing if Deposit is in COURTESY_CALL, caller is TDT owner and FRT owner", async () => {
      await testDeposit.setState(states.COURTESY_CALL)
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })

      await tbtcToken.resetBalance(signerFee, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, signerFee, {
        from: owner,
      })
      block = await web3.eth.getBlock("latest")

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events.length).to.equal(0)
    })

    it("escrows fee and sends correct TBTC if Deposit is in COURTESY_CALL and fee is not escrowed", async () => {
      await testDeposit.setState(states.COURTESY_CALL)
      await feeRebateToken.burn(web3.utils.toBN(testDeposit.address), {
        from: frtHolder,
      })
      await testDeposit.setRedeemerAddress(accounts[9])
      await tbtcToken.resetBalance(depositValue, {from: accounts[9]})
      await tbtcToken.resetAllowance(testDeposit.address, depositValue, {
        from: accounts[9],
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: accounts[9]})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })

      expect(events[0].returnValues.from).to.equal(accounts[9])
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
      expect(events[1].returnValues.from).to.equal(accounts[9])
      expect(events[1].returnValues.to).to.equal(owner)
      expect(events[1].returnValues.value).to.eq.BN(depositValue.sub(signerFee))
    })

    it("transfers 1 TBTC to TDT holder if deposit is in COURTESY_CALL and fee is escrowed", async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await testDeposit.setState(states.COURTESY_CALL)
      await testDeposit.setRedeemerAddress(accounts[9])
      await tbtcToken.resetBalance(depositValue, {from: accounts[9]})
      await tbtcToken.resetAllowance(testDeposit.address, depositValue, {
        from: accounts[9],
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: accounts[9]})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events[0].returnValues.from).to.equal(accounts[9])
      expect(events[0].returnValues.to).to.equal(owner)
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it("transfers signerFee if deposit is pre-term and redeemer is not FRT holder", async () => {
      await tbtcToken.resetBalance(signerFee, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, signerFee, {
        from: owner,
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
    })

    it("burns 1 TBTC if deposit is at-term and Deposit Token owner is Vending Machine", async () => {
      await increaseTime(depositTerm)
      await tbtcDepositToken.transferFrom(
        tdtHolder,
        vendingMachine.address,
        tdtId,
        {from: owner},
      )
      await tbtcToken.resetBalance(depositValue, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, depositValue, {
        from: owner,
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events.length).to.equal(1)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(ZERO_ADDRESS)
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it("sends 1 TBTC to Deposit Token owner if deposit is at-term and fee is escrowed", async () => {
      await increaseTime(depositTerm)
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId, {
        from: owner,
      })
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await increaseTime(depositTerm)

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events.length).to.equal(1)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(accounts[1])
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it("escrows fee and sends correct TBTC if Deposit is at-term and fee is not escrowed", async () => {
      await increaseTime(depositTerm)
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId, {
        from: owner,
      })
      await tbtcToken.resetBalance(depositValue, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, depositValue, {
        from: owner,
      })

      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.performRedemptionTBTCTransfers({from: owner})

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      expect(events.length).to.equal(2)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
      expect(events[1].returnValues.from).to.equal(tdtHolder)
      expect(events[1].returnValues.to).to.equal(accounts[1])
      expect(events[1].returnValues.value).to.eq.BN(depositValue.sub(signerFee))
    })
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
      "0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90"
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
      await testDeposit.setUTXOInfo(valueBytes, 0, outpoint)

      // make sure there is sufficient balance to request redemption. Then approve deposit
      await tbtcToken.resetBalance(requiredBalance, {from: owner})
      await tbtcToken.resetAllowance(testDeposit.address, requiredBalance, {
        from: owner,
      })
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("updates state successfully and fires a RedemptionRequested event", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is 2.86331153 BTC
      await testDeposit.requestRedemption(
        "0x1111111100000000",
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
        "0x1111111100000000",
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

    it("escrows the fee rebate reward for the fee rebate token holder", async () => {
      const block = await web3.eth.getBlock("latest")

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)

      // the fee is 2.86331153 BTC
      const {
        receipt: {blockNumber: transferBlock},
      } = await testDeposit.requestRedemption(
        "0x1111111100000000",
        redeemerOutputScript,
        {from: owner},
      )

      const events = await tbtcToken.getPastEvents("Transfer", {
        fromBlock: transferBlock,
        toBlock: "latest",
      })
      const event = events[0]
      expect(event.returnValues.from).to.equal(tdtHolder)
      expect(event.returnValues.to).to.equal(testDeposit.address)
      expect(event.returnValues.value).to.eq.BN(signerFee)
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

    it("reverts if the fee is low", async () => {
      await expectRevert(
        testDeposit.requestRedemption(
          "0x0011111111111111",
          "0x1976a914" + "33".repeat(20) + "88ac",
          {from: owner},
        ),
        "Fee is too low",
      )
    })

    it("reverts if the output script is non-standard", async () => {
      const block = await web3.eth.getBlock("latest")
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)

      await tbtcDepositToken.transferFrom(tdtHolder, frtHolder, tdtId, {
        from: owner,
      })

      await expectRevert(
        testDeposit.requestRedemption(
          "0x1111111100000000",
          "0x" + "33".repeat(20),
          {from: owner},
        ),
        "Output script must be a standard type.",
      )
    })

    it("reverts if the caller is not the deposit owner", async () => {
      const block = await web3.eth.getBlock("latest")
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)

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

  describe("approveDigest", async () => {
    beforeEach(async () => {
      await testDeposit.setSigningGroupPublicKey("0x00", "0x00")
    })

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
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

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
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - feeIncreaseTimer.toNumber()
      await testDeposit.setDigestApprovedAtTime(
        prevSighash,
        withdrawalRequestTime,
      )
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        initialFee,
        withdrawalRequestTime,
        prevSighash,
      )
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
      const block = await web3.eth.getBlock("latest")
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        redeemerOutputScript,
        initialFee,
        block.timestamp,
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
      // Mint the signer fee so we don't try to transfer nonexistent tokens eh.
      await tbtcToken.forceMint(
        testDeposit.address,
        await testDeposit.signerFee(),
      )
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await testDeposit.setUTXOInfo(
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

    it("updates the state, deconstes struct info, calls TBTC and Keep, and emits a Redeemed event", async () => {
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
      expect(requestInfo[0]).to.equal(ZERO_ADDRESS)
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
  })

  describe("redemptionTransactionChecks", async () => {
    beforeEach(async () => {
      await testDeposit.setUTXOInfo(
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
      await testDeposit.setUTXOInfo(
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

  describe("notifySignatureTimeout", async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getSignatureTimeout.call()
    })

    beforeEach(async () => {
      await createSnapshot()
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await ecdsaKeepStub.burnContractBalance()
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_SIGNATURE)
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        0,
        withdrawalRequestTime,
        bytes32zero,
      )
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("reverts if not awaiting redemption signature", async () => {
      testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifySignatureTimeout(),
        "Not currently awaiting a signature",
      )
    })

    it("reverts if the signature timeout has not elapsed", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        0,
        withdrawalRequestTime + timer + 1,
        bytes32zero,
      )
      await expectRevert(
        testDeposit.notifySignatureTimeout(),
        "Signature timer has not elapsed",
      )
    })

    it("reverts if no funds received as signer bond", async () => {
      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(new BN(0), "no bond should be sent").to.eq.BN(bond)

      await expectRevert(
        testDeposit.notifySignatureTimeout(),
        "No funds received, unexpected",
      )
    })

    it("liquidates the deposit and allows redeemer to withdraw signer bond", async () => {
      const signerBonds = new BN("1000000")
      await ecdsaKeepStub.send(signerBonds, {from: tdtHolder})
      await testDeposit.setRedeemerAddress(accounts[1])

      const initialWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: accounts[1],
      })

      await testDeposit.notifySignatureTimeout()
      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(bond, "Bond not seized as expected").to.eq.BN(new BN(0))

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      const finalWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: accounts[1],
      })

      expect(initialWithdrawable).to.eq.BN(new BN("0"))
      expect(finalWithdrawable).to.eq.BN(signerBonds)
      expect(liquidationTime[0], "Auction should not be initiated").to.eq.BN(0)
    })
  })

  describe("notifyRedemptionProofTimeout", async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getRedemptionProofTimeout.call()
    })

    beforeEach(async () => {
      await createSnapshot()
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await ecdsaKeepStub.burnContractBalance()
      await testDeposit.setState(states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        0,
        withdrawalRequestTime,
        bytes32zero,
      )
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("reverts if not awaiting redemption proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyRedemptionProofTimeout(),
        "Not currently awaiting a redemption proof",
      )
    })

    it("reverts if the proof timeout has not elapsed", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        0,
        withdrawalRequestTime * 5,
        bytes32zero,
      )

      await expectRevert(
        testDeposit.notifyRedemptionProofTimeout(),
        "Proof timer has not elapsed",
      )
    })

    it("reverts if no funds recieved as signer bond", async () => {
      await testDeposit.setRequestInfo(
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        0,
        withdrawalRequestTime,
        bytes32zero,
      )
      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(bond, "no bond should be sent").to.eq.BN(new BN(0))

      await expectRevert(
        testDeposit.notifyRedemptionProofTimeout(),
        "No funds received, unexpected",
      )
    })

    it("liquidates the deposit and allows redeemer to withdraw signer bond", async () => {
      const signerBonds = new BN("1000000")
      await ecdsaKeepStub.send(signerBonds, {from: owner})
      await testDeposit.setRedeemerAddress(accounts[1])

      const initialWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: accounts[1],
      })

      await testDeposit.notifyRedemptionProofTimeout()

      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(bond, "Bond not seized as expected").to.eq.BN(new BN(0))

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      const finalWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: accounts[1],
      })

      expect(initialWithdrawable).to.eq.BN("0")
      expect(finalWithdrawable).to.eq.BN(signerBonds)
      expect(liquidationTime[0], "Auction should not be initiated").to.eq.BN(0)
    })
  })
})
beforeEach(async () => {
  await createSnapshot()
})
