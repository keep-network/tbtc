export class UniswapHelpers {
  static async getDeadline(web3) {
    const block = await web3.eth.getBlock('latest')
    const DEADLINE_FROM_NOW = 300 // TX expires in 300 seconds (5 minutes)
    const deadline = block.timestamp + DEADLINE_FROM_NOW
    return deadline
  }
}
