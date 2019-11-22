contract('Vending', (accounts) => {

  // check if a deposit is qualified. Will be a
  describe('isQualified', async () => {
    it('Returns true for qualified Token', async () => {
      const TokenId = 1
      // mint a new token with given ID to accounts[1]
      DepositOwnerTokenStub.new(TokenId, accounts[1])
      DepositOwnerTokenStub(tokenId).setQualified(true)
      const qualified = Vending.isQualified(tokenId)
      assert.equal(qualified, true)
    })

    it('Returns false for unqualified Token', async () => {
      const TokenId = 1
      // mint a new token with given ID to accounts[1]
      DepositOwnerTokenStub.new(TokenId, accounts[1])
      const qualified = Vending.isQualified(tokenId)
      assert.equal(qualified, false)
    })

    it('Reverts wtih bad token', async () => {
      const TokenId = 2
      await expectThrow(
        Vending.isQualified(TokenId)
      )
    })
  })

  describe('qualifyDepositOwnerToken', async () => {
    it('Qualifies a Token to be used by the vending machine', async () => {
      const TokenId = 1

      // required proofs for qualifying token. This will not be just one variable
      // and should include merkle proof, block headers etc
      const proofRequirements = '0x00'

      DepositOwnerTokenStub.new(TokenId, accounts[1])
      Vending.qualifyDepositOwnerToken(TokenId, proofRequirements)

      const check = DepositOwnerTokenStub.isQualified(TokenId)
      assert.equal(check, true)
    })


    it('Fails to qualify a token with bad proof', async () => {
      const TokenId = 1

      // required proofs for qualifying token. This will not be just one variable
      // and should include merkle proof, block headers etc
      const proofRequirements = '0xBAD'

      DepositOwnerTokenStub.new(TokenId, accounts[1])
      await expectThrow(
        Vending.qualifyDepositOwnerToken(TokenId, proofRequirements)
      )
    })
	})

  describe('getQualificationRequirements', async () => {
    it('Returns qualification requirements for a given Deposit Owner Token', async () => {

			// artificially set vending machine volume. This assumes a vending stub
      Vending.setVolume(X) //

			const actualRequirements = Vending.getQualificationRequirements()
			
			const expectedRequirements = //Calculate expected requirements here

      assert.equal(expectedRequirements, actualRequirements)
    })
	})
	
	describe('GetTBTC unlocked', async () => {
    it('correctly grants TBTC and handles NFT tokens', async () => {

			const TokenId = 1
			DepositOwnerTokenStub.new(TokenId, accounts[0])
			DepositOwnerTokenStub(tokenId).setQualified(false)

			// -- CHECK PRIOR BALANCES

			Vending.getTBTC(false, TokenId)

			// - CHECK BALANCES AGAIN

			assert(TBTC balance increased correctly)
			assert(DOT is now owned by vending machine)
			assert(accounts[0] no ownes correct Deposit Beneficiary Token)
			assert(DOT is still unredeemed)

    })
	})
	describe('GetTBTC locked', async () => {
    it('correctly grants TBTC and handles NFT tokens', async () => {

			const TokenId = 1
			DepositOwnerTokenStub.new(TokenId, accounts[0])
			DepositOwnerTokenStub(tokenId).setQualified(true)

			// -- CHECK PRIOR BALANCES

			Vending.getTBTC(false, TokenId)

			// - CHECK BALANCES AGAIN

			assert(TBTC balance increased correctly)
			assert(DOT is now owned by accounts[0])
			assert(Deposit Beneficiary Token balance is unchanged)
			assert(DOT is redeemed)

    })
	})
	describe('getDepositOwnerNFT', async () => {
    it('correctly grants Deposit Owner Token', async () => {

			const TokenId = 1
			DepositOwnerTokenStub.new(TokenId, accounts[0])
			DepositOwnerTokenStub.transferTo(TokenId, Vending.address)

			TBTCToken.transferTo(accounts[0], requiredTBTC)
			TBTCToken.approve(Vending, requiredTBTC)
			DepositOwnerTokenStub.approve(Vending)

			Vending.getDepositOwnerNFT(1)

			assert(TBTC was correctly burned)
			assert(accounts[0] now owns Deposit OWner Token with ID 1)
			assert(Deposit Owner Token 1 is unredeemed)

		})
		it('Fails to grants Deposit Owner Token with insufficient funds', async () => {

			const TokenId = 1
			DepositOwnerTokenStub.new(TokenId, Vending.address)
			await expectThrow(
				Vending.getDepositOwnerNFT(1)
				)
    })
	})

	describe('redemptionWrapper', async () => {
    it('correctly initiates Deposit redemption', async () => {

			const TokenId = 1
			DepositOwnerTokenStub.new(TokenId, accounts[0])
			DepositOwnerTokenStub.transferTo(TokenId, Vending.address)

			TBTCToken.transferTo(accounts[0], requiredTBTC)
			TBTCToken.approve(Vending, requiredTBTC)
			DepositOwnerTokenStub.approve(Vending)

			Vending.redemptionWrapper(1)

			assert(TBTC was correctly burned)
			assert(DepositOwnerToken is burned)
			assert(Deposit in redemption state)

		})
  })
})
