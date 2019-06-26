pragma solidity 0.4.25;

import {ITBTCSystem} from "./ITBTCSystem.sol";
import {IERC721} from "./IERC721.sol";
import {TBTCToken} from "./TBTCToken.sol";

contract TBTCSystemStub is ITBTCSystem, IERC721 {

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;
    address depositOwner = address(0);
    address _TBTCToken;

    //Update auth with ACL. Can't be set in constructor
    function setExternalAddresses(address tokenAddress) external {
        _TBTCToken = tokenAddress;
    }

    // Price Oracle
    function fetchOraclePrice() external view returns (uint256) {return oraclePrice;}

    // Difficulty Oracle
    // TODO: This is a workaround. It will be replaced by tbtc-difficulty-oracle.
    function fetchRelayCurrentDifficulty() external view returns (uint256) {
        return currentDifficulty;
    }

    function fetchRelayPreviousDifficulty() external view returns (uint256) {
        return previousDifficulty;
    }

    function submitCurrentDifficulty(uint256 _currentDifficulty) public {
        if (currentDifficulty != _currentDifficulty) {
            previousDifficulty = currentDifficulty;
            currentDifficulty = _currentDifficulty;
        }
    }

    //modify access with ACL implementation
    function systemMint(address _account, uint256 _amount) public {
        TBTCToken(_TBTCToken).systemMint(_account, _amount);
    }

    //modify access with ACL implementation
    function systemBurnFrom(address _account, uint256 _amount) public {
        TBTCToken(_TBTCToken).systemBurnFrom(_account, _amount);
    }

    //modify access with ACL implementation
    function systemTransferFrom(address _from, address _to, uint256 _value) public {
        TBTCToken(_TBTCToken).systemTransferFrom(_from, _to, _value);
    }

    // ERC721
    function balanceOf(address _owner) public view returns (uint256 balance) {_owner; balance = 0;}
    function ownerOf(uint256 _tokenId) public view returns (address owner) {_tokenId; owner = depositOwner;}
    function approve(address _to, uint256 _tokenId) public {_to; _tokenId;}
    function getApproved(uint256 _tokenId) public view returns (address operator) {_tokenId; operator = address(8);}
    function setApprovalForAll(address _operator, bool _approved) public {_operator; _approved;}
    function isApprovedForAll(address _owner, address _operator) public view returns (bool) {_owner; _operator;}
    function transferFrom(address _from, address _to, uint256 _tokenId) public {_from; _to; _tokenId;}
    function safeTransferFrom(address _from, address _to, uint256 _tokenId) public {_from; _to; _tokenId;}
    function safeTransferFrom(address _from, address _to, uint256 _tokenId, bytes memory _data) public {_from; _to; _tokenId; _data;}
}
