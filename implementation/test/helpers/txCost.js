const BN = require('bn.js')

/**
 * Gets the cost of a transaction in ETH.
 *
 * This is often useful in the cases when you're trying to assert some eth was refunded to `msg.sender`.
 * Because ETH doesn't emit Transfer events like an ERC20 (unless you're usig W-ETH), we need
 * to manually get balances and account for transaction costs.
 *
 * @param {*} txRes the result object from a web3 contract method call
 * @return {BN} the tx cost in ETH
 */
export async function getTxCost(txRes) {
  const gasUsed = txRes.receipt.gasUsed

  const tx = await web3.eth.getTransaction(txRes.tx)
  const gasPrice = tx.gasPrice

  const txCostEth = new BN(gasUsed).mul(new BN(gasPrice))
  return txCostEth
}
