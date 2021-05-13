const ETHBTCPriceFeedMock = artifacts.require('ETHBTCPriceFeedMock')

async function run() {
  try {
    if (process.argv.length <= 4 || process.argv[4] === 'help') {
      console.log(
        'run `truffle exec multiply-local-eth-btc-price.js <factor>` to update the current price by a factor of <factor>',
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
        `Successfully updated the price by a factor of ${factor} from ${currentPrice} to ${newPrice}`,
      )
    }
    process.exit(0)
  } catch (ex) {
    console.error(ex)
    throw ex
  }
}

module.exports = run
