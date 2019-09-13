pragma solidity 0.5.10;

interface IUniswapExchange {
    // Provide Liquidity
    function addLiquidity(uint256 min_liquidity, uint256 max_tokens, uint256 deadline) external payable returns (uint256);

    // Get Prices
    function getEthToTokenOutputPrice(uint256 tokens_bought) external view returns (uint256 eth_sold);

    // Trade ETH to ERC20
    function ethToTokenSwapOutput(uint256 tokens_bought, uint256 deadline) external payable returns (uint256  eth_sold);
}