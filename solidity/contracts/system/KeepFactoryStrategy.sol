pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";
import "./IKeepFactorySelector.sol";

library KeepFactoryStrategy {

    struct Storage {
        IBondedECDSAKeepFactory regularFactory;

        IBondedECDSAKeepFactory fullyBackedFactory;

        IKeepFactorySelector factorySelector;
    }

    /// @notice Selects the right keep factory.
    /// @return Selected keep factory.
    function selectFactory(
        Storage storage _self
    ) public view returns (IBondedECDSAKeepFactory) {
        if (
            address(_self.fullyBackedFactory) == address(0) ||
            address(_self.factorySelector) == address(0)
        ) {
            return _self.regularFactory;
        }

        return _self.factorySelector.selectFactory(
            _self.regularFactory,
            _self.fullyBackedFactory
        );
    }
}
