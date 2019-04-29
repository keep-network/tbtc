const headerChains = require('./headerchains.json')
const tx = require('./tx.json')
const createHash = require('create-hash')
const BN = require('bn.js')

// genesis header -- diff 1
lowDiffHeader = '0x0100000000000000000000000000000000000000000000000000000000000000000000003BA3EDFD7A7B12B27AC72C3E67768F617FC81BC3888A51323A9FB8AA4B1E5E4A29AB5F49FFFF001D1DAC2B7C'

states = {
  START: new BN(0),
  AWAITING_SIGNER_SETUP: new BN(1),
  AWAITING_BTC_FUNDING_PROOF: new BN(2),
  FRAUD_AWAITING_BTC_FUNDING_PROOF: new BN(3),
  FAILED_SETUP: new BN(4),
  ACTIVE: new BN(5),
  AWAITING_WITHDRAWAL_SIGNATURE: new BN(6),
  AWAITING_WITHDRAWAL_PROOF: new BN(7),
  REDEEMED: new BN(8),
  COURTESY_CALL: new BN(9),
  FRAUD_LIQUIDATION_IN_PROGRESS: new BN(10),
  LIQUIDATION_IN_PROGRESS: new BN(11),
  LIQUIDATED: new BN(12)
}

function hash160 (hexString) {
    let buffer = Buffer.from(hexString, 'hex')
        var t = createHash('sha256').update(buffer).digest()
            var u = createHash('rmd160').update(t).digest()
                return '0x' + u.toString('hex')
}

function chainToProofBytes(chain) {
  byteString = '0x'
  for (let header = 6; header < chain.length; header++) {
    byteString += chain[header].hex
  }
  return byteString
}

async function deploySystem(deploy_list) {
  deployed = {}  // name: contract object
  linkable = {}  // name: linkable address
  for (let i in deploy_list) {
    await deploy_list[i].contract.link(linkable)
    contract = await deploy_list[i].contract.new()
    linkable[deploy_list[i].name] = contract.address
    deployed[deploy_list[i].name] = contract
  }
  return deployed
}

module.exports = {
  address0: '0x' + '00'.repeat(20),
  bytes32zero: '0x' + '00'.repeat(32),
  hash160: hash160,
  states: states,
  LOW_DIFF_HEADER: lowDiffHeader,
  deploySystem: deploySystem,
  HEADER_CHAINS: headerChains,
  TX: tx,
  HEADER_PROOFS: headerChains.map(chainToProofBytes)
}
