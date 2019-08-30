/* solium-disable function-order */
pragma solidity ^0.5.10;

import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {DepositLog} from "../DepositLog.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import {KeepRegistry} from "keep-tecdsa/solidity/contracts/KeepRegistry.sol";
import {ECDSAKeepVendor} from "keep-tecdsa/solidity/contracts/ECDSAKeepVendor.sol";

contract TBTCSystem is ITBTCSystem, ERC721, DepositLog {
    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;

    address keepRegistry;

    //TODO: add: onlyOwner
    function initialize(address _keepRegistry) public {
        keepRegistry = _keepRegistry;
    }

    // Price Oracle
    function fetchOraclePrice() external view returns (uint256) {return oraclePrice;}

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
    /// @param _m Minimum number of honest keep members.
    /// @param _n Number of members in the keep.
    /// @return Address of a new keep.
    function requestNewKeep(uint256 _m, uint256 _n) external payable returns (address _keepAddress){
        address keepVendorAddress = KeepRegistry(keepRegistry)
            .getVendor("ECDSAKeep");

        _keepAddress = ECDSAKeepVendor(keepVendorAddress)
            .openKeep(_n,_m, msg.sender);
    }

    // ERC721

    /// @dev             Function to mint a new token.
    ///                  Reverts if the given token ID already exists.
    /// @param _to       The address that will own the minted token
    /// @param _tokenId  uint256 ID of the token to be minted
    function mint(address _to, uint256 _tokenId) public {
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
