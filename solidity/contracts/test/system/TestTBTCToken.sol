pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../../contracts/system/TBTCToken.sol";

contract TestTBTCToken is TBTCToken{

    constructor(address _vendingMachine)
        TBTCToken(_vendingMachine)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev             We can't call TBTCToken mint function from deposit Test becuase of ACL.
    ///                  This function bypasses ACL and can be called in Deposit tests
    ///                  Mints an amount of the token and assigns it to an account.
    ///                  Uses the internal _mint function.
    /// @param _account  The account that will receive the created tokens.
    /// @param _amount   The amount of tokens that will be created.
    function forceMint(address _account, uint256 _amount) public returns (bool){
        // NOTE: this is a public function with unchecked minting.
        _mint(_account, _amount);
        return true;
    }

    /// @dev             We can't call TBTCToken burn function from deposit Test becuase of ACL.
    ///                  This function bypasses ACL and can be called in Deposit tests.
    ///                  Burns an amount of the token of a given account
    ///                  deducting from the sender's allowance for said account.
    ///                  Bypasses ACL and can be called in Deposit tests
    /// @param _account  The account whose tokens will be burnt.
    /// @param _amount   The amount of tokens that will be burnt.
    function forceBurn(address _account, uint256 _amount) public {
        // NOTE: this uses internal function _burn instead of _burnFrom.
        // This will bypass allowance check for now.
        _burn(_account, _amount);
    }

    /// @dev                Zeros out the balance of the calling address.
    ///                     Does nothing if address has zero balance.
    function zeroBalance() public {
        uint256 currentBalance = balanceOf(msg.sender);
        if (currentBalance > 0){
            _burn(msg.sender, currentBalance);
        }
    }

    /// @dev                Uses exposed token functions to reset caller's balance.
    /// @param _newBalance  New balance to assign to caller
    function resetBalance(uint256 _newBalance) public {
        uint256 currentBalance = balanceOf(msg.sender);
        if(currentBalance > 0){
            forceBurn(msg.sender, currentBalance);
        }
        forceMint(msg.sender, _newBalance);
    }

    /// @dev                   Uses exposed token functions to reset the allowance
    ///                        of a given account.
    /// @param _spender        The allowed account.
    /// @param _newAllowance   New allowance to assign.
    function resetAllowance(address _spender, uint256 _newAllowance) public {
        uint256 currentAllowance = allowance(msg.sender, _spender);
        if (currentAllowance > 0){
            decreaseAllowance(_spender, currentAllowance);
        }
        approve(_spender, _newAllowance);
    }
}
