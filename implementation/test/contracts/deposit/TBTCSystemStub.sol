pragma solidity ^0.5.10;

import {TBTCSystem} from '../../../contracts/system/TBTCSystem.sol';

contract TBTCSystemStub is TBTCSystem {
    address keepAddress = address(7);
    uint256 oraclePrice = 10 ** 12;

    constructor(address _priceFeed, address _relay)
        // Set expected factory address to 0-address.
        // Address is irelevant as test use forceMint function to bypass ACL
        TBTCSystem(_priceFeed, _relay)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    function setOraclePrice(uint256 _oraclePrice) external {
        oraclePrice = _oraclePrice;
    }

    /// @dev Override TBTCSystem.fetchBitcoinPrice, don't call out to the price feed.
    function fetchBitcoinPrice() external view returns (uint256) {
        return oraclePrice;
    }

    // override parent
    function approvedToLog(address) public pure returns (bool) {
        return true;
    }

    function requestNewKeep(uint256, uint256, uint256) external payable returns (address _keepAddress) {
        return keepAddress;
    }
}
