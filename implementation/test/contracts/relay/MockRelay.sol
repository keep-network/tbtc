pragma solidity ^0.5.10;

contract MockRelay {

  uint256 current;
  uint256 previous;

  function setMock(uint256 _current, uint256 _previous) public returns (bool) {
    current = _current;
    previous = _previous;
  }

  function getCurrentEpochDifficulty() external view returns (uint256) {
      return current;
  }

  function getPrevEpochDifficulty() external view returns (uint256) {
    return previous;
  }
}
