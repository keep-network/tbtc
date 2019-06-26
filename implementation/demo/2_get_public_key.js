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
                console.log(`retrieveSignerPubkey failed: ${err}`);
            });

        console.log("retrieveSignerPubkey transaction: ", result.tx)

        let eventList = await depositLog.getPastEvents('RegisteredPubkey', {
            fromBlock: blockNumber,
            toBlock: 'latest'
        })

        let publicKeyX = eventList[0].returnValues._signingGroupPubkeyX
        let publicKeyY = eventList[0].returnValues._signingGroupPubkeyY

        console.log(`Registered public key:\nX: ${publicKeyX}\nY: ${publicKeyY}`)
    }

    await getPublicKey();

    process.exit();
}
