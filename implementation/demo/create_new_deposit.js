var Deposit = artifacts.require("./Deposit.sol");
var KeepMediator = artifacts.require("./KeepBridge.sol");

const keepRegistry = "0x41649ff7E9E4512fbe8b42A51d73f33377D012c5"; // KeepRegistry contract address

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
            "0x0000000000000000000000000000000000000000", // address _TBTCSystem,
            "0x0000000000000000000000000000000000000000", // address _TBTCToken,
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
