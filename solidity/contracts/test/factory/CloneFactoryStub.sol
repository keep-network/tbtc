pragma solidity 0.5.17;

import {CloneFactory} from "../../../contracts/proxy/CloneFactory.sol";

contract CloneFactoryStub is CloneFactory {

    event ContractCloneCreated(address contractCloneAddress);

    function createClone_exposed(address _target) external returns (address){
        address cloneAddress = createClone(_target);
        emit ContractCloneCreated(cloneAddress);
    }

    function isClone_exposed(address _target, address _query)external view returns (bool result){
        return isClone(_target, _query);
    }
}
