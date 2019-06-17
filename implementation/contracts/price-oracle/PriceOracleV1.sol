pragma solidity 0.4.25;

import "../interfaces/IPriceOracle.sol";
// import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import {SafeMath} from "../bitcoin-spv/SafeMath.sol";

/**
 * The price oracle implements a simple price feed, managed
 * by a trusted operator.
 */
contract PriceOracleV1 is IPriceOracle {
    using SafeMath for uint128;

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

    // The time at which the last price update is considered to be stale.
    uint256 private zzz;

    // Trusted user that updates the oracle.
    address private operator;

    event PriceUpdated(uint128 price, uint256 zzz);

    constructor(
        address _operator,
        uint128 _defaultPrice
    ) public {
        operator = _operator;
        price = _defaultPrice;
        zzz = block.timestamp + 1 hours;
    }

    function getPrice() external view returns (uint128) {
        require(block.timestamp < zzz, "stale price");
        return price;
    }

    function updatePrice(uint128 _newPrice) external {
        require(msg.sender == operator, "unauthorised");

        // abs(1 - (p1 / p0)) > 0.01
        // p0 * 0.01 = minimum delta
        // 1% = 0.01 = 1/100
        uint256 minDelta = price.div(100);
        uint256 delta = _newPrice > price
                        ? _newPrice.sub(price)
                        : price.sub(_newPrice);
        require(delta > minDelta, "Price change is negligible (>1%)");

        price = _newPrice;
        zzz = block.timestamp + 6 hours;

        emit PriceUpdated(price, zzz);
    }
}