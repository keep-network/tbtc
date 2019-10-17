pragma solidity ^0.5.10;

import {TBTCSystem} from '../../../contracts/system/TBTCSystem.sol';

contract TBTCSystemStub is TBTCSystem {
    address keepAddress = address(7);

    constructor(address _depositFactory)
        // Set expected factory address to 0-address.
        // Address is irelevant as test use forceMint function to bypass ACL
        TBTCSystem(_depositFactory)
    public {
        // solium-disable-previous-line no-empty-blocks
    }

    // Bypasses the ACL on TBTCSystem.initialize
    // and allows repeat initialization.
    function reinitialize(address _tbtcUniswapExchange) external {
        tbtcUniswapExchange = _tbtcUniswapExchange;
    }

    function setOraclePrice(uint256 _oraclePrice) external {
        oraclePrice = _oraclePrice;
    }

    function setCurrentDiff(uint256 _currentDifficulty) external {
        currentDifficulty = _currentDifficulty;
    }

    function setPreviousDiff(uint256 _previousDifficulty) external {
        previousDifficulty = _previousDifficulty;
    }

    // override parent
    function approvedToLog(address _caller) public view returns (bool) {
        _caller; return true;
    }

    function requestNewKeep(uint256 _m, uint256 _n) external payable returns (address _keepAddress) {
        _m; _n;
        return keepAddress;
    }

    /// @notice          Function to mint a new token.
    /// @dev             We can't call 721 mint function from deposit Test becuase of ACL.
    ///                  This function bypasses ACL and can be called in Deposit tests
    ///                  Reverts if the given token ID already exists.
    /// @param _to       The address that will own the minted token
    /// @param _tokenId  uint256 ID of the token to be minted
    function forceMint(address _to, uint256 _tokenId) public {
        _mint(_to, _tokenId);
    }
}
