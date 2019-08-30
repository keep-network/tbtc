pragma solidity ^0.5.10;

import {TBTCSystem} from '../../../contracts/system/TBTCSystem.sol';

contract TBTCSystemStub is TBTCSystem {
    address keepAddress = address(7);

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

    function requestNewKeep(uint256 _m, uint256 _n) external payable returns (address _keepAddress) {
        _m; _n;
        return keepAddress;
    }
}
