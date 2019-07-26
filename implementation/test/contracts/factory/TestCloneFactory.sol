pragma solidity ^0.5.10;

import {CloneFactory} from '../../../contracts/proxy/CloneFactory.sol';
import {Deposit} from '../../../contracts/deposit/Deposit.sol';

contract TestCloneFactory is CloneFactory {

    event DepositCloneCreated(address depositCloneAddress);

    function getCloneAddress(address _target) external returns (address){
        address cloneAddress = createClone(_target);
        emit DepositCloneCreated(cloneAddress);
    }

    function checkClone(address _target, address _query)external view returns (bool result){
        return isClone(_target, _query);
    }
}