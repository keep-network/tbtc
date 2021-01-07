// Configuration for addresses of externally deployed smart contracts
// prettier-ignore
const BondedECDSAKeepFactoryAddress = "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"

// Medianized price feeds.
// These are deployed and maintained by Maker.
// See: https://github.com/makerdao/oracles-v2#live-mainnet-oracles
const ETHBTCMedianizer = "0xABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDE"

// Ropsten contracts.
const RopstenETHBTCPriceFeed = "0xe9046e086137d2c0ffe60035391f6d7b4ec16733"
const RopstenTestnetRelay = "0xcF3a6246879aab9eb7beC0d936743208EA51d0Ed"

module.exports = {
  BondedECDSAKeepFactoryAddress,
  ETHBTCMedianizer,
  RopstenETHBTCPriceFeed,
  RopstenTestnetRelay,
}
