const BN = require('bn.js');
const truffleConfig = require('../truffle-config');

const ID_SPACE_SIZE = new BN('2', 10).pow(new BN('32', 10));

module.exports = {
    mainnet: {
        genesis: '0x00006020dd02d03c03dbc1f41312a6940e89919ce67fbf99a20307000000000000000000260d70e7ae07c80db07fbf29d09ec1a86d4f788e58098189a6f9021236572a7dd99eb15e397a11178294a823',
        height: 629070,
        epochStart: '0x459ec50d4ea62a89da04eb1ef3e352ec740bca50e8a808000000000000000000',
        getFirstID: (ethNetwork) => ID_SPACE_SIZE.muln(truffleConfig.networks[ethNetwork].network_id),
        contractName: "OnDemandSPV"
    },

    testnet: {
        genesis: '0x0000c0205d1103efc13e6647977e3d65f253c3e762451e9ca9b920517d000000000000008442a07bcde3292a888277ea6337ba5bbdfa808ae01535846f19d843144c8f60478bb15e7b41011a88be5d36',
        height: 1723030,
        epochStart: '0xe2657f702faa9470815005305c45b4be2271c22ade1348e6fe00000000000000',
        getFirstID: (ethNetwork) => ID_SPACE_SIZE.muln(truffleConfig.networks[ethNetwork].network_id + 0x800000),
        contractName: "TestnetRelay"
    }
}