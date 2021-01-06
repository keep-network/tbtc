#!/bin/bash

set -e

if [[ -z $TRUFFLE_NETWORK ]]; then
  echo "one or more required variables are undefined"
  exit 1
fi

mkdir -p /tmp/tbtc/contracts
cd ./solidity
npm ci

echo "<<<<<<START Contract Migration START<<<<<<"
./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
cp ./build/contracts/* /tmp/tbtc/contracts
echo ">>>>>>FINISH Contract Migration FINISH>>>>>>"
