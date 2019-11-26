pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract DepositOwnerToken is ERC721Metadata {

    constructor() ERC721Metadata("Deposit Owner Token", "DOT") public {
        // solium-disable-previous-line no-empty-blocks
    }
}