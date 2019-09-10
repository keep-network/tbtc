
// eslint-disable-next-line no-only-tests/no-only-tests
/**
 * A helper to structure end-to-end (E2E) integration tests.
 * Separates the unit test from the integration test suites.
 */
export const e2e = process.env.TEST_E2E == 'true' ? contract.only : contract.skip
