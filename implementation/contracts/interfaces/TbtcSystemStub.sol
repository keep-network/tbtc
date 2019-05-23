pragma solidity 0.4.25;

import {ITBTCSystem} from '../interfaces/ITBTCSystem.sol';
import {IERC721} from '../interfaces/IERC721.sol';
import {DepositLog} from '../../contracts/DepositLog.sol';

contract TBTCSystemStub is ITBTCSystem, IERC721, DepositLog {

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;
    address depositOwner = address(0);

    // DepositLog
    // Override parent function until authorization is available
    function approvedToLog(address _caller) public view returns (bool) {_caller; return true;}

    // TBTCSystem
    function fetchOraclePrice() external view returns (uint256) {return oraclePrice;}
    function fetchRelayCurrentDifficulty() external view returns (uint256) {return currentDifficulty;}
    function fetchRelayPreviousDifficulty() external view returns (uint256) {return previousDifficulty;}

    // ERC721
    function balanceOf(address owner) public view returns (uint256 balance) {owner; balance = 0;}
    function ownerOf(uint256 tokenId) public view returns (address owner) {tokenId; owner = depositOwner;}
    function approve(address to, uint256 tokenId) public {to; tokenId;}
    function getApproved(uint256 tokenId) public view returns (address operator) {tokenId; operator = address(8);}
    function setApprovalForAll(address operator, bool _approved) public {operator; _approved;}
    function isApprovedForAll(address owner, address operator) public view returns (bool) {owner; operator;}
    function transferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public {from; to; tokenId; data;}
}
