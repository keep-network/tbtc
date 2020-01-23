pragma solidity ^0.5.10;

import {TBTCDepositToken} from "../system/TBTCDepositToken.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";
import {VendingMachine} from "../system/VendingMachine.sol";
import {Deposit} from "../deposit/Deposit.sol";
import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";

/// @notice A one-click script for redeeming TBTC into BTC.
/// @dev Wrapper script for VendingMachine.tbtcToBtc.
contract RedemptionScript {
    using BytesLib for bytes;
    
    TBTCToken tbtcToken;
    VendingMachine vendingMachine;
    FeeRebateToken feeRebateToken;

    constructor(
        address _VendingMachine,
        address _TBTCToken,
        address _FeeRebateToken
    ) public {
        vendingMachine = VendingMachine(_VendingMachine);
        tbtcToken = TBTCToken(_TBTCToken);
        feeRebateToken = FeeRebateToken(_FeeRebateToken);
    }

    /// @notice Receives approval for a TBTC transfer, and calls `VendingMachine.tbtcToBtc` for a user.
    /// @dev Implements the approveAndCall receiver interface.
    /// @param _from The owner of the token who approved them for transfer.
    /// @param _amount Approved TBTC amount for the transfer.
    /// @param _token Token contract address.
    /// @param _extraData Encoded function call to `VendingMachine.tbtcToBtc`.
    function receiveApproval(address _from, uint256 _amount, address _token, bytes memory _extraData) public {
        tbtcToken.transferFrom(_from, address(this), _amount);
        tbtcToken.approve(address(vendingMachine), _amount);

        (bool success, bytes memory returnData) = address(this).call(_extraData);
        require(success, string(returnData));
    }

    /// @notice Redeems a Deposit by purchasing a TDT with TBTC, and using the TDT to redeem corresponding Deposit.
    ///         This function will revert if the Deposit is not in ACTIVE state.
    /// @dev Vending Machine transfers TBTC allowance to Deposit.
    /// @param  _depositAddress     The address of the Deposit to redeem.
    /// @param  _outputValueBytes   The 8-byte Bitcoin transaction output size in Little Endian.
    /// @param  _requesterPKH       The 20-byte Bitcoin pubkeyhash to which to send funds.
    function tbtcToBtc(
        address payable _from,
        address payable _depositAddress,
        bytes8 _outputValueBytes,
        bytes20 _requesterPKH
    ) public{
        address payable requesterAddress = _from;

        vendingMachine.tbtcToTdt(uint256(_depositAddress));

        Deposit _d = Deposit(_depositAddress);

        uint256 tbtcOwed = _d.getRedemptionTbtcRequirement(requesterAddress);

        if(tbtcOwed != 0){
            require(tbtcToken.transferFrom(requesterAddress, address(this), tbtcOwed), "transfer failed");
            tbtcToken.approve(_depositAddress, tbtcOwed);
        }

        _d.requestRedemption(
            _outputValueBytes,
            _requesterPKH,
            requesterAddress
        );

        return;
    }
}