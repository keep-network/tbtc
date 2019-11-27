pragma solidity ^0.5.10;

contract VendingMachine {
    address depositOwnerToken;
    address tbtcToken;

    constructor(
        address _tbtcToken,
        address _depositOwnerToken
    ) public {
        tbtcToken = _tbtcToken;
        depositOwnerToken = _depositOwnerToken;
    }

    /// @notice Qualifies a deposit for minting TBTC.
    function qualifyDeposit(
        address _depositAddress,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public {
        // require(!isQualified(_depositId), "Deposit already qualified");
        // TODO
    }

    /// @notice Determines whether a deposit is qualified for minting TBTC.
    function isQualified(address _depositAddress) public returns (bool) {
        // TODO
        // This is stubbed out for prototyping, separate to the actual qualification logic.
        // However we might remove it later.
        return true;
    }

    /// @notice Pay back the deposit's TBTC and receive the Deposit Owner Token.
    function tbtcToDot(uint256 _dotId) public {
        require(isQualified(address(_dotId)), "Deposit must be qualified");
    }

    /// @notice Trade in the Deposit Owner Token and mint TBTC.
    function dotToTbtc(uint256 _dotId) public {
        require(isQualified(address(_dotId)), "Deposit must be qualified");
        // TODO
    }
}