const BytesLib = artifacts.require('BytesLib')
const BTCUtils = artifacts.require('BTCUtils')
const ValidateSPV = artifacts.require('ValidateSPV')
const CheckBitcoinSigs = artifacts.require('CheckBitcoinSigs')

const Deposit = artifacts.require('./Deposit');
const DepositUtils = artifacts.require('./DepositUtils');

const DepositLog = artifacts.require('./DepositLog');
const OutsourceDepositLogging = artifacts.require('./OutsourceDepositLogging');
const TBTCConstants = artifacts.require('./TBTCConstants')

const IBurnableERC20 = artifacts.require('./IBurnableERC20')
const IERC721 = artifacts.require('./IERC721');
const IKeep = artifacts.require('./IKeep')
const ITBTCSystem = artifacts.require('./ITBTCSystem')


module.exports = (deployer) => {
  deployer.then(async () => {

    await deployer.deploy(BytesLib)

    await deployer.link(BytesLib, [BTCUtils, ValidateSPV])
    await deployer.deploy(BTCUtils)

    await deployer.link(BTCUtils, ValidateSPV)
    await deployer.deploy(ValidateSPV)

    await deployer.deploy(TBTCConstants)
    await deployer.deploy(DepositLog)

    await deployer.link(BytesLib, CheckBitcoinSigs)
    await deployer.link(BTCUtils, CheckBitcoinSigs)
    await deployer.deploy(CheckBitcoinSigs);

    //await deployer.link(OutsourceDepositLogging, [TBTCConstants, DepositLog])
    await deployer.link(TBTCConstants, OutsourceDepositLogging)
    //await deployer.link(DepositLog, OutsourceDepositLogging)
    await deployer.deploy(OutsourceDepositLogging)

    await deployer.link(TBTCConstants, DepositUtils)
    await deployer.link(ValidateSPV, DepositUtils)
    await deployer.link(BTCUtils, DepositUtils)
    await deployer.link(BytesLib, DepositUtils)
    await deployer.deploy(DepositUtils);


    //await deployer.deploy(IBurnableERC20)
    //await deployer.deploy(IERC721)
    //await deployer.deploy(IKeep)
    //await deployer.deploy(ITBTCSystem)

    //await deployer.link(Deposit, [OutsourceDepositLogging, IBurnableERC20, IERC721, IKeep, ITBTCSystem])
   // await deployer.link(BTCUtils, BytesLib, CheckBitcoinSigs, TBTCConstants, ValidateSPV , Deposit)
    await deployer.link(BTCUtils, Deposit)
    await deployer.link(BytesLib, Deposit)
    await deployer.link(CheckBitcoinSigs, Deposit)
    await deployer.link(TBTCConstants, Deposit)
    await deployer.link(ValidateSPV, Deposit)
    await deployer.link(DepositUtils, Deposit)
    await deployer.deploy(Deposit)
  }
  )}
