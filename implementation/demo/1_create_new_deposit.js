var Deposit = artifacts.require("./Deposit.sol");
var KeepBridge = artifacts.require("./KeepBridge.sol");
var TBTCSystem = artifacts.require("./TBTCSystemStub.sol");

const keepRegistry = "0x3F2D36e02FbA6c0704738f2AA09f78AFC86934F8"; // KeepRegistry contract address

module.exports = async function () {
    let deposit;
    let tbtcSystem;
    let keepBridge;

    async function initContracts() {
        deposit = await Deposit.deployed();
        keepBridge = await KeepBridge.deployed();
        tbtcSystem = await TBTCSystem.deployed();

        await keepBridge.initialize(keepRegistry)
            .catch((err) => {
                console.log(`initialization failed: ${err}`);
            });
    }

    async function createNewDeposit() {
        console.log('Call createNewDeposit');
        let result = await deposit.createNewDeposit(
            tbtcSystem.address, // address _TBTCSystem,
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
