pragma solidity 0.5.10;

contract DepositFactoryAuthority {

    address internal _depositFactory;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _factory) public {
        _depositFactory = _factory;
    }

    /// @notice Function modifier ensures modified function is only called by set deeposit factory
    modifier onlyFactory(){
        require(msg.sender == _depositFactory, "Caller must be depositFactory contract");
        _;
    }
}