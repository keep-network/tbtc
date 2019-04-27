const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositLog = artifacts.require('DepositLog')
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

STANDARD_DEPLOY = [
  {name: 'BytesLib', contract: BytesLib},
  {name: 'BTCUtils', contract: BTCUtils},
  {name: 'ValidateSPV', contract: ValidateSPV},
  {name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs},
  {name: 'TBTCConstants', contract: TestTBTCConstants},  // note the name
  {name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging},
  {name: 'DepositLog', contract: DepositLog},
  {name: 'DepositStates', contract: DepositStates},
  {name: 'DepositUtils', contract: DepositUtils},
  {name: 'DepositFunding', contract: DepositFunding},
  {name: 'DepositRedemption', contract: DepositRedemption},
  {name: 'DepositLiquidation', contract: DepositLiquidation},
  {name: 'TestDepositUtils', contract: TestDepositUtils},
  {name: 'KeepStub', contract: KeepStub},
  {name: 'TBTCStub', contract: TBTCStub},
  {name: 'SystemStub', contract: SystemStub}]


contract.only('DepositUtils', accounts => {

  let deployed
  let testUtilsInstance

  before(async () => {
      deployed = await utils.deploySystem(STANDARD_DEPLOY)
      testUtilsInstance = deployed.TestDepositUtils

      await testUtilsInstance.createNewDeposit(
        deployed.SystemStub.address,
        deployed.TBTCStub.address,
        deployed.KeepStub.address,
        1,  //m
        1)  //n
  })

  describe('currentBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      res = await testUtilsInstance.currentBlockDifficulty.call()
      assert(res.eq(new BN(1)))

      await deployed.SystemStub.setCurrentDiff(33)
      res = await testUtilsInstance.currentBlockDifficulty.call()
      assert(res.eq(new BN(33)))
    })
  })

  describe('previousBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      res = await testUtilsInstance.previousBlockDifficulty.call()
      assert(res.eq(new BN(1)))

      await deployed.SystemStub.setPreviousDiff(44)
      res = await testUtilsInstance.previousBlockDifficulty.call()
      assert(res.eq(new BN(44)))
    })
  })

  describe('evaluateProofDifficulty()', async () => {
    it('reverts on unknown difficulty', async () => {
      await deployed.SystemStub.setCurrentDiff(1)
      await deployed.SystemStub.setPreviousDiff(1)
      try {
        res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'not at current or previous difficulty')
      }
    })

    it('evaluates a header proof with previous', async () => {
      await deployed.SystemStub.setPreviousDiff(5646403851534)
      res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])

    })

    it('evaluates a header proof with current', async () => {
      await deployed.SystemStub.setPreviousDiff(1)
      await deployed.SystemStub.setCurrentDiff(5646403851534)
      res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
    })

    it('reverts on low difficulty', async () => {
      try {
        res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0].slice(0, 160 * 4 + 2))
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Insufficient accumulated difficulty in header chain')
      }
    })

    it('reverts on a ValidateSPV error code (3 or lower)', async () => {
      try {
        await deployed.SystemStub.setPreviousDiff(1)
        res = await testUtilsInstance.evaluateProofDifficulty(utils.LOW_DIFF_HEADER)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'ValidateSPV returned an error code')
      }
    })
  })

  describe('checkProof()', async () => {
    it('returns the correct _txid' , async () => {
      await deployed.SystemStub.setCurrentDiff(6379265451411)
      res = await testUtilsInstance.checkProof.call(utils.TX.tx, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
      assert.equal(res, utils.TX.tx_id_le)
    })

    it('fails with a broken proof', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProof.call(utils.TX.tx, utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and tx')
      }
    })

    it('fails with a broken tx', async () => {
      try {
        await deployed.SystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProof.call('0x00', utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })
  })

  describe('auctionValue()', async () => {
    it.skip('is TODO')
  })

  describe('signerFee()', async () => {
    it('returns a derived constant', async () => {
      res = await testUtilsInstance.signerFee.call()
      assert(res.eq(new BN(500000)))
    })
  })

  describe('beneficiaryReward()', async () => {
    it('returns a derived constant', async () => {
      res = await testUtilsInstance.beneficiaryReward.call()
      assert(res.eq(new BN(10 ** 5)))
    })
  })

  describe('determineCompressionPrefix()', async () => {
    it('selects 2 for even', async () => {
      res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(32))
      assert.equal(res, '0x02')
    })
    it('selects 3 for odd', async () => {
      res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(31) + '01')
      assert.equal(res, '0x03')
    })
  })

  describe('compressPubkey()', async () => {
    it('returns a 33 byte array with a prefix', async () => {
      res = await testUtilsInstance.compressPubkey.call('0x' + '00'.repeat(32), '0x' + '00'.repeat(32))
      assert.equal(res, '0x02' + '00'.repeat(32))
    })
  })

  describe('signerPubkey()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      res = await testUtilsInstance.signerPubkey.call()
      assert.equal(res, '0x' + '00'.repeat(64))
    })
  })

  describe('signerPKH()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      res = await testUtilsInstance.signerPKH.call()
      assert.equal(res, utils.hash160('02' + '00'.repeat(32)))
    })
  })

  describe('utxoSize()', async () => {
    it('returns the state\'s utxoSizeBytes as an integer', async () => {
      res = await testUtilsInstance.utxoSize.call()
      assert(res.eq(new BN(0)))

      await testUtilsInstance.setUTXOInfo('0x11223344', 01, '0x')
      res = await testUtilsInstance.utxoSize.call()
      assert(res.eq(new BN('44332211', 16)))
    })
  })

  describe('fetchOraclePrice()', async () => {
    it('calls out to the system', async () => {
      res = await testUtilsInstance.fetchOraclePrice.call()
      assert(res.eq(new BN('1000000000000', 10)))

      await deployed.SystemStub.setOraclePrice(44)
      res = await testUtilsInstance.fetchOraclePrice.call()
      assert(res.eq(new BN(44)))
    })
  })

  describe('fetchBondAmount()', async () => {
    it('calls out to the keep system', async () => {
      res = await testUtilsInstance.fetchBondAmount.call()
      assert(res.eq(new BN(10000)))

      await deployed.KeepStub.setBondAmount(44)
      res = await testUtilsInstance.fetchBondAmount.call()
      assert(res.eq(new BN(44)))
    })
  })

  describe('bytes8LEToUint()', async () => {
    it('interprets bytes as LE uints ', async () => {
      res = await testUtilsInstance.bytes8LEToUint.call('0x' + '00'.repeat(8))
      assert(res.eq(new BN(0)))
      res = await testUtilsInstance.bytes8LEToUint.call('0x11223344')
      assert(res.eq(new BN('44332211', 16)))
      res = await testUtilsInstance.bytes8LEToUint.call('0x1100229933884477')
      assert(res.eq(new BN('7744883399220011', 16)))
    })
  })

  describe('wasDigestApprovedForSigning()', async () => {
    it('calls out to the keep system', async () => {
      res = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert.equal(res, false)

      await deployed.KeepStub.approveDigest(7, '0x' + '00'.repeat(32))
      res = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert(res.eq(new BN(100)))
    })
  })

  describe('depositBeneficiary()', async () => {
    it('calls out to the system', async () => {
      res = await testUtilsInstance.depositBeneficiary.call()
      assert.equal(res, '0x' + '00'.repeat(19) + '07')
    })
  })

  describe('redemptionTeardown()', async () => {
    it('deletes state', async () => {
      await testUtilsInstance.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 5, 6, '0x' + '33'.repeat(32))
      res = await testUtilsInstance.getRequestInfo.call()
      assert.equal(res[4], '0x' + '33'.repeat(32))
      await testUtilsInstance.redemptionTeardown()
      res = await testUtilsInstance.getRequestInfo.call()
      assert.equal(res[4], '0x' + '00'.repeat(32))
    })
  })

  describe('seizeSignerBonds()', async () => {
    it('calls out to the keep system and returns the seized amount', async () => {
      value = 5000
      await deployed.KeepStub.send(value, {from: accounts[0]})
      res = await testUtilsInstance.seizeSignerBonds.call()
      await testUtilsInstance.seizeSignerBonds()
      assert(res.eq(new BN(value)))
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
    it.skip('is TODO')
  })

  describe('pushFundsToKeepGroup()', async () => {
    it('calls out to the keep contract', async () => {
      value = 10000
      await testUtilsInstance.send(value, {from: accounts[0]})
      await testUtilsInstance.pushFundsToKeepGroup(value)
      keepBalance = await web3.eth.getBalance(deployed.KeepStub.address)
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
