const headerChains = require("./headerchains.json")
const tx = require("./tx.json")
const createHash = require("create-hash")
const {web3} = require("@openzeppelin/test-environment")

const ozHelpers = require("@openzeppelin/test-helpers")
const {BN} = ozHelpers
const ozExpectEvent = ozHelpers.expectEvent
const {expect} = require("chai")

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
  const t = createHash("sha256")
    .update(buffer)
    .digest()
  const u = createHash("rmd160")
    .update(t)
    .digest()
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
        param => linkable[param] || param,
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

  return new Promise((resolve, reject) => {
    web3.currentProvider.send(
      {
        jsonrpc: "2.0",
        method: "evm_increaseTime",
        params: [duration],
        id: id,
      },
      err1 => {
        if (err1) return reject(err1)

        web3.currentProvider.send(
          {
            jsonrpc: "2.0",
            method: "evm_mine",
            id: id + 1,
          },
          (err2, res) => {
            return err2 ? reject(err2) : resolve(res)
          },
        )
      },
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
    const contracts =
        Object
            .entries(contractContainer)
            .map(([, value]) => value)
            .filter(_ => _.contract && _.address)

    const { resolved: resolvedLogs } = contracts.reduce(
        ({ raw, resolved }, contract) => {
            const events = contract.contract._jsonInterface.filter(_ => _.type === "event")
            const contractLogs = raw.filter(_ => _.address == contract.address)

            const decoded = contractLogs.map(log => {
                const event = events.find(_ => log.topics.includes(_.signature))
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
                raw: raw.filter(_ => _.address != contract.address),
                resolved: resolved.concat(decoded),
            }
        },
        { raw: receipt.rawLogs, resolved: [] },
    )

    return {
        ...receipt,
        logs: resolvedLogs,
    }
}

/**
 * Wrapper for OpenZeppelin's expectEvent helper that deals with array-of-BN
 * parameters.
 *
 * @param {TxReceipt} receipt The receipt to check for the specified event.
 * @param {string} eventName The name of the event to look for.
 * @param {object} parameters The parameters to look for in the event; unlike
 *         OpenZeppelin's default version, parameters may have an array-of-BN
 *         value and they will be properly validated.
 */
function expectEvent(receipt, eventName, parameters) {
  const bnArrayParameterNames = []
  for (const [parameterName, parameterValue] of Object.entries(parameters)) {
    if (web3.utils.isBN(parameterValue[0])) {
      bnArrayParameterNames.push(parameterName)
    }
  }

  // Use OpenZeppelin helper for all non-array-of-BN parameters.
  const withoutBnArray = Object.assign({}, parameters)
  bnArrayParameterNames.forEach(_ => delete withoutBnArray[_])
  ozExpectEvent(receipt, eventName, withoutBnArray)

  if (bnArrayParameterNames.length > 0) {
    // Check array-of-BN parameters in nested fashion, directly.
    const log = receipt.logs.find(_ => _.event == eventName)
    bnArrayParameterNames.forEach(paramName => {
      log.args[paramName].forEach((bn, i) => {
        expect(bn).to.eq.BN(parameters[paramName][i])
      })
    })
  }

/**
 * Similar to array.reduce, but accumulates promises instead of non-promises and
 * ensures that each step in the reduction waits on the previous one. Standard
 * reduces with async functions fire all the functions at once and require
 * explicit internal waits in the reducer; this function pulls that out and
 * correctly invokes each step of the reducer after the previous one's promise
 * has settled.
 *
 * @typeparam T The type of object in the array.
 * @typeparam A The type of object the reducer accumulates.
 * @param {T[]} array The array to reduce over.
 * @param {(A, T)=>A)} reducer The reducer that combines values of type T into
 *        an accumulator of type A.
 * @param {A} initialValue The initial value for the reduce, passed as the first
 *        parameter to the first call to the reducer.
 */
async function asyncReduce(array, reducer, initialValue) {
    return array.reduce(
        async (previousValue, nextValue) => {
            const realPrev = await previousValue
            return reducer(realPrev, nextValue)
        },
        initialValue,
    )
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
  txInputVector: "0x01913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d0000000000ffffffff",
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
  prevoutOutpoint: "0x913e39197867de39bff2c93c75173e086388ee7e8707c90ce4a02dd23f7d2c0d00000000",
  prevoutValueBytes: "0xf078351d00000000",
}

const legacyFundingTx = {
  difficulty: 6113,
  version: "0x01000000",
  txInputVector:
      "0x01f3002663bbdfe1a0c02bafe6ad6bdc713902bbe678bdffd7e0a77972a4a6c4b2000000006b483045022100dfaa9366f2fd5a233e8029ae2f73e31dd75d85ef8970a74145c3128c0463cd9e0220739c9a57f1266d835b39f3f10d5df77d855ed4110ded0443584f4b7835c6167c012102b6ec2e4df2062066b685fd5cfe47f33c54b29a79a87ba20fcdb66f0b87170006ffffffff",
  txOutputVector:
      "0x027818613b000000001976a914a29c9701293e838fc01ca1c13b010e42cf53242388aca08601000000000017a91414ea6747c79c9c7388eca6c029e6cc81b5e0951c87",
  fundingOutputIndex: 0,
  txLocktime: "0x00000000",
  txIndexInBlock: 1,
  bitcoinHeaders:
    "0x03000000a3ecf52cf38b65970661436cab0b9813543fa949ac7b27db901ee402000000009873fc6b41300fb93957dc4fdec91f72f12e698bc5648a5c57623878a173002e2a9f9e5560b80a1baa34e3b5030000000d801f1aab80f52ca35e510df974c96901d1471c4eec0a377f570700000000004696408b572b53ed3df66479a4a3f69c5b6585daa3778d905483993dfad29360eda39e55ffff001d046ef73b0300000020cfaa259c8e375766f74be88525c7d1e41430d42a13844249a4f70300000000e9929e791caffe690fbb9406d06ef9376604bbf380dde84f64e98b5943e1a400039f9e5560b80a1b0e17707d03000000f818ca6e20d44ecb355e0897bd14eb50dcc35145f8d7554ebf2c040000000000ffe8e2da3ed2ae815e5948cef4addf32f7c599427a7d87e39d8c063c4dd30368cca39e55ffff001de8319459030000005b45b4375c101d7a95764cfda6e6f7b40246867cd54e89855b73200100000000fc074723d089d5338b65800707544f509bf111bc4d17a2b0353978319145fc5a459f9e5560b80a1b3fd1c3d2030000001cfc4d39d6ce1bebb6511ecb2039bb6fc26016eafb4d1e010be8020000000000e00725a1bca5411e7072970dccee9eeab802ca6c9cbcccd8a0aa39c40251bcbd02a49e55ffff001dd247f8d203000000cf778c9ea94a7f0eb51b06c29e65386f55efca5f3714c535b6e68703000000004bb0084076f57bce7a03f74da9d8784979ea0dbeefcce590ebd62d246b20bd6ab5a39e5560b80a1bbd5c1c1403000000e24a5e009fe23639a71f1024c6590abf3e2ffbe0ac92a43207ee0200000000007f740aa5ae4ea28ec5a3bd29151f3d053e10bb0489eb9a53b3f41604d05bfbccdaa39e5560b80a1bf27e379a03000000020b75077b8f84f296ee015471e7d99b0a57ff92acfb7739f62909000000000028cc108626e3c6f6fa12228c927b1ebec587b570bf085dcd24b63bb9b730c569eea39e5560b80a1b17793fc8",
  signerPubkeyX:
    "0x85134688b6be6f0fd93ce6a7f7904e9dba2b22bf3cc65c0feb5067b60cf5d821",
  signerPubkeyY:
    "0x76f80a7d522ea754db0e25ca43fdacfd1f313e4dc2765e2dfcc18fb3d63a66c4",
  concatenatedKeys:
    "0x85134688b6be6f0fd93ce6a7f7904e9dba2b22bf3cc65c0feb5067b60cf5d82176f80a7d522ea754db0e25ca43fdacfd1f313e4dc2765e2dfcc18fb3d63a66c4",
  merkleProof:
    "0x0efe91abb6a919b0f89c4c785460e2c2e50b57f3d9eb8da0885c26b02ddb50cb",
  expectedUTXOOutpoint:
    "0x0ee73932b031135c57e3d8f53db8ed5c97e6023a3e4980ea465f1aa2962d17b200000000",
  _outputValue: 996219000,
  outValueBytes: "0x7818613b00000000",
}

module.exports = {
  address0: "0x" + "00".repeat(20),
  bytes32zero: "0x" + "00".repeat(32),
  hash160: hash160,
  states: states,
  LOW_WORK_HEADER: lowWorkHeader,
  deploySystem: deploySystem,
  HEADER_CHAINS: headerChains,
  increaseTime: increaseTime,
  expectEvent: expectEvent,
  TX: tx,
  HEADER_PROOFS: headerChains.map(chainToProofBytes),
  fundingTx,
  legacyFundingTx,
  resolveAllLogs,
  asyncReduce,
}
