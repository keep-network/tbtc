pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "../../../contracts/external/IMedianizer.sol";

/// @notice A mock implementation of a medianizer price oracle
/// @dev This is used in the Keep testnets only. Mainnet uses the MakerDAO medianizer.
contract MockMedianizer is Ownable, IMedianizer {
    uint128 private value;

    constructor() public {
    // solium-disable-previous-line no-empty-blocks
    }

    function read() external view returns (uint256){
        return value;
    }

    function peek() external view returns (uint256, bool) {
        return (value, value > 0);
    }

    function setValue(uint128 _value) external onlyOwner {
        value = _value;
    }
}
