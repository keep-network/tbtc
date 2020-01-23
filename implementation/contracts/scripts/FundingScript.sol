pragma solidity ^0.5.10;

import {TBTCDepositToken} from "../system/TBTCDepositToken.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";
import {VendingMachine} from "../system/VendingMachine.sol";

/// @notice A one-click script for minting TBTC from an unqualified DOT.
/// @dev Wrapper script for VendingMachine.unqualifiedDepositToTbtc.
contract FundingScript {
    TBTCToken tbtcToken;
    VendingMachine vendingMachine;
    TBTCDepositToken tbtcDepositToken;
    FeeRebateToken feeRebateToken;

    constructor(
        address _VendingMachine,
        address _TBTCToken,
        address _TBTCDepositToken,
        address _FeeRebateToken
    ) public {
        vendingMachine = VendingMachine(_VendingMachine);
        tbtcToken = TBTCToken(_TBTCToken);
        tbtcDepositToken = TBTCDepositToken(_TBTCDepositToken);
        feeRebateToken = FeeRebateToken(_FeeRebateToken);
    }

    /// @notice Receives approval for a DOT transfer, and calls `VendingMachine.unqualifiedDepositToTbtc` for a user.
    /// @dev Implements the approveAndCall receiver interface.
    /// @param _from The owner of the token who approved them for transfer.
    /// @param _tokenId Approved DOT for the transfer.
    /// @param _token Token contract address.
    /// @param _extraData Encoded function call to `VendingMachine.unqualifiedDepositToTbtc`.
    function receiveApproval(address _from, uint256 _tokenId, address _token, bytes memory _extraData) public {
        tbtcDepositToken.transferFrom(_from, address(this), _tokenId);
        tbtcDepositToken.approve(address(vendingMachine), _tokenId);

        // Verify _extraData is a call to unqualifiedDepositToTbtc.
        bytes4 sig;
        /* solium-disable security/no-inline-assembly */
        assembly { sig := mload(add(_extraData, 0x20)) }
        require(sig == vendingMachine.unqualifiedDepositToTbtc.selector, "Invalid method signature encoded in _extraData.");

        // Call the VendingMachine.
        // We could explictly encode the call to vending machine, but this would
        // involve manually parsing _extraData and allocating variables.
        (bool success,) = address(vendingMachine).call(
            _extraData
        );
        require(success, "VendingMachine.unqualifiedDepositToTbtc failed");

        // Transfer the TBTC and feeRebateToken to the user.
        tbtcToken.transfer(_from, tbtcToken.balanceOf(address(this)));
        feeRebateToken.transferFrom(address(this), _from, uint256(_tokenId));
    }
}