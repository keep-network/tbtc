
import BN from 'bn.js'

export class UniswapHelpers {
  static async getDeadline(web3) {
    const block = await web3.eth.getBlock('latest')
    const DEADLINE_FROM_NOW = 300 // TX expires in 300 seconds (5 minutes)
    const deadline = block.timestamp + DEADLINE_FROM_NOW
    return deadline
  }

  /*
   * Add TBTC and ETH liquidity to the Uniswap exchange.
   */
  static async addLiquidity(account, uniswapExchange, tbtcToken, ethAmount, tbtcAmount) {
    // Mint tBTC, mock liquidity
    // supply the equivalent of 10 actors posting liquidity
    const supplyFactor = new BN(10)

    const tbtcSupply = new BN(tbtcAmount).mul(supplyFactor)
    await tbtcToken.forceMint(account, tbtcSupply, { from: account })
    await tbtcToken.approve(uniswapExchange.address, tbtcSupply, { from: account })

    // Uniswap requires a minimum of 1000000000 wei for the initial addLiquidity call.
    // See https://github.com/Uniswap/contracts-vyper/blob/c10c08d81d6114f694baa8bd32f555a40f6264da/contracts/uniswap_exchange.vy#L65
    // for more context on this undocumented constant.
    const UNISWAP_MINIMUM_INITIAL_LIQUIDITY_WEI = new BN('1000000000')
    const ethSupply = new BN(ethAmount).add(UNISWAP_MINIMUM_INITIAL_LIQUIDITY_WEI).mul(supplyFactor)

    const deadline = await this.getDeadline(web3)
    /* eslint-disable no-multi-spaces */
    await uniswapExchange.addLiquidity(
      '0',           // min_liquidity
      tbtcSupply,    // max_tokens
      deadline,      // deadline
      { from: account, value: ethSupply }
    )
    /* eslint-enable no-multi-spaces */
  }
}
