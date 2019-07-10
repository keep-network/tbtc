// bitcoin-spv
const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

// logging
const OutsourceDepositLogging = artifacts.require('OutsourceDepositLogging')
const DepositLog = artifacts.require('DepositLog')

// deposit
const DepositStates = artifacts.require('DepositStates')
const DepositUtils = artifacts.require('DepositUtils')
const DepositFunding = artifacts.require('DepositFunding')
const DepositRedemption = artifacts.require('DepositRedemption')
const DepositLiquidation = artifacts.require('DepositLiquidation')
const Deposit = artifacts.require('Deposit')

// system
const TBTCConstants = artifacts.require('TBTCConstants')
const TBTCSystem = artifacts.require('TBTCSystem')

// keep
const KeepBridge = artifacts.require('KeepBridge')
const KeepRegistryAddress = '0xd04ed7D5C75cCC22DEafFD90A70c5BF932eC235e' // KeepRegistry contract address

const TBTC = artifacts.require('TBTC')

const all = [BytesLib, BTCUtils, ValidateSPV, TBTCConstants, CheckBitcoinSigs,
  OutsourceDepositLogging, DepositLog, DepositStates, DepositUtils,
  DepositFunding, DepositRedemption, DepositLiquidation, Deposit, TBTCSystem,
  KeepBridge]

const uniswap = require('../uniswap')

const UniswapDeployment = artifacts.require('UniswapDeployment')

async function deployUniswap(deployer, network, accounts) {
  async function deployFromBytecode(abi, bytecode) {
    const Contract = new web3.eth.Contract(abi)
    const instance = await Contract
      .deploy({
        data: `0x`+bytecode,
      })
      .send({
        from: accounts[0],
        gas: 4712388,
      })
    return instance
  }

  const exchange = await deployFromBytecode(uniswap.abis.exchange, uniswap.bytecode.exchange)
  const factory = await deployFromBytecode(uniswap.abis.factory, uniswap.bytecode.factory)

  // Required for Uniswap to clone and create factories
  await factory.methods.initializeFactory(exchange.options.address).send({ from: accounts[0] })

  // Save the deployed addresses to the UniswapDeployment contract
  await deployer.deploy(UniswapDeployment, factory.options.address, exchange.options.address)

  return factory.options.address
}

module.exports = (deployer, network, accounts) => {
  deployer.then(async () => {
    try {
      await deployUniswap(deployer, network, accounts)
    } catch (err) {
      throw new Error(`uniswap deployment failed: ${err}`)
    }

    // bitcoin-spv
    await deployer.deploy(BytesLib)
    await deployer.link(BytesLib, all)

    await deployer.deploy(BTCUtils)
    await deployer.link(BTCUtils, all)

    await deployer.deploy(ValidateSPV)
    await deployer.link(ValidateSPV, all)

    await deployer.deploy(CheckBitcoinSigs)
    await deployer.link(CheckBitcoinSigs, all)

    // constants
    await deployer.deploy(TBTCConstants)
    await deployer.link(TBTCConstants, all)

    // logging
    await deployer.deploy(OutsourceDepositLogging)
    await deployer.link(OutsourceDepositLogging, all)

    // deposit
    await deployer.deploy(DepositStates)
    await deployer.link(DepositStates, all)

    await deployer.deploy(DepositUtils)
    await deployer.link(DepositUtils, all)

    await deployer.deploy(DepositLiquidation)
    await deployer.link(DepositLiquidation, all)

    await deployer.deploy(DepositRedemption)
    await deployer.link(DepositRedemption, all)

    await deployer.deploy(DepositFunding)
    await deployer.link(DepositFunding, all)

    await deployer.deploy(Deposit)

    // system
    await deployer.deploy(TBTCSystem)
    await deployer.deploy(TBTC)

    // keep
    await deployer.deploy(KeepBridge).then((instance) => {
      instance.initialize(KeepRegistryAddress)
    })
  })
}
