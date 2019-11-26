pragma solidity ^0.5.10;

contract VendingMachine {
    address depositOwnerToken;

    constructor(
        address _depositOwnerToken
    ) public {
        depositOwnerToken = _depositOwnerToken;
    }

    /// @notice Qualifies a deposit for minting TBTC.
    function qualifyDeposit(
        uint256 _depositId,
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
    function isQualified(uint256 _depositId) public returns (bool) {
        // TODO
        // This is stubbed out for prototyping, separate to the actual qualification logic.
        // However we might remove it later.
        return true;
    }

    /// @notice Pay back the deposit's TBTC and receive the Deposit Owner Token.
    function tbtcToDot(uint256 _depositId) public {
        require(isQualified(_depositId), "Deposit must be qualified");
    }

    /// @notice Trade in the Deposit Owner Token and mint TBTC.
    function dotToTbtc(uint256 _dotId) public {
        require(isQualified(_dotId), "Deposit must be qualified");
        // TODO
    }
}