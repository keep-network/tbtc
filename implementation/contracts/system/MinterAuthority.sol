pragma solidity 0.5.10;

import {TBTCSystem} from "./TBTCSystem.sol";

contract MinterAuthority {

    address internal _TBTCSystem;

    /// @dev Set the address of the System contract on contract initialization
    constructor(address _system) public {
        _TBTCSystem = _system;
    }

    modifier onlyDeposit(){
        require(TBTCSystem(_TBTCSystem).isDeposit(msg.sender), "caller must be a deposit");
        _;
    }
}