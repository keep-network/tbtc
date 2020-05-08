pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

/// @title Keep factory selector.
/// @notice Selects an appropriate keep factory.
interface IKeepFactorySelector {

    /// @notice Selects keep factory.
    /// @param _regularFactory Regular keep factory.
    /// @param _fullyBackedFactory Fully backed keep factory.
    /// @return Selected keep factory.
    function selectFactory(
        IBondedECDSAKeepFactory _regularFactory,
        IBondedECDSAKeepFactory _fullyBackedFactory
    ) external view returns (IBondedECDSAKeepFactory);
}
