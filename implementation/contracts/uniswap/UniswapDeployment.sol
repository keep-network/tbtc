pragma solidity 0.4.25;

/// @notice A contract for tracking the location of Uniswap testnet deployments
/// @dev useful for usage in tests/migrations, but SHOULD NOT be used by other contracts
contract UniswapDeployment {
    address public factory;
    address public exchange;

    constructor(
        address _factory,
        address _exchange
    ) public {
        factory = _factory;
        exchange = _exchange;
    }
}