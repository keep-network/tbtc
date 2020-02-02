pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "./VendingMachineAuthority.sol";

/**
 @dev Interface of recipient contract for approveAndCall pattern.
*/
interface tokenRecipient { function receiveApproval(address _from, uint256 _value, address _token, bytes calldata _extraData) external; }

contract TBTCToken is ERC20Detailed, ERC20, VendingMachineAuthority {
    /// @dev Constructor, calls ERC20Detailed constructor to set Token info
    ///      ERC20Detailed(TokenName, TokenSymbol, NumberOfDecimals)
    constructor(address _VendingMachine)
        ERC20Detailed("Trustless bitcoin", "TBTC", 18)
        VendingMachineAuthority(_VendingMachine)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev             Mints an amount of the token and assigns it to an account.
    ///                  Uses the internal _mint function
    /// @param _account  The account that will receive the created tokens.
    /// @param _amount   The amount of tokens that will be created.
    // TODO: add back ACL when vending machine refactoring is done.
    function mint(address _account, uint256 _amount) public onlyVendingMachine returns (bool) {
        // NOTE: this is a public function with unchecked minting.
        _mint(_account, _amount);
        return true;
    }

    /// @dev             Burns an amount of the token of a given account
    ///                  deducting from the sender's allowance for said account.
    ///                  Uses the internal _burn function.
    /// @param _account  The account whose tokens will be burnt.
    /// @param _amount   The amount of tokens that will be burnt.
    // TODO: add back ACL when vending machine refactoring is done.
    function burnFrom(address _account, uint256 _amount) public {
        _burnFrom(_account, _amount);
    }

    /// @dev Destroys `amount` tokens from `msg.sender`, reducing the
    /// total supply.
    /// @param _amount   The amount of tokens that will be burnt.
    function burn(uint256 _amount) public {
        _burn(msg.sender, _amount);
    }

    /// @notice           Set allowance for other address and notify.
    ///                   Allows `_spender` to spend no more than `_value` tokens
    ///                   on your behalf and then ping the contract about it.
    /// @dev              The `_spender` should implement the `tokenRecipient` interface above
    ///                   to receive approval notifications.
    /// @param _spender   Address of contract authorized to spend.
    /// @param _value     The max amount they can spend.
    /// @param _extraData Extra information to send to the approved contract.
    function approveAndCall(address _spender, uint256 _value, bytes memory _extraData) public returns (bool success) {
        tokenRecipient spender = tokenRecipient(_spender);
        if (approve(_spender, _value)) {
            spender.receiveApproval(msg.sender, _value, address(this), _extraData);
            return true;
        }
    }
}
