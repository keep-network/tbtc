const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')

module.exports = async function() {
  const depositAddress = process.argv[4]

  let deposit
  let depositLog

  try {
    deposit = await Deposit.at(depositAddress)
    depositLog = await TBTCSystem.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  // getPublicKey calls tBTC to fetch signer's public key from the keep dedicated
  // for the deposit.
  async function getPublicKey() {
    console.log(`Call getPublicKey for deposit [${deposit.address}]`)
    const result = await deposit.retrieveSignerPubkey()
      .catch((err) => {
        console.error(`retrieveSignerPubkey failed: ${err}`)
        process.exit(1)
      })

    console.log('retrieveSignerPubkey transaction: ', result.tx)
  }

  async function logEvents(startBlockNumber) {
    const eventList = await depositLog.getPastEvents('RegisteredPubkey', {
      fromBlock: startBlockNumber,
      toBlock: 'latest',
    })

    const publicKeyX = eventList[0].returnValues._signingGroupPubkeyX
    const publicKeyY = eventList[0].returnValues._signingGroupPubkeyY

    console.log(`Registered public key:\nX: ${publicKeyX}\nY: ${publicKeyY}`)
  }

  const startBlockNumber = await web3.eth.getBlock('latest').number

  await getPublicKey()
  await logEvents(startBlockNumber)

  process.exit()
}
