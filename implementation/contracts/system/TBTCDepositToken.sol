pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./ERC721MinterAuthority.sol";


/// @title tBTC Deposit Token for tracking deposit ownership
/// @notice The tBTC Deposit Token, commonly referenced as the TDT, is an
///         ERC721 non-fungible token whose ownership reflects the ownership
///         of its corresponding deposit. Each deposit has one TDT, and vice
///         versa. Owning a TDT is equivalent to owning its corresponding
///         deposit. TDTs can be transferred freely. tBTC's VendingMachine
///         contract takes ownership of TDTs and in exchange returns fungible
//          TBTC tokens whose value is backed 1-to-1 by the corresponding
//          deposit's BTC.
/// @dev Currently, TDTs are minted using the uint256 casting of the
///      corresponding deposit contract's address. That is, the TDTs id is
///      convertible to the deposit's address and vice versa.
contract TBTCDepositToken is ERC721Metadata, ERC721MinterAuthority {

    constructor(address _depositFactory) 
        ERC721Metadata("tBTC Deopsit Token", "TDT")
        ERC721MinterAuthority(_depositFactory) 
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    /// @dev Mints a new token.
    /// Reverts if the given token ID already exists.
    /// @param _to The address that will own the minted token
    /// @param _tokenId uint256 ID of the token to be minted
    function mint(address _to, uint256 _tokenId) public onlyFactory {
        _mint(_to, _tokenId);
    }

    /// @dev Returns whether the specified token exists.
    /// @param _tokenId uint256 ID of the token to query the existence of
    /// @return bool whether the token exists
    function exists(uint256 _tokenId) public view returns (bool) {
        return _exists(_tokenId);
    }
}
