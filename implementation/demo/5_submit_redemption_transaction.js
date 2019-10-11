const Deposit = artifacts.require('./Deposit.sol')
const TBTCSystem = artifacts.require('./TBTCSystem.sol')

const txUtils = require('./tools/Transaction.js')
const BN = require('bn.js')

module.exports = async function() {
  try {
    // Parse arguments
    const depositAddress = process.argv[4]

    let deposit
    let depositLog
    let keepAddress

    // try {
    //   deposit = await Deposit.deployed()
    //   depositLog = await TBTCSystem.deployed()
    // } catch (err) {
    //   console.error(`initialization failed: ${err}`)
    //   process.exit(1)
    // }

    // try {
    //   const depositCreatedEvents = await depositLog.getPastEvents('Created', {
    //     fromBlock: 0,
    //     toBlock: 'latest',
    //     filter: { _depositContractAddress: depositAddress },
    //   })
    //   keepAddress = depositCreatedEvents[0].returnValues._keepAddress
    // } catch (err) {
    //   console.error(`failed to get keep address: ${err}`)
    //   process.exit(1)
    // }

    // const events = await depositLog.getPastEvents(
    //   'RedemptionRequested',
    //   {
    //     filter: { _depositContractAddress: depositAddress },
    //     fromBlock: 0,
    //     toBlock: 'latest',
    //   }
    // ).catch((err) => {
    //   console.error(`failed to get past redemption requested events`)
    //   process.exit(1)
    // })

    // let latestEvent
    // if (events.length > 0) {
    //   latestEvent = events[events.length - 1]
    // } else {
    //   console.error(`events list is empty`)
    //   process.exit(1)
    // }

    // Sample values for testing
    latestEvent = {
      returnValues: {
        _depositContractAddress: '0x582a29981e0691B2E8b8E78E009aA34b0713c282',
        _requester: '0x4f76Eb125610290301a6F70F535Af9838F48DbC1',
        _digest: '0xb68a6378ddb770a82ae4779a915f0a447da7d753630f8dd3b00be8638677dd90',
        _utxoSize: new BN('1229782938247303441', 10),
        _requesterPKH: '0x3333333333333333333333333333333333333333',
        _requestedFee: new BN('1229782937960972288', 10),
        _outpoint: '0x333333333333333333333333333333333333333333333333333333333333333333333333',
      },
    }

    let unsignedTransaction
    try {
      const utxoSize = new BN(latestEvent.returnValues._utxoSize)
      const requesterPKH = Buffer.from(web3.utils.hexToBytes(latestEvent.returnValues._requesterPKH))
      const requestedFee = new BN(latestEvent.returnValues._requestedFee)
      const outpoint = Buffer.from(web3.utils.hexToBytes(latestEvent.returnValues._outpoint))

      const outputValue = utxoSize.sub(requestedFee)

      unsignedTransaction = txUtils.oneInputOneOutputWitnessTX(
        outpoint,
        0, // AS PER https://github.com/summa-tx/bitcoin-spv/blob/2a9d594d9b14080bdbff2a899c16ffbf40d62eef/solidity/contracts/CheckBitcoinSigs.sol#L154
        outputValue,
        requesterPKH
      )
    } catch (err) {
      console.error(`failed to get transaction preimage: ${err}`)
      process.exit(1)
    }
    console.log('Transaction preimage:\n', unsignedTransaction)

    // Get keep public key
    let keepPublicKey
    try {
      // TODO: Get public key from event
      keepPublicKey = Buffer.from('657282135ed640b0f5a280874c7e7ade110b5c3db362e0552e6b7fff2cc8459328850039b734db7629c31567d7fc5677536b7fc504e967dc11f3f2289d3d4051', 'hex')
    } catch (err) {
      console.error(`failed to get public key: ${err}`)
      process.exit(1)
    }

    // Get signature calculated by keep
    let signatureR
    let signatureS
    try {
      // TODO: get signature from event
      signatureR = Buffer.from('9b32c3623b6a16e87b4d3a56cd67c666c9897751e24a51518136185403b1cba2', 'hex')
      signatureS = Buffer.from('90838891021e1c7d0d1336613f24ecab703dee5ff1b6c8881bccc2c011606a35', 'hex')
    } catch (err) {
      console.error(`failed to get signature: ${err}`)
      process.exit(1)
    }

    // Add witness signature to transaction
    let signedTransaction
    try {
      signedTransaction = txUtils.addWitnessSignature(
        unsignedTransaction,
        0,
        signatureR,
        signatureS,
        keepPublicKey
      )
    } catch (err) {
      console.error(`failed to add witness to transaction: ${err}`)
      process.exit(1)
    }

    console.log('Signed transaction:', signedTransaction)

    // Publish transaction to bitcoin chain
    // TODO: Submit raw signer transaction to bitcoin via electrum
  } catch (err) {
    console.error(err)
    process.exit(1)
  }
  process.exit()
}
