#!/bin/bash
set -ex

# KeepRegistryAddress: Migration from keep-network/keep-tecdsa
# KEEP_REGISTRY_CONTRACT_DATA is set in the CircleCI job config
# ETH_NETWORK_ID is set in the CircleCI context for each deployed environment
KEEP_REGISTRY_ADDRESS=""

function fetch_keep_registry_address() {
  gsutil -q cp gs://${CONTRACT_DATA_BUCKET}/keep-tecdsa/${KEEP_REGISTRY_CONTRACT_DATA} ./
  KEEP_REGISTRY_ADDRESS=$(cat ./${KEEP_REGISTRY_CONTRACT_DATA} | jq ".networks[\"${ETH_NETWORK_ID}\"].address" | tr -d '"')
}

function set_keep_registry_address() {
  sed -i -e "/KeepRegistryAddress/s/0x[a-fA-F0-9]\{0,40\}/${KEEP_REGISTRY_ADDRESS}/" ./implementation/migrations/externals.js
}

fetch_keep_registry_address
set_keep_registry_address