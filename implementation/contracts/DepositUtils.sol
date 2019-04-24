pragma solidity 0.4.25;

import {TBTCConstants} from './TBTCConstants.sol';
import {SafeMath} from "../bitcoin-spv/contracts/SafeMath.sol";
import {ITBTCSystem} from './interfaces/ITBTCSystem.sol';
import {BytesLib} from "../bitcoin-spv/contracts/BytesLib.sol";
import {BTCUtils} from "../bitcoin-spv/contracts/BTCUtils.sol";
import {ValidateSPV} from "../bitcoin-spv/contracts/ValidateSPV.sol";


library DepositUtils {

   /// @notice             Determines the prefix to the compressed public key
    /// @dev                The prefix encodes the parity of the Y coordinate
    /// @param  _pubkeyY    The Y coordinate of the public key
    /// @return             The 1-byte prefix for the compressed key
    function determineCompressionPrefix(bytes32 _pubkeyY) public pure returns (bytes) {
        if(uint256(_pubkeyY) & 1 == 1) {
            return hex'03';  // Odd Y
        } else {
            return hex'02';  // Even Y
        }
    }

    /// @notice             Compresses a public key
    /// @dev                Converts the 64-byte key to a 33-byte key, bitcoin-style
    /// @param  _pubkeyX    The X coordinate of the public key
    /// @param  _pubkeyY    The Y coordinate of the public key
    /// @return             The 33-byte compressed pubkey
    function compressPubkey(bytes32 _pubkeyX, bytes32 _pubkeyY) public pure returns (bytes) {
        return abi.encodePacked(determineCompressionPrefix(_pubkeyY), _pubkeyX);
    }
    
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

    /// @notice         Determines the fees due to the signers for work performeds
    /// @dev            Signers are paid based on the TBTC issued
    /// @return         Accumulated fees in smallest TBTC unit (tsat)
    function signerFee() public pure returns (uint256) {
        return SafeMath.div(TBTCConstants.getLotSize(),TBTCConstants.getSignerFeeDivisor());
    }

    /// @notice         Convert a LE bytes8 to a uint256
    /// @dev            Do this by converting to bytes, then reversing endianness, then converting to int
    /// @return         The uint256 represented in LE by the bytes8
    function bytes8LEToUint(bytes8 _b) public pure returns (uint256) {
        return BTCUtils.bytesToUint(BTCUtils.reverseEndianness(abi.encodePacked(_b)));
    }

    /// @notice     calculates the beneficiary reward based on the deposit size
    /// @dev        the amount of extra ether to pay the beneficiary at closing time
    /// @return     the amount of ether in wei to pay the beneficiary
    function beneficiaryReward() public pure returns (uint256) {
        return SafeMath.div(TBTCConstants.getLotSize(),TBTCConstants.getBeneficiaryRewardDivisor());
    }

    /// @notice         Gets the current block difficulty
    /// @dev            Calls the light relay and gets the current block difficulty
    /// @return         The difficulty
    function currentBlockDifficulty() public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(TBTCConstants.getSystemContractAddress());
        return _sys.fetchRelayCurrentDifficulty();
    }

    /// @notice         Gets the previous block difficulty
    /// @dev            Calls the light relay and gets the previous block difficulty
    /// @return         The difficulty
    function previousBlockDifficulty() public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(TBTCConstants.getSystemContractAddress());
        return _sys.fetchRelayPreviousDifficulty();
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
        (_nIns, _ins, _nOuts, _outs, _locktime, _txid) = ValidateSPV.parseTransaction(_bitcoinTx);
        require(_txid != bytes32(0), 'Failed tx parsing');
        require(
            ValidateSPV.prove(
                _txid,
                BytesLib.toBytes32(BTCUtils.extractMerkleRootLE(_bitcoinHeaders)),
                _merkleProof,
                _index),
            'Tx merkle proof is not valid for provided header and tx');

        require(ValidateSPV.validateHeaderChain(_bitcoinHeaders) > SafeMath.mul(currentBlockDifficulty(),6),
                'Insufficient accumulated difficulty in header chain');

        return _txid;
    }
    function SpvFraudProofHelper( 
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders,
        bytes20 _requesterPKH,
        bytes _utxoOutpoint,
        uint256 _initialRedemptionFee,
        uint256 _utxoSize ) public view {
        
        bytes memory _input;
        bytes memory _output;
        bool _inputConsumed;
        uint8 i;

        checkProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);
        for (i = 0; i < BTCUtils.extractNumInputs(_bitcoinTx); i++) {
            _input = BTCUtils.extractInputAtIndex(_bitcoinTx, i);
            if (keccak256(BTCUtils.extractOutpoint(_input)) == keccak256(_utxoOutpoint)) {
                _inputConsumed = true;
            }
        }
        require(_inputConsumed, 'No input spending custodied UTXO found');

        uint256 _permittedFeeBumps = 5;  /* TODO: can we refactor withdrawal flow to improve this? */
        uint256 _requiredOutputSize = SafeMath.sub(_utxoSize, (_initialRedemptionFee * (1 + _permittedFeeBumps)));
        for (i = 0; i < BTCUtils.extractNumOutputs(_bitcoinTx); i++) {
            _output = BTCUtils.extractOutputAtIndex(_bitcoinTx, i);
            if (BTCUtils.extractValue(_output) >= _requiredOutputSize
                && keccak256(BTCUtils.extractHash(_output)) == keccak256(abi.encodePacked(_requesterPKH))) {
                revert('Found an output paying the redeemer as requested');
            }
        }
    }
      /// @notice                 Parses a bitcoin tx to find an output paying the signing group PKH
    /// @dev                    Reverts if no funding output found
    /// @param  _bitcoinTx      The bitcoin tx that should contain the funding output
    /// @return                 The 8-byte LE encoded value, and the index of the output
    function findAndParseFundingOutput(
        bytes _bitcoinTx, bytes20 _signerPKH
    ) internal pure returns (bytes8, uint8) {
        bytes8 _valueBytes;
        bytes memory _output;

        // Find the output paying the signer PKH
        // This will fail if there are more than 256 outputs
        for (uint8 i = 0; i <  BTCUtils.extractNumOutputs(_bitcoinTx); i++) {
            _output = BTCUtils.extractOutputAtIndex(_bitcoinTx, i);
            if (keccak256(BTCUtils.extractHash(_output)) == keccak256(abi.encodePacked(_signerPKH))) {
                return (_valueBytes, i);
            }
        }
        // If we don't return from inside the loop, we failed.
        revert('Did not find output with correct PKH');
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
        bytes _bitcoinHeaders,
        bytes20 _signerPKH
    ) internal view returns (bytes8 _valueBytes, bytes _outpoint) {
        uint8 _outputIndex;
        bytes32 _txid = checkProof(_bitcoinTx, _merkleProof, _index, _bitcoinHeaders);
        (_valueBytes, _outputIndex) = findAndParseFundingOutput(_bitcoinTx, _signerPKH);

        // Don't validate deposits under the lot size
        require(bytes8LEToUint(_valueBytes) >= TBTCConstants.getLotSize(), 'Deposit too small');

        // The outpoint is the LE TXID plus the index of the output as a 4-byte LE int
        // _outputIndex is a uint8, so we know it is only 1 byte
        // Therefore, pad with 3 more bytes
        _outpoint = abi.encodePacked(_txid, _outputIndex, hex'000000');
    }
    
     /// @notice                     Evaluates the header difficulties in a proof
    /// @dev                        Uses the light oracle to source recent difficulty
    /// @param  _bitcoinHeaders     The header chain to evaluate
    /// @return                     True if acceptable, otherwise revert
    function evaluateProofDifficulty(bytes _bitcoinHeaders) public view returns (bool) {
        uint256 _reqDiff;
        uint256 _current = currentBlockDifficulty();
        uint256 _previous = previousBlockDifficulty();
        uint256 _firstHeaderDiff = BTCUtils.calculateDifficulty(BTCUtils.extractTarget(_bitcoinHeaders));

        if (_firstHeaderDiff == _current) {
            _reqDiff = _current;
        } else if (_firstHeaderDiff == _previous) {
            _reqDiff = _previous;
        } else {
            revert('not at current or previous difficulty');
        }

        /* TODO: make this better than 6 */
        require(ValidateSPV.validateHeaderChain(_bitcoinHeaders) > SafeMath.mul(_reqDiff, 6),
                'Insufficient accumulated difficulty in header chain');
        return true;
    }

    /// @notice                 helper function for redemption proof checks
    /// @dev                    The signers will be penalized if this is not called
    /// @param  _bitcoinTx      The bitcoin tx that purportedly contain the redemption output
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block
    /// @param  _index          The index of the tx in the Bitcoin block (1-indexed)
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers
    /// @return                 txid if successful, revert otherwise
    function RedemptionProofChecks(
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _index,
        bytes _bitcoinHeaders,
        bytes _utxoOutpoint,
        bytes20 _requesterPKH
    ) public view returns (bytes32, bytes) {
        bytes memory _nIns;
        bytes memory _ins;
        bytes memory _nOuts;
        bytes memory _outs;
        bytes memory _locktime;
        bytes32 _txid;

        // We don't use checkproof here because we need access to the parse info
        (_nIns, _ins, _nOuts, _outs, _locktime, _txid) = ValidateSPV.parseTransaction(_bitcoinTx);
        // require(_txid != bytes32(0), 'Failed tx parsing');
        
        require(evaluateProofDifficulty(_bitcoinHeaders));
        SPVMerkelCheck(_txid, _bitcoinHeaders, _merkleProof, _index);
        require(keccak256(_locktime) == keccak256(hex'00000000'), 'Wrong locktime set');
        require(keccak256(_nIns) == keccak256(hex'01'), 'Too many ins');
        require(keccak256(_nOuts) == keccak256(hex'01'), 'Too many outs');
        require(keccak256(BTCUtils.extractOutpoint(_ins)) == keccak256(_utxoOutpoint),
                'Tx spends the wrong UTXO');
        require(keccak256(BTCUtils.extractHash(_outs)) == keccak256(abi.encodePacked(_requesterPKH)),
                'Tx sends value to wrong pubkeyhash');
        /* TODO: refactor redemption flow to improve this */
    
        return (_txid, _outs);
    }

    function SPVMerkelCheck(bytes32 _txid, bytes _bitcoinHeaders, bytes _merkleProof, uint256 _index) public pure {
        require(
            ValidateSPV.prove(
                _txid,
                BytesLib.toBytes32(BTCUtils.extractMerkleRootLE(_bitcoinHeaders)),
                _merkleProof,
                _index),
            'Tx merkle proof is not valid for provided header');
    }
}