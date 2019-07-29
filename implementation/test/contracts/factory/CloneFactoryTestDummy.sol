pragma solidity ^0.5.10;

contract CloneFactoryTestDummy{

    uint256 public state;

    function setState(uint256 _state) public {
        state = _state;
    }

    function getState()public view returns (uint256) {
        return state;
    }
}