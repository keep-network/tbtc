const { readFileSync } = require('fs')
const { join } = require('path')

const abis = {
    exchange: require('./contracts-vyper/abi/uniswap_exchange.json'),
    factory:  require('./contracts-vyper/abi/uniswap_factory.json')
}

const bytecode = {
    exchange:        readFileSync(require.resolve('./contracts-vyper/bytecode/exchange.txt'), 'utf-8').slice(2).trim(),
    factory:         readFileSync(require.resolve('./contracts-vyper/bytecode/factory.txt'),  'utf-8').slice(2).trim()
}

module.exports = {
    abis,
    bytecode
}