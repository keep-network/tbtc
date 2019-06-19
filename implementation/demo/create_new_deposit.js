var Deposit        = artifacts.require("./Deposit.sol");
var KeepMediator   = artifacts.require("./KeepBridge.sol");

const keepRegistry = "0xdF534b280A95B614F18E9e2204f2599AB6fba84B"; // KeepRegistry contract address
const TBTCSystem   = "0x39B2F1020e1692A495261D89a58624946DcE825e"; //TBTC System contract address
const TBTCToken    = "0xc4Ba466Cf38B21AE7cCd9D10d546D17485dD50Ca"; //TBTC TOKEN

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
