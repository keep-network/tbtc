pragma solidity 0.4.25;

import {ValidateSPV} from "../bitcoin-spv/ValidateSPV.sol";
import {SafeMath} from "../bitcoin-spv/SafeMath.sol";
import {BTCUtils} from "../bitcoin-spv/BTCUtils.sol";
import {BytesLib} from "../bitcoin-spv/BytesLib.sol";
import {TBTCConstants} from "./TBTCConstants.sol";
import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {IERC721} from "../interfaces/IERC721.sol";
import {IKeep} from "../interfaces/IKeep.sol";
import {IBurnableERC20} from "../interfaces/IBurnableERC20.sol";

library DepositUtils {

    using SafeMath for uint256;
    using BytesLib for bytes;
    using BTCUtils for bytes;
    using BTCUtils for uint256;
    using ValidateSPV for bytes;
    using ValidateSPV for bytes32;

    struct Deposit {

        // SET DURING CONSTRUCTION
        address TBTCSystem;
        address TBTCToken;
        address KeepBridge;
        uint8 currentState;

        // SET ON FRAUD
        uint256 liquidationInitiated;  // Timestamp of when liquidation starts
        uint256 courtesyCallInitiated; // When the courtesy call is issued

        // written when we request a keep
        uint256 keepID;  // The ID of our keep group
        uint256 signingGroupRequestedAt;  // timestamp of signing group request

        // written when we get a keep result
        uint256 fundingProofTimerStart;  // start of the funding proof period. reused for funding fraud proof period
        bytes32 signingGroupPubkeyX;  // The X coordinate of the signing group's pubkey
        bytes32 signingGroupPubkeyY;  // The Y coordinate of the signing group's pubkey

        // INITIALLY WRITTEN BY REDEMPTION FLOW
        address requesterAddress;  // The requester's addr4ess, used as fallback for fraud in redemption
        bytes20 requesterPKH;  // The 20-byte requeser PKH
        uint256 initialRedemptionFee;  // the initial fee as requested
        uint256 withdrawalRequestTime;  // the most recent withdrawal request timestamp
        bytes32 lastRequestedDigest;  // the digest most recently requested for signing

        // written when we get funded
        bytes8 utxoSizeBytes;  // LE uint. the size of the deposit UTXO in satoshis
        uint256 fundedAt; // timestamp when funding proof was received
        bytes utxoOutpoint;  // the 36-byte outpoint of the custodied UTXO
    }

    /// @notice         Gets the current block difficulty
    /// @dev            Calls the light relay and gets the current block difficulty
    /// @return         The difficulty
    function currentBlockDifficulty(Deposit storage _d) public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(_d.TBTCSystem);
        return _sys.fetchRelayCurrentDifficulty();
    }

    /// @notice         Gets the previous block difficulty
    /// @dev            Calls the light relay and gets the previous block difficulty
    /// @return         The difficulty
    function previousBlockDifficulty(Deposit storage _d) public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(_d.TBTCSystem);
        return _sys.fetchRelayPreviousDifficulty();
    }

    /// @notice                     Evaluates the header difficulties in a proof
    /// @dev                        Uses the light oracle to source recent difficulty
    /// @param  _bitcoinHeaders     The header chain to evaluate
    /// @return                     True if acceptable, otherwise revert
    function evaluateProofDifficulty(Deposit storage _d, bytes _bitcoinHeaders) public view {
        uint256 _reqDiff;
        uint256 _current = currentBlockDifficulty(_d);
        uint256 _previous = previousBlockDifficulty(_d);
        uint256 _firstHeaderDiff = _bitcoinHeaders.extractTarget().calculateDifficulty();

        if (_firstHeaderDiff == _current) {
            _reqDiff = _current;
        } else if (_firstHeaderDiff == _previous) {
            _reqDiff = _previous;
        } else {
            revert("not at current or previous difficulty");
        }

        uint256 _observedDiff = _bitcoinHeaders.validateHeaderChain();
        require(_observedDiff > 3, "ValidateSPV returned an error code");
        /* TODO: make this better than 6 */
        require(
            _observedDiff >= _reqDiff.mul(6),
            "Insufficient accumulated difficulty in header chain"
        );
    }

    /// @notice                 Syntactically check an SPV proof for a bitcoin tx
    /// @dev                    Stateless SPV Proof verification documented elsewhere
    /// @param _d               Deposit storage pointer
    /// @param _bitcoinTx       The bitcoin tx that is purportedly included in the header chain
    /// @param _merkleProof     The merkle proof of inclusion of the tx in the bitcoin block
    /// @param _txIndexInBlock  The index of the tx in the Bitcoin block (1-indexed)
    /// @param _bitcoinHeaders  An array of tightly-packed bitcoin headers
    /// @return                 The 32 byte transaction id (little-endian, not block-explorer)
    function checkProofFromTx(
        Deposit storage _d,
        bytes _bitcoinTx,
        bytes _merkleProof,
        uint256 _txIndexInBlock,
        bytes _bitcoinHeaders
    ) public view returns (bytes32) {
        bytes memory _nIns;
        bytes memory _ins;
        bytes memory _nOuts;
        bytes memory _outs;
        bytes memory _locktime;
        bytes32 _txid;
        (_nIns, _ins, _nOuts, _outs, _locktime, _txid) = _bitcoinTx.parseTransaction();
        require(_txid != bytes32(0), "Failed tx parsing");
        checkProofFromTxId(_d, _txid, _merkleProof, _txIndexInBlock, _bitcoinHeaders);
        return _txid;
    }

    /// @notice                 Syntactically check an SPV proof for a bitcoin transaction with its hash (ID)
    /// @dev                    Stateless SPV Proof verification documented elsewhere (see github.com/summa-tx/bitcoin-spv)
    /// @param _d               Deposit storage pointer
    /// @param _txId            The bitcoin txid of the tx that is purportedly included in the header chain
    /// @param _merkleProof     The merkle proof of inclusion of the tx in the bitcoin block
    /// @param _txIndexInBlock  The index of the tx in the Bitcoin block (1-indexed)
    /// @param _bitcoinHeaders  An array of tightly-packed bitcoin headers
    function checkProofFromTxId(
        Deposit storage _d,
        bytes32 _txId,
        bytes _merkleProof,
        uint256 _txIndexInBlock,
        bytes _bitcoinHeaders
    ) public view{
        require(
            _txId.prove(
                _bitcoinHeaders.extractMerkleRootLE().toBytes32(),
                _merkleProof,
                _txIndexInBlock
            ),
            "Tx merkle proof is not valid for provided header and txId");

        evaluateProofDifficulty(_d, _bitcoinHeaders);
    }

    /// @notice                     Find and validate funding output in transaction output vector using the index
    /// @dev                        Gets `_fundingOutputIndex` output from the output vector and validates if it's
    ///                             Public Key Hash matches a Public Key Hash of the deposit.
    /// @param _d                   Deposit storage pointer
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC outputs
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector
    /// @return                     Funding value
    function findAndParseFundingOutput(
        DepositUtils.Deposit storage _d,
        bytes _txOutputVector,
        uint8 _fundingOutputIndex
    ) public view returns (bytes8) {
        bytes8 _valueBytes;
        bytes memory _output;

        // Find the output paying the signer PKH
        _output = extractOutputAtIndex(_txOutputVector, _fundingOutputIndex);
        if (keccak256(_output.extractHash()) == keccak256(abi.encodePacked(signerPKH(_d)))) {
            _valueBytes = bytes8(_output.slice(0, 8).toBytes32());
            return _valueBytes;
        }
        // If we don't return from inside the loop, we failed.
        revert("could not identify output funding the required public key hash");
    }

    /// @notice                     Extracts the output at a given index in _txOutputVector
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC outputs
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @return                     The specified output
    function extractOutputAtIndex(
        bytes _txOutputVector,
        uint8 _fundingOutputIndex
    ) public view returns (bytes) {
        // Transaction outputs vector consists of a number of outputs followed by a list of outputs:
        //
        // |                                  outputs vector                                  |
        // | outputs number |            output 1            |          output 2...           |
        // | outputs number | value | script length | script | value | script length | script |
        //
        // Each output contains value (8 bytes), script length (VarInt) and a script.

        // extract output number to verify that it's not a varint.
        uint256 _n = (_txOutputVector.slice(0, 1)).bytesToUint();
        require(_n < 0xfd, "VarInts not supported, Number of outputs cannot exceed 252");

        // Determine length of first output
        // offset starts at 1 to skip output number varint
        // skip the 8 byte output value to get to length
        // next two bytes used to calculate length
        uint _offset = 1 + 8;
        uint _length = (_txOutputVector.slice(_offset, 2)).determineOutputLength();

        // This loop moves forward, and then gets the len of the next one
        for (uint i = 0; i < _fundingOutputIndex; i++) {
            _offset = _offset + _length;
            _length = (_txOutputVector.slice(_offset, 2)).determineOutputLength();
        }

        // We now have the length and offset of the one we want
        return _txOutputVector.slice(_offset - 8, _length);
    }

    /// @notice                     Validates the funding tx and parses information from it
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding
    /// @param  _d                  Deposit storage pointer
    /// @param _txVersion           Transaction version number (4-byte LE)
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs
    /// @param _txLocktime          Final 4 bytes of the transaction
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed)
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block
    /// @param _txIndexInBlock      Transaction index in the block (1-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first
    /// @return                     The 8-byte LE UTXO size in satoshi, the 36byte outpoint
    function validateAndParseFundingSPVProof(
        DepositUtils.Deposit storage _d,
        bytes _txVersion,
        bytes _txInputVector,
        bytes _txOutputVector,
        bytes _txLocktime,
        uint8 _fundingOutputIndex,
        bytes _merkleProof,
        uint256 _txIndexInBlock,
        bytes _bitcoinHeaders
    ) public view returns (bytes8 _valueBytes, bytes _utxoOutpoint){
        bytes32 txID = abi.encodePacked(_txVersion, _txInputVector, _txOutputVector, _txLocktime).hash256();

        _valueBytes = findAndParseFundingOutput(_d, _txOutputVector, _fundingOutputIndex);
        require(bytes8LEToUint(_valueBytes) >= TBTCConstants.getLotSize(), "Deposit too small");

        checkProofFromTxId(_d, txID, _merkleProof, _txIndexInBlock, _bitcoinHeaders);

        // The utxoOutpoint is the LE txID plus the index of the output as a 4-byte LE int
        // _fundingOutputIndex is a uint8, so we know it is only 1 byte
        // Therefore, pad with 3 more bytes
        _utxoOutpoint = abi.encodePacked(txID, _fundingOutputIndex, hex"000000");
    }
    
    /// @notice     Calculates the amount of value at auction right now
    /// @dev        We calculate the % of the auction that has elapsed, then scale the value up
    /// @param _d   deposit storage pointer
    /// @return     the value to distribute in the auction at the current time
    function auctionValue(Deposit storage _d) public view returns (uint256) {
        uint256 _elapsed = block.timestamp.sub(_d.liquidationInitiated);
        uint256 _available = address(this).balance;
        if (_elapsed > TBTCConstants.getAuctionDuration()) {
            return _available;
        }

        // This should make a smooth flow from base% to 100%
        uint256 _basePercentage = TBTCConstants.getAuctionBasePercentage();
        uint256 _elapsedPercentage = uint256(100).sub(_basePercentage).mul(_elapsed).div(TBTCConstants.getAuctionDuration());
        uint256 _percentage = _basePercentage + _elapsedPercentage;

        return _available.mul(_percentage).div(100);
    }

    /// @notice         Determines the fees due to the signers for work performeds
    /// @dev            Signers are paid based on the TBTC issued
    /// @return         Accumulated fees in smallest TBTC unit (tsat)
    function signerFee() public pure returns (uint256) {
        return TBTCConstants.getLotSize().div(TBTCConstants.getSignerFeeDivisor());
    }

    /// @notice     calculates the beneficiary reward based on the deposit size
    /// @dev        the amount of extra ether to pay the beneficiary at closing time
    /// @return     the amount of ether in wei to pay the beneficiary
    function beneficiaryReward() public pure returns (uint256) {
        return TBTCConstants.getLotSize().div(TBTCConstants.getBeneficiaryRewardDivisor());
    }

    /// @notice         Determines the amount of TBTC paid to redeem the deposit
    /// @dev            This is the amount of TBTC needed to repay to redeem the Deposit
    /// @return         Outstanding debt in smallest TBTC unit (tsat)
    function redemptionTBTCAmount(Deposit storage _d) public view returns (uint256) {
        if (_d.requesterAddress == address(0)) {
            return TBTCConstants.getLotSize().add(signerFee()).add(beneficiaryReward());
        } else {
            return 0;
        }
    }

    /// @notice     Determines the amount of TBTC accepted in the auction
    /// @dev        If requesterAddress is non-0, that means we came from redemption, and no auction should happen
    /// @return     The amount of TBTC that must be paid at auction for the signer's bond
    function auctionTBTCAmount(Deposit storage _d) public view returns (uint256) {
        if (_d.requesterAddress == address(0)) {
            return TBTCConstants.getLotSize();
        } else {
            return 0;
        }
    }

    /// @notice             Determines the prefix to the compressed public key
    /// @dev                The prefix encodes the parity of the Y coordinate
    /// @param  _pubkeyY    The Y coordinate of the public key
    /// @return             The 1-byte prefix for the compressed key
    function determineCompressionPrefix(bytes32 _pubkeyY) public pure returns (bytes) {
        if(uint256(_pubkeyY) & 1 == 1) {
            return hex"03";  // Odd Y
        } else {
            return hex"02";  // Even Y
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

    /// @notice         Returns the packed public key (64 bytes) for the signing group
    /// @dev            We store it as 2 bytes32, (2 slots) then repack it on demand
    /// @return         64 byte public key
    function signerPubkey(Deposit storage _d) public view returns (bytes) {
        return abi.encodePacked(_d.signingGroupPubkeyX, _d.signingGroupPubkeyY);
    }

    /// @notice         Returns the Bitcoin pubkeyhash (hash160) for the signing group
    /// @dev            This is used in bitcoin output scripts for the signers
    /// @return         20-bytes public key hash
    function signerPKH(Deposit storage _d) public view returns (bytes20) {
        bytes memory _pubkey = compressPubkey(_d.signingGroupPubkeyX, _d.signingGroupPubkeyY);
        bytes memory _digest = _pubkey.hash160();
        return bytes20(_digest.toAddress(0));  // dirty solidity hack
    }

    /// @notice         Returns the size of the deposit UTXO in satoshi
    /// @dev            We store the deposit as bytes8 to make signature checking easier
    /// @return         UTXO value in satoshi
    function utxoSize(Deposit storage _d) public view returns (uint256) {
        return bytes8LEToUint(_d.utxoSizeBytes);
    }

    /// @notice     Gets the current oracle price of Bitcoin in Ether
    /// @dev        Polls the oracle via the system contract
    /// @return     The current price of 1 sat in wei
    function fetchOraclePrice(Deposit storage _d) public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(_d.TBTCSystem);
        return _sys.fetchOraclePrice();
    }

    /// @notice     Fetches the Keep's bond amount in wei
    /// @dev        Calls the keep contract to do so
    /// @return     The amount of bonded ETH in wei
    function fetchBondAmount(Deposit storage _d) public view returns (uint256) {
        IKeep _keep = IKeep(_d.KeepBridge);
        return _keep.checkBondAmount(_d.keepID);
    }

    /// @notice         Convert a LE bytes8 to a uint256
    /// @dev            Do this by converting to bytes, then reversing endianness, then converting to int
    /// @return         The uint256 represented in LE by the bytes8
    function bytes8LEToUint(bytes8 _b) public pure returns (uint256) {
        return abi.encodePacked(_b).reverseEndianness().bytesToUint();
    }

    /// @notice         determines whether a digest has been approved for our keep group
    /// @dev            calls out to the keep contract, storing a 256bit int costs the same as a bool
    /// @param  _digest the digest to check approval time for
    /// @return         the time it was approved. 0 if unapproved
    function wasDigestApprovedForSigning(Deposit storage _d, bytes32 _digest) public view returns (uint256) {
        IKeep _keep = IKeep(_d.KeepBridge);
        return _keep.wasDigestApprovedForSigning(_d.keepID, _digest);
    }

    /// @notice         Looks up the deposit beneficiary by calling the tBTC system
    /// @dev            We cast the address to a uint256 to match the 721 standard
    /// @return         The current deposit beneficiary
    function depositBeneficiary(Deposit storage _d) public view returns (address) {
        IERC721 _systemContract = IERC721(_d.TBTCSystem);
        return _systemContract.ownerOf(uint256(address(this)));
    }

    /// @notice     Deletes state after termination of redemption process
    /// @dev        We keep around the requester address so we can pay them out
    function redemptionTeardown(Deposit storage _d) public {
        // don't 0 requesterAddress because we use it to calculate auctionTBTCAmount
        _d.requesterPKH = bytes20(0);
        _d.initialRedemptionFee = 0;
        _d.withdrawalRequestTime = 0;
        _d.lastRequestedDigest = bytes32(0);
    }

    /// @notice     Seize the signer bond from the keep contract
    /// @dev        we check our balance before and after
    /// @return     the amount of ether seized
    function seizeSignerBonds(Deposit storage _d) public returns (uint256) {
        uint256 _preCallBalance = address(this).balance;
        IKeep _keep = IKeep(_d.KeepBridge);
        _keep.seizeSignerBonds(_d.keepID);
        uint256 _postCallBalance = address(this).balance;
        require(_postCallBalance > _preCallBalance, "No funds received, unexpected");
        return _postCallBalance.sub(_preCallBalance);
    }

    /// @notice     Distributes the beneficiary reward to the beneficiary
    /// @dev        We distribute the whole TBTC balance as a convenience,
    ///             whenever this is called we are shutting down.
    function distributeBeneficiaryReward(Deposit storage _d) public {
        IBurnableERC20 _tbtc = IBurnableERC20(_d.TBTCToken);
        /* solium-disable-next-line */
        require(_tbtc.transfer(depositBeneficiary(_d), _tbtc.balanceOf(address(this))));
    }

    /// @notice             pushes ether held by the deposit to the signer group
    /// @dev                useful for returning bonds to the group, or otherwise paying them
    /// @param  _ethValue   the amount of ether to send
    /// @return             true if successful, otherwise revert
    function pushFundsToKeepGroup(Deposit storage _d, uint256 _ethValue) public returns (bool) {
        require(address(this).balance >= _ethValue, "Not enough funds to send");
        IKeep _keep = IKeep(_d.KeepBridge);
        return _keep.distributeEthToKeepGroup.value(_ethValue)(_d.keepID);
    }
}
