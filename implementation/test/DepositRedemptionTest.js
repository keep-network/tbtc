import expectThrow from './helpers/expectThrow'
import increaseTime from './helpers/increaseTime'
import { createSnapshot, restoreSnapshot } from './helpers/snapshot'
import deployTestDeposit from './helpers/deployTestDeposit'

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

contract('DepositRedemption', (accounts) => {
  let tbtcConstants
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

  const tdtHolder = accounts[0]
  const frtHolder = accounts[4]

  before(async () => {
    let deployed
    ({
      tbtcConstants,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      ecdsaKeepStub,
      deployed,
    } = await deployTestDeposit())
    vendingMachine = deployed.VendingMachine

    await testDeposit.setSignerFeeDivisor(new BN('200'))

    await feeRebateToken.forceMint(frtHolder, web3.utils.toBN(testDeposit.address))

    tdtId = await web3.utils.toBN(testDeposit.address)
    await tbtcDepositToken.forceMint(tdtHolder, tdtId)

    depositValue = await testDeposit.lotSizeTbtc.call()
    signerFee = await testDeposit.signerFee.call()
    depositTerm = await tbtcConstants.getDepositTerm.call()
  })

  beforeEach(async () => {
    await testDeposit.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  describe('getOwnerRedemptionTbtcRequirement', async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = '0x' + '33'.repeat(36)
      valueBytes = '0x1111111111111111'
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock('latest')
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('returns signerFee if we are pre-term and owner is not FRT holder', async () => {
      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(accounts[0])
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it('returns zero if deposit is pre-term, owner is FRT holder and signer fee is escrowed', async () => {
      await feeRebateToken.transferFrom(accounts[4], accounts[0], tdtId, { from: accounts[4] })
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(accounts[0])
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it('returns signer fee if deposit is pre-term, owner is FRT holder and signer fee is not escrowed', async () => {
      await feeRebateToken.transferFrom(accounts[4], accounts[0], tdtId, { from: accounts[4] })

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(accounts[0])
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it('returns correct fee if deposit is pre-term, owner is FRT holder and signer fee is partially escrowed', async () => {
      const expectedFee = new BN(100)
      await feeRebateToken.transferFrom(accounts[4], accounts[0], tdtId, { from: accounts[4] })
      await tbtcToken.forceMint(testDeposit.address, signerFee.sub(expectedFee))

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(accounts[0])
      expect(tbtcOwed).to.eq.BN(expectedFee)
    })
  })

  describe('getRedemptionTbtcRequirement', async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = '0x' + '33'.repeat(36)
      valueBytes = '0x1111111111111111'
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock('latest')
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('returns signerFee if we are pre term and redeemer is not FRT holder', async () => {
      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it('returns zero if deposit is pre-term, redeemer is FRT holder and signer fee is escrowed', async () => {
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, { from: frtHolder })
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })

    it('returns full TBTC if we are pre term and we are at COURTESY_CALL', async () => {
      await testDeposit.setState(utils.states.COURTESY_CALL)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(accounts[0])
      expect(tbtcOwed).to.eq.BN(depositValue)
    })

    it('returns signer fee if deposit is pre-term, redeemer is FRT holder and signer fee is not escrowed', async () => {
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, { from: frtHolder })

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it('returns correct fee if deposit is pre-term, owner is FRT holder and signer fee is partially escrowed', async () => {
      const expectedFee = new BN(100)
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, { from: frtHolder })
      await tbtcToken.forceMint(testDeposit.address, signerFee.sub(expectedFee))

      const tbtcOwed = await testDeposit.getOwnerRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(expectedFee)
    })

    it('reverts if deposit is pre-term and redeemer is not Deposit owner', async () => {
      await expectThrow(
        testDeposit.getRedemptionTbtcRequirement.call(accounts[1]),
        'Only TDT owner can redeem unless deposit is at-term or in COURTESY_CALL'
      )
    })

    it('returns full TBTC if we are at-term and caller is not TDT owner', async () => {
      await increaseTime(depositTerm.toNumber())
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId)

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(depositValue)
    })

    it('returns SignerFee if we are at-term, caller is TDT owner, and fee is not escrowed', async () => {
      await increaseTime(depositTerm.toNumber())

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(signerFee)
    })

    it('returns zero if we are at-term, caller is TDT owner and signer fee is escrowed', async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await increaseTime(depositTerm.toNumber())

      const tbtcOwed = await testDeposit.getRedemptionTbtcRequirement.call(tdtHolder)
      expect(tbtcOwed).to.eq.BN(new BN(0))
    })
  })

  describe('performRedemptionTBTCTransfers', async () => {
    let outpoint
    let valueBytes
    let block
    before(async () => {
      outpoint = '0x' + '33'.repeat(36)
      valueBytes = '0x1111111111111111'
    })

    beforeEach(async () => {
      await createSnapshot()
      block = await web3.eth.getBlock('latest')
      await testDeposit.setRedeemerAddress(accounts[0])
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
      await tbtcToken.resetBalance(depositValue)
      await tbtcToken.resetAllowance(testDeposit.address, depositValue)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('does nothing if deposit is pre-term, redeemer is FRT holder and signerFee is escrowed', async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await feeRebateToken.transferFrom(frtHolder, tdtHolder, tdtId, { from: frtHolder })

      block = await web3.eth.getBlock('latest')
      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      assert.equal(events.length, 0)
    })

    it('escrows signerFee if deposit is pre-term, redeemer is FRT holder and signerFee is not escrowed', async () => {
      await feeRebateToken.transferFrom(accounts[4], accounts[0], tdtId, { from: accounts[4] })
      await tbtcToken.resetBalance(signerFee)
      await tbtcToken.resetAllowance(testDeposit.address, signerFee)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()
      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })

      expect(events[0].returnValues.from).to.equal(accounts[0])
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
    })

    it('burns 1 TBTC if deposit is in COURTESY_CALL and TDT owner is the Vending Machine', async () => {
      await tbtcDepositToken.transferFrom(accounts[0], vendingMachine.address, tdtId)
      block = await web3.eth.getBlock('latest')
      await testDeposit.setState(utils.states.COURTESY_CALL)

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })

      expect(events[0].returnValues.from).to.equal(accounts[0])
      expect(events[0].returnValues.to).to.equal(utils.address0)
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it('escrows correct fee if deposit is pre-term, owner is TDT holder and signer fee is partially escrowed', async () => {
      const expectedFee = new BN(100)
      await tbtcToken.forceMint(testDeposit.address, signerFee.sub(expectedFee))
      await feeRebateToken.transferFrom(accounts[4], accounts[0], tdtId, { from: accounts[4] })
      await tbtcToken.resetBalance(expectedFee)
      await tbtcToken.resetAllowance(testDeposit.address, expectedFee)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      expect(events[0].returnValues.from).to.equal(accounts[0])
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(expectedFee)
    })

    it('escrows fee and sends correct TBTC if Deposit is in COURTESY_CALL and fee is not escrowed', async () => {
      block = await web3.eth.getBlock('latest')
      await testDeposit.setState(utils.states.COURTESY_CALL)

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })

      expect(events[0].returnValues.from).to.equal(accounts[0])
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
      expect(events[1].returnValues.from).to.equal(accounts[0])
      expect(events[1].returnValues.to).to.equal(accounts[0])
      expect(events[1].returnValues.value).to.eq.BN(depositValue.sub(signerFee))
    })

    it('transfers 1 TBTC to TDT owner if deposit is in COURTESY_CALL and fee is escrowed', async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await testDeposit.setState(utils.states.COURTESY_CALL)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      expect(events[0].returnValues.from).to.equal(accounts[0])
      expect(events[0].returnValues.to).to.equal(accounts[0])
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it('transfers signerFee if deposit is pre-term and redeemer is not FRT holder', async () => {
      await tbtcToken.resetBalance(signerFee)
      await tbtcToken.resetAllowance(testDeposit.address, signerFee)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
    })

    it('burns 1 TBTC if deposit is at-term and Deposit Token owner is Vending Machine', async () => {
      await increaseTime(depositTerm.toNumber())
      await tbtcDepositToken.transferFrom(tdtHolder, vendingMachine.address, tdtId)
      await tbtcToken.resetBalance(depositValue)
      await tbtcToken.resetAllowance(testDeposit.address, depositValue)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      assert.equal(events.length, 1)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(utils.address0)
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it('sends 1 TBTC to Deposit Token owner if deposit is at-term and fee is escrowed', async () => {
      await increaseTime(depositTerm.toNumber())
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId)
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await increaseTime(depositTerm.toNumber())
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      assert.equal(events.length, 1)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(accounts[1])
      expect(events[0].returnValues.value).to.eq.BN(depositValue)
    })

    it('escrows fee and sends correct TBTC if Deposit is at-term and fee is not escrowed', async () => {
      await increaseTime(depositTerm.toNumber())
      await tbtcDepositToken.transferFrom(tdtHolder, accounts[1], tdtId)
      await tbtcToken.resetBalance(depositValue)
      await tbtcToken.resetAllowance(testDeposit.address, depositValue)
      block = await web3.eth.getBlock('latest')

      await testDeposit.performRedemptionTBTCTransfers()

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      assert.equal(events.length, 2)
      expect(events[0].returnValues.from).to.equal(tdtHolder)
      expect(events[0].returnValues.to).to.equal(testDeposit.address)
      expect(events[0].returnValues.value).to.eq.BN(signerFee)
      expect(events[1].returnValues.from).to.equal(tdtHolder)
      expect(events[1].returnValues.to).to.equal(accounts[1])
      expect(events[1].returnValues.value).to.eq.BN(depositValue.sub(signerFee))
    })
  })

  describe('requestRedemption', async () => {
    // the TX produced will be:
    // 01000000000101333333333333333333333333333333333333333333333333333333333333333333333333000000000001111111110000000016001433333333333333333333333333333333333333330000000000
    // the signer pkh will be:
    // 5eb9b5e445db673f0ed8935d18cd205b214e5187
    // == hash160(023333333333333333333333333333333333333333333333333333333333333333)
    // the sighash preimage will be:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788ac111111111111111100000000e4ca7a168bd64e3123edd7f39e1ab7d670b32311cac2dda8e083822139c7936c0000000001000000
    const sighash = '0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90'
    const outpoint = '0x' + '33'.repeat(36)
    const valueBytes = '0x1111111111111111'
    const keepPubkeyX = '0x' + '33'.repeat(32)
    const keepPubkeyY = '0x' + '44'.repeat(32)
    const redeemerPKH = '0x' + '33'.repeat(20)
    let requiredBalance

    before(async () => {
      requiredBalance = depositValue
    })

    beforeEach(async () => {
      await createSnapshot()
      await testDeposit.setState(utils.states.ACTIVE)
      await testDeposit.setUTXOInfo(valueBytes, 0, outpoint)

      // make sure there is sufficient balance to request redemption. Then approve deposit
      await tbtcToken.resetBalance(requiredBalance)
      await tbtcToken.resetAllowance(testDeposit.address, requiredBalance)
      await ecdsaKeepStub.setSuccess(true)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('updates state successfully and fires a RedemptionRequested event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is ~12,297,829,380 BTC
      await testDeposit.requestRedemption('0x1111111100000000', redeemerPKH)

      const requestInfo = await testDeposit.getRequestInfo()
      assert.equal(requestInfo[1], redeemerPKH)
      assert(!requestInfo[3].eqn(0)) // withdrawalRequestTime is set
      assert.equal(requestInfo[4], sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, sighash)
    })

    it('updates state successfully and fires a RedemptionRequested event from COURTESY_CALL state', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setState(utils.states.COURTESY_CALL)

      // the fee is ~12,297,829,380 BTC
      await testDeposit.requestRedemption('0x1111111100000000', redeemerPKH)

      const requestInfo = await testDeposit.getRequestInfo()
      assert.equal(requestInfo[1], redeemerPKH)
      assert(!requestInfo[3].eqn(0)) // withdrawalRequestTime is set
      assert.equal(requestInfo[4], sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, sighash)
    })

    it('escrows the fee rebate reward for the fee rebate token holder', async () => {
      const block = await web3.eth.getBlock('latest')

      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)

      // the fee is ~12,297,829,380 BTC
      await testDeposit.requestRedemption('0x1111111100000000', redeemerPKH)

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock: block.number, toBlock: 'latest' })
      const event = events[0]
      expect(event.returnValues.from).to.equal(tdtHolder)
      expect(event.returnValues.to).to.equal(testDeposit.address)
      expect(event.returnValues.value).to.eq.BN(signerFee)
    })

    it('reverts if not in Active or Courtesy', async () => {
      await testDeposit.setState(utils.states.LIQUIDATED)

      await expectThrow(
        testDeposit.requestRedemption('0x1111111100000000', '0x' + '33'.repeat(20)),
        'Redemption only available from Active or Courtesy state'
      )
    })

    it('reverts if the fee is low', async () => {
      await expectThrow(
        testDeposit.requestRedemption('0x0011111111111111', '0x' + '33'.repeat(20)),
        'Fee is too low'
      )
    })

    it('reverts if the caller is not the deposit owner', async () => {
      const block = await web3.eth.getBlock('latest')
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)

      await tbtcDepositToken.transferFrom(tdtHolder, frtHolder, tdtId)

      await expectThrow(
        testDeposit.requestRedemption('0x1111111100000000', '0x' + '33'.repeat(20)),
        'Only TDT owner can redeem unless deposit is at-term or in COURTESY_CALL'
      )
    })
  })

  describe('approveDigest', async () => {
    beforeEach(async () => {
      await testDeposit.setSigningGroupPublicKey('0x00', '0x00')
    })

    it('calls keep for signing', async () => {
      const digest = '0x' + '08'.repeat(32)

      await testDeposit.approveDigest(digest)
        .catch((err) => {
          assert.fail(`cannot approve digest: ${err}`)
        })

      const blockNumber = await web3.eth.getBlock('latest').number

      // Check if ECDSAKeep has been called and event emitted.
      const eventList = await ecdsaKeepStub.getPastEvents(
        'SignatureRequested',
        { fromBlock: blockNumber, toBlock: 'latest' },
      )

      assert.equal(
        eventList[0].returnValues.digest,
        digest,
        'incorrect digest in emitted event',
      )
    })

    it('registers timestamp for digest approval', async () => {
      const digest = '0x' + '02'.repeat(32)

      const approvalTx = await testDeposit.approveDigest(digest)
        .catch((err) => {
          assert.fail(`cannot approve digest: ${err}`)
        })

      const block = await web3.eth.getBlock(approvalTx.receipt.blockNumber)
      const expectedTimestamp = block.timestamp

      const timestamp = await testDeposit.wasDigestApprovedForSigning(digest)
        .catch((err) => {
          assert.fail(`cannot check digest approval: ${err}`)
        })

      assert.equal(
        timestamp,
        expectedTimestamp,
        'incorrect registered timestamp',
      )
    })
  })

  describe('provideRedemptionSignature', async () => {
    // signing the sha 256 of '11' * 32
    // signing with privkey '11' * 32
    // using RFC 6979 nonce (libsecp256k1)
    const pubkeyX = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa'
    const pubkeyY = '0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
    const digest = '0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc'
    const v = 27
    const r = '0xd7e83e8687ba8b555f553f22965c74e81fd08b619a7337c5c16e4b02873b537e'
    const s = '0x633bf745cdf7ae303ca8a6f41d71b2c3a21fcbd1aed9e7ffffa295c08918c1b3'

    beforeEach(async () => {
      await testDeposit.setState(utils.states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it('updates the state and logs GotRedemptionSignature', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testDeposit.setSigningGroupPublicKey(pubkeyX, pubkeyY)
      await testDeposit.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 0, 0, digest)

      await testDeposit.provideRedemptionSignature(v, r, s)

      const state = await testDeposit.getState.call()
      expect(state).to.eq.BN(utils.states.AWAITING_WITHDRAWAL_PROOF)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents('GotRedemptionSignature', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._r, r)
      assert.equal(eventList[0].returnValues._s, s)
    })

    it('errors if not awaiting withdrawal signature', async () => {
      await testDeposit.setState(utils.states.START)

      await expectThrow(
        testDeposit.provideRedemptionSignature(v, r, s),
        'Not currently awaiting a signature'
      )
    })

    it('errors on invaid sig', async () => {
      await expectThrow(
        testDeposit.provideRedemptionSignature(28, r, s),
        'Invalid signature'
      )
    })
  })

  describe('increaseRedemptionFee', async () => {
    // the signer pkh will be:
    // 5eb9b5e445db673f0ed8935d18cd205b214e5187
    // == hash160(023333333333333333333333333333333333333333333333333333333333333333)
    // prevSighash preimage is:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788acffffffffffffffff0000000077feba7d287431f43e61c3e716e92d266658eb223ebad0e64d83b747652e2b9b0000000001000000
    // nextSighash preimage is:
    // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788acffffffffffffffff0000000044bf045101f1d83d0f2e017a73bb85857f25137c31cc7382ef363c909659f55a0000000001000000

    const prevSighash = '0xd94b6f3bf19147cc3305ef202d6bd64f9b9a12d4d19cc2d8c7f93ef58fc8fffe'
    const nextSighash = '0xbb56d80cfd71e90215c6b5200c0605b7b80689d3479187bc2c232e756033a560'
    const keepPubkeyX = '0x' + '33'.repeat(32)
    const keepPubkeyY = '0x' + '44'.repeat(32)
    const prevoutValueBytes = '0xffffffffffffffff'
    const previousOutputBytes = '0x0000ffffffffffff'
    const newOutputBytes = '0x0100feffffffffff'
    const initialFee = 0xffff
    const outpoint = '0x' + '33'.repeat(36)
    const redeemerPKH = '0x' + '33'.repeat(20)
    let feeIncreaseTimer

    before(async () => {
      feeIncreaseTimer = await tbtcConstants.getIncreaseFeeTimer.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - feeIncreaseTimer.toNumber()
      await testDeposit.setDigestApprovedAtTime(prevSighash, withdrawalRequestTime)
      await testDeposit.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setRequestInfo(utils.address0, redeemerPKH, initialFee, withdrawalRequestTime, prevSighash)
      await ecdsaKeepStub.setSuccess(true)
    })

    it('approves a new digest for signing, updates the state, and logs RedemptionRequested', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes)

      const requestInfo = await testDeposit.getRequestInfo.call()
      assert.equal(requestInfo[4], nextSighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, nextSighash)
    })

    it('reverts if not awaiting withdrawal proof', async () => {
      await testDeposit.setState(utils.states.START)

      await expectThrow(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        'Fee increase only available after signature provided'
      )
    })

    it('reverts if the increase fee timer has not elapsed', async () => {
      const block = await web3.eth.getBlock('latest')
      await testDeposit.setRequestInfo(utils.address0, redeemerPKH, initialFee, block.timestamp, prevSighash)

      await expectThrow(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        'Fee increase not yet permitted'
      )
    })

    it('reverts if the fee step is not linear', async () => {
      await expectThrow(
        testDeposit.increaseRedemptionFee(previousOutputBytes, '0x1101010101102201'),
        'Not an allowed fee step'
      )
    })

    it('reverts if the previous sighash was not the latest approved', async () => {
      await testDeposit.setRequestInfo(utils.address0, redeemerPKH, initialFee, withdrawalRequestTime, keepPubkeyX)

      // Previous sigHash is not approved for signing.
      await testDeposit.setDigestApprovedAtTime(prevSighash, 0)

      await expectThrow(
        testDeposit.increaseRedemptionFee(previousOutputBytes, newOutputBytes),
        'Provided previous value does not yield previous sighash'
      )
    })
  })

  describe('provideRedemptionProof', async () => {
    // real tx from mainnet bitcoin
    const currentDiff = 6353030562983
    // const txid = '0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f'
    const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c'
    const _version = '0x01000000'
    const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
    const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
    const _txLocktime = '0x4ec10800'
    const proof = '0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049'
    const index = 129
    const headerChain = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
    const outpoint = '0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000'
    const prevoutValueBytes = '0xf078351d00000000'
    const redeemerPKH = '0x86e7303082a6a21d5837176bc808bf4828371ab6'

    beforeEach(async () => {
      await tbtcSystemStub.setCurrentDiff(currentDiff)
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setRequestInfo('0x' + '11'.repeat(20), redeemerPKH, 14544, 0, '0x' + '11' * 32)
    })

    it('updates the state, deconstes struct info, calls TBTC and Keep, and emits a Redeemed event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testDeposit.provideRedemptionProof(_version, _txInputVector, _txOutputVector, _txLocktime, proof, index, headerChain)

      const depositState = await testDeposit.getState.call()
      assert(depositState.eq(utils.states.REDEEMED))

      const requestInfo = await testDeposit.getRequestInfo.call()
      assert.equal(requestInfo[0], '0x' + '11'.repeat(20)) // address is intentionally not cleared
      assert.equal(requestInfo[1], utils.address0)
      assert.equal(requestInfo[4], utils.bytes32zero)

      const eventList = await tbtcSystemStub.getPastEvents('Redeemed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._txid, txidLE)
    })

    it('reverts if not in the redemption flow', async () => {
      await testDeposit.setState(utils.states.START)

      await expectThrow(
        testDeposit.provideRedemptionProof(_version, _txInputVector, _txOutputVector, _txLocktime, proof, index, headerChain),
        'Redemption proof only allowed from redemption flow'
      )
    })

    it('reverts if the merkle proof is not validated successfully', async () => {
      await expectThrow(
        testDeposit.provideRedemptionProof(_version, _txInputVector, _txOutputVector, _txLocktime, proof, 0, headerChain),
        'Tx merkle proof is not valid for provided header'
      )
    })
  })

  describe('redemptionTransactionChecks', async () => {
    const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
    const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
    const outputValue = 490029088
    const outpoint = '0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000'
    const prevoutValueBytes = '0xf078351d00000000'
    const redeemerPKH = '0x86e7303082a6a21d5837176bc808bf4828371ab6'

    beforeEach(async () => {
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testDeposit.setRequestInfo('0x' + '11'.repeat(20), redeemerPKH, 14544, 0, '0x' + '11' * 32)
    })

    it('returns the output value', async () => {
      const redemptionChecks = await testDeposit.redemptionTransactionChecks.call(_txInputVector, _txOutputVector)
      assert.equal(redemptionChecks, outputValue)
    })

    it('reverts if bad input vector is provided', async () => {
      await expectThrow(
        testDeposit.redemptionTransactionChecks('0x00', _txOutputVector),
        'invalid input vector provided'
      )
    })

    it('reverts if bad output vector is provided', async () => {
      await expectThrow(
        testDeposit.redemptionTransactionChecks(_txInputVector, '0x00'),
        'invalid output vector provided'
      )
    })

    it('reverts if the tx spends the wrong utxo', async () => {
      await testDeposit.setUTXOInfo(prevoutValueBytes, 0, '0x' + '33'.repeat(36))

      await expectThrow(
        testDeposit.redemptionTransactionChecks.call(_txInputVector, _txOutputVector),
        'Tx spends the wrong UTXO'
      )
    })

    it('reverts if the tx sends value to the wrong pkh', async () => {
      await testDeposit.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 14544, 0, '0x' + '11' * 32)

      await expectThrow(
        testDeposit.redemptionTransactionChecks.call(_txInputVector, _txOutputVector),
        'Tx sends value to wrong pubkeyhash'
      )
    })
  })

  describe('notifySignatureTimeout', async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getSignatureTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await ecdsaKeepStub.burnContractBalance()
      await testDeposit.setState(utils.states.AWAITING_WITHDRAWAL_SIGNATURE)
      await testDeposit.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
    })

    it('reverts if not awaiting redemption signature', async () => {
      testDeposit.setState(utils.states.START)

      await expectThrow(
        testDeposit.notifySignatureTimeout(),
        'Not currently awaiting a signature'
      )
    })

    it('reverts if the signature timeout has not elapsed', async () => {
      await testDeposit.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime + timer + 1, utils.bytes32zero)
      await expectThrow(
        testDeposit.notifySignatureTimeout(),
        'Signature timer has not elapsed'
      )
    })

    it('reverts if no funds received as signer bond', async () => {
      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      assert.equal(0, bond, 'no bond should be sent')

      await expectThrow(
        testDeposit.notifySignatureTimeout(),
        'No funds received, unexpected'
      )
    })

    it('starts abort liquidation', async () => {
      await ecdsaKeepStub.send(1000000, { from: tdtHolder })
      await testDeposit.notifySignatureTimeout()

      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyRedemptionProofTimeout', async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getRedemptionProofTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await ecdsaKeepStub.burnContractBalance()
      await testDeposit.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testDeposit.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
    })

    it('reverts if not awaiting redemption proof', async () => {
      await testDeposit.setState(utils.states.START)

      await expectThrow(
        testDeposit.notifyRedemptionProofTimeout(),
        'Not currently awaiting a redemption proof'
      )
    })

    it('reverts if the proof timeout has not elapsed', async () => {
      await testDeposit.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime * 5, utils.bytes32zero)

      await expectThrow(
        testDeposit.notifyRedemptionProofTimeout(),
        'Proof timer has not elapsed'
      )
    })

    it('reverts if no funds recieved as signer bond', async () => {
      await testDeposit.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      assert.equal(0, bond, 'no bond should be sent')

      await expectThrow(
        testDeposit.notifyRedemptionProofTimeout(),
        'No funds received, unexpected'
      )
    })

    it('starts abort liquidation', async () => {
      await ecdsaKeepStub.send(1000000, { from: accounts[0] })
      await testDeposit.notifyRedemptionProofTimeout()

      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })
})
