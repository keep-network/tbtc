pragma solidity ^0.5.10;

import {TBTCSystem} from '../../../contracts/system/TBTCSystem.sol';
import {TBTCToken} from "../../../contracts/system/TBTCToken.sol";

contract TBTCSystemStub is TBTCSystem {

    address _TBTCToken;

    function setExternalAddresses(address _tokenAddress) public {
        _TBTCToken = _tokenAddress;
    }

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
}
