const headerChains = require("./headerchains.json")
const tx = require("./tx.json")
const createHash = require("create-hash")
const { web3 } = require("@openzeppelin/test-environment")

const ozHelpers = require("@openzeppelin/test-helpers")
const { BN } = ozHelpers
const { expect } = require("chai")

// Header with insufficient work. It's used for negative scenario tests when we
// want to validate invalid header which hash (work) doesn't meet requirement of
// the target.
const lowWorkHeader =
  "0xbbbbbbbb7777777777777777777777777777777777777777777777777777777777777777e0e333d0fd648162d344c1a760a319f2184ab2dce1335353f36da2eea155f97fcccccccc7cd93117e85f0000bbbbbbbbcbee0f1f713bdfca4aa550474f7f252581268935ef8948f18d48ec0a2b4800008888888888888888888888888888888888888888888888888888888888888888cccccccc7cd9311701440000bbbbbbbbfe6c72f9b42e11c339a9cbe1185b2e16b74acce90c8316f4a5c8a6c0a10f00008888888888888888888888888888888888888888888888888888888888888888dccccccc7cd9311730340000"

const states = {
  START: new BN(0),
  AWAITING_SIGNER_SETUP: new BN(1),
  AWAITING_BTC_FUNDING_PROOF: new BN(2),
  FAILED_SETUP: new BN(3),
  ACTIVE: new BN(4),
  AWAITING_WITHDRAWAL_SIGNATURE: new BN(5),
  AWAITING_WITHDRAWAL_PROOF: new BN(6),
  REDEEMED: new BN(7),
  COURTESY_CALL: new BN(8),
  FRAUD_LIQUIDATION_IN_PROGRESS: new BN(9),
  LIQUIDATION_IN_PROGRESS: new BN(10),
  LIQUIDATED: new BN(11),
}

function hash160(hexString) {
  const buffer = Buffer.from(hexString, "hex")
  const t = createHash("sha256").update(buffer).digest()
  const u = createHash("rmd160").update(t).digest()
  return "0x" + u.toString("hex")
}

function chainToProofBytes(chain) {
  let byteString = "0x"
  for (let header = 6; header < chain.length; header++) {
    byteString += chain[header].hex
  }
  return byteString
}

// Links and deploys contracts. contract Name and address are used along with
// an optional constructorParams parameter. The constructorParams will be
// passed as constructor parameter for the deployed contract. If the
// constructorParams contain the name of a previously deployed contract, that
// contract's address is passed as the constructor parameter instead.
async function deploySystem(deployList) {
  const deployed = {} // name: contract object
  const linkable = {} // name: linkable address
  const names = []
  const addresses = []
  for (let i = 0; i < deployList.length; ++i) {
    await deployList[i].contract.detectNetwork()
    for (let q = 0; q < names.length; q++) {
      await deployList[i].contract.link(names[q], addresses[q])
    }
    let contract
    if (deployList[i].constructorParams == undefined) {
      contract = await deployList[i].contract.new()
    } else {
      const resolvedConstructorParams = deployList[i].constructorParams.map(
        (param) => linkable[param] || param
      )
      contract = await deployList[i].contract.new(...resolvedConstructorParams)
    }
    linkable[deployList[i].name] = contract.address
    deployed[deployList[i].name] = contract
    names.push(deployList[i].name)
    addresses.push(contract.address)
  }
  return deployed
}

function increaseTime(duration) {
  const id = Date.now()

  // Convert from BN if necessary. Throws if the BN is too big.
  const durationNumber = BN.isBN(duration) ? duration.toNumber() : duration
  return new Promise((resolve, reject) => {
    web3.currentProvider.send(
      {
        jsonrpc: "2.0",
        method: "evm_increaseTime",
        params: [durationNumber],
        id: id,
      },
      (err1) => {
        if (err1) return reject(err1)

        web3.currentProvider.send(
          {
            jsonrpc: "2.0",
            method: "evm_mine",
            id: id + 1,
          },
          (err2, res) => {
            return err2 ? reject(err2) : resolve(res)
          }
        )
      }
    )
  })
}

function getCurrentTime() {
  const id = Date.now()

  return new Promise((resolve, reject) => {
    web3.currentProvider.send(
      {
        jsonrpc: "2.0",
        method: "eth_getBlockByNumber",
        params: ["latest", true],
        id: id,
      },
      (err1, { result: tx }) => {
        if (err1) return reject(err1)

        resolve(parseInt(tx.timestamp))
      }
    )
  })
}

