#!/bin/bash

set -ex

uniswap() {
    # Clone Uniswap for deployment
    if [ ! -d uniswap ]; then
        git clone https://github.com/keep-network/uniswap
    fi

    # Setup repo
    cd uniswap
    npm i
    export ETH_RPC_URL="http://localhost:8545"
    
    # Run migration
    UNISWAP_DEPLOYMENT=$(npm run migrate)

    # Get address of UniswapFactory
    FACTORY=$(echo "$UNISWAP_DEPLOYMENT" | sed -n /Factory/p | cut -d' ' -f2)

    # Update UniswapFactoryAddress in migration
    sed -i -e "/UniswapFactoryAddress/s/0x[a-fA-F0-9]\{0,40\}/$FACTORY/" ../../migrations/2_deploy_contracts.js
    if [[ -e ../../migrations/2_deploy_contracts.js-e ]]; then
        rm ../../migrations/2_deploy_contracts.js-e
    fi
}

uniswap