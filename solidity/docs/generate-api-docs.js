const fs = require('fs')

let doc = `:toc: true
:toclevels: 2

= tBTC API Documentation

Welcome to the tBTC API Documentation. The primary contracts involved in tBTC
are listed below, along with their public methods.

toc::[]
`

const jsonFiles = [
  '../../solidity/build/contracts/TBTCSystem.json',
  // "../../solidity/build/contracts/TBTCConstants.json",
  '../../solidity/build/contracts/Deposit.json',
  '../../solidity/build/contracts/DepositStates.json',
  '../../solidity/build/contracts/DepositFunding.json',
  '../../solidity/build/contracts/DepositRedemption.json',
  '../../solidity/build/contracts/DepositLiquidation.json',
  '../../solidity/build/contracts/ISatWeiPriceFeed.json',
  '../../solidity/build/contracts/IBondedECDSAKeep.json',
  '../../solidity/build/contracts/ITokenRecipient.json',
]

jsonFiles.forEach((file) => {
  const json = JSON.parse(fs.readFileSync(file, { encoding: 'utf8' }))

  let section = '== `' + json.contractName + '`\n\n'

  for (const signature in json.userdoc.methods) {
    const props = json.userdoc.methods[signature]

    let subsection = '=== `' + signature + '`\n\n'

    const userDocs = json.userdoc.methods[signature]
    const devDocs = json.devdoc.methods[signature]

    if (userDocs.notice) {
      subsection += `${userDocs.notice}\n\n`
    }

    if (devDocs) {
      subsection += '==== Developers\n\n'

      if (devDocs.details) {
        subsection += `${devDocs.details}\n\n`
      }

      if (devDocs.params) {
        for (const paramName in devDocs.params) {
          const paramDoc = devDocs.params[paramName]

          subsection += `\`${paramName}\`:: ` + paramDoc + '\n'
        }
      }
      if (devDocs.return) {
        subsection += `Returns:: ${devDocs['return']}`
      }

      subsection += '\n\n'
    }

    section += subsection
  }

  doc += section
})

console.log(doc)