/**
 * Uses the ABIs of all contracts in the `contractContainer` to resolve any
 * events they may have emitted into the given `receipt`'s logs. Typically
 * Truffle only resolves the events on the calling contract; this function
 * resolves all of the ones that can be resolved.
 *
 * @param {TruffleReceipt} receipt The receipt of a contract function call
 *        submission.
 * @param {ContractContainer} contractContainer An object that contains
 *        properties that are TruffleContracts. Not all properties in the
 *        container need be contracts, nor do all contracts need to have events
 *        in the receipt.
 *
 * @return {TruffleReceipt} The receipt, with its `logs` property updated to
 *         include all resolved logs.
 */
function resolveAllLogs(receipt, contractContainer) {
  const contracts = Object.entries(contractContainer)
    .map(([, value]) => value)
    .filter((_) => _.contract && _.address)

  const { resolved: resolvedLogs } = contracts.reduce(
    ({ raw, resolved }, contract) => {
      const events = contract.contract._jsonInterface.filter(
        (_) => _.type === "event"
      )
      const contractLogs = raw.filter((_) => _.address == contract.address)

      const decoded = contractLogs.map((log) => {
        const event = events.find((_) => log.topics.includes(_.signature))
        const decoded = web3.eth.abi.decodeLog(
          event.inputs,
          log.data,
          log.topics.slice(1)
        )

        return Object.assign({}, log, {
          event: event.name,
          args: decoded,
        })
      })

      return {
        raw: raw.filter((_) => _.address != contract.address),
        resolved: resolved.concat(decoded),
      }
    },
    { raw: receipt.rawLogs, resolved: [] }
  )

  return Object.assign({}, receipt, {
    logs: resolvedLogs,
  })
}

/**
 * Wrapper for OpenZeppelin's expectEvent helper that deals with array-of-BN
 * parameters.
 *
 * @param {TxReceipt} receipt The receipt to check for the specified event.
 * @param {string} eventName The name of the event to look for.
 * @param {object} [parameters] The parameters to look for in the event; unlike
 *         OpenZeppelin's default version, parameters may have an array-of-BN
 *         value and they will be properly validated.
 * @param {object} [customMessage] The message to report if a matching event is
 *        not found.
 */
function expectEvent(receipt, eventName, parameters, customMessage) {
  parameters = parameters || {}

  const matchingLogs = receipt.logs.filter((_) => _.event == eventName)
  expect(matchingLogs).not.to.have.lengthOf(
    0,
    customMessage ||
      `Found no events named ${eventName} in ${JSON.stringify(receipt.logs)}`
  )
  expect(
    matchingLogs,
    customMessage ||
      `Could not find parameter match for\n${JSON.stringify(
        parameters
      )}\nin\n` +
        `${matchingLogs.map((_) => JSON.stringify(_.args)).join("\n")}.`
  ).to.satisfy((matchingLogs) => {
    // Check array-of-BN parameters in nested fashion, directly.
    return matchingLogs.some((log) =>
      Object.entries(log.args).every(([argName, argValue]) => {
        if (!parameters.hasOwnProperty(argName)) {
          return true
        } else if (web3.utils.isBN(parameters[argName])) {
          return parameters[argName].toString() == web3.utils.isBN(argValue)
            ? argValue.toString()
            : argValue
        } else if (web3.utils.isBN(parameters[argName][0])) {
          return parameters[argName].every((bn, i) => {
            return bn.toString() == web3.utils.isBN(argValue[i])
              ? argValue[i].toString()
              : argValue[i]
          })
        } else {
          return parameters[argName] == argValue
        }
      })
    )
  })
}

/**
 * Verifies that the given `receipt` has *no* event of with the given name and
 * matching the given parameters. The assertion errors if an event with this
 * name and matching the specified parameters exists in the receipt. Any
 * parameters that are not explicitly enumerated match any value; that is, if
 * an event with the given name is found that matches the parameters passed,
 * and that event has an additional parameter, this assertion will still fail,
 * despite the additional parameter not being explicitly mentioned.
 *
 * @param {TxReceipt} receipt The receipt to check for the specified event.
 * @param {string} eventName The name of the event to look for.
 * @param {object} [parameters] The parameters to look for in the event; unlike
 *         OpenZeppelin's default version, parameters may have an array-of-BN
 *         value and they will be properly validated.
 */
