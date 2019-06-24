export class UniswapHelpers {
    static getDeadline() {
        const DEADLINE_FROM_NOW = 300    // TX expires in 300 seconds (5 minutes)  
        const deadline = Math.ceil(Date.now() / 1000) + DEADLINE_FROM_NOW
        return deadline
    }
}