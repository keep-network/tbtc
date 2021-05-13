const { deployAndLinkAll } = require('./helpers/testDeployer.js')
const { contract, web3, accounts } = require('@openzeppelin/test-environment')
const { states, fundingTx } = require('./helpers/utils.js')
const { BN, constants, expectRevert } = require('@openzeppelin/test-helpers')
const { ZERO_ADDRESS } = constants
const { expect } = require('chai')

const ECDSAKeepStub = contract.fromArtifact('ECDSAKeepStub')
const Deposit = contract.fromArtifact('Deposit')
const TestDeposit = contract.fromArtifact('TestDeposit')
const TBTCSystem = contract.fromArtifact('TBTCSystem')

describe('DepositFactory', async function() {
  let openKeepFee
  const fullBtc = 100000000

  describe('createDeposit()', async () => {
    let depositFactory
    let ecdsaKeepFactoryStub
    let mockSatWeiPriceFeed

    before(async () => {
      // To properly test createDeposit, we deploy the real Deposit contract and
      // make sure we don't get hit by the ACL hammer.
      ;({
        depositFactory,
        tbtcSystemStub,
        ecdsaKeepStub,
        ecdsaKeepFactoryStub,
      } = await deployAndLinkAll([], {
        TestDeposit: Deposit,
      }))
      openKeepFee = await tbtcSystemStub.getNewDepositFeeEstimate()
    })

    it('creates new clone instances', async () => {
      await tbtcSystemStub.setKeepAddress(ecdsaKeepStub.address)
      await ecdsaKeepStub.send(openKeepFee)
      await ecdsaKeepStub.setBondAmount(openKeepFee)

      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })

      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })

      const eventList = await depositFactory.getPastEvents(
        'DepositCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' },
      )

      expect(eventList.length).to.equal(2)

      expect(
        web3.utils.isAddress(eventList[0].returnValues.depositCloneAddress),
      ).to.be.true
      expect(
        web3.utils.isAddress(eventList[1].returnValues.depositCloneAddress),
      ).to.be.true

      expect(
        eventList[0].returnValues.depositCloneAddress,
        'clone addresses should not be equal',
      ).to.not.equal(eventList[1].returnValues.depositCloneAddress)
    })

    it('correctly forwards value to keep factory', async () => {
      ;({
        ecdsaKeepFactoryStub,
        depositFactory,
        mockSatWeiPriceFeed,
      } = await deployAndLinkAll([], {
        TestDeposit: Deposit,
        TBTCSystemStub: TBTCSystem,
      }))
      await ecdsaKeepFactoryStub.setKeepAddress(ecdsaKeepStub.address)

      await mockSatWeiPriceFeed.setPrice(new BN('1000000000000', 10))

      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })
      expect(
        await web3.eth.getBalance(ecdsaKeepFactoryStub.address),
        'Factory did not correctly forward value on Deposit creation',
      ).to.eq.BN(openKeepFee)
    })

    it('reverts if insufficient fee is provided', async () => {
      await tbtcSystemStub.setKeepAddress(ecdsaKeepStub.address)

      const badOpenKeepFee = openKeepFee.sub(new BN('1'))
      await expectRevert(
        depositFactory.createDeposit(fullBtc, { value: badOpenKeepFee }),
        'Insufficient value for new keep creation',
      )
    })
  })

  describe('clone state', async () => {
    let mockRelay
    let tbtcSystemStub
    let tbtcToken
    let tbtcDepositToken
    let testDeposit
    let depositFactory

    const publicKey =
      '0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1'

    before(async () => {
      ;({
        mockRelay,
        tbtcSystemStub,
        tbtcToken,
        tbtcDepositToken,
        testDeposit,
        depositFactory,
        ecdsaKeepStub,
        ecdsaKeepFactoryStub,
      } = await deployAndLinkAll([]))
      openKeepFee = await tbtcSystemStub.getNewDepositFeeEstimate()
    })

    it('is not affected by state changes to other clone', async () => {
      const keep1 = await ECDSAKeepStub.new()
      const keep2 = await ECDSAKeepStub.new()
      const blockNumber = await web3.eth.getBlockNumber()

      await tbtcSystemStub.setKeepAddress(keep1.address)
      await keep1.send(openKeepFee)
      await keep1.setBondAmount(openKeepFee)
      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })

      await tbtcSystemStub.setKeepAddress(keep2.address)
      await keep2.send(openKeepFee)
      await keep2.setBondAmount(openKeepFee)
      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })

      const eventList = await depositFactory.getPastEvents(
        'DepositCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' },
      )

      const clone1 = eventList[0].returnValues.depositCloneAddress
      const clone2 = eventList[1].returnValues.depositCloneAddress
      const deposit1 = await TestDeposit.at(clone1)
      const deposit2 = await TestDeposit.at(clone2)

      await deposit1.setKeepAddress(keep1.address)

      await deposit2.setKeepAddress(keep2.address)

      await keep1.setPublicKey(publicKey)
      await keep2.setPublicKey(publicKey)

      await deposit1.retrieveSignerPubkey()
      await deposit2.retrieveSignerPubkey()
      await mockRelay.setCurrentEpochDifficulty(fundingTx.difficulty)
      await mockRelay.setPrevEpochDifficulty(fundingTx.difficulty)
      await deposit2.provideBTCFundingProof(
        fundingTx.version,
        fundingTx.txInputVector,
        fundingTx.txOutputVector,
        fundingTx.txLocktime,
        fundingTx.fundingOutputIndex,
        fundingTx.merkleProof,
        fundingTx.txIndexInBlock,
        fundingTx.bitcoinHeaders,
      )

      // deposit1 should be AWAITING_BTC_FUNDING_PROOF (2)
      // deposit2 should be ACTIVE (5)
      const deposit1state = await deposit1.currentState()
      const deposit2state = await deposit2.currentState()

      expect(
        deposit1state,
        'Deposit 1 should be in AWAITING_BTC_FUNDING_PROOF',
      ).to.eq.BN(states.AWAITING_BTC_FUNDING_PROOF)
      expect(deposit2state, 'Deposit 2 should be in ACTIVE').to.eq.BN(
        states.ACTIVE,
      )
    })

    it('is not affected by state changes to master', async () => {
      const keep = await ECDSAKeepStub.new()

      await tbtcDepositToken.forceMint(accounts[0], testDeposit.address)
      await testDeposit.initializeDeposit(
        tbtcSystemStub.address,
        tbtcToken.address,
        tbtcDepositToken.address,
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        fullBtc,
        { value: openKeepFee },
      )

      await testDeposit.setKeepAddress(keep.address)

      await keep.setPublicKey(publicKey)

      await testDeposit.retrieveSignerPubkey()

      // master deposit should now be in AWAITING_BTC_FUNDING_PROOF
      const masterState = await testDeposit.currentState()

      const blockNumber = await web3.eth.getBlockNumber()

      await depositFactory.createDeposit(fullBtc, { value: openKeepFee })

      const eventList = await depositFactory.getPastEvents(
        'DepositCloneCreated',
        { fromBlock: blockNumber, toBlock: 'latest' },
      )
      const cloneNew = eventList[0].returnValues.depositCloneAddress
      const depositNew = await TestDeposit.at(cloneNew)

      // should be behind Master, at AWAITING_SIGNER_SETUP
      const newCloneState = await depositNew.currentState()

      expect(
        masterState,
        'Master deposit should be in AWAITING_BTC_FUNDING_PROOF',
      ).to.eq.BN(states.AWAITING_BTC_FUNDING_PROOF)
      expect(
        newCloneState,
        'New clone should be in AWAITING_SIGNER_SETUP',
      ).to.eq.BN(states.AWAITING_SIGNER_SETUP)
    })
  })
})
