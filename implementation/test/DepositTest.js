import expectThrow from './helpers/expectThrow'

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
const TestDeposit = artifacts.require('TestDeposit')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
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
  { name: 'TestDeposit', contract: TestDeposit },
  { name: 'KeepStub', contract: KeepStub },
  { name: 'TBTCStub', contract: TBTCStub },
  { name: 'SystemStub', contract: SystemStub }]

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'


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

contract('Deposit', (accounts) => {
  let deployed
  let testInstance
  let withdrawalRequestTime
  let fundingProofTimerStart

  before(async () => {
    deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)
    testInstance = deployed.TestDeposit
    testInstance.setExteroriorAddresses(deployed.SystemStub.address, deployed.TBTCStub.address, deployed.KeepStub.address)
  })

  beforeEach(async () => {
    await testInstance.reset()
  })

  describe('getCurrentState', async () => {
    it.skip('seems too small to test. maybe later')
  })

  describe('createNewDeposit', async () => {
    it('runs and updates state and fires a created event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.createNewDeposit(
        deployed.SystemStub.address,
        deployed.TBTCStub.address,
        deployed.KeepStub.address,
        1, // m
        1)

      // state updates
      const depositState = await testInstance.getState.call()
      expect(depositState, 'state not as expected').to.eq.BN(utils.states.AWAITING_SIGNER_SETUP)

      const keepInfo = await testInstance.getKeepInfo.call()
      expect(keepInfo[0], 'keepID not as expected').to.eq.BN(7)
      expect(keepInfo[1], 'signing group timestamp not as expected').not.to.eq.BN(0)

      // fired an event
      const eventList = await deployed.SystemStub.getPastEvents('Created', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._keepID, '7')
    })

    it('reverts if not in the start state', async () => {
      await testInstance.setState(utils.states.REDEEMED)
      try {
        await testInstance.createNewDeposit.call(
          deployed.SystemStub.address,
          deployed.TBTCStub.address,
          deployed.KeepStub.address,
          1, // m
          1)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Deposit setup already requested')
      }
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
    const requesterPKH = '0x' + '33'.repeat(20)

    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await testInstance.setUTXOInfo(valueBytes, 0, outpoint)
    })

    it('updates state successfully and fires a RedemptionRequested event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.setKeepInfo(0, 0, 0, keepPubkeyX, keepPubkeyY)

      // the fee is ~12,297,829,380 BTC
      await testInstance.requestRedemption('0x1111111100000000', requesterPKH)

      const requestInfo = await testInstance.getRequestInfo()
      assert.equal(requestInfo[1], requesterPKH)
      assert(!requestInfo[3].eqn(0)) // withdrawalRequestTime is set
      assert.equal(requestInfo[4], sighash)

      // fired an event
      const eventList = await deployed.SystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, sighash)
    })

    it('reverts if not in Active or Courtesy', async () => {
      await testInstance.setState(utils.states.LIQUIDATED)
      try {
        await testInstance.requestRedemption('0x1111111100000000', '0x' + '33'.repeat(20))
      } catch (e) {
        assert.include(e.message, 'Redemption only available from Active or Courtesy state')
      }
    })

    it('reverts if the fee is low', async () => {
      try {
        await testInstance.requestRedemption('0x0011111111111111', '0x' + '33'.repeat(20))
      } catch (e) {
        assert.include(e.message, 'Fee is too low')
      }
    })

    it('reverts if the keep returns false', async () => {
      try {
        await deployed.KeepStub.setSuccess(false)
        await testInstance.requestRedemption('0x1111111100000000', '0x' + '33'.repeat(20))
      } catch (e) {
        await deployed.KeepStub.setSuccess(true)
        assert.include(e.message, 'Keep returned false')
      }
    })

    it('calls Keep to approve the digest', async () => {
      // test relies on a side effect
      const approved = await deployed.KeepStub.wasDigestApprovedForSigning.call(0, sighash)
      assert(!approved.eqn(0), 'digest was not approved')
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
      await testInstance.setState(utils.states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it('updates the state and logs GotRedemptionSignature', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.setKeepInfo(0, 0, 0, pubkeyX, pubkeyY)
      await testInstance.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 0, 0, digest)

      await testInstance.provideRedemptionSignature(v, r, s)

      const state = await testInstance.getState.call()
      expect(state).to.eq.BN(utils.states.AWAITING_WITHDRAWAL_PROOF)

      // fired an event
      const eventList = await deployed.SystemStub.getPastEvents('GotRedemptionSignature', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._r, r)
      assert.equal(eventList[0].returnValues._s, s)
    })

    it('errors if not awaiting withdrawal signature', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.provideRedemptionSignature(v, r, s)
      } catch (e) {
        assert.include(e.message, 'Not currently awaiting a signature')
      }
    })

    it('errors on invaid sig', async () => {
      try {
        await testInstance.provideRedemptionSignature(28, r, s)
      } catch (e) {
        assert.include(e.message, 'Invalid signature')
      }
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
    const requesterPKH = '0x' + '33'.repeat(20)
    let feeIncreaseTimer

    before(async () => {
      feeIncreaseTimer = await deployed.TBTCConstants.getIncreaseFeeTimer.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - feeIncreaseTimer.toNumber()
      await deployed.KeepStub.setDigestApprovedAtTime(prevSighash, withdrawalRequestTime)
      await testInstance.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testInstance.setKeepInfo(0, 0, 0, keepPubkeyX, keepPubkeyY)
      await testInstance.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testInstance.setRequestInfo(utils.address0, requesterPKH, initialFee, withdrawalRequestTime, prevSighash)
    })

    it('approves a new digest for signing, updates the state, and logs RedemptionRequested', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.increaseRedemptionFee(previousOutputBytes, newOutputBytes)

      const requestInfo = await testInstance.getRequestInfo.call()
      assert.equal(requestInfo[4], nextSighash)

      // fired an event
      const eventList = await deployed.SystemStub.getPastEvents('RedemptionRequested', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._digest, nextSighash)
    })

    it('reverts if not awaiting withdrawal proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.increaseRedemptionFee(previousOutputBytes, newOutputBytes)
      } catch (e) {
        assert.include(e.message, 'Fee increase only available after signature provided')
      }
    })

    it('reverts if the increase fee timer has not elapsed', async () => {
      try {
        const block = await web3.eth.getBlock('latest')
        await testInstance.setRequestInfo(utils.address0, requesterPKH, initialFee, block.timestamp, prevSighash)
        await testInstance.increaseRedemptionFee(previousOutputBytes, newOutputBytes)
      } catch (e) {
        assert.include(e.message, 'Fee increase not yet permitted')
      }
    })

    it('reverts if the fee step is not linear', async () => {
      try {
        await testInstance.increaseRedemptionFee(previousOutputBytes, '0x1101010101102201')
      } catch (e) {
        assert.include(e.message, 'Not an allowed fee step')
      }
    })

    it('reverts if the previous sighash was not the latest approved', async () => {
      try {
        await testInstance.setRequestInfo(utils.address0, requesterPKH, initialFee, withdrawalRequestTime, keepPubkeyX)
        await testInstance.increaseRedemptionFee(previousOutputBytes, newOutputBytes)
      } catch (e) {
        assert.include(e.message, 'Provided previous value does not yield previous sighash')
      }
    })

    it('reverts if the keep returned false', async () => {
      try {
        await deployed.KeepStub.setSuccess(false)
        await testInstance.increaseRedemptionFee(previousOutputBytes, newOutputBytes)
      } catch (e) {
        await deployed.KeepStub.setSuccess(true)
        assert.include(e.message, 'Keep returned false')
      }
    })
  })

  describe('provideRedemptionProof', async () => {
    // real tx from mainnet bitcoin
    const currentDiff = 6353030562983
    // const txid = '0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f'
    const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c'
    const tx = '0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800'
    const proof = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe04995ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b5095548'
    const index = 130
    const headerChain = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
    const outpoint = '0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000'
    const prevoutValueBytes = '0xf078351d00000000'
    const requesterPKH = '0x86e7303082a6a21d5837176bc808bf4828371ab6'

    beforeEach(async () => {
      await deployed.SystemStub.setCurrentDiff(currentDiff)
      await testInstance.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testInstance.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testInstance.setRequestInfo('0x' + '11'.repeat(20), requesterPKH, 14544, 0, '0x' + '11' * 32)
    })

    it('updates the state, deconstes struct info, calls TBTC and Keep, and emits a Redeemed event', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.provideRedemptionProof(tx, proof, index, headerChain)

      const depositState = await testInstance.getState.call()
      assert(depositState.eq(utils.states.REDEEMED))

      const requestInfo = await testInstance.getRequestInfo.call()
      assert.equal(requestInfo[0], '0x' + '11'.repeat(20)) // address is intentionally not cleared
      assert.equal(requestInfo[1], utils.address0)
      assert.equal(requestInfo[4], utils.bytes32zero)

      const eventList = await deployed.SystemStub.getPastEvents('Redeemed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._txid, txidLE)
    })

    it('reverts if not in the redemption flow', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.provideRedemptionProof(tx, proof, index, headerChain)
      } catch (e) {
        assert.include(e.message, 'Redemption proof only allowed from redemption flow')
      }
    })

    it('reverts if the merkle proof is not validated successfully', async () => {
      try {
        await testInstance.provideRedemptionProof(tx, proof, 0, headerChain)
      } catch (e) {
        assert.include(e.message, 'Tx merkle proof is not valid for provided header')
      }
    })
  })

  describe('redemptionTransactionChecks', async () => {
    const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c'
    const tx = '0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800'
    const badtx = '0x05' + tx.slice(4)
    const outputValue = 490029088
    const outpoint = '0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000'
    const prevoutValueBytes = '0xf078351d00000000'
    const requesterPKH = '0x86e7303082a6a21d5837176bc808bf4828371ab6'

    beforeEach(async () => {
      await testInstance.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await testInstance.setRequestInfo('0x' + '11'.repeat(20), requesterPKH, 14544, 0, '0x' + '11' * 32)
    })

    it('returns the little-endian txid and output value', async () => {
      const redemptionChecks = await testInstance.redemptionTransactionChecks.call(tx)
      assert.equal(redemptionChecks[0], txidLE)
      expect(redemptionChecks[1]).to.eq.BN(outputValue)
    })

    it('reverts if tx parsing fails', async () => {
      try {
        await testInstance.redemptionTransactionChecks(badtx)
      } catch (e) {
        assert.include(e.message, 'Failed tx parsing')
      }
    })

    it('reverts if the tx spends the wrong utxo', async () => {
      try {
        await testInstance.setUTXOInfo(prevoutValueBytes, 0, '0x' + '33'.repeat(36))
        await testInstance.redemptionTransactionChecks.call(tx)
      } catch (e) {
        assert.include(e.message, 'Tx spends the wrong UTXO')
      }
    })
    it('reverts if the tx sends value to the wrong pkh', async () => {
      try {
        await testInstance.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 14544, 0, '0x' + '11' * 32)
        await testInstance.redemptionTransactionChecks.call(tx)
      } catch (e) {
        assert.include(e.message, 'Tx sends value to wrong pubkeyhash')
      }
    })
  })

  describe('notifySignatureTimeout', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getSignatureTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await deployed.KeepStub.burnContractBalance()
      await testInstance.setState(utils.states.AWAITING_WITHDRAWAL_SIGNATURE)
      await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
    })

    it('reverts if not awaiting redemption signature', async () => {
      try {
        await testInstance.setState(utils.states.START)
      } catch (e) {
        assert.include(e.message, 'Not currently awaiting a signature')
      }
    })

    it('reverts if the signature timeout has not elapsed', async () => {
      await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime + timer + 1, utils.bytes32zero)
      try {
        await testInstance.notifySignatureTimeout()
      } catch (e) {
        assert.include(e.message, 'Signature timer has not elapsed')
      }
    })

    it('reverts if no funds recieved as signer bond', async () => {
      await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime * 5, utils.bytes32zero)
      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(0, bond, 'no bond should be sent')
      try {
        await expectThrow(testInstance.notifySignatureTimeout())
      } catch (e) {
        assert.include(e.message, 'No funds should be received as signer bond')
      }
    })

    it('starts abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifySignatureTimeout()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyRedemptionProofTimeout', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getRedepmtionProofTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      withdrawalRequestTime = blockTimestamp - timer.toNumber() - 1
      await deployed.KeepStub.burnContractBalance()
      await testInstance.setState(utils.states.AWAITING_WITHDRAWAL_PROOF)
      await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
    })

    it('reverts if not awaiting redemption proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
      } catch (e) {
        assert.include(e.message, 'Not currently awaiting a redemption proof')
      }
    })

    it('reverts if the proof timeout has not elapsed', async () => {
      try {
        await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime * 5, utils.bytes32zero)
        await testInstance.notifyRedemptionProofTimeout()
      } catch (e) {
        assert.include(e.message, 'Proof timer has not elapsed')
      }
    })
    it('reverts if no funds recieved as signer bond', async () => {
      await testInstance.setRequestInfo(utils.address0, utils.address0, 0, withdrawalRequestTime, utils.bytes32zero)
      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(0, bond, 'no bond should be sent')
      await expectThrow(testInstance.notifyRedemptionProofTimeout())
    })

    it('starts abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifyRedemptionProofTimeout()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifySignerSetupFailure', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getSigningGroupFormationTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1
      await testInstance.setState(utils.states.AWAITING_SIGNER_SETUP)
      await testInstance.setKeepInfo(0, fundingProofTimerStart, 0, utils.bytes32zero, utils.bytes32zero)
    })

    it('updates state to setup failed, deconstes state, and logs SetupFailed', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.notifySignerSetupFailure()

      const keepInfo = await testInstance.getKeepInfo.call()
      expect(keepInfo[1], 'signingGroupRequestedAt should be 0').to.eq.BN(0)
      expect(keepInfo[2], 'fundingProofTimerStart should be 0').to.eq.BN(0)

      const eventList = await deployed.SystemStub.getPastEvents('SetupFailed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1, 'Event list is the wrong length')
    })

    it('reverts if not awaiting signer setup', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifySignerSetupFailure()
      } catch (e) {
        assert.include(e.message, 'Not awaiting setup')
      }
    })

    it('reverts if the timer has not yet elapsed', async () => {
      try {
        await testInstance.setKeepInfo(0, fundingProofTimerStart * 5, 0, utils.bytes32zero, utils.bytes32zero)
        await testInstance.notifySignerSetupFailure()
      } catch (e) {
        assert.include(e.message, 'Signing group formation timeout not yet elapsed')
      }
    })

    it('returns funder bond', async () => {
      const beneficiary = accounts[4]
      const bond = 1000000000000

      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.send(bond, { from: beneficiary })

      const initialBalance = await web3.eth.getBalance(beneficiary)

      await testInstance.setKeepInfo(0, fundingProofTimerStart, 0, utils.bytes32zero, utils.bytes32zero)
      await deployed.KeepStub.setBondAmount(bond)
      await testInstance.notifySignerSetupFailure()

      const balance = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(new BN(bond))
      expect(balance, 'bond not returned to funder').to.eq.BN(balanceCheck)
    })
  })

  describe('retrieveSignerPubkey', async () => {
    const pubkeyX = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa'
    const pubkeyY = '0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'

    beforeEach(async () => {
      await testInstance.setState(utils.states.AWAITING_SIGNER_SETUP)
    })

    it('updates the pubkey X and Y, changes state, and logs RegisteredPubkey', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.retrieveSignerPubkey()

      const keepInfo = await testInstance.getKeepInfo.call()
      assert.equal(keepInfo[3], pubkeyX)
      assert.equal(keepInfo[4], pubkeyY)

      const eventList = await deployed.SystemStub.getPastEvents('RegisteredPubkey', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._signingGroupPubkeyX, pubkeyX, 'Logged X is wrong')
      assert.equal(eventList[0].returnValues._signingGroupPubkeyY, pubkeyY, 'Logged Y is wrong')
    })

    it('reverts if not awaiting signer setup', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.retrieveSignerPubkey()
      } catch (e) {
        assert.include(e.message, 'Not currently awaiting signer setup')
      }
    })

    it('reverts if either half of the pubkey is 0', async () => {
      try {
        await deployed.KeepStub.setPubkey('0x' + '00'.repeat(64))
        await testInstance.retrieveSignerPubkey()
      } catch (e) {
        await deployed.KeepStub.setPubkey('0x' + '00')
        assert.include(e.message, 'Keep returned bad pubkey')
      }
    })
  })

  describe('notifyFundingTimeout', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1
      await testInstance.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart, utils.bytes32zero, utils.bytes32zero)
    })

    it('updates the state to failed setup, deconstes funding info, and logs SetupFailed', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyFundingTimeout()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.FAILED_SETUP)

      const eventList = await deployed.SystemStub.getPastEvents('SetupFailed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not awaiting a funding proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyFundingTimeout()
      } catch (e) {
        assert.include(e.message, 'Funding timeout has not started')
      }
    })

    it('reverts if the timeout has not elapsed', async () => {
      try {
        await testInstance.setKeepInfo(0, 0, fundingProofTimerStart * 5, utils.bytes32zero, utils.bytes32zero)
        await testInstance.notifyFundingTimeout()
      } catch (e) {
        assert.include(e.message, 'Funding timeout has not elapsed')
      }
    })

    it('distributes the funder bond to the keep group', async () => {
      const beneficiary = accounts[4]
      const bond = 1000000
      await testInstance.setKeepInfo(0, 0, withdrawalRequestTime, utils.bytes32zero, utils.bytes32zero)
      await testInstance.send(bond, { from: beneficiary })
      await testInstance.notifyFundingTimeout()
      const keepBalance = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(keepBalance, new BN(bond), 'funder bond not distributed properly')
    })
  })

  describe('provideFundingECDSAFraudProof', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1 // has elapsed
      await deployed.KeepStub.setSuccess(true)
      await testInstance.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart, utils.bytes32zero, utils.bytes32zero)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('updates to awaiting fraud funding proof and logs FraudDuringSetup if the timer has not elapsed', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart * 5, utils.bytes32zero, utils.bytes32zero) // timer has not elapsed
      await testInstance.provideFundingECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.FRAUD_AWAITING_BTC_FUNDING_PROOF)

      const keepState = await testInstance.getKeepInfo.call()
      assert(keepState[2].gtn(fundingProofTimerStart), 'fundingProofTimerStart did not increase')

      const eventList = await deployed.SystemStub.getPastEvents('FraudDuringSetup', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('updates to failed, logs FraudDuringSetup, and burns value if the timer has elapsed', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.provideFundingECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.FAILED_SETUP)

      const eventList = await deployed.SystemStub.getPastEvents('FraudDuringSetup', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not awaiting funding proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.provideFundingECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        assert.include(e.message, 'Signer fraud during funding flow only available while awaiting funding')
      }
    })

    it('reverts if the signature is not fraud', async () => {
      try {
        await deployed.KeepStub.setSuccess(false)
        await testInstance.provideFundingECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        await deployed.KeepStub.setSuccess(true)
        assert.include(e.message, 'Signature is not fraudulent')
      }
    })

    it('returns the funder bond if the timer has not elapsed', async () => {
      const beneficiary = accounts[4]
      const bond = 10000000000
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.send(bond, { from: beneficiary })
      const initialBalance = await web3.eth.getBalance(beneficiary)
      const signerBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await deployed.KeepStub.setBondAmount(bond)
      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart * 6, utils.bytes32zero, utils.bytes32zero)
      await testInstance.provideFundingECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')

      const finalBalance = await web3.eth.getBalance(beneficiary)
      const eventList = await deployed.SystemStub.getPastEvents('FraudDuringSetup', { fromBlock: blockNumber, toBlock: 'latest' })
      const balanceCheck = new BN(initialBalance).add(new BN(bond).add(new BN(signerBalance)))

      assert.equal(eventList.length, 1)
      expect(finalBalance, 'funder and signer bond should be included in final result').to.eq.BN(balanceCheck)
    })
  })

  describe('notifyFraudFundingTimeout', async () => {
    let timer

    before(async () => {
      timer = await deployed.TBTCConstants.getFraudFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1 // timer has elapsed
      await testInstance.setState(utils.states.FRAUD_AWAITING_BTC_FUNDING_PROOF)
      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart, utils.bytes32zero, utils.bytes32zero)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('updates state to setup failed, logs SetupFailed, and deconstes state', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyFraudFundingTimeout()

      const keepState = await testInstance.getKeepInfo.call()
      assert(keepState[0].eqn(0), 'Keep id not deconsted')
      assert(keepState[1].eqn(0), 'signingGroupRequestedAt not deconsted')
      assert(keepState[2].eqn(0), 'fundingProofTimerStart not deconsted')
      assert.equal(keepState[3], utils.bytes32zero) // pubkey X
      assert.equal(keepState[4], utils.bytes32zero) // pubkey Y

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.FAILED_SETUP)

      const eventList = await deployed.SystemStub.getPastEvents('SetupFailed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not awaiting fraud funding proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyFraudFundingTimeout()
      } catch (e) {
        assert.include(e.message, 'Not currently awaiting fraud-related funding proof')
      }
    })

    it('reverts if the timer has not elapsed', async () => {
      try {
        await testInstance.setKeepInfo(0, 0, fundingProofTimerStart * 5, utils.bytes32zero, utils.bytes32zero)
        await testInstance.notifyFraudFundingTimeout()
      } catch (e) {
        assert.include(e.message, 'Fraud funding proof timeout has not elapsed')
      }
    })

    it('asserts that it partially slashes signers', async () => {
      const beneficiary = accounts[4]
      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      const initialBalance = await web3.eth.getBalance(beneficiary)
      const toSeize = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.setKeepInfo(0, 0, fundingProofTimerStart, utils.bytes32zero, utils.bytes32zero)
      await testInstance.notifyFraudFundingTimeout()

      const divisor = await deployed.TBTCConstants.getFundingFraudPartialSlashDivisor.call()
      const slash = new BN(toSeize).div(new BN(divisor))
      const balanceAfter = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(slash)

      assert.equal(balanceCheck, balanceAfter, 'partial slash not correctly awarded to funder')
    })
  })

  describe('provideFraudBTCFundingProof', async () => {
    beforeEach(async () => {
      await testInstance.setKeepInfo(0, 0, 0, _signerPubkeyX, _signerPubkeyY)
      await deployed.SystemStub.setCurrentDiff(currentDifficulty)
      await testInstance.setState(utils.states.FRAUD_AWAITING_BTC_FUNDING_PROOF)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('updates to setup failed, logs SetupFailed, ', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.provideFraudBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      const keepState = await testInstance.getKeepInfo.call()
      assert(keepState[0].eqn(0), 'Keep id not deconsted')
      assert(keepState[1].eqn(0), 'signingGroupRequestedAt not deconsted')
      assert(keepState[2].eqn(0), 'fundingProofTimerStart not deconsted')
      assert.equal(keepState[3], utils.bytes32zero) // pubkey X
      assert.equal(keepState[4], utils.bytes32zero) // pubkey Y

      const depostState = await testInstance.getState.call()
      expect(depostState).to.eq.BN(utils.states.FAILED_SETUP)

      const eventList = await deployed.SystemStub.getPastEvents('SetupFailed', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not awaiting a funding proof during setup fraud', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.provideFraudBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
      } catch (e) {
        assert.include(e.message, 'Not awaiting a funding proof during setup fraud')
      }
    })

    it('assert distribute signer bonds to funder', async () => {
      const beneficiary = accounts[4]
      await deployed.SystemStub.setDepositOwner(0, beneficiary)

      const initialBalance = await web3.eth.getBalance(beneficiary)
      const signerBond = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.provideFraudBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      const balanceAfter = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(new BN(signerBond))

      assert.equal(balanceCheck, balanceAfter, 'partial slash not correctly awarded to funder')
    })
  })

  describe('provideBTCFundingProof', async () => {
    beforeEach(async () => {
      await testInstance.setKeepInfo(0, 0, 0, _signerPubkeyX, _signerPubkeyY)
      await deployed.SystemStub.setCurrentDiff(currentDifficulty)
      await testInstance.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('updates to active, stores UTXO info, deconstes funding info, logs Funded', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.provideBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      const UTXOInfo = await testInstance.getUTXOInfo.call()
      assert.equal(UTXOInfo[0], _outValueBytes)
      assert.equal(UTXOInfo[2], _expectedUTXOoutpoint)

      const keepState = await testInstance.getKeepInfo.call()
      assert(keepState[1].eqn(0), 'signingGroupRequestedAt not deconsted')
      assert(keepState[2].eqn(0), 'fundingProofTimerStart not deconsted')

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.ACTIVE)

      const eventList = await deployed.SystemStub.getPastEvents('Funded', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not awaiting funding proof', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.provideBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)
      } catch (e) {
        assert.include(e.message, 'Not awaiting funding')
      }
    })

    it('returns funder bonds and mints tokens', async () => {
      const beneficiary = accounts[4]
      const signerBond = 10000000000
      const initialTokenBalance = await deployed.TBTCStub.getBalance(beneficiary)

      await testInstance.send(signerBond, { from: beneficiary })
      await deployed.SystemStub.setDepositOwner(0, beneficiary)

      const initialBalance = await web3.eth.getBalance(beneficiary)

      await testInstance.provideBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      const balanceAfter = await web3.eth.getBalance(beneficiary)
      const balanceCheck = new BN(initialBalance).add(new BN(signerBond))

      assert.equal(balanceCheck, balanceAfter, 'funder bond not currectly returned')

      const endingTokenBalancce = await deployed.TBTCStub.getBalance(beneficiary)
      const lotSize = await deployed.TBTCConstants.getLotSize.call()
      const toMint = lotSize.mul(new BN(95)).div(new BN(100))
      const tokenCheck = initialTokenBalance.add(new BN(toMint))
      expect(tokenCheck, 'incorrect amount minted').to.eq.BN(endingTokenBalancce)
    })
  })

  describe('provideECDSAFraudProof', async () => {
    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('executes', async () => {
      await testInstance.provideECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
    })

    it('reverts if in the funding flow', async () => {
      try {
        await testInstance.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
        await testInstance.provideECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        assert.include(e.message, 'Use provideFundingECDSAFraudProof instead')
      }
    })

    it('reverts if already in signer liquidation', async () => {
      try {
        await testInstance.setState(utils.states.LIQUIDATION_IN_PROGRESS)
        await testInstance.provideECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        assert.include(e.message, 'Signer liquidation already in progress')
      }
    })

    it('reverts if the contract has halted', async () => {
      try {
        await testInstance.setState(utils.states.REDEEMED)
        await testInstance.provideECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        assert.include(e.message, 'Contract has halted')
      }
    })

    it('reverts if signature is not fraud according to Keep', async () => {
      try {
        await deployed.KeepStub.setSuccess(false)
        await testInstance.provideECDSAFraudProof(0, utils.bytes32zero, utils.bytes32zero, utils.bytes32zero, '0x00')
      } catch (e) {
        await deployed.KeepStub.setSuccess(true)
        assert.include(e.message, 'Signature is not fraud')
      }
    })

    it.skip('TODO: full test for startSignerFraudLiquidation', async () => {})
  })

  describe('provideSPVFraudProof', async () => {
    // real tx from mainnet bitcoin
    const currentDiff = 6353030562983
    // const txid = '0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f'
    // const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c'
    const tx = '0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800'
    const proof = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe04995ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b5095548'
    const index = 130
    const headerChain = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
    const outpoint = '0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000'
    const prevoutValueBytes = '0xf078351d00000000'
    const requesterPKH = '0x86e7303082a6a21d5837176bc808bf4828371ab6'

    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.SystemStub.setCurrentDiff(currentDiff)
      await testInstance.setUTXOInfo(prevoutValueBytes, 0, outpoint)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('executes', async () => {
      await testInstance.provideSPVFraudProof(tx, proof, index, headerChain)
    })

    it('reverts if in the funding flow', async () => {
      try {
        await testInstance.setState(utils.states.AWAITING_BTC_FUNDING_PROOF)
        await testInstance.provideSPVFraudProof('0x00', '0x00', 0, '0x00')
      } catch (e) {
        assert.include(e.message, 'SPV Fraud proofs not valid before Active state')
      }
    })

    it('reverts if already in signer liquidation', async () => {
      try {
        await testInstance.setState(utils.states.LIQUIDATION_IN_PROGRESS)
        await testInstance.provideSPVFraudProof('0x00', '0x00', 0, '0x00')
      } catch (e) {
        assert.include(e.message, 'Signer liquidation already in progress')
      }
    })

    it('reverts if the contract has halted', async () => {
      try {
        await testInstance.setState(utils.states.REDEEMED)
        await testInstance.provideSPVFraudProof('0x00', '0x00', 0, '0x00')
      } catch (e) {
        assert.include(e.message, 'Contract has halted')
      }
    })

    it('reverts if it can\'t verify the Deposit UTXO was consumed', async () => {
      try {
        await testInstance.setUTXOInfo(prevoutValueBytes, 0, '0x' + '00'.repeat(36))
        await testInstance.provideSPVFraudProof(tx, proof, index, headerChain)
      } catch (e) {
        assert.include(e.message, 'No input spending custodied UTXO found')
      }
    })

    it('reverts if it finds an output paying the redeemer', async () => {
      try {
        await testInstance.setRequestInfo(utils.address0, requesterPKH, 100, 0, utils.bytes32zero)
        await testInstance.provideSPVFraudProof(tx, proof, index, headerChain)
      } catch (e) {
        assert.include(e.message, 'Found an output paying the redeemer as requested')
      }
    })

    it.skip('TODO: full test for startSignerFraudLiquidation', async () => {})
  })

  describe('purchaseSignerBondsAtAuction', async () => {
    beforeEach(async () => {
      await testInstance.setState(utils.states.LIQUIDATION_IN_PROGRESS)
    })

    it('sets state to liquidated, logs Liquidated, ', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number
      await testInstance.purchaseSignerBondsAtAuction()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.LIQUIDATED)

      const eventList = await deployed.SystemStub.getPastEvents('Liquidated', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in a liquidation auction', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.purchaseSignerBondsAtAuction()
      } catch (e) {
        assert.include(e.message, 'No active auction')
      }
    })

    it('reverts if TBTC balance is insufficient', async () => {
      try {
        await deployed.TBTCStub.setReturnUint(0)
        await testInstance.purchaseSignerBondsAtAuction()
      } catch (e) {
        await deployed.TBTCStub.setReturnUint(new BN('1000000000000000000', 10))
        assert.include(e.message, 'Not enough TBTC to cover outstanding debt')
      }
    })

    it(`burns msg.sender's tokens`, async () => {
      const caller = accounts[4]
      const beneficiary = accounts[5]
      const lotSize = await deployed.TBTCConstants.getLotSize.call()
      const initialTokenBalance = await deployed.TBTCStub.getBalance(caller)

      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await deployed.TBTCStub.getBalance(caller)
      const tokenCheck = new BN(finalTokenBalance).add( new BN(lotSize) )
      expect(tokenCheck, 'tokens not burned correctly').to.eq.BN(initialTokenBalance)
    })

    it('distributes beneficiary reward', async () => {
      const caller = accounts[4]
      const beneficiary = accounts[5]
      const initialTokenBalance = await deployed.TBTCStub.getBalance(beneficiary)
      const returned = await deployed.TBTCStub.balanceOf.call(caller)

      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalTokenBalance = await deployed.TBTCStub.getBalance(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add( new BN(returned))

      expect(finalTokenBalance, 'tokens not returned to beneficiary correctly').to.eq.BN(tokenCheck)
    })

    it('distributes value to the caller', async () => {
      const value = 1000000000000
      const caller = accounts[4]
      const beneficiary = accounts[5]
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const initialBalance = await web3.eth.getBalance(caller)

      await testInstance.send(value, { from: accounts[0] })
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(caller)

      expect(new BN(finalBalance), 'caller balance should increase').to.be.gte.BN(initialBalance)
    })

    it('returns keep funds if not fraud', async () => {
      const value = 1000000000000
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const caller = accounts[4]
      const beneficiary = accounts[5]
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      assert(new BN(finalBalance).gtn(new BN(initialBalance)), 'caller balance should increase')
    })

    it('burns if fraud', async () => {
      const value = 1000000000000
      const block = await web3.eth.getBlock('latest')
      const notifiedTime = block.timestamp
      const caller = accounts[4]
      const beneficiary = accounts[5]
      const initialBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      await testInstance.send(value, { from: accounts[0] })
      await testInstance.setState(utils.states.FRAUD_LIQUIDATION_IN_PROGRESS)
      await testInstance.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await deployed.SystemStub.setDepositOwner(0, beneficiary)
      await testInstance.purchaseSignerBondsAtAuction({ from: caller })

      const finalBalance = await web3.eth.getBalance(deployed.KeepStub.address)

      expect(new BN(finalBalance)).to.eq.BN(initialBalance)
    })
  })

  describe('notifyCourtesyCall', async () => {
    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.SystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('sets courtesy call state, sets the timestamp, and logs CourtesyCalled', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.COURTESY_CALL)

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await deployed.SystemStub.getPastEvents('CourtesyCalled', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in active state', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Can only courtesy call from active state')
      }
    })

    it('reverts if sufficiently collateralized', async () => {
      try {
        await deployed.KeepStub.setBondAmount(1000)
        await deployed.SystemStub.setOraclePrice(new BN(1))
        await testInstance.notifyCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Signers have sufficient collateral')
      }
    })
  })

  describe('exitCourtesyCall', async () => {
    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      const notifiedTime = blockTimestamp // not expired
      const fundedTime = blockTimestamp // not expired
      await deployed.KeepStub.setBondAmount(new BN('1000000000000000000000000', 10))
      await deployed.SystemStub.setOraclePrice(new BN('1', 10))
      await testInstance.setState(utils.states.COURTESY_CALL)
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), fundedTime, '0x' + '00'.repeat(36))
      await testInstance.setLiquidationAndCourtesyInitated(0, notifiedTime)
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.SystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('transitions to active, and logs ExitedCourtesyCall', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.exitCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.ACTIVE)

      const eventList = await deployed.SystemStub.getPastEvents('ExitedCourtesyCall', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in courtesy call state', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Not currently in courtesy call')
      }
    })

    it('reverts if the deposit term is expiring anyway', async () => {
      try {
        await testInstance.setUTXOInfo('0x' + '00'.repeat(8), 0, '0x' + '00'.repeat(36))
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is expiring')
      }
    })

    it('reverts if the deposit is still undercollateralized', async () => {
      try {
        await deployed.SystemStub.setOraclePrice(new BN('1000000000000', 10))
        await deployed.KeepStub.setBondAmount(0)
        await testInstance.exitCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is still undercollateralized')
      }
    })
  })

  describe('notifyUndercollateralizedLiquidation', async () => {
    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await deployed.KeepStub.setBondAmount(0)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    afterEach(async () => {
      await deployed.KeepStub.setBondAmount(1000)
      await deployed.SystemStub.setOraclePrice(new BN('1000000000000', 10))
    })

    it('executes', async () => {
      await testInstance.notifyUndercollateralizedLiquidation()
    })

    it('reverts if not in active or courtesy call', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyUndercollateralizedLiquidation()
      } catch (e) {
        assert.include(e.message, 'Deposit not in active or courtesy call')
      }
    })

    it('reverts if the deposit is not severely undercollateralized', async () => {
      try {
        await deployed.KeepStub.setBondAmount(10000000)
        await deployed.SystemStub.setOraclePrice(new BN(1))
        await testInstance.notifyUndercollateralizedLiquidation()
      } catch (e) {
        assert.include(e.message, 'Deposit has sufficient collateral')
      }
    })

    it('assert starts signer abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifyUndercollateralizedLiquidation()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyCourtesyTimeout', async () => {
    let courtesyTime
    let timer
    before(async () => {
      timer = await deployed.TBTCConstants.getCourtesyCallTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      courtesyTime = blockTimestamp - timer.toNumber() // has not expired
      await testInstance.setState(utils.states.COURTESY_CALL)
      await testInstance.setLiquidationAndCourtesyInitated(0, courtesyTime)
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
    })

    it('executes', async () => {
      await testInstance.notifyCourtesyTimeout()
    })

    it('reverts if not in a courtesy call period', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyCourtesyTimeout()
      } catch (e) {
        assert.include(e.message, 'Not in a courtesy call period')
      }
    })

    it('reverts if the period has not elapsed', async () => {
      try {
        await testInstance.setLiquidationAndCourtesyInitated(0, courtesyTime * 5)
        await testInstance.notifyCourtesyTimeout()
      } catch (e) {
        assert.include(e.message, 'Courtesy period has not elapsed')
      }
    })

    it('assert starts signer abort liquidation', async () => {
      await deployed.KeepStub.send(1000000, { from: accounts[0] })
      await testInstance.notifyCourtesyTimeout()

      const bond = await web3.eth.getBalance(deployed.KeepStub.address)
      assert.equal(bond, 0, 'Bond not seized as expected')

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[0], 'liquidation timestamp not recorded').not.to.eq.BN(0)
    })
  })

  describe('notifyDepositExpiryCourtesyCall', async () => {
    let timer
    let fundedTime

    before(async () => {
      timer = await deployed.TBTCConstants.getCourtesyCallTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock('latest')
      const blockTimestamp = block.timestamp
      fundedTime = blockTimestamp - timer.toNumber() - 1 // has expired
      await testInstance.setState(utils.states.ACTIVE)
      await testInstance.setUTXOInfo('0x' + '00'.repeat(8), 0, '0x' + '00'.repeat(36))
    })

    it('sets courtesy call state, stores the time, and logs CourtesyCalled', async () => {
      const blockNumber = await web3.eth.getBlock('latest').number

      await testInstance.notifyDepositExpiryCourtesyCall()

      const depositState = await testInstance.getState.call()
      expect(depositState).to.eq.BN(utils.states.COURTESY_CALL)

      const liquidationTime = await testInstance.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await deployed.SystemStub.getPastEvents('CourtesyCalled', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList.length, 1)
    })

    it('reverts if not in active', async () => {
      try {
        await testInstance.setState(utils.states.START)
        await testInstance.notifyDepositExpiryCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit is not active')
      }
    })

    it('reverts if deposit not yet expiring', async () => {
      try {
        await testInstance.setUTXOInfo('0x' + '00'.repeat(8), fundedTime * 5, '0x' + '00'.repeat(36))
        await testInstance.notifyDepositExpiryCourtesyCall()
      } catch (e) {
        assert.include(e.message, 'Deposit term not elapsed')
      }
    })
  })
})
