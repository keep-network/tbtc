#!/bin/bash
set -ex

# Connect to the local node
ETH_RPC_URL=http://127.0.0.1:8545

# Use the unlocked RPC accounts
ETH_RPC_ACCOUNTS=yes

# Get the 1st account and use it
ETH_FROM=$(seth accounts | head -n1 | awk '{ print $1 }')

# 1/ Deploy Exchange
EXCHANGE_ADDR=$(seth send --create $(cat ./contracts-vyper/bytecode/exchange.txt) --status 2>/dev/null)

# 2/ Deploy Factory
FACTORY_ADDR=$(seth send --create $(cat ./contracts-vyper/bytecode/factory.txt) --status 2>/dev/null)

# 3/ 
seth send $FACTORY_ADDR "initializeFactory(address)" $EXCHANGE_ADDR


echo "$EXCHANGE_ADDR" > deployments/Exchange
echo "$FACTORY_ADDR" > deployments/Factory