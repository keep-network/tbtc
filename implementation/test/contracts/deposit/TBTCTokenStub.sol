pragma solidity 0.4.25;

import {TBTCToken} from '../../../contracts/system/TBTCToken.sol';

contract TBTCTokenStub is TBTCToken{
    
    function clearBalance(address _of)external{
        uint256 currentBalance = balanceOf(_of);
        burnFrom(_of, currentBalance);
    }
}
