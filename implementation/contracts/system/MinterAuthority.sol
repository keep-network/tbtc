pragma solidity 0.5.10;

import {TBTCSystem} from "./TBTCSystem.sol";

contract MinterAuthority {

    address internal _systemAddress;

    /// @dev Set the address of the System contract on contract initialization
    constructor(address _system) public {
        _systemAddress = _system;
    }

    modifier isApproved(){
        require(TBTCSystem(_systemAddress).isDeposit(msg.sender), "caller must be a deposit");
        _;
    }
}