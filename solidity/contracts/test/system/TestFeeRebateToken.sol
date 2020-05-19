pragma solidity 0.5.17;

import "../../../contracts/system/FeeRebateToken.sol";

contract TestFeeRebateToken is FeeRebateToken {

    constructor(address _vendingMachine)
        FeeRebateToken(_vendingMachine)
    public{}

    /// @dev             We can't call FeeRebateToken mint function from tests because of ACL.
    ///                  This function bypasses ACL and can be called in tests.
    ///                  Mints a token and assigns it to an account.
    ///                  Uses the internal _mint function.
    /// @param _account  The account that will receive the token.
    /// @param _frtId    The Fee Rebate Token ID.
    function forceMint(address _account, uint256 _frtId) public returns (bool) {
        _mint(_account, _frtId);
        return true;
    }

    /// @dev             Delete a feeRebateToken. This function is used to
    ///                  test cases where the existence of an FRT impacts test outcomes.
    ///                  Must be called my the token owner.
    /// @param _frtId    The Fee Rebate Token ID.
    function burn(uint256 _frtId) public {
        _burn(msg.sender, _frtId);
    }
}
