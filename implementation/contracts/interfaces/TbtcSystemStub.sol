pragma solidity 0.4.25;

import {ITBTCSystem} from '../interfaces/ITBTCSystem.sol';
import {IERC721} from '../interfaces/IERC721.sol';
import {DepositLog} from '../../contracts/DepositLog.sol';

contract TBTCSystemStub is ITBTCSystem, IERC721, DepositLog {

    uint256 _currentDifficulty = 1;
    uint256 _previousDifficulty = 1;
    uint256 _oraclePrice = 10 ** 12;
    address _depositOwner = address(0);

    // DepositLog
    // Override parent function until authorization is available
    function approvedToLog(address caller) public view returns (bool) {caller; return true;}

    // TBTCSystem
    function fetchOraclePrice() external view returns (uint256) {return _oraclePrice;}
    function fetchRelayCurrentDifficulty() external view returns (uint256) {return _currentDifficulty;}
    function fetchRelayPreviousDifficulty() external view returns (uint256) {return _previousDifficulty;}

    // ERC721
    function balanceOf(address owner) public view returns (uint256 balance) {owner; balance = 0;}
    function ownerOf(uint256 tokenId) public view returns (address owner) {tokenId; owner = _depositOwner;}
    function approve(address to, uint256 tokenId) public {to; tokenId;}
    function getApproved(uint256 tokenId) public view returns (address operator) {tokenId; operator = address(8);}
    function setApprovalForAll(address operator, bool approved) public {operator; approved;}
    function isApprovedForAll(address owner, address operator) public view returns (bool) {owner; operator;}
    function transferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public {from; to; tokenId; data;}
}
