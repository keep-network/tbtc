pragma solidity ^0.5.10;

import "../deposit/Deposit.sol";
import {DepositOwnerToken} from "./DepositOwnerToken.sol";
import {DepositBeneficiaryToken} from "./DepositBeneficiaryToken.sol";
import {DepositFactory} from "../proxy/DepositFactory.sol"; // TODO use TBTCSystem instead
import {TBTCToken} from "../system/TBTCToken.sol";

contract Vending is RedemptionQueue {

	address depositOwnerToken;
	uint256 current;
	uint256 length;

	constructor(address _depositOwnerToken) public{
		depositOwnerToken = _depositOwnerToken;
	}

	// Check if a Deposit Owner Token (DOT) is qualified
	function isQualified(tokenId) public view returns (bool){
		uint256 currentState = Deposit(address(tokenID)).getCurrentState();
        return currentState == 5 || currentState == 6;
    }

	/// @notice                     Fully qualify a deposit.
	///								Anyone may notify the deposit of a funding proof to activate the deposit
    ///                             This is the happy-path of the funding flow. It means that we have succeeded
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding
    /// @param  tokenID             ID of the Deposit Owner Token linked to this deposit
    /// @param _txVersion           Transaction version number (4-byte LE)
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param _txLocktime          Final 4 bytes of the transaction
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first
    /// @return                     True if no errors are thrown
	function qualifyDeposit(
		uint256 tokenID,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public {
		TBTCToken _token = TBTCToken(_d.getTokenAddress);
		DepositOwnerToken _dot = DepositOwnerToken(depositOwnerToken);
        require(
            Deposit(address(tokenID)).provideBTCFundingProof(
                _txVersion,
                _txInputVector,
                _txOutputVector,
                _txLocktime,
                _fundingOutputIndex,
                _merkleProof,
                _txIndexInBlock,
                _bitcoinHeaders),
            "Failed funding proof"
        );

		// Deposit state is now ACTIVE_LOCKED.  DOT(tokenID) is now qualified and unredeemed.
		address rightfulOwner = DepositFactory.getOwner(tokenId);
		if(_dot.ownerOf(tokenID) == address(this)){
			_dot.transferFrom(address(this), rightfulOwner);
		}
	
		// mint singer fee to Deposit
		uint256 toDeposit;
		(, toDeposit) = _d.toMintTBTC();
		_token.mint(address(tokenId), toDeposit);
	}

	
    /// @notice                     Provide work for 1 confimraiton to get a Deposit Owner Token. This does not qualify the deposit. 
    ///                             Deposit remains in AWAITING_BTC_FUNDING_PROOF untill the full proof is satisfied.
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding
    /// @param  tokenID             ID of the Deposit Owner Token linked to this deposit
    /// @param _txVersion           Transaction version number (4-byte LE)
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param _txLocktime          Final 4 bytes of the transaction
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 2 80-byte bitcoin headers
    /// @return                     True if no errors are thrown
	function provideSingleConfirmation(
		uint256 tokenID,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public {
		DepositOwnerToken _dot = DepositOwnerToken(depositOwnerToken);
        require(
            Deposit(address(tokenID)).providePrequalificaticationWork(
                _txVersion,
                _txInputVector,
                _txOutputVector,
                _txLocktime,
                _fundingOutputIndex,
                _merkleProof,
                _txIndexInBlock,
                _bitcoinHeaders),
            "Failed funding proof"
        );

		// Deposit is not activated. DOT is transfered pending furnther proof
		address depositCreator = DepositFactory.getOwner(tokenId);
		if(_dot.ownerOf(tokenID) == address(this)){
			_dot.transferFrom(address(this), depositCreator, tokenID);
		}
	}

	/// @notice Returns a varibale number of confimratons needed to qualify a deposit
	///			based on volume of opened deposits within a given period.
	/// @return number of confirmations needed to qualify a deposit
	function getQualificationRequirements() public view returns (uint256){
		// TODO implement variable block requirement
		return 6;
	}

	/// @notice 		Use a Deposit Owner Token (DOT) to obtain TBTC
	/// @dev 			
	/// @param locked	Boolean, if true, retain exclusive redemption rights,
	///				 	keep control of DOT and receive TBTC. The DOT is marked as redeemed.
	///					Ownership of the DOT now represents the right to redeem the specific custodied UTXO. 
	/// 				If locked bool == false, Forfeit exclusive redemption rights, Swap DOT
	/// 				for Deposit Beneficiary NFT, and receive TBTC.
    /// @param  tokenID ID of the Deposit Owner Token linked to this deposit
	///                 
	function GetTBTC(bool locked, uint256 tokenID) public {
        uint256 toBeneficiary;
		uint256 toDeposit;
		Deposit _d = Deposit(address(tokenID));
		DepositOwnerToken _dot = DepositOwnerToken(depositOwnerToken);
		TBTCToken _token = TBTCToken(_d.getTokenAddress);
		(toBeneficiary, toDeposit) = _d.toMintTBTC();

		require(isQualified(tokenId));
		require(_dot.tokenURI(tokenID) != "REDEEMED");

		if(!locked){
			DepositOwnerToken.transferFrom(msg.sender, address(this), tokenID);

			_dot.setTokenURI(tokenId, "UNREDEEMED");

			//mintTBTC
			_token.mint(msg.sender, toBeneficiary);
			_token.mint(address(tokenId), toDeposit);

			//mint deposit beneficary token
			DepositBeneficiaryToken.mint(msg.sender, tokenId);
			return;
		}
		//mintTBTC
		require(_dot.ownerOf(tokenId) == msg.sender);

		_dot.setTokenURI(tokenId, "REDEEMED");

		_token.mint(msg.sender, toBeneficiary);
	}

	/// @notice 		Pay required TBTC to retrieve given DOT (in unredeemed state)
	/// @param tokenID  ID of the Deposit Owner Token linked to this deposit
	function getDepositOwnerNFT(uint256 tokenId) public returns (address) {
		DepositOwnerToken _dot = DepositOwnerToken(depositOwnerToken);
		TBTCToken _token = TBTCToken(_d.getTokenAddress);

		require(_dot.tokenURI(tokenId) == "UNREDEEMED");
		require(_dot.ownerOf() == address(this));

		uint256 toBeneficiary;
		(toBeneficiary, ) = _d.toMintTBTC();
		_token.burnFrom(msg.sender, toBeneficiary);

		_dot.transferFrom(address(this), depositCreator, tokenID);
		// TODO - allow for non-specific Token to be received
			// decide which NFT to return. oldest/random
			// some market potential here? (note async here will break redemption wrapper)
			// Market can be external system impacting queue order of UTXO redemption.
			// could offer an alternative to signer self redemption by paying for redemption priority.
			// Paying for DOT queue priority would have to somehow lock the DOT to prevent getDepositOwnerNFT()
			// on the specific DOT (since Ownership of the DOT does not guarantee imminent redemption) 
			// only in redemption Wrapper: getDepositOwnerNFT(RedemptionQueue.Next()) -- not async, cool. 
	}

	/// @notice 		redeem a Deposit with TBTC, specific DOT preference.
	/// @param _TokenID ID of the Deposit Owner Token linked to this deposit
	function redemptionWrapper(uint256 _TokenID) public {
		uint256 TokenId;
		if(ownerOf(_TokenId) == msg.sender){
			tokenId = _TokenID;
		}
		else{
			TokenId = getDepositOwnerNFT(_TokenID);
		}
		redeemDepositOwnerToken(TokenId);
	}

	/// @notice 		redeem a Dposit using an unredeemed DOT.
	///					Only TBTC requirement is fees.
	/// @param _TokenID ID of the Deposit Owner Token linked to this deposit.
	function redeemDepositOwnerToken(uint256 TokenID) public {
		DepositOwnerToken _dot = DepositOwnerToken(depositOwnerToken);
		TBTCToken _token = TBTCToken(_d.getTokenAddress);
		uint256 redemptionAmount = 

		require(_dot.ownerOf() == msg.sender);
		require(_dot.tokenURI(tokenId) != "REDEEMED");

		// TODO: handle beneficiary reward for locked deposits cleanly
 		_tbtc.burnFrom(msg.sender,  _dot.signerFee());
        _tbtc.transferFrom(msg.sender, address(tokenId), _dot.beneficiaryReward());

		Deposit(address(tokenId)).requestRedemption(); //temp
	}
}