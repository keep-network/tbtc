/*
This script is used to update client configuration file with latest deployed
contracts addresses.

Example:
BTC_NETWORK=testnet \
  CONFIG_FILE_PATH=~go/src/github.com/keep-network/tbtc/relay/config/config.toml \
  npx truffle exec scripts/lcl-client-config.js --network local
*/
const fs = require("fs")
const toml = require("toml")
const tomlify = require("tomlify-j0.4")
const btcNetworks = require('../migrations/btc-networks')

const Relay = artifacts.require(
    btcNetworks[process.env.BTC_NETWORK].contractName
)

module.exports = async function () {
  try {
    const configFilePath = process.env.CONFIG_FILE_PATH

    let relayAddress
    try {
      const relay = await Relay.deployed()
      relayAddress = relay.address
    } catch (err) {
      console.error("failed to get deployed contracts", err)
      process.exit(1)
    }

    try {
      const fileContent = toml.parse(
          fs.readFileSync(configFilePath, "utf8")
      )

      fileContent.ethereum.ContractAddresses.Relay = relayAddress

      const formattedConfigFile = tomlify.toToml(fileContent)

      fs.writeFileSync(configFilePath, formattedConfigFile, (err) => {
        if (err) throw err
      })

      console.log(`relay config written to ${configFilePath}`)
    } catch (err) {
      console.error("failed to update relay config", err)
      process.exit(1)
    }
  } catch (err) {
    console.error(err)
    process.exit(1)
  }
  process.exit(0)
}
