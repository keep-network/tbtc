/* solium-disable function-order */
pragma solidity ^0.5.10;

import {IKeepRegistry} from "@keep-network/keep-ecdsa/contracts/api/IKeepRegistry.sol";
import {IECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeepVendor.sol";
import {VendingMachine} from "./VendingMachine.sol";
import {DepositFactory} from "../proxy/DepositFactory.sol";

import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {IBTCETHPriceFeed} from "../interfaces/IBTCETHPriceFeed.sol";
import {DepositLog} from "../DepositLog.sol";

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract TBTCSystem is Ownable, ITBTCSystem, DepositLog {

    using SafeMath for uint256;

    event LotSizesUpdated(uint256[] _lotSizes);
    event AllowNewDepositsUpdated(bool _allowNewDeposits);
    event SignerFeeDivisorUpdated(uint256 _signerFeeDivisor);
    event CollateralizationThresholdsUpdated(
        uint256 _undercollateralizedThresholdPercent,
        uint256 _severelyUndercollateralizedThresholdPercent
    );

    bool _initialized = false;
    uint256 pausedTimestamp;
    uint256 pausedDuration = 10 days;

    uint256 currentDifficulty = 1;
    uint256 previousDifficulty = 1;

    address public keepRegistry;
    address public priceFeed;

    // Parameters governed by the TBTCSystem owner
    bool private allowNewDeposits = false;
    uint256 private signerFeeDivisor = 200; // 1/200 == 50bps == 0.5% == 0.005
    uint128 private undercollateralizedThresholdPercent = 140;  // percent
    uint128 private severelyUndercollateralizedThresholdPercent = 120; // percent
    uint256[] lotSizesSatoshis = [10**5, 10**6, 10**7, 20**7, 50**7, 10**8]; // [0.001, 0.01, 0.1, 0.2, 0.5, 1.0] BTC

    constructor(address _priceFeed) public {
        priceFeed = _priceFeed;
    }

    function initialize(
        address _keepRegistry,
        address _depositFactory,
        address _masterDepositAddress,
        address _tbtcToken,
        address _tbtcDepositToken,
        address _feeRebateToken,
        address _vendingMachine,
        uint256 _keepThreshold,
        uint256 _keepSize
    ) external onlyOwner {
        require(!_initialized, "already initialized");

        keepRegistry = _keepRegistry;
        VendingMachine(_vendingMachine).setExternalAddresses(
            _tbtcToken,
            _tbtcDepositToken,
            _feeRebateToken
        );
        DepositFactory(_depositFactory).setExternalDependencies(
            _masterDepositAddress,
            address(this),
            _tbtcToken,
            _tbtcDepositToken,
            _feeRebateToken,
            _vendingMachine,
            _keepThreshold,
            _keepSize
        );
        _initialized = true;
        allowNewDeposits = true;
    }

    /// @notice gets whether new deposits are allowed
    function getAllowNewDeposits() external view returns (bool) { return allowNewDeposits; }

    /// @notice One-time-use emergency function to disallow future deposit creation for 10 days. 
    function emergencyPauseNewDeposits() external onlyOwner returns (bool) { 
        require(pausedTimestamp == 0, "emergencyPauseNewDeposits can only be called once");
        pausedTimestamp = block.timestamp;
        allowNewDeposits = false;
        emit AllowNewDepositsUpdated(false);
    }

    /// @notice Anyone can reactivate deposit creations after the pause duration is over.
    function resumeNewDeposits() public {
        require(allowNewDeposits == false, "New deposits are currently allowed");
        require(pausedTimestamp != 0, "Deposit has not been paused");
        require(block.timestamp.sub(pausedTimestamp) >= pausedDuration, "Deposits are still paused");
        allowNewDeposits = true;
        emit AllowNewDepositsUpdated(true);
    }

    function getRemainingPauseTerm() public view returns (uint256) {
        require(allowNewDeposits == false, "New deposits are currently allowed");
        return (block.timestamp.sub(pausedTimestamp) >= pausedDuration)?
            0:
            pausedDuration.sub(block.timestamp.sub(pausedTimestamp));
    }

    /// @notice Set the system signer fee divisor.
    /// @param _signerFeeDivisor The signer fee divisor.
    function setSignerFeeDivisor(uint256 _signerFeeDivisor)
        external onlyOwner
    {
        require(_signerFeeDivisor > 9, "Signer fee divisor must be greater than 9, for a signer fee that is <= 10%.");
        signerFeeDivisor = _signerFeeDivisor;
        emit SignerFeeDivisorUpdated(_signerFeeDivisor);
    }

    /// @notice Gets the system signer fee divisor.
    /// @return The signer fee divisor.
    function getSignerFeeDivisor() external view returns (uint256) { return signerFeeDivisor; }

    /// @notice Set the allowed deposit lot sizes.
    /// @dev    Lot sizes should be 
    /// @param _lotSizes Array of allowed lot sizes.
    function setLotSizes(uint256[] calldata _lotSizes) external onlyOwner {
        for( uint i = 0; i < _lotSizes.length; i++){
            if (_lotSizes[i] == 10**8){
                lotSizesSatoshis = _lotSizes;
                emit LotSizesUpdated(_lotSizes);
                return;
            }
        }
        revert("Lot size array must always contain 1BTC");
    }

    /// @notice Gets the allowed lot sizes
    /// @return Uint256 array of allowed lot sizes 
    function getAllowedLotSizes() external view returns (uint256[] memory){
        return lotSizesSatoshis;
    }

    /// @notice Check if a lot size is allowed.
    /// @param _lotSize Lot size to check.
    /// @return True if lot size is allowed, false otherwise. 
    function isAllowedLotSize(uint256 _lotSize) external view returns (bool){
        for( uint i = 0; i < lotSizesSatoshis.length; i++){
            if (lotSizesSatoshis[i] == _lotSize){
                return true;
            }
        }
        return false;
    }

    /// @notice Set the system collateralization levels
    /// @param _undercollateralizedThresholdPercent first undercollateralization trigger
    /// @param _severelyUndercollateralizedThresholdPercent second undercollateralization trigger
    function setCollateralizationThresholds(
        uint128 _undercollateralizedThresholdPercent,
        uint128 _severelyUndercollateralizedThresholdPercent
    ) external onlyOwner {
        require(
            _undercollateralizedThresholdPercent > _severelyUndercollateralizedThresholdPercent,
            "Severe undercollateralized threshold must be > undercollateralized threshold"
        );
        undercollateralizedThresholdPercent = _undercollateralizedThresholdPercent;
        severelyUndercollateralizedThresholdPercent = _severelyUndercollateralizedThresholdPercent;
        emit CollateralizationThresholdsUpdated(
            _undercollateralizedThresholdPercent,
            _severelyUndercollateralizedThresholdPercent
        );
    }

    /// @notice Get the system undercollateralization level for new deposits
    function getUndercollateralizedThresholdPercent() external view returns (uint128) {
        return undercollateralizedThresholdPercent;
    }

    /// @notice Get the system severe undercollateralization level for new deposits
    function getSeverelyUndercollateralizedThresholdPercent() external view returns (uint128) {
        return severelyUndercollateralizedThresholdPercent;
    }

    // Price Feed
    function fetchBitcoinPrice() external view returns (uint256) {
        return IBTCETHPriceFeed(priceFeed).getPrice();
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
