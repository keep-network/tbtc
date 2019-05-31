import expectThrow from './helpers/expectThrow';

const PriceOracleV1 = artifacts.require('PriceOracleV1')

contract('PriceOracleV1', function () {
    describe("#constructor", async () => {
        it('deploys', async () => {
            let instance = await PriceOracleV1.deployed()
        })
    })

    describe("Methods", async () => {
        let instance;

        beforeEach(async () => {
            instance = await PriceOracleV1.deployed()
        })

        describe("#getPrice", async () => {
            it('passes', () => {})
        })
    
        describe("#updatePrice", async () => {
            it('passes', () => {})            
        })
    })
})