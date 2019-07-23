const FundingProof = require('./scripts/FundingProof.js')

const TX_ID = '72e7fd57c2adb1ed2305c4247486ff79aec363296f02ec65be141904f80d214e'
const CONFIRMATIONS = 6

async function main(txID, confirmations) {
  const result = await FundingProof.getTransactionProof(txID, confirmations)
    .catch((err) => {
      console.error(`failed to get transaction proof: [${err}]`)
      process.exit(1)
    })

  result.fundingOutputIndex = 0 // TODO: Find index in transaction based on deposit's public key

  console.log('Funding Proof:', result)
  process.exit(0)
}

main(TX_ID, CONFIRMATIONS)
