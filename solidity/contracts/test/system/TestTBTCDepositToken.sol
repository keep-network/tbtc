pragma solidity 0.5.17;

import "../../../contracts/system/TBTCDepositToken.sol";

contract TestTBTCDepositToken is TBTCDepositToken {

    constructor(address _factory)
        TBTCDepositToken(_factory)
    public{}

    /// @dev             We can't call TBTCDepositToken mint function from tests because of ACL.
    ///                  This function bypasses ACL and can be called in tests.
    ///                  Mints a token and assigns it to an account.
    ///                  Uses the internal _mint function.
    /// @param _account  The account that will receive the token.
    /// @param _tdtId    The tBTC Deposit Token ID.
    function forceMint(address _account, uint256 _tdtId) public returns (bool) {
        if(_exists(_tdtId)){
            _burn(_tdtId);
        }
        _mint(_account, _tdtId);
        return true;
    }
}
