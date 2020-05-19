pragma solidity 0.5.17;

// Import of `TestnetRelay` will build artifacts for `TestnetRelay` and `OnDemandSPVProof`
// that are required by the migration script but are not used for local contracts
// deployment.
import "@summa-tx/relay-sol/contracts/TestnetRelay.sol";

contract MockRelay {
    uint256 current = 1;
    uint256 previous = 1;

    function getCurrentEpochDifficulty() external view returns (uint256) {
        return current;
    }

    function getPrevEpochDifficulty() external view returns (uint256) {
        return previous;
    }

    function setCurrentEpochDifficulty(uint256 _current) public {
        current = _current;
    }

    function setPrevEpochDifficulty(uint256 _previous) public {
        previous = _previous;
    }
}
