const DepositFactory = artifacts.require('./DepositFactory.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')
const TBTCToken = artifacts.require('./TBTCToken.sol')

module.exports = async function() {
  let tbtcSystem
  let tbtcToken

  try {
    tbtcToken = await TBTCToken.deployed()
    tbtcSystem = await TBTCSystem.deployed()
    depositFactory = await DepositFactory.deployed()
  } catch (err) {
    console.error(`initialization failed: ${err}`)
    process.exit(1)
  }

  async function createNewDeposit() {
    result = await depositFactory.createDeposit(
      tbtcSystem.address, // address _TBTCSystem
      tbtcToken.address, // address _TBTCToken
      5, // uint256 _m
      10, // uint256 _n
      { value: 100000 }
    ).catch((err) => {
      console.error(`call to factory failed: ${err}`)
      process.exit(1)
    })

    console.log('createNewDeposit transaction: ', result.tx)
  }

  async function logEvents(startBlockNumber) {
    const eventList = await depositFactory.getPastEvents('DepositCloneCreated', {
      fromBlock: startBlockNumber,
      toBlock: 'latest',
    })

    const depositAddress = eventList[0].returnValues.depositCloneAddress

    console.log('new deposit deployed: ', depositAddress)
  }

  const startBlockNumber = await web3.eth.getBlock('latest').number

  await createNewDeposit()
    .catch((err) => {
      console.error(`create new deposit failed: ${err}`)
      process.exit(1)
    })


  await logEvents(startBlockNumber)
    .catch((err) => {
      console.error(`log events failed: ${err}`)
      process.exit(1)
    })

  process.exit()
}
