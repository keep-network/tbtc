pragma solidity 0.5.17;

import {TBTCDepositToken} from "../system/TBTCDepositToken.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";
import {VendingMachine} from "../system/VendingMachine.sol";

/// @notice A one-click script for minting TBTC from an unqualified TDT.
/// @dev Wrapper script for VendingMachine.unqualifiedDepositToTbtc
/// This contract implements receiveApproval() and can therefore use
/// approveAndCall(). This pattern combines TBTC Token approval and
/// vendingMachine.unqualifiedDepositToTbtc() in a single transaction.
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

    /// @notice Receives approval for a TDT transfer, and calls `VendingMachine.unqualifiedDepositToTbtc` for a user.
    /// @dev Implements the approveAndCall receiver interface.
    /// @param _from The owner of the token who approved them for transfer.
    /// @param _tokenId Approved TDT for the transfer.
    /// @param _token Token contract address.
    /// @param _extraData Encoded function call to `VendingMachine.unqualifiedDepositToTbtc`.
    function receiveApproval(address _from, uint256 _tokenId, address _token, bytes memory _extraData) public {
        tbtcDepositToken.transferFrom(_from, address(this), _tokenId);
        tbtcDepositToken.approve(address(vendingMachine), _tokenId);

        // Verify _extraData is a call to unqualifiedDepositToTbtc.
        bytes4 functionSignature;
        assembly {
            functionSignature := and(mload(add(_extraData, 0x20)), not(0xff))
        }
        require(
            functionSignature == vendingMachine.unqualifiedDepositToTbtc.selector,
            "Bad _extraData signature. Call must be to unqualifiedDepositToTbtc."
        );

        // Call the VendingMachine.
        // We could explictly encode the call to vending machine, but this would
        // involve manually parsing _extraData and allocating variables.
        (bool success, bytes memory returnData) = address(vendingMachine).call(
            // solium-disable-previous-line security/no-low-level-calls
            _extraData
        );
        require(success, string(returnData));

        // Transfer the TBTC and feeRebateToken to the user.
        tbtcToken.transfer(_from, tbtcToken.balanceOf(address(this)));
        feeRebateToken.transferFrom(address(this), _from, uint256(_tokenId));
    }
}
