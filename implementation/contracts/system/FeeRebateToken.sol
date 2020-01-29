pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./VendingMachineAuthority.sol";

contract FeeRebateToken is ERC721Metadata, VendingMachineAuthority {

    constructor(address _vendingMachine)
        ERC721Metadata("Fee Rebate Token", "FRT")
        VendingMachineAuthority(_vendingMachine)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev Mints a new token.
    /// Reverts if the given token ID already exists.
    /// @param _to The address that will own the minted token
    /// @param _tokenId uint256 ID of the token to be minted
    function mint(address _to, uint256 _tokenId) public onlyVendingMachine {
        _mint(_to, _tokenId);
    }

    /// @dev Returns whether the specified token exists.
    /// @param _tokenId uint256 ID of the token to query the existence of
    /// @return bool whether the token exists
    function exists(uint256 _tokenId) public view returns (bool) {
        return _exists(_tokenId);
    }
}