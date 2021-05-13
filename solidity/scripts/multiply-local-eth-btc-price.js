const ETHBTCPriceFeedMock = artifacts.require('ETHBTCPriceFeedMock')

async function run() {
  try {
    if (process.argv.length <= 4 || process.argv[4] === 'help') {
      console.log(
        `
        |truffle exec multiply-local-eth-btc-price.js <factor>
        |  Updates the local mock ETHBTC price by the specified factor.
        |
        |  <factor>
        |    A factor to multiply the current ETHBTC price by. ETHBTC price
        |    is in wei per satoshi, so a factor of 2 means the price of one satoshi
        |    has doubled in terms of wei, i.e. it takes twice as many wei to buy
        |    one satoshi. In terms of collateralization, this means a deposit that
        |    was 100% collateralized is now 50% collateralized.
        |`.replaceAll(/[ \t]+\|/g, '').trim()
      )
    } else if (isNaN(parseFloat(process.argv[4]))) {
      console.log('<factor> must be a number.')
    } else {
      const ethBtcPriceFeedMock = await ETHBTCPriceFeedMock.deployed()
      const currentPrice = parseInt(await ethBtcPriceFeedMock.read())
      // the arguments are [node, truffle, exec, <script name>, arguments...]
      const factor = parseFloat(process.argv[4])
      const newPrice = Math.round(currentPrice * factor)
      await ethBtcPriceFeedMock.setValue(newPrice.toString())
      console.log(
        `Successfully updated the ETHBTC price by a factor of ${factor} from ${currentPrice} wei/satoshi to ${newPrice} wei/satoshi`,
      )
    }
    process.exit(0)
  } catch (ex) {
    console.error(ex)
    throw ex
  }
}

module.exports = run
