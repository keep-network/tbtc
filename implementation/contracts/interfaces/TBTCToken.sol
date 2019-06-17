pragma solidity 0.4.25;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Mintable.sol";

contract TBTCToken is ERC20Detailed, ERC20Mintable {
    constructor(address _newMinter) ERC20Detailed("Trustless bitcoin", "TBTC", 18) public{
        addMinter(_newMinter);
    }
     function burnFrom(address account, uint256 amount) public onlyMinter{
        _burn(account, amount);
    }
}
