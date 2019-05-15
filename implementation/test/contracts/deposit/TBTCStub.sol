pragma solidity 0.4.25;

import {IBurnableERC20} from '../../../contracts/interfaces/IBurnableERC20.sol';

contract TBTCStub is IBurnableERC20 {

    bool returnBool = true;
    uint256 returnUint = 10 ** 18;
    mapping (address => uint) internal balances;
    function setReturnBool(bool _res) public {returnBool = _res;}
    function setReturnUint(uint _res) public {returnUint = _res;}
    function transfer(address to, uint256 value) external returns (bool) {
        to; value;
         if(balances[to] + value > balances[to]){
            balances[to] += value;
        }
     return returnBool;}
    function approve(address spender, uint256 value) external returns (bool) {spender; value; return returnBool;}
    function transferFrom(address from, address to, uint256 value) external returns (bool) {from; to; value; return returnBool;}
    function totalSupply() external view returns (uint256) {return returnUint;}
    function balanceOf(address who) external view returns (uint256) {who; return returnUint;}
    function allowance(address owner, address spender) external view returns (uint256) {owner; spender; return returnUint;}
    function getBalance(address _of) external view returns (uint256) {
        return balances[_of];
    }
    function burnFrom(address from, uint256 value) external {
        if(balances[from] - value < balances[from]){
            balances[from] -= value;
        }
    }
    function burn(uint256 value) external {
        if(balances[msg.sender] - value < balances[msg.sender]){
            balances[msg.sender] -= value;
        }
    }
    function mint(address to, uint256 value) external returns (bool) {balances[to] += value; return returnBool;}

}
