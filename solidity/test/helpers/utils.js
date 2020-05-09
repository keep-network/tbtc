const headerChains = require("./headerchains.json")
const tx = require("./tx.json")
const createHash = require("create-hash")
const BN = require("bn.js")
const {web3} = require("@openzeppelin/test-environment")

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

                return {
                    ...log,
                    event: event.name,
                    args: decoded,
                }
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

module.exports = {
  address0: "0x" + "00".repeat(20),
  bytes32zero: "0x" + "00".repeat(32),
  hash160: hash160,
  states: states,
  LOW_WORK_HEADER: lowWorkHeader,
  deploySystem: deploySystem,
  HEADER_CHAINS: headerChains,
  increaseTime: increaseTime,
  TX: tx,
  HEADER_PROOFS: headerChains.map(chainToProofBytes),
  resolveAllLogs,
  asyncReduce,
}
