pragma solidity ^0.5.10;

import Deposit, TBTCToken, DepositOwnerNFT, DepositBEenficiaryNFT;

contract Vending {

	// Check if a Deposit Owner Token (DOT) is qualified
	function isQualified(tokenId) public view returns (bool);

	// Qualify a DOT
	function qualifyDepositOwnerToken(tokenID, proofRequirements) public view returns (bool){
		require(tokenId has not been redeemed);
		require(within timeout) // expand on timeout in this context
		//proof checks out
		//proofRequirements <> getQualificationRequirements() 
	}

	// get qualification requirements, currently just number of
	// confirmations required. This is 6 + X. X depends on volume
	function getQualificationRequirements() public view returns (uint256){
		// return current number of blocks needed for qualification
	}

	// Use a DOT to obtain TBTC (ERC20)
	// If locked bool == true, maintain exclusive redemption rights,
	// keep control of DOT and receive TBTC. The DOT is marked as redeemed. Ownership of the DOT
	// now represents the right to redeem the specific custodied UTXO. 
	// If locked bool == false, Forfeit exclusive redemption rights, Swap DOT
	// for Deposit Beneficiary NFT, and receive TBTC.
	// DOTs controlled by Vending machine do not need to be marked as redeemed due to 
	// the TBTC burning requirement enforced by getDepositOwnerNFT()
	// this function can be separated into getTBTCLocked() and getTBTCUnlocked()
	function GetTBTC(bool locked, uint256 tokenID, proofRequirements ) public {
		require(isQualified(tokenId, proofRequirements))
		require(within timeout) // expand on timeout in this context
		if(!locked){
			// swap Deposit Owner NFT for TBTC
			// supply Deposit beneficiary NFT
			// return
		}
		// supply TBTC, but no deposit beneficiary NFT
	}

	// pay required TBTC to retrieve given DOT (in unredeemed state)
	function getDepositOwnerNFT(uint256 tokenid) public returns (address) {
		require(sended approved sufficient TBTC to be burned);
		if(tokenId == 0){
			// decide which NFT to return. oldest/random
			// some market potential here? (note async here will break redemption wrapper)
			// Market can be external system impacting queue order of UTXO redemption.
			// could offer an alternative to signer self redemption by paying for redemption priority.
			// Paying for DOT queue priority would have to somehow lock the DOT to prevent getDepositOwnerNFT()
			// on the specific DOT (since Ownership of the DOT does not guarantee imminent redemption) 
			// only in redemption Wrapper: getDepositOwnerNFT(RedemptionQueue.Next()) -- not async, cool. 
		}
		// check vending machine for given NFT, revert if it doesn't exist.
		// transfer Deposit owner NFT
		// return tokenID
	}

	// pay required TBTC to retrieve non-specific DOT (in unredeemed state)
	function getDepositOwnerNFT() public returns (address) {
			getDepositOwnerNFT(0);
	}

	// redeem a Deposit with TBTC, specific DOT preference.
	function redemptionWrapper(uint256 _TokenID) public {
		uint256 TokenId;
		if(ownerOf(_TokenId) == msg.sender){
			tokenId = _TokenID;
		}
		else{
			TokenId = getDepositOwnerNFT(_TokenID);
		}
		redeemDepositOwnerToken(TokenId)
	}

	// redeem a Deposit using TBTC. No Specific DOT preference
	function redemptionWrapper() public{
		uint256 TokenId = getDepositOwnerNFT(0);
		redeemDepositOwnerToken(TokenId)
	}

	// redeem a Dposit using an unredeemed DOT. not TBTC requirement
	function redeemDepositOwnerToken(uint256 TokenID) public {
		require(msg.sender == tokenId owner);
		require(TokenID is valid and unredeemed)
		Deposit(address(tokenId)).requestRedemption() //temp
			// msg.sender issues here? 
			// will need an approval chain for TBTC burn approval requirements
			// might need restructure on Deposit to send
			// beneficiary reward to correct msg.sender instead of address(this) pass var?
	}
}