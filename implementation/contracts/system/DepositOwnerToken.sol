pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract DepositOwnerToken is ERC721Metadata {

    constructor() ERC721Metadata("Deposit Owner Token", "DOT") public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev Returns whether the specified token exists.
    /// @param tokenId uint256 ID of the token to query the existence of
    /// @return bool whether the token exists
    function exists(uint256 tokenId) public view returns (bool) {
        return _exists(tokenId);
    }
}