let Deposit = artifacts.require("./Deposit.sol");
let TBTCSystem = artifacts.require("./TBTCSystemStub.sol");

const FundingProof = require('./tools/FundingProof')

module.exports = async function () {
    let deposit
    let depositLog

    let txID = process.argv[4]
    let headerLen = process.argv[5]

    async function initContracts() {
        try {
            deposit = await Deposit.deployed()
            depositLog = await TBTCSystem.deployed();
        } catch (err) {
            console.error(`call failed: ${err}`)
            process.exit(1)
        }
    }

    async function callFundingProof(fundingProof) {
        console.log('Submit funding proof...');

        let blockNumber = await web3.eth.getBlock("latest").number

        let result = await deposit.provideBTCFundingProof(
            fundingProof.version,
            fundingProof.txInVector,
            fundingProof.txOutVector,
            fundingProof.locktime,
            fundingProof.fundingOutputIndex,
            fundingProof.merkleProof,
            fundingProof.txInBlockIndex,
            fundingProof.chainHeaders
        ).catch((err) => {
            console.error(`provideBTCFundingProof failed: ${err}`)
        });

        console.log("provideBTCFundingProof transaction: ", result.tx)

        let eventList = await depositLog.getPastEvents('Funded', {
            fromBlock: blockNumber,
            toBlock: 'latest'
        })

        console.log("Funding proof accepted for the deposit: ", eventList[0].returnValues._depositContractAddress)
    }

    await initContracts()
    await FundingProof.getTransactionProof(txID, headerLen, callFundingProof)
}
