pragma solidity 0.4.25;

import {IUniswapExchange} from "../../../contracts/uniswap/IUniswapExchange.sol";
import {IBurnableERC20} from '../../../contracts/interfaces/IBurnableERC20.sol';

contract UniswapExchangeStub is IUniswapExchange {
    IBurnableERC20 tbtc;

    constructor(address tbtc) {
        tbtc = IBurnableERC20(tbtc);
    }

    // Give some mock liquidity to the stub
    function mockLiquidity() {
        require(msg.value > 0, "requires ETH for liquidity");
        uint ONE_TBTC = (10 ** 8);
        tbtc.transfer(address(this), ONE_TBTC * 10);
    }

    function getEthToTokenInputPrice(uint256 eth_sold)
        external view
        returns (uint256 tokens_bought)
    {

    }

    function ethToTokenSwapOutput(uint256 tokens_bought, uint256 deadline)
        external payable
        returns (uint256 eth_sold)
    {

    }
}