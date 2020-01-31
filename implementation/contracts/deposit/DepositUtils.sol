pragma solidity ^0.5.10;

import {ValidateSPV} from "@summa-tx/bitcoin-spv-sol/contracts/ValidateSPV.sol";
import {SafeMath} from "@summa-tx/bitcoin-spv-sol/contracts/SafeMath.sol";
import {BTCUtils} from "@summa-tx/bitcoin-spv-sol/contracts/BTCUtils.sol";
import {BytesLib} from "@summa-tx/bitcoin-spv-sol/contracts/BytesLib.sol";
import {DepositStates} from "./DepositStates.sol";
import {TBTCConstants} from "./TBTCConstants.sol";
import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {IERC721} from "openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import {IECDSAKeep} from "@keep-network/keep-ecdsa/contracts/api/IECDSAKeep.sol";
import {IBondedECDSAKeep} from "../external/IBondedECDSAKeep.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";

library DepositUtils {

    using SafeMath for uint256;
    using BytesLib for bytes;
    using BTCUtils for bytes;
    using BTCUtils for uint256;
    using ValidateSPV for bytes;
    using ValidateSPV for bytes32;
    using DepositStates for DepositUtils.Deposit;

    struct Deposit {

        // SET DURING CONSTRUCTION
        address TBTCSystem;
        address TBTCToken;
        address TBTCDepositToken;
        address FeeRebateToken;
        address VendingMachine;
        uint256 lotSizeSatoshis;
        uint8 currentState;
        uint256 signerFeeDivisor;
        uint128 undercollateralizedThresholdPercent;
        uint128 severelyUndercollateralizedThresholdPercent;

        // SET ON FRAUD
        uint256 liquidationInitiated;  // Timestamp of when liquidation starts
        uint256 courtesyCallInitiated; // When the courtesy call is issued
        address payable liquidationInitiator;

        // written when we request a keep
        address keepAddress;  // The address of our keep contract
        uint256 signingGroupRequestedAt;  // timestamp of signing group request

        // written when we get a keep result
        uint256 fundingProofTimerStart;  // start of the funding proof period. reused for funding fraud proof period
        bytes32 signingGroupPubkeyX;  // The X coordinate of the signing group's pubkey
        bytes32 signingGroupPubkeyY;  // The Y coordinate of the signing group's pubkey

        // INITIALLY WRITTEN BY REDEMPTION FLOW
        address payable redeemerAddress;  // The redeemer's address, used as fallback for fraud in redemption
        bytes20 redeemerPKH;  // The 20-byte redeemer PKH
        uint256 initialRedemptionFee;  // the initial fee as requested
        uint256 withdrawalRequestTime;  // the most recent withdrawal request timestamp
        bytes32 lastRequestedDigest;  // the digest most recently requested for signing

        // written when we get funded
        bytes8 utxoSizeBytes;  // LE uint. the size of the deposit UTXO in satoshis
        uint256 fundedAt; // timestamp when funding proof was received
        bytes utxoOutpoint;  // the 36-byte outpoint of the custodied UTXO

        /// @notice Map of timestamps for transaction digests approved for signing
        /// @dev Holds a timestamp from the moment when the transaction digest
        /// was approved for signing
        mapping (bytes32 => uint256) approvedDigests;
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
    function evaluateProofDifficulty(Deposit storage _d, bytes memory _bitcoinHeaders) public view {
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

        require(_observedDiff != ValidateSPV.getErrBadLength(), "Invalid length of the headers chain");
        require(_observedDiff != ValidateSPV.getErrInvalidChain(), "Invalid headers chain");
        require(_observedDiff != ValidateSPV.getErrLowWork(), "Insufficient work in a header");

        /* TODO: make this better than 6 */
        require(
            _observedDiff >= _reqDiff.mul(TBTCConstants.getTxProofDifficultyFactor()),
            "Insufficient accumulated difficulty in header chain"
        );
    }

    /// @notice                 Syntactically check an SPV proof for a bitcoin transaction with its hash (ID)
    /// @dev                    Stateless SPV Proof verification documented elsewhere (see github.com/summa-tx/bitcoin-spv)
    /// @param _d               Deposit storage pointer
    /// @param _txId            The bitcoin txid of the tx that is purportedly included in the header chain
    /// @param _merkleProof     The merkle proof of inclusion of the tx in the bitcoin block
    /// @param _txIndexInBlock  The index of the tx in the Bitcoin block (0-indexed)
    /// @param _bitcoinHeaders  An array of tightly-packed bitcoin headers
    function checkProofFromTxId(
        Deposit storage _d,
        bytes32 _txId,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public view{
        require(
            _txId.prove(
                _bitcoinHeaders.extractMerkleRootLE().toBytes32(),
                _merkleProof,
                _txIndexInBlock
            ),
            "Tx merkle proof is not valid for provided header and txId");
        // TODO: Update for variable confirmation requirements via Vending Machine.
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
        bytes memory _txOutputVector,
        uint8 _fundingOutputIndex
    ) public view returns (bytes8) {
        bytes8 _valueBytes;
        bytes memory _output;

        // Find the output paying the signer PKH
        _output = _txOutputVector.extractOutputAtIndex(_fundingOutputIndex);

        if (keccak256(_output.extractHash()) == keccak256(abi.encodePacked(signerPKH(_d)))) {
            _valueBytes = bytes8(_output.slice(0, 8).toBytes32());
            return _valueBytes;
        }
        // If we don't return from inside the loop, we failed.
        revert("could not identify output funding the required public key hash");
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
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed)
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first
    /// @return                     The 8-byte LE UTXO size in satoshi, the 36byte outpoint
    function validateAndParseFundingSPVProof(
        DepositUtils.Deposit storage _d,
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public view returns (bytes8 _valueBytes, bytes memory _utxoOutpoint){
        require(_txInputVector.validateVin(), "invalid input vector provided");
        require(_txOutputVector.validateVout(), "invalid output vector provided");

        bytes32 txID = abi.encodePacked(_txVersion, _txInputVector, _txOutputVector, _txLocktime).hash256();

        _valueBytes = findAndParseFundingOutput(_d, _txOutputVector, _fundingOutputIndex);

        require(bytes8LEToUint(_valueBytes) >= _d.lotSizeSatoshis, "Deposit too small");

        checkProofFromTxId(_d, txID, _merkleProof, _txIndexInBlock, _bitcoinHeaders);

        // The utxoOutpoint is the LE txID plus the index of the output as a 4-byte LE int
        // _fundingOutputIndex is a uint8, so we know it is only 1 byte
        // Therefore, pad with 3 more bytes
        _utxoOutpoint = abi.encodePacked(txID, _fundingOutputIndex, hex"000000");
    }

    /// @notice Retreive the remaining term of the deposit
    /// @dev    The value is not guaranteed since block.timestmap can be lightly manipulated by miners.
    /// @return The remaining term of the deposit in seconds. 0 if already at term
    function remainingTerm(DepositUtils.Deposit storage _d) public view returns(uint256){
        uint256 endOfTerm = _d.fundedAt + TBTCConstants.getDepositTerm();
        if(block.timestamp < endOfTerm ) {
            return endOfTerm - block.timestamp;
        }
        return 0;
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

    /// @notice         Gets the lot size in erc20 decimal places (max 18) 
    /// @return         uint256 lot size in erc20 
    function lotSizeTbtc(Deposit storage _d) public view returns (uint256){
        return _d.lotSizeSatoshis * TBTCConstants.getSatoshiMultiplier();
    }

    /// @notice         Determines the fees due to the signers for work performed
    /// @dev            Signers are paid based on the TBTC issued
    /// @return         Accumulated fees in smallest TBTC unit (tsat)
    function signerFee(Deposit storage _d) public view returns (uint256) {
        return lotSizeTbtc(_d).div(_d.signerFeeDivisor);
    }

    /// @notice     Determines the amount of TBTC accepted in the auction
    /// @dev        If redeemerAddress is non-0, that means we came from redemption, and no auction should happen
    /// @return     The amount of TBTC that must be paid at auction for the signer's bond
    function auctionTBTCAmount(Deposit storage _d) public view returns (uint256) {
        if (_d.redeemerAddress == address(0)) {
            return lotSizeTbtc(_d);
        } else {
            return 0;
        }
    }

    /// @notice             Determines the prefix to the compressed public key
    /// @dev                The prefix encodes the parity of the Y coordinate
    /// @param  _pubkeyY    The Y coordinate of the public key
    /// @return             The 1-byte prefix for the compressed key
    function determineCompressionPrefix(bytes32 _pubkeyY) public pure returns (bytes memory) {
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
    function compressPubkey(bytes32 _pubkeyX, bytes32 _pubkeyY) public pure returns (bytes memory) {
        return abi.encodePacked(determineCompressionPrefix(_pubkeyY), _pubkeyX);
    }

    /// @notice         Returns the packed public key (64 bytes) for the signing group
    /// @dev            We store it as 2 bytes32, (2 slots) then repack it on demand
    /// @return         64 byte public key
    function signerPubkey(Deposit storage _d) public view returns (bytes memory) {
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

    /// @notice     Gets the current price of Bitcoin in Ether
    /// @dev        Polls the price feed via the system contract
    /// @return     The current price of 1 sat in wei
    function fetchBitcoinPrice(Deposit storage _d) public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(_d.TBTCSystem);
        return _sys.fetchBitcoinPrice();
    }

    /// @notice     Fetches the Keep's bond amount in wei
    /// @dev        Calls the keep contract to do so
    /// @return     The amount of bonded ETH in wei
    function fetchBondAmount(Deposit storage _d) public view returns (uint256) {
        IBondedECDSAKeep _keep = IBondedECDSAKeep(_d.keepAddress);
        return _keep.checkBondAmount(_d.keepAddress);
    }

    /// @notice         Convert a LE bytes8 to a uint256
    /// @dev            Do this by converting to bytes, then reversing endianness, then converting to int
    /// @return         The uint256 represented in LE by the bytes8
    function bytes8LEToUint(bytes8 _b) public pure returns (uint256) {
        return abi.encodePacked(_b).reverseEndianness().bytesToUint();
    }

    /// @notice         Gets timestamp of digest approval for signing
    /// @dev            Identifies entry in the recorded approvals by keep ID and digest pair
    /// @param _digest  Digest to check approval for
    /// @return         Timestamp from the moment of recording the digest for signing.
    ///                 Returns 0 if the digest was not approved for signing
    function wasDigestApprovedForSigning(Deposit storage _d, bytes32 _digest) public view returns (uint256) {
        return _d.approvedDigests[_digest];
    }

    /// @notice         Looks up the Fee Rebate Token holder.
    /// @return         The current token holder if the Token exists.
    ///                 address(0) if the token does not exist.
    function feeRebateTokenHolder(Deposit storage _d) public view returns (address payable) {
        FeeRebateToken _feeRebateToken = FeeRebateToken(_d.FeeRebateToken);
        address tokenHolder;
        if(_feeRebateToken.exists(uint256(address(this)))){
            tokenHolder = address(uint160(_feeRebateToken.ownerOf(uint256(address(this)))));
        }
        return address(uint160(tokenHolder));
    }

    /// @notice         Looks up the deposit beneficiary by calling the tBTC system
    /// @dev            We cast the address to a uint256 to match the 721 standard
    /// @return         The current deposit beneficiary
    function depositOwner(Deposit storage _d) public view returns (address payable) {
        IERC721 _tbtcDepositToken = IERC721(_d.TBTCDepositToken);
        return address(uint160(_tbtcDepositToken.ownerOf(uint256(address(this)))));
    }

    /// @notice     Deletes state after termination of redemption process
    /// @dev        We keep around the redeemer address so we can pay them out
    function redemptionTeardown(Deposit storage _d) public {
        // don't 0 redeemerAddress because we use it to calculate auctionTBTCAmount
        _d.redeemerPKH = bytes20(0);
        _d.initialRedemptionFee = 0;
        _d.withdrawalRequestTime = 0;
        _d.lastRequestedDigest = bytes32(0);
    }

    /// @notice     Seize the signer bond from the keep contract
    /// @dev        we check our balance before and after
    /// @return     the amount of ether seized
    function seizeSignerBonds(Deposit storage _d) internal returns (uint256) {
        uint256 _preCallBalance = address(this).balance;
        IBondedECDSAKeep _keep = IBondedECDSAKeep(_d.keepAddress);
        _keep.seizeSignerBonds(_d.keepAddress);
        uint256 _postCallBalance = address(this).balance;
        require(_postCallBalance > _preCallBalance, "No funds received, unexpected");
        return _postCallBalance.sub(_preCallBalance);
    }

    /// @notice     Distributes the fee rebate to the Fee Rebate Token owner
    ///             whenever this is called we are shutting down.
    function distributeFeeRebate(Deposit storage _d) internal {
        TBTCToken _tbtc = TBTCToken(_d.TBTCToken);

        address rebateTokenHolder = feeRebateTokenHolder(_d);

        // We didn't escrow a rebate if the redeemer is also the Fee Rebate Token holder
        if(_d.redeemerAddress == rebateTokenHolder) return;

        // pay out the rebate if it is available
        if(_tbtc.balanceOf(address(this)) >= signerFee(_d)) {
            _tbtc.transfer(rebateTokenHolder, signerFee(_d));
        }
    }

    /// @notice             pushes ether held by the deposit to the signer group
    /// @dev                useful for returning bonds to the group, or otherwise paying them
    /// @param  _ethValue   the amount of ether to send
    /// @return             true if successful, otherwise revert
    function pushFundsToKeepGroup(Deposit storage _d, uint256 _ethValue) internal returns (bool) {
        require(address(this).balance >= _ethValue, "Not enough funds to send");
        IECDSAKeep _keep = IECDSAKeep(_d.keepAddress);
        _keep.distributeETHToMembers.value(_ethValue)();
        return true;
    }

    /// @notice             Get TBTC amount required for redemption assuming _redeemer
    ///                     is this deposit's TDT owner.
    /// @param _redeemer    The assumed owner of the deposit's TDT 
    /// @return             The amount in TBTC needed to redeem the deposit.
    function getOwnerRedemptionTbtcRequirement(DepositUtils.Deposit storage _d, address _redeemer) internal view returns(uint256) {
        uint256 fee = signerFee(_d);
        bool inCourtesy = _d.inCourtesyCall();
        if(remainingTerm(_d) > 0 && !inCourtesy){
            if(feeRebateTokenHolder(_d) != _redeemer) {
                return fee;
            }
        }
        uint256 contractTbtcBalance = TBTCToken(_d.TBTCToken).balanceOf(address(this));
        if(contractTbtcBalance < fee) {
            return fee.sub(contractTbtcBalance);
        }
        return 0;
    }

    /// @notice             Get TBTC amount required by redemption by a specified _redeemer
    /// @dev                Will revert if redemption is not possible by msg.sender.
    /// @param _redeemer    The deposit redeemer. 
    /// @return             The amount in TBTC needed to redeem the deposit.
    function getRedemptionTbtcRequirement(DepositUtils.Deposit storage _d, address _redeemer) internal view returns(uint256) {
        bool inCourtesy = _d.inCourtesyCall();
        if (depositOwner(_d) == _redeemer && !inCourtesy) {
            return getOwnerRedemptionTbtcRequirement(_d, _redeemer);
        }
        require(remainingTerm(_d) == 0 || inCourtesy, "Only TDT owner can redeem unless deposit is at-term or in COURTESY_CALL");
        return lotSizeTbtc(_d);
    }
}
