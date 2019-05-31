var Deposit = artifacts.require("./Deposit.sol");
var KeepMediator = artifacts.require("./KeepBridge.sol");

const keepRegistry = "0xc84Ba13B9a6f3506557dF7A78A9a63d87e1f1FcC"; // KeepRegistry contract address
const TBTCSystem = "0xDb3af74FDB714B09184c09a5783012ED42Fce817";
const TBTCToken = "0x8d119BD689BA0a3bE499BfCe7AaaD2a6fDd90a16"

module.exports = async function () {
    let deposit = await Deposit.deployed();
    let keepBridge;

    async function initContracts() {
        deposit = await Deposit.deployed();
        keepBridge = await KeepMediator.deployed();

        await keepBridge.initialize(keepRegistry)
            .catch((err) => {
                console.log(`initialization failed: ${err}`);
            });
    }

    async function createNewDeposit() {
        console.log('Call createNewDeposit');
        let result = await deposit.createNewDeposit(
            TBTCSystem, // address _TBTCSystem,
            TBTCToken, // address _TBTCToken,
            keepBridge.address, // address _KeepBridge,
            5, // uint256 _m,
            10 // uint256 _n
        ).catch((err) => {
            console.log(`creation failed: ${err}`);
        });

        console.log('deposit creation tx: ' + result.tx);
    }

    await initContracts();
    await createNewDeposit();

    process.exit();
}
