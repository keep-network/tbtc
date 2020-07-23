const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {states} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, expectRevert} = require("@openzeppelin/test-helpers")
const {expect} = require("chai")

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0xs234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

describe("DepositLiquidation", async function() {
  let tbtcConstants
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let feeRebateToken
  let testDeposit
  let ecdsaKeepStub
  let beneficiary

  before(async () => {
    await createSnapshot()
  })

  after(async () => {
    await restoreSnapshot()
  })

  before(async () => {
    ;({
      tbtcConstants,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      ecdsaKeepStub,
    } = await deployAndLinkAll())

    beneficiary = accounts[4]

    const underThreshold = await tbtcSystemStub.getUndercollateralizedThresholdPercent()
    const severeThreshold = await tbtcSystemStub.getSeverelyUndercollateralizedThresholdPercent()

    await testDeposit.setUndercollateralizedThresholdPercent(
      new BN(underThreshold),
    )
    await testDeposit.setSeverelyUndercollateralizedThresholdPercent(
      new BN(severeThreshold),
    )
    await testDeposit.setSignerFeeDivisor(new BN("200"))

    await tbtcDepositToken.forceMint(
      beneficiary,
      web3.utils.toBN(testDeposit.address),
    )
    await feeRebateToken.forceMint(
      beneficiary,
      web3.utils.toBN(testDeposit.address),
    )
    await testDeposit.reset()
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("getCollateralizationPercentage", async () => {
    const ETHprice = new BN(190)
    const BTCPrice = new BN(7700)
    const satwei = BTCPrice.div(ETHprice).mul(new BN(10000000000))

    it("returns correct collatearlization value (> 100%)", async () => {
      const lotSize = await testDeposit.lotSizeSatoshis.call()

      await tbtcSystemStub.setOraclePrice(satwei)
      // expect 200% collateralization by setting bond = 2 * lotValue
      await ecdsaKeepStub.setBondAmount(satwei.mul(lotSize).mul(new BN(2)))

      const collateralization = await testDeposit.getCollateralizationPercentage()

      expect(collateralization).to.eq.BN(new BN(200))
    })

    it("returns correct collatearlization value (= 100%)", async () => {
      const lotSize = await testDeposit.lotSizeSatoshis.call()

      await tbtcSystemStub.setOraclePrice(satwei)
      // expect 100% collateralization by setting bond = lotValue.
      await ecdsaKeepStub.setBondAmount(satwei.mul(lotSize))

      const collateralization = await testDeposit.getCollateralizationPercentage()

      expect(collateralization).to.eq.BN(new BN(100))
    })

    it("returns correct collatearlization value (< 100%)", async () => {
      const lotSize = await testDeposit.lotSizeSatoshis.call()

      await tbtcSystemStub.setOraclePrice(satwei)
      // send 1/5 of value, expect 20% collateralization.
      await ecdsaKeepStub.setBondAmount(satwei.mul(lotSize).div(new BN(5)))

      const collateralization = await testDeposit.getCollateralizationPercentage()

      expect(collateralization).to.eq.BN(new BN(20))
    })

    it("returns correct collatearlization value (< 1%)", async () => {
      const lotSize = await testDeposit.lotSizeSatoshis.call()

      await tbtcSystemStub.setOraclePrice(satwei)
      // set less than 1% of bond, expect to receive a 0% collateralization (no decimals)
      await ecdsaKeepStub.setBondAmount(satwei.mul(lotSize).div(new BN(101)))

      const collateralization = await testDeposit.getCollateralizationPercentage()

      expect(collateralization).to.eq.BN(new BN(0))
    })
  })

  describe("purchaseSignerBondsAtAuction", async () => {
    let lotSize
    let buyer

    before(async () => {
      lotSize = await testDeposit.lotSizeTbtc.call()
      await testDeposit.setInitialCollateralizedPercent(new BN(150))
      buyer = accounts[1]
    })

    beforeEach(async () => {
      await testDeposit.setState(states.LIQUIDATION_IN_PROGRESS)
      for (let i = 0; i < 2; i++) {
        await tbtcToken.resetBalance(lotSize, {from: accounts[i]})
        await tbtcToken.resetAllowance(testDeposit.address, lotSize, {
          from: accounts[i],
        })
      }
    })

    it("sets state to liquidated, logs Liquidated, ", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.purchaseSignerBondsAtAuction({from: owner})

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATED)

      const eventList = await tbtcSystemStub.getPastEvents("Liquidated", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not in a liquidation auction", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        "No active auction",
      )
    })

    it("reverts if TBTC balance is insufficient", async () => {
      // burn 1 from buyer to make balance insufficient
      await tbtcToken.forceBurn(owner, 1)

      await expectRevert(
        testDeposit.purchaseSignerBondsAtAuction(),
        "Not enough TBTC to cover outstanding debt",
      )
    })

    it(`burns msg.sender's tokens`, async () => {
      const initialTokenBalance = await tbtcToken.balanceOf(buyer)

      await testDeposit.purchaseSignerBondsAtAuction({from: buyer})

      const finalTokenBalance = await tbtcToken.balanceOf(buyer)
      const tokenCheck = new BN(finalTokenBalance).add(new BN(lotSize))
      expect(tokenCheck, "tokens not burned correctly").to.eq.BN(
        initialTokenBalance,
      )
    })

    it("distributes reward to FRT holder", async () => {
      // Make sure Deposit has enough to cover beneficiary reward
      const beneficiaryReward = await testDeposit.signerFeeTbtc.call()
      await tbtcToken.forceMint(testDeposit.address, beneficiaryReward)

      const initialTokenBalance = await tbtcToken.balanceOf(beneficiary)

      await testDeposit.purchaseSignerBondsAtAuction({from: buyer})

      const finalTokenBalance = await tbtcToken.balanceOf(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add(
        new BN(beneficiaryReward).add(lotSize),
      )
      expect(
        finalTokenBalance,
        "tokens not returned to beneficiary correctly",
      ).to.eq.BN(tokenCheck)
    })

    it("awards withdrawable value to the buyer", async () => {
      const value = new BN("10000000000000000")
      const block = await web3.eth.getBlock("latest")
      const notifiedTime = block.timestamp

      await testDeposit.send(value, {from: accounts[8]})

      await testDeposit.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testDeposit.setState(states.LIQUIDATION_IN_PROGRESS)
      const auctionValue = await testDeposit.auctionValue.call()

      await testDeposit.purchaseSignerBondsAtAuction({from: buyer})

      // calculate the split of the un-purchased signer bond
      const split = value.sub(auctionValue).div(new BN(2))

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: buyer,
      })

      const depositBalance = await web3.eth.getBalance(testDeposit.address)

      expect(depositBalance).to.eq.BN(auctionValue.add(split))
      expect(withdrawable, "buyer should have a withdrawable balance").to.eq.BN(
        auctionValue,
      )
    })

    it("splits funds between liquidation triggerer and signers if not fraud", async () => {
      const liquidationInitiator = accounts[4]
      const block = await web3.eth.getBlock("latest")
      const notifiedTime = block.timestamp
      const value = 1000000000000
      const basePercentage = await testDeposit.getAuctionBasePercentage.call()

      await testDeposit.send(value, {from: accounts[8]})

      const initialSignerBalance = await web3.eth.getBalance(
        ecdsaKeepStub.address,
      )

      await testDeposit.setLiquidationInitiator(liquidationInitiator)
      await testDeposit.setLiquidationAndCourtesyInitated(notifiedTime, 0)

      const auctionValue = await testDeposit.auctionValue.call()
      // Buy auction immediately. No scaling takes place. Auction value is base
      // percentage of signer bond.
      await testDeposit.purchaseSignerBondsAtAuction({from: buyer})

      const finalSignerBalance = await web3.eth.getBalance(
        ecdsaKeepStub.address,
      )

      const signerBalanceDiff = new BN(finalSignerBalance).sub(
        new BN(initialSignerBalance),
      )

      const totalReward = (value * (100 - basePercentage)) / 100
      const split = totalReward / 2

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: liquidationInitiator,
      })
      const depositBalance = await web3.eth.getBalance(testDeposit.address)

      expect(depositBalance).to.eq.BN(auctionValue.add(new BN(split)))
      expect(new BN(split)).to.eq.BN(withdrawable)
      expect(new BN(split)).to.eq.BN(signerBalanceDiff)
    })

    it("transfers full ETH balance to liquidation triggerer if fraud", async () => {
      const block = await web3.eth.getBlock("latest")
      const notifiedTime = block.timestamp
      const liquidationInitiator = accounts[2]
      const value = 1000000000000
      const basePercentage = await testDeposit.getAuctionBasePercentage.call()

      await testDeposit.send(value, {from: accounts[8]})

      const initialSignerBalance = await web3.eth.getBalance(
        ecdsaKeepStub.address,
      )

      await testDeposit.setState(states.FRAUD_LIQUIDATION_IN_PROGRESS)
      await testDeposit.setLiquidationAndCourtesyInitated(notifiedTime, 0)
      await testDeposit.setLiquidationInitiator(liquidationInitiator)
      // Buy auction immediately. No scaling taken place. Auction value is base percentage of signer bond.
      await testDeposit.purchaseSignerBondsAtAuction({from: buyer})

      const finalSignerBalance = await web3.eth.getBalance(
        ecdsaKeepStub.address,
      )

      const signerBalanceDiff = new BN(finalSignerBalance).sub(
        new BN(initialSignerBalance),
      )

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: liquidationInitiator,
      })

      const totalReward = (value * (100 - basePercentage)) / 100
      const depositBalance = await web3.eth.getBalance(testDeposit.address)

      expect(depositBalance).to.eq.BN(new BN(value))
      expect(new BN(signerBalanceDiff)).to.eq.BN(0)
      expect(new BN(withdrawable)).to.eq.BN(totalReward)
    })
  })

  describe("notifyCourtesyCall", async () => {
    let oraclePrice
    let lotSize
    let lotValue
    let undercollateralizedPercent

    before(async () => {
      await tbtcSystemStub.setOraclePrice(new BN("1000000000000", 10))

      oraclePrice = await tbtcSystemStub.fetchBitcoinPrice.call()
      lotSize = await testDeposit.lotSizeSatoshis.call()
      lotValue = lotSize.mul(oraclePrice)

      undercollateralizedPercent = await testDeposit.getUndercollateralizedThresholdPercent.call()
    })

    beforeEach(async () => {
      await testDeposit.setState(states.ACTIVE)
      await ecdsaKeepStub.setBondAmount(0)
    })

    it("sets courtesy call state, sets the timestamp, and logs CourtesyCalled", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we subtract `1` to test collateralization less than undercollateralized
      // threshold (140%).
      const bondValue = undercollateralizedPercent
        .mul(lotValue)
        .div(new BN(100))
        .sub(new BN(1))
      await ecdsaKeepStub.setBondAmount(bondValue)

      await testDeposit.notifyCourtesyCall()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.COURTESY_CALL)

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      expect(liquidationTime[1]).not.to.eq.BN(0)

      const eventList = await tbtcSystemStub.getPastEvents("CourtesyCalled", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not in active state", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyCourtesyCall(),
        "Can only courtesy call from active state",
      )
    })

    it("reverts if sufficiently collateralized", async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization equal undercollateralized threshold (140%).
      const bondValue = undercollateralizedPercent
        .mul(lotValue)
        .div(new BN(100))
      await ecdsaKeepStub.setBondAmount(bondValue)
      await expectRevert(
        testDeposit.notifyCourtesyCall(),
        "Signers have sufficient collateral",
      )
    })
  })

  describe("exitCourtesyCall", async () => {
    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      const notifiedTime = blockTimestamp // not expired
      const fundedTime = blockTimestamp // not expired
      await ecdsaKeepStub.setBondAmount(new BN("1000000000000000000000000", 10))
      await tbtcSystemStub.setOraclePrice(new BN("1", 10))
      await testDeposit.setState(states.COURTESY_CALL)
      await testDeposit.setFundingInfo(
        "0x" + "00".repeat(8),
        fundedTime,
        "0x" + "00".repeat(36),
      )
      await testDeposit.setLiquidationAndCourtesyInitated(0, notifiedTime)
    })

    afterEach(async () => {
      await ecdsaKeepStub.setBondAmount(1000)
      await tbtcSystemStub.setOraclePrice(new BN("1000000000000", 10))
    })

    it("transitions to active, and logs ExitedCourtesyCall", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.exitCourtesyCall()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents(
        "ExitedCourtesyCall",
        {fromBlock: blockNumber, toBlock: "latest"},
      )
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not in courtesy call state", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.exitCourtesyCall(),
        "Not currently in courtesy call",
      )
    })

    it("reverts if the deposit is still undercollateralized", async () => {
      await tbtcSystemStub.setOraclePrice(new BN("1000000000000", 10))
      await ecdsaKeepStub.setBondAmount(0)

      await expectRevert(
        testDeposit.exitCourtesyCall(),
        "Deposit is still undercollateralized",
      )
    })
  })

  describe("notifyUndercollateralizedLiquidation", async () => {
    let oraclePrice
    let lotSize
    let lotValue
    let severelyUndercollateralizedPercent

    before(async () => {
      await tbtcSystemStub.setOraclePrice(new BN("1000000000000", 10))

      oraclePrice = await tbtcSystemStub.fetchBitcoinPrice.call()
      lotSize = await testDeposit.lotSizeSatoshis.call()
      lotValue = lotSize.mul(oraclePrice)

      severelyUndercollateralizedPercent = await tbtcSystemStub.getSeverelyUndercollateralizedThresholdPercent.call()
    })

    beforeEach(async () => {
      await testDeposit.setState(states.ACTIVE)
      await ecdsaKeepStub.setBondAmount(0)
      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    it("executes and moves state to LIQUIDATION_IN_PROGRESS", async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization less than severely undercollateralized
      // threshold (120%).
      const bondValue = severelyUndercollateralizedPercent
        .mul(lotValue)
        .div(new BN(100))
        .sub(new BN(1))
      await ecdsaKeepStub.setBondAmount(bondValue)

      await testDeposit.notifyUndercollateralizedLiquidation()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
      // TODO: Add validations or cover with `reverts if the deposit is not
      // severely undercollateralized` test case.
    })

    it("reverts if not in active or courtesy call", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyUndercollateralizedLiquidation(),
        "Deposit not in active or courtesy call",
      )
    })

    it("reverts if the deposit is not severely undercollateralized", async () => {
      // Bond value is calculated as:
      // `bondValue = collateralization * (lotSize * oraclePrice) / 100`
      // Here we test collateralization equal severely undercollateralized threshold (120%).
      const bondValue = severelyUndercollateralizedPercent
        .mul(lotValue)
        .div(new BN(100))
      await ecdsaKeepStub.setBondAmount(bondValue)

      await expectRevert(
        testDeposit.notifyUndercollateralizedLiquidation(),
        "Deposit has sufficient collateral",
      )
    })

    it("assert starts signer abort liquidation", async () => {
      await ecdsaKeepStub.send(1000000, {from: owner})
      await testDeposit.notifyUndercollateralizedLiquidation()

      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(new BN(bond), "Bond not seized as expected").to.eq.BN("0")

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      expect(
        liquidationTime[0],
        "liquidation timestamp not recorded",
      ).not.to.eq.BN(0)
    })
  })

  describe("notifyCourtesyTimeout", async () => {
    let courtesyTime
    let timer
    before(async () => {
      timer = await tbtcConstants.getCourtesyCallTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      courtesyTime = blockTimestamp - timer.toNumber() // has not expired
      await testDeposit.setState(states.COURTESY_CALL)
      await testDeposit.setLiquidationAndCourtesyInitated(0, courtesyTime)
      await ecdsaKeepStub.send(1000000, {from: owner})
    })

    it("executes and moves state to LIQUIDATION_IN_PROGRESS", async () => {
      await testDeposit.notifyCourtesyTimeout()
      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.LIQUIDATION_IN_PROGRESS)
    })

    it("reverts if not in a courtesy call period", async () => {
      await testDeposit.setState(states.START)
      await expectRevert(
        testDeposit.notifyCourtesyTimeout(),
        "Not in a courtesy call period",
      )
    })

    it("reverts if the period has not elapsed", async () => {
      await testDeposit.setLiquidationAndCourtesyInitated(0, courtesyTime * 5)
      await expectRevert(
        testDeposit.notifyCourtesyTimeout(),
        "Courtesy period has not elapsed",
      )
    })

    it("assert starts signer abort liquidation", async () => {
      await ecdsaKeepStub.send(1000000, {from: owner})
      await testDeposit.notifyCourtesyTimeout()

      const bond = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(new BN(bond), "Bond not seized as expected").to.eq.BN("0")

      const liquidationTime = await testDeposit.getLiquidationAndCourtesyInitiated.call()
      expect(
        liquidationTime[0],
        "liquidation timestamp not recorded",
      ).not.to.eq.BN(0)
    })
  })
})
