var Deposit = artifacts.require("./Deposit.sol");
let TBTCSystemStub = artifacts.require("./TbtcSystemStub.sol");

module.exports = async function () {
    let deposit = await Deposit.deployed();
    let depositLog = await TBTCSystemStub.deployed();

    async function getPublicKey() {
        let blockNumber = await web3.eth.getBlock("latest").number

        console.log('Call getPublicKey');
        let result = await deposit.retrieveSignerPubkey()
            .catch((err) => {
                console.log(`creation failed: ${err}`);
            });

        console.log(result.tx)

        let eventList = await depositLog.getPastEvents('RegisteredPubkey', {
            fromBlock: blockNumber,
            toBlock: 'latest'
        })

        console.log(eventList[0].returnValues)
    }

    await getPublicKey();

    process.exit();
}
