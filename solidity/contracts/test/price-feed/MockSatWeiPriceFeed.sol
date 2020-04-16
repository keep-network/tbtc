pragma solidity 0.5.17;

contract MockSatWeiPriceFeed{

    uint256 private price;

    /// @notice Get the current price of satoshi in wei.
    /// @return The price of one satoshi in wei.
    function getPrice()
        external view returns (uint256)
    {
        return price;
    }

    function setPrice(uint256 _price) public {
        price = _price;
    }

}
