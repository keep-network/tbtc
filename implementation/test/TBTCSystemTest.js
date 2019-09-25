import expectThrow from './helpers/expectThrow'

const TBTCSystem = artifacts.require('TBTCSystem')

contract('TBTCSystem', (accounts) => {
  let tbtcSystem

  before(async () => {
    // Create new TBTCSystem instance where only accounts[0] can mint ERC721 tokens
    // accounts[0] is taking the place of deposit factory address
    tbtcSystem = await TBTCSystem.new(accounts[0])
  })

  describe('mint()', async () => {
    it('correctly mints 721 token with approved caller', async () => {
      const tokenId = 11111
      const mintTo = accounts[1]

      tbtcSystem.mint(mintTo, tokenId)

      const tokenOwner = await tbtcSystem.ownerOf(tokenId).catch((err) => {
        assert.fail(`Token not minted properly: ${err}`)
      })

      assert.equal(mintTo, tokenOwner, 'Token not minted to correct address')
    })

    it('fails to mint 721 token with bad caller', async () => {
      const tokenId = 22222
      const mintTo = accounts[1]

      await expectThrow(
        tbtcSystem.mint(mintTo, tokenId, { from: accounts[1] }),
        'Caller must be depositFactory contract'
      )
    })
  })
})
