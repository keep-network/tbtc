/**
 * Snapshots are a feature of some EVM implementations ([1]) for improved dev UX.
 * They allow us to snapshot the entire state of the chain, and restore it at a
 * later point.
 *
 * We can use snapshots to create state for an `it` test, perform some transactions,
 * and then revert that state after. However, **Ganache snapshots can only be restored ONCE.**
 *
 * What this means in practice, is that snapshots must be created before EVERY test in
 * a `beforeEach` block, and thereafter reverted in `afterEach`.
 *
 * Example usage:
 * ```js
 * before(async () => {
 *  await createSnapshot()
 * })
 *
 * after(async () => {
 *  await restoreSnapshot()
 * })
 *
 * beforeEach(async () => {
 *  await createSnapshot()
 * })
 *
 * afterEach(async () => {
 *  await restoreSnapshot()
 * })
 *
 * it('#something', async () => { ... })
 * ```
 *
 * This pattern is adopted from the 0x codebase - see their BlockchainLifecycle [2], which
 * implements it for Ganache and Geth. And in their tests [3], running it in both before
 * and beforeEach blocks.
 *
 * [1]: https://github.com/trufflesuite/ganache-core/blob/master/README.md#custom-methods
 * [2]: https://sourcegraph.com/github.com/0xProject/0x-monorepo@ec92cea5982375fa2fa7ba8445b5e8af589b75bd/-/blob/packages/dev-utils/src/blockchain_lifecycle.ts#L21
 * [3]: https://sourcegraph.com/github.com/0xProject/0x-monorepo@ec92cea/-/blob/contracts/asset-proxy/test/authorizable.ts#L27
 */


const snapshotIdsStack = []

export async function createSnapshot() {
  const snapshotId = await new Promise((res, rej) => {
    web3.currentProvider.send({
      jsonrpc: '2.0',
      method: 'evm_snapshot',
      id: 12345,
    }, function(err, result) {
      if (err) rej(err)
      else res(result.result)
    })
  })

  snapshotIdsStack.push(snapshotId)
}

export async function restoreSnapshot() {
  const snapshotId = snapshotIdsStack.pop()

  try {
    await new Promise((res, rej) => {
      web3.currentProvider.send({
        jsonrpc: '2.0',
        method: 'evm_revert',
        id: 12345,
        params: [snapshotId],
      }, function(err, result) {
        if (err) rej(err)
        else res(result.result)
      })
    })
  } catch (ex) {
    throw new Error(`Snapshot with id #${snapshotId} failed to revert`)
  }
}