function expectNoEvent(receipt, eventName, parameters) {
  parameters = parameters || {}
  const matchingLogs = receipt.logs
    .filter((_) => _.event == eventName)
    .filter((log) => {
      for (const [parameterName, parameterValue] of Object.entries(
        parameters
      )) {
        if (!log.args.hasOwnProperty(parameterName)) {
          // If the parameter doesn't exist, this doesn't match.
          return false
        }

        let matching
        const eventValue = log.args[parameterName]
        if (web3.utils.isBN(parameterValue[0])) {
          matching =
            parameterValue.filter((bn, i) => bn.eq(eventValue[i])).length ==
            parameterValue.length
        } else if (web3.utils.isBN(parameterValue)) {
          matching = parameterValue.eq(eventValue)
        } else {
          matching = parameterValue == eventValue
        }

        if (!matching) {
          return false
        }
      }

      return true
    })

  expect(matchingLogs, `Unexpected event: ${JSON.stringify(matchingLogs[0])}`)
    .to.be.empty
}

// real tx from mainnet bitcoin, interpreted as funding tx
// tx source: https://www.blockchain.com/btc/tx/7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f
const fundingTx = {
  tx:
    "0x01000000000101913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab602473044022046c3c852a2042ee01ffd7d8d252f297ccc67ae2aa1fac4170f50e8a90af5398802201585ffbbed6e812fb60c025d2f82ae115774279369610b0c76165b6c7132f2810121020c67643b5c862a1aa1afe0a77a28e51a21b08396a0acae69965b22d2a403fd1c4ec10800",
  txid: "0x7c48181cb5c030655eea651c5e9aa808983f646465cbe9d01c227d99cfbc405f",
  txidLE: "0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c",
  difficulty: 6353030562983,
  version: "0x01000000",
  txInputVector:
    "0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff",
  txOutputVector:
    "0x012040351d0000000016001486e7303082a6a21d5837176bc808bf4828371ab6",
  fundingOutputIndex: 0,
  txLocktime: "0x4ec10800",
  txIndexInBlock: 129,
  bitcoinHeaders:
    "0x00e0ff3fd877ad23af1d0d3e0eb6a700d85b692975dacd36e47b1b00000000000000000095ba61df5961d7fa0a45cd7467e11f20932c7a0b74c59318e86581c6b509554876f6c65c114e2c17e42524d300000020994d3802da5adf80345261bcff2eb87ab7b70db786cb0000000000000000000003169efc259f6e4b5e1bfa469f06792d6f07976a098bff2940c8e7ed3105fdc5eff7c65c114e2c170c4dffc30000c020f898b7ea6a405728055b0627f53f42c57290fe78e0b91900000000000000000075472c91a94fa2aab73369c0686a58796949cf60976e530f6eb295320fa15a1b77f8c65c114e2c17387f1df00000002069137421fc274aa2c907dbf0ec4754285897e8aa36332b0000000000000000004308f2494b702c40e9d61991feb7a15b3be1d73ce988e354e52e7a4e611bd9c2a2f8c65c114e2c1740287df200000020ab63607b09395f856adaa69d553755d9ba5bd8d15da20a000000000000000000090ea7559cda848d97575cb9696c8e33ba7f38d18d5e2f8422837c354aec147839fbc65c114e2c175cf077d6000000200ab3612eac08a31a8fb1d9b5397f897db8d26f6cd83a230000000000000000006f4888720ecbf980ff9c983a8e2e60ad329cc7b130916c2bf2300ea54e412a9ed6fcc65c114e2c17d4fbb88500000020d3e51560f77628a26a8fad01c88f98bd6c9e4bc8703b180000000000000000008e2c6e62a1f4d45dd03be1e6692df89a4e3b1223a4dbdfa94cca94c04c22049992fdc65c114e2c17463edb5e",
  signerPubkeyX:
    "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6e",
  signerPubkeyY:
    "0xe8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1",
  concatenatedKeys:
    "0xd4aee75e57179f7cd18adcbaa7e2fca4ff7b1b446df88bf0b4398e4a26965a6ee8bfb23428a4efecb3ebdc636139de9a568ed427fff20d28baa33ed48e9c44e1",
  merkleProof:
    "0x886f7da48f4ccfe49283c678dedb376c89853ba46d9a297fe39e8dd557d1f8deb0fb1a28c03f71b267f3a33459b2566975b1653a1238947ed05edca17ef64181b1f09d858a6e25bae4b0e245993d4ea77facba8ed0371bb9b8a6724475bcdc9edf9ead30b61cf6714758b7c93d1b725f86c2a66a07dd291ef566eaa5a59516823d57fd50557f1d938cc2fb61fe0e1acee6f9cb618a9210688a2965c52feabee66d660a5e7f158e363dc464fca2bb1cc856173366d5d20b5cd513a3aab8ebc5be2bd196b783b8773af2472abcea3e32e97938283f7b454769aa1c064c311c3342a755029ee338664999bd8d432080eafae3ca86b52ad2e321e9e634a46c1bd0d174e38bcd4c59a0f0a78c5906c015ef4daf6beb0500a59f4cae00cd46069ce60db2182e74561028e4462f59f639c89b8e254602d6ad9c212b7c2af5db9275e48c467539c6af678d6f09214182df848bd79a06df706f7c3fddfdd95e6f27326c6217ee446543a443f82b711f48c173a769ae8d1e92a986bc76fca732f088bbe049",
  expectedUTXOOutpoint:
    "0x5f40bccf997d221cd0e9cb6564643f9808a89a5e1c65ea5e6530c0b51c18487c00000000",
  outputValue: 490029088,
  outValueBytes: "0x2040351d00000000",
  prevoutOutpoint:
    "0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000",
  prevoutValueBytes: "0xf078351d00000000",
}

