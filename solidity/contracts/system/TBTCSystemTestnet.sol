pragma solidity 0.5.17;

import {TBTCSystem} from "./TBTCSystem.sol";

contract TBTCSystemTestnet is TBTCSystem {
    constructor(address _priceFeed, address _relay)
        public
        TBTCSystem(_priceFeed, _relay)
    {
        // solium-disable-previous-line no-empty-blocks
        lotSizesSatoshis = [
            10**3,
            10**4,
            10**5,
            10**6,
            10**7,
            2 * 10**7,
            5 * 10**7,
            10**8
        ]; // [0.00001, 0.0001, 0.001, 0.01, 0.1, 0.2, 0.5, 1.0] BTC
    }
}
