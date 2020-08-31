#!/bin/bash

set -e

if [[ -z $GOOGLE_PROJECT_NAME || -z $BUILD_TAG || -z $TRUFFLE_NETWORK || -z $TENDERLY_TOKEN || -z $ETH_NETWORK_ID ]]; then
  echo "one or more required variables are undefined"
  exit 1
fi

echo "<<<<<<START Tenderly Intallation START<<<<<<"
curl https://raw.githubusercontent.com/Tenderly/tenderly-cli/master/scripts/install-linux.sh | sh
echo "<<<<<<FINISH Tenderly Installation FINISH<<<<<<"

mkdir -p /tmp/tbtc/contracts
npm ci

echo "<<<<<<START Contract Migration START<<<<<<"
./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
cp ./build/contracts/* /tmp/tbtc/contracts
echo ">>>>>>FINISH Contract Migration FINISH>>>>>>"

# We are not deploying TBTCConstants to Ropsten/Dev Environment. For those we
# use TBTCDevelopmentConstants. We remove TBTCConstants before we push to
# Tenderly because Tenderly push fails on non-deployed contracts.
echo "<<<<<<START Tenderly Push Preparation START<<<<<<"
rm ./build/contracts/TBTCConstants.json
echo ">>>>>>FINISH Tenderly Push Preparation FINISH>>>>>>"

echo "<<<<<<START Tenderly Push START<<<<<<"
tenderly login --authentication-method access-key --access-key $TENDERLY_TOKEN
tenderly push --networks $ETH_NETWORK_ID --tag tbtc \ 
  --tag $GOOGLE_PROJECT_NAME --tag $BUILD_TAG || echo "tendery push failed :("
echo "<<<<<<FINISH Tenderly Push FINISH<<<<<<"