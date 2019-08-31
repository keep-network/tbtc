const { readFileSync } = require('fs')

export const abis = {
  exchange: require('./contracts-vyper/abi/uniswap_exchange.json'),
  factory: require('./contracts-vyper/abi/uniswap_factory.json'),
}

export const bytecode = {
  exchange: readFileSync(require.resolve('./contracts-vyper/bytecode/exchange.txt'), 'utf-8').slice(2).trim(),
  factory: readFileSync(require.resolve('./contracts-vyper/bytecode/factory.txt'), 'utf-8').slice(2).trim(),
}

async function deployFromBytecode(web3, accounts, abi, bytecode) {
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

export async function deployUniswap(web3, accounts) {
  const exchange = await deployFromBytecode(web3, accounts, abis.exchange, bytecode.exchange)
  const factory = await deployFromBytecode(web3, accounts, abis.factory, bytecode.factory)

  // Required for Uniswap to clone and create factories
  await factory.methods.initializeFactory(exchange.options.address).send({ from: accounts[0] })

  return { exchange, factory }
}
