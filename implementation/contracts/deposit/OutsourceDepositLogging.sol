pragma solidity 0.4.25;

import {DepositLog} from "../DepositLog.sol";
import {DepositUtils} from "./DepositUtils.sol";

library OutsourceDepositLogging {


    /// @notice               Fires a Created event
    /// @dev                  We append the sender, which is the deposit contract that called
    /// @param  _keepAddress  The address of the associated keep request
    function logCreated(DepositUtils.Deposit storage _d, address _keepAddress) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logCreated(_keepAddress);
    }

    /// @notice                 Fires a RedemptionRequested event
    /// @dev                    This is the only event without an explicit timestamp
    /// @param  _requester      The ethereum address of the requester
    /// @param  _digest         The calculated sighash digest
    /// @param  _utxoSize       The size of the utxo in sat
    /// @param  _requesterPKH   The requester's 20-byte bitcoin pkh
    /// @param  _requestedFee   The requester or bump-system specified fee
    /// @param  _outpoint       The 36 byte outpoint
    /// @return                 True if successful, else revert
    function logRedemptionRequested(
        DepositUtils.Deposit storage _d,
        address _requester,
        bytes32 _digest,
        uint256 _utxoSize,
        bytes20 _requesterPKH,
        uint256 _requestedFee,
        bytes _outpoint
    ) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logRedemptionRequested(
            _requester,
            _digest,
            _utxoSize,
            _requesterPKH,
            _requestedFee,
            _outpoint);
    }

    /// @notice         Fires a GotRedemptionSignature event
    /// @dev            We append the sender, which is the deposit contract that called
    /// @param  _digest signed digest
    /// @param  _r      signature r value
    /// @param  _s      signature s value
    /// @return         True if successful, else revert
    function logGotRedemptionSignature(
        DepositUtils.Deposit storage _d,
        bytes32 _digest,
        bytes32 _r,
        bytes32 _s
    ) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logGotRedemptionSignature(
            _digest,
            _r,
            _s
        );
    }

    /// @notice     Fires a RegisteredPubkey event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logRegisteredPubkey(
        DepositUtils.Deposit storage _d,
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logRegisteredPubkey(
            _signingGroupPubkeyX,
            _signingGroupPubkeyY);
    }

    /// @notice     Fires a SetupFailed event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logSetupFailed(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logSetupFailed();
    }

    /// @notice     Fires a FraudDuringSetup event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logFraudDuringSetup(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logFraudDuringSetup();
    }

    /// @notice     Fires a Funded event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logFunded(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logFunded();
    }

    /// @notice     Fires a CourtesyCalled event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logCourtesyCalled(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logCourtesyCalled();
    }

    /// @notice             Fires a StartedLiquidation event
    /// @dev                We append the sender, which is the deposit contract that called
    /// @param _wasFraud    True if liquidating for fraud
    function logStartedLiquidation(DepositUtils.Deposit storage _d, bool _wasFraud) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logStartedLiquidation(_wasFraud);
    }

    /// @notice     Fires a Redeemed event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logRedeemed(DepositUtils.Deposit storage _d, bytes32 _txid) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logRedeemed(_txid);
    }

    /// @notice     Fires a Liquidated event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logLiquidated(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logLiquidated();
    }

    /// @notice     Fires a ExitedCourtesyCall event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logExitedCourtesyCall(DepositUtils.Deposit storage _d) external {
        DepositLog _logger = DepositLog(_d.TBTCSystem);
        _logger.logExitedCourtesyCall();
    }
}
