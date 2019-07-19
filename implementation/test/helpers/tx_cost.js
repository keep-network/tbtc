const BN = require('bn.js')

// This is often useful in the cases when you're
// trying to assert some eth was refunded to msg.sender.
// Because ETH doesn't emit Transfer events like an ERC20 (unless you're usig W-ETH)
// we need to manually get balances and account for transaction costs.
export async function getTxCost(tx) {
  const gasUsed = tx.receipt.gasUsed

  const tx2 = await web3.eth.getTransaction(tx.tx)
  const gasPrice = tx2.gasPrice

  const txCostEth = new BN(gasUsed).mul(new BN(gasPrice))
  return txCostEth
}
