pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import {DepositOwnerToken} from "./DepositOwnerToken.sol";
import {TBTCToken} from "./TBTCToken.sol";
import {TBTCConstants} from "../deposit/TBTCConstants.sol";
import {Deposit} from "../deposit/Deposit.sol";
import {DepositUtils} from "../deposit/DepositUtils.sol";

contract VendingMachine {
    using SafeMath for uint256;
    
    TBTCToken tbtcToken;
    DepositOwnerToken depositOwnerToken;

    constructor(
        address _tbtcToken,
        address _depositOwnerToken
    ) public {
        tbtcToken = TBTCToken(_tbtcToken);
        depositOwnerToken = DepositOwnerToken(_depositOwnerToken);
    }

    /// @notice Qualifies a deposit for minting TBTC.
    /// @dev Calls out to Deposit to verify funding proof, mints the TBTC signer fee.
    function qualifyDeposit(
        address payable _depositAddress,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public {
        // require(!isQualified(_depositId), "Deposit already qualified");
        require(
            Deposit(_depositAddress).provideBTCFundingProof(
                _txVersion,
                _txInputVector,
                _txOutputVector,
                _txLocktime,
                _fundingOutputIndex,
                _merkleProof,
                _txIndexInBlock,
                _bitcoinHeaders
            ),
            "failed to provide funding proof");

        tbtcToken.mint(_depositAddress, DepositUtils.signerFee());

        // TODO mark deposit as qualified
    }

    /// @notice Determines whether a deposit is qualified for minting TBTC.
    /// @param _depositAddress the address of the deposit
    function isQualified(address _depositAddress) public returns (bool) {
        // TODO
        // This is stubbed out for prototyping, separate to the actual qualification logic.
        // However we might remove it later.
        return true;
    }

    /// @notice Pay back the deposit's TBTC and receive the Deposit Owner Token.
    /// @dev    Burns TBTC, transfers DOT from vending machine to caller
    /// @param _dotId ID of Deposit Owner Token to buy
    function tbtcToDot(uint256 _dotId) public {
        require(isQualified(address(_dotId)), "Deposit must be qualified");

        require(tbtcToken.balanceOf(msg.sender) >= getDepositValueLessSignerFee(), "Not enough TBTC for DOT exchange");
        tbtcToken.burnFrom(msg.sender, getDepositValueLessSignerFee());

        // TODO do we need the owner check below? transferFrom can be approved for a user, which might be an interesting use case.
        require(depositOwnerToken.ownerOf(_dotId) == address(this), "Deposit is locked");
        depositOwnerToken.transferFrom(address(this), msg.sender, _dotId);
    }

    /// @notice Trade in the Deposit Owner Token and mint TBTC.
    /// @dev    Transfers DOT from caller to vending machine, and mints TBTC to caller
    /// @param _dotId ID of Deposit Owner Token to sell
    function dotToTbtc(uint256 _dotId) public {
        require(isQualified(address(_dotId)), "Deposit must be qualified");

        depositOwnerToken.transferFrom(msg.sender, address(this), _dotId);
        tbtcToken.mint(msg.sender, getDepositValueLessSignerFee());
    }

    // TODO temporary helper function
    /// @notice Gets the lot size less signer fees
    function getDepositValueLessSignerFee() internal returns (uint) {
        uint256 _multiplier = TBTCConstants.getSatoshiMultiplier();
        uint256 _signerFee = DepositUtils.signerFee();
        uint256 _totalValue = TBTCConstants.getLotSize().mul(_multiplier);
        return _totalValue.sub(_signerFee);
    }
}