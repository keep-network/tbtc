// A collection of integration tests to run as part of a test suite
export const tests = []

function integrationTest(runner, title, fn) {
  const test = () => runner(title, fn)
  tests.push(test)
}

// eslint-disable-next-line no-only-tests/no-only-tests
/**
 * A helper to structure integration/E2E tests.
 * Tests run with this helper will be added to the integration test suite, and
 * only run during `npm run integration-test`.
 *
 * The API mirrors `describe`, so usage should look like:
 * ```
 * integration('Title', function() {
 *  it('behaviour')
 * })
 * ```
 * @param {string} title
 * @param {Function} fn
 * @return {Suite}
 */
const integration = (title, fn) => integrationTest(contract, title, fn)
integration.skip = (title, fn) => integrationTest(contract.skip, title, fn)
integration.only = (title, fn) => integrationTest(contract.only, title, fn)

export { integration }
