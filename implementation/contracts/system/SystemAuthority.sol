pragma solidity 0.5.10;

contract SystemAuthority {

    address internal _tbtcSystem;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _systemAddress) public {
        _tbtcSystem = _systemAddress;
    }

    /// @notice Function modifier ensures modified function is only called by set deeposit factory
    modifier onlySystem(){
        require(msg.sender == _tbtcSystem, "Caller must be depositFactory contract");
        _;
    }
}