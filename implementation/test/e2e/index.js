const uniswapTest = require('./UniswapTest')

// eslint-disable-next-line no-only-tests/no-only-tests
const e2eTests = process.env.TEST_E2E == 'true' ? describe.only : describe.skip

e2eTests('E2E tests', function() {
  uniswapTest()
})
