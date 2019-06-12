const fs = require('fs');

var Deposit = artifacts.require("./Deposit.sol");
var KeepMediator = artifacts.require("./KeepBridge.sol");

const keepRegistry = "0xc84Ba13B9a6f3506557dF7A78A9a63d87e1f1FcC"; // KeepRegistry contract address
const TBTCSystem = "0x3790354Be419094B6af6c80b1dB6521c82B75192"; //TBTC System contract address

module.exports = async function () {
    let deposit = await Deposit.deployed();
    let TXID = process.argv[4]
    let headerLen = process.argv[5]
    let keepBridge, inputs
    let version, vin, vout, locktime, merkleProof, index, fundingOutputIndex, bitcoinHeaders

    async function initContracts() {
        deposit = await Deposit.deployed();
        keepBridge = await KeepMediator.deployed();

        await keepBridge.initialize(keepRegistry)
            .catch((err) => {
                console.log(`initialization failed: ${err}`);
            });
    }

    async function callMerkle() {
        if(TXID == undefined || TXID.length < 64){
            console.log('missing TxId argument \nexiting...');
            process.exit();
        } 
    
        const spawn = require("child_process").spawn;
        await spawn('python3',["merkle.py", TXID, headerLen]);
    }

    async function getData(){
        let rawdata = fs.readFileSync('spv_out.json');  
        inputs = JSON.parse(rawdata);
        process.stdout.write('waiting for json');
       
        let i=0;
        while(inputs.txid !== TXID){
            if(i > 15){
                console.log("json didn't update, check merkle.py")
                process.exit();
            }
            process.stdout.write('.');
            await sleep(300);
            rawdata = fs.readFileSync('spv_out.json'); 
            inputs = JSON.parse(rawdata);
            i++
        }
        console.log("");
        console.log(inputs);  
    }

    async function parseData(){
        tx = inputs.tx;
        merkleProof = inputs.proof;
        index = inputs.index;
        bitcoinHeaders = inputs.chain;
        version = await getVersion(tx);
        vin = await getVin(tx);
        vout = await getVout(tx);
        locktime =  await getLocktime(tx);
        fundingOutputIndex = await getOutIndex(tx);

        console.log(version, locktime);
    }

    async function callFundingProof(){

    }

    async function getVersion(tx){ return '0x'+tx.slice(0,8)}
    async function getLocktime(tx){ return '0x'+ tx.slice(tx.length - 8)}

    async function getVin(tx){
        let incount = parseInt('0x'+tx.slice(12,14))
    }

    async function getVout(tx){
        return 0x0
    }

    async function getOutIndex(tx){
        return 0x0
    }

    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
      }

    //await initContracts();
    //await callMerkle();
    await getData();
    await parseData();
    process.exit();

    // const _bitcoinTx = '0x01000000000101179137fe810d6712b6b00c196c1b76a93240f396c53d15f4cd523995dad6f6880000000000ffffffff02581b000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4e8030000000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac0247304402204edb2a92223cd854755ec3e34ce9a0cbeb48ad88a603d4d268203a81245258440220465a1412b8676f1504fae0026c89df342b04b741fcbb55604930339242a5a1630121028896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2e00000000';
    // const txID = '0xd8973d0405cda733f60f92362eec3ae02d4fd52896153f5c8a9f6089862368cb'

    // const _version = '0x01000000'
    // const _vin = `0x01179137fe810d6712b6b00c196c1b76a93240f396c53d15f4cd523995dad6f6880000000000ffffffff`
    // const _vout = '0x02581b000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4e8030000000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac'
    // const _fundingOutputIndex = 0;
    // const _locktime = '0x00000000'
    // const _merkleProof = '0xcb68238689609f8a5c3f159628d54f2de03aec2e36920ff633a7cd05043d97d8529898acab187401f4d1eed73d60cf2efb15361b48c82d84451980c3e645a8bbbc5d20f5aaed58e535d68761a166437bb92cc0c2fff5f9e7e5381e945c95f49cf971736c3bf71a38eaa2022d84f00a3492375bfc68e538ee292ae89e873eac62157d5ca1290a3b412e2adcdadf746c18952be2ddaa66bad24467be0f212b9c45a31bd04725f962b4b18f619e22cc1ab7e906c2cdcaa205d8e51cf1cf199fcf1de2da374c31f1fecb0f9ae6d8366df16da24a63c7c6504713f34bcac7fab3b12227415c10b23066c79980d7081cc31cd38009fa7b14fa807d862d81f9e6989ba6eb018246bbea71dabe571c4fb88b13bcc046569dfad297f406beceaa42ee2766ec5cdfbf35ff7ab6cd6c16744269966a2042dc6ce0b81ecf0a5f51ecd3655ed7';
    // const _merkleValid = '0xcb68238689609f8a5c3f159628d54f2de03aec2e36920ff633a7cd05043d97d800000000'
    // const _index = 71;
    // const _bitcoinHeaders = '0x00000020aaefe576c20c7680e8da9baf4db2fac1f59043b9afaf87361701000000000000ec5cdfbf35ff7ab6cd6c16744269966a2042dc6ce0b81ecf0a5f51ecd3655ed7bd81ee5c453e011a26018de900000020c3188a1be4ff654d8ff147e6f264200ecaeaef5bfe12d3e97000000000000000b3063a20a0567d88dfba3cd6ca4cf0b3bb6d03da5e2a3335a7f0495ce17004cc8086ee5cffff001d37da48a90000002092a7a8318c5e84369c7391b80ab80de961ad30e995244110b73f1600000000000b1d95fd860b17fb38e3c3ce63b623212f12b353dcfd892cfbc92ec9c9ed2684f189ee5c453e011ac750f246000000203991cb8866955373cbab078f25bf38d50a02587a2b342f1a7c00000000000000e865281766dee7ef0ed5786b56e9c05e57f41a50f95c08530e62ffe2f64fc65b598aee5c453e011aa90ed39000000020cc08a7635123b495dee29b1b795bdc9d08ac555cc194519f2f01000000000000b298a563189c2139e0d7e8cb9aa49b681d85b3dd8b930580c33a4af8e0a65b651d8fee5cffff001defe70bcd000000206f04e29a46f3bb4bd3d92f299f5c4fe1bfd4c7f6b6f72dc1b80c180000000000a1ab92186e8980b5a47446578647acbaa54030db04dc4d9d0e5672b9747ae98ecf90ee5c453e011a51e3a3b200000020dd9f42060a303fffd81b3c2cb9ae45897782ccf4e44fc7dfd700000000000000f67fe1436f88023fe2102f3a850c197715d8f394d6887350d95e08f223081911b591ee5c453e011a23a70f1400000020f8d36c426b6c9ddb8b902a33e62a08ba876988ca95ee89ec530000000000000044c9ac396e66df9dfe42f4f85acd871a1e82000f4e3025c20973fc29656851b03792ee5c453e011ad8c1e70d';
    // const _signerX = '0x8896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2e'
    // const _signerY = '0xf687f923b5896a409cb7e2b5ae456f61ac61862305c6ec86bd7421b5bce115e0'
    // const _value = 7000;
    // const _outValueBytes = '0x581b000000000000'
}
