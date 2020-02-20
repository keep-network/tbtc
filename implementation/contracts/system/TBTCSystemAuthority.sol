pragma solidity 0.5.10;

/// @title  TBTC System Authority.
/// @notice Contract to secure function calls to the TBTC System contract.
/// @dev    The `TBTCSystem` contract address is passed as a constructor parameter.
contract TBTCSystemAuthority {

    address internal tbtcSystem;

    /// @notice Set the address of the System contract on contract initialization.
    constructor(address _tbtcSystem) public {
        tbtcSystem = _tbtcSystem;
    }

    /// @notice Function modifier ensures modified function is only called by TBTCSystem.
    modifier onlyTbtcSystem(){
        require(msg.sender == tbtcSystem, "Caller must be tbtcSystem contract");
        _;
    }
}