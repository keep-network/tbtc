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

const TestTBTCConstants = artifacts.require('TestTBTCConstants')
const TestDeposit = artifacts.require('TestDeposit')

TO_DEPLOY = [
  {name: 'BytesLib', contract: BytesLib},
  {name: 'BTCUtils', contract: BTCUtils},
  {name: 'ValidateSPV', contract: ValidateSPV},
  {name: 'CheckBitcoinSigs', contract: CheckBitcoinSigs},
  {name: 'TBTCConstants', contract: TestTBTCConstants},  // note the name
  {name: 'OutsourceDepositLogging', contract: OutsourceDepositLogging},
  {name: 'DepositLog', contract: DepositLog},
  {name: 'DepositStates', contract: DepositStates},
  {name: 'DepositUtils', contract: DepositUtils},
  {name: 'DepositFunding', contract: DepositFunding},
  {name: 'DepositRedemption', contract: DepositRedemption},
  {name: 'DepositLiquidation', contract: DepositLiquidation},
  {name: 'TestDeposit', contract: TestDeposit}]


async function deployTestDepositContract () {
  deployed = {}
  for (let i in TO_DEPLOY) {
    await TO_DEPLOY[i].contract.link(deployed)
    contract = await TO_DEPLOY[i].contract.new()
    deployed[TO_DEPLOY[i].name] = contract.address
  }
  return contract // TestDeposit is last
}

contract('Deposit', accounts => {

  describe('deployment', async () => {
    it('deploys', async () => {
      instance = await deployTestDepositContract()
    })
  })

})
