pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

/// @title Bonded ECDSA keep factory selection strategy.
/// @notice The strategy defines the algorithm for selecting a factory. tBTC
/// uses two bonded ECDSA keep factories, selecting one of them for each new
/// deposit being opened.
interface KeepFactorySelector {

    /// @notice Selects keep factory for the new deposit.
    /// @param _seed Request seed.
    /// @param _keepStakeFactory Regular, KEEP-stake based keep factory.
    /// @param _ethStakeFactory Fully backed, ETH-stake based keep factory.
    /// @return The selected keep factory.
    function selectFactory(
        uint256 _seed,
        IBondedECDSAKeepFactory _keepStakeFactory,
        IBondedECDSAKeepFactory _ethStakeFactory
    ) external view returns (IBondedECDSAKeepFactory);
}

/// @title Bonded ECDSA keep factory selection library.
/// @notice tBTC uses two bonded ECDSA keep factories: one based on KEEP stake
/// and ETH bond, and another based on ETH stake and ETH bond. The library holds
/// a reference to both factories as well as a reference to a selection strategy
/// deciding which factory to choose for the new deposit being opened.
library KeepFactorySelection {

    struct Storage {
        uint256 requestCounter;

        IBondedECDSAKeepFactory selectedFactory;

        KeepFactorySelector factorySelector;

        // Standard ECDSA keep factory: KEEP stake and ETH bond.
        IBondedECDSAKeepFactory keepStakeFactory;

        // Fully backed ECDSA keep factory: ETH stake and ETH bond.
        IBondedECDSAKeepFactory ethStakeFactory;
    }

    /// @notice Returns the currently selected keep factory.
    /// @return Selected keep factory.
    function selectFactory(
        Storage storage _self
    ) public view returns (IBondedECDSAKeepFactory) {
        if (address(_self.selectedFactory) == address(0)) {
            return _self.keepStakeFactory;
        }

        return _self.selectedFactory;
    }

    /// @notice Returns the currently selected keep factory and
    /// performs selection of the new factory.
    /// @return Selected keep factory.
    function selectFactoryAndRefresh(
        Storage storage _self
    ) public returns (IBondedECDSAKeepFactory) {
        IBondedECDSAKeepFactory factory = selectFactory(_self);
        refreshFactory(_self);

        return factory;
    }

    /// @notice Performs selection of new keep factory.
    function refreshFactory(Storage storage _self) internal {
        if (
            address(_self.ethStakeFactory) == address(0) ||
            address(_self.factorySelector) == address(0)
        ) {
            _self.selectedFactory = _self.keepStakeFactory;
            return;
        }

        _self.requestCounter++;
        uint256 seed = uint256(
            keccak256(abi.encodePacked(address(this), _self.requestCounter))
        );
        _self.selectedFactory = _self.factorySelector.selectFactory(
            seed,
            _self.keepStakeFactory,
            _self.ethStakeFactory
        );
    }
}
