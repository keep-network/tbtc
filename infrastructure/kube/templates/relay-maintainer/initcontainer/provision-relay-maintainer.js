const fs = require('fs')
const Web3 = require('web3')
const HDWalletProvider = require("@truffle/hdwallet-provider")

const { URL } = require('url')

// ETH host info
const ethRPCUrl = process.env.ETH_RPC_URL
const ethNetworkId = process.env.ETH_NETWORK_ID

// Contract owner info
const contractOwnerAddress = process.env.CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
const purse = contractOwnerAddress

const contractOwnerPrivateKey = process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY
const contractOwnerProvider = new HDWalletProvider(
  contractOwnerPrivateKey,
  ethRPCUrl,
)

const operatorAddress = process.env.RELAY_MAINTAINER_ETH_ACCOUNT_ADDRESS
const operatorKey = process.env.RELAY_MAINTAINER_ETH_ACCOUNT_PRIVATE_KEY

/*
We override transactionConfirmationBlocks and transactionBlockTimeout because they're
25 and 50 blocks respectively at default.  The result of this on small private testnets
is long wait times for scripts to execute.
*/
const web3_options = {
  defaultBlock: 'latest',
  defaultGas: 4712388,
  transactionBlockTimeout: 25,
  transactionConfirmationBlocks: 3,
  transactionPollingTimeout: 480
}

const web3 = new Web3(contractOwnerProvider, null, web3_options)

async function provisionRelayMaintainer() {
  console.log('###########  Provisioning relay maintainer! ###########')

  console.log(`\n<<<<<<<<<<<< Funding Operator Account ${operatorAddress} >>>>>>>>>>>>`)
  await fundOperator(operatorAddress, purse, '10')

  console.log('\n<<<<<<<<<<<< Creating relay maintainer env config file >>>>>>>>>>>>')
  await createRelayMaintainerConfig()

  console.log("\n########### keep-ecdsa Provisioning Complete! ###########")
}

async function fundOperator(operatorAddress, purse, requiredEtherBalance) {
  let requiredBalance = web3.utils.toBN(web3.utils.toWei(requiredEtherBalance, 'ether'))

  const currentBalance = web3.utils.toBN(await web3.eth.getBalance(operatorAddress))
  if (currentBalance.gte(requiredBalance)) {
    console.log(`Operator address is already funded, current balance: ${web3.utils.fromWei(currentBalance)}`)
    return
  }

  const transferAmount = requiredBalance.sub(currentBalance)

  console.log(`Funding account ${operatorAddress} with ${web3.utils.fromWei(transferAmount)} ether from purse ${purse}`)
  await web3.eth.sendTransaction({ from: purse, to: operatorAddress, value: transferAmount })
  console.log(`Account ${operatorAddress} funded!`)
}

async function createRelayMaintainerConfig() {
  const envTemplate = fs.readFileSync('/tmp/env-template', 'utf8')

  const relayContractAddress = process.env.RELAY_CONTRACT_ADDRESS
  const ethURL = new URL(process.env.ETH_RPC_URL)

  const ethNetworkName = process.env.ETH_NETWORK_NAME
  // Do not include the protocol for BCOIN_HOST.
  // It's added on by the app.
  const bcoinHost = process.env.BCOIN_HOST
  const bcoinPort = process.env.BCOIN_PORT
  const bcoinApiKey = process.env.BCOIN_API_KEY
  const infuraKey = ""

  const finalEnv =
    envTemplate
      .replace(/^(SUMMA_RELAY_ETHER_HOST=).*$/m, `$1${ethURL.hostname}`)
      .replace(/^(SUMMA_RELAY_ETHER_PORT=).*$/m, `$1${ethURL.port}`)
      .replace(/^(SUMMA_RELAY_OPERATOR_KEY=).*$/m, `$1${operatorKey}`)
      .replace(/^(SUMMA_RELAY_ETH_NETWORK=).*$/m, `$1${ethNetworkName}`)
      .replace(/^(SUMMA_RELAY_ETH_CHAIN_ID=).*$/m, `$1${ethNetworkId}`)
      .replace(/^(SUMMA_RELAY_BCOIN_HOST=).*$/m, `$1${bcoinHost}`)
      .replace(/^(SUMMA_RELAY_BCOIN_PORT=).*$/m, `$1${bcoinPort}`)
      .replace(/^(SUMMA_RELAY_BCOIN_API_KEY=).*$/m, `$1${bcoinApiKey}`)
      .replace(/^(SUMMA_RELAY_INFURA_KEY=).*$/m, `$1${infuraKey}`)
      .replace(/^(SUMMA_RELAY_CONTRACT=).*$/m, `$1${relayContractAddress}`)

  fs.writeFileSync('/mnt/relay-maintainer/.env', finalEnv)
  console.log('relay maintainer .env file written to /mnt/relay-maintainer/.env')
}

provisionRelayMaintainer()
  .catch(error => {
    console.error(error)
    process.exit(1)
  })
  .then(()=> {
    process.exit(0)
  })
