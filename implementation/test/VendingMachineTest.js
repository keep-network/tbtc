
import expectThrow from './helpers/expectThrow'

const DepositUtils = artifacts.require('DepositUtils')
const TestTBTCConstants = artifacts.require('TestTBTCConstants')

const TestVendingMachine = artifacts.require('TestVendingMachine')
const TestToken = artifacts.require('TestToken')
const TestDepositOwnerToken = artifacts.require('TestDepositOwnerToken')

const BN = require('bn.js')
const utils = require('./utils')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

const TEST_DEPOSIT_DEPLOY = [
    { name: 'DepositUtils', contract: DepositUtils },
    { name: 'TBTCConstants', contract: TestTBTCConstants }, // note the name

    { name: 'TestVendingMachine', contract: TestVendingMachine }
]


contract('VendingMachine', (accounts) => {
    let vendingMachine
    let depositOwnerToken
    let tbtcToken

    before(async () => {
        // VendingMachine relies on linked libraries, hence we use deploySystem for consistency.
        let deployed = await utils.deploySystem(TEST_DEPOSIT_DEPLOY)

        tbtcToken = await TestToken.new(utils.address0)
        depositOwnerToken = await TestDepositOwnerToken.new()
        
        vendingMachine = deployed.TestVendingMachine
        await vendingMachine.setExteriorAddresses(tbtcToken.address, depositOwnerToken.address)

        console.log(
            depositOwnerToken.address,
            tbtcToken.address,
            vendingMachine.address
        )
    })

    describe('#dotToTbtc', async () => {
        it('converts DOT to TBTC', async () => {

        })
    })
    
})