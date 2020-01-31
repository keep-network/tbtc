pragma solidity 0.5.10;

contract TBTCSystemAuthority {

    address internal _tbtcSystem;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _tbtcSystem) public {
        _tbtcSystem = _tbtcSystem;
    }

    /// @notice Function modifier ensures modified function is only called by set deeposit factory
    modifier onlyTbtcSystem(){
        require(msg.sender == _tbtcSystem, "Caller must be depositFactory contract");
        _;
    }
}