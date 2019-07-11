pragma solidity ^0.4.25;

import "./CloneFactory.sol";
import "../deposit/Deposit.sol";

contract DepositFactory is CloneFactory{

    //holds the address of the deposit contract
    address public implementation;

    //holds list of all clones for easy access
    address[] public cloneContracts;

    event DepositCloneCreated(address clonedContract);

    //we set the deposit address on contract initialization
    constructor(address _implementation) public {
        implementation = _implementation;
    }

    //creates a deposit clone, calls deposit `init` method
    function createDeposit (
        address _TBTCSystem,
        address _TBTCToken,
        address _KeepBridge,
        uint256 _m,
        uint256 _n
        ) public {

        //Note 1: we don't offer pure createClone.
        //This means the only way to create a clone is by also calling the deposit init function
        //Deposits created this way will never pass by state 0 (START)

        //Note 2: we can avoid using the cloneContract array by listening to events
        //This will forgo the cost asociated with writing to storage, but keeping it for testing temporarily.

        address clone = createClone(implementation);
        Deposit(clone).createNewDeposit(_TBTCSystem, _TBTCToken, _KeepBridge, _m, _n);

        emit DepositCloneCreated(clone);

        //store clone address so we can access later
        cloneContracts.push(clone);
    }

    //retreive clone address from array using index
    function getAddress(uint _index) public view returns (address) {
        require(_index < cloneContracts.length, "index out of bounds");
        return cloneContracts[_index];
    }
}