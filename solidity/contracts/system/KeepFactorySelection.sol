pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";
import "./IKeepFactorySelector.sol";

library KeepFactorySelection {

    struct Storage {
        // Refresh requests counter;
        uint256 refreshRequestCounter;

        // Currently selected factory.
        IBondedECDSAKeepFactory selectedFactory;

        // Regular ECDSA keep factory.
        IBondedECDSAKeepFactory regularFactory;

        // Fully backed ECDSA keep factory.
        IBondedECDSAKeepFactory fullyBackedFactory;

        // Keep factory selector.
        IKeepFactorySelector factorySelector;
    }

    /// @notice Returns the currently selected keep factory.
    /// @return Selected keep factory.
    function selectFactory(
        Storage storage _self
    ) public view returns (IBondedECDSAKeepFactory) {
        if (address(_self.selectedFactory) == address(0)) {
            return _self.regularFactory;
        }

        return _self.selectedFactory;
    }

    /// @notice Returns the currently selected keep factory and
    /// performs selection of the new factory.
    /// @return Selected keep factory.
    function selectFactoryAndRefresh(
        Storage storage _self
    ) public returns (IBondedECDSAKeepFactory) {
        IBondedECDSAKeepFactory currentlySelectedFactory = selectFactory(_self);

        refreshFactory(_self);

        return currentlySelectedFactory;
    }

    /// @notice Performs selection of new keep factory.
    function refreshFactory(Storage storage _self) internal {
        if (
            address(_self.fullyBackedFactory) == address(0) ||
            address(_self.factorySelector) == address(0)
        ) {
            _self.selectedFactory = _self.regularFactory;
            return;
        }

        _self.refreshRequestCounter++;

        uint256 seed = uint256(
            keccak256(abi.encodePacked(address(this), _self.refreshRequestCounter))
        );

        _self.selectedFactory =  _self.factorySelector.selectFactory(
            seed,
            _self.regularFactory,
            _self.fullyBackedFactory
        );
    }
}
