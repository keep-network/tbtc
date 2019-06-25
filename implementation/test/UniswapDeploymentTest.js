const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTC = artifacts.require('TBTC')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

const uniswap = require('../uniswap')

const truffleAssert = require('truffle-assertions'); 
const Web3 = require('web3');
const web3 = new Web3()

const UniswapDeployment = artifacts.require('UniswapDeployment')


import { UniswapHelpers } from './helpers/uniswap'
import expectThrow from './helpers/expectThrow'

// Tests the Uniswap deployment

contract('Uniswap', (accounts) => {
    let tbtcSystem;
    let tbtc;

    describe('deployment', async () => {
        it('deployed the uniswap factory and exchange', async () => {
            let tbtcSystem = await TBTCSystemStub.deployed();

            let tbtc = await tbtcSystem.tbtc()
            expect(tbtc).to.not.be.empty;

            let uniswapFactoryAddr = await tbtcSystem.uniswapFactory()
            expect(uniswapFactoryAddr).to.not.be.empty;

            let uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr);
            let tbtcExchangeAddr = await uniswapFactory.getExchange(tbtc)
            expect(tbtcExchangeAddr).to.not.be.empty;
        })

        it('has liquidity by default', async () => {})
    })

    describe('TBTC Uniswap Exchange', () => {
        let tbtc;
        let tbtcExchange;

        beforeEach(async () => {
            tbtc = await TBTC.new()

            // We rely on the already pre-deployed Uniswap factory here.
            let tbtcSystem = await TBTCSystemStub.deployed();
            const uniswapDeployment = await UniswapDeployment.deployed()
            let uniswapFactoryAddr = await uniswapDeployment.factory()

            let uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr);
        
            let res = await uniswapFactory.createExchange(tbtc.address)
            let tbtcExchangeAddr = await uniswapFactory.getExchange.call(tbtc.address)

            tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddr)
        })

        it('has no liquidity by default', async () => {
            await expectThrow(
                tbtcExchange.getTokenToEthInputPrice.call(1)
            );
        }) 

        describe.only('e2e testing of a trade', () => {        
            it('adds liquidity and trades ETH for TBTC', async () => {
                // Both tokens use 18 decimal places, so we can use toWei here.
                const TBTC_AMT = web3.utils.toWei('50', 'ether');
                const ETH_AMT = web3.utils.toWei('1', 'ether');

                // Mint TBTC
                await tbtc.mint(
                    accounts[0],
                    TBTC_AMT
                )
                await tbtc.mint(
                    accounts[1],
                    TBTC_AMT
                )
                
                await tbtc.approve(tbtcExchange.address, TBTC_AMT, { from: accounts[0] })
                await tbtc.approve(tbtcExchange.address, TBTC_AMT, { from: accounts[1] })

                // min_liquidity, max_tokens, deadline
                const TBTC_ADDED = web3.utils.toWei('10', 'ether')
                await tbtcExchange.addLiquidity(
                    '0', 
                    TBTC_ADDED, 
                    UniswapHelpers.getDeadline(),
                    { value: ETH_AMT }
                )

                // it will be at an exchange rate of 
                // 10 TBTC : 1 ETH
                const TBTC_BUY_AMT = web3.utils.toWei('1', 'ether');

                // rough price - we don't think about slippage
                // we are testing that Uniswap works, not testing the exact
                // formulae of the price invariant
                // when they come out with uniswap.js, this code could be made better
                let priceEth = await tbtcExchange.getTokenToEthInputPrice.call(TBTC_BUY_AMT)
                expect(priceEth.toString()).to.eq('90661089388014913')

                const buyer = accounts[1];
                await tbtcExchange.ethToTokenSwapInput(
                    TBTC_BUY_AMT,
                    UniswapHelpers.getDeadline(),
                    { value: '90661089388014914', from: buyer }
                )
                
                let balance = await tbtc.balanceOf(buyer)
                expect(balance.gt(TBTC_BUY_AMT));
            })
        })
    })

})


