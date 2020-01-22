pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

/**
 @dev Interface of recipient contract for approveAndCall pattern.
*/
interface tokenRecipient { function receiveApproval(address _from, uint256 _tokenId, address _token, bytes calldata _extraData) external; }

contract DepositOwnerToken is ERC721Metadata {

    constructor() ERC721Metadata("Deposit Owner Token", "DOT") public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev Mints a new token.
    /// Reverts if the given token ID already exists.
    /// @param _to The address that will own the minted token
    /// @param _tokenId uint256 ID of the token to be minted
    function mint(address _to, uint256 _tokenId) public {
        _mint(_to, _tokenId);
    }

    /// @dev Returns whether the specified token exists.
    /// @param _tokenId uint256 ID of the token to query the existence of
    /// @return bool whether the token exists
    function exists(uint256 _tokenId) public view returns (bool) {
        return _exists(_tokenId);
    }

    /// @notice           Set allowance for other address and notify.
    ///                   Allows `_spender` to transfer the token with ID, `_tokenId`.
    ///                   on your behalf and then ping the contract about it.
    /// @dev              The `_spender` should implement the `tokenRecipient` interface above
    ///                   to receive approval notifications.
    /// @param _spender   Address of contract authorized to spend.
    /// @param _tokenId   uint256 ID of the token to approve access to
    /// @param _extraData Extra information to send to the approved contract.
    function approveAndCall(address _spender, uint256 _tokenId, bytes memory _extraData) public returns (bool success) {
        tokenRecipient spender = tokenRecipient(_spender);
        approve(_spender, _tokenId);
        spender.receiveApproval(msg.sender, _tokenId, address(this), _extraData);
        return true;
    }
}