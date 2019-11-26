/* solium-disable function-order */
pragma solidity ^0.5.10;

import {IKeepRegistry} from "@keep-network/keep-ecdsa/contracts/api/IKeepRegistry.sol";
import {IECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeepVendor.sol";

import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {DepositLog} from "../DepositLog.sol";

import "./DepositOwnerToken.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./ERC721MinterAuthority.sol";

contract TBTCSystem is Ownable, ITBTCSystem, DepositOwnerToken, ERC721MinterAuthority, DepositLog {

    bool _initialized = false;

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;

    address public keepRegistry;
    address public tbtcUniswapExchange;

    constructor(address _depositFactory)
        ERC721MinterAuthority(_depositFactory)
        public
    {
            // solium-disable-previous-line no-empty-blocks
    }

    function initialize(
        address _keepRegistry,
        address _tbtcUniswapExchange
    ) external onlyOwner {
        require(!_initialized, "already initialized");

        keepRegistry = _keepRegistry;
        tbtcUniswapExchange = _tbtcUniswapExchange;

        _initialized = true;
    }

    function getTBTCUniswapExchange() external view returns (address) {
        return tbtcUniswapExchange;
    }

    // Price Oracle
    function fetchOraclePrice() external view returns (uint256) {
        return oraclePrice;
    }

    // Difficulty Oracle
    // TODO: This is a workaround. It will be replaced by tbtc-difficulty-oracle.
    function fetchRelayCurrentDifficulty() external view returns (uint256) {
        return currentDifficulty;
    }

    function fetchRelayPreviousDifficulty() external view returns (uint256) {
        return previousDifficulty;
    }

    function submitCurrentDifficulty(uint256 _currentDifficulty) public {
        if (currentDifficulty != _currentDifficulty) {
            previousDifficulty = currentDifficulty;
            currentDifficulty = _currentDifficulty;
        }
    }

    /// @notice Request a new keep opening.
    /// @param _m Minimum number of honest keep members required to sign.
    /// @param _n Number of members in the keep.
    /// @return Address of a new keep.
    function requestNewKeep(uint256 _m, uint256 _n)
        external
        payable
        returns (address _keepAddress)
    {
        address keepVendorAddress = IKeepRegistry(keepRegistry)
            .getVendor("ECDSAKeep");

        _keepAddress = IECDSAKeepVendor(keepVendorAddress)
            .openKeep(_n,_m, msg.sender);
    }

    /// @notice          Function to mint a new token.
    /// @dev             Reverts if the given token ID already exists.
    ///                  This function can only be called by depositFactory
    /// @param _to       The address that will own the minted token
    /// @param _tokenId  uint256 ID of the token to be minted
    function mint(address _to, uint256 _tokenId) public onlyFactory {
        _mint(_to, _tokenId);
    }

    /// @notice  Checks if an address is a deposit.
    /// @dev     Verifies if Deposit ERC721 token with given address exists.
    /// @param _depositAddress  The address to check
    /// @return  True if deposit with given value exists, false otherwise.
    function isDeposit(address _depositAddress) public returns (bool){
        return _exists(uint256(_depositAddress));
    }
}
