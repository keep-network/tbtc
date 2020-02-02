pragma solidity ^0.5.10;

import "../../../contracts/system/VendingMachine.sol";

contract TestVendingMachine is VendingMachine {

    constructor(address _tbtcSystem)
        VendingMachine(_tbtcSystem)
    public {}

    function setExternalAddresses(
        address _tbtcToken,
        address _tbtcDepositToken,
        address _feeRebateToken
    ) public {
        tbtcToken = TBTCToken(_tbtcToken);
        tbtcDepositToken = TBTCDepositToken(_tbtcDepositToken);
        feeRebateToken = FeeRebateToken(_feeRebateToken);
    }
}
