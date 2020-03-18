#!/bin/bash
set -ex

# Dafault inputs.
KEEP_ETHEREUM_PASSWORD_DEFAULT="password"
KEEP_ECDSA_PATH_DEFAULT=$(realpath -m $(dirname $0)/../../../keep-tecdsa)

# Run script.
LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

read -p "Enter path to the keep-ecdsa project [$KEEP_ECDSA_PATH_DEFAULT]: " keep_ecdsa_path
KEEP_ECDSA_PATH=$(realpath ${keep_ecdsa_path:-$KEEP_ECDSA_PATH_DEFAULT})

printf "${LOG_START}Starting installation...${LOG_END}"
KEEP_ECDSA_SOL_ARTIFACTS_PATH=$(realpath $KEEP_ECDSA_PATH/solidity/build/contracts)

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

printf "${LOG_START}Finding current ethereum network ID...${LOG_END}"

output=$(truffle exec ./scripts/get-network-id.js --network development)
NETWORKID=$(echo "$output" | tail -1)
printf "Current network ID: ${NETWORKID}\n"

printf "${LOG_START}Fetching external contracts addresses...${LOG_END}"
KEEP_ECDSA_SOL_ARTIFACTS_PATH=$KEEP_ECDSA_SOL_ARTIFACTS_PATH \
NETWORKID=$NETWORKID \
    ./scripts/lcl-provision-external-contracts.sh

printf "${LOG_START}Migrating contracts...${LOG_END}"
truffle migrate --reset --network development