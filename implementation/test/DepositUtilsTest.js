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
const TBTCSystemStub = artifacts.require('TBTCSystemStub')

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
  { name: 'TBTCSystemStub', contract: TBTCSystemStub }]

// real tx from mainnet bitcoin, interpreted as funding tx
// tx source: https://www.blockchain.com/btc/tx/7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f
// const tx = '0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800'
// const txid = '0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f';
// const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c';
const currentDifficulty = 6353030562983
const _version = '0x01000000'
const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
const _fundingOutputIndex = 0
const _txLocktime = '0x4ec10800'
const _txIndexInBlock = 130
const _bitcoinHeaders = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
const _signerPubkeyX = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e'
const _signerPubkeyY = '0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'
const _merkleProof = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe04995ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b5095548'
const _expectedUTXOoutpoint = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c00000000'
// const _outputValue = 490029088;
const _outValueBytes = '0x2040351d00000000'

contract('DepositUtils', (accounts) => {
  let deployed
  let testUtilsInstance
  before(async () => {
    deployed = await utils.deploySystem(TEST_DEPOSIT_UTILS_DEPLOY)
    testUtilsInstance = deployed.TestDepositUtils

    await testUtilsInstance.createNewDeposit(
      deployed.TBTCSystemStub.address,
      deployed.TBTCStub.address,
      deployed.KeepStub.address,
      1, // m
      1) // n
  })

  describe('currentBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      const blockDifficulty = await testUtilsInstance.currentBlockDifficulty.call()
      assert(blockDifficulty.eq(new BN(1)))

      await deployed.TBTCSystemStub.setCurrentDiff(33)
      const newBlockDifficulty = await testUtilsInstance.currentBlockDifficulty.call()
      assert(newBlockDifficulty.eq(new BN(33)))
    })
  })

  describe('previousBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      const blockDifficulty = await testUtilsInstance.previousBlockDifficulty.call()
      assert(blockDifficulty.eq(new BN(1)))

      await deployed.TBTCSystemStub.setPreviousDiff(44)
      const newBlockDifficulty = await testUtilsInstance.previousBlockDifficulty.call()
      assert(newBlockDifficulty.eq(new BN(44)))
    })
  })

  describe('evaluateProofDifficulty()', async () => {
    it('reverts on unknown difficulty', async () => {
      await deployed.TBTCSystemStub.setCurrentDiff(1)
      await deployed.TBTCSystemStub.setPreviousDiff(1)
      try {
        await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'not at current or previous difficulty')
      }
    })

    it('evaluates a header proof with previous', async () => {
      await deployed.TBTCSystemStub.setPreviousDiff(5646403851534)
      await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
    })

    it('evaluates a header proof with current', async () => {
      await deployed.TBTCSystemStub.setPreviousDiff(1)
      await deployed.TBTCSystemStub.setCurrentDiff(5646403851534)
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
        await deployed.TBTCSystemStub.setPreviousDiff(1)
        await testUtilsInstance.evaluateProofDifficulty(utils.LOW_DIFF_HEADER)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'ValidateSPV returned an error code')
      }
    })
  })

  describe('checkProofFromTx()', async () => {
    it('returns the correct _txid', async () => {
      await deployed.TBTCSystemStub.setCurrentDiff(6379265451411)
      const res = await testUtilsInstance.checkProof.call(utils.TX.tx, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
      assert.equal(res, utils.TX.tx_id_le)
    })

    it('fails with a broken proof', async () => {
      try {
        await deployed.TBTCSystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProof.call(utils.TX.tx, utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and tx')
      }
    })

    it('fails with a broken tx', async () => {
      try {
        await deployed.TBTCSystemStub.setCurrentDiff(6379265451411)
        await testUtilsInstance.checkProof.call('0x00', utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })
  })

  describe('checkProofFromTxId()', async () => {
    before(async () => {
      await deployed.SystemStub.setCurrentDiff(utils.TX.difficulty)
    })

    it('does not error', async () => {
      try {
        await testUtilsInstance.checkProofFromTxId.call(utils.TX.tx_id_le, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
        assert(true, 'passes proof validation')
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })

    it('fails with a broken proof', async () => {
      try {
        await testUtilsInstance.checkProofFromTxId.call(utils.TX.tx_id_le, utils.TX.proof, 0, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and txId')
      }
    })

    it('fails with bad difficulty', async () => {
      await deployed.SystemStub.setCurrentDiff(1)
      try {
        await testUtilsInstance.checkProofFromTxId.call(utils.TX.tx_id_le, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'not at current or previous difficulty')
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

    it('fails with incorrect signer pubKey', async () => {
      await testUtilsInstance.setPubKey('0x' + '11'.repeat(20), '0x' + '11'.repeat(20))
      try {
        await testUtilsInstance.findAndParseFundingOutput.call(_txOutputVector, _fundingOutputIndex)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'could not identify output funding the required public key hash')
      }
    })
  })

  describe('extractOutputAtIndex()', async () => {
    it('extracts outputs at specified indicex (vector length 1)', async () => {
      const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
      const res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector, 0)
      assert.equal(res, '0x2040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6')
    })

    it('extracts outputs at specified indices (vector length 2)', async () => {
      let res
      const _txOutputVector1 = '0x024897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211'
      const _txOutputVector2 = '0x024db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f040420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922'
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector1, 0)
      assert.equal(res, '0x4897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c18')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector1, 1)
      assert.equal(res, '0x0000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector2, 0)
      assert.equal(res, '0x4db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f0')
      res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector2, 1)
      assert.equal(res, '0x40420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922')
    })

    it('extracts outputs at specified index (vecor length 4)', async () => {
      const _txOutputVector = '0x044897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b078952114db6000000000000160014455c0ea778752831d6fc25f6f8cf55dc49d335f040420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922'
      const res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector, 3)
      assert.equal(res, '0x40420f0000000000220020aedad4518f56379ef6f1f52f2e0fed64608006b3ccaff2253d847ddc90c91922')
    })

    it('fails to extract output from bad index', async () => {
      const _txOutputVector= '0x024897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211'
      try {
        res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector, 2)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Slice out of bounds')
      }
    })

    it('fails to extract output from varint prepended vector', async () => {
      // we don't need to include the number of outputs suggested by the varint
      const _txOutputVector= '0xfe123412344897070000000000220020a4333e5612ab1a1043b25755c89b16d55184a42f81799e623e6bc39db8539c180000000000000000166a14edb1b5c2f39af0fec151732585b1049b07895211'
      try {
        res = await testUtilsInstance.extractOutputAtIndex.call(_txOutputVector, 2)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'VarInts not supported, Number of outputs cannot exceed 252')
      }
    })
  })

  describe('validateAndParseFundingSPVProof()', async () => {
    before(async () => {
      await testUtilsInstance.setPubKey(_signerPubkeyX, _signerPubkeyY)
      await deployed.SystemStub.setCurrentDiff(currentDifficulty)
    })

    it('returns currect value and outpoint', async () => {
      const parseResults = await testUtilsInstance.validateAndParseFundingSPVProof.call(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
      assert.equal(parseResults[0], _outValueBytes)
      assert.equal(parseResults[1], _expectedUTXOoutpoint)
    })

    it('fails with bad _txInputVector', async () => {
      try {
        await testUtilsInstance.validateAndParseFundingSPVProof.call(_version, '0x' + '00'.repeat(32), _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header and txId')
      }
    })

    it('fails with insufficient difficulty', async () => {
      const _badheaders = `0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6`
      try {
        await testUtilsInstance.validateAndParseFundingSPVProof.call(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _badheaders)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Insufficient accumulated difficulty in header chain')
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
      const _signerPubkeyX = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e'
      const _signerPubkeyY = '0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'
      const _concatinatedKeys = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'

      await testUtilsInstance.setPubKey(_signerPubkeyX, _signerPubkeyY)
      const signerPubkey = await testUtilsInstance.signerPubkey.call()
      assert.equal(signerPubkey, _concatinatedKeys)
    })

    it('returns base value for unset public key', async () => {
      const newTestUtilsInstance = await TestDepositUtils.new()
      const signerPubkey = await newTestUtilsInstance.signerPubkey.call()
      console.log(signerPubkey)
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

      await deployed.TBTCSystemStub.setOraclePrice(44)
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
      await deployed.TBTCSystemStub.setDepositOwner(0, beneficiary)

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
