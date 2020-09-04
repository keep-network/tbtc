const DepositFactory = artifacts.require("DepositFactory.sol")
const Deposit = artifacts.require("Deposit.sol")
const BondedECDSAKeepJson = require("@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json")

const truffleContract = require("@truffle/contract")

const deploymentBlock = 8594983

module.exports = async function() {
    try {
        const factory = await DepositFactory.at(
            "0x4CEE725584e38413603373C9D5df593a33560293"
        )

        const depositCreatedEvents = await factory.getPastEvents(
            "DepositCloneCreated",
            {
                fromBlock: deploymentBlock,
                toBlock: "latest",
            }
        )

        console.log(`Number of created deposits: ${depositCreatedEvents.length}`)

        const depositAddresses = []
        depositCreatedEvents.forEach((event) => 
          depositAddresses.push(event.args.depositCloneAddress)
        )

        for (i = 0; i < depositAddresses.length; i++) {
            const deposit = await Deposit.at(depositAddresses[i])
            const state = await deposit.currentState()
            let stateString = ""
            switch(state.toString()) {
                case "0": stateString = "START"; break
                case "1": stateString = "AWAITING_SIGNER_SETUP"; break
                case "2": stateString = "AWAITING_BTC_FUNDING_PROOF"; break
                case "3": stateString = "FAILED_SETUP"; break
                case "4": stateString = "ACTIVE"; break
                case "5": stateString = "AWAITING_WITHDRAWAL_SIGNATURE"; break
                case "6": stateString = "AWAITING_WITHDRAWAL_PROOF"; break
                case "7": stateString = "REDEEMED"; break
                case "8": stateString = "COURTESY_CALL"; break
                case "9": stateString = "FRAUD_LIQUIDATION_IN_PROGRESS"; break
                case "10": stateString = "LIQUIDATION_IN_PROGRESS"; break
                case "11": stateString = "LIQUIDATED"; break
                default: stateString = "<< UNKNOWN >>"; break
            }
            const keepAddress = await deposit.keepAddress()
            const lotSizeSatoshis = await deposit.lotSizeSatoshis()
            const lotSizeTbtc = await deposit.lotSizeTbtc()

            console.log(`deposit address: ${depositAddresses[i]}`)
            console.log(`deposit index:   ${i}`)
            console.log(`deposit state:   ${stateString}`)
            console.log(`keep address:    ${keepAddress}`)
            console.log(`lot size [sat]:  ${lotSizeSatoshis}`)
            console.log(`lot size [tbtc]: ${lotSizeTbtc}`)

            if (stateString == "AWAITING_SIGNER_SETUP") {
                await retrieveSignerPubkey(deposit)
            } else if (stateString == "AWAITING_BTC_FUNDING_PROOF") {            
                await notifyFundingTimedOut(deposit)
            } else if (stateString == "AWAITING_WITHDRAWAL_SIGNATURE") {
                await provideRedemptionSignature(deposit, keepAddress)
            } else if (stateString == "AWAITING_WITHDRAWAL_PROOF") {
                await notifyRedemptionProofTimeout(deposit)
            }

            console.log(``)
        }

        process.exit()
    } catch (error) {
      console.log(error)
      process.exit()
    }
}

async function notifySignerSetupFailed(deposit) {
    console.log(`Notifying signer setup failed...`)
    try {
        const receipt = await deposit.notifySignerSetupFailed()
        console.log(`TX: ${receipt.tx}`)
    } catch (error) {
        console.log(`Failed with: ${JSON.stringify(error)}`)
    }
}

async function retrieveSignerPubkey(deposit) {
    console.log(`Retrieving signer pubkey...`)
    try {
        const receipt = await deposit.retrieveSignerPubkey()
        console.log(`TX: ${receipt.tx}`)
    } catch (error) {
        console.log(`Failed with: ${JSON.stringify(error)}`)
    }
}

async function notifyFundingTimedOut(deposit) {
    console.log(`Notifying funding timeout...`)
    try {
        const receipt = await deposit.notifyFundingTimedOut()
        console.log(`TX: ${receipt.tx}`)
    } catch (error) {
        console.log(`Failed with: ${JSON.stringify(error)}`)
    }
}

async function provideRedemptionSignature(deposit, keepAddress) {
    const BondedECDSAKeep = truffleContract(BondedECDSAKeepJson)
    BondedECDSAKeep.setProvider(web3.currentProvider)

    const keep = await BondedECDSAKeep.at(keepAddress)

    console.log(`Looking for the signature from keep ${keepAddress}...`)

    console.log(await keep.digest())
    const signatureSubmittedEvents = await keep.getPastEvents(
        "SignatureSubmitted",
        {
            fromBlock: deploymentBlock,
            toBlock: "latest",
        }
    )

    const signaturesCount = signatureSubmittedEvents.length
    console.log(`Number of signatures submitted: ${signaturesCount}`)
   
    if (signaturesCount == 0) {
        // no signatures, nothing to do here
        return
    }
    const lastSignature = signatureSubmittedEvents[signaturesCount - 1]

    const {digest, r, s, recoveryID } = lastSignature.returnValues
    // A constant in the Ethereum ECDSA signature scheme, used for public key recovery [1]
    // Value is inherited from Bitcoin's Electrum wallet [2]
    // [1] https://bitcoin.stackexchange.com/questions/38351/ecdsa-v-r-s-what-is-v/38909#38909
    // [2] https://github.com/ethereum/EIPs/issues/155#issuecomment-253810938
    const ETHEREUM_ECDSA_RECOVERY_V = web3.utils.toBN(27)
    const v = web3.utils.toBN(recoveryID).add(ETHEREUM_ECDSA_RECOVERY_V)

    console.log(`digest = ${digest}, r = ${r}, s = ${s}, v = ${v}`)
    
    try {
        console.log(`Providing redemption signature...`)
        const receipt = await deposit.provideRedemptionSignature(
            v,
            r.toString(),
            s.toString()
        )
        console.log(`TX: ${receipt.tx}`)
    } catch (error) {
        console.log(`Failed with: ${JSON.stringify(error)}`)
    }
}

async function notifyRedemptionProofTimeout(deposit) {
    console.log(`Notifying redemption proof timeout...`)
    try {
        const receipt = await deposit.notifyRedemptionProofTimedOut()
        console.log(`TX: ${receipt.tx}`)
    } catch (error) {
        console.log(`Failed with: ${JSON.stringify(error)}`)
    }
}