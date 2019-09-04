// Snapshots are a feature of some EVM implementations (Ganache, but Geth is working on it) for improved dev UX.
// They allow us to snapshot the entire state of the chain, and restore it at a later point.
// https://github.com/trufflesuite/ganache-core/blob/master/README.md#custom-methods

const snapshotIdsStack = []


/**
 * Create a snapshot
 */
export async function createSnapshot() {
  return await new Promise((res, rej) => {
    web3.currentProvider.send({
      jsonrpc: '2.0',
      method: 'evm_snapshot',
      id: 12345,
    }, function(err, result) {
      if (err) rej(err)
      const snapshotId = result.result
      snapshotIdsStack.push(snapshotId)
      res()
    })
  })
}

/**
   * Restores the chain to a snapshot
   */
export async function restoreSnapshot() {
  const snapshotId = snapshotIdsStack.pop()
  return await new Promise((res, rej) => {
    web3.currentProvider.send({
      jsonrpc: '2.0',
      method: 'evm_revert',
      id: 12345,
      params: [snapshotId],
    }, function(err, result) {
      if (err) rej(err)
      else res()
    })
  })
}
