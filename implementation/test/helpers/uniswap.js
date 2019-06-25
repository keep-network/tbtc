const BN = require('bn.js')

export class UniswapHelpers {
    // see https://github.com/Uniswap/contracts-vyper/issues/27#issuecomment-476148467
    static GAS_ESTIMATE = {
        
    }

    static getDeadline() {
        const DEADLINE_FROM_NOW = 300    // TX expires in 300 seconds (5 minutes)  
        const deadline = Math.ceil(Date.now() / 1000) + DEADLINE_FROM_NOW
        return deadline
    }

    // Adds 3% Uniswap fee into amount
    static calcWithFee(amt_) {
        const amt = new BN(amt_)
        return amt.add(
            amt.mul(new BN(1003)).div(new BN(1000))
        )
    }
}