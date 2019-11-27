pragma solidity ^0.5.10;

contract VendingMachineAuthority {
    address internal VendingMachine;

    constructor(address _vendingMachine) public {
        VendingMachine = _vendingMachine;
    }

    /// @notice Function modifier ensures modified function caller address is the vending machine
    modifier onlyVendingMachine() {
        require(msg.sender == VendingMachine, "caller must be the vending machine");
        _;
    }
}