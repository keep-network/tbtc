pragma solidity 0.4.25;

contract MinterAuthority {

    address internal _systemAddress;

    /// @dev Set the address of the System contract on contract initialization
    constructor(address _system) public {
        _systemAddress = _system;
    }

    modifier onlySystem(){
        require(msg.sender == _systemAddress, "Caller must be TBTC System contract");
        _;
    }
}                   