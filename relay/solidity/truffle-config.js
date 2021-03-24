/**
 * Use this file to configure your truffle project.
 *
 * More information about configuration can be found at:
 * truffleframework.com/docs/advanced/configuration
 *
 */

require("@babel/register")
require("@babel/polyfill")

module.exports = {
    /**
     * Networks define how you connect to your ethereum client and let you set
     * the defaults web3 uses to send transactions. You can ask a truffle
     * command to use a specific network from the command line, e.g
     *
     * $ truffle test --network <network-name>
     */
    networks: {
        local: {
            host: "127.0.0.1",
            port: 8546,
            network_id: 1101,
            websockets: true,
        },
    },
    // Configure your compilers
    compilers: {
        solc: {
            // Fetch exact version from solc-bin (default: truffle's version)
            version: "0.5.17",
        },
    },
}
