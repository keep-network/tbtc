pragma solidity ^0.5.10;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract DepositOwnerToken is ERC721Metadata {
    address vendingMachine;
    constructor(address _vendingMachine) 
    ERC721Metadata("Deposit Owner Token", "DOT") 
    public {
        vendingMachine = _vendingMachine;
    }

    modifier onlyVendingMachine(){
        require(msg.sender == vendingMachine, "caller must be the vending machine");
        _;
    }
    function setTokenURI(uint256 tokenId, string memory uri)public onlyVendingMachine {
        _setTokenURI(tokenId, uri);
    }
}