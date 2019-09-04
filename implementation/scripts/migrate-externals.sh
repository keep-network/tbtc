#!/bin/bash

set -ex

uniswap() {
    if [ ! -d uniswap ]; then
        git clone https://github.com/keep-network/uniswap
    fi
    cd uniswap
    npm i
    export ETH_RPC_URL="http://localhost:8545"
    UNISWAP_DEPLOYMENT=$(npm run migrate)
    ADDR=$(echo "$UNISWAP_DEPLOYMENT" | sed -n /Factory/p | cut -d' ' -f2)
    sed -i '' -e "/UniswapFactoryAddress/s/0x[a-fA-F0-9]\{0,40\}/$ADDR/" ../../migrations/2_deploy_contracts.js
}

uniswap