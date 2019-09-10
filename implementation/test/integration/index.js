import { tests } from '../helpers/integration'

/*
 * The entire suite is only run when `INTEGRATION_TEST` is set.
 * Otherwise it is skipped.
 */
// eslint-disable-next-line no-only-tests/no-only-tests
export const integrationSuite = process.env.INTEGRATION_TEST ? describe.only : describe.skip

integrationSuite('Integration tests', function() {
  tests.map((test) => test())
})
