pragma solidity 0.4.25;

import "../interfaces/IPriceOracle.sol";

/**
 * The price oracle implements a simple price feed, managed
 * by a trusted operator.
 */
contract PriceOracleV1 is IPriceOracle {
    // Price of BTC expressed in ETH, denominated in weis to satoshis.
    // A bitcoin has 8 decimal places, the smallest unit a satoshi,
    // meaning 100 000 000 satoshis = 1 bitcoin
    // An ether by contract, has 18 decimal places, the smallest unit a wei,
    // meaning 1,000,000,000,000,000,000 wei = 1 ether
    // 
    // eg. 1 BTC : 32.32 ETH (Jun 2019), which represents:
    //     100 000 000 satoshis : 32 320 000 000 000 000 000 wei
    //     or simplified:
    //     1 : 32 320 000 000 00
    //     price = 3232000000000
    // 
    // uint128 can store an int of max size 2^128 (3.403 e38)
    // 38 decimal places should be enough... :)
    uint128 private price;

    constructor(uint128 _defaultPrice) public {
        price = _defaultPrice;
    }

    function getPrice() external view returns (uint128) {
        return price;
    }

    function updatePrice(uint128 _price) public {
        price = _price;
    }
}