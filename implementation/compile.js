const path = require('path')
const solc = require('solc')
const fs = require('fs-extra')

// Path to compiled json contracts
const buildPath = path.resolve(__dirname, 'build')
fs.removeSync(buildPath)
fs.ensureDirSync(buildPath)

// Paths to solidity contracts
const BTCUtilsPath = path.resolve(__dirname, 'bitcoin-spv', 'contracts', 'BTCUtils.sol')
const ValidateSPVPath = path.resolve(__dirname, 'bitcoin-spv', 'contracts', 'ValidateSPV.sol')
const BytesPath = path.resolve(__dirname, 'bitcoin-spv', 'contracts', 'BytesLib.sol')
const SafeMathPath = path.resolve(__dirname, 'bitcoin-spv', 'contracts', 'SafeMath.sol')
const SigCheckPath = path.resolve(__dirname, 'bitcoin-spv', 'contracts', 'SigCheck.sol')
const IBurnableERC20Path = path.resolve(__dirname, 'contracts', 'interfaces', 'IBurnableERC20.sol')
const IERC721Path = path.resolve(__dirname, 'contracts', 'interfaces', 'IERC721.sol')
const ITBTCSystemPath = path.resolve(__dirname, 'contracts', 'interfaces', 'ITBTCSystem.sol')
const IKeepPath = path.resolve(__dirname, 'contracts', 'interfaces', 'IKeep.sol')
const TBTCConstantsPath = path.resolve(__dirname, 'contracts', 'TBTCConstants.sol')
const DepositLogPath = path.resolve(__dirname, 'contracts', 'DepositLog.sol')
const OutsourceDepositLoggingPath = path.resolve(__dirname, 'contracts', 'OutsourceDepositLogging.sol')
const DepositPath = path.resolve(__dirname, 'contracts', 'Deposit.sol')

let input = {
    'BTCUtils.sol': fs.readFileSync(BTCUtilsPath, 'utf8'),
    'ValidateSPV.sol': fs.readFileSync(ValidateSPVPath, 'utf8'),
    'BytesLib.sol': fs.readFileSync(BytesPath, 'utf8'),
    'SafeMath.sol': fs.readFileSync(SafeMathPath, 'utf8'),
    'SigCheck.sol': fs.readFileSync(SigCheckPath, 'utf8'),
    'IBurnableERC20.sol': fs.readFileSync(IBurnableERC20Path, 'utf8'),
    'IERC721.sol': fs.readFileSync(IERC721Path, 'utf8'),
    'ITBTCSystem.sol': fs.readFileSync(ITBTCSystemPath, 'utf8'),
    'IKeep.sol': fs.readFileSync(IKeepPath, 'utf8'),
    'TBTCConstants.sol': fs.readFileSync(TBTCConstantsPath, 'utf8'),
    'DepositLog.sol': fs.readFileSync(DepositLogPath, 'utf8'),
    'OutsourceDepositLogging.sol': fs.readFileSync(OutsourceDepositLoggingPath, 'utf8'),
    'Deposit.sol': fs.readFileSync(DepositPath, 'utf8')
}

const output = solc.compile({sources: input}, 1)

// log errors
if (output.errors) {
  for (let e in output.errors) {
    console.log(output.errors[e])
  }
}

// Save compiled contracts to json files
for (let contract in output.contracts) {
    contract_name = contract.split(':')
    fs.outputJsonSync(
        path.resolve(buildPath, contract_name[1] + '.json'),
        output.contracts[contract]
    )
}
