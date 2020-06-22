// @ts-check
const {BN} = require("@openzeppelin/test-helpers")
const {increaseTime, resolveAllLogs} = require("../helpers/utils.js")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {expect} = require("chai")

/** @typedef { import("mocha/lib/suite") } MochaSuite */
/** @typedef {object} TruffleReceipt */
/**
 * @typedef {object} TruffleTransaction
 * @property {TruffleReceipt} receipt
 */

/**
 * @typedef {(state: UnresolvedState)=>Promise<ResolvedState>} StateDependencyResolver
 * @template UnresolvedState Starting state.
 * @template ResolvedState Starting state with resolved dependencies included.
 */
/**
 * @typedef {(baseState: A, receipt: TruffleReceipt)=>Promise<T>} StateTransitionResolver
 *
 * @template A Dependency-resolved starting state.
 * @template T Starting state with resolved transition variables.
 */
/**
 * @typedef {{[x: string]: StateTransitionResolver<any,any>}} StateTransitionResolvers
 */
/**
 * @typedef {{state: string, tx: Promise<TruffleTransaction>, [x: string]: any}} StateTransitionResult
 */
/**
 * @typedef {(state: A)=>Promise<StateTransitionResult>} StateTransition
 * @template A Starting state.
 */
/**
 * @typedef {(initialState: InitialState, receipt: TruffleReceipt, transitionedState: TransitionedState)=>Promise<any>} StateTransitionExpectation
 * @template InitialState Starting state of the transition.
 * @template TransitionedState Transitioned state after the transition.
 */
/**
 * @typedef {(initialState: InitialState, error: Error)=>Promise<any>} StateTransitionFailureExpectation
 * @template InitialState Starting state of the transition.
 */
/**
 * @typedef {Object<string,StateDependencyResolver<any,any>>} StateDependencies
 */
/**
 * @typedef {object} StateTransitionDefinition
 * @property {(state: object)=>Promise<BN>} after A Promise to a BN delay to
 *           simulate before running a state transition.
 * @property {StateTransition<object>} transition A function that transitions
 *           from this state to another.
 * @property {StateTransitionExpectation<any,any>} expect A function that verifies the
 *           transition from the `transition` property via mocha expectations.
 * @property {StateTransitionFailureExpectation<any>} expectError A function
 *           that verifies an error from a transition from the `transition`
 *           property via mocha expectations.
 */
/**
 * @typedef {Object<string,StateTransitionDefinition>} StateTransitionDefinitions
 */

/**
 * @typedef {object} StateDefinition
 * @property {string} name The name of the state.
 * @property {StateDependencies} dependencies
 * @property {StateTransitionDefinitions} next
 * @property {StateTransitionDefinitions} failNext
 * @template PreviousState
 */

/**
 * StateRunner provides the tools to run a state machine description. The main
 * entry point is runStatePath, which takes a set of state definitions and runs
 * a set path through the state machine, stopping at each state and verifying
 * all of its expectations of possible (next) and impossible (failNext)
 * transitions.
 */
