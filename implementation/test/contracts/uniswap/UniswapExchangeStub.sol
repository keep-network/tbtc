pragma solidity 0.4.25;

import {IUniswapExchange} from "../../../contracts/uniswap/IUniswapExchange.sol";
import {IBurnableERC20} from '../../../contracts/interfaces/IBurnableERC20.sol';

contract UniswapExchangeStub {
    IBurnableERC20 tbtc;

    uint256 ethToTokenPrice = 0;

    constructor(address tbtc) {
        tbtc = IBurnableERC20(tbtc);
    }

    function() public {
        revert("unimplemented");
    }

    // Give some mock liquidity to the stub
    function mockLiquidity() {
        require(msg.value > 0, "requires ETH for liquidity");
        uint ONE_TBTC = (10 ** 8);
        tbtc.transfer(address(this), ONE_TBTC * 10);
    }

    function setEthToTokenInputPrice(uint256 _price) {
        ethToTokenPrice = _price;
    }

    function getEthToTokenInputPrice(uint256 eth_sold)
        external view
        returns (uint256 tokens_bought)
    {
        return ethToTokenPrice;
    }

    function ethToTokenSwapOutput(uint256 tokens_bought, uint256 deadline)
        external payable
        returns (uint256 eth_sold)
    {
        deadline;
        require(msg.value == ethToTokenPrice, "incorrect eth sent");
        tbtc.transfer(msg.sender, tokens_bought);
        
        return msg.value;
    }
}