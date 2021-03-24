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
CONFIG_DIR_PATH_DEFAULT="$RELAY_PATH/config"
CONFIG_DIR_PATH=$(realpath "${CONFIG_DIR_PATH:-$CONFIG_DIR_PATH_DEFAULT}")

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"
   echo -e "\nEnvironment variables:"
   echo -e "\tBTC_NETWORK: Which BTC network should be used." \
           "Default value is 'testnet'."
   echo -e "\tETH_NETWORK: Which ETH network should be used." \
           "Default value is 'local'."
   echo -e "\tCONFIG_DIR_PATH: Location of relay config file(s)." \
           "Default value is 'config' dir placed under project root."
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
go generate ./...
go build -a -o relay .

printf "${LOG_START}Updating relay config files...${LOG_END}"
cd $RELAY_SOL_PATH
for CONFIG_FILE in $CONFIG_DIR_PATH/*.toml
do
  BTC_NETWORK=$BTC_NETWORK \
    CONFIG_FILE_PATH=$CONFIG_FILE \
    npx truffle exec scripts/lcl-client-config.js --network $ETH_NETWORK
done

printf "${DONE_START}Installation completed!${DONE_END}"