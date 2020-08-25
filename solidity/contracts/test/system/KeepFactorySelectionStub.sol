pragma solidity 0.5.17;

import {
    IBondedECDSAKeepVendor
} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";
import {
    IBondedECDSAKeepFactory
} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

import "../../system/KeepFactorySelection.sol";

contract KeepFactorySelectionStub {
    using KeepFactorySelection for KeepFactorySelection.Storage;
    KeepFactorySelection.Storage keepFactorySelection;

    function initialize(IBondedECDSAKeepVendor _defaultVendor) public {
        keepFactorySelection.initialize(_defaultVendor);
    }

    function selectFactory() public view returns (IBondedECDSAKeepFactory) {
        return keepFactorySelection.selectFactory();
    }

    function selectFactoryAndRefresh() public returns (IBondedECDSAKeepFactory) {
        return keepFactorySelection.selectFactoryAndRefresh();
    }

    function setFullyBackedKeepVendor(address _fullyBackedVendor) public {
        keepFactorySelection.setFullyBackedKeepVendor(_fullyBackedVendor);
    }

    function setKeepFactorySelector(address _factorySelector) public {
        keepFactorySelection.setKeepFactorySelector(_factorySelector);
    }

    function keepStakeFactory() public view returns (address) {
        return address(keepFactorySelection.keepStakeFactory);
    }

    function ethStakeFactory() public view returns (address) {
        return address(keepFactorySelection.ethStakeFactory);
    }

    function factoriesVersionsLock() public view returns (bool) {
        return keepFactorySelection.factoriesVersionsLock;
    }

    function lockFactoriesVersions(
        address _expectedKeepStakeFactory,
        address _expectedFullyBackedFactory
    ) public {
        keepFactorySelection.lockFactoriesVersions(
            _expectedKeepStakeFactory,
            _expectedFullyBackedFactory
        );
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
