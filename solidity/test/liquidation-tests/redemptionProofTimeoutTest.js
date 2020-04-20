const { deployAndLinkAll } = require("./liquidation-test-utils/setup.js")
const { states, increaseTime } = require("../helpers/utils.js")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { accounts, web3, contract } = require("@openzeppelin/test-environment")
const [owner] = accounts
const { BN, expectRevert } = require("@openzeppelin/test-helpers")
const { expect } = require("chai")

const TestDeposit = contract.fromArtifact("TestDeposit")

describe("Redemption-proof timeout", async function () {
  const lotSize = new BN("100000000");
  const lotSizeTbtc = new BN("10000000000").mul(lotSize);
  const satwei = new BN("466666666666")
  const fee = new BN("5000000000000000")
  const depositInitiator = accounts[2]
  const auctionBuyer = accounts[3]
  const liqInitiator = accounts[4]
  let testDeposit
  let openKeepFee
  let minSufficientlyCollateralized

  before(async () => {
    const currentDifficulty = 6353030562983
    const _version = "0x01000000"
    const _txInputVector = `0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff`
    const _txOutputVector =
      "0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6"
    const _fundingOutputIndex = 0
    const _txLocktime = "0x4ec10800"
    const _txIndexInBlock = 129
    const _bitcoinHeaders =
      "0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e"
    const _merkleProof =
      "0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049"
    // const _outputValue = 490029088;
    // set up a new Deposit
    const publicKey =
      "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1"
    // redemption-values
    const redeemerOutputScript = "0x160014" + "33".repeat(20)
    const blockNumber = await web3.eth.getBlockNumber()



      ; ({
        mockRelay,
        depositFactory,
        ecdsaKeepFactoryStub,
        tbtcSystem,
        mockSatWeiPriceFeed,
        ecdsaKeepStub,
        tbtcConstants,
        tbtcToken
      } = await deployAndLinkAll())

    await ecdsaKeepFactoryStub.setKeepAddress(ecdsaKeepStub.address)

    await mockRelay.setCurrentEpochDifficulty(currentDifficulty)
    await mockRelay.setPrevEpochDifficulty(currentDifficulty)
    // 466666666666 wei per sat
    mockSatWeiPriceFeed.setPrice(satwei)

    const undercollateralized = await tbtcSystem.getUndercollateralizedThresholdPercent.call()

    minSufficientlyCollateralized = lotSize.mul(satwei).mul(undercollateralized).div(new BN(100))

    await ecdsaKeepStub.setBondAmount(minSufficientlyCollateralized)
    await ecdsaKeepStub.send(minSufficientlyCollateralized)

    openKeepFee = await tbtcSystem.createNewDepositFeeEstimate.call()

    await depositFactory.createDeposit(lotSize, { value: openKeepFee, from: depositInitiator })

    const eventList = await depositFactory.getPastEvents(
      "DepositCloneCreated",
      { fromBlock: blockNumber, toBlock: "latest" },
    )
    //  Deposit at clone address received from depositFactory event
    testDeposit = await TestDeposit.at(eventList[0].returnValues.depositCloneAddress)

    // retrieve pubkey and move to AWAITING_FUNDING_PROOF.
    // can call on ecdsaKeepStub directly since we set it up before
    await ecdsaKeepStub.setPublicKey(publicKey)

    await testDeposit.retrieveSignerPubkey()

    // provide the funding proof and move to ACTIVE state.
    await testDeposit.provideBTCFundingProof(
      _version,
      _txInputVector,
      _txOutputVector,
      _txLocktime,
      _fundingOutputIndex,
      _merkleProof,
      _txIndexInBlock,
      _bitcoinHeaders,
    )

    await tbtcToken.resetBalance(fee, { from: depositInitiator })
    await tbtcToken.resetAllowance(testDeposit.address, fee, {
      from: depositInitiator,
    })

    //  ACTIVE -> AWAITING_WITHDRAWAL_SIGNATURE
    await testDeposit.requestRedemption(
      "0x1111111100000000",
      redeemerOutputScript,
      { from: depositInitiator },// only TDT owner can redeem pre-term
    )

  })

  describe("Signature-timeout", async () => {

    it("unable to start liquidation with sufficient collateral", async () => {
      await expectRevert(
        testDeposit.notifyRedemptionProofTimeout(),
        "Not currently awaiting a redemption proof",
      )
      const depositState = await testDeposit.getCurrentState.call()
      expect(depositState).to.eq.BN(states.AWAITING_WITHDRAWAL_SIGNATURE)
    })

    it("liquidates correctly", async () => {
      const pubkeyX =
        "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
      const pubkeyY =
        "0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"
      const digest =
        "0x02d449a31fbb267c8f352e9968a79e3e5fc95c1bbeaa502fd6454ebde5a4bedc"
      const v = 27
      const r =
        "0xd7e83e8687ba8b555f553f22965c74e81fd08b619a7337c5c16e4b02873b537e"
      const s =
        "0x633bf745cdf7ae303ca8a6f41d71b2c3a21fcbd1aed9e7ffffa295c08918c1b3"
      await testDeposit.setSigningGroupPublicKey(pubkeyX, pubkeyY)
      await testDeposit.setRequestInfo(
        depositInitiator,
        "0x" + "11".repeat(20),
        0,
        0,
        digest,
      )

      //  AWAITING_WITHDRAWAL_SIGNATURE -> AWAITING_WITHDRAWAL_PROOF
      await testDeposit.provideRedemptionSignature(v, r, s)
      const requirement = await testDeposit.getOwnerRedemptionTbtcRequirement.call(depositInitiator)
      const timer = await tbtcConstants.getRedemptionProofTimeout.call()
      await increaseTime(timer.toNumber())

      await testDeposit.notifyRedemptionProofTimeout()

      const depositState = await testDeposit.getCurrentState.call()
      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: depositInitiator,
      })
     

      expect(depositState).to.eq.BN(states.LIQUIDATED)
      expect(requirement).to.eq.BN(fee)
      expect(withdrawable).to.eq.BN(minSufficientlyCollateralized)
    })
  })
})
