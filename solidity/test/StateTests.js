const states = require('./states/deposit.js')
const { deployAndLinkAll } = require('./helpers/fullDeployer.js')
const { runStatePath } = require('./states/run.js')

describe.only("tBTC states", async () => {
    it("when verifying paths", async () => {
        const setup = async () => {
            const deployed = await deployAndLinkAll()

            return deployed
        }

        await runStatePath(
            states,
            await setup(),
            "start",
            "awaitingSignerSetup",
        )
    })
})
