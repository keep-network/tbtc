pragma solidity ^0.5.10;

import "../../../contracts/system/VendingMachine.sol";

contract TestVendingMachine is VendingMachine {
    constructor()
        public VendingMachine(address(0), address(0))
    {
        // solium-disable-previous-line no-empty-blocks
    }

    function setExteriorAddresses(
        address _tbtcToken,
        address _depositOwnerToken
    ) external {
        tbtcToken = TBTCToken(_tbtcToken);
        depositOwnerToken = DepositOwnerToken(_depositOwnerToken);
    }
}