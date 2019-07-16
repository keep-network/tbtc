pragma solidity 0.4.25;

contract UniswapFactoryStub {
    address private tbtcExchange;

    constructor() public {
        return;
    }

    function setTbtcExchange(address _tbtcExchange) public {
        tbtcExchange = _tbtcExchange;
    }

    function getExchange(address token) external view returns (address exchange) {
        token;
        return tbtcExchange;
    }
}