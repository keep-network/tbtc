#!/bin/bash
set -e pipefail

# Run script.
LOG_START='\n\e[1;36m'  # new line + bold + cyan
LOG_END='\n\e[0m'       # new line + reset
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

printf "${LOG_START}Starting installation...${LOG_END}"
TBTC_PATH=$(realpath $(dirname $0)/../)
TBTC_SOL_PATH=$(realpath $TBTC_PATH/solidity)

cd $TBTC_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm ci
npm link @keep-network/keep-ecdsa

printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
    npx truffle exec scripts/unlock-eth-accounts.js --network development

printf "${LOG_START}Migrating contracts...${LOG_END}"
npm run clean
npx truffle migrate --reset --network development

printf "${LOG_START}Creating links...${LOG_END}"
ln -sf build/contracts artifacts
npm link

printf "${DONE_START}Installation completed!${DONE_END}"
