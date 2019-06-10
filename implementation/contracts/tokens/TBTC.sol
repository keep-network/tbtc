pragma solidity 0.4.25;

import {ERC20Mintable} from "openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";
import {ERC20Burnable} from "openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
import {ERC20Detailed} from "openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";

contract TBTC is ERC20Detailed, ERC20Mintable, ERC20Burnable {
    constructor()
        ERC20Detailed("Trustless Bitcoin", "TBTC", 18)
        public
    {
    }
}