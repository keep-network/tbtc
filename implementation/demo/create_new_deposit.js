var Deposit = artifacts.require("./Deposit.sol");
var KeepMediator = artifacts.require("./KeepBridge.sol");

const keepRegistry = "0xc84Ba13B9a6f3506557dF7A78A9a63d87e1f1FcC"; // KeepRegistry contract address
const TBTCSystem = "0x3c58e11260c0Fe0e091ceb19c518d169Fcbc350c";
const ERC20 = "0x507adB87B2192544C3e79C846b4408B7C988356d"

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
            ERC20, // address _TBTCToken,
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
