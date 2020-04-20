pragma solidity 0.5.17;

import {TBTCDepositToken} from "../system/TBTCDepositToken.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";
import {VendingMachine} from "../system/VendingMachine.sol";
import {Deposit} from "../deposit/Deposit.sol";
import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";

/// @notice A one-click script for redeeming TBTC into BTC.
/// @dev Wrapper script for VendingMachine.tbtcToBtc
/// This contract implements receiveApproval() and can therefore use
/// approveAndCall(). This pattern combines TBTC Token approval and
/// vendingMachine.tbtcToBtc() in a single transaction.
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
    /// @param _extraData Encoded function call to `VendingMachine.tbtcToBtc`.
    function receiveApproval(address _from, uint256 _amount, address, bytes memory _extraData) public {
        tbtcToken.transferFrom(_from, address(this), _amount);
        tbtcToken.approve(address(vendingMachine), _amount);

        // Verify _extraData is a call to tbtcToBtc.
        bytes4 functionSignature;
        assembly { functionSignature := mload(add(_extraData, 0x20)) }
        require(
            functionSignature == vendingMachine.tbtcToBtc.selector,
            "Bad _extraData signature. Call must be to tbtcToBtc."
        );

        (bool success, bytes memory returnData) = address(vendingMachine).call(_extraData);
        // solium-disable-previous-line security/no-low-level-calls
        // By default, `address.call`  will catch any revert messages.
        // Converting the `returnData` to a string will effectively forward any revert messages.
        // https://ethereum.stackexchange.com/questions/69133/forward-revert-message-from-low-level-solidity-call
        require(success, string(returnData));
    }
}
