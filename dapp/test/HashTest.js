const Hash = require('../scripts/utils/Hash.js')
const ByteUtils = require('../scripts/utils/ByteUtils.js')
const chai = require('chai')
const assert = chai.assert

describe('Hash', async () => {
  describe('hash256', async () => {
    it('message "abc"', async () => {
      const message = 'abc'
      expectedResult = ByteUtils.fromHex('4f8b42c22dd3729b519ba6f68d2da7cc5b2d606d05daed5ad5128cc03e6c6358')

      result = Hash.hash256(message)

      assert.deepEqual(result, expectedResult)
    })

    it('message "x00"', async () => {
      const message = '\x00'
      expectedResult = ByteUtils.fromHex('1406e05881e299367766d313e26c05564ec91bf721d31726bd6e46e60689539a')

      result = Hash.hash256(message)

      assert.deepEqual(result, expectedResult)
    })

    it('long message', async () => {
      const message = 'The quick brown fox jumps over the lazy dog'
      expectedResult = ByteUtils.fromHex('6d37795021e544d82b41850edf7aabab9a0ebe274e54a519840c4666f35b3937')
      result = Hash.hash256(message)

      assert.deepEqual(result, expectedResult)
    })

    it('bytes buffer message', async () => {
      const message = ByteUtils.fromHex('4e210df8041914be65ec026f2963c3ae79ff867424c40523edb1adc257fde77252846fd232df9ac2952dbdff1981c904abeae46ff4d4fa70bf2767df5bbb5b8b')
      expectedResult = ByteUtils.fromHex('3fbba4ea975db9fec28550e56ec605ba3fbc8e6a1b41c0a94701a183034d2eac')
      result = Hash.hash256(message)

      assert.deepEqual(result, expectedResult)
    })
  })

  describe('sha256', async () => {
    it('empty message', async () => {
      const message = ''
      expectedResult = ByteUtils.fromHex('e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855')

      result = Hash.sha256(message)

      assert.deepEqual(result, expectedResult)
    })

    it('message "abc"', async () => {
      const message = 'abc'
      const expectedResult = ByteUtils.fromHex('ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad')

      result = Hash.sha256(message)

      assert.deepEqual(result, expectedResult)
    })

    it('bytes buffer message', async () => {
      const message = ByteUtils.fromHex('4e210df8041914be65ec026f2963c3ae79ff867424c40523edb1adc257fde77252846fd232df9ac2952dbdff1981c904abeae46ff4d4fa70bf2767df5bbb5b8b')
      expectedResult = ByteUtils.fromHex('1446cbe4d9951074abcf3a6e49368562651072b81a69c7dfbf1ece2fc2cacc04')
      result = Hash.sha256(message)

      assert.deepEqual(result, expectedResult)
    })
  })
})
