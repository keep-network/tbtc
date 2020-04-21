const {deployAndLinkAll} = require("./helpers/testDeployer.js")
const {
  HEADER_PROOFS,
  TX,
  LOW_WORK_HEADER,
  increaseTime,
} = require("./helpers/utils.js")
const {states} = require("./helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const {accounts, contract, web3} = require("@openzeppelin/test-environment")
const [owner] = accounts
const {BN, constants, expectRevert} = require("@openzeppelin/test-helpers")
const {ZERO_ADDRESS} = constants
const {expect} = require("chai")
const TestDepositUtils = contract.fromArtifact("TestDepositUtils")
const TestDepositUtilsSPV = contract.fromArtifact("TestDepositUtilsSPV")

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
const _concatenatedKeys =
  "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1"
const _merkleProof =
  "0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049"
const _expectedUTXOoutpoint =
  "0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c00000000"
// const _outputValue = 490029088;
const _outValueBytes = "0x2040351d00000000"

describe("DepositUtils", async function() {
  let beneficiary
  const funderBondAmount = new BN("10").pow(new BN("5"))
  const fullBtc = 100000000
  let tbtcConstants
  let mockRelay
  let tbtcSystemStub
  let tbtcToken
  let tbtcDepositToken
  let feeRebateToken
  let testDeposit
  let ecdsaKeepStub
  let depositUtils

  before(async () => {
    ;({
      tbtcConstants,
      mockRelay,
      tbtcSystemStub,
      tbtcToken,
      tbtcDepositToken,
      feeRebateToken,
      testDeposit,
      ecdsaKeepStub,
      depositUtils,
    } = await deployAndLinkAll([], {TestDeposit: TestDepositUtils}))

    beneficiary = accounts[2]

    feeRebateToken.forceMint(beneficiary, web3.utils.toBN(testDeposit.address))
    tbtcDepositToken.forceMint(
      beneficiary,
      web3.utils.toBN(testDeposit.address),
    )

    await testDeposit.createNewDeposit(
      tbtcSystemStub.address,
      tbtcToken.address,
      tbtcDepositToken.address,
      feeRebateToken.address,
      ZERO_ADDRESS,
      1, // m
      1, // n
      fullBtc,
      {value: funderBondAmount},
    )
  })

  beforeEach(async () => {
    await ecdsaKeepStub.reset()
    await testDeposit.setKeepAddress(ecdsaKeepStub.address)
  })

  describe("currentBlockDifficulty()", async () => {
    it("calls out to the system", async () => {
      const blockDifficulty = await testDeposit.currentBlockDifficulty.call()
      expect(blockDifficulty).to.eq.BN(1)

      await mockRelay.setCurrentEpochDifficulty(33)
      const newBlockDifficulty = await testDeposit.currentBlockDifficulty.call()
      expect(newBlockDifficulty).to.eq.BN(33)
    })
  })

  describe("previousBlockDifficulty()", async () => {
    it("calls out to the system", async () => {
      const blockDifficulty = await testDeposit.previousBlockDifficulty.call()
      expect(blockDifficulty).to.eq.BN(1)

      await mockRelay.setPrevEpochDifficulty(44)
      const newBlockDifficulty = await testDeposit.previousBlockDifficulty.call()
      expect(newBlockDifficulty).to.eq.BN(44)
    })
  })

  describe("evaluateProofDifficulty()", async () => {
    it("reverts on unknown difficulty", async () => {
      await mockRelay.setPrevEpochDifficulty(1)

      await expectRevert(
        testDeposit.evaluateProofDifficulty(HEADER_PROOFS[0]),
        "not at current or previous difficulty",
      )
    })

    it("evaluates a header proof with previous", async () => {
      await mockRelay.setPrevEpochDifficulty(5646403851534)
      await testDeposit.evaluateProofDifficulty(HEADER_PROOFS[0])
    })

    it("evaluates a header proof with current", async () => {
      await mockRelay.setCurrentEpochDifficulty(5646403851534)
      await testDeposit.evaluateProofDifficulty(HEADER_PROOFS[0])
    })

    it("reverts on low difficulty", async () => {
      await expectRevert(
        testDeposit.evaluateProofDifficulty(
          HEADER_PROOFS[0].slice(0, 160 * 4 + 2),
        ),
        "Insufficient accumulated difficulty in header chain",
      )
    })

    describe("reverts on a ValidateSPV errors", async () => {
      before(async () => {
        await mockRelay.setCurrentEpochDifficulty(5646403851534)
      })

      it("bad headers chain length work", async () => {
        // Cut one byte from a last header. We slice a hexadecimal representation
        // of concatenated chain headers bytes, followed by the `0x` prefix (add
        // `2`). Each header is expected to be `160` characters, we take
        // `4` headers and cut out a last byte from the last header (subtract
        // `2`).
        const badLengthChain = HEADER_PROOFS[0].slice(0, 2 + 160 * 4 - 2)

        await expectRevert(
          testDeposit.evaluateProofDifficulty(badLengthChain),
          "Invalid length of the headers chain",
        )
      })

      it("invalid headers chain", async () => {
        // Cut out one header from the headers chain. Take out second header
        // from the headers chain.
        const invalidChain =
          HEADER_PROOFS[0].slice(0, 2 + 160) +
          HEADER_PROOFS[0].slice(2 + 160 * 2, 2 + 160 * 4)

        await expectRevert(
          testDeposit.evaluateProofDifficulty(invalidChain),
          "Invalid headers chain",
        )
      })

      it("insufficient work in a header", async () => {
        await expectRevert(
          testDeposit.evaluateProofDifficulty(LOW_WORK_HEADER),
          "Insufficient work in a header",
        )
      })
    })
  })

  describe("checkProofFromTxId()", async () => {
    let testDeposit
    before(async () => {
      ;({mockRelay, testDeposit} = await deployAndLinkAll([], {
        TestDeposit: TestDepositUtilsSPV,
      }))

      await mockRelay.setCurrentEpochDifficulty(TX.difficulty)
    })

    it("does not error", async () => {
      await testDeposit.checkProofFromTxId.call(
        TX.tx_id_le,
        TX.proof,
        TX.index,
        HEADER_PROOFS.slice(-1)[0],
      )
    })

    it("fails with a broken proof", async () => {
      await expectRevert(
        testDeposit.checkProofFromTxId.call(
          TX.tx_id_le,
          TX.proof,
          0,
          HEADER_PROOFS.slice(-1)[0],
        ),
        "Tx merkle proof is not valid for provided header and txId",
      )
    })

    it("fails with bad difficulty", async () => {
      await mockRelay.setCurrentEpochDifficulty(1)
      await mockRelay.setPrevEpochDifficulty(1)

      await expectRevert(
        testDeposit.checkProofFromTxId.call(
          TX.tx_id_le,
          TX.proof,
          TX.index,
          HEADER_PROOFS.slice(-1)[0],
        ),
        "not at current or previous difficulty",
      )
    })
  })

  describe("findAndParseFundingOutput()", async () => {
    const _txOutputVector =
      "0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6"
    const _fundingOutputIndex = 0
    const _outValueBytes = "0x2040351d00000000"

    it("correctly returns valuebytes", async () => {
      await testDeposit.setPubKey(_signerPubkeyX, _signerPubkeyY)
      const valueBytes = await testDeposit.findAndParseFundingOutput.call(
        _txOutputVector,
        _fundingOutputIndex,
      )
      expect(
        _outValueBytes,
        "Got incorrect value bytes from funding output",
      ).to.equal(valueBytes)
    })

    it("fails with incorrect signer pubKey", async () => {
      await testDeposit.setPubKey(
        "0x" + "11".repeat(20),
        "0x" + "11".repeat(20),
      )

      await expectRevert(
        testDeposit.findAndParseFundingOutput.call(
          _txOutputVector,
          _fundingOutputIndex,
        ),
        "could not identify output funding the required public key hash",
      )
    })
  })

  describe("validateAndParseFundingSPVProof()", async () => {
    let tbtcSystemStub
    let tbtcToken
    let tbtcDepositToken
    let feeRebateToken
    let testDeposit

    before(async () => {
      ;({
        mockRelay,
        tbtcSystemStub,
        tbtcToken,
        tbtcDepositToken,
        feeRebateToken,
        testDeposit,
      } = await deployAndLinkAll([], {TestDeposit: TestDepositUtilsSPV}))

      beneficiary = accounts[2]

      tbtcDepositToken.forceMint(
        beneficiary,
        web3.utils.toBN(testDeposit.address),
      )
      feeRebateToken.forceMint(
        beneficiary,
        web3.utils.toBN(testDeposit.address),
      )

      await testDeposit.createNewDeposit(
        tbtcSystemStub.address,
        tbtcToken.address,
        tbtcDepositToken.address,
        feeRebateToken.address,
        ZERO_ADDRESS,
        1, // m
        1, // n
        fullBtc,
        {value: funderBondAmount},
      )

      await testDeposit.setPubKey(_signerPubkeyX, _signerPubkeyY)
      await mockRelay.setCurrentEpochDifficulty(currentDifficulty)
    })

    it("returns correct value and outpoint", async () => {
      const parseResults = await testDeposit.validateAndParseFundingSPVProof.call(
        _version,
        _txInputVector,
        _txOutputVector,
        _txLocktime,
        _fundingOutputIndex,
        _merkleProof,
        _txIndexInBlock,
        _bitcoinHeaders,
      )
      expect(parseResults[0]).to.equal(_outValueBytes)
      expect(parseResults[1]).to.equal(_expectedUTXOoutpoint)
    })

    it("fails with bad _txInputVector", async () => {
      await expectRevert(
        testDeposit.validateAndParseFundingSPVProof.call(
          _version,
          "0x" + "00".repeat(32),
          _txOutputVector,
          _txLocktime,
          _fundingOutputIndex,
          _merkleProof,
          _txIndexInBlock,
          _bitcoinHeaders,
        ),
        "invalid input vector provided",
      )
    })

    it("fails with bad _txOutputVector", async () => {
      await expectRevert(
        testDeposit.validateAndParseFundingSPVProof.call(
          _version,
          _txInputVector,
          "0x" + "00".repeat(32),
          _txLocktime,
          _fundingOutputIndex,
          _merkleProof,
          _txIndexInBlock,
          _bitcoinHeaders,
        ),
        "invalid output vector provided",
      )
    })

    it("fails with bad _merkleProof", async () => {
      await expectRevert(
        testDeposit.validateAndParseFundingSPVProof.call(
          _version,
          _txInputVector,
          _txOutputVector,
          _txLocktime,
          _fundingOutputIndex,
          "0x" + "00".repeat(32),
          _txIndexInBlock,
          _bitcoinHeaders,
        ),
        "Tx merkle proof is not valid for provided header and txId",
      )
    })

    it("fails with insufficient difficulty", async () => {
      const _badheaders = `0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6`

      await expectRevert(
        testDeposit.validateAndParseFundingSPVProof.call(
          _version,
          _txInputVector,
          _txOutputVector,
          _txLocktime,
          _fundingOutputIndex,
          _merkleProof,
          _txIndexInBlock,
          _badheaders,
        ),
        "Insufficient accumulated difficulty in header chain",
      )
    })
  })

  describe("getAuctionBasePercentage()", async () => {
    it("returns correct base percentage", async () => {
      const basePercentage = await testDeposit.getAuctionBasePercentage.call()

      const initialCollateralization = await testDeposit.getInitialCollateralizedPercent()

      // 10000 to avoid losing value to truncating in solidity.
      const expected = new BN(10000).div(initialCollateralization)
      expect(basePercentage).to.eq.BN(expected) // 66 for 150
    })
  })

  describe("auctionValue()", async () => {
    let duration
    let basePercentage
    before(async () => {
      duration = await tbtcConstants.getAuctionDuration.call()
      basePercentage = await testDeposit.getAuctionBasePercentage.call()
      auctionValue = new BN(100000000)
    })
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("returns base value if no time has elapsed", async () => {
      ecdsaKeepStub.pushFundsFromKeep(testDeposit.address, {
        value: auctionValue,
      })
      const block = await web3.eth.getBlock("latest")

      await testDeposit.setLiquidationAndCourtesyInitated(block.timestamp, 0)

      const value = await testDeposit.auctionValue.call()
      expect(value).to.eq.BN(auctionValue.mul(basePercentage).div(new BN(100)))
    })

    it("returns full value if auction Duration has elapsed ", async () => {
      ecdsaKeepStub.pushFundsFromKeep(testDeposit.address, {
        value: auctionValue,
      })
      const block = await web3.eth.getBlock("latest")

      await testDeposit.setLiquidationAndCourtesyInitated(block.timestamp, 0)
      await increaseTime(duration.toNumber())

      const value = await testDeposit.auctionValue.call()
      expect(value).to.eq.BN(auctionValue)
    })

    it("scales auction value correctly", async () => {
      ecdsaKeepStub.pushFundsFromKeep(testDeposit.address, {
        value: auctionValue,
      })
      const block = await web3.eth.getBlock("latest")

      await testDeposit.setLiquidationAndCourtesyInitated(block.timestamp, 0)

      const elapsedTime = duration.div(new BN(2))
      const elapsedPercent = new BN(100)
        .sub(basePercentage)
        .mul(elapsedTime)
        .div(duration)
      const percentage = basePercentage.add(elapsedPercent)

      // elapse half the auction time
      await increaseTime(elapsedTime.toNumber())

      const value = await testDeposit.auctionValue.call()
      expect(value).to.eq.BN(auctionValue.mul(percentage).div(new BN(100)))
    })
  })

  describe("signerFee()", async () => {
    it("returns a derived constant", async () => {
      const signerFee = await testDeposit.signerFee.call()
      expect(signerFee).to.eq.BN(5000000000000000)
    })
  })

  describe("determineCompressionPrefix()", async () => {
    it("selects 2 for even", async () => {
      const res = await depositUtils.determineCompressionPrefix.call(
        "0x" + "00".repeat(32),
      )
      expect(res).to.equal("0x02")
    })

    it("selects 3 for odd", async () => {
      const res = await depositUtils.determineCompressionPrefix.call(
        "0x" + "00".repeat(31) + "01",
      )
      expect(res).to.equal("0x03")
    })
  })

  describe("compressPubkey()", async () => {
    it("returns a 33 byte array with a prefix", async () => {
      const compressed = await depositUtils.compressPubkey.call(
        "0x" + "00".repeat(32),
        "0x" + "00".repeat(32),
      )
      expect(compressed).to.equal("0x02" + "00".repeat(32))
    })
  })

  describe("signerPubkey()", async () => {
    it("returns the concatenated signer X and Y coordinates", async () => {
      await testDeposit.setPubKey(_signerPubkeyX, _signerPubkeyY)
      const signerPubkey = await testDeposit.signerPubkey.call()
      expect(signerPubkey).to.equal(_concatenatedKeys)
    })

    it("returns base value for unset public key", async () => {
      const newTestUtilsInstance = await TestDepositUtils.new()
      const signerPubkey = await newTestUtilsInstance.signerPubkey.call()
      expect(signerPubkey).to.equal("0x" + "00".repeat(64))
    })
  })

  describe("signerPKH()", async () => {
    it("returns the concatenated signer X and Y coordinates", async () => {
      const expectedSignerPKH = "0xa99c23add58e3d0712278b2873c3c0bd21657115"
      await testDeposit.setPubKey(_signerPubkeyX, _signerPubkeyX)
      const signerPKH = await testDeposit.signerPKH.call()
      expect(signerPKH).to.equal(expectedSignerPKH)
    })
  })

  describe("utxoSize()", async () => {
    it("returns the state's utxoSizeBytes as an integer", async () => {
      const utxoSize = await testDeposit.utxoSize.call()
      expect(utxoSize).to.eq.BN(0)

      await testDeposit.setUTXOInfo("0x11223344", 1, "0x")
      const newUtxoSize = await testDeposit.utxoSize.call()
      expect(newUtxoSize).to.eq.BN(new BN("44332211", 16))
    })
  })

  describe("fetchBitcoinPrice()", async () => {
    it("calls out to the system", async () => {
      const oraclePrice = await testDeposit.fetchBitcoinPrice.call()
      expect(oraclePrice).to.eq.BN(1000000000000)

      await tbtcSystemStub.setOraclePrice(44)
      const newOraclePrice = await testDeposit.fetchBitcoinPrice.call()
      expect(newOraclePrice).to.eq.BN(44)
    })
  })

  describe("fetchBondAmount()", async () => {
    it("calls out to the keep system", async () => {
      const bondAmount = await testDeposit.fetchBondAmount.call()
      expect(bondAmount).to.eq.BN(10000)

      await ecdsaKeepStub.setBondAmount(44)
      const newBondAmount = await testDeposit.fetchBondAmount.call()
      expect(newBondAmount).to.eq.BN(44)
    })
  })

  describe("bytes8LEToUint()", async () => {
    it("interprets bytes as LE uints ", async () => {
      let res = await depositUtils.bytes8LEToUint.call("0x" + "00".repeat(8))
      expect(res).to.eq.BN(0)

      res = await depositUtils.bytes8LEToUint.call("0x11223344")
      expect(res).to.eq.BN(new BN("44332211", 16))

      res = await depositUtils.bytes8LEToUint.call("0x1100229933884477")
      expect(res).to.eq.BN(new BN("7744883399220011", 16))
    })
  })

  describe("wasDigestApprovedForSigning()", async () => {
    it("returns 0 when digest has not been approved", async () => {
      const digest = "0x" + "71".repeat(32)
      const expectedApprovalTime = new BN(0)

      const approvalTime = await testDeposit.wasDigestApprovedForSigning.call(
        digest,
      )

      expect(approvalTime).to.eq.BN(expectedApprovalTime)
    })

    it("returns approval time for approved digest", async () => {
      const digest = "0x" + "02".repeat(32)
      const expectedApprovalTime = new BN(100)

      await testDeposit.setDigestApprovedAtTime(digest, expectedApprovalTime)

      const approvalTime = await testDeposit.wasDigestApprovedForSigning.call(
        digest,
      )

      expect(approvalTime).to.eq.BN(expectedApprovalTime)
    })
  })

  describe("feeRebateTokenHolder()", async () => {
    it("calls out to the system", async () => {
      const res = await testDeposit.feeRebateTokenHolder.call()
      expect(res).to.equal(accounts[2])
    })
  })

  describe("redemptionTeardown()", async () => {
    it("deletes state", async () => {
      await testDeposit.setRequestInfo(
        "0x" + "11".repeat(20),
        "0x" + "11".repeat(20),
        5,
        6,
        "0x" + "33".repeat(32),
      )
      const requestInfo = await testDeposit.getRequestInfo.call()
      expect(requestInfo[4]).to.equal("0x" + "33".repeat(32))

      await testDeposit.redemptionTeardown()
      const newRequestInfo = await testDeposit.getRequestInfo.call()
      expect(newRequestInfo[4]).to.equal("0x" + "00".repeat(32))
    })
  })

  describe("seizeSignerBonds()", async () => {
    it("calls out to the keep system and returns the seized amount", async () => {
      const value = 5000
      await ecdsaKeepStub.send(value, {from: owner})

      const seized = await testDeposit.seizeSignerBonds.call()
      await testDeposit.seizeSignerBonds()

      expect(seized).to.eq.BN(value)
    })

    it("errors if no funds were seized", async () => {
      await expectRevert(
        testDeposit.seizeSignerBonds.call(),
        "No funds received, unexpected",
      )
    })
  })

  describe("enableWithdrawal()", async () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("correctly adds value to withdrawalAllowances", async () => {
      const value = new BN(10000)
      const beneficiary = accounts[1]

      const initialWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: beneficiary,
      })

      await testDeposit.enableWithdrawal(accounts[1], value)

      const finalWithdrawable = await testDeposit.getWithdrawAllowance.call({
        from: beneficiary,
      })

      expect(initialWithdrawable).to.eq.BN(new BN(0))
      expect(finalWithdrawable).to.eq.BN(value)
    })
  })

  describe("withdrawFunds()", async () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("correctly withdraws withdrawable balance at end-states", async () => {
      const endStates = [
        states.LIQUIDATED,
        states.REDEEMED,
        states.FAILED_SETUP,
      ]

      for (let i = 0; i < endStates.length; i++) {
        await createSnapshot()
        await testDeposit.setState(endStates[i])
        const value = new BN(web3.utils.toWei("0.1"))
        const beneficiary = accounts[1]
        const initialBalance = await web3.eth.getBalance(beneficiary)

        await ecdsaKeepStub.pushFundsFromKeep(testDeposit.address, {
          value: value,
        })
        await testDeposit.enableWithdrawal(beneficiary, value)
        const tx = await testDeposit.withdrawFunds({from: beneficiary})
        const finalBalance = await web3.eth.getBalance(beneficiary)

        const withdrawable = await testDeposit.getWithdrawAllowance.call({
          from: beneficiary,
        })
        const gasPrice = await web3.eth.getGasPrice()
        const gasUSed = tx.receipt.cumulativeGasUsed
        const totalTxCost = new BN(gasPrice).mul(new BN(gasUSed))

        expect(withdrawable).to.eq.BN(0)
        expect(new BN(finalBalance).add(totalTxCost)).to.eq.BN(
          new BN(initialBalance).add(value),
        )
        await restoreSnapshot()
      }
    })

    it("Reverts if not in end-state", async () => {
      const notEndStates = [
        states.START,
        states.AWAITING_SIGNER_SETUP,
        states.AWAITING_BTC_FUNDING_PROOF,
        states.ACTIVE,
        states.AWAITING_WITHDRAWAL_SIGNATURE,
        states.AWAITING_WITHDRAWAL_PROOF,
        states.COURTESY_CALL,
        states.FRAUD_LIQUIDATION_IN_PROGRESS,
        states.LIQUIDATION_IN_PROGRESS,
      ]

      for (let i = 0; i < notEndStates.length; i++) {
        await testDeposit.setState(notEndStates[i])

        await expectRevert(
          testDeposit.withdrawFunds(),
          "Contract not yet terminated",
        )
      }
    })
    it("Reverts if there is no available balance", async () => {
      await testDeposit.setState(states.LIQUIDATED)

      await expectRevert(testDeposit.withdrawFunds(), "Nothing to withdraw")
    })
  })

  describe("distributeFeeRebate()", async () => {
    it("checks that beneficiary is rewarded", async () => {
      // min an arbitrary reward value to the funding contract
      const reward = await testDeposit.signerFee.call()
      await tbtcToken.forceMint(testDeposit.address, reward)

      const initialTokenBalance = await tbtcToken.balanceOf(beneficiary)

      await testDeposit.distributeFeeRebate()

      const finalTokenBalance = await tbtcToken.balanceOf(beneficiary)
      const tokenCheck = new BN(initialTokenBalance).add(new BN(reward))
      expect(
        finalTokenBalance,
        "tokens not rewarded to beneficiary correctly",
      ).to.eq.BN(tokenCheck)
    })
  })

  describe("pushFundsToKeepGroup()", async () => {
    it("calls out to the keep contract", async () => {
      const value = 10000
      await ecdsaKeepStub.pushFundsFromKeep(testDeposit.address, {value: value})
      await testDeposit.pushFundsToKeepGroup(value)
      const keepBalance = await web3.eth.getBalance(ecdsaKeepStub.address)
      expect(keepBalance).to.eq.BN(new BN(value))
    })

    it("reverts if insufficient value", async () => {
      await expectRevert(
        testDeposit.pushFundsToKeepGroup.call(10000000000000),
        "Not enough funds to send",
      )
    })
  })

  describe("remainingTerm", async () => {
    const prevoutValueBytes = "0xffffffffffffffff"
    const outpoint = "0x" + "33".repeat(36)
    let depositTerm

    before(async () => {
      depositTerm = await tbtcConstants.getDepositTerm.call()
    })

    it("returns remaining term from current block", async () => {
      const block = await web3.eth.getBlock("latest")
      // Set Deposit.fundedAt to current block.
      await testDeposit.setUTXOInfo(
        prevoutValueBytes,
        block.timestamp,
        outpoint,
      )
      const remainingTerm = await testDeposit.remainingTerm.call()
      const finalBlock = await web3.eth.getBlock("latest")
      const expectedRemainder = new BN(block.timestamp)
        .add(depositTerm)
        .sub(new BN(finalBlock.timestamp))

      expect(remainingTerm).to.eq.BN(expectedRemainder)
    })

    it("returns 0 if deposit is at term", async () => {
      const block = await web3.eth.getBlock("latest")
      await testDeposit.setUTXOInfo(
        prevoutValueBytes,
        block.timestamp,
        outpoint,
      )
      const increaseTimeErrorMargin = 10
      await increaseTime(depositTerm.toNumber() + increaseTimeErrorMargin)

      const remainingTerm = await testDeposit.remainingTerm.call()
      expect(remainingTerm).to.eq.BN(0)
    })
  })
})
