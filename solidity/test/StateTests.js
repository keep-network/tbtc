const states = require('./states/deposit.js')
const { deployAndLinkAll } = require('./helpers/fullDeployer.js')
const { runStatePath } = require('./states/run.js')
const { expect } = require('chai')

describe.only("tBTC states", async function() {
    before(async () => {
        await runStatePath(
            this, // pass the mocha suite to properly set tests up
            states,
            await deployAndLinkAll(),
            "start",
            "awaitingSignerSetup",
        )
    })

    it("should have run all the other tests", async () => {
        expect(true)
    })
})
