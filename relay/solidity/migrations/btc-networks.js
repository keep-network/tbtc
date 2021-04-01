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
        genesis: '0x0000453faa4dcc559d47ed24be33b5539638b797e2dbd7527e61d8c13f00000000000000c4cb592bba9489c40ae5b12efa0f39e4d98dd6d58b794b65f6c50e450ee9957081b96160c0ff3f19015ae4ae',
        height: 1968421,
        epochStart: '0xb2e552027cbe48d7f85d5e1b7f0a49227173ba960f420e5d2500000000000000',
        getFirstID: (ethNetwork) => ID_SPACE_SIZE.muln(truffleConfig.networks[ethNetwork].network_id + 0x800000),
        contractName: "TestnetRelay"
    }
}