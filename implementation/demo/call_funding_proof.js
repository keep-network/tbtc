let Deposit = artifacts.require("./Deposit.sol");

const FundingProof = require('./get_funding_proof')

module.exports = async function () {
    let deposit

    let txID = process.argv[4]
    let headerLen = process.argv[5]
    let fundingOutputIndex = process.argv[6]

    async function initContracts() {
        deposit = await Deposit.deployed().catch((err) => {
            console.error(`call failed: ${err}`)
            process.exit(1)
        });
    }

    async function callFundingProof(fundingProof) {
        console.log('Submit funding proof...');
        console.log(fundingProof);

        let result = await deposit.provideBTCFundingProof(
            fundingProof.version,
            fundingProof.txInVector,
            fundingProof.txOutVector,
            fundingProof.locktime,
            fundingOutputIndex,
            fundingProof.merkleProof,
            fundingProof.txInBlockIndex,
            fundingProof.chainHeaders
        ).catch((err) => {
            console.error(`call failed: ${err}`)
        });

        console.log('funding call tx: ' + result.tx)
    }

    await initContracts()
    await FundingProof.getTransactionProof(txID, headerLen, callFundingProof)
}
