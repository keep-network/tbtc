pragma solidity ^0.5.10;

import {TestToken} from '../deposit/TestToken.sol';
import {IUniswapExchange} from '../../../contracts/external/IUniswapExchange.sol';

contract UniswapExchangeStub is IUniswapExchange {
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

    function addLiquidity(uint256 min_liquidity, uint256 max_tokens, uint256 deadline)
        external payable
        returns (uint256)
    {
        require(msg.value > 0, "ETH missing from addLiquidity");
        tbtc.forceMint(address(this), max_tokens);
        // Stub doesn't implement the internal Uniswap token (UNI),
        // so return 0 here for total minted UNI.
        return 0;
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
        require(msg.value == ethPrice, "incorrect eth sent");
        require(tbtc.balanceOf(address(this)) >= tokens_bought, "not enough TBTC liquidity mocked");
        tbtc.transfer(msg.sender, tokens_bought);
        return msg.value;
    }
}