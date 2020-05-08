pragma solidity 0.5.17;

import "../../system/IKeepFactorySelector.sol";

contract KeepFactorySelectorStub is IKeepFactorySelector {

    bool internal regularMode = true;

    function selectFactory(
        uint256 _seed,
        IBondedECDSAKeepFactory _regularFactory,
        IBondedECDSAKeepFactory _fullyBackedFactory
    ) external view returns (IBondedECDSAKeepFactory) {
        if (regularMode) {
            return _regularFactory;
        } else {
            return _fullyBackedFactory;
        }
    }

    function setRegularMode() public {
        regularMode = true;
    }

    function setFullyBackedMode() public {
        regularMode = false;
    }
}