// This is pulled from Bitcoin mainnet transaction data, and the transaction
// data in question has 6 confirmations as opposed to the 1 we use for testnet
// transactions to deal with the unstable difficulty on Bitcoin mainnet.
const depositRoundTrip = {
  fundingTx: {
    difficulty: 16104807485529,
    version: "0x01000000",
    txInputVector:
      "0x01f401c5aed1e842101320b73ae5410f1e927135ca7e002463d8e89484d2c0bda50000000000ffffff00",
    txOutputVector:
      "0x02a0860100000000001600146ff21a38ce66b538e2dfbf78e268763c8aa105bdafafdc1d0000000016001425258f81f2e400e47ef80365e2bde9c93345db3d",
    txLocktime: "0x00000000",
    fundingOutputIndex: 0,
    merkleProof:
      "0x3881a52c9bb1652fa09625d97974e8e2b19a3ad1e5c7ce41e3fa827e4a2dafa9c2e6fd123afcf06a73dbcd8611eb1f9cc1a1d08862ccdb3dc4ad14dede55d108d3e33b3739d6d340d8844bae7e8db2022f590aa6fe2915e056a0003477c8083b68f8e1906e72b81b070c8b929a98ed74ba19fb3e8298457948a1d2174d29926a0003a513b93318055827bc971680e75ba2265543036dcfe06d46eca99c8282e442cbe0adf95bf23a961e76384d08aa455f72d376d1946c85e41e0b68669084b562e862c5833ff76eae5e3f22cd15b1720827a3a327e9affe508dcc4c48a139e24a55a1fdc70dfbe1456562d2847cb684c3b3fdf1fc10830c89ce3becbdc03a75483e1301e5de2d5a4d143a73465700bbc6b4f67650a51b29481601709938436be39a519dd161c56f7cd2661b28006d878d5fff61749368bfac9522afcc1ea196635c577b0ee5c92bb76490eb83dc6fb5aed206f7a50e6af1a3705193bcee5876726ae4f28ee93c8f5f360838591075b30f5e86810afaf016fa048ba3c4e7404a",
    txIndexInBlock: "102",
    bitcoinHeaders:
      "0x0000802081bf9fde653c5831b396f1100d4581ee10a9bf92e85a0500000000000000000065fb243f6fbfe87b315f09df6c1d103c09aed448405ed57ac5948a332c2eb38fcd07bf5e397a1117450be633000080209b8b7afc8ffa3320cb89ba1aa874d5fb5a3f745e3acc090000000000000000000dc3deadaa3b28a0c4f07fd3144998412bba636be4ae03afd50473312d28e5cd1008bf5e397a11173a5d164600000020be7e0bbbd59079dd4965e584d478267424733744d6ad05000000000000000000a3c6d2bf63f06f649bba3f70ec900ef90ff5bfe641088b40903143301ae80075d50dbf5e397a11172755494200004020c24fb2df05de3bee29d31168cc82d008a904684179390000000000000000000074f9ff1d3f11aef31795f21d646d009d6059ff94957067aec6a7ee6a0272c04be50ebf5e397a11171befabf100008020c4ec41da448bab7c1be4b5e5357fa5d28322972c81b70b000000000000000000015654f89c6bb323f17ef92a31832bf1d6d9e926550680a57c571b739358bd1f000fbf5e397a11171ad46a68000000206d5bc32c9a37c05891ee9a260a73edb22647ef29aa1807000000000000000000fe7b1b8fb96e889f6aed650b0c19ae426aefbaccc2040e2a840095dcc573ea3d4211bf5e397a111784dc4846",
  },
  redemptionTx: {
    difficulty: 16104807485529,
    version: "0x01000000",
    outputScript: "0x16001480768eb6cb7b50c8da5c1178665d71e8c290c161",
    outputValueBytes: "0xd07e010000000000",
    txInputVector:
      "0x01dfca5cff90459a5c3add9d64a7fa3db9be9f00ab939aa1e0f0a8b4242d0651a3000000000000000000",
    txOutputVector:
      "0x01d07e01000000000016001480768eb6cb7b50c8da5c1178665d71e8c290c161",
    txLocktime: "0x00000000",
    merkleProof:
      "0x8a0873e58a1335ec458287e5a2df523ebb900ebd67bcffcf32f87df6d1537e372cc9bd8a2e001d4cc8beef2cf56a6960dd62293a2b09571af1850406cd5dcb620a65ee63ab9205bad53c3570348d8b0d8020025f07518091fab53ff249c1eb902a5a3cbb57fe21fa68826ec8f3e61b1ad6a6c7635c72dbfcfdb398c09d28c933c4121e75d58f49d0c4f3d2c0da121752e06e2b280108e78d89b1c64b307c702aac5667d775c18f54c239d4b62dfbd4157bf5a4d03206d9dc3f018dc8a614addc85766bb574bdb308520a6c48efff78969559547885aecc5385908043b81f48d3c9dec5f905b3816c2b96280f8eb463658f20ada1b025a7498c2954fc8ad9209804ba85bb923ce65c98bbd35a81fc1e4413e89c77dc6d4e6bf1574efb9df79ba237ab53843a6b1714f174466d838b45bf181ca88e83eef4f2a7d4958459cebafeded6ed6b373123c63ece31720e854c14d802ca687f933bec19efd6b8a7447d8bc182b315e94295a1b353f6ddad7eeada4dcd403fc7b8b56f3b500249ae66cf9d",
    txIndexInBlock: "2010",
    bitcoinHeaders:
      "0x0000002037663aa575820cb0808f0aa43a71dd33476449f20044040000000000000000001865a9657f89d65a533836ef7ce76d2685ad5763efcaa2604ae1d08f94ee80496e89bf5e397a1117287de35400000020878951b7d5ad6a5e2389be975548a794db04b06d38dd050000000000000000000807a61d298355301971d13527d2aa5859c063be213dbb77bba07bb4878605abd489bf5e397a11174ddefb30000040204a6da5bc5c5296dbd5d3bb3d01a79413c6efc4353b2e0d0000000000000000003cb45755636f5d481d6dbafabc2347756481224b61c9c7ba9fc22b727915979d468abf5e397a1117459044de00e0ff370b0a718aebc9f08f9db9f51e455f1aa8a06c9d62247206000000000000000000d271e46a103502a7d3ed928d9a83bdba4e46000a006bca944a88d8d1e1168f26478fbf5e397a11170528ff5100200020c7b8641da255cec9a3074adb5ed310dd4f2f42cb77e206000000000000000000994bb5f00f875e91844938f295920794314a6a5f3334fe1281340366c0afa8726a92bf5e397a111745d0c9430000c0209518ab1842a42d181fc027139019ae2304e14778eaf90800000000000000000076fb6dc5bb1e4f3d3d9ee30cfdc71a991f39386f1edfdfacb06fe564f9d2eba0d49cbf5e397a111730b17d590000002051cb8a1c57ba8fa4490a3c0a76e3c4a829ce659233e30c000000000000000000a8a9d2e475e386889ef239426017848bfd19a6d90d500824182146828fcb8d6e2ba4bf5e397a11178558f372",
  },
  signerPubkey: {
    x: "0x025df4a0c6e84e385ae5bbcdfe455d31b42989cdd9d3c1ddd47853a54f0717be",
    y: "0xf84082821e1bff2c50e2b253b5e4318a3cc33c1fde640c895fccf92909dc03bd",
    concatenated:
      "0x025df4a0c6e84e385ae5bbcdfe455d31b42989cdd9d3c1ddd47853a54f0717bef84082821e1bff2c50e2b253b5e4318a3cc33c1fde640c895fccf92909dc03bd",
  },
}

module.exports = {
  address0: "0x" + "00".repeat(20),
  bytes32zero: "0x" + "00".repeat(32),
  hash160: hash160,
  states: states,
  LOW_WORK_HEADER: lowWorkHeader,
  deploySystem: deploySystem,
  HEADER_CHAINS: headerChains,
  increaseTime,
  getCurrentTime,
  expectEvent,
  expectNoEvent,
  TX: tx,
  HEADER_PROOFS: headerChains.map(chainToProofBytes),
  fundingTx,
  depositRoundTrip,
  resolveAllLogs,
}
