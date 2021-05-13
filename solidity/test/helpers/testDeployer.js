const { contract } = require('@openzeppelin/test-environment')
const { BN } = require('@openzeppelin/test-helpers')
const { deploySystem } = require('./utils.js')

const BytesLib = contract.fromArtifact('BytesLib')
const BTCUtils = contract.fromArtifact('BTCUtils')
const ValidateSPV = contract.fromArtifact('ValidateSPV')
const CheckBitcoinSigs = contract.fromArtifact('CheckBitcoinSigs')
const OutsourceDepositLogging = contract.fromArtifact('OutsourceDepositLogging')
const DepositStates = contract.fromArtifact('DepositStates')
const DepositUtils = contract.fromArtifact('DepositUtils')
const DepositFunding = contract.fromArtifact('DepositFunding')
const DepositRedemption = contract.fromArtifact('DepositRedemption')
const DepositLiquidation = contract.fromArtifact('DepositLiquidation')
const ECDSAKeepStub = contract.fromArtifact('ECDSAKeepStub')
const ECDSAKeepFactoryStub = contract.fromArtifact('ECDSAKeepFactoryStub')
const TestTBTCToken = contract.fromArtifact('TestTBTCToken')
const MockRelay = contract.fromArtifact('MockRelay')
const MockSatWeiPriceFeed = contract.fromArtifact('MockSatWeiPriceFeed')
const KeepFactorySelection = contract.fromArtifact('KeepFactorySelection')
const KeepFactorySelectorStub = contract.fromArtifact('KeepFactorySelectorStub')
const TBTCSystemStub = contract.fromArtifact('TBTCSystemStub')
const TBTCDepositToken = contract.fromArtifact('TestTBTCDepositToken')
const FeeRebateToken = contract.fromArtifact('TestFeeRebateToken')
const DepositFactory = contract.fromArtifact('DepositFactory')
const VendingMachine = contract.fromArtifact('VendingMachine')
const TestTBTCConstants = contract.fromArtifact('TestTBTCConstants')
const TestDeposit = contract.fromArtifact('TestDeposit')
const RedemptionScript = contract.fromArtifact('RedemptionScript')
const FundingScript = contract.fromArtifact('FundingScript')

const TEST_DEPOSIT_DEPLOY = [
  { name: 'KeepFactorySelection', contract: KeepFactorySelection },
  { name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging },
  { name: 'MockRelay', contract: MockRelay },
  { name: 'MockSatWeiPriceFeed', contract: MockSatWeiPriceFeed },
  {
    name: 'TBTCSystemStub',
    contract: TBTCSystemStub,
    constructorParams: ['MockSatWeiPriceFeed', 'MockRelay'],
  },
  {
    name: 'DepositFactory',
    contract: DepositFactory,
    constructorParams: ['TBTCSystemStub'],
  },
  {
    name: 'VendingMachine',
    contract: VendingMachine,
    constructorParams: ['TBTCSystemStub'],
  },
  { name: 'DepositStates', contract: DepositStates },
  { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name
  { name: 'DepositUtils', contract: DepositUtils },
  { name: 'DepositRedemption', contract: DepositRedemption },
  { name: 'DepositLiquidation', contract: DepositLiquidation },
  { name: 'DepositFunding', contract: DepositFunding },
  { name: 'TestDeposit', contract: TestDeposit },
  { name: 'BytesLib', contract: BytesLib },
  { name: 'BTCUtils', contract: BTCUtils },
  { name: 'ValidateSPV', contract: ValidateSPV },
  { name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs },
  { name: 'ECDSAKeepFactoryStub', contract: ECDSAKeepFactoryStub },
  {
    name: 'TBTCDepositToken',
    contract: TBTCDepositToken,
    constructorParams: ['DepositFactory'],
  },
  {
    name: 'FeeRebateToken',
    contract: FeeRebateToken,
    constructorParams: ['VendingMachine'],
  },
  {
    name: 'ECDSAKeepStub',
    contract: ECDSAKeepStub,
    constructorParams: ['VendingMachine', ''],
  },
  {
    name: 'KeepFactorySelectorStub',
    contract: KeepFactorySelectorStub,
    constructorParams: ['ECDSAKeepFactoryStub'],
  },
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
 *    - ecdsaKeepStub
 *    - depositFactory
 *    Additionally, the object contains a `deployed` property that holds
 *    references to all deployed contracts by specified name.
 */
async function deployAndLinkAll(additions = [], substitutions = {}) {
  const deployment = TEST_DEPOSIT_DEPLOY.concat(additions)
  for (let i = 0; i < deployment.length; ++i) {
    const substitution = substitutions[deployment[i].name]
    if (substitution) {
      deployment[i] = {
        name: deployment[i].name,
        contract: substitution,
        constructorParams: deployment[i].constructorParams,
      }
    }
  }

  const deployed = await deploySystem(deployment)

  const tbtcConstants = deployed.TBTCConstants

  const vendingMachine = deployed.VendingMachine

  const mockRelay = deployed.MockRelay
  await mockRelay.setPrevEpochDifficulty(1)

  const tbtcSystemStub = deployed.TBTCSystemStub
  const keepFactorySelectorStub = deployed.KeepFactorySelectorStub
  const ecdsaKeepFactoryStub = deployed.ECDSAKeepFactoryStub

  const tbtcToken = await TestTBTCToken.new(vendingMachine.address)
  const testDeposit = deployed.TestDeposit

  const tbtcDepositToken = deployed.TBTCDepositToken
  const feeRebateToken = deployed.FeeRebateToken
  const depositUtils = deployed.DepositUtils
  const ecdsaKeepStub = deployed.ECDSAKeepStub
  const depositFactory = deployed.DepositFactory
  const mockSatWeiPriceFeed = deployed.MockSatWeiPriceFeed
  const redemptionScript = await RedemptionScript.new(
    vendingMachine.address,
    tbtcToken.address,
    feeRebateToken.address,
  )
  const fundingScript = await FundingScript.new(
    vendingMachine.address,
    tbtcToken.address,
    tbtcDepositToken.address,
    feeRebateToken.address,
  )
  if (testDeposit.setExteriorAddresses) {
    // Test setup if this is in fact a TestDeposit. If it's been substituted
    // with e.g. Deposit, we don't set it up.
    await testDeposit.setExteriorAddresses(
      tbtcSystemStub.address,
      tbtcToken.address,
      tbtcDepositToken.address,
      feeRebateToken.address,
      vendingMachine.address,
    )

    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
    await testDeposit.setLotSize(new BN('100000000'))
  }

  await tbtcSystemStub.initialize(
    ecdsaKeepFactoryStub.address,
    depositFactory.address,
    testDeposit.address,
    tbtcToken.address,
    tbtcDepositToken.address,
    feeRebateToken.address,
    vendingMachine.address,
    1,
    1,
  )

  return {
    tbtcConstants,
    mockRelay,
    mockSatWeiPriceFeed,
    tbtcSystemStub,
    tbtcToken,
    tbtcDepositToken,
    feeRebateToken,
    testDeposit,
    depositUtils,
    ecdsaKeepFactoryStub,
    ecdsaKeepStub,
    depositFactory,
    deployed,
    redemptionScript,
    fundingScript,
    keepFactorySelectorStub,
  }
}

module.exports.deployAndLinkAll = deployAndLinkAll
