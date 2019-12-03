/**
 * The floodgate mitigates the risk of reorg attacks, by gating the volume of Bitcoin transferred through the system.
 */
contract BitcoinFloodgate {
    uint gatedVolume = 0;

    uint BLOCK_REWARD = 12.5 * 10^8;

    constructor() public {
    }

    /**
     * Releases bitcoin from the system, failing if the level poses a risk
     * @param _amount Amount of bitcoin to be released
     * @param _numConfirmations Number of confirmations associated with this bitcoin transaction
     */
    function release(uint _amount, bytes memory _numConfirmations) public {
        uint minConfirmations = BLOCK_REWARD / (gatedVolume + 1);
        require(_numConfirmations >= minConfirmations, "bitcoin not released: not enough confirmations");
    }
    
    /**
     * Increases the volume of Bitcoin gated for release
     * @param _amount Amount of bitcoin gated for release
     */
    function gate(uint _amount) {
        gatedVolume += _amount;
    }
}