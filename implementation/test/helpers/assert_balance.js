export class AssertBalanceHelpers {
  constructor(tbtc) {
    this.tbtcInstance = tbtc
  }

  async tbtc(account, amount) {
    const balance = (await this.tbtcInstance.balanceOf(account)).toString()
    expect(balance).to.eq(amount)
  }

  async eth(account, amount) {
    const balance = await web3.eth.getBalance(account)
    expect(balance).to.eq(amount)
  }
}
