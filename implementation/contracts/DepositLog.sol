pragma solidity ^0.5.10;

contract DepositLog {

    /*
    TODO: review events, see what new information should be added
    Logging philosophy:
      Every state transition should fire a log
      That log should have ALL necessary info for off-chain actors
      Everyone should be able to ENTIRELY rely on log messages
    */

    // This event is fired when we init the deposit
    event Created(
        address indexed _depositContractAddress,
        address indexed _keepAddress,
        uint256 _timestamp
    );

    // This log event contains all info needed to rebuild the redemption tx
    // We index on request and signers and digest
    event RedemptionRequested(
        address indexed _depositContractAddress,
        address indexed _requester,
        bytes32 indexed _digest,
        uint256 _utxoSize,
        bytes20 _requesterPKH,
        uint256 _requestedFee,
        bytes _outpoint
    );

    // This log event contains all info needed to build a witnes
    // We index the digest so that we can search events for the other log
    event GotRedemptionSignature(
        address indexed _depositContractAddress,
        bytes32 indexed _digest,
        bytes32 _r,
        bytes32 _s,
        uint256 _timestamp
    );

    // This log is fired when the signing group returns a public key
    event RegisteredPubkey(
        address indexed _depositContractAddress,
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY,
        uint256 _timestamp
    );

    // This event is fired when we enter the SETUP_FAILED state for any reason
    event SetupFailed(address indexed _depositContractAddress, uint256 _timestamp);

    // This event is fired when we detect an ECDSA fraud before seeing a funding proof
    event FraudDuringSetup(address indexed _depositContractAddress, uint256 _timestamp);

    // This event is fired when we enter the ACTIVE state
    event Funded(address indexed _depositContractAddress, uint256 _timestamp);

    // This event is called when we enter the COURTESY_CALL state
    event CourtesyCalled(address indexed _depositContractAddress, uint256 _timestamp);

    // This event is fired when we go from COURTESY_CALL to ACTIVE
    event ExitedCourtesyCall(address indexed _depositContractAddress, uint256 _timestamp);

    // This log event is fired when liquidation
    event StartedLiquidation(
        address indexed _depositContractAddress,
        bool _wasFraud,
        uint256 _timestamp
    );

    // This event is fired when the Redemption SPV proof is validated
    event Redeemed(
        address indexed _depositContractAddress,
        bytes32 indexed _txid,
        uint256 _timestamp
    );

    // This event is fired when Liquidation is completed
    event Liquidated(
        address indexed _depositContractAddress,
        uint256 _timestamp
    );

    //
    // AUTH
    //

    /// @notice             Checks if an address is an allowed logger
    /// @dev                Calls the system to check if the caller is a Deposit
    ///                     We don't require this, so deposits are not bricked if the system borks
    /// @param  _caller     The address of the calling contract
    /// @return             True if approved, otherwise false
    /* solium-disable-next-line no-empty-blocks */
    function approvedToLog(address _caller) public pure returns (bool) {
        /* TODO: auth via system */
        _caller;
        return true;
    }

    //
    // Logging
    //

    /// @notice               Fires a Created event
    /// @dev                  We append the sender, which is the deposit contract that called
    /// @param  _keepAddress  The address of the associated keep
    /// @return               True if successful, else revert
    function logCreated(address _keepAddress) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit Created(
            msg.sender,
            _keepAddress,
            block.timestamp);
        return true;
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
        bytes memory _outpoint
    ) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit RedemptionRequested(
            msg.sender,
            _requester,
            _digest,
            _utxoSize,
            _requesterPKH,
            _requestedFee,
            _outpoint);
        return true;
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
    ) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit GotRedemptionSignature(
            msg.sender,
            _digest,
            _r,
            _s,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a RegisteredPubkey event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logRegisteredPubkey(
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY
    ) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit RegisteredPubkey(
            msg.sender,
            _signingGroupPubkeyX,
            _signingGroupPubkeyY,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a SetupFailed event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logSetupFailed() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit SetupFailed(
            msg.sender,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a FraudDuringSetup event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logFraudDuringSetup() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit FraudDuringSetup(
            msg.sender,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a Funded event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logFunded() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit Funded(
            msg.sender,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a CourtesyCalled event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logCourtesyCalled() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit CourtesyCalled(
            msg.sender,
            block.timestamp);
        return true;
    }

    /// @notice             Fires a StartedLiquidation event
    /// @dev                We append the sender, which is the deposit contract that called
    /// @param _wasFraud    True if liquidating for fraud
    /// @return             True if successful, else revert
    function logStartedLiquidation(bool _wasFraud) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit StartedLiquidation(
            msg.sender,
            _wasFraud,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a Redeemed event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logRedeemed(bytes32 _txid) public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit Redeemed(
            msg.sender,
            _txid,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a Liquidated event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logLiquidated() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit Liquidated(
            msg.sender,
            block.timestamp);
        return true;
    }

    /// @notice     Fires a ExitedCourtesyCall event
    /// @dev        We append the sender, which is the deposit contract that called
    ///             returns false if not approved, to prevent accidentally halting Deposit
    /// @return     True if successful, else false
    function logExitedCourtesyCall() public returns (bool) {
        if (!approvedToLog(msg.sender)) return false;
        emit ExitedCourtesyCall(
            msg.sender,
            block.timestamp);
        return true;
    }
}
