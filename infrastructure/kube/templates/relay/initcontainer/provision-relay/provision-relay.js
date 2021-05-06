const fs = require('fs')
const HDWalletProvider = require("@truffle/hdwallet-provider")
const toml = require("toml")
const tomlify = require("tomlify-j0.4")
const Web3 = require('web3')

const hostChain = process.env.HOST_CHAIN || "ethereum"

// Host chain info.
const hostChainWsUrl = process.env.HOST_CHAIN_WS_URL
const hostChainRpcUrl = process.env.HOST_CHAIN_RPC_URL
const relayContractAddress = process.env.RELAY_CONTRACT_ADDRESS

// Relay operator info.
const operatorAddress = process.env.RELAY_ACCOUNT_ADDRESS
const operatorKeyFile = process.env.RELAY_ACCOUNT_KEY_FILE

// Contract owner info.
const contractOwnerAddress = process.env.CONTRACT_OWNER_ACCOUNT_ADDRESS
const contractOwnerKey = process.env.CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY

// Bitcoin chain info.
const bitcoinUrl = process.env.BITCOIN_URL
const bitcoinUsername = process.env.BITCOIN_USERNAME
const bitcoinPassword = process.env.BITCOIN_PASSWORD

// Metrics info.
const metricsPort = Number(process.env.METRICS_PORT || 0)
const metricsChainTick = Number(process.env.METRICS_CHAIN_TICK || 0)
const metricsNodeTick = Number(process.env.METRICS_NODE_TICK || 0)

/*
We override transactionConfirmationBlocks and transactionBlockTimeout because
they're 25 and 50 blocks respectively at default.  The result of this on small
private testnets is long wait times for scripts to execute.
*/
const web3Options = {
    defaultBlock: 'latest',
    defaultGas: 4712388,
    transactionBlockTimeout: 25,
    transactionConfirmationBlocks: 3,
    transactionPollingTimeout: 480
}

async function provisionRelayMaintainer() {
    console.log('###########  Provisioning relay! ###########')

    if (contractOwnerAddress) {
        console.log(`\n<<<<<<<<<<<< Funding operator account ${operatorAddress} >>>>>>>>>>>>`)
        await fundOperator(operatorAddress, contractOwnerAddress, '10')
    }

    console.log('\n<<<<<<<<<<<< Creating relay config file >>>>>>>>>>>>')
    await createRelayConfig()

    console.log("\n########### Relay provisioning complete! ###########")
}

async function fundOperator(operatorAddress, purse, requiredEtherBalance) {
    const contractOwnerProvider = new HDWalletProvider(contractOwnerKey, hostChainRpcUrl)
    const web3 = new Web3(contractOwnerProvider, null, web3Options)

    const requiredBalance = web3.utils.toBN(
        web3.utils.toWei(requiredEtherBalance, 'ether')
    )

    const currentBalance = web3.utils.toBN(
        await web3.eth.getBalance(operatorAddress)
    )

    if (currentBalance.gte(requiredBalance)) {
        console.log(
            `Operator address is already funded, ` +
            `current balance: ${web3.utils.fromWei(currentBalance)}`
        )
        return
    }

    const transferAmount = requiredBalance.sub(currentBalance)

    console.log(
        `Funding account ${operatorAddress} with ` +
        `${web3.utils.fromWei(transferAmount)} ether from purse ${purse}`
    )

    await web3.eth.sendTransaction({
        from: purse,
        to: operatorAddress,
        value: transferAmount
    })

    console.log(`Account ${operatorAddress} funded!`)
}

async function createRelayConfig() {
    const configFile = toml.parse(
        fs.readFileSync("/tmp/relay-config.toml", "utf8")
    )

    configFile[hostChain].URL = hostChainWsUrl
    configFile[hostChain].URLRPC = hostChainRpcUrl
    configFile[hostChain].account.KeyFile = operatorKeyFile
    configFile[hostChain].ContractAddresses.Relay = relayContractAddress

    configFile.bitcoin.URL = bitcoinUrl
    configFile.bitcoin.Username = bitcoinUsername
    configFile.bitcoin.Password = bitcoinPassword

    configFile.Metrics.Port = metricsPort
    configFile.Metrics.ChainMetricsTick = metricsChainTick
    configFile.Metrics.NodeMetricsTick = metricsNodeTick

    // tomlify.toToml() writes integer values as a float. Here we format the
    // default rendering to write the config file with integer values as needed.
    const formattedConfigFile = tomlify.toToml(configFile, {
        space: 2,
        replace: (key, value) => {
            // Find keys that match exactly `Port` or end with `MetricsTick`.
            const matcher = /(^Port|MetricsTick)$/

            return typeof key === "string" && key.match(matcher) ?
                value.toFixed(0) :
                false
        },
    })

    const configWritePath = "/mnt/relay/config/relay-config.toml"
    fs.writeFileSync(configWritePath, formattedConfigFile)
    console.log(`relay config written to ${configWritePath}`)
}

provisionRelayMaintainer()
    .catch(error => {
        console.error(error)
        process.exit(1)
    })
    .then(()=> {
        process.exit(0)
    })
