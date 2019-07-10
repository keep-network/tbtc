const TBTCSystemStub = artifacts.require('TBTCSystemStub')
const TBTC = artifacts.require('TBTC')
const IUniswapFactory = artifacts.require('IUniswapFactory')
const IUniswapExchange = artifacts.require('IUniswapExchange')

const UniswapDeployment = artifacts.require('UniswapDeployment')

import { UniswapHelpers } from './helpers/uniswap'
import expectThrow from './helpers/expectThrow'

// Tests the Uniswap deployment

contract('Uniswap', (accounts) => {
  let tbtcSystem
  let tbtc

  describe('deployment', async () => {
    it('deployed the uniswap factory and exchange', async () => {
      tbtcSystem = await TBTCSystemStub.deployed()

      tbtc = await tbtcSystem.tbtc()
      expect(tbtc).to.not.be.empty

      const uniswapFactoryAddr = await tbtcSystem.uniswapFactory()
      expect(uniswapFactoryAddr).to.not.be.empty

      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)
      const tbtcExchangeAddr = await uniswapFactory.getExchange(tbtc)
      expect(tbtcExchangeAddr).to.not.be.empty
    })

    it('has liquidity by default', async () => {})
  })

  describe('TBTC Uniswap Exchange', () => {
    let tbtc
    let tbtcExchange


    beforeEach(async () => {
      /* eslint-disable no-unused-vars */
      tbtc = await TBTC.new()

      // We rely on the already pre-deployed Uniswap factory here.
      const tbtcSystem = await TBTCSystemStub.deployed()
      const uniswapDeployment = await UniswapDeployment.deployed()
      const uniswapFactoryAddr = await uniswapDeployment.factory()

      const uniswapFactory = await IUniswapFactory.at(uniswapFactoryAddr)

      const res = await uniswapFactory.createExchange(tbtc.address)
      const tbtcExchangeAddr = await uniswapFactory.getExchange.call(tbtc.address)


      tbtcExchange = await IUniswapExchange.at(tbtcExchangeAddr)
      /* eslint-enable no-unused-vars */
    })


    it('has no liquidity by default', async () => {
      await expectThrow(
        tbtcExchange.getTokenToEthInputPrice.call(1)
      )
    })

    describe.only('e2e testing of a trade', () => {
      it('adds liquidity and trades ETH for TBTC', async () => {
        // This avoids rabbit-hole debugging
        // stemming from the fact Vyper is new and they don't do REVERT's
        expect(
          await web3.eth.getBalance(accounts[0])
        ).to.not.eq('0')

        expect(
          await web3.eth.getBalance(accounts[1])
        ).to.not.eq('0')

        // Both tokens use 18 decimal places, so we can use toWei here.
        const TBTC_AMT = web3.utils.toWei('50', 'ether')
        const ETH_AMT = web3.utils.toWei('1', 'ether')

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
        const TBTC_BUY_AMT = web3.utils.toWei('1', 'ether')

        // rough price - we don't think about slippage
        // we are testing that Uniswap works, not testing the exact
        // formulae of the price invariant
        // when they come out with uniswap.js, this code could be made better
        const priceEth = await tbtcExchange.getTokenToEthInputPrice.call(TBTC_BUY_AMT)
        expect(priceEth.toString()).to.eq('90661089388014913')

        const buyer = accounts[1]

        // def ethToTokenSwapInput(min_tokens: uint256, deadline: timestamp) -> uint256:
        await tbtcExchange.ethToTokenSwapInput(
          TBTC_BUY_AMT,
          UniswapHelpers.getDeadline(),
          { value: UniswapHelpers.calcWithFee(priceEth), from: buyer }
        )

        const balance = await tbtc.balanceOf(buyer)
        expect(balance.gt(TBTC_BUY_AMT))
      })
    })
  })
})


/**


Some notes on the cryptoeconomics of Uniswap and TBTC, time arbitrage between chains

 1. Deposit abort flow is started
 2. Bond goes for automatic liquidation via Uniswap
 3. We liquidate the signer bonds `150% * fetchOraclePrice()`. Potential options:
  a. Uniswap 100% of ETH
  b. Uniswap 50% ETH and 50% falling price auction

I think it's worth documenting the additional complexity stemming from automatic liquidation. 3a won't necessarily buy up enough TBTC to burn, depending on the pool size, but the arbitrage incentives are going to keep the price stable to its price-oracle reported value (with a little delta).
*/

/**
 *

source sethenv.sh

# deploy uniswap
export FACTORY=$(seth send --create $(cat uniswap/contracts-vyper/bytecode/factory.txt))
export EXCHANGE_TMPL=$(seth send --create $(cat uniswap/contracts-vyper/bytecode/exchange.txt))
seth send $FACTORY "initializeFactory(address)" $EXCHANGE_TMPL


#export TBTC=$(jq -r '.networks["5777"].address' build/contracts/TBTC.json)
export UNISWAP_DEPLOY=$(jq -r '.networks["5777"].address' build/contracts/UniswapDeployment.json)
export FACTORY=$(seth call $UNISWAP_DEPLOY "factory()(address)")

export EXCHANGE=$(seth call $FACTORY "getExchange(address)(address)" $TBTC)
seth balance $EXCHANGE

export BUYER=$(seth accounts | sed -n 2p | awk '{ print $1 }')
seth balance $BUYER


# mint some tbtc
export TBTC=$(seth send --create $(cat build/contracts/TBTC.json | jq -r .bytecode))
seth send $TBTC "mint(address,uint256)" $BUYER $(seth --to-uint256 10000000000)
seth call $TBTC "balanceOf(address)(uint256)" $BUYER

# create exchange + approve
seth send $FACTORY "createExchange(address)(address)" $TBTC
EXCHANGE=$(seth call $FACTORY "getExchange(address)(address)" $TBTC)
seth send --from $BUYER $TBTC "approve(address,uint)(bool)" $EXCHANGE $(seth --to-uint256 10000000000)
seth call $TBTC "allowance(address,address)(uint256)" $BUYER $EXCHANGE


# add to exchange
# value must be above 1000000000 wei
export DEADLINE=$(bc <<< "$(perl -MTime::HiRes=time -e 'printf "%d\n", time') + 200000")
seth send --from $BUYER --value 100 $EXCHANGE "addLiquidity(uint256,uint256,uint256)(uint256)" 1 "$(seth --to-uint256 100000000)" $DEADLINE

seth call $TBTC "balanceOf(address)(uint256)" $EXCHANGE
seth balance $EXCHANGE

# now estimate a proper gas margin
export ETH_SOLD=$(seth --to-uint256 500000)

export PRICE=$(seth --to-dec $(seth call $EXCHANGE "getTokenToEthInputPrice(uint256)(uint256)" $ETH_SOLD))
seth send --from $ETH_FROM --value $PRICE $EXCHANGE "ethToTokenSwapInput(uint256,uint256)" $BUY_AMT $DEADLINE


 */
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
