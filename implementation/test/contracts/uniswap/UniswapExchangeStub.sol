pragma solidity ^0.5.10;

import {TestToken} from '../deposit/TestToken.sol';

contract UniswapExchangeStub {
    TestToken tbtc;


    // The below returns an absurdly large price for tBTC
    // such that attemptToLiquidateOnchain will return early, from not being funded enough
    uint256 ethPrice = 10**8;

    constructor(address _tbtc) public {
        tbtc = TestToken(_tbtc);
    }

    function setEthPrice(uint256 _ethPrice) public {
        ethPrice = _ethPrice;
    }

    function getEthToTokenOutputPrice(uint256 tokens_sold)
        external view
        returns (uint256)
    {
        tokens_sold;
        return ethPrice;
    }

    function ethToTokenSwapOutput(uint256 tokens_bought, uint256 deadline)
        external payable
        returns (uint256 eth_sold)
    {
        deadline;
        require(msg.value     == ethPrice, "incorrect eth sent");
        tbtc.forceMint(msg.sender, tokens_bought);
        return msg.value;
    }
}