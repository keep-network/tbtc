pragma solidity ^0.5.10;

contract UniswapFactoryStub {
    address private exchange;

    constructor() public {
        return;
    }

    function setExchange(address _exchange) public {
        exchange = _exchange;
    }

    function getExchange(address token) external view returns (address) {
        token;
        return exchange;
    }
}