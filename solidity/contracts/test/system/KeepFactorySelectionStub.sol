pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

import "../../system/KeepFactorySelection.sol";

contract KeepFactorySelectionStub {
    using KeepFactorySelection for KeepFactorySelection.Storage;
    KeepFactorySelection.Storage keepFactorySelection;

    function initialize(IBondedECDSAKeepFactory _defaultFactory) public {
        keepFactorySelection.initialize(_defaultFactory);
    }

    function selectFactory() public view returns (IBondedECDSAKeepFactory) {
        return keepFactorySelection.selectFactory();
    }

    function selectFactoryAndRefresh() public returns (IBondedECDSAKeepFactory) {
        return keepFactorySelection.selectFactoryAndRefresh();
    }

    function setFullyBackedKeepFactory(address _fullyBackedFactory) public {
        keepFactorySelection.setFullyBackedKeepFactory(_fullyBackedFactory);
    }

    function setKeepFactorySelector(address _factorySelector) public {
        keepFactorySelection.setKeepFactorySelector(_factorySelector);
    }
}

contract KeepFactorySelectorStub is KeepFactorySelector {

    bool internal defaultMode = true;
    bool internal maliciousMode = false;

    function selectFactory(
        uint256 _seed,
        IBondedECDSAKeepFactory _defaultFactory,
        IBondedECDSAKeepFactory _fullyBackedFactory
    ) external view returns (IBondedECDSAKeepFactory) {
        if (maliciousMode) {
            return IBondedECDSAKeepFactory(0xaFacEbadfAceCffeeFaceacebacEAfAceCfFEeA0);
        } else if (defaultMode) {
            return _defaultFactory;
        } else {
            return _fullyBackedFactory;
        }
    }

    function setMaliciousMode() public {
        maliciousMode = true;
    }

    function unsetMaliciousMode() public {
        maliciousMode = false;
    }

    function setDefaultMode() public {
        defaultMode = true;
    }

    function setFullyBackedMode() public {
        defaultMode = false;
    }
}
