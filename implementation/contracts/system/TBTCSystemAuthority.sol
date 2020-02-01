pragma solidity 0.5.10;

contract TBTCSystemAuthority {

    address internal tbtcSystem;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _tbtcSystem) public {
        tbtcSystem = _tbtcSystem;
    }

    /// @notice Function modifier ensures modified function is only called by set deeposit factory
    modifier onlyTbtcSystem(){
        require(msg.sender == tbtcSystem, "Caller must be tbtcSystem contract");
        _;
    }
}