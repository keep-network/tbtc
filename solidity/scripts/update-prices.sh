#!/bin/sh
set -ex
truffle exec ./scripts/get-prices.js --network mainnet | tail -n +3 > ./migrations/prices.json