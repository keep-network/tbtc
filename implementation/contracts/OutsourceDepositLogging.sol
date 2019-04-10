pragma solidity 0.4.25;

import {DepositLog} from './DepositLog.sol';
import {TBTCConstants} from './TBTCConstants.sol';

contract OutsourceDepositLogging {

    function logCreated(uint256 _keepID) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logCreated(_keepID);
    }

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
    function logSetupFailed() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logSetupFailed();
    }
    function logFraudDuringSetup() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logFraudDuringSetup();
    }
    function logFunded() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logFunded();
    }
    function logCourtesyCalled() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logCourtesyCalled();
    }
    function logStartedLiquidation(bool _wasFraud) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logStartedLiquidation(_wasFraud);
    }
    function logRedeemed(bytes32 _txid) internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logRedeemed(_txid);
    }
    function logLiquidated() internal {
        DepositLog _logger = DepositLog(TBTCConstants.getSystemContractAddress());
        _logger.logLiquidated();
    }
}
