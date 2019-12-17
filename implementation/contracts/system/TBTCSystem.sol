/* solium-disable function-order */
pragma solidity ^0.5.10;

import {IKeepRegistry} from "@keep-network/keep-ecdsa/contracts/api/IKeepRegistry.sol";
import {IECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeepVendor.sol";

import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {DepositLog} from "../DepositLog.sol";

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract TBTCSystem is Ownable, ITBTCSystem, DepositLog {

    bool _initialized = false;

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;
    uint256 oraclePrice = 10 ** 12;

    address public keepRegistry;

    // Governed parameters by the TBTCSystem owner
    bool private allowNewDeposits = true;
    uint256 private signerFeeDivisor = 200; // 1/200 == 50bps == 0.5% == 0.005

    function initialize(
        address _keepRegistry
    ) external onlyOwner {
        require(!_initialized, "already initialized");

        keepRegistry = _keepRegistry;
        _initialized = true;
    }

    /// @notice Enables/disables new deposits from being created.
    /// @param _allowNewDeposits Whether to allow new deposits.
    function setAllowNewDeposits(bool _allowNewDeposits)
        external onlyOwner
    {
        allowNewDeposits = _allowNewDeposits;
    }

    /// @notice Gets whether new deposits are allowed.
    function getAllowNewDeposits() public view returns (bool) { return allowNewDeposits; }

    /// @notice Set the system signer fee divisor.
    /// @param _signerFeeDivisor The signer fee divisor.
    function setSignerFeeDivisor(uint256 _signerFeeDivisor)
        external onlyOwner
    {
        require(_signerFeeDivisor > 1, "Signer fee must be lower than 100%");
        signerFeeDivisor = _signerFeeDivisor;
    }

    /// @notice Gets the system signer fee divisor.
    /// @return The signer fee divisor.
    function getSignerFeeDivisor() public view returns (uint256) { return signerFeeDivisor; }

    // Price Oracle
    function fetchOraclePrice() external view returns (uint256) {
        return oraclePrice;
    }

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

    /// @notice Request a new keep opening.
    /// @param _m Minimum number of honest keep members required to sign.
    /// @param _n Number of members in the keep.
    /// @return Address of a new keep.
    function requestNewKeep(uint256 _m, uint256 _n)
        external
        payable
        returns (address _keepAddress)
    {
        address keepVendorAddress = IKeepRegistry(keepRegistry)
            .getVendor("ECDSAKeep");

        _keepAddress = IECDSAKeepVendor(keepVendorAddress)
            .openKeep(_n,_m, msg.sender);
    }
}
