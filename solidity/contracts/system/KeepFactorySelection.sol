pragma solidity 0.5.17;

import {IBondedECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";
import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

/// @title Bonded ECDSA keep factory selection strategy.
/// @notice The strategy defines the algorithm for selecting a factory. tBTC
/// uses two bonded ECDSA keep factories, selecting one of them for each new
/// deposit being opened.
interface KeepFactorySelector {

    /// @notice Selects keep factory for the new deposit.
    /// @param _seed Request seed.
    /// @param _keepStakeFactory Regular, KEEP-stake based keep factory.
    /// @param _fullyBackedFactory Fully backed, ETH-bond-only based keep factory.
    /// @return The selected keep factory.
    function selectFactory(
        uint256 _seed,
        IBondedECDSAKeepFactory _keepStakeFactory,
        IBondedECDSAKeepFactory _fullyBackedFactory
    ) external view returns (IBondedECDSAKeepFactory);
}

/// @title Bonded ECDSA keep factory selection library.
/// @notice tBTC uses two bonded ECDSA keep factories: one based on KEEP stake
/// and ETH bond and another based only on ETH bond. Factories addresses
/// are obtained through calls to respective vendor contracts. The library holds
/// a reference to both vendors and factories as well as a reference to a selection
/// strategy deciding which factory to choose for the new deposit being opened.
/// Factories addresses can be locked, so the vendors won't be called anymore to
/// get the latest factories addresses.
library KeepFactorySelection {

    struct Storage {
        uint256 requestCounter;

        IBondedECDSAKeepFactory selectedFactory;

        KeepFactorySelector factorySelector;

        // Standard ECDSA keep vendor and factory: KEEP stake and ETH bond.
        // Guaranteed to be set for initialized factory.
        IBondedECDSAKeepVendor keepStakeVendor;
        IBondedECDSAKeepFactory keepStakeFactory;

        // Fully backed ECDSA keep vendor and factory: ETH bond only.
        IBondedECDSAKeepVendor fullyBackedVendor;
        IBondedECDSAKeepFactory fullyBackedFactory;

        // Lock for factories versions freeze. When set to true vendor won't be
        // called to obtain a new factory address but a version of a factory
        // from the moment of freezing will be used. Once locked the factory
        // address won't be able to update anymore.
        bool factoriesVersionsLock;
    }

    /// @notice Initializes the library with the default KEEP-stake-based
    /// vendor. The default vendor is guaranteed to be set and this function
    /// must be called when creating contract using this library. It calls the
    /// vendor to obtain the default factory address.
    /// @dev This function can be called only one time.
    function initialize(
        Storage storage _self,
        IBondedECDSAKeepVendor _defaultVendor
    ) public {
        require(
            address(_self.keepStakeVendor) == address(0),
            "Already initialized"
        );

        _self.keepStakeVendor = IBondedECDSAKeepVendor(_defaultVendor);
        _self.selectedFactory = getKeepStakedFactory(_self);
    }

    /// @notice Returns the selected keep factory.
    /// This function guarantees that the same factory is returned for every
    /// call until selectFactoryAndRefresh is executed. This lets to evaluate
    /// open keep fee estimate on the same factory that will be used later for
    /// opening a new keep (fee estimate and open keep requests are two
    /// separate calls).
    /// @return Selected keep factory. The same vale will be returned for every
    /// call of this function until selectFactoryAndRefresh is executed.
    function selectFactory(
        Storage storage _self
    ) public view returns (IBondedECDSAKeepFactory) {
        return _self.selectedFactory;
    }

    /// @notice Returns the selected keep factory and refreshes the choice
    /// for the next select call. The value returned by this function has been
    /// evaluated during the previous call. This lets to return the same value
    /// from selectFactory and selectFactoryAndRefresh, thus, allowing to use
    /// the same factory for which open keep fee estimate was evaluated (fee
    /// estimate and open keep requests are two separate calls).
    /// @return Selected keep factory.
    function selectFactoryAndRefresh(
        Storage storage _self
    ) public returns (IBondedECDSAKeepFactory) {
        IBondedECDSAKeepFactory factory = selectFactory(_self);
        refreshFactory(_self);

        return factory;
    }

    /// @notice Sets the minimum bondable value required from the operator to
    /// join the sortition pool for tBTC.
    /// @param _minimumBondableValue The minimum bond value the application
    /// requires from a single keep.
    /// @param _groupSize Number of signers in the keep.
    /// @param _honestThreshold Minimum number of honest keep signers.
    function setMinimumBondableValue(
        Storage storage _self,
        uint256 _minimumBondableValue,
        uint256 _groupSize,
        uint256 _honestThreshold
    ) public {
        if (address(_self.keepStakeFactory) != address(0)) {
            _self.keepStakeFactory.setMinimumBondableValue(
                _minimumBondableValue,
                _groupSize,
                _honestThreshold
            );
        }
        if (address(_self.fullyBackedFactory) != address(0)) {
            _self.fullyBackedFactory.setMinimumBondableValue(
                _minimumBondableValue,
                _groupSize,
                _honestThreshold
            );
        }
    }

    /// @notice Refreshes the keep factory choice. If either ETH-bond-only vendor
    /// or selection strategy is not set, KEEP-stake factory is selected.
    /// Otherwise, calls selection strategy providing addresses of both
    /// factories to make a choice. Additionally, passes the selection seed
    /// evaluated from the current request counter value.
    /// Unless lock is set, it calls vendors to obtain the latests factories
    /// versions.
    function refreshFactory(Storage storage _self) internal {
        IBondedECDSAKeepFactory keepStakeFactory = getKeepStakedFactory(_self);

        if (
            address(_self.fullyBackedVendor) == address(0) ||
            address(_self.factorySelector) == address(0)
        ) {
            // KEEP-stake factory is guaranteed to be there. If the selection
            // can not be performed, this is the default choice.
            _self.selectedFactory = keepStakeFactory;
            return;
        }

        IBondedECDSAKeepFactory fullyBackedFactory = getFullyBackedFactory(_self);

        _self.requestCounter++;
        uint256 seed = uint256(
            keccak256(abi.encodePacked(address(this), _self.requestCounter))
        );
        _self.selectedFactory = _self.factorySelector.selectFactory(
            seed,
            keepStakeFactory,
            fullyBackedFactory
        );

        require(
            _self.selectedFactory == keepStakeFactory ||
                _self.selectedFactory == fullyBackedFactory,
            "Factory selector returned unknown factory"
        );
    }

    /// @notice Returns KEEP stake based factory address. If factories lock is not set
    /// it calls the KEEP staked vendor to obtain the latest versions of the factory.
    function getKeepStakedFactory(Storage storage _self)
        internal
        view
        returns (IBondedECDSAKeepFactory)
    {
        if (_self.factoriesVersionsLock) {
            return _self.keepStakeFactory;
        } else {
            IBondedECDSAKeepFactory keepStakeFactory = IBondedECDSAKeepFactory(
                _self.keepStakeVendor.selectFactory()
            );

            require(
                address(keepStakeFactory) != address(0),
                "Vendor returned invalid factory address"
            );

            return keepStakeFactory;
        }
    }

    /// @notice Returns ETH-bond-only based factory address. If factories lock is not set
    /// it calls the ETH staked vendor to obtain the latest versions of the factory.
    function getFullyBackedFactory(Storage storage _self)
        internal
        view
        returns (IBondedECDSAKeepFactory)
    {
        if (_self.factoriesVersionsLock) {
            return _self.fullyBackedFactory;
        } else {
            IBondedECDSAKeepFactory fullyBackedFactory = IBondedECDSAKeepFactory(
                _self.fullyBackedVendor.selectFactory()
            );

            require(
                address(fullyBackedFactory) != address(0),
                "Vendor returned invalid factory address"
            );

            return fullyBackedFactory;
        }
    }

    /// @notice Sets the address of the fully backed, ETH-bond-only based keep
    /// factory. KeepFactorySelection can work without the fully-backed keep
    /// factory set, always selecting the default KEEP-stake-based factory.
    /// Once both fully-backed keep factory and factory selection strategy are
    /// set, KEEP-stake-based factory is no longer the default choice and it is
    /// up to the selection strategy to decide which factory should be chosen.
    /// @dev Can be called only one time! The function calls the vendor to confirm
    /// it correctly returns a factory address.
    /// @param _fullyBackedVendor Address of the fully-backed, ETH-bond-only based
    /// keep vendor.
    function setFullyBackedKeepVendor(
        Storage storage _self,
        address _fullyBackedVendor
    ) internal {
        require(
            address(_self.fullyBackedVendor) == address(0),
            "Fully backed vendor already set"
        );
        require(_fullyBackedVendor != address(0), "Invalid address");

        IBondedECDSAKeepVendor fullyBackedVendor = IBondedECDSAKeepVendor(
            _fullyBackedVendor
        );

        require(
            fullyBackedVendor.selectFactory() != address(0),
            "Vendor returned invalid factory address"
        );

        _self.fullyBackedVendor = fullyBackedVendor;
    }

    /// @notice Sets the address of the keep factory selection strategy contract.
    /// KeepFactorySelection can work without the keep factory selection
    /// strategy set, always selecting the default KEEP-stake-based factory.
    /// Once both fully-backed keep factory and factory selection strategy are
    /// set, KEEP-stake-based factory is no longer the default choice and it is
    /// up to the selection strategy to decide which factory should be chosen.
    /// @dev Can be called only one time!
    /// @param _factorySelector Address of the keep factory selection strategy.
    function setKeepFactorySelector(
        Storage storage _self,
        address _factorySelector
    ) internal {
        require(
            address(_self.factorySelector) == address(0),
            "Factory selector already set"
        );
        require(
            address(_factorySelector) != address(0),
            "Invalid address"
        );

        _self.factorySelector = KeepFactorySelector(_factorySelector);
    }

    /// @notice Locks versions of factories. When lock is set vendor contracts
    /// won't be called anymore to obtain the latest factories versions.
    /// It requires expected factories addresses to be provided to protect from
    /// locking on unexpected addresses.
    /// @param _expectedKeepStakeFactory Expected KEEP-staked factory address
    /// @param _expectedFullyBackedFactory Expected ETH-bond-only factory address
    function lockFactoriesVersions(
        Storage storage _self,
        address _expectedKeepStakeFactory,
        address _expectedFullyBackedFactory
    ) internal {
        require(!_self.factoriesVersionsLock, "Already locked");

        require(
            address(_self.keepStakeVendor) != address(0),
            "KEEP backed vendor not set"
        );
        require(
            address(_self.fullyBackedVendor) != address(0),
            "Fully backed vendor not set"
        );

        IBondedECDSAKeepFactory latestKeepStakeFactory = getKeepStakedFactory(
            _self
        );
        IBondedECDSAKeepFactory latestFullyBackedFactory = getFullyBackedFactory(
            _self
        );

        require(
            address(latestKeepStakeFactory) == _expectedKeepStakeFactory,
            "Unexpected KEEP backed factory"
        );
        require(
            address(latestFullyBackedFactory) == _expectedFullyBackedFactory,
            "Unexpected fully backed factory"
        );

        _self.keepStakeFactory = latestKeepStakeFactory;
        _self.fullyBackedFactory = latestFullyBackedFactory;

        _self.factoriesVersionsLock = true;
    }
}
