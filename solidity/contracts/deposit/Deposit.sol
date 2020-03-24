pragma solidity ^0.5.10;

import {DepositLiquidation} from "./DepositLiquidation.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {DepositFunding} from "./DepositFunding.sol";
import {DepositRedemption} from "./DepositRedemption.sol";
import {DepositStates} from "./DepositStates.sol";
import "../system/DepositFactoryAuthority.sol";

/// @title  Deposit.
/// @notice This is the main contract for tBTC. It is the state machine
/// that (through various libraries) handles bitcoin funding,
/// bitcoin-spv proofs, redemption, liquidation, and fraud logic.
/// @dev This is the execution context for libraries:
/// `DepositFunding`, `DepositLiquidaton`, `DepositRedemption`,
/// `DepositStates`, `DepositUtils`, `OutsourceDepositLogging`, and `TBTCConstants`.
contract Deposit is DepositFactoryAuthority {

    using DepositRedemption for DepositUtils.Deposit;
    using DepositFunding for DepositUtils.Deposit;
    using DepositLiquidation for DepositUtils.Deposit;
    using DepositUtils for DepositUtils.Deposit;
    using DepositStates for DepositUtils.Deposit;

    DepositUtils.Deposit self;

    // We separate the constructor from createNewDeposit to make proxy factories easier.
    /* solium-disable-next-line no-empty-blocks */
    constructor () public {
        initialize(address(0));
    }

    function () external payable {
        require(msg.sender == self.keepAddress, "Deposit contract can only receive ETH from rom underlying keep");
    }

    /// @notice     Get the integer representing the current state.
    /// @dev        We implement this because contracts don't handle foreign enums well.
    ///             see DepositStates for more info on states.
    /// @return     The 0-indexed state from the DepositStates enum.
    function getCurrentState() public view returns (uint256) {
        return uint256(self.currentState);
    }

    /// @notice     Check if the Deposit is in ACTIVE state.
    /// @return     True if state is ACTIVE, fale otherwise.
    function inActive() public view returns (bool) {
        return self.inActive();
    }

    /// @notice Retrieve the remaining term of the deposit in seconds.
    /// @dev    The value accuracy is not guaranteed since block.timestmap can
    ///         be lightly manipulated by miners.
    /// @return The remaining term of the deposit in seconds. 0 if already at term.
    function remainingTerm() public view returns(uint256){
        return self.remainingTerm();
    }

    /// @notice     Get the signer fee for the Deposit.
    /// @dev        This is the one-time fee required by the signers to perform .
    ///             the tasks needed to maintain a decentralized and trustless
    ///             model for tBTC. It is a percentage of the lotSize (deposit size).
    /// @return     Fee amount in tBTC.
    function signerFee() public view returns (uint256) {
        return self.signerFee();
    }

    /// @notice     Get the deposit's BTC lot size in satoshi.
    /// @return     uint256 lot size in satoshi.
    function lotSizeSatoshis() public view returns (uint64){
        return self.lotSizeSatoshis;
    }

    /// @notice     Get the Deposit ERC20 lot size.
    /// @dev        This is the same as lotSizeSatoshis(),
    ///             but is multiplied to scale to ERC20 decimal.
    /// @return     uint256 lot size in erc20 decimal (max 18 decimal places).
    function lotSizeTbtc() public view returns (uint256){
        return self.lotSizeTbtc();
    }

    /// @notice     Get the size of the funding UTXO.
    /// @dev        This will only return 0 unless
    ///             the funding transaction has been confirmed on-chain.
    ///             See `provideBTCFundingProof` for more info on the funding proof.
    /// @return     Uint256 UTXO size in satoshi.
    ///             0 if no funding proof has been provided.
    function utxoSize() public view returns (uint256){
        return self.utxoSize();
    }

    // THIS IS THE INIT FUNCTION
    /// @notice        The Deposit Factory can spin up a new deposit.
    /// @dev           Only the Deposit factory can call this.
    /// @param _TBTCSystem        `TBTCSystem` address. More info in `VendingMachine`.
    /// @param _TBTCToken         `TBTCToken` address. More info in TBTCToken`.
    /// @param _TBTCDepositToken  `TBTCDepositToken` (TDT) address. More info in `TBTCDepositToken`.
    /// @param _FeeRebateToken    `FeeRebateToken` (FRT) address. More info in `FeeRebateToken`.
    /// @param _VendingMachine    `VendingMachine` address. More info in `VendingMachine`.
    /// @param _m           Signing group honesty threshold.
    /// @param _n           Signing group size.
    /// @param _lotSizeSatoshis The minimum amount of satoshi the funder is required to send.
    ///                         This is also the amount of TBTC the TDT holder will receive:
    ///                         (10**7 satoshi == 0.1 BTC == 0.1 TBTC).
    /// @return             True if successful, otherwise revert.
    function createNewDeposit(
        address _TBTCSystem,
        address _TBTCToken,
        address _TBTCDepositToken,
        address _FeeRebateToken,
        address _VendingMachine,
        uint16 _m,
        uint16 _n,
        uint64 _lotSizeSatoshis
    ) public onlyFactory payable returns (bool) {
        self.TBTCSystem = _TBTCSystem;
        self.TBTCToken = _TBTCToken;
        self.TBTCDepositToken = _TBTCDepositToken;
        self.FeeRebateToken = _FeeRebateToken;
        self.VendingMachine = _VendingMachine;
        self.createNewDeposit(_m, _n, _lotSizeSatoshis);
        return true;
    }

    /// @notice                     Deposit owner (TDT holder) can request redemption.
    ///                             Once redemption is requested
    ///                             a proof with sufficient accumulated difficulty is
    ///                             required to complete redemption.
    /// @dev                        The redeemer specifies details about the Bitcoin redemption tx.
    /// @param  _outputValueBytes   The 8-byte Little Endian output size.
    /// @param  _redeemerOutputScript The redeemer's length-prefixed output script.
    /// @return                     True if successful, otherwise revert.
    function requestRedemption(
        bytes8 _outputValueBytes,
        bytes memory _redeemerOutputScript
    ) public returns (bool) {
        self.requestRedemption(_outputValueBytes, _redeemerOutputScript);
        return true;
    }

    /// @notice                     Deposit owner (TDT holder) can request redemption.
    ///                             Once redemption is requested a proof with
    ///                             sufficient accumulated difficulty is required
    ///                             to complete redemption.
    /// @dev                        The caller specifies details about the Bitcoin redemption tx and pays
    ///                             for the redemption. The TDT (deposit ownership) is transfered to _finalRecipient, and
    ///                             _finalRecipient is marked as the deposit redeemer.
    /// @param  _outputValueBytes   The 8-byte LE output size.
    /// @param  _redeemerOutputScript The redeemer's length-prefixed output script.
    /// @param  _finalRecipient     The address to receive the TDT and later be recorded as deposit redeemer.
    function transferAndRequestRedemption(
        bytes8 _outputValueBytes,
        bytes memory _redeemerOutputScript,
        address payable _finalRecipient
    ) public returns (bool) {
        self.transferAndRequestRedemption(
            _outputValueBytes,
            _redeemerOutputScript,
            _finalRecipient
        );
        return true;
    }

    /// @notice             Get TBTC amount required for redemption by a specified _redeemer.
    /// @dev                Will revert if redemption is not possible by _redeemer.
    /// @param _redeemer    The deposit redeemer.
    /// @return             The amount in TBTC needed to redeem the deposit.
    function getRedemptionTbtcRequirement(address _redeemer) public view returns(uint256){
        return self.getRedemptionTbtcRequirement(_redeemer);
    }

    /// @notice             Get TBTC amount required for redemption assuming _redeemer
    ///                     is this deposit's owner (TDT holder).
    /// @param _redeemer    The assumed owner of the deposit's TDT .
    /// @return             The amount in TBTC needed to redeem the deposit.
    function getOwnerRedemptionTbtcRequirement(address _redeemer) public view returns(uint256){
        return self.getOwnerRedemptionTbtcRequirement(_redeemer);
    }

    /// @notice     Anyone may provide a withdrawal signature if it was requested.
    /// @dev        The signers will be penalized if this (or provideRedemptionProof) is not called.
    /// @param  _v  Signature recovery value.
    /// @param  _r  Signature R value.
    /// @param  _s  Signature S value. Should be in the low half of secp256k1 curve's order.
    /// @return     True if successful, False if prevented by timeout, otherwise revert.
    function provideRedemptionSignature(
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) public returns (bool) {
        self.provideRedemptionSignature(_v, _r, _s);
        return true;
    }

    /// @notice                             Anyone may notify the contract that a fee bump is needed.
    /// @dev                                This sends us back to AWAITING_WITHDRAWAL_SIGNATURE.
    /// @param  _previousOutputValueBytes   The previous output's value.
    /// @param  _newOutputValueBytes        The new output's value.
    /// @return                             True if successful, False if prevented by timeout, otherwise revert.
    function increaseRedemptionFee(
        bytes8 _previousOutputValueBytes,
        bytes8 _newOutputValueBytes
    ) public returns (bool) {
        return self.increaseRedemptionFee(_previousOutputValueBytes, _newOutputValueBytes);
    }

    /// @notice                 Anyone may provide a withdrawal proof to prove redemption.
    /// @dev                    The signers will be penalized if this is not called.
    /// @param  _txVersion      Transaction version number (4-byte Little Endian).
    /// @param  _txInputVector  All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs.
    /// @param  _txOutputVector All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs.
    /// @param  _txLocktime     Final 4 bytes of the transaction.
    /// @param  _merkleProof    The merkle proof of inclusion of the tx in the bitcoin block.
    /// @param  _txIndexInBlock The index of the tx in the Bitcoin block (0-indexed).
    /// @param  _bitcoinHeaders An array of tightly-packed bitcoin headers.
    function provideRedemptionProof(
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public returns (bool) {
        self.provideRedemptionProof(
            _txVersion,
            _txInputVector,
            _txOutputVector,
            _txLocktime,
            _merkleProof,
            _txIndexInBlock,
            _bitcoinHeaders
        );
        return true;
    }

    /// @notice     Anyone may notify the contract that the signers have failed to produce a signature.
    /// @dev        This is considered fraud, and is punished.
    /// @return     True if successful, otherwise revert.
    function notifySignatureTimeout() public returns (bool) {
        self.notifySignatureTimeout();
        return true;
    }

    /// @notice     Anyone may notify the contract that the signers have failed to produce a redemption proof.
    /// @dev        This is considered fraud, and is punished.
    /// @return     True if successful, otherwise revert.
    function notifyRedemptionProofTimeout() public returns (bool) {
        self.notifyRedemptionProofTimeout();
        return true;
    }

    //
    // FUNDING FLOW
    //

    /// @notice     Anyone may notify the contract that signing group setup has timed out.
    /// @dev        We rely on the keep system punishes the signers in this case.
    /// @return     True if successful, otherwise revert.
    function notifySignerSetupFailure() public returns (bool) {
        self.notifySignerSetupFailure();
        return true;
    }

    /// @notice             Poll the Keep contract to retrieve our pubkey.
    /// @dev                Store the pubkey as 2 bytestrings, X and Y.
    /// @return             True if successful, otherwise revert.
    function retrieveSignerPubkey() public returns (bool) {
        self.retrieveSignerPubkey();
        return true;
    }

    /// @notice     Anyone may notify the contract that the funder has failed to send BTC.
    /// @dev        This is considered a funder fault, and we revoke their bond.
    /// @return     True if successful, otherwise revert.
    function notifyFundingTimeout() public returns (bool) {
        self.notifyFundingTimeout();
        return true;
    }

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud during funding.
    /// @dev                    Calls out to the keep to verify if there was fraud.
    /// @param  _v              Signature recovery value.
    /// @param  _r              Signature R value.
    /// @param  _s              Signature S value.
    /// @param _signedDigest    The digest signed by the signature vrs tuple.
    /// @param _preimage        The sha256 preimage of the digest.
    /// @return                 True if successful, otherwise revert.
    function provideFundingECDSAFraudProof(
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes memory _preimage
    ) public returns (bool) {
        self.provideFundingECDSAFraudProof(_v, _r, _s, _signedDigest, _preimage);
        return true;
    }

    /// @notice     Anyone may notify the contract no funding proof was submitted during funding fraud.
    /// @dev        This is not a funder fault. The signers have faulted, so the funder shouldn't fund.
    /// @return     True if successful, otherwise revert.
    function notifyFraudFundingTimeout() public returns (bool) {
        self.notifyFraudFundingTimeout();
        return true;
    }

    /// @notice                     Anyone may notify the deposit of a funding proof during funding fraud.
    //                              We reward the funder the entire bond if this occurs.
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding.
    /// @param _txVersion           Transaction version number (4-byte LE).
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs.
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs.
    /// @param _txLocktime          Final 4 bytes of the transaction.
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed).
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block.
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed).
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first.
    /// @return                     True if no errors are thrown.
    function provideFraudBTCFundingProof(
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public returns (bool) {
        self.provideFraudBTCFundingProof(
            _txVersion,
            _txInputVector,
            _txOutputVector,
            _txLocktime,
            _fundingOutputIndex,
            _merkleProof,
            _txIndexInBlock,
            _bitcoinHeaders
        );
        return true;
    }

    /// @notice                     Anyone may notify the deposit of a funding proof to activate the deposit.
    ///                             This is the happy-path of the funding flow. It means that we have succeeded.
    /// @dev                        Takes a pre-parsed transaction and calculates values needed to verify funding.
    /// @param _txVersion           Transaction version number (4-byte LE).
    /// @param _txInputVector       All transaction inputs prepended by the number of inputs encoded as a VarInt, max 0xFC(252) inputs.
    /// @param _txOutputVector      All transaction outputs prepended by the number of outputs encoded as a VarInt, max 0xFC(252) outputs.
    /// @param _txLocktime          Final 4 bytes of the transaction.
    /// @param _fundingOutputIndex  Index of funding output in _txOutputVector (0-indexed).
    /// @param _merkleProof         The merkle proof of transaction inclusion in a block.
    /// @param _txIndexInBlock      Transaction index in the block (0-indexed).
    /// @param _bitcoinHeaders      Single bytestring of 80-byte bitcoin headers, lowest height first.
    /// @return                     True if no errors are thrown.
    function provideBTCFundingProof(
        bytes4 _txVersion,
        bytes memory _txInputVector,
        bytes memory _txOutputVector,
        bytes4 _txLocktime,
        uint8 _fundingOutputIndex,
        bytes memory _merkleProof,
        uint256 _txIndexInBlock,
        bytes memory _bitcoinHeaders
    ) public returns (bool) {
        self.provideBTCFundingProof(
            _txVersion,
            _txInputVector,
            _txOutputVector,
            _txLocktime,
            _fundingOutputIndex,
            _merkleProof,
            _txIndexInBlock,
            _bitcoinHeaders
        );
        return true;
    }

    //
    // FRAUD
    //

    /// @notice                 Anyone can provide a signature that was not requested to prove fraud.
    /// @dev                    Calls out to the keep to verify if there was fraud.
    /// @param  _v              Signature recovery value.
    /// @param  _r              Signature R value.
    /// @param  _s              Signature S value.
    /// @param _signedDigest    The digest signed by the signature vrs tuple.
    /// @param _preimage        The sha256 preimage of the digest.
    /// @return                 True if successful, otherwise revert.
    function provideECDSAFraudProof(
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        bytes32 _signedDigest,
        bytes memory _preimage
    ) public returns (bool) {
        self.provideECDSAFraudProof(_v, _r, _s, _signedDigest, _preimage);
        return true;
    }

    //
    // LIQUIDATION
    //

    /// @notice     Closes an auction and purchases the signer bonds. Payout to buyer, funder, then signers if not fraud.
    /// @dev        For interface, reading auctionValue will give a past value. the current is better.
    /// @return     True if successful, revert otherwise.
    function purchaseSignerBondsAtAuction() public returns (bool) {
        self.purchaseSignerBondsAtAuction();
        return true;
    }

    /// @notice     Notify the contract that the signers are undercollateralized.
    /// @dev        Calls out to the system for oracle info.
    /// @return     True if successful, otherwise revert.
    function notifyCourtesyCall() public returns (bool) {
        self.notifyCourtesyCall();
        return true;
    }

    /// @notice     Goes from courtesy call to active.
    /// @dev        Only callable if collateral is sufficient and the deposit is not expiring.
    /// @return     True if successful, otherwise revert.
    function exitCourtesyCall() public returns (bool) {
        self.exitCourtesyCall();
        return true;
    }

    /// @notice     Notify the contract that the signers are undercollateralized.
    /// @dev        Calls out to the system for oracle info.
    /// @return     True if successful, otherwise revert.
    function notifyUndercollateralizedLiquidation() public returns (bool) {
        self.notifyUndercollateralizedLiquidation();
        return true;
    }

    /// @notice     Notifies the contract that the courtesy period has elapsed.
    /// @dev        This is treated as an abort, rather than fraud.
    /// @return     True if successful, otherwise revert.
    function notifyCourtesyTimeout() public returns (bool) {
        self.notifyCourtesyTimeout();
        return true;
    }
}
