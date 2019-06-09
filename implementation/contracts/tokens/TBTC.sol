pragma solidity 0.4.25;

import {IBurnableERC20} from "../interfaces/IBurnableERC20.sol";

contract TBTCToken is IBurnableERC20 {

    bool returnBool = true;
    uint256 returnUint = 10 ** 18;

    function transfer(address to, uint256 value) external returns (bool) {to; value; return returnBool;}
    function approve(address spender, uint256 value) external returns (bool) {spender; value; return returnBool;}
    function transferFrom(address from, address to, uint256 value) external returns (bool) {from; to; value; return returnBool;}
    function totalSupply() external view returns (uint256) {return returnUint;}
    function balanceOf(address who) external view returns (uint256) {who; return returnUint;}
    function allowance(address owner, address spender) external view returns (uint256) {owner; spender; return returnUint;}
    function burnFrom(address from, uint256 value) external {from; value;}
    function burn(uint256 value) external {value;}
    function mint(address to, uint256 value) external returns (bool) {to; value;return returnBool;}

}
