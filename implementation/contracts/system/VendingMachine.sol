pragma solidity ^0.5.10;

import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import {DepositOwnerToken} from "./DepositOwnerToken.sol";
import {FeeRebateToken} from "./FeeRebateToken.sol";
import {TBTCToken} from "./TBTCToken.sol";
import {TBTCConstants} from "../deposit/TBTCConstants.sol";
import {DepositUtils} from "../deposit/DepositUtils.sol";
import "../deposit/Deposit.sol";

contract VendingMachine {
    using SafeMath for uint256;

    TBTCToken tbtcToken;
    DepositOwnerToken depositOwnerToken;
    FeeRebateToken feeRebateToken;

    constructor(
        address _tbtcToken,
        address _depositOwnerToken,
        address _feeRebateToken
    ) public {
        tbtcToken = TBTCToken(_tbtcToken);
        depositOwnerToken = DepositOwnerToken(_depositOwnerToken);
        feeRebateToken = FeeRebateToken(_feeRebateToken);
    }

    /// @notice Determines whether a deposit is qualified for minting TBTC.
    /// @param _depositAddress the address of the deposit
    function isQualified(address payable _depositAddress) public returns (bool) {
        return Deposit(_depositAddress).inActive();
    }

    /// @notice Pay back the deposit's TBTC and receive the Deposit Owner Token.
    /// @dev    Burns TBTC, transfers DOT from vending machine to caller
    /// @param _dotId ID of Deposit Owner Token to buy
    function tbtcToDot(uint256 _dotId) public {
        require(depositOwnerToken.exists(_dotId), "Deposit Owner Token does not exist");
        require(isQualified(address(_dotId)), "Deposit must be qualified");

        uint256 depositValue =  TBTCConstants.getLotSizeTbtc();
        require(tbtcToken.balanceOf(msg.sender) >= depositValue, "Not enough TBTC for DOT exchange");
        tbtcToken.burnFrom(msg.sender, depositValue);

        // TODO do we need the owner check below? transferFrom can be approved for a user, which might be an interesting use case.
        require(depositOwnerToken.ownerOf(_dotId) == address(this), "Deposit is locked");
        depositOwnerToken.transferFrom(address(this), msg.sender, _dotId);
    }

    /// @notice Trade in the Deposit Owner Token and mint TBTC.
    /// @dev    Transfers DOT from caller to vending machine, and mints TBTC to caller
    /// @param _dotId ID of Deposit Owner Token to sell
    function dotToTbtc(uint256 _dotId) public {
        require(depositOwnerToken.exists(_dotId), "Deposit Owner Token does not exist");
        require(isQualified(address(_dotId)), "Deposit must be qualified");

        depositOwnerToken.transferFrom(msg.sender, address(this), _dotId);

        // If the backing Deposit does not have a signer fee in escrow, mint it.
        Deposit deposit = Deposit(address(uint160(_dotId)));
        uint256 signerFee = deposit.signerFee();
        uint256 depositValue = TBTCConstants.getLotSizeTbtc();

        if(tbtcToken.balanceOf(address(_dotId)) < signerFee) {
            tbtcToken.mint(msg.sender, depositValue.sub(signerFee));
            tbtcToken.mint(address(_dotId), signerFee);
        }
        else{
            tbtcToken.mint(msg.sender, depositValue);
        }

        // owner of the DOT during first TBTC mint receives the FRT
        if(!feeRebateToken.exists(_dotId)){
            feeRebateToken.mint(msg.sender, _dotId);
        }
    }

    // WRAPPERS

    /// @notice Qualifies a deposit and mints TBTC.
    /// @dev User must allow VendingManchine to transfer DOT
    function unqualifiedDepositToTbtc(
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
        Deposit _d = Deposit(_depositAddress);
        require(
            _d.provideBTCFundingProof(
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

        dotToTbtc(uint256(_depositAddress));
    }
}