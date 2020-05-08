pragma solidity 0.5.17;

import {DepositLiquidation} from "./DepositLiquidation.sol";
import {DepositUtils} from "./DepositUtils.sol";
import {DepositFunding} from "./DepositFunding.sol";
import {DepositRedemption} from "./DepositRedemption.sol";
import {DepositStates} from "./DepositStates.sol";
import {ITBTCSystem} from "../interfaces/ITBTCSystem.sol";
import {IERC721} from "openzeppelin-solidity/contracts/token/ERC721/IERC721.sol";
import {TBTCToken} from "../system/TBTCToken.sol";
import {FeeRebateToken} from "../system/FeeRebateToken.sol";

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
        require(msg.sender == self.keepAddress, "Deposit contract can only receive ETH from underlying keep");
    }

    /// @notice     Get the keep contract address associated with the Deposit.
    /// @dev        The keep contract address is saved on Deposit initialization.
    /// @return     Address of the Keep contract.
    function getKeepAddress() public view returns (address) {
        return self.keepAddress;
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
    /// @return     uint64 lot size in satoshi.
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
    /// @param _tbtcSystem        `TBTCSystem` contract. More info in `TBTCSystem`.
    /// @param _tbtcToken         `TBTCToken` contract. More info in TBTCToken`.
    /// @param _tbtcDepositToken  `TBTCDepositToken` (TDT) contract. More info in `TBTCDepositToken`.
    /// @param _feeRebateToken    `FeeRebateToken` (FRT) contract. More info in `FeeRebateToken`.
    /// @param _vendingMachineAddress    `VendingMachine` address. More info in `VendingMachine`.
    /// @param _m           Signing group honesty threshold.
    /// @param _n           Signing group size.
    /// @param _lotSizeSatoshis The minimum amount of satoshi the funder is required to send.
    ///                         This is also the amount of TBTC the TDT holder will receive:
    ///                         (10**7 satoshi == 0.1 BTC == 0.1 TBTC).
    /// @return             True if successful, otherwise revert.
    function createNewDeposit(
        ITBTCSystem _tbtcSystem,
        TBTCToken _tbtcToken,
        IERC721 _tbtcDepositToken,
        FeeRebateToken _feeRebateToken,
        address _vendingMachineAddress,
        uint16 _m,
        uint16 _n,
        uint64 _lotSizeSatoshis
    ) public onlyFactory payable returns (bool) {
        self.tbtcSystem = _tbtcSystem;
        self.tbtcToken = _tbtcToken;
        self.tbtcDepositToken = _tbtcDepositToken;
        self.feeRebateToken = _feeRebateToken;
        self.vendingMachineAddress = _vendingMachineAddress;
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

    /// @notice Requests a funder abort for a failed-funding deposit; that is,
    ///         requests the return of a sent UTXO to _abortOutputScript. It
    ///         imposes no requirements on the signing group. Signers should
    ///         send their UTXO to the requested output script, but do so at
    ///         their discretion and with no penalty for failing to do so. This
    ///         can be used for example when a UTXO is sent that is the wrong
    ///         size for the lot.
    /// @dev This is a self-admitted funder fault, and is only be callable by
    ///      the TDT holder. This function emits the FunderAbortRequested event,
    ///      but stores no additional state.
    /// @param _abortOutputScript The output script the funder wishes to request
    ///        a return of their UTXO to.
    function requestFunderAbort(bytes memory _abortOutputScript) public {
        require(
            self.depositOwner() == msg.sender,
            "Only TDT holder can request funder abort"
        );

        self.requestFunderAbort(_abortOutputScript);
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

    /// @notice Get the current collateralization level for this Deposit.
    /// @dev    This value represents the percentage of the backing BTC value the signers
    ///         currently must hold as bond.
    /// @return The current collateralization level for this deposit.
    function getCollateralizationPercentage() public view returns (uint256) {
        return self.getCollateralizationPercentage();
    }

    /// @notice Get the initial collateralization level for this Deposit.
    /// @dev    This value represents the percentage of the backing BTC value the signers hold initially.
    /// @return The initial collateralization level for this deposit.
    function getInitialCollateralizedPercent() public view returns (uint16) {
        return self.initialCollateralizedPercent;
    }

    /// @notice Get the undercollateralization level for this Deposit.
    /// @dev    This collateralization level is semi-critical. If the collateralization level falls
    ///         below this percentage the Deposit can get courtesy-called. This value represents the percentage
    ///         of the backing BTC value the signers must hold as bond in order to not be undercollateralized.
    /// @return The undercollateralized level for this deposit.
    function getUndercollateralizedThresholdPercent() public view returns (uint16) {
        return self.undercollateralizedThresholdPercent;
    }

    /// @notice Get the severe undercollateralization level for this Deposit.
    /// @dev    This collateralization level is critical. If the collateralization level falls
    ///         below this percentage the Deposit can get liquidated. This value represents the percentage
    ///         of the backing BTC value the signers must hold as bond in order to not be severely undercollateralized.
    /// @return The severely undercollateralized level for this deposit.
    function getSeverelyUndercollateralizedThresholdPercent() public view returns (uint16) {
        return self.severelyUndercollateralizedThresholdPercent;
    }

    /// @notice     Calculates the amount of value at auction right now.
    /// @dev        We calculate the % of the auction that has elapsed, then scale the value up.
    /// @return     The value in wei to distribute in the auction at the current time.
    function auctionValue() public view returns (uint256) {
        return self.auctionValue();
    }

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

    /// @notice     Withdraw caller's allowance.
    /// @dev        Withdrawals can only happen when a contract is in an end-state.
    /// @return     True if successful, otherwise revert.
    function withdrawFunds() public returns (bool) {
        self.withdrawFunds();
        return true;
    }

    /// @notice     Get caller's withdraw allowance.
    /// @return     The withdraw allowance in wei.
    function getWithdrawAllowance() public view returns (uint256) {
        return self.getWithdrawAllowance();
    }
}
