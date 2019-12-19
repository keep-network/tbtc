pragma solidity ^0.5.10;

import "../../../contracts/system/DepositOwnerToken.sol";

contract TestDepositOwnerToken is DepositOwnerToken {
    /// @dev             We can't call DepositOwnerToken mint function from tests because of ACL.
    ///                  This function bypasses ACL and can be called in tests.
    ///                  Mints a token and assigns it to an account.
    ///                  Uses the internal _mint function.
    /// @param _account  The account that will receive the token.
    /// @param _dotId    The Deposit Owner Token ID
    function forceMint(address _account, uint256 _dotId) public returns (bool) {
        _mint(_account, _dotId);
        return true;
    }
}