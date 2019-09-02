pragma solidity 0.5.10;

import {TBTCSystem} from "./TBTCSystem.sol";

contract MinterAuthority {

    address internal _TBTCSystem;

    /// @notice Set the address of the System contract on contract initialization
    constructor(address _tbtcSystem) public {
        _TBTCSystem = _tbtcSystem;
    }

    /// @notice Function modifier ensures modified function caller address exists as an ERC721 token
    modifier onlyDeposit(){
        require(TBTCSystem(_TBTCSystem).isDeposit(msg.sender), "caller must be a deposit");
        _;
    }
}