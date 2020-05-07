pragma solidity 0.5.17;

import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

library KeepFactoryStrategy {

    struct Storage {
        IBondedECDSAKeepFactory regularFactory;

        IBondedECDSAKeepFactory fullyBackedFactory;
    }

    function chooseFactory(
        Storage storage _self,
        address _application
    ) public view returns (IBondedECDSAKeepFactory) {
        require(
            address(_self.regularFactory) != address(0),
            "Regular factory address must be set"
        );

        if (address(_self.fullyBackedFactory) == address(0)) {
            return _self.regularFactory;
        }

        uint256 regularFactoryWeight = _self.regularFactory.getSortitionPoolWeight(_application);
        uint256 fullyBackedFactoryWeight = _self.fullyBackedFactory.getSortitionPoolWeight(_application);

        if (shouldUseRegularFactory(regularFactoryWeight, fullyBackedFactoryWeight)) {
            return _self.regularFactory;
        } else {
            return _self.fullyBackedFactory;
        }
    }

    function shouldUseRegularFactory(
        uint256 _regularFactoryWeight,
        uint256 _fullyBackedFactoryWeight
    ) internal view returns (bool) {
        // TODO: Implementation of the actual strategy.
        return true;
    }
}
