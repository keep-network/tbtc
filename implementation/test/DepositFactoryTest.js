import deployTestDeposit from './helpers/deployTestDeposit'

import BN from 'bn.js'
import utils from './utils'
import chai, { expect } from 'chai'
import bnChai from 'bn-chai'
chai.use(bnChai(BN))

const ECDSAKeepStub = artifacts.require('ECDSAKeepStub')
const TestDeposit = artifacts.require('TestDeposit')

contract('DepositFactory', () => {
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let testDeposit
  let depositFactory

  const funderBondAmount = new BN('10').pow(new BN('5'))
  const fullBtc = 100000000

  before(async () => {
    ({
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      testDeposit,
      depositFactory,
    } = await deployTestDeposit())
  })

  describe('createDeposit()', async () => {
    it('creates new clone instances', async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )

      const eventList = await depositFactory.getPastEvents('DepositCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })

      assert.equal(eventList.length, 2)
      assert(web3.utils.isAddress(eventList[0].returnValues.depositCloneAddress))
      assert(web3.utils.isAddress(eventList[1].returnValues.depositCloneAddress))
      assert.notEqual(eventList[0].returnValues.depositCloneAddress, eventList[1].returnValues.depositCloneAddress, 'clone addresses should not be equal')
    })

    it('correctly forwards value to Deposit', async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )
      const eventList = await depositFactory.getPastEvents(
        'DepositCloneCreated',
        {
          fromBlock: blockNumber,
          toBlock: 'latest',
        })

      const depositAddress = eventList[eventList.length - 1].returnValues.depositCloneAddress

      const balance = await web3.eth.getBalance(depositAddress)
      assert.equal(balance, funderBondAmount, 'Factory did not correctly forward value on Deposit creation')
    })
  })

  describe('clone state', async () => {
    const publicKey = '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'

    it('is not affected by state changes to other clone', async () => {
      const keep1 = await ECDSAKeepStub.new()
      const keep2 = await ECDSAKeepStub.new()
      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )

      const eventList = await depositFactory.getPastEvents('DepositCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })

      const clone1 = eventList[0].returnValues.depositCloneAddress
      const clone2 = eventList[1].returnValues.depositCloneAddress
      const deposit1 = await TestDeposit.at(clone1)
      const deposit2 = await TestDeposit.at(clone2)

      await deposit1.setKeepAddress(keep1.address)

      await deposit2.setKeepAddress(keep2.address)

      const currentDifficulty = 6353030562983
      const _version = '0x01000000'
      const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
      const _txOutputVector = '0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6'
      const _fundingOutputIndex = 0
      const _txLocktime = '0x4ec10800'
      const _txIndexInBlock = 129
      const _bitcoinHeaders = '0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e'
      const _merkleProof = '0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049'
      // const _outputValue = 490029088;

      await keep1.setPublicKey(publicKey)
      await keep2.setPublicKey(publicKey)

      await deposit1.retrieveSignerPubkey()
      await deposit2.retrieveSignerPubkey()
      await tbtcSystemStub.setCurrentDiff(currentDifficulty)
      await deposit2.provideBTCFundingProof(_version, _txInputVector, _txOutputVector, _txLocktime, _fundingOutputIndex, _merkleProof, _txIndexInBlock, _bitcoinHeaders)

      // deposit1 should be AWAITING_BTC_FUNDING_PROOF (2)
      // deposit2 should be ACTIVE (5)
      const deposit1state = await deposit1.getCurrentState()
      const deposit2state = await deposit2.getCurrentState()

      expect(deposit1state, 'Deposit 1 should be in AWAITING_BTC_FUNDING_PROOF').to.eq.BN(utils.states.AWAITING_BTC_FUNDING_PROOF)
      expect(deposit2state, 'Deposit 2 should be in ACTIVE').to.eq.BN(utils.states.ACTIVE)
    })

    it('is not affected by state changes to master', async () => {
      const keep = await ECDSAKeepStub.new()

      await testDeposit.createNewDeposit(
        tbtcSystemStub.address,
        tbtcToken.address,
        tbtcDepositToken.address,
        utils.address0,
        utils.address0,
        1,
        1,
        fullBtc,
        { value: funderBondAmount }
      )

      await testDeposit.setKeepAddress(keep.address)

      await keep.setPublicKey(publicKey)

      await testDeposit.retrieveSignerPubkey()

      // master deposit should now be in AWAITING_BTC_FUNDING_PROOF
      const masterState = await testDeposit.getCurrentState()

      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(
        fullBtc,
        { value: funderBondAmount }
      )

      const eventList = await depositFactory.getPastEvents('DepositCloneCreated', { fromBlock: blockNumber, toBlock: 'latest' })
      const cloneNew = eventList[0].returnValues.depositCloneAddress
      const depositNew = await TestDeposit.at(cloneNew)

      // should be behind Master, at AWAITING_SIGNER_SETUP
      const newCloneState = await depositNew.getCurrentState()

      expect(masterState, 'Master deposit should be in AWAITING_BTC_FUNDING_PROOF').to.eq.BN(utils.states.AWAITING_BTC_FUNDING_PROOF)
      expect(newCloneState, 'New clone should be in AWAITING_SIGNER_SETUP').to.eq.BN(utils.states.AWAITING_SIGNER_SETUP)
    })
  })
})
