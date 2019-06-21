pragma solidity 0.4.25;

import {ITBTCSystem} from "./ITBTCSystem.sol";
import {IERC721} from "./IERC721.sol";
import {IUniswapFactory} from "../uniswap/IUniswapFactory.sol";
import {TBTC} from "../tokens/TBTC.sol";


contract TBTCSystemStub is ITBTCSystem, IERC721 {

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;
    address depositOwner = address(0);
    IUniswapFactory public uniswapFactory;
    TBTC public tbtc;

    function setup(
        address _uniswapFactory,
        address _tbtc
    ) external {
        uniswapFactory = IUniswapFactory(_uniswapFactory);
        tbtc = TBTC(_tbtc);
        uniswapFactory.createExchange(_tbtc);
    }

    function getTBTCUniswapExchange() external view returns (address) {
        return uniswapFactory.getExchange(address(tbtc));
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
