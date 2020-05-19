const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {states} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")
const ECDSAKeepStub = contract.fromArtifact("ECDSAKeepStub")
const TBTCSystem = contract.fromArtifact("TBTCSystem")

// spare signature:
// signing with privkey '11' * 32
// const preimage = '0x' + '33'.repeat(32)
// const digest = '0xdeb0e38ced1e41de6f92e70e80c418d2d356afaaa99e26f5939dbc7d3ef4772a'
// const pubkey = '0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1'
// const v = 28
// const r = '0x9a40a074721355f427762f5e6d5cb16a0a9ada06011984e49fc81b3ce89cab6d'
// const s = '0x234e909713e74a9a49bf9484a69968dabcb1953bf091fa3e31d48531695cf293'

// real tx from mainnet bitcoin, interpreted as funding tx
// tx source: https://www.blockchain.com/btc/tx/7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f
// const tx = '0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800'
// const txid = '0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f';
// const txidLE = '0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c';
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
const _signerPubkeyX =
  "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e"
const _signerPubkeyY =
  "0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1"
const _merkleProof =
  "0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049"
const _expectedUTXOoutpoint =
  "0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c00000000"
// const _outputValue = 490029088;
const _outValueBytes = "0x2040351d00000000"

describe("DepositFunding", async function() {
  let tbtcConstants
  let mockRelay
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let testDeposit
  let ecdsaKeepStub
  let ecdsaKeepFactory

  let fundingProofTimerStart
  let beneficiary

  const funderBondAmount = new BN("10").pow(new BN("5"))
  const fullBtc = 100000000

  before(async () => {
    ;({
      tbtcConstants,
      mockRelay,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      testDeposit,
      ecdsaKeepStub,
      ecdsaKeepFactoryStub,
    } = await deployAndLinkAll())

    await tbtcSystemStub.setKeepAddress(ecdsaKeepStub.address)
    ecdsaKeepFactory = ecdsaKeepFactoryStub
    beneficiary = accounts[4]
    await tbtcDepositToken.forceMint(
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

  describe("createNewDeposit", async () => {
    it("runs and updates state and fires a created event", async () => {
      const expectedKeepAddress = ecdsaKeepStub.address
      const depositFee = await tbtcSystemStub.getNewDepositFeeEstimate()

      await ecdsaKeepStub.send(depositFee)
      await ecdsaKeepStub.setBondAmount(depositFee)

      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.createNewDeposit(
        tbtcSystemStub.address,
        tbtcToken.address,
        tbtcDepositToken.address,
        ZERO_ADDRESS,
        ZERO_ADDRESS,
        1, // m
        1,
        fullBtc,
        {value: funderBondAmount},
      )

      // state updates
      const depositState = await testDeposit.getState.call()
      expect(depositState, "state not as expected").to.eq.BN(
        states.AWAITING_SIGNER_SETUP,
      )

      const systemSignerFeeDivisor = await tbtcSystemStub.getSignerFeeDivisor()
      const signerFeeDivisor = await testDeposit.getSignerFeeDivisor.call()
      expect(signerFeeDivisor).to.eq.BN(systemSignerFeeDivisor)

      const keepAddress = await testDeposit.getKeepAddress.call()
      expect(keepAddress, "keepAddress not as expected").to.equal(
        expectedKeepAddress,
      )

      const stakeLockDuration = await ecdsaKeepFactory.stakeLockDuration.call()
      expect(
        stakeLockDuration,
        "stake lock duration not as expected",
      ).not.to.eq.BN(await tbtcConstants.getDepositTerm.call())

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signing group timestamp not as expected",
      ).not.to.eq.BN(0)

      // fired an event
      const eventList = await tbtcSystemStub.getPastEvents("Created", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList[0].returnValues._keepAddress).to.equal(
        expectedKeepAddress,
      )
    })

    it("reverts if bond is insufficient to cover a deposit creation fee refund", async () => {
      const depositFee = await tbtcSystemStub.getNewDepositFeeEstimate()

      await ecdsaKeepStub.send(depositFee - 1)
      await ecdsaKeepStub.setBondAmount(depositFee - 1)

      await expectRevert(
        testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          fullBtc,
        ),
        "Insufficient signer bonds to cover setup fee",
      )
    })

    it("reverts if not in the start state", async () => {
      await testDeposit.setState(states.REDEEMED)

      await expectRevert(
        testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          fullBtc,
        ),
        "Deposit setup already requested",
      )
    })

    it("fails if new deposits are disabled", async () => {
      await tbtcSystemStub.emergencyPauseNewDeposits()

      await expectRevert(
        testDeposit.createNewDeposit.call(
          tbtcSystemStub.address,
          tbtcToken.address,
          tbtcDepositToken.address,
          ZERO_ADDRESS,
          ZERO_ADDRESS,
          1, // m
          1,
          fullBtc,
        ),
        "Opening new deposits is currently disabled.",
      )
    })
  })
  describe("notifySignerSetupFailure", async () => {
    let timer
    let owner
    let openKeepFee

    before(async () => {
      ;({
        tbtcConstants,
        mockRelay,
        tbtcSystemStub,
        tbtcToken,
        tbtcDepositToken,
        testDeposit,
        ecdsaKeepStub,
      } = await deployAndLinkAll([], {TBTCSystemStub: TBTCSystem}))

      openKeepFee = await ecdsaKeepFactory.openKeepFeeEstimate.call()
      await testDeposit.setKeepSetupFee(openKeepFee)
      owner = accounts[1]
      await tbtcDepositToken.forceMint(
        owner,
        web3.utils.toBN(testDeposit.address),
      )
      timer = await tbtcConstants.getSigningGroupFormationTimeout.call()
    })

    beforeEach(async () => {
      await createSnapshot()

      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      const value = openKeepFee

      await ecdsaKeepStub.send(value)

      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_SIGNER_SETUP)

      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("updates state to setup failed, deconstes state, logs SetupFailed, and refunds TDT owner", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.notifySignerSetupFailure({from: owner})

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()

      const withdrawable = await testDeposit.getWithdrawAllowance.call({
        from: owner,
      })

      const depositBalance = await web3.eth.getBalance(testDeposit.address)
      expect(withdrawable).to.eq.BN(new BN(openKeepFee))
      expect(withdrawable).to.eq.BN(new BN(depositBalance))

      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt should be 0",
      ).to.eq.BN(0)

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "fundingProofTimerStart should be 0",
      ).to.eq.BN(0)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length, "Event list is the wrong length").to.equal(1)
    })

    it("reverts if not awaiting signer setup", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifySignerSetupFailure(),
        "Not awaiting setup",
      )
    })

    it("reverts if the timer has not yet elapsed", async () => {
      await testDeposit.setSigningGroupRequestedAt(fundingProofTimerStart * 5)

      await expectRevert(
        testDeposit.notifySignerSetupFailure(),
        "Signing group formation timeout not yet elapsed",
      )
    })
  })

  describe("retrieveSignerPubkey", async () => {
    const publicKey =
      "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"
    const pubkeyX =
      "0x4f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"
    const pubkeyY =
      "0x385b6b1b8ead809ca67454d9683fcf2ba03456d6fe2c4abe2b07f0fbdbb2f1c1"

    let ecdsaKeepStub

    before(async () => {
      ecdsaKeepStub = await ECDSAKeepStub.new()
    })

    beforeEach(async () => {
      await testDeposit.setState(states.AWAITING_SIGNER_SETUP)
      await testDeposit.setKeepAddress(ecdsaKeepStub.address)
      await ecdsaKeepStub.setPublicKey(publicKey)
    })

    it("updates the pubkey X and Y, changes state, and logs RegisteredPubkey", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.retrieveSignerPubkey()

      const signingGroupPublicKey = await testDeposit.getSigningGroupPublicKey.call()
      expect(signingGroupPublicKey[0]).to.equal(pubkeyX)
      expect(signingGroupPublicKey[1]).to.equal(pubkeyY)

      const eventList = await tbtcSystemStub.getPastEvents("RegisteredPubkey", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(
        eventList[0].returnValues._signingGroupPubkeyX,
        "Logged X is wrong",
      ).to.equal(pubkeyX)
      expect(
        eventList[0].returnValues._signingGroupPubkeyY,
        "Logged Y is wrong",
      ).to.equal(pubkeyY)
    })

    it("reverts if not awaiting signer setup", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "Not currently awaiting signer setup",
      )
    })

    it("reverts when public key is not 64-bytes long", async () => {
      await ecdsaKeepStub.setPublicKey("0x" + "00".repeat(63))

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "public key not set or not 64-bytes long",
      )
    })

    it("reverts if either half of the pubkey is 0", async () => {
      await ecdsaKeepStub.setPublicKey("0x" + "00".repeat(64))

      await expectRevert(
        testDeposit.retrieveSignerPubkey(),
        "Keep returned bad pubkey",
      )
    })
  })

  describe("notifyFundingTimeout", async () => {
    let timer

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    it("updates the state to failed setup, deconsts funding info, and logs SetupFailed", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.notifyFundingTimeout()

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.FAILED_SETUP)

      const eventList = await tbtcSystemStub.getPastEvents("SetupFailed", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting a funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.notifyFundingTimeout(),
        "Funding timeout has not started",
      )
    })

    it("reverts if the timeout has not elapsed", async () => {
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart * 5)

      await expectRevert(
        testDeposit.notifyFundingTimeout(),
        "Funding timeout has not elapsed",
      )
    })
  })

  describe("requestFunderAbort", async () => {
    let timer
    let owner

    before(async () => {
      timer = await tbtcConstants.getFundingTimeout.call()
      owner = accounts[1]
      await tbtcDepositToken.forceMint(
        owner,
        web3.utils.toBN(testDeposit.address),
      )
    })

    beforeEach(async () => {
      const block = await web3.eth.getBlock("latest")
      const blockTimestamp = block.timestamp
      fundingProofTimerStart = blockTimestamp - timer.toNumber() - 1

      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setFundingProofTimerStart(fundingProofTimerStart)
    })

    it("fails if the deposit has not failed setup", () => {
      expectRevert(
        testDeposit.requestFunderAbort("0x1234", {from: owner}),
        "The deposit has not failed funding",
      )
    })

    it("emits a FunderAbortRequested event", async () => {
      const blockNumber = await web3.eth.getBlockNumber()
      await testDeposit.notifyFundingTimeout()

      const outputScript = "0x012345"
      await testDeposit.requestFunderAbort(outputScript, {from: owner})

      const eventList = await tbtcSystemStub.getPastEvents(
        "FunderAbortRequested",
        {
          fromBlock: blockNumber,
          toBlock: "latest",
        },
      )
      expect(eventList.length).to.equal(1)
      expect(eventList[0].name == "FunderAbortRequested")
      expect(eventList[0].returnValues).to.contain({
        _depositContractAddress: testDeposit.address,
        _abortOutputScript: outputScript,
      })
    })
  })

  describe("provideBTCFundingProof", async () => {
    beforeEach(async () => {
      await mockRelay.setCurrentEpochDifficulty(currentDifficulty)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(_signerPubkeyX, _signerPubkeyY)
      await ecdsaKeepStub.send(1000000, {from: accounts[0]})
    })

    it("updates to active, stores UTXO info, deconstes funding info, logs Funded", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      const {
        receipt: {blockNumber: proofBlock},
      } = await testDeposit.provideBTCFundingProof(
        _version,
        _txInputVector,
        _txOutputVector,
        _txLocktime,
        _fundingOutputIndex,
        _merkleProof,
        _txIndexInBlock,
        _bitcoinHeaders,
      )
      const expectedFundedAt = (await web3.eth.getBlock(proofBlock)).timestamp

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      expect(UTXOInfo[0]).to.equal(_outValueBytes)
      expect(UTXOInfo[1]).to.eq.BN(new BN(expectedFundedAt))
      expect(UTXOInfo[2]).to.equal(_expectedUTXOoutpoint)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(
        signingGroupRequestedAt,
        "signingGroupRequestedAt not deconsted",
      ).to.not.equal(0)

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(
        fundingProofTimerStart,
        "fundingProofTimerStart not deconsted",
      ).to.not.equal(0)

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents("Funded", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.equal(1)
    })

    it("reverts if not awaiting funding proof", async () => {
      await testDeposit.setState(states.START)

      await expectRevert(
        testDeposit.provideBTCFundingProof(
          _version,
          _txInputVector,
          _txOutputVector,
          _txLocktime,
          _fundingOutputIndex,
          _merkleProof,
          _txIndexInBlock,
          _bitcoinHeaders,
        ),
        "Not awaiting funding",
      )
    })
  })

  describe("provideBTCFundingProof Legacy", async () => {
    const currentDifficultyLegacy = 6113
    const _versionLegacy = "0x01000000"
    const _txInputVectorLegacy =
      "0x01f3002663bbdfe1a0c02bafe6ad6bdc713902bbe678bdffd7e0a77972a4a6c4b2000000006b483045022100dfaa9366f2fd5a233e8029ae2f73e31dd75d85ef8970a74145c3128c0463cd9e0220739c9a57f1266d835b39f3f10d5df77d855ed4110ded0443584f4b7835c6167c012102b6ec2e4df2062066b685fd5cfe47f33c54b29a79a87ba20fcdb66f0b87170006ffffffff"
    const _txOutputVectorLegacy =
      "0x027818613b000000001976a914a29c9701293e838fc01ca1c13b010e42cf53242388aca08601000000000017a91414ea6747c79c9c7388eca6c029e6cc81b5e0951c87"
    const _txLocktimeLegacy = "0x00000000"
    const _fundingOutputIndexLegacy = 0
    const _merkleProofLegacy =
      "0x0efe91abb6a919b0f89c4c785460e2c2e50b57f3d9eb8da0885c26b02ddb50cb"
    const _txIndexInBlockLegacy = 1
    const _bitcoinHeadersLegacy =
      "0x03000000a3ecf52cf38b65970661436cab0b9813543fa949ac7b27db901ee402000000009873fc6b41300fb93957dc4fdec91f72f12e698bc5648a5c57623878a173002e2a9f9e5560b80a1baa34e3b5030000000d801f1aab80f52ca35e510df974c96901d1471c4eec0a377f570700000000004696408b572b53ed3df66479a4a3f69c5b6585daa3778d905483993dfad29360eda39e55ffff001d046ef73b0300000020cfaa259c8e375766f74be88525c7d1e41430d42a13844249a4f70300000000e9929e791caffe690fbb9406d06ef9376604bbf380dde84f64e98b5943e1a400039f9e5560b80a1b0e17707d03000000f818ca6e20d44ecb355e0897bd14eb50dcc35145f8d7554ebf2c040000000000ffe8e2da3ed2ae815e5948cef4addf32f7c599427a7d87e39d8c063c4dd30368cca39e55ffff001de8319459030000005b45b4375c101d7a95764cfda6e6f7b40246867cd54e89855b73200100000000fc074723d089d5338b65800707544f509bf111bc4d17a2b0353978319145fc5a459f9e5560b80a1b3fd1c3d2030000001cfc4d39d6ce1bebb6511ecb2039bb6fc26016eafb4d1e010be8020000000000e00725a1bca5411e7072970dccee9eeab802ca6c9cbcccd8a0aa39c40251bcbd02a49e55ffff001dd247f8d203000000cf778c9ea94a7f0eb51b06c29e65386f55efca5f3714c535b6e68703000000004bb0084076f57bce7a03f74da9d8784979ea0dbeefcce590ebd62d246b20bd6ab5a39e5560b80a1bbd5c1c1403000000e24a5e009fe23639a71f1024c6590abf3e2ffbe0ac92a43207ee0200000000007f740aa5ae4ea28ec5a3bd29151f3d053e10bb0489eb9a53b3f41604d05bfbccdaa39e5560b80a1bf27e379a03000000020b75077b8f84f296ee015471e7d99b0a57ff92acfb7739f62909000000000028cc108626e3c6f6fa12228c927b1ebec587b570bf085dcd24b63bb9b730c569eea39e5560b80a1b17793fc8"
    const _signerPubkeyXLegacy =
      "0x85134688b6be6f0fd93ce6a7f7904e9dba2b22bf3cc65c0feb5067b60cf5d821"
    const _signerPubkeyYLegacy =
      "0x76f80a7d522ea754db0e25ca43fdacfd1f313e4dc2765e2dfcc18fb3d63a66c4"

    const _expectedUTXOutpointLegacy =
      "0x0ee73932b031135c57e3d8f53db8ed5c97e6023a3e4980ea465f1aa2962d17b200000000"
    // const _outputValue = 996219000;
    const _outValueBytesLegacy = "0x7818613b00000000"

    beforeEach(async () => {
      await mockRelay.setCurrentEpochDifficulty(currentDifficultyLegacy)
      await testDeposit.setState(states.AWAITING_BTC_FUNDING_PROOF)
      await testDeposit.setSigningGroupPublicKey(
        _signerPubkeyXLegacy,
        _signerPubkeyYLegacy,
      )
      await ecdsaKeepStub.send(1000000, {from: accounts[0]})
    })

    it("updates to active, stores UTXO info, deconsts funding info, logs Funded", async () => {
      const blockNumber = await web3.eth.getBlockNumber()

      await testDeposit.provideBTCFundingProof(
        _versionLegacy,
        _txInputVectorLegacy,
        _txOutputVectorLegacy,
        _txLocktimeLegacy,
        _fundingOutputIndexLegacy,
        _merkleProofLegacy,
        _txIndexInBlockLegacy,
        _bitcoinHeadersLegacy,
      )

      const UTXOInfo = await testDeposit.getUTXOInfo.call()
      expect(UTXOInfo[0]).to.eql(_outValueBytesLegacy)
      expect(UTXOInfo[2]).to.eql(_expectedUTXOutpointLegacy)

      const signingGroupRequestedAt = await testDeposit.getSigningGroupRequestedAt.call()
      expect(signingGroupRequestedAt).to.eq.BN(
        0,
        "signingGroupRequestedAt not deconsted",
      )

      const fundingProofTimerStart = await testDeposit.getFundingProofTimerStart.call()
      expect(fundingProofTimerStart).to.eq.BN(
        0,
        "fundingProofTimerStart not deconsted",
      )

      const depositState = await testDeposit.getState.call()
      expect(depositState).to.eq.BN(states.ACTIVE)

      const eventList = await tbtcSystemStub.getPastEvents("Funded", {
        fromBlock: blockNumber,
        toBlock: "latest",
      })
      expect(eventList.length).to.eql(1)
    })
  })
})
