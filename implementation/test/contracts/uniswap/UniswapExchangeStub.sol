pragma solidity ^0.5.10;

import {TBTCToken} from '../../../contracts/system/TBTCToken.sol';

contract UniswapExchangeStub {
    TBTCToken tbtc;


    // The below returns an absurdly large price for tBTC
    // such that attemptToLiquidateOnchain will return early, from not being funded enough
    uint256 ethPrice = 10**8;

    constructor(address _tbtc) public {
        tbtc = TBTCToken(_tbtc);
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
        // require(tokens_bought == tbtcPrice, "incorrect tbtc ask");
        // require(tbtc.balanceOf(address(this)) >= tokens_bought, "not enough tbtc for trade");
        tbtc.mint(msg.sender, tokens_bought);
        // require(tbtc.transferFrom(address(this), msg.sender, tokens_bought), "transfer failed");
        return msg.value;
    }
}