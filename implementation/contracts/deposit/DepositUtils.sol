pragma solidity 0.4.25;

import {ValidateSPV} from "../bitcoin-spv/ValidateSPV.sol";
import {SafeMath} from "../bitcoin-spv/SafeMath.sol";
import {BTCUtils} from "../bitcoin-spv/BTCUtils.sol";
import {BytesLib} from "../bitcoin-spv/BytesLib.sol";
import {TBTCConstants} from './TBTCConstants.sol';
import {ITBTCSystem} from '../interfaces/ITBTCSystem.sol';
import {IERC721} from '../interfaces/IERC721.sol';
import {IKeep} from '../interfaces/IKeep.sol';
import {IBurnableERC20} from '../interfaces/IBurnableERC20.sol';

library DepositUtils {

    using SafeMath for uint256;
    using BytesLib for bytes;
    using BTCUtils for bytes;
    using BTCUtils for uint256;
    using ValidateSPV for bytes;
    using ValidateSPV for bytes32;

    struct Deposit {

        // SET DURING CONSTRUCTION
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
    ) public view returns (bytes32) {
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

        require(_bitcoinHeaders.validateHeaderChain() > currentBlockDifficulty().mul(6),
                'Insufficient accumulated difficulty in header chain');

        return _txid;
    }

    /// @notice         Check if the caller is the tBTC system contract
    /// @dev            Stored as a constant in the config library
    /// @param _caller  The address of the caller to compare to the tbtc system constant
    /// @return         True if the caller is approved, else False
    function isTBTCSystemContract(address _caller) public pure returns (bool) {
        return _caller == TBTCConstants.getSystemContractAddress();
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
        return lotSize().div(TBTCConstants.getSignerFeeDivisor());
    }

    /// @notice     calculates the beneficiary reward based on the deposit size
    /// @dev        the amount of extra ether to pay the beneficiary at closing time
    /// @return     the amount of ether in wei to pay the beneficiary
    function beneficiaryReward() public pure returns (uint256) {
        return lotSize().div(TBTCConstants.getBeneficiaryRewardDivisor());
    }

    /// @notice         Determines the amount of TBTC paid to redeem the deposit
    /// @dev            This is the amount of TBTC needed to repay to redeem the Deposit
    /// @return         Outstanding debt in smallest TBTC unit (tsat)
    function redemptionTBTCAmount(Deposit storage _d) public view returns (uint256) {
        if (_d.requesterAddress == address(0)) {
            return lotSize().add(signerFee()).add(beneficiaryReward());
        } else {
            return 0;
        }
    }

    /// @notice     Determines the amount of TBTC accepted in the auction
    /// @dev        If requesterAddress is non-0, that means we came from redemption, and no auction should happen
    /// @return     The amount of TBTC that must be paid at auction for the signer's bond
    function auctionTBTCAmount(Deposit storage _d) public view returns (uint256) {
        if (_d.requesterAddress == address(0)) {
            return lotSize();
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

    /// @notice         Returns the size of the standard lot
    /// @dev            This is the amount of TBTC issued, and the minimum amount of BTC in the utxo
    /// @return         lot size value in tsat
    function lotSize() public pure returns (uint256) {
        return TBTCConstants.getLotSize();
    }

    /// @notice         Returns the size of the deposit UTXO in satoshi
    /// @dev            We store the deposit as bytes8 to make signature checking easier
    /// @return         UTXO value in satoshi
    function utxoSize(Deposit storage _d) public view returns (uint256) {
        return bytes8LEToUint(_d.utxoSizeBytes);
    }

    /// @notice     Looks up the size of the funder bond
    /// @dev        This is stored as a constant
    /// @return     The refundable portion of the funder bond
    function funderBondAmount() public pure returns (uint256) {
        TBTCConstants.getFunderBondAmount();
    }


    /// @notice     Gets the current oracle price of Bitcoin in Ether
    /// @dev        Polls the oracle via the system contract
    /// @return     The current price of 1 sat in wei
    function fetchOraclePrice() public view returns (uint256) {
        ITBTCSystem _sys = ITBTCSystem(TBTCConstants.getSystemContractAddress());
        return _sys.fetchOraclePrice();
    }

    /// @notice     Fetches the Keep's bond amount in wei
    /// @dev        Calls the keep contract to do so
    /// @return     The amount of bonded ETH in wei
    function fetchBondAmount(Deposit storage _d) public view returns (uint256) {
        IKeep _keep = IKeep(TBTCConstants.getKeepContractAddress());
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
        IKeep _keep = IKeep(TBTCConstants.getKeepContractAddress());
        return _keep.wasDigestApprovedForSigning(_d.keepID, _digest);
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

    /// @notice                     Evaluates the header difficulties in a proof
    /// @dev                        Uses the light oracle to source recent difficulty
    /// @param  _bitcoinHeaders     The header chain to evaluate
    /// @return                     True if acceptable, otherwise revert
    function evaluateProofDifficulty(bytes _bitcoinHeaders) public view returns (bool) {
        uint256 _reqDiff;
        uint256 _current = currentBlockDifficulty();
        uint256 _previous = previousBlockDifficulty();
        uint256 _firstHeaderDiff = _bitcoinHeaders.extractTarget().calculateDifficulty();

        if (_firstHeaderDiff == _current) {
            _reqDiff = _current;
        } else if (_firstHeaderDiff == _previous) {
            _reqDiff = _previous;
        } else {
            revert('not at current or previous difficulty');
        }

        /* TODO: make this better than 6 */
        require(_bitcoinHeaders.validateHeaderChain() > _reqDiff.mul(6),
                'Insufficient accumulated difficulty in header chain');
        return true;
    }

    /// @notice         Looks up the deposit beneficiary by calling the tBTC system
    /// @dev            We cast the address to a uint256 to match the 721 standard
    /// @return         The current deposit beneficiary
    function depositBeneficiary() public view returns (address) {
        IERC721 _systemContract = IERC721(TBTCConstants.getSystemContractAddress());
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
    function seizeSignerBonds(DepositUtils.Deposit storage _d) public returns (uint256) {
        uint256 _preCallBalance = address(this).balance;
        IKeep _keep = IKeep(TBTCConstants.getKeepContractAddress());
        _keep.seizeSignerBonds(_d.keepID);
        uint256 _postCallBalance = address(this).balance;
        require(_postCallBalance > _preCallBalance, 'No funds received, unexpected');
        return _postCallBalance.sub(_preCallBalance);
    }

    /// @notice     Distributes the beneficiary reward to the beneficiary
    /// @dev        We distribute the whole TBTC balance as a convenience,
    ///             whenever this is called we are shutting down.
    function distributeBeneficiaryReward() public {
        IBurnableERC20 _tbtc = IBurnableERC20(TBTCConstants.getTokenContractAddress());
        require(_tbtc.transfer(depositBeneficiary(), _tbtc.balanceOf(address(this))));
    }

    /// @notice             pushes ether held by the deposit to the signer group
    /// @dev                useful for returning bonds to the group, or otherwise paying them
    /// @param  _ethValue   the amount of ether to send
    /// @return             true if successful, otherwise revert
    function pushFundsToKeepGroup(DepositUtils.Deposit storage _d, uint256 _ethValue) public returns (bool) {
        require(address(this).balance >= _ethValue, 'Not enough funds to send');
        IKeep _keep = IKeep(TBTCConstants.getKeepContractAddress());
        return _keep.distributeEthToKeepGroup.value(_ethValue)(_d.keepID);
    }
}
