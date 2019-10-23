const BN = require('bn.js')
const chai = require('chai')
const expect = chai.expect
const bnChai = require('bn-chai')
chai.use(bnChai(BN))

export class AssertBalance {
  constructor(tbtc) {
    this.tbtcInstance = tbtc
  }

  async tbtc(account, amount) {
    const balance = await this.tbtcInstance.balanceOf(account)
    expect(balance).to.eq.BN(amount)
  }

  async eth(account, amount) {
    const balance = await web3.eth.getBalance(account)
    expect(balance).to.equal(amount)
  }
}
