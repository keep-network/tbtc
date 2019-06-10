pragma solidity 0.4.25;

/**
 * @title ERC20 interface
 * @dev see https://eips.ethereum.org/EIPS/eip-20
 */
interface IBurnableERC20 {

    function totalSupply() external view returns (uint256);

    function balanceOf(address _who) external view returns (uint256);

    function transfer(address _to, uint256 _value) external returns (bool);

    function allowance(address _owner, address _spender) external view returns (uint256);

    function approve(address _spender, uint256 _value) external returns (bool);

    function transferFrom(address _from, address _to, uint256 _value) external returns (bool);

    function burnFrom(address _from, uint256 _value) external;

    function burn(uint256 _value) external;

    function mint(address _to, uint256 _value) external returns (bool);

    event Transfer(address indexed _from, address indexed _to, uint256 _value);

    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}
