pragma solidity ^0.5.10;

import "./CloneFactory.sol";
import "../deposit/Deposit.sol";
import "../system/TBTCSystem.sol";
import {DepositOwnerToken} from "../system/DepositOwnerToken.sol";

/// @title Deposit Factory
/// @notice Factory for the creation of new deposit clones.
/// @dev We avoid redeployment of deposit contract by using the clone factory.
/// Proxy delegates calls to Deposit and therefore does not affect deposit state.
/// This means that we only need to deploy the deposit contracts once.
/// The factory provides clean state for every new deposit clone.
contract DepositFactory is CloneFactory{

    // Holds the address of the deposit contract
    // which will be used as a master contract for cloning.
    address public masterDepositAddress;

    event DepositCloneCreated(address depositCloneAddress);

    /// @dev                          Set the master deposit contract address
    ///                               on contract initialization
    /// @param _masterDepositAddress  The address of the master deposit contract
    constructor(address _masterDepositAddress) public {
        masterDepositAddress = _masterDepositAddress;
    }

    /// @notice                Creates a new deposit instance
    /// @dev                   Calls createNewDeposit from deposit contract as init method.
    ///                        We don't offer pure createClone, meaning that the only way
    ///                        to create a clone is by also calling createNewDeposit()
    ///                        Deposits created this way will never pass by state 0 (START)
    /// @param _TBTCSystem     Address of system contract
    /// @param _TBTCToken      Address of TBTC token contract
    /// @param _keepThreshold  Minimum number of honest keep members
    /// @param _keepSize       Number of all members in a keep
    /// @return                True if successful, otherwise revert
    function createDeposit (
        address _TBTCSystem,
        address _TBTCToken,
        address _DepositOwnerToken,
        uint256 _keepThreshold,
        uint256 _keepSize
    ) public payable returns(address) {
        address cloneAddress = createClone(masterDepositAddress);

        Deposit(address(uint160(cloneAddress))).createNewDeposit.value(msg.value)(
            _TBTCSystem,
            _TBTCToken,
            _DepositOwnerToken,
            _keepThreshold,
            _keepSize);

        DepositOwnerToken(_depositOwnerToken).mint(msg.sender, uint256(cloneAddress));

        emit DepositCloneCreated(cloneAddress);

        return cloneAddress;
    }
}
