const {BN} = require("@openzeppelin/test-helpers")
const {web3} = require("@openzeppelin/test-environment")
const {increaseTime} = require("../helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const Test = require('mocha/lib/test')
const Suite = require('mocha/lib/suite')

async function asyncReduce(array, reducer, initialValue) {
    return array.reduce(
        async (previousValue, nextValue) => {
            const realPrev = await previousValue
            return reducer(realPrev, nextValue)
        },
        initialValue,
    )
}

const resolveAllLogs = (receipt, state) => {
    const contracts =
        Object
            .entries(state)
            .map(([, value]) => value)
            .filter(_ => _.contract && _.address)

    const { resolved: resolvedLogs } = contracts.reduce(
        ({ raw, resolved }, contract) => {
            const events = contract.contract._jsonInterface.filter(_ => _.type === "event")
            const contractLogs = raw.filter(_ => _.address == contract.address)

            const decoded = contractLogs.map(log => {
                const event = events.find(_ => log.topics.includes(_.signature))
                const decoded = web3.eth.abi.decodeLog(
                    event.inputs,
                    log.data,
                    log.topics.slice(1)
                )

                return {
                    ...log,
                    event: event.name,
                    args: decoded,
                }
            })

            return {
                raw: raw.filter(_ => _.address != contract.address),
                resolved: resolved.concat(decoded),
            }
        },
        { raw: receipt.rawLogs, resolved: [] },
    )

    return {
        ...receipt,
        logs: resolvedLogs,
    }
}

const runner = {
    resolveDependencies: async (baseState, dependencies) => {
        const resolved = {}
        for (let [name, resolver] of Object.entries(dependencies)) {
            resolved[name] = await resolver(baseState)
        }

        return {
            ...baseState,
            ...resolved,
        }
    },
    resolveResults: async (initialState, transitionResult) => {
        const receipt =
            await transitionResult.tx
                .then(_ => resolveAllLogs(_.receipt, initialState))
        const resolved = {}
        for (let [property, resolver] of Object.entries(transitionResult)) {
            if (property.startsWith("resolve")) {
                let resolvedProperty = property.replace(/^resolve/, '')
                resolvedProperty =
                    resolvedProperty[0].toLowerCase() + resolvedProperty.substring(1)

                resolved[resolvedProperty] = await resolver(
                    initialState,
                    receipt,
                )
            }
        }

        return {
            ...initialState,
            ...resolved,
        }
    },
    verifyStateTransitions: async (mochaSuite, baseState, stateDefinition) => {
        const resolvedInitialState =
            await runner.resolveDependencies(baseState, stateDefinition.dependencies)

        const stateSuite = Suite.create(
            mochaSuite,
            `should transition from ${stateDefinition.name}`,
        )
        for (let [nextStateName, { transition, after, expect }] of
                Object.entries(stateDefinition.next)) {
            stateSuite.addTest(new Test(
                `to ${nextStateName}`,
                async () => {
                    await createSnapshot()
                    if (after) {
                        await increaseTime((await after(resolvedInitialState)).add(new BN(1)))
                    }
                    const transitionResult = await transition(resolvedInitialState)
                    const receipt = await transitionResult.tx.then(_ => resolveAllLogs(_.receipt, resolvedInitialState))
                    if (expect) {
                        await expect(
                            resolvedInitialState,
                            receipt,
                            await runner.resolveResults(resolvedInitialState, transitionResult),
                        )
                    }
                    await restoreSnapshot()
                }
            ))
        }
    },
    advanceToState: async (baseState, stateDefinition, nextStateName) => {
        const resolvedInitialState =
            await runner.resolveDependencies(baseState, stateDefinition.dependencies)


        const { transition, after } = stateDefinition.next[nextStateName]
        if (after) {
            await increaseTime((await after(resolvedInitialState)).add(new BN(1)))
        }
        const transitionResult = await transition(resolvedInitialState)
        return await runner.resolveResults(resolvedInitialState, transitionResult)
    },
    runStatePath: async (
        mochaSuite,
        stateDefinitions,
        baseState,
        firstStateName,
        ...path
    ) => {
        return asyncReduce(
            path.concat([null]),
            async ({ previousState, definition }, nextStateName) => {
                await runner.verifyStateTransitions(mochaSuite, previousState, definition)

                // Terminating condition.
                if (nextStateName !== null) {
                    return {
                        previousState: await runner.advanceToState(previousState, definition, nextStateName),
                        definition: stateDefinitions[nextStateName],
                    }
                }
            },
            {
                previousState: baseState,
                definition: stateDefinitions[firstStateName],
            },
        )
    }
}

module.exports = runner
