pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";
import "./IKeepFactorySelector.sol";

library KeepFactoryStrategy {

    struct Storage {
        IBondedECDSAKeepFactory regularFactory;

        IBondedECDSAKeepFactory fullyBackedFactory;

        IKeepFactorySelector factorySelector;
    }

    /// @notice Selects the right keep factory for the given application.
    /// @param _application Application address.
    /// @return Selected keep factory.
    function selectFactory(
        Storage storage _self,
        address _application
    ) public view returns (IBondedECDSAKeepFactory) {
        require(
            address(_self.regularFactory) != address(0),
            "Regular factory address must be set"
        );

        if (
            address(_self.fullyBackedFactory) == address(0) ||
            address(_self.factorySelector) == address(0)
        ) {
            return _self.regularFactory;
        }

        return _self.factorySelector.selectFactory(
            _application,
            _self.regularFactory,
            _self.fullyBackedFactory
        );
    }
}
