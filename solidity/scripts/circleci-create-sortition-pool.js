/*
For each tbtc migration we should create a sortition pool for the new TBTCSystem application.
keep-ecdsa nodes can't function against a TBTCSystem instance if a sortition pool doesn't exist.
This script is intended to run via CircleCI after the migrate_contracts job is complete, 
however this script can be executed manually if needed.
*/

const Web3 = require("web3")
const HDWalletProvider = require("@truffle/hdwallet-provider")

// Ethereum host info.
const ethereumHost = process.env.ETH_HOST
const ethereumNetworkId = process.env.ETH_NETWORK_ID

// Contract owner info.
const contractOwnerAddress = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
const contractOwnerPrivateKey =
  process.env._CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY
const contractOwnerProvider = new HDWalletProvider(
  `${contractOwnerPrivateKey}`,
  `${ethereumHost}`,
)

// We override transactionConfirmationBlocks and transactionBlockTimeout because they're
// 25 and 50 blocks respectively at default.  The result of this on small private testnets
// is long wait times for scripts to execute.
const web3Options = {
  defaultBlock: "latest",
  defaultGas: 4712388,
  transactionBlockTimeout: 25,
  transactionConfirmationBlocks: 3,
  transactionPollingTimeout: 480,
}

// Setup web3 provider.  We use the keepContractOwner since it needs to sign the approveAndCall transaction.
const web3 = new Web3(contractOwnerProvider, null, web3Options)

// We lean on the keep-ecdsa package imported when installing the tbtc package file.
const BondedECDSAKeepFactoryJson = require("@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json")

/* 
Rather than use npm packages we directly access the compiled TBTCSystem contract that's generated
as part of the CircleCI migrate_contracts step.  This file is what's ultimately wrapped up in an
npm package, so it saves us a step.
*/
const TBTCSystemJson = require("./build/contracts/TBTCSystem.json")

// Boilerplate setup for access contract functions.
const BondedECDSAKeepFactoryJsonAbi = BondedECDSAKeepFactoryJson.abi
const BondedECDSAKeepFactoryJsonAddress =
  BondedECDSAKeepFactoryJsonJson.networks[ethereumNetworkId].address
const BondedECDSAKeepFactory = new web3.eth.Contract(
  BondedECDSAKeepFactoryJsonAbi,
  BondedECDSAKeepFactoryJsonAddress,
)
BondedECDSAKeepFactory.options.handleRevert = true

// We need the TBTCSystem address to pass to the
const TBTCSystemContractAddress =
  TBTCSystemJson.networks[ethereumNetworkId].address

async function createTBTCSystemSortitionPool() {
  try {
    console.log(
      `creating sortition pool for BondedECDSAKeepFactory: [${BondedECDSAKeepFactoryJsonAddress}]`,
    )
    // Create the pool for TBTCSystem.
    await BondedECDSAKeepFactory.methods
      .createSortitionPool(TBTCSystemContractAddress)
      .send({from: contractOwnerAddress})
      .on("transactionHash", hash => {
        console.log(`sortition pool creation tx: [${hash}]`)
      })
      .on("error", error => {
        console.error(`error creating sortition pool: [${error}]`)
      })
    // Log the new sortition pool address.
    sortitionPoolContractAddress = await BondedECDSAKeepFactory.methods
      .getSortitionPool(TBTCSystemContractAddress)
      .call()
    console.log(`sortition pool address: [${sortitionPoolContractAddress}]`)
    process.exit(0)
  } catch (error) {
    console.error(error)
    process.exit(1)
  }
}

createTBTCSystemSortitionPool().catch(error => {
  console.error(error)
  process.exit(1)
})
