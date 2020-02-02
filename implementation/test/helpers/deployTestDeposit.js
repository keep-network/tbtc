import utils from '../utils'
import BN from 'bn.js'

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

const ECDSAKeepStub = artifacts.require('ECDSAKeepStub')
const KeepRegistryStub = artifacts.require('KeepRegistryStub')

const TestToken = artifacts.require('TestToken')
const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTCDepositToken = artifacts.require('TestTBTCDepositToken')
const FeeRebateToken = artifacts.require('TestFeeRebateToken')
const DepositFactory = artifacts.require('DepositFactory')
const VendingMachine = artifacts.require('VendingMachine')

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')

export const TEST_DEPOSIT_DEPLOY = [
  { name: 'TBTCSystemStub', contract: TBTCSystemStub, constructorParam: utils.address0 },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name
  { name: 'DepositFactory', contract: DepositFactory, constructorParam: 'TBTCSystemStub' },
  { name: 'VendingMachine', contract: VendingMachine, constructorParam: 'TBTCSystemStub' },
  { name: 'BytesLib', contract: BytesLib },
  { name: 'BTCUtils', contract: BTCUtils },
  { name: 'ValidateSPV', contract: ValidateSPV },
  { name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs },
  { name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'TestDeposit', contract: TestDeposit, constructorParam: 'DepositFactory' },
  { name: 'TBTCDepositToken', contract: TBTCDepositToken, constructorParam: 'DepositFactory' },
  { name: 'FeeRebateToken', contract: FeeRebateToken, constructorParam: 'VendingMachine' },
  { name: 'ECDSAKeepStub', contract: ECDSAKeepStub },
]

/**
 * Deploys a test deposit setup and returns a resulting set of contracts. Two
 * methods are provided to customize the deployment functionality: a way to
 * substitute a different contract (for example, a mock or test version) for a
 * given contract name, and a way to specify additional contracts to deploy.
 *
 * If a VendingMachine contract is specified as an additional contract, its
 * address is passed to the TBTCSystem contract. Additional contracts are always
 * deployed _after_ the default contracts this function deploys.
 *
 * @param {*} additions A list of {name, contract} objects to deploy in addition
 *   to the default contracts deployed by this function. If a contract named
 *   VendingMachine is specified, its address is used for the deposit as well.
 * @param {*} substitutions An object mapping a contract name to a different
 *   contract type than the default one used by this function. For
 *   composability, substitutions are applied after the additions are added to
 *   the list of contracts to deploy.
 *
 * @return {object} An object with properties for all deployed contracts listed here:
 *    - tbtcConstants
 *    - mockRelay
 *    - tbtcSystemStub
 *    - tbtcToken
 *    - tbtcDepositToken
 *    - feeRebateToken
 *    - testDeposit
 *    - depositUtils
 *    - keepRegistryStub
 *    - ecdsaKeepStub
 *    - depositFactory
 *    Additionally, the object contains a `deployed` property that holds
 *    references to all deployed contracts by specified name.
 */
export default async function deployTestDeposit(
  additions = [],
  substitutions = {},
) {
  const deployment = TEST_DEPOSIT_DEPLOY.concat(additions)
  for (let i = 0; i < deployment.length; ++i) {
    const substitution = substitutions[deployment[i].name]
    if (substitution) {
      deployment[i] = {
        name: deployment[i].name,
        contract: substitution,
        constructorParam: deployment[i].constructorParam,
      }
    }
  }

  const deployed = await utils.deploySystem(deployment)

  const tbtcConstants = deployed.TBTCConstants

  const vendingMachine = deployed.VendingMachine

  const tbtcSystemStub = deployed.TBTCSystemStub
  const keepRegistryStub = await KeepRegistryStub.new()

  const tbtcToken = await TestToken.new(vendingMachine.address)
  const testDeposit = deployed.TestDeposit

  const tbtcDepositToken = deployed.TBTCDepositToken
  const feeRebateToken = deployed.FeeRebateToken
  const depositUtils = deployed.DepositUtils
  const ecdsaKeepStub = deployed.ECDSAKeepStub
  const depositFactory = deployed.DepositFactory

  await testDeposit.setExteriorAddresses(
    tbtcSystemStub.address,
    tbtcToken.address,
    tbtcDepositToken.address,
    feeRebateToken.address,
    vendingMachine.address,
  )

  await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  await testDeposit.setLotSize(new BN('100000000'))

  await tbtcSystemStub.initialize(
    keepRegistryStub.address,
    depositFactory.address,
    testDeposit.address,
    tbtcToken.address,
    tbtcDepositToken.address,
    feeRebateToken.address,
    vendingMachine.address,
    1,
    1
  )

  return {
    tbtcConstants,
    tbtcSystemStub,
    tbtcToken,
    tbtcDepositToken,
    feeRebateToken,
    testDeposit,
    depositUtils,
    keepRegistryStub,
    ecdsaKeepStub,
    depositFactory,
    deployed,
  }
}
