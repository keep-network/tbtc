module.exports = async function () {
    let txID = process.argv[4]
    let headerLen = 6 // process.argv[5]

    const merkleScript = "/Users/jakub/workspace/bitcoin-spv/scripts/merkle.py"
    let fundingProof = {}

    // async function installPython() {
    //     console.log("Install python environment...\n")

    //     const { spawn } = require('child_process');

    //     const spawnProcess = spawn('pipenv', ['install'])

    //     spawnProcess.stdout.on('data', (data) => {
    //         console.log(`${data}`);
    //     });

    //     spawnProcess.stderr.on('data', (data) => {
    //         console.error(`Failure:\n${data}`)
    //         process.exit(1)
    //     });

    //     spawnProcess.on('close', (code) => {
    //         console.log(`child process exited with code ${code}`);
    //     });
    // }

    async function getTransactionProof() {
        console.log("Get transaction proof...\n")

        if (txID == undefined || txID.length < 64) {
            console.error('missing txID argument');
        }

        console.log(`Transaction ID: ${txID}`)

        const { spawn } = require('child_process');
        const spawnProcess = spawn('pipenv', ['run', 'python', merkleScript, txID, headerLen])

        spawnProcess.stdout.on('data', data => {
            console.log(`Received data from bitcoin-spv`);
            spvProof = parseBitcoinSPVOutput(data.toString())

            fundingProof.merkleProof = spvProof.merkleProof
            fundingProof.txInBlockIndex = spvProof.txInBlockIndex
            fundingProof.chainHeaders = spvProof.chainHeaders

            parseTransaction(spvProof.tx)
        });

        spawnProcess.stderr.on('data', (data) => {
            console.error(`Failure:\n${data}`)
            process.exit(1)
        });

        spawnProcess.on('close', (code) => {
            console.log(`child process exited with code ${code}`);
            return
        });
    }

    function parseBitcoinSPVOutput(output) {
        console.log("Parse bitcoin-spv output...\n")

        let tx = output.match(/(^-* TX -*$\n)^(.*)$/m)[2]
        let merkleProof = output.match(/(^-* PROOF -*$\n)^(.*)$/m)[2]
        let txInBlockIndex = output.match(/(^-* INDEX -*$\n)^(.*)$/m)[2]
        let chainHeaders = output.match(/(^-* CHAIN -*$\n)^(.*)$/m)[2]

        return {
            tx: tx,
            merkleProof: '0x' + merkleProof,
            txInBlockIndex: txInBlockIndex,
            chainHeaders: '0x' + chainHeaders
        };
    }

    async function parseTransaction(tx) {
        console.log(`Parse transaction...\nTX: ${tx}`)

        let bytes = web3.utils.hexToBytes('0x' + tx)

        fundingProof.version = web3.utils.bytesToHex(getVersion(bytes))
        fundingProof.txInVector = web3.utils.bytesToHex(getTxInputVector(bytes))
        fundingProof.txOutVector = web3.utils.bytesToHex(getTxOutputVector(bytes))
        fundingProof.locktime = web3.utils.bytesToHex(getLocktime(bytes))
        fundingProof.fundingOutputIndex // TODO: find output index

        console.log(`Result\n${JSON.stringify(fundingProof)}`)
    }

    function getPrefix(tx) {
        if (tx.slice(5, 7) == "0001") {
            return tx.slice(0, 7)
        }
        return tx.slice(0, 4)
    }

    function getVersion(tx) {
        return tx.slice(0, 4)
    }

    function isFlagPresent(tx) {
        // TODO: Check for witness transaction
        if (tx.slice(5, 7) == "0001") {
            return true
        }
        return false
    }

    function getTxInputVector(tx) {
        let txInVectorStartPosition = getPrefix(tx).length
        let txInVectorEndPosition

        if (isFlagPresent(tx)) {
            // TODO: Implement for witness transaction
            console.error("witness not supported")
            process.exit(1)
        } else {
            inputCount = parseInt('0x' + tx.slice(txInVectorStartPosition, txInVectorStartPosition + 1))

            if (inputCount != 1) {
                // TODO: Support multiple inputs
                console.error(`exactly one input is required, but got ${inputCount}`);
                process.exit()
            } else {
                startPos = txInVectorStartPosition + 1

                previousHash = tx.slice(startPos, startPos + 32).reverse()

                previousOutIndex = parseInt(tx.slice(startPos + 32, startPos + 36))

                scriptLength = parseInt(tx.slice(startPos + 36, startPos + 37))
                if (scriptLength >= 253) {
                    console.error(`VarInts not supported`);
                    process.exit(1)
                }

                script = tx.slice(startPos + 37, startPos + 37 + scriptLength)

                sequenceNumber = tx.slice(startPos + 37 + scriptLength, startPos + 37 + scriptLength + 4)

                txInVectorEndPosition = startPos + 37 + scriptLength + 4
            }
        }
        return tx.slice(txInVectorStartPosition, txInVectorEndPosition)
    }

    function getTxOutputVector(tx) {
        txPrefix = getPrefix(tx)
        txInput = getTxInputVector(tx)

        outStartPosition = txPrefix.length + txInput.length

        outputsCount = tx.slice(outStartPosition, outStartPosition + 1)

        startPosition = outStartPosition + 1
        let outEndPosition

        for (let i = 0; i < outputsCount; i++) {
            value = tx.slice(startPosition, startPosition + 8)
            scriptLength = parseInt(tx.slice(startPosition + 8, startPosition + 8 + 1))
            if (scriptLength >= 253) {
                console.error(`VarInts not supported`);
                process.exit()
            }

            script = tx.slice(startPosition + 8 + 1, startPosition + 8 + 1 + scriptLength)

            outEndPosition = startPosition + 8 + 1 + scriptLength
            startPosition = outEndPosition
        }

        return tx.slice(outStartPosition, outEndPosition)
    }

    function getLocktime(tx) {
        return tx.slice(tx.length - 8)
    }

    await getTransactionProof();
}
