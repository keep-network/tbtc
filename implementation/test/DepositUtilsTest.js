const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')

const KeepStub = artifacts.require('KeepStub')
const TBTCStub = artifacts.require('TBTCStub')
const SystemStub = artifacts.require('SystemStub')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDepositUtils = artifacts.require('TestDepositUtils')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_UTILS_DEPLOY = [
  { name: 'BytesLib', contract: BytesLib },
  { name: 'BTCUtils', contract: BTCUtils },
  { name: 'ValidateSPV', contract: ValidateSPV },
  { name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name
  { name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'TestDepositUtils', contract: TestDepositUtils },
  { name: 'KeepStub', contract: KeepStub },
  { name: 'TBTCStub', contract: TBTCStub },
  { name: 'SystemStub', contract: SystemStub }]


contract('DepositUtils', (accounts) => {
  let deployed
  let testUtilsInstance
  before(async () => {
    deployed = await utils.deploySystem(TEST_DEPOSIT_UTILS_DEPLOY)
    testUtilsInstance = deployed.TestDepositUtils

    await testUtilsInstance.createNewDeposit(
      deployed.SystemStub.address,
      deployed.TBTCStub.address,
      deployed.KeepStub.address,
      1, // m
      1) // n
  })

  describe('currentBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      const blockDifficulty = await testUtilsInstance.currentBlockDifficulty.call()
      assert(blockDifficulty.eq(new BN(1)))

      await deployed.SystemStub.setCurrentDiff(33)
      const newBlockDifficulty = await testUtilsInstance.currentBlockDifficulty.call()
      assert(newBlockDifficulty.eq(new BN(33)))
    })
  })

  describe('previousBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      const blockDifficulty = await testUtilsInstance.previousBlockDifficulty.call()
      assert(blockDifficulty.eq(new BN(1)))

      await deployed.SystemStub.setPreviousDiff(44)
      const newBlockDifficulty = await testUtilsInstance.previousBlockDifficulty.call()
      assert(newBlockDifficulty.eq(new BN(44)))
    })
  })

  describe('evaluateProofDifficulty()', async () => {
    it('reverts on unknown difficulty', async () => {
      await deployed.SystemStub.setCurrentDiff(1)
      await deployed.SystemStub.setPreviousDiff(1)
      try {
        await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'not at current or previous difficulty')
      }
    })

    it('evaluates a header proof with previous', async () => {
      await deployed.SystemStub.setPreviousDiff(5646403851534)
      await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
    })

    it('evaluates a header proof with current', async () => {
      await deployed.SystemStub.setPreviousDiff(1)
      await deployed.SystemStub.setCurrentDiff(5646403851534)
      await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
    })

    it('reverts on low difficulty', async () => {
      try {
        await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0].slice(0, 160 * 4 + 2))
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Insufficient accumulated difficulty in header chain')
      }
    })

    it('reverts on a ValidateSPV error code (3 or lower)', async () => {
      try {
        await deployed.SystemStub.setPreviousDiff(1)
        await testUtilsInstance.evaluateProofDifficulty(utils.LOW_DIFF_HEADER)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'ValidateSPV returned an error code')
      }
    })
  })

  describe('checkProofFromTx()', async () => {
    it('returns the correct _txid', async () => {
      await deployed.SystemStub.setCurrentDiff(6379265451411)
      const res = await testUtilsInstance.checkProofFromTx.call(utils.TX.tx, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
      assert.equal(res, utils.TX.tx_id_le)
    })

    it('fails with a broken proof', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProofFromTx.call(utils.TX.tx, utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and tx')
      }
    })

    it('fails with a broken tx', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProofFromTx.call('0x00', utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })
  })

  describe('checkProofFromTxId()', async () => {
    it('does not error', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProofFromTxId.call(utils.TX.tx_id_le, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
        assert(true, 'passes proof validation')
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })

    it('fails with a broken proof', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProofFromTxId.call(utils.TX.tx_id_le, utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and txId')
      }
    })

    it('fails with a broken txId', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProofFromTxId.call('0x00', utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and txId')
      }
    })
  })

  describe('findAndParseFundingOutput()', async () => {
    const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
    const _fundingOutputIndex = 0
    const _outValueBytes = '0x2040351d00000000'
    const _signerPubkeyX = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e'
    const _signerPubkeyY = '0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'
    it('correctly returns valuebytes', async () => {
      await testUtilsInstance.setPubKey(_signerPubkeyX, _signerPubkeyY)
      const valueBytes = await testUtilsInstance.findAndParseFundingOutput.call(_txOutputVector, _fundingOutputIndex)
      assert.equal(_outValueBytes, valueBytes, 'Got incorrect value bytes from funding output')
    })

    it('fails with incorrect singer pubKey', async () => {
      await testUtilsInstance.setPubKey('0x' + '11'.repeat(20), '0x' + '11'.repeat(20))
      try {
        await testUtilsInstance.findAndParseFundingOutput.call(_txOutputVector, _fundingOutputIndex)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Did not find output with correct PKH')
      }
    })
  })

  describe('extractOutputAtIndex()', async () => {
    it('extracts outputs at specified indices', async () => {
      let res
      const _txOutputVector1 = '0x024897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211'
      const _txOutputVector2 = '0x024db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f040420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922'
      const _txOutputVector3 = '0x044897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b078952114db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f040420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922'
      const _txOutputVector4 = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector1, 0)
      assert.equal(res, '0x4897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c18')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector1, 1)
      assert.equal(res, '0x0000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector2, 0)
      assert.equal(res, '0x4db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f0')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector2, 1)
      assert.equal(res, '0x40420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector3, 3)
      assert.equal(res, '0x40420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector4, 0)
      assert.equal(res, '0x2040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6')
      try {
        res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector1, 2)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Slice out of bounds')
      }
    })
  })

  describe('auctionValue()', async () => {
    it.skip('is TODO')
  })

  describe('signerFee()', async () => {
    it('returns a derived constant', async () => {
      const signerFee = await testUtilsInstance.signerFee.call()
      assert(signerFee.eq(new BN(500000)))
    })
  })

  describe('beneficiaryReward()', async () => {
    it('returns a derived constant', async () => {
      const beneficiaryReward = await testUtilsInstance.beneficiaryReward.call()
      assert(beneficiaryReward.eq(new BN(10 ** 5)))
    })
  })

  describe('determineCompressionPrefix()', async () => {
    it('selects 2 for even', async () => {
      const res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(32))
      assert.equal(res, '0x02')
    })
    it('selects 3 for odd', async () => {
      const res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(31) + '01')
      assert.equal(res, '0x03')
    })
  })

  describe('compressPubkey()', async () => {
    it('returns a 33 byte array with a prefix', async () => {
      const compressed = await testUtilsInstance.compressPubkey.call('0x' + '00'.repeat(32), '0x' + '00'.repeat(32))
      assert.equal(compressed, '0x02' + '00'.repeat(32))
    })
  })

  describe('signerPubkey()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      await testUtilsInstance.setPubKey('0x' + '00'.repeat(32), '0x' + '00'.repeat(32))
      const signerPubkey = await testUtilsInstance.signerPubkey.call()
      assert.equal(signerPubkey, '0x' + '00'.repeat(64))
    })
  })

  describe('signerPKH()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      await testUtilsInstance.setPubKey('0x' + '00'.repeat(32), '0x' + '00'.repeat(32))
      const signerPKH = await testUtilsInstance.signerPKH.call()
      assert.equal(signerPKH, utils.hash160('02' + '00'.repeat(32)))
    })
  })

  describe('utxoSize()', async () => {
    it('returns the state\'s utxoSizeBytes as an integer', async () => {
      const utxoSize = await testUtilsInstance.utxoSize.call()
      assert(utxoSize.eq(new BN(0)))

      await testUtilsInstance.setUTXOInfo('0x11223344', 1, '0x')
      const newUtxoSize = await testUtilsInstance.utxoSize.call()
      assert(newUtxoSize.eq(new BN('44332211', 16)))
    })
  })

  describe('fetchOraclePrice()', async () => {
    it('calls out to the system', async () => {
      const oraclePrice = await testUtilsInstance.fetchOraclePrice.call()
      assert(oraclePrice.eq(new BN('1000000000000', 10)))

      await deployed.SystemStub.setOraclePrice(44)
      const newOraclePrice = await testUtilsInstance.fetchOraclePrice.call()
      assert(newOraclePrice.eq(new BN(44)))
    })
  })

  describe('fetchBondAmount()', async () => {
    it('calls out to the keep system', async () => {
      const bondAmount = await testUtilsInstance.fetchBondAmount.call()
      assert(bondAmount.eq(new BN(10000)))

      await deployed.KeepStub.setBondAmount(44)
      const newBondAmount = await testUtilsInstance.fetchBondAmount.call()
      assert(newBondAmount.eq(new BN(44)))
    })
  })

  describe('bytes8LEToUint()', async () => {
    it('interprets bytes as LE uints ', async () => {
      let res = await testUtilsInstance.bytes8LEToUint.call('0x' + '00'.repeat(8))
      assert(res.eq(new BN(0)))
      res = await testUtilsInstance.bytes8LEToUint.call('0x11223344')
      assert(res.eq(new BN('44332211', 16)))
      res = await testUtilsInstance.bytes8LEToUint.call('0x1100229933884477')
      assert(res.eq(new BN('7744883399220011', 16)))
    })
  })

  describe('wasDigestApprovedForSigning()', async () => {
    it('calls out to the keep system', async () => {
      const approved = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert.equal(approved, false)

      await deployed.KeepStub.approveDigest(7, '0x' + '00'.repeat(32))
      const newApproved = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert(newApproved.eq(new BN(100)))
    })
  })

  describe('depositBeneficiary()', async () => {
    it('calls out to the system', async () => {
      const res = await testUtilsInstance.depositBeneficiary.call()
      assert.equal(res, '0x' + '00'.repeat(19) + '00')
    })
  })

  describe('redemptionTeardown()', async () => {
    it('deletes state', async () => {
      await testUtilsInstance.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 5, 6, '0x' + '33'.repeat(32))
      const requestInfo = await testUtilsInstance.getRequestInfo.call()
      assert.equal(requestInfo[4], '0x' + '33'.repeat(32))
      await testUtilsInstance.redemptionTeardown()
      const newRequestInfo = await testUtilsInstance.getRequestInfo.call()
      assert.equal(newRequestInfo[4], '0x' + '00'.repeat(32))
    })
  })

  describe('seizeSignerBonds()', async () => {
    it('calls out to the keep system and returns the seized amount', async () => {
      const value = 5000
      await deployed.KeepStub.send(value, { from: accounts[0] })
      const seized = await testUtilsInstance.seizeSignerBonds.call()
      await testUtilsInstance.seizeSignerBonds()
      assert(seized.eq(new BN(value)))
    })

    it('errors if no funds were seized', async () => {
      try {
        await testUtilsInstance.seizeSignerBonds.call()
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'No funds received, unexpected')
      }
    })
  })

  describe('distributeBeneficiaryReward()', async () => {
    it('checks that beneficiary is rewarded', async () => {
      const beneficiary = accounts[5]
      // reward should == 10**18. This is a stub value. parameter address is irrelevant
      const reward = await deployed.TBTCStub.balanceOf.call(accounts[0])
      const initialTokenBalance = await deployed.TBTCStub.getBalance(beneficiary)
      await deployed.SystemStub.setDepositOwner(0, beneficiary)

      await testUtilsInstance.distributeBeneficiaryReward()

      const finalTokenBalance = await deployed.TBTCStub.getBalance(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add( new BN(reward))
      expect(finalTokenBalance, 'tokens not rewarded to beneficiary correctly').to.eq.BN(tokenCheck)
    })
  })

  describe('pushFundsToKeepGroup()', async () => {
    it('calls out to the keep contract', async () => {
      const value = 10000
      await testUtilsInstance.send(value, { from: accounts[0] })
      await testUtilsInstance.pushFundsToKeepGroup(value)
      const keepBalance = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(keepBalance, value) // web3 balances are integers I guess
    })

    it('reverts if insufficient value', async () => {
      try {
        await testUtilsInstance.pushFundsToKeepGroup.call(10000000000000)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, ('Not enough funds to send'))
      }
    })
  })
})
