pragma solidity ^0.5.10;

import {DepositFactory} from '../../../contracts/proxy/DepositFactory.sol';

contract TestDepositFactory is DepositFactory {

    constructor(address _systemAddress)
        DepositFactory(_systemAddress)
    public{}

    function setExternalDependencies(
        address _masterDepositAddress,
        address _tbtcSystem,
        address _tbtcToken,
        address _depositOwnerToken,
        address _feeRebateToken,
        address _vendingMachine,
        uint256 _keepThreshold,
        uint256 _keepSize
    ) public {
        masterDepositAddress = _masterDepositAddress;
        tbtcSystem = _tbtcSystem;
        tbtcToken = _tbtcToken;
        tbtcDepositToken = _depositOwnerToken;
        feeRebateToken = _feeRebateToken;
        vendingMachine = _vendingMachine;
        keepThreshold = _keepThreshold;
        keepSize = _keepSize;
    }

}