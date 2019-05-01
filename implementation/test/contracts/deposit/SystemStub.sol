pragma solidity 0.4.25;

import {ITBTCSystem} from '../../../contracts/interfaces/ITBTCSystem.sol';
import {IERC721} from '../../../contracts/interfaces/IERC721.sol';
import {DepositLog} from '../../../contracts/DepositLog.sol';

contract SystemStub is ITBTCSystem, IERC721, DepositLog {

    uint256 current = 1;
    uint256 past = 1;
    uint256 oraclePrice = 10 ** 12;

    function setOraclePrice(uint256 _oraclePrice) external {
        oraclePrice = _oraclePrice;
    }

    function setCurrentDiff(uint256 _current) external {
        current = _current;
    }

    function setPreviousDiff(uint256 _past) external {
        past = _past;
    }

    function fetchOraclePrice() external view returns (uint256) {
        return oraclePrice;
    }

    // passthrough requests for the oracle
    function fetchRelayCurrentDifficulty() external view returns (uint256) {
        return current;
    }

    function fetchRelayPreviousDifficulty() external view returns (uint256) {
        return current;
    }

    // 721
    function balanceOf(address owner) public view returns (uint256 balance) {owner; balance = 0;}
    function ownerOf(uint256 tokenId) public view returns (address owner) {tokenId; owner = address(7);}
    function approve(address to, uint256 tokenId) public {to; tokenId;}
    function getApproved(uint256 tokenId) public view returns (address operator) {tokenId; operator = address(8);}
    function setApprovalForAll(address operator, bool _approved) public {operator; _approved;}
    function isApprovedForAll(address owner, address operator) public view returns (bool) {owner; operator;}
    function transferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId) public {from; to; tokenId;}
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory data) public {from; to; tokenId; data;}
}
