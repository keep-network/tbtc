#!/bin/bash

set -euo pipefail

LOG_START='\n\e[1;36m'  # new line + bold + cyan
LOG_END='\n\e[0m'       # new line + reset
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

RELAY_PATH=$(realpath $(dirname $0)/../)
RELAY_SOL_PATH=$(realpath $RELAY_PATH/solidity)

# Defaults, can be overwritten by env variables/input parameters
BTC_NETWORK=${BTC_NETWORK:-"testnet"}
ETH_NETWORK=${ETH_NETWORK:-"local"}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"
   echo -e "\nEnvironment variables:"
   echo -e "\tBTC_NETWORK: Which BTC network should be used." \
           "Default value is 'testnet'."
   echo -e "\tETH_NETWORK: Which ETH network should be used." \
           "Default value is 'local'."
   exit 1 # Exit script after printing help
}

while getopts "h" opt
do
   case "$opt" in
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done

printf "${LOG_START}BTC network: $BTC_NETWORK ${LOG_END}"
printf "${LOG_START}ETH network: $ETH_NETWORK ${LOG_END}"

printf "${LOG_START}Starting installation...${LOG_END}"

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
cd $RELAY_SOL_PATH
npm install

printf "${LOG_START}Migrating contracts...${LOG_END}"
BTC_NETWORK=$BTC_NETWORK npx truffle migrate --reset --network $ETH_NETWORK

printf "${LOG_START}Building relay node client...${LOG_END}"
cd $RELAY_PATH
#go generate ./...
go build -a -o relay .

# TODO: set contract address in client config

printf "${DONE_START}Installation completed!${DONE_END}"