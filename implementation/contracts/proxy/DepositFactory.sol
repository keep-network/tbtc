pragma solidity ^0.4.25;

import "./CloneFactory.sol";
import "../deposit/Deposit.sol";


// We avoid redeployment of deposit contract by using the cloneFactory.
// Proxy delegates calls to Deposit and therefore does not effect deposit state
// This means that we only need to deploy the deposit contracts once.
// The factory provides clean state for every new deposit clone.

contract DepositFactory is CloneFactory{

    // holds the address of the deposit contract
    address public implementation;

    event DepositCloneCreated(address clonedContract);

    /// @dev                    Set the deposit address on contract initialization
    /// @param _implementation  The address of the deployed Deposit contract
    constructor(address _implementation) public {
        implementation = _implementation;
    }

    /// @notice             Creates a new deposit instance
    /// @dev                Calls createNewDeposit from deposit contract as init method. We don't offer pure createClone,
    ///                     meaning that the only way to create a clone is by also calling createNewDeposit()
    ///                     Deposits created this way will never pass by state 0 (START)
    /// @param _TBTCSystem  Address of system contract
    /// @param _TBTCToken   Address of Token contract
    /// @param _KeepBridge  Address of Keep contract
    /// @param _m           m for m-of-n
    /// @param _n           n for m-of-n
    /// @return             True if successful, otherwise revert
    function createDeposit (
        address _TBTCSystem,
        address _TBTCToken,
        address _KeepBridge,
        uint256 _m,
        uint256 _n
    ) public {

        address clone = createClone(implementation);
        Deposit(clone).createNewDeposit(_TBTCSystem, _TBTCToken, _KeepBridge, _m, _n);

        emit DepositCloneCreated(clone);
    }
}