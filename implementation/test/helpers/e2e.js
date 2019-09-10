// A collection of E2E tests to run as part of a test suite
export const e2eTests = []

function e2eTest(runner, title, tests) {
  const test = () => runner(title, tests)
  e2eTests.push(test)
}

// eslint-disable-next-line no-only-tests/no-only-tests
/**
 * A helper to structure end-to-end (E2E) integration tests.
 * Tests run with this helper will be added to the integration test suite, and
 * only run during `npm run test:e2e`.
 *
 * The API mirrors `describe`, so usage should look like:
 * ```
 * e2e('Title', function() {
 *  it('behaviour')
 * })
 * ```
 * @param {string} title
 * @param {Function} tests
 * @return {Suite}
 */
const e2e = (title, tests) => e2eTest(contract, title, tests)
e2e.skip = (title, tests) => e2eTest(contract.skip, title, tests)
e2e.only = (title, tests) => e2eTest(contract.only, title, tests)

export { e2e }
