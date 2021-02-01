/*
Here we assume that the passphrase for unlocking all the accounts on
some private testnet is the same.  This is intended for use with
truffle.  Example:

KEEP_ACCOUNT_PASSWORD=password \
  truffle exec ./scripts/unlock-eth-accounts.js \
  --network keep_dev
*/

const password = process.env.KEEP_ACCOUNT_PASSWORD || "password"

module.exports = async function() {
  const accounts = await web3.eth.getAccounts()

  console.log(`Total accounts: ${accounts.length}`)
  console.log(`---------------------------------`)

  for (let i = 0; i < accounts.length; i++) {
    const account = accounts[i]

    try {
      console.log(`\nUnlocking account: ${account}`)
      await web3.eth.personal.unlockAccount(account, password, 150000)
      console.log(`Account unlocked!`)
    } catch (error) {
      console.log(`\nAccount: ${account} not unlocked!`)
      console.error(error)
    }
    console.log(`\n---------------------------------`)
  }
  process.exit(0)
}
