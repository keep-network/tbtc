#!/bin/bash
set -e

# Fetch addresses of contacts migrated from keep-network/keep-ecdsa project.
# The `keep-ecdsa` contracts have to be migrated before running this script.
# It requires `KEEP_ECDSA_SOL_ARTIFACTS_PATH` variable to pointing to a directory where
# contracts artifacts after migrations are located. It also expects NETWORK_ID
# variable to be set to the ID of the network where contract were deployed.
# 
# Sample command:
# KEEP_ECDSA_SOL_ARTIFACTS_PATH=~/go/src/github.com/keep-network/keep-ecdsa/solidity/build/contracts \
# NETWORK_ID=1801 \
#   ./lcl-provision-external-contracts.sh

FACTORY_CONTRACT_DATA="BondedECDSAKeepFactory.json"
FACTORY_PROPERTY="BondedECDSAKeepFactory"

DESTINATION_FILE=$(realpath $(dirname $0)/../migrations/externals.js)

ADDRESS_REGEXP=^0[xX][0-9a-fA-F]{40}$

# Query to get address of the deployed contract for the first network on the list.
JSON_QUERY=".networks.\"${NETWORKID}\".address"

SED_SUBSTITUTION_REGEXP="['\"][a-zA-Z0-9]*['\"]"

FAILED=false

function fetch_bonded_ecdsa_keep_factory_contract_address() {
  echo "Fetching value for ${FACTORY_PROPERTY}..."
  local contractDataPath=$(realpath $KEEP_ECDSA_SOL_ARTIFACTS_PATH/$FACTORY_CONTRACT_DATA)
  local ADDRESS=$(cat ${contractDataPath} | jq "${JSON_QUERY}" | tr -d '"')

  if [[ !($ADDRESS =~ $ADDRESS_REGEXP) ]]; then
    echo "Invalid address: ${ADDRESS}"
    FAILED=true
  else
    echo "Found value for ${FACTORY_PROPERTY} = ${ADDRESS}"
    sed -i -e "/${FACTORY_PROPERTY}/s/${SED_SUBSTITUTION_REGEXP}/\"${ADDRESS}\"/" $DESTINATION_FILE
  fi
}

fetch_bonded_ecdsa_keep_factory_contract_address

if $FAILED; then
echo "Failed to fetch external contract addresses!"
  exit 1
fi
