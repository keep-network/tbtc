pragma solidity ^0.5.10;

contract MockRelay {

  uint256 current = 1;
  uint256 previous = 1;

  function setCurrentEpochDifficulty(uint256 _current) public {
    current = _current;
  }

  function getCurrentEpochDifficulty() external view returns (uint256) {
      return current;
  }

  function setPrevEpochDifficulty(uint256 _previous) public {
    previous = _previous;
  }

  function getPrevEpochDifficulty() external view returns (uint256) {
    return previous;
  }
}
