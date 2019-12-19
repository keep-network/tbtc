const fs = require('fs')

let doc =
`:toc: true
:toclevels: 2

= tBTC API Documentation

Welcome to the tBTC API Documentation. The primary contracts involved in tBTC
are listed below, along with their public methods.

toc::[]
`

let jsonFiles = [
    "../../implementation/build/contracts/TBTCSystem.json",
    // "../../implementation/build/contracts/TBTCConstants.json",
    "../../implementation/build/contracts/Deposit.json",
    "../../implementation/build/contracts/DepositStates.json",
    "../../implementation/build/contracts/DepositFunding.json",
    "../../implementation/build/contracts/DepositRedemption.json",
    "../../implementation/build/contracts/DepositLiquidation.json",
    "../../implementation/build/contracts/IBTCETHPriceFeed.json",
    "../../implementation/build/contracts/IBondedECDSAKeep.json",
]

jsonFiles.forEach(file => {
    const json = JSON.parse(fs.readFileSync(file, { encoding: 'utf8' }))

    let section = "== `" + json.contractName + "`\n\n"

    for (const signature in json.userdoc.methods) {
        const props = json.userdoc.methods[signature]

        let subsection = "=== `" + signature + "`\n\n"

        let userDocs = json.userdoc.methods[signature]
        let devDocs = json.devdoc.methods[signature]

        if (userDocs.notice) {
            subsection += `${userDocs.notice}\n\n`
        }

        if (devDocs) {
            subsection += "==== Developers\n\n"
            
            if (devDocs.details) {
                subsection += `${devDocs.details}\n\n`
            }

            if (devDocs.params) {
                for (const paramName in devDocs.params) {
                    const paramDoc = devDocs.params[paramName]

                    subsection += `\`${paramName}\`:: ` + paramDoc + "\n"
                }
            }
            if (devDocs.return) {
                subsection += `Returns:: ${devDocs["return"]}`
            }
            
            subsection += "\n\n"
        }

        section += subsection
    }

    doc += section
})

console.log(doc)
