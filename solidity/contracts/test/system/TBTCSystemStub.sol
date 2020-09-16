pragma solidity 0.5.17;

import {TBTCSystem} from "../../../contracts/system/TBTCSystem.sol";

contract TBTCSystemStub is TBTCSystem {
    address keepAddress = address(7);
    uint256 oraclePrice = 10 ** 12;

    constructor(address _priceFeed, address _relay)
        // Set expected factory address to 0-address.
        // Address is irelevant as test use forceMint function to bypass ACL
        TBTCSystem(_priceFeed, _relay)
    public {
        // solium-disable-previous-line no-empty-blocks
        lotSizesSatoshis = [10**5, 10**6, 10**7, 2 * 10**7, 5 * 10**7, 10**8]; // [0.01, 0.1, 0.2, 0.5, 1.0] BTC
    }

    function setOraclePrice(uint256 _oraclePrice) external {
        oraclePrice = _oraclePrice;
    }

    /// @dev Override TBTCSystem.fetchBitcoinPrice, don't call out to the price feed.
    function fetchBitcoinPrice() external view returns (uint256) {
        return oraclePrice;
    }

    function setKeepAddress(address _keepAddress) external {
        keepAddress = _keepAddress;
    }

    function requestNewKeep(uint64, uint256) external payable returns (address _keepAddress) {
        return keepAddress;
    }
}
