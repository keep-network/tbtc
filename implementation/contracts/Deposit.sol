pragma solidity 0.4.25;

import {SafeMath} from "./SafeMath.sol";
import {BytesLib} from "./BytesLib.sol";
import {BTCUtils} from "./BTCUtils.sol";
import {ValidateSPV} from "./ValidateSPV.sol";
import {CheckBitcoinSigs} from './SigCheck.sol';
import {TBTCConstants} from './TBTCConstants.sol';
import {IBurnableERC20} from './IBurnableERC20.sol';
import {IERC721} from './IERC721.sol';

contract Deposit {

    using BytesLib for bytes;
    using BTCUtils for bytes;
    using SafeMath for uint256;
    using ValidateSPV for bytes;
    using ValidateSPV for bytes32;
    using CheckBitcoinSigs for bytes;

    enum DepositStates {
        // DOES NOT EXIST YET
        START,

        // FUNDING FLOW
        AWAITING_SIGNER_SETUP,
        AWAITING_BTC_FUNDING_PROOF,

        // FAILED SETUP
        FRAUD_AWAITING_BTC_FUNDING_PROOF,
        FAILED_SETUP,

        // ACTIVE
        ACTIVE,

        // REDEMPTION FLOW
        AWAITING_WITHDRAWAL_SIGNATURE,
        AWAITING_WITHDRAWAL_PROOF,
        REDEEMED,

        // SIGNER LIQUIDATION FLOW
        FRAUD_LIQUIDATION_IN_PROGRESS,
        LIQUIDATION_IN_PROGRESS,
        LIQUIDATED
    }

    /*
    TODO: should logging be part of the system contract or this one?
    TODO: More events

    Logging philosophy:
      Every state transition should fire a unique log
      That log should have ALL necessary info for off-chain actors
      Everyone should be able to ENTIRELY rely on log messages
    */

    // This log event contains all info needed to rebuild the redemption tx
    // We index on request and signers and digest
    event RedemptionRequested(
        address indexed _requester,
        bytes20 indexed _signerPKH,
        bytes32 indexed _digest,
        uint256 _utxoSize,
        bytes20 _requesterPKH,
        uint256 _requestedFee,
        bytes _outpoint);

    // This log event contains all info needed to build a witnes
    // We index the digest so that we can search events for the other log
    event GotRedemptionSignature(
        bytes32 indexed _digest,
        bytes32 _x,
        bytes32 _y,
        bytes32 _r,
        bytes32 _s);

    // This log is fired when the signing group returns a public key
    event RegisteredPubkey(
        address _signingGroupAccount,
        bytes32 _signingGroupPubkeyX,
        bytes32 _signingGroupPubkeyY);

    // This log event is fired when liquidation
    event StartedLiquidation(
        uint256 _timestamp);

    // This event is fired when the Redemption SPV proof is validated
    event Redeemed(
        bytes32 indexed_txid,
        uint256 _timestamp);

    // SET DURING CONSTRUCTION
    DepositStates currentState;

    // SET ON FRAUD
    uint256 liquidationInitiated;  // Timestamp of when liquidation starts

    // INITIALLY WRITTEN BY FUNDING FLOW
    uint256 signingGroupRequestedAt;  // timestamp of signing group request
    bytes32 signingGroupPubkeyX;  // The X coordinate of the signing group's pubkey
    bytes32 signingGroupPubkeyY;  // The Y coordinate of the signing group's pubkey
    uint256 fundingProofTimerStart;  // start of the funding proof period. reused for funding fraud proof period
    bytes8 depositSizeBytes;  // LE uint. the size of the deposit UTXO in satoshis
    bytes utxoOutpoint;  // the 36-byte outpoint of the custodied UTXO

    // INITIALLY WRITTEN BY REDEMPTION FLOW
    bytes20 requesterPKH;  // The 20-byte requeser PKH
    uint256 initialRedemptionFee;  // the initial fee as requested
    uint256 withdrawalRequestTime;  // the most recent withdrawal request timestamp
    bytes32 lastRequestedDigest;  // the digest most recently requested for signing
    mapping(bytes32 => uint256) public wasRequested; // Signatures that have been requested and the timestamp

    // We separate the constructor from createNewDeposit to make proxy factories easier
    constructor () public {}

    function () public payable { require(false); }  // don't do that

    // THIS IS THE INIT FUNCTION
    /// @notice                         The system can spin up a new deposit
    /// @dev                            This should be called by an approved contract, not a developer
    /// @param      _approvedDigests    Allow the system to set a list of non-slashable messages
    /// @return                         True if succesful, otherwise revert
    function createNewDeposit(
        bytes32[] _approvedDigests
    ) public returns (bool) {
        require(currentState == DepositStates.START, 'Deposit setup already requested');
        require(isApprovedDepositCreator(msg.sender), 'Calling account not allowed to create deposits');
        currentState = DepositStates.AWAITING_SIGNER_SETUP;

        // This allows the system to set digests that the signers may sign
        for (uint i = 0; i < _approvedDigests.length; i++) {
            wasRequested[_approvedDigests[i]] = block.timestamp;
        }
        signingGroupRequestedAt = block.timestamp;
        return true;
    }



    //
    // CHECKING STATES
    //

    /// @notice     Check if the contract is currently in the funding flow
    /// @dev        This checks on the funding flow happy path, not the fraud path
    /// @return     True if contract is currently in the funding flow else False
    function inFunding() public view returns (bool) {
        return (currentState == DepositStates.AWAITING_SIGNER_SETUP
             || currentState == DepositStates.AWAITING_BTC_FUNDING_PROOF);
    }

    /// @notice     Check if the contract is currently in the funding faud flow
    /// @dev        This checks for the flow, not the SETUP_FAILED termination state
    /// @return     True if contract is currently in the funding fraud flow else False
    function inFundingFailure() public view returns (bool) {
        return (currentState == DepositStates.FRAUD_AWAITING_BTC_FUNDING_PROOF);
    }

    /// @notice     Check if the contract is currently in the signer liquidation flow
    /// @dev        This could be caused by fraud, or by an unfilled margin call
    /// @return     True if contract is currently in the liquidaton flow else False
    function inSignerLiquidation() public view returns (bool) {
        return (currentState == DepositStates.LIQUIDATION_IN_PROGRESS
             || currentState == DepositStates.FRAUD_LIQUIDATION_IN_PROGRESS);
    }

    /// @notice     Check if the contract is currently in the redepmtion flow
    /// @dev        This checks on the redemption flow, not the REDEEMED termination state
    /// @return     True if contract is currently in the redemption flow else False
    function inRedemption() public view returns (bool) {
        return (currentState == DepositStates.AWAITING_WITHDRAWAL_SIGNATURE
             || currentState == DepositStates.AWAITING_WITHDRAWAL_PROOF);
    }

    /// @notice     Check if the contract has halted
    /// @dev        This checks on any halt state, regardless of triggering circumstances
    /// @return     True if contract has halted permanently
    function inEndState() public view returns (bool) {
        return (currentState == DepositStates.LIQUIDATED
             || currentState == DepositStates.REDEEMED
             || currentState == DepositStates.FAILED_SETUP);
    }

    /// @notice     Get the integer representing the current state
    /// @dev        We implement this because contracts don't handle foreign enums well
    /// @return     The 0-indexed state from the DepositStates enum
    function getCurrentState() public view returns (uint256) {
        return uint256(currentState);
    }

    //
    // CHEKCING PERMISSIONS
    //

    /// @notice         Check if the caller is an approved deposit creator
    /// @dev            A deposit must be deployed and initated by the system, not a user
    /// @param _caller  The address of the caller to compare to approved creators
    /// @return         True if the caller is approved, else False
    function isApprovedDepositCreator(address _caller) public pure returns (bool) {
        return isTBTCSystemContract(_caller);
    }

    /// @notice         Check if the caller is the tBTC system contract
    /// @dev            Stored as a constant in the config library
    /// @param _caller  The address of the caller to compare to the tbtc system constant
    /// @return         True if the caller is approved, else False
    function isTBTCSystemContract(address _caller) public pure returns (bool) {
        return _caller == TBTCConstants.getSystemContractAddress();
    }

    /// @notice         Check if the caller is the beneficiary
    /// @dev            Simple address comparison
    /// @param _caller  The caller to compare to the beneficiary
    /// @return         True if the caller is the beneficiary, else False
    function isBeneficiary(address _caller) public view returns (bool) {
        return _caller == depositBeneficiary();
    }

    //
    // DERIVED PROPERTIES
    //

    /// @notice         Determines the fees due to the signers for work performeds
    /// @dev            Signers are paid based on the TBTC issued
    /// @return         Accumulated fees in smallest TBTC unit (satoshi)
    function signerFee() public view returns (uint256) {
        return depositSize().div(TBTCConstants.getSignerFeeDivisor());
    }

    // TODO: Docs
    function beneficiaryReward() public view returns (uint256) {
        return depositSize().div(TBTCConstants.getBeneficiaryRewardDivisor());
    }

    /// @notice         Determines the outstanding TBTC
    /// @dev            This is the amount of TBTC needed to repay to redeem the Deposit
    /// @return         Outstanding debt in smallest TBTC unit
    function outstandingTBTC() public view returns (uint256) {
        return depositSize().add(signerFee()).add(beneficiaryReward());
    }

    /// @notice         Returns the packed public key (64 bytes) for the signing group
    /// @dev            We store it as 2 bytes32, (2 slots) then repack it on demand
    /// @return         64 byte public key
    function signerPubkey () public view returns (bytes) {
        return abi.encodePacked(signingGroupPubkeyX, signingGroupPubkeyY);
    }

    /// @notice         Returns the Bitcoin pubkeyhash (hash160) for the signing group
    /// @dev            This is used in bitcoin output scripts for the signers
    /// @return         20-bytes public key hash
    function signerPKH() public view returns (bytes20) {
        bytes memory _pubkey = abi.encodePacked(hex'04', signerPubkey());
        bytes memory _digest = _pubkey.hash160();
        return bytes20(_digest.toAddress(0));  // dirty solidity hack
    }

    /// @notice         Returns the Ethereum account address for the signing group
    /// @dev            This is used in Ethereum signature checking
    /// @return         20-bytes public key hash substring (Ethereum address)
    function signerAccount() public view returns (address) {
        return signerPubkey().accountFromPubkey();
    }

    /// @notice         Returns the size of the deposit UTXO in satoshi
    /// @dev            We store the deposit as bytes8 to make signature checking easier
    /// @return         UTXO value in satoshi
    function depositSize() public view returns (uint256) {
        return bytes8LEToUint(depositSizeBytes);
    }

    // TODO: document
    function funderBondRefundAmount() public pure returns (uint256) {
        TBTCConstants.getFunderBondRefundAmount();
    }

    //
    // EXTERNAL CALLS
    //

    /* TODO: call keep and push ether/tbtc out to keep members */
    function pushFundsToKeepGroup(uint256 _ethValue) internal returns (bool) { /* TODO */ }

    /* TODO: call keep and return the value of the bond seized */
    function seizeSignerBonds() internal returns (uint256) { /* TODO */ }

    /* TODO: calls out to the keep contract to get the result*/
    function getKeepPubkeyResult() public view returns (bytes) { /* TODO */ }

    /* TODO: should we return the current diff here so that stateless proofs can be checked? */
    /// @notice         Gets the current block difficulty
    /// @dev            Calls the light relay and gets the current block difficulty
    /// @return         The proof difficulty requirement
    function currentBlockDifficulty() public view returns (uint256) { /* TODO */ }

    /// @notice         Looks up the deposit beneficiary by calling the tBTC system
    /// @dev            We cast the address to a uint256 to match the 721 standard
    /// @return         The current deposit beneficiary
    function depositBeneficiary() public view returns (address) {
        IERC721 _systemContract = IERC721(TBTCConstants.getSystemContractAddress());
        return _systemContract.ownerOf(uint256(address(this)));
    }

    /// @notice     Tries to liquidate the position on-chain using the signer bond
    /// @dev        Calls out to other contracts, watch for re-entrance
    /// @return     True if Liquidated, False otherwise
    function attmeptToLiquidateOnChain() internal returns (bool) { /* TODO */ }

    //
    // REUSABLE STATE TRANSITIONS
    //

    /// @notice         Get the difficulty required for the contract to accept an SPV proof
    /// @dev            Calls the light relay and gets the current block difficulty
    /// @param _fraud   True if fraud, false otherwise
    /// @return         The proof difficulty requirement
    function startSignerLiquidation(bool _fraud) internal {
        emit StartedLiquidation(
            block.timestamp);

        // Reclaim used state for gas savings
        fundingTeardown();
        redemptionTeardown();
        seizeSignerBonds();
        attmeptToLiquidateOnChain();

        bool _liquidated = attmeptToLiquidateOnChain();
        if (_liquidated && !_fraud) {
            returnRemainingSignerBond();
            currentState = DepositStates.LIQUIDATED;
        }
        if (_liquidated && _fraud) {
            address(0).transfer(address(this).balance);  /* TODO: is this what we want? */
            currentState = DepositStates.LIQUIDATED;
        }
        if (!_liquidated && !_fraud) {
            liquidationInitiated = block.timestamp;  // Store the timestamp for auction
            currentState = DepositStates.LIQUIDATION_IN_PROGRESS;
        }
        if (!_liquidated && _fraud) {
            liquidationInitiated = block.timestamp;  // Store the timestamp for auction
            currentState = DepositStates.FRAUD_LIQUIDATION_IN_PROGRESS;  // Essentially storing a bool in the state
        }
    }

    /* TODO: When we exit the funding flow, we should delete any set state vars */
    function fundingTeardown() internal { /* TODO */ }

    /* TODO: When we exit the funding fraud flow, we should delete any set state vars */
    function fundingFraudTeardown() internal { /* TODO */ }

    /* TODO: When we exit the redemption flow, we should delete any set state vars */
    function redemptionTeardown() internal { /* TODO */ }

    function revokeFunderBond() internal {
        pushFundsToKeepGroup(funderBondRefundAmount());
    }

    function returnFunderBond() internal {
        depositBeneficiary().transfer(funderBondRefundAmount());
    }

    function partiallySlashForFraudInFunding() internal {
        uint256 _seized = seizeSignerBonds();
        uint256 _slash = _seized.div(TBTCConstants.getFundingFraudPartialSlashDivisor());
        pushFundsToKeepGroup(_seized - _slash);
        depositBeneficiary().transfer(_slash);  // TODO: is this what we want?
    }

    // TODO: Docs
    function distributeSignerBondsToFunder() internal {
        uint256 _seized = seizeSignerBonds();
        depositBeneficiary().transfer(_seized);
    }

    // TODO: Docs
    function returnRemainingSignerBond() internal {
        pushFundsToKeepGroup(address(this).balance);
    }

    //
    // UTILS
    //

    /// @notice         Convert a LE bytes8 to a uint256
    /// @dev            Do this by converting to bytes, then reversing endianness, then converting to int
    /// @return         The uint256 represented in LE by the bytes8
    function bytes8LEToUint(bytes8 _b) public pure returns (uint256) {
        return abi.encodePacked(_b).reverseEndianness().bytesToUint();
    }

    /// @notice                 Syntactically check an SPV proof for a bitcoin tx
    /// @dev                    Stateless SPV Proof verification documented elsewhere
    /// @param  _bitcoinTx      The bitcoin tx that is purportedly included in the header chain
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 The 32 byte transaction id (little-endian, not block-explorer)
    function checkProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) internal view returns (bytes32) {
        bytes memory _nIns;
        bytes memory _ins;
        bytes memory _nOuts;
        bytes memory _outs;
        bytes memory _locktime;
        bytes32 _txid;
        (_nIns, _ins, _nOuts, _outs, _locktime, _txid) = _bitcoinTx.parseTransaction();
        require(_txid != bytes32(0), 'Failed tx parsing');
        require(
            _txid.prove(
                _bitcoinHeaders.extractMerkleRootLE().toBytes32(),
                _merkleProof,
                _index),
            'Tx merkle proof is not valid for provided header and tx');

        /*
        TODO: Does the currentBlockDifficulty return the total diff?
              Should it return a single-header diff so we can check all headers
              have that diff?
         */
        require(_bitcoinHeaders.validateHeaderChain() > currentBlockDifficulty(),
                'Insufficient accumulated difficulty in header chain');

        return _txid;
    }

    /// @notice                 Parses a bitcoin tx to find an output paying the signing group PKH
    /// @dev                    Reverts if no funding output found
    /// @param  _bitcoinTx      The bitcoin tx that should contain the funding output
    /// @return                 The 8-byte LE encoded value, and the index of the output
    function findAndParseFundingOutput(
        bytes _bitcoinTx
    ) internal view returns (bytes8, uint8) {
        /* TODO */
        bytes8 _valueBytes;
        bytes memory _output;

        // Find the output paying the signer PKH
        // This will fail if there are more than 256 outputs
        for (uint8 i = 0; i <  _bitcoinTx.extractNumOutputs(); i++) {
            _output = _bitcoinTx.extractOutputAtIndex(i);
            if (keccak256(_output.extractHash()) == keccak256(abi.encodePacked(signerPKH()))) {
                return (_valueBytes, i);
            }
        }
        // If we don't return from inside the loop, we failed.
        require(false, 'Did not find output with correct PKH');
    }

    /// @notice                 Validates the funding tx and parses information from it
    /// @dev                    Stateless SPV Proof & Bitcoin tx format documented elsewhere
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contain the funding output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 The 8-byte LE UTXO size in satoshi, the 36byte outpoint
    function validateAndParseFundingSPVProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) internal view returns (bytes8 _valueBytes, bytes _outpoint) {
        uint8 _outputIndex;
        bytes32 _txid = checkProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);
        (_valueBytes, _outputIndex) = findAndParseFundingOutput(_bitcoinTx);

        // Don't validate deposits under the lot size
        require(bytes8LEToUint(_valueBytes) >= TBTCConstants.getLotSize(), 'Deposit too small');

        // The outpoint is the LE TXID plus the index of the output as a 4-byte LE int
        // _outputIndex is a uint8, so we know it is only 1 byte
        // Therefore, pad with 3 more bytes
        _outpoint = abi.encodePacked(_txid, _outputIndex, hex'000000');
    }

    //
    // STATE TRANSITIONS
    //

    /// @notice                     Anyone can request redemption
    /// @dev                        The redeemer specifies details about the Bitcoin redemption tx
    /// @param  _outputValueBytes   The 8-byte LE output size
    /// @param  _requesterPKH       The 20-byte Bitcoin pubkeyhash to which to send funds
    /// @return                     True if succesful, otherwise revert
    function requestRedemption(
        bytes8 _outputValueBytes,
        bytes20 _requesterPKH
    ) public returns (bool) {
        require(currentState == DepositStates.ACTIVE, 'Redemption only available from Active state');

        currentState = DepositStates.AWAITING_WITHDRAWAL_SIGNATURE;

        // Burn the redeemer's TBTC plus enough extra to cover outstanding debt
        // Requires user to approve first
        // TODO: implement such that it calls the system to burn TBTC
        IBurnableERC20 TBTCContract = IBurnableERC20(TBTCConstants.getTokenContractAddress());
        require(TBTCContract.balanceOf(msg.sender) > outstandingTBTC(), 'Not enough TBTC to cover outstanding debt');
        TBTCContract.burnFrom(msg.sender, outstandingTBTC());

        // Convert the 8-byte LE ints to uint256
        uint256 _outputValue = abi.encodePacked(_outputValueBytes).reverseEndianness().bytesToUint();
        uint256 _requestedFee = depositSize().sub(_outputValue);
        require(_requestedFee >= TBTCConstants.getMinimumRedemptionFee());

        // Calculate the sighash
        bytes32 _sighash = CheckBitcoinSigs.oneInputOneOutputSighash(
            utxoOutpoint,
            signerPKH(),
            depositSizeBytes,
            _outputValueBytes,
            _requesterPKH);

        emit RedemptionRequested(
            msg.sender,
            signerPKH(),
            _sighash,
            depositSize(),
            _requesterPKH,
            _requestedFee,
            utxoOutpoint);

        // write all request details
        requesterPKH = _requesterPKH;
        initialRedemptionFee = _requestedFee;
        withdrawalRequestTime = block.timestamp;
        lastRequestedDigest = _sighash;
        wasRequested[_sighash] = block.timestamp;

        return true;
    }

    /// @notice     Anyone may provide a withdrawal signature if it was requested
    /// @dev        The signers will be penalized if this (or provideRedemptionProof) is not called
    /// @param  _v  Signature recovery value
    /// @param  _r  Signature R value
    /// @param  _s  Signature S value
    /// @return     True if succesful, False if prevented by timeout, otherwise revert
    function provideRedemptionSignature(
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) public returns (bool) {
        require(currentState == DepositStates.AWAITING_WITHDRAWAL_SIGNATURE, 'Not currently awaiting a signature');
        // A signature has been provided, now we wait for fee bump or redemption
        currentState = DepositStates.AWAITING_WITHDRAWAL_PROOF;

        // If we're outside of the signature window, signers have aborted, initate punishment
        /* TODO: discussion: should we instead require an explicit call to notifySignatureTimeout? */
        if (block.timestamp > withdrawalRequestTime + TBTCConstants.getSignatureTimeout()) {
            startSignerLiquidation(false); // Not fraud, just failure
            return false;  // We return instead of reverting so that the above transition takes place
        }

        // The signature must be valid on the pubkey
        require(signerPubkey().checkSig(
            lastRequestedDigest,
            _v,
            _r,
            _s));

        emit GotRedemptionSignature(
            lastRequestedDigest,
            signingGroupPubkeyX,
            signingGroupPubkeyY,
            _r,
            _s);

        return true;
    }

    /// @notice                             Anyone may notify the contract that a fee bump is needed
    /// @dev                                This sends us back to AWAITING_WITHDRAWAL_SIGNATURE
    /// @param  _previousOutputValueBytes   The previous output's value
    /// @param  _newOutputValueBytes        The new output's value
    /// @return                             True if succesful, False if prevented by timeout, otherwise revert
    function increaseRedemptionFee(
        bytes8 _previousOutputValueBytes,
        bytes8 _newOutputValueBytes
    ) public returns (bool) {
        require(currentState == DepositStates.AWAITING_WITHDRAWAL_PROOF);
        require(block.timestamp >= withdrawalRequestTime + TBTCConstants.getIncreaseFeeTimer(), 'Fee increase not yet permitted');

        // If we should have gotten a redemption proof by now, something fishy is going on
        if (block.timestamp > withdrawalRequestTime + TBTCConstants.getRedepmtionProofTimeout()) {
            startSignerLiquidation(false);
            return false;  // We return instead of reverting so that the above transition takes place
        }

        // Calculate the previous one so we can check that it really is the previous one
        bytes32 _previousSighash = CheckBitcoinSigs.oneInputOneOutputSighash(
            utxoOutpoint,
            signerPKH(),
            depositSizeBytes,
            _previousOutputValueBytes,
            requesterPKH);
        require(wasRequested[_previousSighash] == withdrawalRequestTime, 'Provided previous value does not yield previous sighash');

        // Check that we're incrementing the fee by exactly the requester's initial fee
        uint256 _previousOutputValue = abi.encodePacked(_previousOutputValueBytes).reverseEndianness().bytesToUint();
        uint256 _newOutputValue = abi.encodePacked(_newOutputValueBytes).reverseEndianness().bytesToUint();
        require(_previousOutputValue.sub(_newOutputValue) == initialRedemptionFee, 'Not an allowed fee step');

        // Calculate the next sighash
        bytes32 _sighash = CheckBitcoinSigs.oneInputOneOutputSighash(
            utxoOutpoint,
            signerPKH(),
            depositSizeBytes,
            _newOutputValueBytes,
            requesterPKH);

        emit RedemptionRequested(
            msg.sender,
            signerPKH(),
            _sighash,
            depositSize(),
            requesterPKH,
            depositSize() - _newOutputValue,
            utxoOutpoint);

        // Ratchet the signature and redemption proof timeouts
        withdrawalRequestTime = block.timestamp;
        wasRequested[_sighash] = block.timestamp;
        lastRequestedDigest = _sighash;

        // Go back to waiting for a signature
        currentState = DepositStates.AWAITING_WITHDRAWAL_SIGNATURE;

        return true;
    }

    /// @notice                 Anyone may provide a withdrawal proof to prove redemption
    /// @dev                    The signers will be penalized if this is not called
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contain the redemption output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 True if succesful, False if prevented by timeout, otherwise revert
    function provideRedemptionProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) public returns (bool) {
        bytes memory _nIns;
        bytes memory _ins;
        bytes memory _nOuts;
        bytes memory _outs;
        bytes memory _locktime;
        bytes32 _txid;
        require(inRedemption(), 'Redemption proof only allowed from redemption flow');

        // We don't use checkproof here because we need access to the parse info
        (_nIns, _ins, _nOuts, _outs, _locktime, _txid) = _bitcoinTx.parseTransaction();
        require(_txid != bytes32(0), 'Failed tx parsing');
        require(
            _txid.prove(
                _bitcoinHeaders.extractMerkleRootLE().toBytes32(),
                _merkleProof,
                _index),
            'Tx merkle proof is not valid for provided header');

        uint256 _currentDiff = currentBlockDifficulty();
        require(_bitcoinHeaders.validateHeaderChain() > _currentDiff * 6, // TODO
                'Insufficient accumulated difficulty in header chain');

        require(keccak256(_locktime) == keccak256(hex'00000000'), 'Wrong locktime set');
        require(keccak256(_nIns) == keccak256(hex'01'), 'Too many ins');
        require(keccak256(_nOuts) == keccak256(hex'01'), 'Too many outs');
        require(keccak256(_ins.extractOutpoint()) == keccak256(utxoOutpoint),
                'Tx spends the wrong UTXO');
        require(keccak256(_outs.extractHash()) == keccak256(abi.encodePacked(requesterPKH)),
                'Tx sends value to wrong pubkeyhash');
        /* TODO: refactor redemption flow to improve this */
        require((depositSize() - uint256(_outs.extractValue())) <= initialRedemptionFee * 5, 'Fee unexpectedly very high');

        emit Redeemed(
            _txid,
            block.timestamp);

        // We're done yey!
        currentState = DepositStates.REDEEMED;
        redemptionTeardown();
        return true;
    }

    /// @notice     Anyone may notify the contract that the signers have failed to produce a signature
    /// @dev        This is considered fraud, and is punished
    /// @return     True if succesful, otherwise revert
    function notifySignatureTimeout() public returns (bool) {
        require(currentState == DepositStates.AWAITING_WITHDRAWAL_SIGNATURE);
        require(block.timestamp > withdrawalRequestTime + TBTCConstants.getSignatureTimeout());
        startSignerLiquidation(false);  // not fraud, just failure
        return true;
    }

    /// @notice     Anyone may notify the contract that the signers have failed to produce a redemption proof
    /// @dev        This is considered fraud, and is punished
    /// @return     True if succesful, otherwise revert
    function notifyRedemptionProofTimeout() public returns (bool) {
        require(currentState == DepositStates.AWAITING_WITHDRAWAL_PROOF);
        require(block.timestamp > withdrawalRequestTime + TBTCConstants.getRedepmtionProofTimeout());
        startSignerLiquidation(false);  // not fraud, just failure
        return true;
    }

    //
    // FUNDING FLOW
    //

    /* TODO: is this a safe assumption? */
    /// @notice     Anyone may notify the contract that signing group setup has timed out
    /// @dev        We assume that the keep system punishes the signers
    /// @return     True if succesful, otherwise revert
    function notifySignerSetupFailure() public returns (bool) {
        require(currentState == DepositStates.AWAITING_SIGNER_SETUP, 'Not awaiting setup');
        require(block.timestamp > signingGroupRequestedAt + TBTCConstants.getSigningGroupFormationTimeout(),
                'Signing group formation timeout not yet elapsed');
        currentState = DepositStates.FAILED_SETUP;

        returnFunderBond();
        fundingTeardown();

        return true;
    }

    /* TODO: Will the Keep contract call us? or do we need to call it? */
    /// @notice             The Keep contract notifies the Deposit of the signing group's key
    /// @dev                We store the pubkey as 2 bytestrings, X and Y.
    /// @return             True if succesful, otherwise revert
    function retrieveSignerPubkey() public returns (bool) {
        /* TODO INCOMPLETE*/
        bytes memory _keepResult = getKeepPubkeyResult();

        signingGroupPubkeyX = _keepResult.slice(0, 32).toBytes32();
        signingGroupPubkeyY = _keepResult.slice(32, 32).toBytes32();
        fundingProofTimerStart = block.timestamp;

        emit RegisteredPubkey(
            signerAccount(),
            signingGroupPubkeyX,
            signingGroupPubkeyY
        );

        return true;
    }

    /// @notice     Anyone may notify the contract that the funder has failed to send BTC
    /// @dev        This is considered a funder fault, and we revoke their bond
    /// @return     True if succesful, otherwise revert
    function notifyFundingTimeout() public returns (bool) {
        require(currentState == DepositStates.AWAITING_BTC_FUNDING_PROOF, 'Funding timeout has not started');
        require(block.timestamp > fundingProofTimerStart + TBTCConstants.getFundingTimeout(),
                'Funding timeout has not elapsed.');
        currentState = DepositStates.FAILED_SETUP;

        revokeFunderBond();
        fundingTeardown();

        return true;
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud during funding
    /// @dev                    ECDSA is NOT SECURE unless you verify the digest
    /// @param  _v              Signature recovery value
    /// @param  _r              Signature R value
    /// @param  _s              Signature S value
    /// @param _signedDigest    The digest signed by the signature vrs tuple
    /// @param _preimage        The sha256 preimage of the digest
    /// @return                 True if succesful, False if prevented by timeout, otherwise revert
    function provideFundingECDSAFraudProof(
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) public returns (bool) {
        require(currentState == DepositStates.AWAITING_BTC_FUNDING_PROOF,
                'Signer fraud during funding flow only available while awaiting funding');

        // If the funding timeout has elapsed, punish the funder
        if (block.timestamp > fundingProofTimerStart + TBTCConstants.getFundingTimeout()) {
            currentState = DepositStates.FAILED_SETUP;
            revokeFunderBond();
            fundingTeardown();
            return false;
        }

        currentState = DepositStates.FRAUD_AWAITING_BTC_FUNDING_PROOF;

        /* TODO: outsource this to the Keep contract? */
        bool _valid = signerPubkey().checkSig(_signedDigest, _v, _r, _s);
        _valid = _valid && _preimage.isSha256Preimage(_signedDigest);
        require(_valid, 'Signature is not valid');  // Invalid signatures error

        /* NB: This is reuse of the variable */
        fundingProofTimerStart = block.timestamp;

        returnFunderBond();
        fundingTeardown();  /* TODO: is this right? */

        return true;
    }

    /// @notice     Anyone may notify the contract no funding proof was submitted during funding fraud
    /// @dev        This is not a funder fault. The signers have faulted, so the funder shouldn't fund
    /// @return     True if succesful, otherwise revert
    function notifyFraudFundingTimeout() public returns (bool) {
        require(currentState == DepositStates.FRAUD_AWAITING_BTC_FUNDING_PROOF,
                'Not currently awaiting fraud-related funding proof');
        require(block.timestamp > fundingProofTimerStart + TBTCConstants.getFraudFundingTimeout(),
                'Fraud funding proof timeout has not elapsed');
        currentState = DepositStates.FAILED_SETUP;

        partiallySlashForFraudInFunding();
        fundingFraudTeardown();

        return true;
    }

    /// @notice                 Anyone may notify the deposit of a funding proof during funding fraud
    /// @dev                    We reward the funder the entire bond if this occurs
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contains the funding output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 True if succesful, False if prevented by timeout, otherwise revert
    function provideFraudBTCFundingProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) public returns (bool) {
        bytes8 _valueBytes;
        bytes memory _outpoint;
        require(currentState == DepositStates.FRAUD_AWAITING_BTC_FUNDING_PROOF);
        currentState = DepositStates.FAILED_SETUP;

        (_valueBytes, _outpoint) = validateAndParseFundingSPVProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);

        // If the proof is accepted, update to failed, and distribute signer bonds
        distributeSignerBondsToFunder();
        fundingFraudTeardown();

        return true;
    }

    /// @notice                 Anyone may notify the deposit of a funding proof to activate the deposit
    /// @dev                    This is the happy-path of the funding flow. It means that we have suecceeded
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contains the funding output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 True if succesful, False if prevented by timeout, otherwise revert
    function provideBTCFundingProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) public returns (bool) {
        bytes8 _valueBytes;
        bytes memory _outpoint;

        require(currentState == DepositStates.AWAITING_BTC_FUNDING_PROOF);

        // Design decision:
        // We COULD revoke the funder bond here if the funding proof timeout has elapsed
        // HOWEVER, that would only create a situation where the funder loses eerything
        // It would be a large punishment for a small crime (being slightly late)
        // So if the funder manages to call this before anyone notifies of timeout
        // We let them have a freebie
        currentState = DepositStates.ACTIVE;

        (_valueBytes, _outpoint) = validateAndParseFundingSPVProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);

        // Write down the UTXO info and set to active. Congratulations :)
        depositSizeBytes = _valueBytes;
        utxoOutpoint = _outpoint;

        // The funder has performed their duty
        returnFunderBond();
        fundingTeardown();

        return true;
    }

    //
    // FRAUD
    //

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud
    /// @dev                    ECDSA is NOT SECURE unless you verify the digest
    /// @param  _v              Signature recovery value
    /// @param  _r              Signature R value
    /// @param  _s              Signature S value
    /// @param _signedDigest    The digest signed by the signature vrs tuple
    /// @param _preimage        The sha256 preimage of the digest
    /// @return                 True if succesful, otherwise revert
    function provideECDSAFraudProof(
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes _preimage
    ) public returns (bool) {
        require(!inFunding() && !inFundingFailure(),
                'Use provideFundingECDSAFraudProof instead');
        require(!inSignerLiquidation(),
                'Signer liquidation already in progress');

        // Check that the signature is valid
        /* TODO: Outsource this to the keep contract? */
        bool _valid = signerPubkey().checkSig(_signedDigest, _v, _r, _s);
        _valid = _valid && _preimage.isSha256Preimage(_signedDigest);
        require(_valid, 'Signature is not valid');

        if (inRedemption()) {
            require(wasRequested[_signedDigest] == 0, 'Digest is approved for signing');
        }

        startSignerLiquidation(true);
        return true;
    }

    /* TODO: Can we cut this and rely entirely on ECDSA fraud? */
    /// @notice                 Anyone may notify the deposit of fraud via an SPV proof
    /// @dev                    We strong prefer ECDSA fraud proofs
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contains the funding output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 True if succesful, False if prevented by timeout, otherwise revert
    function provideSPVFraudProof(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders
    ) public returns (bool) {
        bytes memory _input;
        bytes memory _output;
        bool _inputConsumed;
        uint8 i;

        require(!inFunding() && !inFundingFailure(),
                'SPV Fraud proofs not valid before Active state.');
        require(!inSignerLiquidation(),
                'Signer liquidation already in progress');

        checkProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);
        for (i = 0; i < _bitcoinTx.extractNumInputs(); i++) {
            _input = _bitcoinTx.extractInputAtIndex(i);
            if (keccak256(_input.extractOutpoint()) == keccak256(utxoOutpoint)) {
                _inputConsumed = true;
            }
        }
        require(_inputConsumed, 'No input spending custodied UTXO found');

        uint256 _permittedFeeBumps = 5;  /* TODO: can we refactor withdrawal flow to improve this? */
        uint256 _requiredOutputSize = depositSize() - (initialRedemptionFee * (1 + _permittedFeeBumps));
        for (i = 0; i < _bitcoinTx.extractNumOutputs(); i++) {
            _output = _bitcoinTx.extractOutputAtIndex(i);
            if (_output.extractValue() >= _requiredOutputSize
                && keccak256(_output.extractHash()) == keccak256(abi.encodePacked(requesterPKH))) {
                require(false, 'Found an output paying the redeemer as requested');
            }
        }

        startSignerLiquidation(true);
        return true;
    }
}
