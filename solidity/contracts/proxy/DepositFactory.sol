pragma solidity 0.5.17;

import "./CloneFactory.sol";
import "../deposit/Deposit.sol";
import "../system/TBTCSystem.sol";
import "../system/TBTCToken.sol";
import "../system/FeeRebateToken.sol";
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
    address payable public masterDepositAddress;
    TBTCDepositToken tbtcDepositToken;
    TBTCSystem public tbtcSystem;
    TBTCToken public tbtcToken;
    FeeRebateToken public feeRebateToken;
    address public vendingMachineAddress;
    uint16 public keepThreshold;
    uint16 public keepSize;

    constructor(address _systemAddress)
        TBTCSystemAuthority(_systemAddress)
    public {}

    /// @dev                          Set the required external variables.
    /// @param _masterDepositAddress  The address of the master deposit contract.
    /// @param _tbtcSystem            Tbtc system contract.
    /// @param _tbtcToken             TBTC token contract.
    /// @param _tbtcDepositToken      TBTC Deposit Token contract.
    /// @param _feeRebateToken        AFee Rebate Token contract.
    /// @param _vendingMachineAddress Address of the Vending Machine contract.
    /// @param _keepThreshold         Minimum number of honest keep members.
    /// @param _keepSize              Number of all members in a keep.
    function setExternalDependencies(
        address payable _masterDepositAddress,
        TBTCSystem _tbtcSystem,
        TBTCToken _tbtcToken,
        TBTCDepositToken _tbtcDepositToken,
        FeeRebateToken _feeRebateToken,
        address _vendingMachineAddress,
        uint16 _keepThreshold,
        uint16 _keepSize
    ) public onlyTbtcSystem {
        masterDepositAddress = _masterDepositAddress;
        tbtcDepositToken = _tbtcDepositToken;
        tbtcSystem = _tbtcSystem;
        tbtcToken = _tbtcToken;
        feeRebateToken = _feeRebateToken;
        vendingMachineAddress = _vendingMachineAddress;
        keepThreshold = _keepThreshold;
        keepSize = _keepSize;
    }

    event DepositCloneCreated(address depositCloneAddress);

    /// @notice                Creates a new deposit instance and mints a TDT.
    ///                        This function is currently the only way to create a new deposit.
    /// @dev                   Calls `Deposit.createNewDeposit` to initialize the instance.
    ///                        Mints the TDT to the function caller.
    //                         (See `TBTCDepositToken` for more info on TDTs).
    /// @return                True if successful, otherwise revert.
    function createDeposit (uint64 _lotSizeSatoshis) public payable returns(address) {
        address cloneAddress = createClone(masterDepositAddress);

        TBTCDepositToken(tbtcDepositToken).mint(msg.sender, uint256(cloneAddress));

        Deposit deposit = Deposit(address(uint160(cloneAddress)));
        deposit.initialize(address(this));
        deposit.createNewDeposit.value(msg.value)(
                tbtcSystem,
                tbtcToken,
                tbtcDepositToken,
                feeRebateToken,
                vendingMachineAddress,
                keepThreshold,
                keepSize,
                _lotSizeSatoshis
            );

        emit DepositCloneCreated(cloneAddress);

        return cloneAddress;
    }
}
