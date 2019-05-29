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
    function approvedToLog(address _caller) public view returns (bool) {_caller; return true;}

    // TBTCSystem
    function fetchOraclePrice() external view returns (uint256) {return _oraclePrice;}
    function fetchRelayCurrentDifficulty() external view returns (uint256) {return _currentDifficulty;}
    function fetchRelayPreviousDifficulty() external view returns (uint256) {return _previousDifficulty;}

    // ERC721
    function balanceOf(address _owner) public view returns (uint256 balance) {_owner; balance = 0;}
    function ownerOf(uint256 _tokenId) public view returns (address owner) {_tokenId; owner = _depositOwner;}
    function approve(address to, uint256 _tokenId) public {to; _tokenId;}
    function getApproved(uint256 _tokenId) public view returns (address operator) {_tokenId; operator = address(8);}
    function setApprovalForAll(address _operator, bool _approved) public {_operator; _approved;}
    function isApprovedForAll(address _owner, address _operator) public view returns (bool) {_owner; _operator;}
    function transferFrom(address _from, address _to, uint256 _tokenId) public {_from; _to; _tokenId;}
    function safeTransferFrom(address _from, address _to, uint256 _tokenId) public {_from; _to; _tokenId;}
    function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes memory _data) public {_from; _to; _tokenId; _data;}
}
