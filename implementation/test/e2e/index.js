import { e2eTests } from '../helpers/e2e'

/*
 * The entire suite is only run when `TEST_E2E` is set.
 * Otherwise it is skipped.
 */
// eslint-disable-next-line no-only-tests/no-only-tests
export const e2eSuite = process.env.TEST_E2E ? describe.only : describe.skip

e2eSuite('Integration tests', function() {
  e2eTests.map((test) => test())
})
