const { readFileSync } = require('fs')
const { join } = require('path')

// const abis = {
//     exchange: require('./contracts-vyper/abi/uniswap_exchange.json'),
//     factory:  require('./contracts-vyper/abi/uniswap_factory.json')
// }

// const bytecode = {
//     exchange: readFileSync(require.resolve('./contracts-vyper/bytecode/exchange.txt'), 'utf-8').slice(2).trim(),
//     factory:  readFileSync(require.resolve('./contracts-vyper/bytecode/factory.txt'),  'utf-8').slice(2).trim()
// }

function getDeployment(name) {
    return readFileSync(
        join(__dirname, `deployments/${name}`), 
        'utf-8'
    ).trim()
}

const deployments = {
    Exchange: getDeployment('Exchange'),
    Factory: getDeployment('Factory')
}

// const bytecode = {
//     exchange: readFileSync(require.resolve('./contracts-vyper/bytecode/exchange.txt'), 'utf-8').slice(2).trim(),
//     factory:  readFileSync(require.resolve('./contracts-vyper/bytecode/factory.txt'),  'utf-8').slice(2).trim()
// }
    


// const Web3 = require('web3')
// var contract = require("truffle-contract");

// function createTruffleContract(abi, bytecode) {
    
// }

// async function deploy(deployer, accounts) {
//     const web3 = new Web3(
//         'http://localhost:8545',
//         null, {}
//     );
//     // const accounts = await web3.eth.getAccounts()
    
//     // https://github.com/Uniswap/contracts-vyper/blob/master/tests/exchange/test_factory.py

//     // Deploy Factory and Exchange

//     // Call Factory.initializeFactory(exchangeAddress)

//     // And then Factory.createExchange

//     const defaultOpts = {
//         from: accounts[0],
//         gas: '5000000',
//         gasPrice: 1
//     }
//     web3.eth.defaultAccount = accounts[0];

//     // async function deployOne(deployer, abi) {
//     //     let x = await deployer.provider.send({
//     //         ...defaultOpts,
//     //         data: abi
//     //     })
//     //     console.log(x)
//     // }

//     // var exchange = contract({
//     //     abi: abis.exchange,
//     //     deployedBytecode: bytecode.exchange
//     // })
//     // exchange.setProvider(deployer.provider);
//     // await deployer.deploy(exchange)


//     const Exchange = new web3.eth.Contract(abis.exchange, null, { data: bytecode.exchange });
//     // const Factory = new web3.eth.Contract(abis.factory, null, { data: bytecode.factory });

//     let exchangeInstance = await Exchange
//         .deploy()
//         .send(defaultOpts)
    
//     // deployOne(exchangeInstance.encodeABI())
    
//     // console.log(123)
    
//     console.log(`Exchange ${exchangeInstance.options.address}`)

//     // let factoryInstance = await Factory
//     //     .deploy()
//     //     .send(defaultOpts);
    
//     // console.log(`Factory ${factoryInstance.options.address}`)

//     // await factoryInstance
//     //     .initializeFactory(Exchange.address)
//     //     .send(defaultOpts)

// }

module.exports = {
    deployments
}