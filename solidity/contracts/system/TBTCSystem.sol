/* solium-disable function-order */
pragma solidity 0.5.17;

import {IBondedECDSAKeepVendor} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepVendor.sol";
import {IBondedECDSAKeepFactory} from "@keep-network/keep-ecdsa/contracts/api/IBondedECDSAKeepFactory.sol";

import {VendingMachine} from "./VendingMachine.sol";
import {DepositFactory} from "../proxy/DepositFactory.sol";

import {IRelay} from "@summa-tx/relay-sol/contracts/Relay.sol";
import "@summa-tx/relay-sol/contracts/TestnetRelay.sol";
import "../external/IMedianizer.sol";

import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {ISatWeiPriceFeed} from "../interfaces/ISatWeiPriceFeed.sol";
import {DepositLog} from "../DepositLog.sol";

import {TBTCDepositToken} from "./TBTCDepositToken.sol";
import "./TBTCToken.sol";
import "./FeeRebateToken.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

/// @title  TBTC System.
/// @notice This contract acts as a central point for access control,
///         value governance, and price feed.
/// @dev    Governable values should only affect new deposit creation.
contract TBTCSystem is Ownable, ITBTCSystem, DepositLog {

    using SafeMath for uint256;

    event EthBtcPriceFeedAdditionStarted(address _priceFeed, uint256 _timestamp);
    event LotSizesUpdateStarted(uint64[] _lotSizes, uint256 _timestamp);
    event SignerFeeDivisorUpdateStarted(uint16 _signerFeeDivisor, uint256 _timestamp);
    event CollateralizationThresholdsUpdateStarted(
        uint16 _initialCollateralizedPercent,
        uint16 _undercollateralizedThresholdPercent,
        uint16 _severelyUndercollateralizedThresholdPercent,
        uint256 _timestamp
    );

    event EthBtcPriceFeedAdded(address _priceFeed);
    event LotSizesUpdated(uint64[] _lotSizes);
    event AllowNewDepositsUpdated(bool _allowNewDeposits);
    event SignerFeeDivisorUpdated(uint16 _signerFeeDivisor);
    event CollateralizationThresholdsUpdated(
        uint16 _initialCollateralizedPercent,
        uint16 _undercollateralizedThresholdPercent,
        uint16 _severelyUndercollateralizedThresholdPercent
    );


    bool _initialized = false;
    uint256 pausedTimestamp;
    uint256 constant pausedDuration = 10 days;

    ISatWeiPriceFeed public  priceFeed;
    IBondedECDSAKeepVendor public keepVendor;
    IRelay public relay;

    // Parameters governed by the TBTCSystem owner
    bool private allowNewDeposits = false;
    uint16 private signerFeeDivisor = 200; // 1/200 == 50bps == 0.5% == 0.005
    uint16 private initialCollateralizedPercent = 150; // percent
    uint16 private undercollateralizedThresholdPercent = 125;  // percent
    uint16 private severelyUndercollateralizedThresholdPercent = 110; // percent
    uint64[] lotSizesSatoshis = [10**5, 10**6, 10**7, 2 * 10**7, 5 * 10**7, 10**8]; // [0.001, 0.01, 0.1, 0.2, 0.5, 1.0] BTC

    uint256 constant governanceTimeDelay = 6 hours;

    uint256 private signerFeeDivisorChangeInitiated;
    uint256 private lotSizesChangeInitiated;
    uint256 private collateralizationThresholdsChangeInitiated;

    uint16 private newSignerFeeDivisor;
    uint64[] newLotSizesSatoshis;
    uint16 private newInitialCollateralizedPercent;
    uint16 private newUndercollateralizedThresholdPercent;
    uint16 private newSeverelyUndercollateralizedThresholdPercent;

    // price feed
    uint256 priceFeedGovernanceTimeDelay = 90 days;
    uint256 appendEthBtcFeedTimer;
    IMedianizer nextEthBtcFeed;

    constructor(address _priceFeed, address _relay) public {
        priceFeed = ISatWeiPriceFeed(_priceFeed);
        relay = IRelay(_relay);
    }

    /// @notice        Initialize contracts
    /// @dev           Only the Deposit factory should call this, and only once.
    /// @param _keepVendor        ECDSA keep vendor.
    /// @param _depositFactory    Deposit Factory. More info in `DepositFactory`.
    /// @param _masterDepositAddress  Master Deposit address. More info in `Deposit`.
    /// @param _tbtcToken         TBTCToken. More info in `TBTCToken`.
    /// @param _tbtcDepositToken  TBTCDepositToken (TDT). More info in `TBTCDepositToken`.
    /// @param _feeRebateToken    FeeRebateToken (FRT). More info in `FeeRebateToken`.
    /// @param _vendingMachine    Vending Machine. More info in `VendingMachine`.
    /// @param _keepThreshold     Signing group honesty threshold.
    /// @param _keepSize          Signing group size.
    function initialize(
        IBondedECDSAKeepVendor _keepVendor,
        DepositFactory _depositFactory,
        address payable _masterDepositAddress,
        TBTCToken _tbtcToken,
        TBTCDepositToken _tbtcDepositToken,
        FeeRebateToken _feeRebateToken,
        VendingMachine _vendingMachine,
        uint16 _keepThreshold,
        uint16 _keepSize
    ) external onlyOwner {
        require(!_initialized, "already initialized");
        keepVendor = _keepVendor;
        _vendingMachine.setExternalAddresses(
            _tbtcToken,
            _tbtcDepositToken,
            _feeRebateToken
        );
        _depositFactory.setExternalDependencies(
            _masterDepositAddress,
            this,
            _tbtcToken,
            _tbtcDepositToken,
            _feeRebateToken,
            address(_vendingMachine),
            _keepThreshold,
            _keepSize
        );
        setTbtcDepositToken(_tbtcDepositToken);
        _initialized = true;
        allowNewDeposits = true;
    }

    /// @notice gets whether new deposits are allowed.
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

    /// @notice Gets the system signer fee divisor.
    /// @return The signer fee divisor.
    function getSignerFeeDivisor() external view returns (uint16) { return signerFeeDivisor; }

    /// @notice Set the system signer fee divisor.
    /// @dev    This can be finalized by calling `finalizeSignerFeeDivisorUpdate`
    ///         Anytime after `governanceTimeDelay` has elapsed.
    /// @param _signerFeeDivisor The signer fee divisor.
    function beginSignerFeeDivisorUpdate(uint16 _signerFeeDivisor)
        external onlyOwner
    {
        require(
            _signerFeeDivisor > 9,
            "Signer fee divisor must be greater than 9, for a signer fee that is <= 10%."
        );
        require(
            _signerFeeDivisor < 2000,
            "Signer fee divisor must be less than 2000, for a signer fee that is > 0.05%."
        );

        newSignerFeeDivisor = _signerFeeDivisor;
        signerFeeDivisorChangeInitiated = block.timestamp;
        emit SignerFeeDivisorUpdateStarted(_signerFeeDivisor, block.timestamp);
    }

    /// @notice Set the allowed deposit lot sizes.
    /// @dev    Lot size array should always contain 10**8 satoshis (1 BTC) and
    ///         cannot contain values less than 50000 satoshis (0.0005 BTC) or
    ///         greater than 10**9 satoshis (10 BTC).
    ///         This can be finalized by calling `finalizeLotSizesUpdate`
    ///         anytime after `governanceTimeDelay` has elapsed.
    /// @param _lotSizes Array of allowed lot sizes.
    function beginLotSizesUpdate(uint64[] calldata _lotSizes)
        external onlyOwner
    {
        bool hasSingleBitcoin = false;
        for (uint i = 0; i < _lotSizes.length; i++) {
            if (_lotSizes[i] == 10**8) {
                hasSingleBitcoin = true;
            } else if (_lotSizes[i] < 50 * 10**3) {
                // Failed the minimum requirement, break on out.
                revert("Lot sizes less than 0.0005 BTC are not allowed");
            } else if (_lotSizes[i] > 10 * 10**8) {
                // Failed the maximum requirement, break on out.
                revert("Lot sizes greater than 10 BTC are not allowed");
            }
        }

        require(hasSingleBitcoin, "Lot size array must always contain 1 BTC");

        lotSizesSatoshis = _lotSizes;
        emit LotSizesUpdateStarted(_lotSizes, block.timestamp);
        newLotSizesSatoshis = _lotSizes;
        lotSizesChangeInitiated = block.timestamp;
    }

    /// @notice Set the system collateralization levels
    /// @dev    This can be finalized by calling `finalizeCollateralizationThresholdsUpdate`
    ///         Anytime after `governanceTimeDelay` has elapsed.
    /// @param _initialCollateralizedPercent default signing bond percent for new deposits
    /// @param _undercollateralizedThresholdPercent first undercollateralization trigger
    /// @param _severelyUndercollateralizedThresholdPercent second undercollateralization trigger
    function beginCollateralizationThresholdsUpdate(
        uint16 _initialCollateralizedPercent,
        uint16 _undercollateralizedThresholdPercent,
        uint16 _severelyUndercollateralizedThresholdPercent
    ) external onlyOwner {
        require(
            _initialCollateralizedPercent <= 300,
            "Initial collateralized percent must be <= 300%"
        );
        require(
            _initialCollateralizedPercent > 100,
            "Initial collateralized percent must be >= 100%"
        );
        require(
            _initialCollateralizedPercent > _undercollateralizedThresholdPercent,
            "Undercollateralized threshold must be < initial collateralized percent"
        );
        require(
            _undercollateralizedThresholdPercent > _severelyUndercollateralizedThresholdPercent,
            "Severe undercollateralized threshold must be < undercollateralized threshold"
        );

        newInitialCollateralizedPercent = _initialCollateralizedPercent;
        newUndercollateralizedThresholdPercent = _undercollateralizedThresholdPercent;
        newSeverelyUndercollateralizedThresholdPercent = _severelyUndercollateralizedThresholdPercent;
        collateralizationThresholdsChangeInitiated = block.timestamp;
        emit CollateralizationThresholdsUpdateStarted(
            _initialCollateralizedPercent,
            _undercollateralizedThresholdPercent,
            _severelyUndercollateralizedThresholdPercent,
            block.timestamp
        );
    }

    modifier onlyAfterDelay(uint256 _changeInitializedTimestamp) {
        require(_changeInitializedTimestamp > 0, "Change not initiated");
        require(
            block.timestamp.sub(_changeInitializedTimestamp) >=
                governanceTimeDelay,
            "Timer not elapsed"
        );
        _;
    }

    /// @notice Finish setting the system signer fee divisor.
    /// @dev `beginSignerFeeDivisorUpdate` must be called first, once `governanceTimeDelay`
    ///       has passed, this function can be called to set the signer fee divisor to the
    ///       value set in `beginSignerFeeDivisorUpdate`
    function finalizeSignerFeeDivisorUpdate()
        external
        onlyOwner
        onlyAfterDelay(signerFeeDivisorChangeInitiated)
    {
        signerFeeDivisor = newSignerFeeDivisor;
        emit SignerFeeDivisorUpdated(newSignerFeeDivisor);
        newSignerFeeDivisor = 0;
        signerFeeDivisorChangeInitiated = 0;
    }
    /// @notice Finish setting the accepted system lot sizes.
    /// @dev `beginLotSizesUpdate` must be called first, once `governanceTimeDelay`
    ///       has passed, this function can be called to set the lot sizes to the
    ///       value set in `beginLotSizesUpdate`
    function finalizeLotSizesUpdate()
        external
        onlyOwner
        onlyAfterDelay(lotSizesChangeInitiated) {

        lotSizesSatoshis = newLotSizesSatoshis;
        emit LotSizesUpdated(newLotSizesSatoshis);
        lotSizesChangeInitiated = 0;
        newLotSizesSatoshis.length = 0;
    }

    /// @notice Gets the allowed lot sizes
    /// @return Uint64 array of allowed lot sizes
    function getAllowedLotSizes() external view returns (uint64[] memory){
        return lotSizesSatoshis;
    }

    /// @notice Check if a lot size is allowed.
    /// @param _lotSizeSatoshis Lot size to check.
    /// @return True if lot size is allowed, false otherwise.
    function isAllowedLotSize(uint64 _lotSizeSatoshis) external view returns (bool){
        for( uint i = 0; i < lotSizesSatoshis.length; i++){
            if (lotSizesSatoshis[i] == _lotSizeSatoshis){
                return true;
            }
        }
        return false;
    }

    /// @notice Finish setting the system collateralization levels
    /// @dev `beginCollateralizationThresholdsUpdate` must be called first, once `governanceTimeDelay`
    ///       has passed, this function can be called to set the collateralization thresholds to the
    ///       value set in `beginCollateralizationThresholdsUpdate`
    function finalizeCollateralizationThresholdsUpdate()
        external
        onlyOwner
        onlyAfterDelay(collateralizationThresholdsChangeInitiated) {

        initialCollateralizedPercent = newInitialCollateralizedPercent;
        undercollateralizedThresholdPercent = newUndercollateralizedThresholdPercent;
        severelyUndercollateralizedThresholdPercent = newSeverelyUndercollateralizedThresholdPercent;

        emit CollateralizationThresholdsUpdated(
            newInitialCollateralizedPercent,
            newUndercollateralizedThresholdPercent,
            newSeverelyUndercollateralizedThresholdPercent
        );

        newInitialCollateralizedPercent = 0;
        newUndercollateralizedThresholdPercent = 0;
        newSeverelyUndercollateralizedThresholdPercent = 0;
        collateralizationThresholdsChangeInitiated = 0;
    }

    /// @notice Get the system undercollateralization level for new deposits
    function getUndercollateralizedThresholdPercent() external view returns (uint16) {
        return undercollateralizedThresholdPercent;
    }

    /// @notice Get the system severe undercollateralization level for new deposits
    function getSeverelyUndercollateralizedThresholdPercent() external view returns (uint16) {
        return severelyUndercollateralizedThresholdPercent;
    }

    /// @notice Get the system initial collateralized level for new deposits.
    function getInitialCollateralizedPercent() external view returns (uint16) {
        return initialCollateralizedPercent;
    }

    /// @notice Get the time remaining until the collateralization thresholds can be updated.
    function getRemainingCollateralizationUpdateTime() external view returns (uint256) {
        return getRemainingChangeTime(collateralizationThresholdsChangeInitiated);
    }

    /// @notice Get the time remaining until the lot sizes can be updated.
    function getRemainingLotSizesUpdateTime() external view returns (uint256) {
        return getRemainingChangeTime(lotSizesChangeInitiated);
    }

    /// @notice Get the time remaining until the signer fee divisor can be updated.
    function geRemainingSignerFeeDivisorUpdateTime() external view returns (uint256) {
        return getRemainingChangeTime(signerFeeDivisorChangeInitiated);
    }

    /// @notice Get the time remaining until the function parameter timer value can be updated.
    function getRemainingChangeTime(uint256 _changeTimestamp) internal view returns (uint256){
        require(_changeTimestamp > 0, "Update not initiated");
        uint256 elapsed = block.timestamp.sub(_changeTimestamp);
        return (elapsed >= governanceTimeDelay)?
        0:
        governanceTimeDelay.sub(elapsed);
    }

    function getGovernanceTimeDelay() public pure returns (uint256) {
        return governanceTimeDelay;
    }

    function getPriceFeedGovernanceTimeDelay() public view returns (uint256) {
        return priceFeedGovernanceTimeDelay;
    }

    // Price Feed

    /// @notice Get the price of one satoshi in wei.
    /// @dev Reverts if the price of one satoshi is 0 wei, or
    ///      if the price of one satoshi is 1 ether.
    /// @return The price of one satoshi in wei.
    function fetchBitcoinPrice() external view returns (uint256) {
        uint256 price = priceFeed.getPrice();
        if (price == 0 || price > 10 ** 18) {
            // This is if a sat is worth 0 wei, or is worth >1 ether. Revert at
            // once.
            revert("System returned a bad price");
        }
        return price;
    }

    /// @notice Initialize the addition of a new ETH/BTC price feed contract to the priecFeed.
    /// @dev `FinalizeAddEthBtcFeed` must be called to finalize.
    function beginAddEthBtcFeed(IMedianizer _ethBtcFeed)
        external
        onlyOwner {
        nextEthBtcFeed = _ethBtcFeed;
        appendEthBtcFeedTimer = block.timestamp + priceFeedGovernanceTimeDelay;
        emit EthBtcPriceFeedAdditionStarted(address(_ethBtcFeed), block.timestamp);
    }

    /// @notice Finish adding a new price feed contract to the priceFeed.
    /// @dev `InitializeAddEthBtcFeed` must be called first, once `appendEthBtcFeedTimer`
    ///       has passed, this function can be called to append a new price feed.
    function finalizeAddEthBtcFeed()
        external
        onlyOwner {
        require(block.timestamp > appendEthBtcFeedTimer, "Timeout not yet elapsed");
        priceFeed.addEthBtcFeed(nextEthBtcFeed);
        emit EthBtcPriceFeedAdded(address(nextEthBtcFeed));
    }

    // Difficulty Oracle
    function fetchRelayCurrentDifficulty() external view returns (uint256) {
        return relay.getCurrentEpochDifficulty();
    }

    function fetchRelayPreviousDifficulty() external view returns (uint256) {
        return relay.getPrevEpochDifficulty();
    }

    /// @notice Gets a fee estimate for creating a new Deposit.
    /// @return Uint256 estimate.
    function getNewDepositFeeEstimate()
        external
        view
        returns (uint256)
    {
        IBondedECDSAKeepFactory _keepFactory = IBondedECDSAKeepFactory(keepVendor.selectFactory());
        return _keepFactory.openKeepFeeEstimate();
    }

    /// @notice Request a new keep opening.
    /// @param _m Minimum number of honest keep members required to sign.
    /// @param _n Number of members in the keep.
    /// @param _maxSecuredLifetime Duration of stake lock in seconds.
    /// @return Address of a new keep.
    function requestNewKeep(
        uint256 _m,
        uint256 _n,
        uint256 _bond,
        uint256 _maxSecuredLifetime
    )
        external
        payable
        returns (address)
    {
        require(tbtcDepositToken.exists(uint256(msg.sender)), "Caller must be a Deposit contract");
        IBondedECDSAKeepFactory _keepFactory = IBondedECDSAKeepFactory(keepVendor.selectFactory());
        return _keepFactory.openKeep.value(msg.value)(_n, _m, msg.sender, _bond, _maxSecuredLifetime);
    }
}