const StateRunner = {
    /**
     * Given a state and an object mapping dependency names to functions that
     * resolve those dependencies, returns an updated state with the resolved
     * values. For example, a base state that is an empty object and a
     * dependencies object that contains `{ name: resolveName, age: resolveAge }`
     * would call `resolveName` and `resolveAge`, and return an object with
     * `name` and `age` properties set to their respective return values.
     *
     * Resolver functions are passed the base state (without any other resolved
     * dependencies). Resolvers can return promises and those promises will be
     * resolved before they are put into the final state. The original state
     * object is left unmodified; a new copy is returned with the resolved
     * properties.
     *
     * @param {InitialState} baseState Any initial state, as long as it is an object.
     * @param {StateDependencies} dependencies A set of dependencies that
     *        maps the property name for the dependency to the function that
     *        resolves the dependency's value. Resolver functions receive the
     *        `baseState`, and the returned value is included in the returned
     *        state as the value of the property with the same name.
     *
     * @return {Promise<object>} A promise to the final resolved state, which is
     *         the `baseState` plus the results of resolving all `dependencies`,
     *         including any returned promises.
     *
     * @template InitialState The starting state.
     */
    resolveDependencies: async (baseState, dependencies) => {
        const resolved = {}
        for (const [name, resolver] of Object.entries(dependencies)) {
            resolved[name] = await resolver(baseState)
        }

        return Object.assign({}, baseState, resolved)
    },
    // We use some fancy TypeScript in our JSDoc types and the ESLint JSDoc
    // syntax validator no likey.
    // eslint-disable-next-line valid-jsdoc
    /**
     * Takes a start state and a transition result, which can be an arbitrary
     * object but must include a `tx` property that contains the Truffle
     * transaction transitioning the contract from one state to the next, and
     * invokes resolver functions in that result to update the start state with
     * new properties.
     *
     * Transition result properties are searched for those that start with
     * `resolve`; these properties are expected to map to functions that will
     * resolve a value for that property in the update state. For example, if
     * the transition result looks like:
     *
     * ```
     * {
     *   state: "test",
     *   otherProperty: 5,
     *   resolveName: () => "hello"
     * }
     * ```
     *
     * The updated state will include a `name` property whose value will be
     * `hello`. Note that the `resolve` part is dropped in the updated state.
     *
     * Resolver functions are passed two things:
     *
     * - The initial state passed to `resolveResults`.
     * - The receipt from the transition result, after resolving all logs in any
     *   contracts included in `initialState` using `resolveAllLogs`.
     *
     * The original state object is left unmodified; a new copy is returned
     * with the resolved properties.
     *
     * @param {InitialState} initialState The state at the start of a state
     *        transition. Should include any contracts that may have emitted
     *        events in the transition.
     * @param {StateTransitionResult & StateTransitionResolvers} transitionResult
     *        The result of running a state transition transaction.
     *
     * @return {Promise<object>} A promise to an object that is the
     *         `initialState` supplemented by the resolved properties from the
     *         `transitionResult`. Resolved properties lose their `resolve`
     *         prefix, so e.g. the result of `resolveDeposit` is placed in a
     *         `deposit` property in the returned state.
     *
     * @template InitialState The starting state.
     */
    resolveResults: async (initialState, transitionResult) => {
        const receipt =
            await transitionResult.tx
                .then(_ => resolveAllLogs(_.receipt, initialState))
        const resolved = {}
        for (const [property, resolver] of Object.entries(transitionResult)) {
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

        return Object.assign({}, initialState, resolved)
    },
    /**
     * Given a mocha suite, a base state, and the definition for that state's
     * future transitions,
     *
     * @param {Promise<object>} baseStatePromise Promise to the base state for
     *        these tessts, after any prior state transitions.
     * @param {StateDefinition<BaseState>} stateDefinition
     *
     * @return {Promise<void>} A promise to the completion of the tests set up
     *         by this function. It will be resolved after tests have finished
     *         running, and allows the chaining of future state tests.
     *
     * @template BaseState The current state.
     */
    setUpStateTests: (baseStatePromise, stateDefinition) => {
        // These are used to resolve promises in the async before hook. This
        // ensures this test's setup runs after any previous tests AND state
        // transitions have completed.
        let baseState
        let resolvedInitialState

        // We resolve these promises below to indicate to the caller that all
        // tests here have completed. Here, capture the resolver functions for
        // that purpose.
        /** @type {()=>void} */
        let positiveTestsCompleted
        /** @type {()=>void} */
        let negativeTestsCompleted
        const positiveTestsCompletePromise = new Promise((resolve) => {
            positiveTestsCompleted = resolve
        })
        const negativeTestsCompletePromise = new Promise((resolve) => {
            negativeTestsCompleted = resolve
        })

        const expectedSuccesses = Object.entries(stateDefinition.next || {})
        if (expectedSuccesses.length) {
            describe(`should transition from ${stateDefinition.name}`, () => {
                before(async () => {
                    baseState = await baseStatePromise
                    resolvedInitialState =
                        await StateRunner.resolveDependencies(
                            baseState,
                            stateDefinition.dependencies,
                        )
                })
                beforeEach(async () => { await createSnapshot() })
                afterEach(async () => { await restoreSnapshot() })
                // After all tests have run, resolve the test promise so that the
                // next state transition can take place.
                after(() => positiveTestsCompleted())

                for (const [nextStateName, { transition, after, expect }] of
                        Object.entries(stateDefinition.next)) {
                    it(`to ${nextStateName}`, async () => {
                        if (after) {
                            await increaseTime((await after(resolvedInitialState)).add(new BN(1)))
                        }
                        const transitionResult = await transition(resolvedInitialState)
                        const receipt = await transitionResult.tx.then(_ => resolveAllLogs(_.receipt, resolvedInitialState))

                        if (expect) {
                            await expect(
                                resolvedInitialState,
                                receipt,
                                await StateRunner.resolveResults(resolvedInitialState, transitionResult),
                            )
                        }
                    })
                }
            })
        } else {
            positiveTestsCompleted()
        }

        const expectedFailures = Object.entries(stateDefinition.failNext || {})
        if (expectedFailures.length > 0) {
            describe(`should NOT transition from ${stateDefinition.name}`, () => {
                before(async () => {
                    // Wait for positive tests to complete so our snapshots don't
                    // walk all over each other.
                    // await positiveTestsCompletePromise
                    baseState = await baseStatePromise
                    resolvedInitialState =
                        await StateRunner.resolveDependencies(
                            baseState,
                            stateDefinition.dependencies,
                        )
                })
                beforeEach(async () => { await createSnapshot() })
                afterEach(async () => { await restoreSnapshot() })
                // After all tests have run, resolve the test promise so that the
                // next state transition can take place.
                after(() => negativeTestsCompleted())

                for (const [nextStateName, { transition, after, expectError }] of expectedFailures) {
                    it(`to ${nextStateName}`, async () => {
                        try {
                            if (after) {
                                await increaseTime((await after(resolvedInitialState)).add(new BN(1)))
                            }
                            const transitionResult = await transition(resolvedInitialState)

                            const tx = await transitionResult.tx

                            expect.fail(
                                "Transition should have failed but succeeded " +
                                `with events: ${
                                    JSON.stringify(
                                        resolveAllLogs(
                                            tx.receipt,
                                            resolvedInitialState
                                        ).logs.map(({ event, args }) => ({ event, args }))
                                    )
                                }`)
                        } catch (e) {
                            if (e.message.match(/Transition should have failed\./)) {
                                throw e
                            } else {
                                await expectError(resolvedInitialState, e)
                            }
                        }
                    })
                }
            })
        } else {
            negativeTestsCompleted()
        }

        // Settle when the negative tests promise does; note that negative tests
        // will only run after positive tests do.
        return positiveTestsCompletePromise.then(() => negativeTestsCompletePromise)
    },
    /**
     * Advances from the given `baseState` to the given `nextStateName`, using
     * `stateDefinition` to resolve dependencies and find the transition to the
     * next state.
     *
     * @param {Promise<object>} baseStatePromise
     * @param {StateDefinition<object>} stateDefinition
     * @param {string} nextStateName
     *
     * @return {Promise<object>} The state produced by the transition to
     *         `nextStateName`.
     */
    advanceToState: async (baseStatePromise, stateDefinition, nextStateName) => {
        const baseState = await baseStatePromise
        const resolvedInitialState =
            await StateRunner.resolveDependencies(baseState, stateDefinition.dependencies)

        const { transition, after } = stateDefinition.next[nextStateName]
        if (after) {
            await increaseTime((await after(resolvedInitialState)).add(new BN(1)))
        }
        const transitionResult = await transition(resolvedInitialState)

        const resolved = await StateRunner.resolveResults(resolvedInitialState, transitionResult)
        return resolved
    },
    /**
     * Runs a state path.
     *
     * @param {Object.<string,StateDefinition<object>>} stateDefinitions
     * @param {object} baseState
     * @param {string} firstStateName
     * @param {...string} path
     *
     * @return {Promise<void>} A promise to the full completion of all tests
     *         in the state path.
     */
    runStatePath: (
        stateDefinitions,
        baseState,
        firstStateName,
        ...path
    ) => {
        return path.concat([null]).reduce(
            ({ currentState, definition }, nextStateName) => {
                const stateTestsPromise = StateRunner.setUpStateTests(
                    currentState,
                    definition,
                )

                // Terminating condition.
                if (nextStateName !== null) {
                    return {
                        currentState: stateTestsPromise.then(() => StateRunner.advanceToState(
                            currentState,
                            definition,
                            nextStateName
                        )),
                        definition: stateDefinitions[nextStateName],
                    }
                } else {
                    return { currentState, definition }
                }
            },
            {
                currentState: Promise.resolve(baseState),
                definition: stateDefinitions[firstStateName],
            },
        ).currentState
    }
}

module.exports = StateRunner