/**
 * 
 * 
 * 
 * 


export ETH_FROM=$(seth accounts | head -n1 | awk '{ print $1 }')
seth send --create $(cat uniswap/contracts-vyper/bytecode/factory.txt)
# 0xf0d61329351932e6e83517e5d0cb666dacb9589d
seth send --create $(cat uniswap/contracts-vyper/bytecode/exchange.txt)
# 0x61af707a9179e171bfc6d6e27dca18c5ab64c577

seth send 0xf0d61329351932e6e83517e5d0cb666dacb9589d "initializeFactory(address)" 0x61af707a9179e171bfc6d6e27dca18c5ab64c577


seth send --create $(cat build/contracts/TBTC.json | jq -r .bytecode)
export TBTC=0xb902b9899d8652219b8c62bc714020e2a0fa99e6
seth send 0xf0d61329351932e6e83517e5d0cb666dacb9589d "createExchange(address)" 0xb902b9899d8652219b8c62bc714020e2a0fa99e6

seth call 0xf0d61329351932e6e83517e5d0cb666dacb9589d "getExchange(address)" 0xb902b9899d8652219b8c62bc714020e2a0fa99e6 | xargs seth --abi-decode "x()(address)"

# exchange:
export EXCHANGE=0x5823ff557798ab14378c44ebf1fe0f327048de10
export MINT_AMT=10000000

seth send $TBTC "mint(address,uint256)" $ETH_FROM $MINT_AMT
seth call $TBTC "balanceOf(address)" $ETH_FROM
seth send $TBTC "approve(address,uint)" $EXCHANGE $MINT_AMT

seth call $TBTC "allowance(address,address)" $ETH_FROM $EXCHANGE

seth send --from $ETH_FROM --value 402343423423423 $EXCHANGE "addLiquidity(uint256,uint256,uint256)" 0 $MINT_AMT 1000000000000
seth send --from $ETH_FROM --value 402343423423423 $EXCHANGE "addLiquidity(uint256,uint256,uint256)" 0 $MINT_AMT 100000000000000

402343423423423
100000000000000

seth call $EXCHANGE "getTokenToEthInputPrice(uint)(uint)" 1
0000000000000000000000000000000000000000000000000000000026218fd9
0000000000000000000000000000000000000000000000000000000015770098




export BUYER=$(seth accounts | sed -n 2p | awk '{ print $1 }')

# function ethToTokenSwapInput(uint256 min_tokens, uint256 deadline) external payable returns (uint256  tokens_bought);
seth send --from $BUYER --value 402343423423423 $EXCHANGE "ethToTokenSwapInput(uint256,uint256)" 3 1000000000000


some issue with the block timestamp, probs encoded as hex



0x9e9073f8d4e877d88da842a039a3698fb5267723 http://127.0.0.1:8545
0x86841052ae15beb5a1b95148999ef9c640de7463 http://127.0.0.1:8545
0x69a97fe59cfab7b9f8fd2854f99a895f00e1f5c6 http://127.0.0.1:8545
0x1583e8a226c926ded00fa583c30ff8d8a6bc9452 http://127.0.0.1:8545


exchange.methods.addLiquidity(
    min_liquidity, 
    max_tokens, 
    deadline 
)


The ethAmount sent to addLiquidity is the exact amount of ETH that will be deposited into the liquidity reserves. It should be 50% of the total value a liquidity provider wishes to deposit into the reserves. 

value = 10000

the Uniswap smart contracts use ethAmount to determine the amount of ERC20 tokens that must be deposited. This token amount is the remaining 50% of total value a liquidity provider wishes to deposit. 

Since exchange rate can change between when a transaction is signed and when it is executed on Ethereum, max_tokens is used to bound the amount this rate can fluctuate. 

For the first liquidity provider, max_tokens is the exact amount of tokens deposited. 

Liquidity tokens are minted to track the relative proportion of total reserves that each liquidity provider has contributed. 

min_liquidity is used in combination with max_tokens and ethAmount to bound the rate at which liquidity tokens are minted. 

For the first liquidity provider, min_liquidity does not do anything and can be set to 0. 

 */