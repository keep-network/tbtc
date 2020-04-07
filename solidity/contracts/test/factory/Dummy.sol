pragma solidity 0.5.17;

/// @dev Dummy contract to used to test clone factory.
///      This contract will be cloned
contract Dummy{

    uint256 public state;

    function setState(uint256 _state) public {
        state = _state;
    }

    function getState()public view returns (uint256) {
        return state;
    }
}
