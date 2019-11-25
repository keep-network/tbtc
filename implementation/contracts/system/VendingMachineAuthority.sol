pragma solidity ^0.5.10;

import {TBTCSystem} from "./TBTCSystem.sol";

contract VendingMachineAuthority {

    address internal VendingMachine;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _vendingMachine) public {
        VendingMachine = _vendingMachine;
    }

    /// @notice Function modifier ensures modified function caller address is the vending machine
    modifier onlyVendingMachine(){
        require(msg.sender == VendingMachine , "caller must be the vending machine");
        _;
    }
}