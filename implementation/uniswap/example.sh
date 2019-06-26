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
export DEADLINE=$(bc <<< "$(perl -MTime::HiRes=time -e 'printf "%d\n", time') + 36000")

seth send --from $BUYER --value 402343423423423 $EXCHANGE "addLiquidity(uint256,uint256,uint256)(uint256)" 0 "$(seth --to-uint256 100000000)" $DEADLINE
