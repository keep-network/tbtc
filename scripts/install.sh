#!/bin/bash
set -euo pipefail

# Dafault inputs.
KEEP_ECDSA_PATH_DEFAULT=$(realpath -m $(dirname $0)/../../keep-ecdsa)
KEEP_ACCOUNT_PASSWORD_DEFAULT="password"
NETWORK_DEFAULT="development"
CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=""

# Run script.
LOG_START='\n\e[1;36m'  # new line + bold + cyan
LOG_END='\n\e[0m'       # new line + reset
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

help()
{
   echo "Usage: $0"\
        "--keep-ecdsa-path <path>"\
        "--account-password <password>"\
        "--private-key <private key>"\
        "--network <network>"
   echo -e "\t--keep-ecdsa-path: Path to the keep-ecdsa project"
   echo -e "\t--account-password: Account password"
   echo -e "\t--private-key: Contract owner's account private key"
   echo -e "\t--network: Connection network"
   exit 1 # Exit script after printing help
}

if [ "$0" == "-help" ]; then
  help
fi

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--keep-ecdsa-path")   set -- "$@" "-d" ;;
    "--account-password")  set -- "$@" "-p" ;;
    "--private-key")       set -- "$@" "-k" ;;
    "--network")           set -- "$@" "-n" ;;
    *)                     set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "d:p:k:n:" opt
do
   case "$opt" in
      d ) keep_ecdsa_path="$OPTARG" ;;
      p ) account_password="$OPTARG" ;;
      k ) private_key="$OPTARG" ;;
      n ) network="$OPTARG" ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

printf "${LOG_START}Starting installation...${LOG_END}"
KEEP_ECDSA_PATH=$(realpath ${keep_ecdsa_path:-$KEEP_ECDSA_PATH_DEFAULT})
TBTC_PATH=$(realpath $(dirname $0)/../)
TBTC_SOL_PATH=$(realpath $TBTC_PATH/solidity)
KEEP_ECDSA_SOL_PATH=$(realpath $KEEP_ECDSA_PATH/solidity)
KEEP_ECDSA_SOL_ARTIFACTS_PATH=$(realpath $KEEP_ECDSA_SOL_PATH/build/contracts)
NETWORK=${network:-$NETWORK_DEFAULT}
ACCOUNT_PRIVATE_KEY=${private_key:-""}

cd $TBTC_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

if [ "$NETWORK" != "alfajores" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ACCOUNT_PASSWORD=$KEEP_ACCOUNT_PASSWORD_DEFAULT \
        npx truffle exec scripts/unlock-eth-accounts.js --network $NETWORK
fi

printf "${LOG_START}Finding current network ID...${LOG_END}"

output=$(CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY npx truffle exec ./scripts/get-network-id.js --network $NETWORK)
NETWORKID=$(echo "$output" | tail -1)

printf "${LOG_START}Fetching external contracts addresses...${LOG_END}"
KEEP_ECDSA_SOL_ARTIFACTS_PATH=$KEEP_ECDSA_SOL_ARTIFACTS_PATH \
NETWORKID=$NETWORKID \
    ./scripts/lcl-provision-external-contracts.sh

printf "${LOG_START}Migrating contracts...${LOG_END}"
npm run clean
CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY \
    npx truffle migrate --reset --network $NETWORK

printf "${DONE_START}Installation completed!${DONE_END}"
