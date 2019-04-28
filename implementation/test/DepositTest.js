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
const TestDeposit = artifacts.require('TestDeposit')

const BN = require('bn.js')
const utils = require('./utils')

TEST_DEPOSIT_DEPLOY = [
  {name: 'BytesLib', contract: BytesLib},
  {name: 'BTCUtils', contract: BTCUtils},
  {name: 'ValidateSPV', contract: ValidateSPV},
  {name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs},
  {name: 'TBTCConstants', contract: TestTBTCConstants},  // note the name
  {name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging},
  {name: 'DepositStates', contract: DepositStates},
  {name: 'DepositUtils', contract: DepositUtils},
  {name: 'DepositFunding', contract: DepositFunding},
  {name: 'DepositRedemption', contract: DepositRedemption},
  {name: 'DepositLiquidation', contract: DepositLiquidation},
  {name: 'TestDeposit', contract: TestDeposit},
  {name: 'KeepStub', contract: KeepStub},
  {name: 'TBTCStub', contract: TBTCStub},
  {name: 'SystemStub', contract: SystemStub}]


contract.only('Deposit', accounts => {

  let deployed
  let testInstance

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
      let blockNumber = await web3.eth.getBlock("latest").number
      await testInstance.createNewDeposit(
        deployed.SystemStub.address,
        deployed.TBTCStub.address,
        deployed.KeepStub.address,
        1,  //m
        1)

      // state updates
      res = await testInstance.getState.call()
      assert(res.eq(utils.states.AWAITING_SIGNER_SETUP), 'state not as expected')
      res = await testInstance.getKeepInfo.call()
      assert(res[0].eq(new BN(7)), 'keepID not as expected')
      assert(!res[1].eq(new BN(0)), 'signing group timestamp not as expected')  // signingGroupRequestedAt

      // fired an event
      eventList = await deployed.SystemStub.getPastEvents('Created', { fromBlock: blockNumber, toBlock: 'latest' })
      assert.equal(eventList[0].returnValues._keepID, '7')
    })

    it('reverts if not in the start state', async () => {
      await testInstance.setState(utils.states.REDEEMED)
      try {
        await testInstance.createNewDeposit.call(
          deployed.SystemStub.address,
          deployed.TBTCStub.address,
          deployed.KeepStub.address,
          1,  //m
          1)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Deposit setup already requested')
      }
    })
  })

  describe('requestRedemption', async () => {
    let sighash = '0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90'

    beforeEach(async () => {
      await testInstance.setState(utils.states.ACTIVE)
      await testInstance.setUTXOInfo('0x1111111111111111', 0, '0x' + '33'.repeat(36))
    })

    it('updates state successfully and fires a RedemptionRequested event', async () => {
      // the TX produced will be:
      // 01000000000101333333333333333333333333333333333333333333333333333333333333333333333333000000000001111111110000000016001433333333333333333333333333333333333333330000000000
      // the signer pkh will be:
      // 5eb9b5e445db673f0ed8935d18cd205b214e5187
      // == hash160(023333333333333333333333333333333333333333333333333333333333333333)
      // the sighash preimage will be:
      // 010000003fc8fd9fada5a3573744477d5e35b0d4d0645e42285e3dec25aac02078db0f838cb9012517c817fead650287d61bdd9c68803b6bf9c64133dcab3e65b5a50cb93333333333333333333333333333333333333333333333333333333333333333333333331976a9145eb9b5e445db673f0ed8935d18cd205b214e518788ac111111111111111100000000e4ca7a168bd64e3123edd7f39e1ab7d670b32311cac2dda8e083822139c7936c0000000001000000
      let blockNumber = await web3.eth.getBlock("latest").number
      await testInstance.setKeepInfo(0, 0, 0, '0x' + '33'.repeat(32), '0x' + '44'.repeat(32))

      // the fee is ~12,297,829,380 BTC
      await testInstance.requestRedemption('0x1111111100000000', '0x' + '33'.repeat(20))

      res = await testInstance.getRequestInfo()
      assert.equal(res[1], '0x' + '33'.repeat(20))  // requester pkh
      assert(!res[3].eqn(0)) // withdrawalRequestTime is set
      assert.equal(res[4], sighash)  // digest
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
      res = await deployed.KeepStub.wasDigestApprovedForSigning.call(0, sighash)
      assert(!res.eqn(0), 'digest was not approved')
    })

  })

  describe('provideRedemptionSignature', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('increaseRedemptionFee', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideRedemptionProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifySignatureTimeout', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyRedemptionProofTimeout', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifySignerSetupFailure', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('retrieveSignerPubkey', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyFundingTimeout', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideFundingECDSAFraudProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyFraudFundingTimeout', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideFraudBTCFundingProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideBTCFundingProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideECDSAFraudProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('provideSPVFraudProof', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('purchaseSignerBondsAtAuction', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyCourtesyCall', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('exitCourtesyCall', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyUndercollateralizedLiquidation', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyCourtesyTimeout', async () => {
    it.skip('is TODO', async () => {})
  })

  describe('notifyDepositExpiryCourtesyCall', async () => {
    it.skip('is TODO', async () => {})
  })
})
