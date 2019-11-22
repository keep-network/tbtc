contract VendingMachine {
    function isQualified(deposit, proof) {
        // check deposit qualified	
		totalValue = 0
        unqualifiedDeposits = [...]
        for deposit in unqualifiedDeposits:
            totalValue += deposit.lotSize
		
		blockReward = 12.5
		minBlocksX = totalValue / blockReward

		proof = verifyProof(proof)
        
        // qualifier = 6 + n
        // where 6 is the minimum number of confs to mint tbtc and
        //       n is the security margin that fluctuates according to opened deposit volume
		require(proof.blocks >= Math.min(6, minBlocksX));
    }
}