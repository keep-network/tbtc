import "./DepositOwnerToken.sol";

contract VendingMachine {
	// The volume of TBTC yet to pass the qualifier.
	uint public unqualifiedVolume = 0;

	// Mapping of deposit to qualification status. 
	mapping(uint256 => bool) qualified;

	// Constants
	uint256 public constant BLOCK_REWARD = 125 * 10^8; // satoshis, TODO var is constant for sake of example.

	modifier onlyDepositOwner(uint256 _depositOwnerTokenId) {
		require(DepositOwnerToken.ownerOf(_depositOwnerTokenId) == msg.sender);
	}

	// After 1 conf, the DOT token can be dispensed
	function dispenseDot(
		uint256 _depositId,
		bytes memory _bitcoinHeaders
	) public {
		require(deposit.inFundingState(), "deposit in wrong state");

		// Verify minimum 1 confirmation
		uint minimumDifficultyFactor = TBTCConstants.getTxProofDifficultyFactor();
		uint _observedDiff = evaluateProofDifficulty(_d, _bitcoinHeaders);
        require(
            _observedDiff >= _reqDiff.mul(minimumDifficultyFactor),
            "Insufficient accumulated difficulty in header chain"
        );

		// Mint the DOT
		DepositOwnerToken.mint(msg.sender, _depositId);
		unqualifiedVolume += deposit.lotSize();

		// Set deposit, active?
		deposit.setActive();
	}
	
	/**
	 * There are two flows here:
	 * 1) qualify, mint tbtc, destroy dot
	 * 2) qualify, keep dot
	 */
	function qualifyDot(
		uint256 _depositOwnerTokenId,
		bytes memory _bitcoinHeaders,
		bool mintTbtc
	) public onlyDepositOwner(_depositOwnerTokenId) {
		require(!qualified[_depositOwnerTokenId], "already qualified");

		// Check qualification
		Deposit deposit = DepositOwnerToken.getDeposit(_depositOwnerTokenId);
		require(isQualified(deposit, _bitcoinHeaders), "deposit doesn't qualify for minting");

		qualified[_depositOwnerTokenId] = true;
		unqualifiedVolume -= deposit.lotSize();

		if(mintTbtc) {
			// Relinquish DOT
			dotToTbtc(_depositOwnerTokenId);
		} else {
			// Keep DOT
		}
	}

	function redeemTbtc(uint256 _depositId) public {
		tbtcToDot();
		deposit.redeem();
	}


	// Repay TBTC drawn and receive DOT token
	function tbtcToDot(
		uint256 _depositId
	) public {
		require(DepositOwnerToken.ownerOf(_depositOwnerTokenId) == address(this), "DOT not available for vending");

		require(tbtc.burn(msg.sender, 1 tbtc));
		DepositOwnerToken.transferFrom(address(this), msg.sender, _depositId);
	}

	// Relinquish DOT token and receive TBTC
	function dotToTbtc(
		uint256 _depositOwnerTokenId,
		bytes memory _bitcoinHeaders
	) internal onlyDepositOwner(_depositOwnerTokenId) {
		require(qualified[_depositOwnerTokenId], "dot not qualified");

		Deposit deposit = DepositOwnerToken.getDeposit(_depositOwnerTokenId);

		require(DepositOwnerToken.transferFrom(msg.sender, _depositOwnerTokenId, this), "no permission to claim");
		mintTbtc(msg.sender, deposit.lotSize());
		mintBeneficiaryToken(msg.sender, deposit);
	}

	// Returns whether a DOT is qualified to be minted into TBTC
	function isQualified(Deposit storage _d, bytes memory _bitcoinHeaders) public view {
		uint minimumDifficultyFactor = TBTCConstants.getTxProofDifficultyFactor();
		// The stopgate is an additional security margin, which limits the volume of TBTC minted per-block,
		// as to disincentivise cross-chain miner arbitrage.
		uint stopgateFactor = unqualifiedVolume / BLOCK_REWARD;

		// TODO modify evaluateProofDifficulty
		uint _observedDiff = evaluateProofDifficulty(_d, _bitcoinHeaders);

		// qualifier = 6 + n
		// where 6 is the minimum number of confs to mint tbtc and
		//       n is the security margin (stopgate) that fluctuates according to opened deposit volume
        require(
            _observedDiff >= _reqDiff.mul(minimumDifficultyFactor + stopgateFactor),
            "Insufficient accumulated difficulty in header chain"
        );
	}
}