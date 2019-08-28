pragma solidity ^0.5.10;

import {TBTCToken} from '../../../contracts/system/TBTCToken.sol';

contract UniswapExchangeStub {
    TBTCToken tbtc;

    uint256 ethPrice = 10**8;
    // uint256 tbtcPrice = 0;

    constructor(address _tbtc) public {
        tbtc = TBTCToken(_tbtc);
    }

    // Give some mock liquidity to the stub
    // function mockLiquidity(uint256 _tbtcAmount) public payable {
    //     require(msg.value > 0, "requires ETH for liquidity");
    //     tbtc.transferFrom(msg.sender, address(this), _tbtcAmount);
    // }

    // function setPrices(uint256 _ethPrice, uint256 _tbtcPrice) public {
    //     ethPrice = _ethPrice;
    //     tbtcPrice = _tbtcPrice;
    // }

    function getTokenToEthInputPrice(uint256 tokens_sold)
        external view
        returns (uint256)
    {
        tokens_sold;
        // TODO(liamz): fix this
        // return msg.sender.balance;
        return 2**250;
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