const OnDemandSPV = artifacts.require('@summa-tx/relay-sol/OnDemandSPV');
const TestnetRelay = artifacts.require('@summa-tx/relay-sol/TestnetRelay');
const btcNetworks = require('./btc-networks');

module.exports = async (deployer, network) => {
  const btcNetwork = btcNetworks[process.env.BTC_NETWORK];

  if (!btcNetwork) {
    console.error(
        "BTC_NETWORK environment variable must point " +
        "either to `mainnet` or `testnet` value"
    );
    process.exit(1)
  }

  const contract = {OnDemandSPV, TestnetRelay}[btcNetwork.contractName]
  const { genesis, height, epochStart, getFirstID } = btcNetwork;
  const firstID = getFirstID(network)

  console.log("\nDeployment info\n")
  console.log(`genesis: ${genesis}`)
  console.log(`height: ${height}`)
  console.log(`epochStart: ${epochStart}`)
  console.log(`firstID: ${firstID}`)

  await deployer.deploy(contract, genesis, height, epochStart, firstID);
};