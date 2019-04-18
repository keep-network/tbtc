pragma solidity 0.4.25;

import {DepositLog} from './DepositLog.sol';
import {TBTCConstants} from './TBTCConstants.sol';

contract OutsourceDepositLogging {


    /// @notice             Fires a Created event
    /// @dev                We append the sender, which is the deposit contract that called
    /// @param  _keepID     The ID of the associated keep request
    function logCreated(uint256 _keepID) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logCreated(_keepID);
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
        address _requester,
        bytes32 _digest,
        uint256 _utxoSize,
        bytes20 _requesterPKH,
        uint256 _requestedFee,
        bytes _outpoint
    ) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
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
        bytes32 _digest,
        bytes32 _r,
        bytes32 _s
    ) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logGotRedemptionSignature(
            _digest,
            _r,
            _s);
    }

    /// @notice     Fires a RegisteredPubkey event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logRegisteredPubkey(
        address _signingGroupAccount,
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logRegisteredPubkey(
            _signingGroupAccount,
            _signingGroupPubkeyX,
            _signingGroupPubkeyY);
    }

    /// @notice     Fires a SetupFailed event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logSetupFailed() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logSetupFailed();
    }

    /// @notice     Fires a FraudDuringSetup event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logFraudDuringSetup() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logFraudDuringSetup();
    }

    /// @notice     Fires a Funded event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logFunded() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logFunded();
    }

    /// @notice     Fires a CourtesyCalled event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logCourtesyCalled() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logCourtesyCalled();
    }

    /// @notice             Fires a StartedLiquidation event
    /// @dev                We append the sender, which is the deposit contract that called
    /// @param _wasFraud    True if liquidating for fraud
    function logStartedLiquidation(bool _wasFraud) internal {aw
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logStartedLiquidation(_wasFraud);
    }

    /// @notice     Fires a Redeemed event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logRedeemed(bytes32 _txid) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logRedeemed(_txid);
    }

    /// @notice     Fires a Liquidated event
    /// @dev        The logger is on a system contract, so all logs from all deposits are from the smae addres
    function logLiquidated() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logLiquidated();
    }
}
