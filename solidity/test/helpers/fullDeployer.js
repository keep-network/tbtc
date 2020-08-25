const {contract} = require("@openzeppelin/test-environment")
const {deploySystem} = require("./utils.js")

const Deposit = contract.fromArtifact("Deposit")
const BytesLib = contract.fromArtifact("BytesLib")
const BTCUtils = contract.fromArtifact("BTCUtils")
const KeepFactorySelection = contract.fromArtifact("KeepFactorySelection")
const ValidateSPV = contract.fromArtifact("ValidateSPV")
const CheckBitcoinSigs = contract.fromArtifact("CheckBitcoinSigs")
const OutsourceDepositLogging = contract.fromArtifact("OutsourceDepositLogging")
const DepositStates = contract.fromArtifact("DepositStates")
const DepositUtils = contract.fromArtifact("DepositUtils")
const DepositFunding = contract.fromArtifact("DepositFunding")
const DepositRedemption = contract.fromArtifact("DepositRedemption")
const DepositLiquidation = contract.fromArtifact("DepositLiquidation")
const ECDSAKeepStub = contract.fromArtifact("ECDSAKeepStub")
const ECDSAKeepVendorStub = contract.fromArtifact("ECDSAKeepVendorStub")
const ECDSAKeepFactoryStub = contract.fromArtifact("ECDSAKeepFactoryStub")
const TestTBTCToken = contract.fromArtifact("TestTBTCToken")
const MockRelay = contract.fromArtifact("MockRelay")
const MockSatWeiPriceFeed = contract.fromArtifact("MockSatWeiPriceFeed")
const TBTCSystem = contract.fromArtifact("TBTCSystem")
const TBTCDepositToken = contract.fromArtifact("TestTBTCDepositToken")
const FeeRebateToken = contract.fromArtifact("TestFeeRebateToken")
const DepositFactory = contract.fromArtifact("DepositFactory")
const VendingMachine = contract.fromArtifact("VendingMachine")
const TBTCConstants = contract.fromArtifact("TBTCConstants")
const TestDeposit = contract.fromArtifact("TestDeposit")

const TEST_DEPOSIT_DEPLOY = [
  {name: "OutsourceDepositLogging", contract: OutsourceDepositLogging},
  {name: "MockRelay", contract: MockRelay},
  {name: "MockSatWeiPriceFeed", contract: MockSatWeiPriceFeed},
  {name: "KeepFactorySelection", contract: KeepFactorySelection},
  {
    name: "TBTCSystem",
    contract: TBTCSystem,
    constructorParams: ["MockSatWeiPriceFeed", "MockRelay"],
  },
  {
    name: "DepositFactory",
    contract: DepositFactory,
    constructorParams: ["TBTCSystem"],
  },
  {
    name: "VendingMachine",
    contract: VendingMachine,
    constructorParams: ["TBTCSystem"],
  },
  {name: "DepositStates", contract: DepositStates},
  {name: "TBTCConstants", contract: TBTCConstants}, // note the name
  {name: "DepositUtils", contract: DepositUtils},
  {name: "DepositRedemption", contract: DepositRedemption},
  {name: "DepositLiquidation", contract: DepositLiquidation},
  {name: "DepositFunding", contract: DepositFunding},
  {name: "TestDeposit", contract: TestDeposit},
  {name: "BytesLib", contract: BytesLib},
  {name: "BTCUtils", contract: BTCUtils},
  {name: "ValidateSPV", contract: ValidateSPV},
  {name: "CheckBitcoinSigs", contract: CheckBitcoinSigs},
  {
    name: "TBTCDepositToken",
    contract: TBTCDepositToken,
    constructorParams: ["DepositFactory"],
  },
  {
    name: "FeeRebateToken",
    contract: FeeRebateToken,
    constructorParams: ["VendingMachine"],
  },
  {
    name: "ECDSAKeepStub",
    contract: ECDSAKeepStub,
    constructorParams: ["VendingMachine", ""],
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
 *    - keepVendorStub
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

  const vendingMachine = deployed.VendingMachine

  const mockRelay = deployed.MockRelay
  await mockRelay.setPrevEpochDifficulty(1)

  const tbtcSystem = deployed.TBTCSystem
  const ecdsaKeepFactoryStub = await ECDSAKeepFactoryStub.new()
  await ecdsaKeepFactoryStub.setKeepAddress(deployed.ECDSAKeepStub.address)
  const ecdsaKeepVendorStub = await ECDSAKeepVendorStub.new(ecdsaKeepFactoryStub.address)

  const TestTBTCTokenInstance = await TestTBTCToken.new(vendingMachine.address)
  const testDeposit = deployed.TestDeposit

  const tbtcDepositToken = deployed.TBTCDepositToken
  const feeRebateToken = deployed.FeeRebateToken
  const ecdsaKeepStub = deployed.ECDSAKeepStub
  const depositFactory = deployed.DepositFactory
  const mockSatWeiPriceFeed = deployed.MockSatWeiPriceFeed
  await mockSatWeiPriceFeed.setPrice(100000000)

  await tbtcSystem.initialize(
    ecdsaKeepVendorStub.address,
    depositFactory.address,
    testDeposit.address,
    TestTBTCTokenInstance.address,
    tbtcDepositToken.address,
    feeRebateToken.address,
    vendingMachine.address,
    1,
    1,
  )

  const setupFee = await tbtcSystem.getNewDepositFeeEstimate()
  await ecdsaKeepStub.setBondAmount(setupFee)
  await ecdsaKeepStub.send(setupFee)

  return Object.assign({}, deployed, {
    Deposit,
    TestTBTCToken: TestTBTCTokenInstance,
    mockRelay,
    mockSatWeiPriceFeed,
    tbtcDepositToken,
    feeRebateToken,
    ecdsaKeepFactoryStub,
    ecdsaKeepStub,
  })
}

module.exports.deployAndLinkAll = deployAndLinkAll
