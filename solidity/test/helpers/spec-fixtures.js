// @ts-check
const AsciidoctorCore = require('@asciidoctor/core')
const asciidoctor = AsciidoctorCore.default()
const path = require('path')
const BN = require('bn.js')

/**
 * @typedef {Object} DListItems
 * @property {()=>[AsciidoctorCore.Asciidoctor.ListItem[],AsciidoctorCore.Asciidoctor.ListItem?][]} getItems
 */

 const docsPath = path.join(__dirname, "..", "..", "..", "docs")

const redemptionPaymentAppendixPath = path.join(
    docsPath,
    "appendix",
    "disbursal",
    "index.adoc",
)
const redemptionPaymentTable =
    /** @type { AsciidoctorCore.Asciidoctor.Table } */
    (asciidoctor
        .loadFile(redemptionPaymentAppendixPath, { base_dir: docsPath })
        .findBy({ context: 'table' })[0])

const redemptionPaymentColumns = redemptionPaymentTable.getHeadRows()[0].map(_ => _.getText())

/**
 * @typedef {number?} AccountIndex A 0-based index that can be used to index
 *          into a testing `accounts` array. Undefined or null if the field of
 *          this type has no associated account.
 */

/**
 * @typedef {object} RedemptionPaymentFixture
 * @property {number} row The row in the disbursal table this fixture represents.
 * @property {string} description A human-readable description of this fixture's
 *           system state.
 * @property {boolean} preTerm Whether this scenario is for a pre-term deposit.
 * @property {boolean} courtesyCall Whether this scenario is for a courtesy-call
 *           deposit.
 * @property {boolean} [redemptionShouldRevert] If true, this redemption should
 *           revert rather than complete successfully.
 * @property {AccountIndex} [tdtHolder] The TDT holder account index.
 * @property {AccountIndex} [frtHolder] The FRT holder account index.
 * @property {AccountIndex} [redeemer] The redeemer account index.
 * @property {BN} repaymentAmount The repayment amount that the deposit should
 *           require at requestRedemption time from the redeemer, in TBTC *
 *           10**18.
 * @property {Object.<string,BN>} disbursalAmounts The amount due
 *           to each of a set of parties. Where a party can be converted to an
 *           `AccountIndex`, it is. Signers are referenced as `signers`,
 *           however, since they don't correspond to a single account.
 */

/**
 * Returns a set of redemption payment fixtures for the specified lot size and
 * signer fee divisor. An aray of objects is returned that represent various
 * repayment scenarios, pulled directly from the tBTC specification.
 *
 * @param {BN} lotSize The lot size in TBTC * 10**18.
 * @param {BN} signerFee The signer fee for the deposit in TBTC * 10** 18.
 * @return {RedemptionPaymentFixture[]} An array of redemption payment
 *         fixtures.
 */
function redemptionPaymentFixturesFor(lotSize, signerFee) {
    /**
     * Function to convert a party's name in the disbursal table to an account
     * index. Returns the name if it can't be converted.
     * @param {string} partyName The name of the party from the disbursal table;
     *        typically a letter, a `-` to indicate no party, or `signers` for
     *        the signers.
     * @return {string | AccountIndex | undefined} If the party name is a
     *         letter, an AccountIndex for that letter. If there is no party,
     *         undefined. Otherwise, the party name is returned, lower-cased.
     */
    function asAccountIndex(partyName) {
        switch(partyName) {
            case "A": return 0
            case "B": return 1
            case "C": return 2
            case "-": return undefined
            default: return partyName.toLowerCase()
        }
    }

    /**
     * Looks up the given value variable as a factor of lot size and signer fee.
     * @param {string} valueVariable The asciidoctor variable reference that
     *        should be resolved to a specific value.
     * @return {BN} The value for the provided variable, or 0 if the variable
     *         has no known value.
     */
    function lookupValue(valueVariable) {
        const variableMatch = valueVariable.match(/{[^}]+}/) || []
        switch(variableMatch[0]) {
            case "{signer-fee}": return signerFee
            case "{tbtc-lot-size}": return lotSize
            case "{tbtc-lot-size-less-signer-fee}": return lotSize.sub(signerFee)
            case "{btc-lot-size}": return new BN(0) // Nothing due inside the system.
            default: return new BN(0)
        }
    }

    return redemptionPaymentTable.getBodyRows().flatMap((row, index) => {
            /** @type {RedemptionPaymentFixture} */
            const fixture = {
                row: index + 1,
                description: '',
                preTerm: true,
                courtesyCall: false,
                repaymentAmount: new BN(0),
                disbursalAmounts: {},
            }

            for (let i = 0; i < redemptionPaymentColumns.length; ++i) {
                const columnName = redemptionPaymentColumns[i]
                const columnText = row[i].getText()

                // Below, breaking breaks the loop, not the switch; use continue
                // instead.
                switch(columnName) {
                case "Deposit state":
                    fixture.preTerm = columnText == "Pre-term"
                    continue
                case "Repayment Amount":
                    fixture.redemptionShouldRevert = !!columnText.match(/N\/A/)
                    fixture.repaymentAmount = lookupValue(columnText)
                    continue
                case "Disbursal Amounts":
                    fixture.redemptionShouldRevert = !!columnText.match(/N\/A/)

                    if (!fixture.redemptionShouldRevert) {
                        const disbursalList =
                            /** @type {DListItems} */ (
                                /** @type {unknown} trust me */
                            (row[i].getInnerDocument().findBy({ context: 'dlist' })[0]))
                        disbursalList.getItems().forEach(([recipient, amount]) => {
                            // Definition lists can have multiple terms per
                            // definition, though this list doesn't do that.
                            fixture.disbursalAmounts[
                                asAccountIndex(recipient[0].getText())
                            ] = lookupValue(amount.getText())
                        })
                    }
                    continue
                default:
                    const camelCaseName =
                        columnName
                            .replace(/^[A-Z]+/, (_) => _.toLowerCase())
                            .replace(/ [a-z]/g, (_) => _.toUpperCase())
                            .replace(/ /g, '')

                    fixture.description += ` ${columnName}: ${columnText}`

                    fixture[camelCaseName] = asAccountIndex(columnText)
                }
            }

            fixture.description =
                (fixture.preTerm
                    ? "Pre-term"
                    : "At-term ") + fixture.description

            if (! fixture.preTerm) {
                // For at-term fixtures, create a copy that is the same fixture
                // setup, but at courtesy call instead. This is because
                // courtesy-call scenarios are identical to at-term scenarios
                // for disbursal purposes, but their starting state is
               // different.
                const courtesyCallCopy =
                    Object.assign(
                        {
                            courtesyCall: true,
                            description: fixture.description.replace(
                                "At-term ",
                                "Courtesy",
                            ),
                        },
                        fixture
                    )

                return [fixture, courtesyCallCopy]
            } else {
                return fixture
            }
        },
        [],
    )
}

module.exports = {
    redemptionPaymentFixturesFor,
}
