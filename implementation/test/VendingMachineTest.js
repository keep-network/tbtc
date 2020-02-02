const FundingScript = artifacts.require('FundingScript')
const RedemptionScript = artifacts.require('RedemptionScript')

import expectThrow from './helpers/expectThrow'
import { AssertBalance } from './helpers/assertBalance'
import { createSnapshot, restoreSnapshot } from './helpers/snapshot'
import deployTestDeposit from './helpers/deployTestDeposit'

import BN from 'bn.js'
import utils from './utils'
import chai, { expect } from 'chai'
import bnChai from 'bn-chai'
chai.use(bnChai(BN))

const currentDifficulty = 6353030562983
const _version = '0x01000000'
const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
const _fundingOutputIndex = 0
const _txLocktime = '0x4ec10800'
const _txIndexInBlock = 129
const _bitcoinHeaders = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
const _signerPubkeyX = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e'
const _signerPubkeyY = '0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'
const _merkleProof = '0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049'
const _expectedUTXOoutpoint = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c00000000'
// const _outputValue = 490029088;
const _outValueBytes = '0x2040351d00000000'

contract('VendingMachine', (accounts) => {
  let vendingMachine

  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let feeRebateToken
  let testDeposit
  let ecdsaKeepStub

  let assertBalance
  let tdtId
  // this is the amount of TBTC exchanged for a TDT.
  let depositValue
  let signerFee

  const signerFeeDivisor = new BN('200')

  before(async () => {
    let deployed
    ({
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      ecdsaKeepStub,
      deployed,
    } = await deployTestDeposit())
    vendingMachine = deployed.VendingMachine

    assertBalance = new AssertBalance(tbtcToken)

    await testDeposit.setSignerFeeDivisor(signerFeeDivisor)

    tdtId = await web3.utils.toBN(testDeposit.address)

    depositValue = await testDeposit.lotSizeTbtc()
    signerFee = depositValue.div(signerFeeDivisor)
  })

  describe('#isQualified', async () => {
    it('returns true if deposit is in ACTIVE State', async () => {
      await testDeposit.setState(utils.states.ACTIVE)
      const qualified = await vendingMachine.isQualified.call(testDeposit.address)
      expect(qualified).to.be.true
    })

    it('returns false if deposit is not in ACTIVE State', async () => {
      await testDeposit.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      const qualified = await vendingMachine.isQualified.call(testDeposit.address)
      expect(qualified).to.be.false
    })
  })

  describe('#tdtToTbtc', async () => {
    before(async () => {
      await testDeposit.setState(utils.states.ACTIVE)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('converts TDT to TBTC', async () => {
      await tbtcDepositToken.forceMint(accounts[0], tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, { from: accounts[0] })

      await vendingMachine.tdtToTbtc(tdtId)

      await assertBalance.tbtc(accounts[0], depositValue.sub(signerFee))
    })

    it('mints full lot size if backing deposit has signer fee escrowed', async () => {
      await tbtcToken.forceMint(testDeposit.address, signerFee)

      await tbtcDepositToken.forceMint(accounts[0], tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, { from: accounts[0] })

      await vendingMachine.tdtToTbtc(tdtId)

      await assertBalance.tbtc(accounts[0], depositValue)
    })

    it('fails if deposit not qualified', async () => {
      await testDeposit.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectThrow(
        vendingMachine.tdtToTbtc(tdtId),
        'Deposit must be qualified'
      )
    })

    it(`fails if TDT doesn't exist`, async () => {
      await expectThrow(
        vendingMachine.tdtToTbtc(new BN(123345)),
        'tBTC Deposit Token does not exist'
      )
    })

    it(`fails if TDT transfer not approved`, async () => {
      await tbtcDepositToken.forceMint(accounts[0], tdtId)

      await expectThrow(
        vendingMachine.tdtToTbtc(tdtId),
        'ERC721: transfer caller is not owner nor approved.'
      )
    })
  })

  describe('#tbtcToTdt', async () => {
    before(async () => {
      await testDeposit.setState(utils.states.ACTIVE)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('converts TBTC to TDT', async () => {
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
      await tbtcToken.forceMint(accounts[0], depositValue)
      await tbtcToken.approve(vendingMachine.address, depositValue, { from: accounts[0] })

      const fromBlock = await web3.eth.getBlockNumber()
      await vendingMachine.tbtcToTdt(tdtId)

      const events = await tbtcToken.getPastEvents('Transfer', { fromBlock, toBlock: 'latest' })
      const tbtcBurntEvent = events[0]
      expect(tbtcBurntEvent.returnValues.from).to.equal(accounts[0])
      expect(tbtcBurntEvent.returnValues.to).to.equal(utils.address0)
      expect(tbtcBurntEvent.returnValues.value).to.equal(depositValue.toString())

      expect(
        await tbtcDepositToken.ownerOf(tdtId)
      ).to.equal(accounts[0])
    })

    it('fails if deposit not qualified', async () => {
      await testDeposit.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectThrow(
        vendingMachine.tbtcToTdt(tdtId),
        'Deposit must be qualified'
      )
    })

    it(`fails if caller hasn't got enough TBTC`, async () => {
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)

      await expectThrow(
        vendingMachine.tbtcToTdt(tdtId),
        'Not enough TBTC for TDT exchange.'
      )
    })

    it(`fails if TDT doesn't exist`, async () => {
      await expectThrow(
        vendingMachine.tdtToTbtc(new BN(123345)),
        'tBTC Deposit Token does not exist'
      )
    })

    it(`fails if deposit is locked`, async () => {
      // Deposit is locked if the tBTC Deposit Token is not owned by the vending machine
      const depositOwner = accounts[1]
      await tbtcDepositToken.forceMint(depositOwner, tdtId)
      await tbtcToken.forceMint(accounts[0], depositValue)
      await tbtcToken.approve(vendingMachine.address, depositValue, { from: accounts[0] })

      await expectThrow(
        vendingMachine.tbtcToTdt(tdtId),
        'Deposit is locked.'
      )
    })
  })

  describe('#unqualifiedDepositToTbtc', async () => {
    before(async () => {
      await tbtcSystemStub.setCurrentDiff(currentDifficulty)
      await testDeposit.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(_signerPubkeyX, _signerPubkeyY)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('qualifies a Deposit', async () => {
      await tbtcDepositToken.forceMint(accounts[0], tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, { from: accounts[0] })
      const blockNumber = await web3.eth.getBlock('latest').number

      await vendingMachine.unqualifiedDepositToTbtc(testDeposit.address, _version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      assert.equal(UTXOInfo[0], _outValueBytes)
      assert.equal(UTXOInfo[2], _expectedUTXOoutpoint)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      assert(signingGroupRequestedAt.eqn(0), 'signingGroupRequestedAt not updated')

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      assert(fundingProofTimerStart.eqn(0), 'fundingProofTimerStart not updated')

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(utils.states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents('Funded', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('mints TBTC to the TDT owner and siger fee to Deposit', async () => {
      await tbtcDepositToken.forceMint(accounts[0], tdtId)
      await tbtcDepositToken.approve(vendingMachine.address, tdtId, { from: accounts[0] })

      await vendingMachine.unqualifiedDepositToTbtc(testDeposit.address, _version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      await assertBalance.tbtc(accounts[0], depositValue.sub(signerFee))
      await assertBalance.tbtc(testDeposit.address, signerFee)
    })
  })

  describe('#tbtcToBtc', async () => {
    const sighash = '0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90'
    const outpoint = '0x' + '33'.repeat(36)
    const valueBytes = '0x1111111111111111'
    const keepPubkeyX = '0x' + '33'.repeat(32)
    const keepPubkeyY = '0x' + '44'.repeat(32)
    const requesterPKH = '0x' + '33'.repeat(20)
    let requiredBalance
    let block

    before(async () => {
      requiredBalance = await testDeposit.lotSizeTbtc.call()

      block = await web3.eth.getBlock('latest')
      await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
      await tbtcToken.resetBalance(requiredBalance)
      await tbtcToken.resetAllowance(vendingMachine.address, requiredBalance)
      await ecdsaKeepStub.setSuccess(true)
      await testDeposit.setState(utils.states.ACTIVE)
      await testDeposit.setUTXOInfo(valueBytes, block.timestamp, outpoint)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('successfully redeems via wrapper', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await tbtcToken.forceMint(testDeposit.address, signerFee)
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is ~12,297,829,380 BTC
      await feeRebateToken.forceMint(accounts[0], tdtId)
      await vendingMachine.tbtcToBtc(testDeposit.address, '0x1111111100000000', requesterPKH, accounts[0])
      const requestInfo = await testDeposit.getRequestInfo()

      assert.equal(requestInfo[1], requesterPKH)
      assert(!requestInfo[3].eqn(0)) // withdrawalRequestTime is set
      assert.equal(requestInfo[4], sighash)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, sighash)
    })

    it('fails to redeem with insufficient balance', async () => {
      await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

      // the fee is ~12,297,829,380 BTC
      // requester does not own the FRT, and therefore owes an additional SignerFee
      await feeRebateToken.forceMint(accounts[1], tdtId)

      await expectThrow(
        vendingMachine.tbtcToBtc(testDeposit.address, '0x1111111100000000', requesterPKH, accounts[0]),
        'SafeMath: subtraction overflow.'
      )
    })

    describe('RedemptionScript', async () => {
      let redemptionScript

      before(async () => {
        redemptionScript = await RedemptionScript.new(vendingMachine.address, tbtcToken.address, feeRebateToken.address)
      })

      beforeEach(async () => {
        await createSnapshot()
      })

      afterEach(async () => {
        await restoreSnapshot()
      })

      it('successfully requests redemption', async () => {
        await testDeposit.setState(utils.states.ACTIVE)
        await tbtcDepositToken.forceMint(vendingMachine.address, tdtId)
        await tbtcToken.forceMint(accounts[0], depositValue.add(signerFee))
        await feeRebateToken.forceMint(accounts[0], tdtId)

        const blockNumber = await web3.eth.getBlock('latest').number

        await testDeposit.setSigningGroupPublicKey(keepPubkeyX, keepPubkeyY)

        const tbtcToBtc = vendingMachine.abi.filter((x) => x.name == 'tbtcToBtc')[0]
        const calldata = web3.eth.abi.encodeFunctionCall(
          tbtcToBtc, [testDeposit.address, '0x1111111100000000', requesterPKH, accounts[0]]
        )

        await tbtcToken.approveAndCall(
          redemptionScript.address,
          depositValue.add(signerFee),
          calldata
        )

        const requestInfo = await testDeposit.getRequestInfo()
        assert.equal(requestInfo[1], requesterPKH)
        assert(!requestInfo[3].eqn(0)) // withdrawalRequestTime is set
        assert.equal(requestInfo[4], sighash)

        // fired an event
        const eventList = await tbtcSystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
        assert.equal(eventList[0].returnValues._digest, sighash)
      })

      it('reverts for unknown function calls encoded in _extraData', async () => {
        const unknownFunctionSignature = '0xCAFEBABE'
        await tbtcToken.forceMint(accounts[0], depositValue.add(signerFee))

        await expectThrow(
          tbtcToken.approveAndCall(
            redemptionScript.address,
            depositValue.add(signerFee),
            unknownFunctionSignature
          ),
          'Bad _extraData signature. Call must be to tbtcToBtc.'
        )
      })
    })
  })

  describe('FundingScript', async () => {
    let fundingScript

    before(async () => {
      await tbtcSystemStub.setCurrentDiff(currentDifficulty)
      await testDeposit.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(_signerPubkeyX, _signerPubkeyY)
      await tbtcToken.zeroBalance(accounts[0])
      await tbtcDepositToken.forceMint(accounts[0], tdtId)
      fundingScript = await FundingScript.new(vendingMachine.address, tbtcToken.address, tbtcDepositToken.address, feeRebateToken.address)
    })

    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it('calls unqualifiedDepositToTbtcABI', async () => {
      const unqualifiedDepositToTbtcABI = vendingMachine.abi.filter((x) => x.name == 'unqualifiedDepositToTbtc')[0]
      const calldata = web3.eth.abi.encodeFunctionCall(
        unqualifiedDepositToTbtcABI,
        [testDeposit.address, _version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders]
      )

      await tbtcDepositToken.approveAndCall(
        fundingScript.address,
        tdtId,
        calldata
      )

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      assert.equal(UTXOInfo[0], _outValueBytes)
      assert.equal(UTXOInfo[2], _expectedUTXOoutpoint)

      await assertBalance.tbtc(accounts[0], depositValue.sub(signerFee))
      expect(await tbtcDepositToken.ownerOf(tdtId)).to.equal(vendingMachine.address)
      expect(await feeRebateToken.ownerOf(tdtId)).to.equal(accounts[0])
    })

    it('reverts for unknown function calls encoded in _extraData', async () => {
      const unknownFunctionSignature = '0xCAFEBABE'

      await expectThrow(
        tbtcDepositToken.approveAndCall(
          fundingScript.address,
          tdtId,
          unknownFunctionSignature
        ),
        'Bad _extraData signature. Call must be to unqualifiedDepositToTbtc.'
      )
    })
  })
})
