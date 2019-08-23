pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../../contracts/system/TBTCToken.sol";

contract TestToken is TBTCToken{

    constructor(address _system)
        TBTCToken(_system)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev             Mints an amount of the token and assigns it to an account.
    ///                  Uses the internal _mint function
    /// @param _account  The account that will receive the created tokens.
    /// @param _amount   The amount of tokens that will be created.
    function forceMint(address _account, uint256 _amount) public returns (bool){
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
    function forceBurn(address _account, uint256 _amount) public {
        // NOTE: this uses internal function _burn instead of _burnFrom.
        // This will bypass allowance check for now.
        // TODO: enforce calling authority.
        _burn(_account, _amount);
    }
}
