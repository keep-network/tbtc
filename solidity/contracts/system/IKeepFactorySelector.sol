pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

/// @title Keep factory selector.
/// @notice Selects an appropriate keep factory for the given application.
interface IKeepFactorySelector {

    /// @notice Selects keep factory for the given application.
    /// @param _application Application address.
    /// @param _regularFactory Regular keep factory.
    /// @param _fullyBackedFactory Fully backed keep factory.
    /// @return Selected keep factory.
    function selectFactory(
        address _application,
        IBondedECDSAKeepFactory _regularFactory,
        IBondedECDSAKeepFactory _fullyBackedFactory
    ) external view returns (IBondedECDSAKeepFactory);
}
