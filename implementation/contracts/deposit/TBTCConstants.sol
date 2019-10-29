pragma solidity ^0.5.10;

library TBTCConstants {

    // This is intended to make it easy to update system params
    // During testing swap this out with another constats contract

    // System Parameters
    uint256 public constant LOT_SIZE = 1000; // satoshi TODO: decreased for testing, set it back original value `10 ** 8`
    uint256 public constant SATOSHI_MULTIPLIER = 10 ** 10; // multiplier to convert satoshi to TBTC token units
    uint256 public constant SIGNER_FEE_DIVISOR = 200; // 1/200 == 50bps == 0.5% == 0.005
    uint256 public constant BENEFICIARY_FEE_DIVISOR = 1000;  // 1/1000 = 10 bps = 0.1% = 0.001
    uint256 public constant FUNDING_FRAUD_PARTIAL_SLASH_DIVISOR = 2;  // 1/2 = 5000bps = 50% = 0.5
    uint256 public constant DEPOSIT_TERM_LENGTH = 180 * 24 * 60 * 60; // 180 days in seconds
    uint256 public constant UNDERCOLLATERALIZED_THRESHOLD_PERCENT = 140;  // percent
    uint256 public constant SEVERELY_UNDERCOLLATERALIZED_THRESHOLD_PERCENT = 120; // percent
    uint256 public constant TX_PROOF_DIFFICULTY_FACTOR = 1; // TODO: decreased for testing, original value: `6`

    // Redemption Flow
    uint256 public constant REDEMPTION_SIGNATURE_TIMEOUT = 2 * 60 * 60;  // seconds
    uint256 public constant INCREASE_FEE_TIMER = 4 * 60 * 60;  // seconds
    uint256 public constant REDEMPTION_PROOF_TIMEOUT = 6 * 60 * 60;  // seconds
    uint256 public constant MINIMUM_REDEMPTION_FEE = 150; // satoshi TODO: decreased for testing, set it back original value `2000`

    // Funding Flow
    uint256 public constant FUNDING_PROOF_TIMEOUT = 3 * 60 * 60; // seconds
    uint256 public constant FORMATION_TIMEOUT = 3 * 60 * 60; // seconds
    uint256 public constant FRAUD_FUNDING_PROOF_TIMEOUT = 3 * 60 * 60; // seconds
    uint256 public constant FUNDER_BOND_AMOUNT_WEI = 10 ** 5; // wei

    // Liquidation Flow
    uint256 public constant COURTESY_CALL_DURATION = 6 * 60 * 60; // seconds
    uint256 public constant AUCTION_DURATION = 24 * 60 * 60; // seconds
    uint256 public constant AUCTION_BASE_PERCENTAGE = 90; // percents
    uint256 public constant PERMITTED_FEE_BUMPS = 5; // number of times the fee can be increased

    // Getters for easy access
    function getLotSize() public pure returns (uint256) { return LOT_SIZE; }
    function getSatoshiMultiplier() public pure returns (uint256) { return SATOSHI_MULTIPLIER; }
    function getSignerFeeDivisor() public pure returns (uint256) { return SIGNER_FEE_DIVISOR; }
    function getBeneficiaryRewardDivisor() public pure returns (uint256) { return BENEFICIARY_FEE_DIVISOR; }
    function getFundingFraudPartialSlashDivisor() public pure returns (uint256) { return FUNDING_FRAUD_PARTIAL_SLASH_DIVISOR; }
    function getDepositTerm() public pure returns (uint256) { return DEPOSIT_TERM_LENGTH; }
    function getUndercollateralizedPercent() public pure returns (uint256) { return UNDERCOLLATERALIZED_THRESHOLD_PERCENT; }
    function getSeverelyUndercollateralizedPercent() public pure returns (uint256) { return SEVERELY_UNDERCOLLATERALIZED_THRESHOLD_PERCENT; }
    function getTxProofDifficultyFactor() public pure returns (uint256) { return TX_PROOF_DIFFICULTY_FACTOR; }

    function getSignatureTimeout() public pure returns (uint256) { return REDEMPTION_SIGNATURE_TIMEOUT; }
    function getIncreaseFeeTimer() public pure returns (uint256) { return INCREASE_FEE_TIMER; }
    function getRedepmtionProofTimeout() public pure returns (uint256) { return REDEMPTION_PROOF_TIMEOUT; }
    function getMinimumRedemptionFee() public pure returns (uint256) { return MINIMUM_REDEMPTION_FEE; }

    function getFundingTimeout() public pure returns (uint256) { return FUNDING_PROOF_TIMEOUT; }
    function getSigningGroupFormationTimeout() public pure returns (uint256) { return FORMATION_TIMEOUT; }
    function getFraudFundingTimeout() public pure returns (uint256) { return FRAUD_FUNDING_PROOF_TIMEOUT; }
    function getFunderBondAmount() public pure returns (uint256) { return FUNDER_BOND_AMOUNT_WEI; }

    function getCourtesyCallTimeout() public pure returns (uint256) { return COURTESY_CALL_DURATION; }
    function getAuctionDuration() public pure returns (uint256) { return AUCTION_DURATION; }
    function getAuctionBasePercentage() public pure returns (uint256) { return AUCTION_BASE_PERCENTAGE; }
    function getPermittedFeeBumps() public pure returns (uint256) {return PERMITTED_FEE_BUMPS; }
}
