const fs = require('fs')
const path = require('path')
const Web3 = require('web3')

// These contracts are deployed on Ethereum mainnet.
const ETHUSDMedianizer = '0x64de91f5a373cd4c28de3600cb34c7c6ce410c85'
const BTCUSDMedianizer = '0xe0F30cb149fAADC7247E953746Be9BbBB6B5751f'

async function run() {
    try {
        /**
         * Reads the price from a MakerDAO price oracle.
         * Based off their Bash code detailed here:
         * https://github.com/makerdao/oracles-v2/tree/f795b66d9e928c2f4f6c701a0e6170080de95b7c#query-oracle-contracts
         * @param {string} address Address of price oracle
         */
        async function getMedianizerValue(address) {
            // We retrieve from contract storage, instead of calling IMedianizer.read(), 
            // to save on a conversion to bytes below.
            let value = await web3.eth.getStorageAt(address, 1)
            // The storage value is a uint256, but the first 128 bits 
            // are not related to the price value.
            let valueLower128 = value.slice(34, 66)
            return web3.utils.fromWei('0x' + valueLower128)
        }

        let eth = await getMedianizerValue(ETHUSDMedianizer)
        let btc = await getMedianizerValue(BTCUSDMedianizer)

        let config = {
            BTCUSD: btc,
            ETHUSD: eth,
        }
        console.log(JSON.stringify(config, null, 4))

        process.exit(0)

    } catch(ex) {
        console.error(ex)
        throw ex
    }
}

module.exports = run