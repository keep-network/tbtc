pragma solidity 0.4.25;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

contract TBTCToken is ERC20Detailed, ERC20 {
    /// @dev Constructor, calls ERC20Detailed constructor to set Token info
    ///      ERC20Detailed(TokenName, TokenSymbol, NumberOfDecimals)
    constructor() ERC20Detailed("Trustless bitcoin", "TBTC", 18) public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev             Mints an amount of the token and assigns it to an account.
    ///                  Uses the internal _mint function
    /// @param _account  The account that will receive the created tokens.
    /// @param _amount   The amount that will be created.
    function mint(address _account, uint256 _amount) public returns (bool){
        // NOTE: this is a public function with unchecked minting.
        // TODO: enforce calling authority.
        _mint(_account, _amount);
        return true;
    }

    /// @dev             Burns an amount of the token of a given account
    ///                  deducting from the sender's allowance for said account.
    ///                  Uses the internal _burn function.
    /// @param _account  The account whose tokens will be burnt.
    /// @param _amount   The amount of tokens that will be burnt.
    function burnFrom(address _account, uint256 _amount) public {
        // NOTE: this uses internal function _burn instead of _burnFrom.
        // This will bypass allowance check for now.
        // TODO: enforce calling authority.
        _burn(_account, _amount);
    }

    /// @dev           Transfer tokens from one address to another
    ///                Uses the internal _transfer function.
    /// @param _from   The address to send tokens from
    /// @param _to     The address to transfer to
    /// @param _value  The amount of tokens to be transferred
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool) {
        // NOTE: this overrides transferFrom in openZeppelin ERC20.sol
        // in order to bypass allowance check for now.
        // TODO: enforce calling authority.
        _transfer(_from, _to, _value);
        return true;
    }
}
