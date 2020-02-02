pragma solidity ^0.5.10;

import "./CloneFactory.sol";
import "../deposit/Deposit.sol";
import "../system/TBTCSystem.sol";
import "../system/TBTCSystemAuthority.sol";
import {TBTCDepositToken} from "../system/TBTCDepositToken.sol";


/// @title Deposit Factory
/// @notice Factory for the creation of new deposit clones.
/// @dev We avoid redeployment of deposit contract by using the clone factory.
/// Proxy delegates calls to Deposit and therefore does not affect deposit state.
/// This means that we only need to deploy the deposit contracts once.
/// The factory provides clean state for every new deposit clone.
contract DepositFactory is CloneFactory, TBTCSystemAuthority{

    // Holds the address of the deposit contract
    // which will be used as a master contract for cloning.
    address public masterDepositAddress;
    address public tbtcSystem;
    address public tbtcToken;
    address public tbtcDepositToken;
    address public feeRebateToken;
    address public vendingMachine;
    uint256 public keepThreshold;
    uint256 public keepSize;

    constructor(address _systemAddress) 
        TBTCSystemAuthority(_systemAddress)
    public {}

    /// @dev                          Set the required external variables
    /// @param _masterDepositAddress  The address of the master deposit contract
    /// @param _tbtcSystem            Address of system contract
    /// @param _tbtcToken             Address of TBTC token contract
    /// @param _depositOwnerToken     Address of the Deposit Owner Token contract
    /// @param _feeRebateToken        Address of the Fee Rebate Token contract
    /// @param _vendingMachine        Address of the Vending Machine contract
    /// @param _keepThreshold         Minimum number of honest keep members
    /// @param _keepSize              Number of all members in a keep
    function setExternalDependencies(
        address _masterDepositAddress,
        address _tbtcSystem,
        address _tbtcToken,
        address _depositOwnerToken,
        address _feeRebateToken,
        address _vendingMachine,
        uint256 _keepThreshold,
        uint256 _keepSize
    ) public onlyTbtcSystem{
        masterDepositAddress = _masterDepositAddress;
        tbtcSystem = _tbtcSystem;
        tbtcToken = _tbtcToken;
        tbtcDepositToken = _depositOwnerToken;
        feeRebateToken = _feeRebateToken;
        vendingMachine = _vendingMachine;
        keepThreshold = _keepThreshold;
        keepSize = _keepSize;
    }

    event DepositCloneCreated(address depositCloneAddress);

    /// @notice                Creates a new deposit instance
    /// @dev                   Calls createNewDeposit from deposit contract as init method.
    ///                        We don't offer pure createClone, meaning that the only way
    ///                        to create a clone is by also calling createNewDeposit()
    ///                        Deposits created this way will never pass by state 0 (START)
    /// @return                True if successful, otherwise revert
    function createDeposit (uint256 _lotSize) public payable returns(address) {
        address cloneAddress = createClone(masterDepositAddress);

        Deposit(address(uint160(cloneAddress))).createNewDeposit.value(msg.value)(
            tbtcSystem,
            tbtcToken,
            tbtcDepositToken,
            feeRebateToken,
            vendingMachine,
            keepThreshold,
            keepSize,
            _lotSize
        );

        TBTCDepositToken(tbtcDepositToken).mint(msg.sender, uint256(cloneAddress));

        emit DepositCloneCreated(cloneAddress);

        return cloneAddress;
    }
}
