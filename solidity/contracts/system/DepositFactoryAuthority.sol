pragma solidity 0.5.17;

/// @title  Deposit Factory Authority
/// @notice Contract to secure function calls to the Deposit Factory.
/// @dev    Secured by setting the depositFactory address and using the onlyFactory
///         modifier on functions requiring restriction.
contract DepositFactoryAuthority {

    bool internal _initialized = false;
    address internal _depositFactory;

    /// @notice Set the address of the System contract on contract initialization.
    function initialize(address _factory) public {
        require(! _initialized, "Factory can only be initialized once.");

        _depositFactory = _factory;
        _initialized = true;
    }

    /// @notice Function modifier ensures modified function is only called by set deposit factory.
    modifier onlyFactory(){
        require(_initialized, "Factory initialization must have been called.");
        require(msg.sender == _depositFactory, "Caller must be depositFactory contract");
        _;
    }
}
