pragma solidity 0.4.25;

import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20Detailed.sol";
import "../../node_modules/openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";

contract TBTCToken is ERC20Detailed, ERC20 {
    address TBTCSystem;

    constructor(address _TBTCSystem) ERC20Detailed("Trustless bitcoin", "TBTC", 18) public {
        TBTCSystem = _TBTCSystem;
    }

    modifier onlySystem(){
        require(msg.sender == TBTCSystem);
        _;
    }

    function systemMint(address _account, uint256 _amount) public onlySystem {
        _mint(_account, _amount);
    }

    function systemBurnFrom(address _account, uint256 _amount) public onlySystem {
        _burn(_account, _amount);
    }

    function systemTransferFrom(address _from, address _to, uint256 _value) public onlySystem {
        _transfer(_from, _to, _value);
    }
}
