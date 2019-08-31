pragma solidity ^0.5.10;

import {TBTCSystem} from '../../../contracts/system/TBTCSystem.sol';

contract TBTCSystemStub is TBTCSystem {

    function setOraclePrice(uint256 _oraclePrice) external {
        oraclePrice = _oraclePrice;
    }

    function setCurrentDiff(uint256 _currentDifficulty) external {
        currentDifficulty = _currentDifficulty;
    }

    function setPreviousDiff(uint256 _previousDifficulty) external {
        previousDifficulty = _previousDifficulty;
    }

    // override parent
    function approvedToLog(address _caller) public view returns (bool) {
        _caller; return true;
    }


    address depositOwner;
    
    function setDepositOwner(address user) public returns (address) {
        depositOwner = user;
    }

    function ownerOf(uint _id) public view returns (address) {
        _id;
        return depositOwner;
    }
}
