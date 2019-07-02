var Deposit = artifacts.require("./Deposit.sol");
var KeepBridge = artifacts.require("./KeepBridge.sol");
var TBTCSystem = artifacts.require("./TBTCSystemStub.sol");

const keepRegistry = "0x01E6Ba85b279fE53BC7AeA57A41dA4740aE649E4"; // KeepRegistry contract address

module.exports = async function () {
    let deposit;
    let tbtcSystem;
    let keepBridge;

    try {
        deposit = await Deposit.deployed()
        keepBridge = await KeepBridge.deployed()
        tbtcSystem = await TBTCSystem.deployed()

        await keepBridge.initialize(keepRegistry)
    } catch (err) {
        console.error(`initialization failed: ${err}`);
        process.exit(1)
    }

    async function createNewDeposit() {
        let result = await deposit.createNewDeposit(
            tbtcSystem.address, // address _TBTCSystem,
            "0x0000000000000000000000000000000000000000", // address _TBTCToken,
            keepBridge.address, // address _KeepBridge,
            5, // uint256 _m,
            10 // uint256 _n
        ).catch((err) => {
            console.error(`createNewDeposit failed: ${err}`)
            process.exit(1)
        });

        console.log("createNewDeposit transaction: ", result.tx)
    }

    await createNewDeposit();

    process.exit();
}
