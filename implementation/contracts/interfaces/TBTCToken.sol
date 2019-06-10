pragma solidity 0.4.25;

import {IBurnableERC20} from './IBurnableERC20.sol';
import {SafeMath} from "../bitcoin-spv/SafeMath.sol";

/**
 * @title ERC20 implementation
 * @dev This is the TBTC token contract.
 */
contract TBTCToken is IBurnableERC20 {

    using SafeMath for uint256;

    /**
     * @dev Emitted when `_value` tokens are moved from one account (`_from`) to
     * another (`_to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    /**
     * @dev Emitted when the allowance of a `_spender` for an `_owner` is set by
     * a call to `_approve`. `_value` is the new allowance.
     */
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);

    /// @dev holds user balances
    mapping (address => uint256) private _balances;

    /// @dev holds user allowances
    mapping (address => mapping (address => uint256)) private _allowed;

    //the current total supply of the token
    uint256 private _totalSupply;

    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256) {
        return _totalSupply;
    }

     /**
     * @dev Moves `_value` tokens from the caller's account to `_to`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a `Transfer` event.
     */
    function transfer(address _to, uint256 _value) external returns (bool){
        _transfer(msg.sender, _to, _value);
        return true;
    }

    /**
     * @dev Sets `_value` as the allowance of `_spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * > Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an `Approval` event.
     */
    function approve(address _spender, uint256 _value) external returns (bool){
        _approve(msg.sender, _spender, _value);
        return true;
    }
     /**
     * @dev Moves `_value` tokens from `_from` to `_to` using the
     * allowance mechanism. `_value` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a `Transfer` event.
     */
    function transferFrom(address _from, address _to, uint256 _value) external returns (bool){
        _transfer(_from, _to, _value);
        _approve(_from, msg.sender, _allowed[_from][msg.sender].sub(_value));
        return true;
    }

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address _who) external view returns (uint256){
        return _balances[_who];
    }

    /**
     * @dev Returns the remaining number of tokens that `_spender` will be
     * allowed to spend on behalf of `_owner` through `transferFrom`. This is
     * zero by default.
     *
     * This value changes when `approve` or `transferFrom` are called.
     */
    function allowance(address _owner, address _spender) external view returns (uint256){
        return _allowed[_owner][_spender];
    }

    /**
     * @dev Destoys `_value` tokens from `_from`.`_value` is then deducted
     * from the caller's allowance.
     
     * See `_burnFrom`.
     */
    function burnFrom(address _from, uint256 _value) external{
        _burnFrom(_from, _value);
    }

     /**
     * @dev Destroys `amount` tokens from the caller.
     *
     * See `_burn`.
     */
    function burn(uint256 _value) external{
        _burn(msg.sender, _value);
    }

    function mint(address _to, uint256 _value) external returns (bool){
        //TODO: Minter Roles
        _mint(_to, _value);
        return true;
    }

    /** @dev Creates `_amount` tokens and assigns them to `_account`, increasing
     * the total supply.
     *
     * Emits a `Transfer` event with `_from` set to the zero address.
     *
     * Requirements
     *
     * - `_to` cannot be the zero address.
     */
    function _mint(address _account, uint256 _amount) internal {
        require(_account != address(0), "ERC20: mint to the zero address");

        _totalSupply = _totalSupply.add(_amount);
        _balances[_account] = _balances[_account].add(_amount);
        emit Transfer(address(0), _account, _amount);
    }
    
    /**
     * @dev Destoys `_amount` tokens from `_account`, reducing the
     * total supply.
     *
     * Emits a `Transfer` event with `_to` set to the zero address.
     *
     * Requirements
     *
     * - `_account` cannot be the zero address.
     * - `_account` must have at least `amount` tokens.
     */
    function _burn(address _account, uint256 _value) internal {
        require(_account != address(0), "ERC20: burn from the zero address");

        _totalSupply = _totalSupply.sub(_value);
        _balances[_account] = _balances[_account].sub(_value);
        emit Transfer(_account, address(0), _value);
    }

    /**
     * @dev Moves tokens `_amount` from `_sender` to `_recipient`.
     *
     * This is internal function is equivalent to `transfer`, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a `Transfer` event.
     *
     * Requirements:
     *
     * - `_sender` cannot be the zero address.
     * - `_recipient` cannot be the zero address.
     * - `_sender` must have a balance of at least `_amount`.
     */
    function _transfer(address _sender, address _recipient, uint256 _amount) internal {
        require(_sender != address(0), "ERC20: transfer from the zero address");
        require(_recipient != address(0), "ERC20: transfer to the zero address");

        _balances[_sender] = _balances[_sender].sub(_amount);
        _balances[_recipient] = _balances[_recipient].add(_amount);
        emit Transfer(_sender, _recipient, _amount);
    }

    /**
     * @dev Sets `_value` as the allowance of `_spender` over the `_owner`s tokens.
     *
     * This is internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an `Approval` event.
     *
     * Requirements:
     *
     * - `_owner` cannot be the zero address.
     * - `_spender` cannot be the zero address.
     */
    function _approve(address _owner, address _spender, uint256 _value) internal {
        require(_owner != address(0), "ERC20: approve from the zero address");
        require(_spender != address(0), "ERC20: approve to the zero address");

        _allowed[_owner][_spender] = _value;
        emit Approval(_owner, _spender, _value);
    }

     /**
     * @dev Destoys `_amount` tokens from `_account`.`_amount` is then deducted
     * from the caller's allowance.
     *
     * See `_burn` and `_approve`.
     */
    function _burnFrom(address _account, uint256 _amount) internal {
        _burn(_account, _amount);
        _approve(_account, msg.sender, _allowed[_account][msg.sender].sub(_amount));
    }
}
