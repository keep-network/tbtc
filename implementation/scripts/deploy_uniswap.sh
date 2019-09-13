#!/bin/bash

set -ex

SED="sed"

if [[ "$OSTYPE" == "darwin"* ]]; then
    if ! [ -x "$(command -v gsed)" ]; then
        # See also: https://unix.stackexchange.com/questions/13711/differences-between-sed-on-mac-osx-and-other-standard-sed
        echo 'Error: gsed is not installed.' >&2
        echo 'GNU sed is required due to difference in functionality from OSX/FreeBSD sed.'
        echo 'Install with: brew install gnu-sed'
        exit 1
    fi
    SED="gsed"
fi


# Clone Uniswap for deployment
if [ ! -d uniswap ]; then
    git clone https://github.com/keep-network/uniswap
fi

# Setup repo
cd uniswap
git checkout v0.1.0
npm i
export ETH_RPC_URL="http://localhost:8545"

# Run migration
UNISWAP_DEPLOYMENT=$(npm run migrate)

# Get address of UniswapFactory
FACTORY=$(echo "$UNISWAP_DEPLOYMENT" | sed -n /Factory/p | cut -d' ' -f2)

# Update UniswapFactoryAddress in migration
$SED -i -e "/UniswapFactoryAddress/s/0x[a-fA-F0-9]\{0,40\}/$FACTORY/" ../../migrations/3_initialize.js

# Clean up
rm -rf uniswap