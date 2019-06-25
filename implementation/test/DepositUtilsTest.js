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
const TBTC = artifacts.require('TBTC')
const SystemStub = artifacts.require('SystemStub')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDepositUtils = artifacts.require('TestDepositUtils')

const UniswapDeployment = artifacts.require('UniswapDeployment')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

const BN = require('bn.js')
const utils = require('./utils')

const TEST_DEPOSIT_UTILS_DEPLOY = [
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
  {name: 'TestDepositUtils', contract: TestDepositUtils},
  {name: 'KeepStub', contract: KeepStub},
  {name: 'TBTCStub', contract: TBTCStub},
  {name: 'SystemStub', contract: SystemStub}];


import { UniswapHelpers } from './helpers/uniswap'

contract('DepositUtils', accounts => {

  let deployed
  let testUtilsInstance

  before(async () => {
      deployed = await utils.deploySystem(TEST_DEPOSIT_UTILS_DEPLOY)
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
      let res = await testUtilsInstance.currentBlockDifficulty.call()
      assert(res.eq(new BN(1)))

      await deployed.SystemStub.setCurrentDiff(33)
      res = await testUtilsInstance.currentBlockDifficulty.call()
      assert(res.eq(new BN(33)))
    })
  })

  describe('previousBlockDifficulty()', async () => {
    it('calls out to the system', async () => {
      let res = await testUtilsInstance.previousBlockDifficulty.call()
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
        let res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'not at current or previous difficulty')
      }
    })

    it('evaluates a header proof with previous', async () => {
      await deployed.SystemStub.setPreviousDiff(5646403851534)
      let res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])

    })

    it('evaluates a header proof with current', async () => {
      await deployed.SystemStub.setPreviousDiff(1)
      await deployed.SystemStub.setCurrentDiff(5646403851534)
      let res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0])
    })

    it('reverts on low difficulty', async () => {
      try {
        let res = await testUtilsInstance.evaluateProofDifficulty(utils.HEADER_PROOFS[0].slice(0, 160 * 4 + 2))
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'Insufficient accumulated difficulty in header chain')
      }
    })

    it('reverts on a ValidateSPV error code (3 or lower)', async () => {
      try {
        await deployed.SystemStub.setPreviousDiff(1)
        let res = await testUtilsInstance.evaluateProofDifficulty(utils.LOW_DIFF_HEADER)
        assert(false, 'Test call did not error as expected')
      } catch (e) {
        assert.include(e.message, 'ValidateSPV returned an error code')
      }
    })
  })

  describe('checkProof()', async () => {
    it('returns the correct _txid' , async () => {
      await deployed.SystemStub.setCurrentDiff(6379265451411)
      let res = await testUtilsInstance.checkProof.call(utils.TX.tx, utils.TX.proof, utils.TX.index, utils.HEADER_PROOFS.slice(-1)[0])
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
      let res = await testUtilsInstance.signerFee.call()
      assert(res.eq(new BN(500000)))
    })
  })

  describe('beneficiaryReward()', async () => {
    it('returns a derived constant', async () => {
      let res = await testUtilsInstance.beneficiaryReward.call()
      assert(res.eq(new BN(10 ** 5)))
    })
  })

  describe('determineCompressionPrefix()', async () => {
    it('selects 2 for even', async () => {
      let res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(32))
      assert.equal(res, '0x02')
    })
    it('selects 3 for odd', async () => {
      let res = await testUtilsInstance.determineCompressionPrefix.call('0x' + '00'.repeat(31) + '01')
      assert.equal(res, '0x03')
    })
  })

  describe('compressPubkey()', async () => {
    it('returns a 33 byte array with a prefix', async () => {
      let res = await testUtilsInstance.compressPubkey.call('0x' + '00'.repeat(32), '0x' + '00'.repeat(32))
      assert.equal(res, '0x02' + '00'.repeat(32))
    })
  })

  describe('signerPubkey()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      let res = await testUtilsInstance.signerPubkey.call()
      assert.equal(res, '0x' + '00'.repeat(64))
    })
  })

  describe('signerPKH()', async () => {
    it('returns the concatenated signer X and Y coordinates', async () => {
      let res = await testUtilsInstance.signerPKH.call()
      assert.equal(res, utils.hash160('02' + '00'.repeat(32)))
    })
  })

  describe('utxoSize()', async () => {
    it('returns the state\'s utxoSizeBytes as an integer', async () => {
      let res = await testUtilsInstance.utxoSize.call()
      assert(res.eq(new BN(0)))

      await testUtilsInstance.setUTXOInfo('0x11223344', 1, '0x')
      res = await testUtilsInstance.utxoSize.call()
      assert(res.eq(new BN('44332211', 16)))
    })
  })

  describe('fetchOraclePrice()', async () => {
    it('calls out to the system', async () => {
      let res = await testUtilsInstance.fetchOraclePrice.call()
      assert(res.eq(new BN('1000000000000', 10)))

      await deployed.SystemStub.setOraclePrice(44)
      res = await testUtilsInstance.fetchOraclePrice.call()
      assert(res.eq(new BN(44)))
    })
  })

  describe('fetchBondAmount()', async () => {
    it('calls out to the keep system', async () => {
      let res = await testUtilsInstance.fetchBondAmount.call()
      assert(res.eq(new BN(10000)))

      await deployed.KeepStub.setBondAmount(44)
      res = await testUtilsInstance.fetchBondAmount.call()
      assert(res.eq(new BN(44)))
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
      let res = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert.equal(res, false)

      await deployed.KeepStub.approveDigest(7, '0x' + '00'.repeat(32))
      res = await testUtilsInstance.wasDigestApprovedForSigning.call('0x' + '00'.repeat(32))
      assert(res.eq(new BN(100)))
    })
  })

  describe('depositBeneficiary()', async () => {
    it('calls out to the system', async () => {
      let res = await testUtilsInstance.depositBeneficiary.call()
      assert.equal(res, '0x' + '00'.repeat(19) + '00')
    })
  })

  describe('redemptionTeardown()', async () => {
    it('deletes state', async () => {
      await testUtilsInstance.setRequestInfo('0x' + '11'.repeat(20), '0x' + '11'.repeat(20), 5, 6, '0x' + '33'.repeat(32))
      let res = await testUtilsInstance.getRequestInfo.call()
      assert.equal(res[4], '0x' + '33'.repeat(32))
      await testUtilsInstance.redemptionTeardown()
      res = await testUtilsInstance.getRequestInfo.call()
      assert.equal(res[4], '0x' + '00'.repeat(32))
    })
  })

  describe('seizeSignerBonds()', async () => {
    it('calls out to the keep system and returns the seized amount', async () => {
      let value = 5000
      await deployed.KeepStub.send(value, {from: accounts[0]})
      let res = await testUtilsInstance.seizeSignerBonds.call()
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
      let value = 10000
      await testUtilsInstance.send(value, {from: accounts[0]})
      await testUtilsInstance.pushFundsToKeepGroup(value)
      let keepBalance = await web3.eth.getBalance(deployed.KeepStub.address)
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


  describe('attemptToLiquidateOnchain', async () => {
    // see https://github.com/Uniswap/contracts-vyper/issues/27#issuecomment-476148467
    const UNISWAP_GAS_ESTIMATE = {
      addLiquidity: 21115 + 10000
    }

    it.only('works', async () => {
      // Contracts
      const deposit = deployed.TestDepositUtils;
      const keep = deployed.KeepStub;
      const tbtc = await TBTC.new()
      
      // Actors
      const tbtcSeller = accounts[1]

      // Token amounts
      const ethLiquidated = web3.utils.toWei('10', 'ether')
      const tbtcAmt = web3.utils.toWei('1', 'ether') // TODO(liamz): decimal places of TBTC
      const tbtcBought = tbtcAmt


      // Setup
      await tbtc.mint(tbtcSeller, tbtcAmt)
      
      await testUtilsInstance.setExteroriorAddresses(
        deployed.SystemStub.address,
        tbtc.address,
        deployed.KeepStub.address,
      )

      // Deploy Uniswap exchange for TBTC
      const uniswapDeployment = await UniswapDeployment.deployed()
      const uniswapFactory = await IUniswapFactory.at(await uniswapDeployment.factory())
      await uniswapFactory.createExchange(tbtc.address)

      const exchange = await IUniswapExchange.at(
        await uniswapFactory.getExchange(tbtc.address)
      )

      expect(
        await web3.eth.getBalance(exchange.address)
      ).to.eq('0')
      
      // Now the tbtcSeller adds liquidity
      await tbtc.approve(
        exchange.address, 
        web3.utils.toWei('10000', 'ether'), 
        { from: tbtcSeller }
      )

      await exchange.addLiquidity(
        '0',
        tbtcAmt,
        UniswapHelpers.getDeadline(),
        { from: tbtcSeller, value: ethLiquidated, gas: UNISWAP_GAS_ESTIMATE.addLiquidity }
      )

      
      // add some TBTC to the exchange
      // getLotSize
      // let lotSize = 1;

      // // Send some ETH to the Keep
      // // TODO(liamz): move this into test of startSignerFraudLiquidation etc.
      // //              for now it's a good placeholder for assumptions
      // // const keepBondAmount = await keep.checkBondAmount()
      const bondAmount = await deposit.fetchBondAmount()
      assert(bondAmount.gt(0))
      // // expect(await web3.eth.getBalance(keep.address)).to.eq.BN(0)
      // await keep.send(bondAmount, {from: accounts[0]})
      await deposit.send(bondAmount, {from: accounts[0]})

      // // Pretend the deposit has now seized the signer bonds
      // await deposit.setState(utils.states.ACTIVE)
      
      // // Liquidate
      // // TODO(liamz): make a mock uniswap exchange contract + set the price
      // // 1000000000000000000
      await deposit.attemptToLiquidateOnchain()
      
      // // Assert balances of 
      // // beneficiary     who initially deposited btc for tbtc
      // // tbtcSeller      who performs liquidation
      // // signers         who are the keep group
      // // deposit         the deposit contract
      expect(
        await tbtc.balanceOf(deposit.address)
      ).to.eq(tbtcBought)

      expect(
        await web3.eth.getBalance(deposit.address)
      ).to.eq('0')


      expect(
        await tbtc.balanceOf(tbtcSeller)
      ).to.eq(new BN(0))

      expect(
        web3.eth.getBalance(tbtcSeller)
      ).to.eq('0')


      /**
       * deposit:
       *  0 ETH
       *  X TBTC
       * 
       * keep/signers
       *  -X ETH
       *  0 TBTC
       * 
       * tbtcSeller:
       *   X ETH
       *  -X TBTC
       * 
       * beneficiary:
       *   0 ETH
       *   0 TBTC
       */
    })
  })
})
