const FundingProof = require('./FundingProof')

contract('FundingProof', (accounts) => {
  describe('parseTransaction()', async () => {
    it('returns details of a transaction with a flag', async () => {
      const hex = '02000000000101f767a8374661b93c2d4b06b0462e400d227f88550644fbeeaf4e97db79577c9300000000000f965080012fbf0300000000002200201ac5f59fe0afb9bbc9a555dd7fbe45a8a9a91aad968ad5f234a7984e3cd8cda504004730440220456879167d9ad65ad1deeef6715bd59a89cd5afa79ee85ede94f78af12819cd40220555cadb02827d73b0394adb4a04efb74b5e15dd604b443307ef5c971b3a3847701483045022100b1b406c837e52a59c05d03ab1aeefc7e11ee725804094c5a3eadefdc8e7e679b02202d918469b0a868e975b956b50390c613d8ce6ad5dfea402f338aab3413e00a8901475221023e42ab8dc26b6126bf70317aad082b2ee3a28207df125a65708d316f3f85e70b210356038ffd47f875faaddd28ddcb0f84f91c4eeb86eaace0241df281b462c034d852ae14a3f620'

      const tx = Buffer.from(hex, 'hex')

      const result = FundingProof.parseTransaction(tx)

      const expectedResult = {
        version: Buffer.from('02000000', 'hex'),
        txInVector: Buffer.from('01f767a8374661b93c2d4b06b0462e400d227f88550644fbeeaf4e97db79577c9300000000000f965080', 'hex'),
        txOutVector: Buffer.from('012fbf0300000000002200201ac5f59fe0afb9bbc9a555dd7fbe45a8a9a91aad968ad5f234a7984e3cd8cda5', 'hex'),
        locktime: Buffer.from('14a3f620', 'hex'),
        fundingOutputIndex: 0,
      }

      assert.deepEqual(
        result,
        expectedResult,
        'unexpected result',
      )
    })

    it('returns details of a transaction with no flag', async () => {
      const hex = '0100000001aea2e43133a533669b942d335ef6ebef7528f01c3ed1d43b4ccff1e9590d44c9010000006a4730440220785a31ce8bf2c63c5fbda079dea98f2740eaa81dfd09d6987b7ba9a4d2a5ccb702204ef4ff2f852a25fb4f75c2b16d61c09fcd411eea962f51e2ceec630de6e3cc8f0121028896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2effffffff022823000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4f0840d00000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac00000000'

      const tx = Buffer.from(hex, 'hex')

      const result = FundingProof.parseTransaction(tx)

      const expectedResult = {
        version: Buffer.from('01000000', 'hex'),
        txInVector: Buffer.from('01aea2e43133a533669b942d335ef6ebef7528f01c3ed1d43b4ccff1e9590d44c9010000006a4730440220785a31ce8bf2c63c5fbda079dea98f2740eaa81dfd09d6987b7ba9a4d2a5ccb702204ef4ff2f852a25fb4f75c2b16d61c09fcd411eea962f51e2ceec630de6e3cc8f0121028896955d043b5a43957b21901f2cce9f0bfb484531b03ad6cd3153e45e73ee2effffffff', 'hex'),
        txOutVector: Buffer.from('022823000000000000160014d849b1e1cede2ac7d7188cf8700e97d6975c91c4f0840d00000000001976a914d849b1e1cede2ac7d7188cf8700e97d6975c91c488ac', 'hex'),
        locktime: Buffer.from('00000000', 'hex'),
        fundingOutputIndex: 0,
      }

      assert.deepEqual(
        result,
        expectedResult,
        'unexpected result',
      )
    })
  })
})
