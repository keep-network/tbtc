pragma solidity 0.4.25;

import {IBurnableERC20} from '../interfaces/IBurnableERC20.sol';
import {SafeMath} from "../bitcoin-spv/SafeMath.sol";

/**
 * @title ERC20 implementation
 * @dev This is the TBTC token contract.
 */
contract TBTCToken is IBurnableERC20 {

    using SafeMath for uint256;

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    mapping (address => uint256) private _balances;
    mapping (address => mapping (address => uint256)) private _allowed;

    uint256 private _totalSupply;

    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

    function transfer(address _to, uint256 _value) external returns (bool){
        _transfer(msg.sender, _to, _value);
        return true;
    }

    function approve(address _spender, uint256 _value) external returns (bool){
        _approve(msg.sender, _spender, _value);
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) external returns (bool){
        _transfer(_from, _to, _value);
        _approve(_from, msg.sender, _allowed[_from][msg.sender].sub(_value));
        return true;
    }

    function balanceOf(address _who) external view returns (uint256){
        return _balances[_who];
    }

    function allowance(address _owner, address _spender) external view returns (uint256){
        return _allowed[_owner][_spender];
    }

    function burnFrom(address _from, uint256 _value) external{
        _burnFrom(_from, _value);
    }

    function burn(uint256 _value) external{
        _burn(msg.sender, _value);
    }

    function mint(address _to, uint256 _value) external returns (bool){
        //TODO: Minter Roles
        _mint(_to, _value);
        return true;
    }

    function _mint(address _account, uint256 _amount) internal {
        require(_account != address(0), "ERC20: mint to the zero address");

        _totalSupply = _totalSupply.add(_amount);
        _balances[_account] = _balances[_account].add(_amount);
        emit Transfer(address(0), _account, _amount);
    }
 
    function _burn(address _account, uint256 _value) internal {
        require(_account != address(0), "ERC20: burn from the zero address");

        _totalSupply = _totalSupply.sub(_value);
        _balances[_account] = _balances[_account].sub(_value);
        emit Transfer(_account, address(0), _value);
    }

    function _transfer(address _sender, address _recipient, uint256 _amount) internal {
        require(_sender != address(0), "ERC20: transfer from the zero address");
        require(_recipient != address(0), "ERC20: transfer to the zero address");

        _balances[_sender] = _balances[_sender].sub(_amount);
        _balances[_recipient] = _balances[_recipient].add(_amount);
        emit Transfer(_sender, _recipient, _amount);
    }

    function _approve(address _owner, address _spender, uint256 _value) internal {
        require(_owner != address(0), "ERC20: approve from the zero address");
        require(_spender != address(0), "ERC20: approve to the zero address");

        _allowed[_owner][_spender] = _value;
        emit Approval(_owner, _spender, _value);
    }

    function _burnFrom(address _account, uint256 _amount) internal {
        _burn(_account, _amount);
        _approve(_account, msg.sender, _allowed[_account][msg.sender].sub(_amount));
    }
}
