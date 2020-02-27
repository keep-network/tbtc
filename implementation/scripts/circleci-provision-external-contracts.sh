#!/bin/bash
set -ex

# BondedECDSAKeepVendorAddress: Migration from keep-network/keep-tecdsa
# BONDED_ECDSA_KEEP_VENDOR_CONTRACT_DATA is set in the CircleCI job config
# ETH_NETWORK_ID is set in the CircleCI context for each deployed environment
BONDED_ECDSA_KEEP_VENDOR_ADDRESS=""

function fetch_bonded_keep_vendor_address() {
  gsutil -q cp gs://${CONTRACT_DATA_BUCKET}/keep-tecdsa/${BONDED_ECDSA_KEEP_VENDOR_CONTRACT_DATA} ./
  BONDED_ECDSA_KEEP_VENDOR_ADDRESS=$(cat ./${BONDED_ECDSA_KEEP_VENDOR_CONTRACT_DATA} | jq ".networks[\"${ETH_NETWORK_ID}\"].address" | tr -d '"')
}

function set_bonded_keep_vendor_address() {
  # TODO: Replace file we store external addresses by a `json` file and use `jq` to update it.
  sed -i -e "/BondedECDSAKeepVendorAddress/s/0x[a-zA-Z0-9]\{0,40\}/${BONDED_ECDSA_KEEP_VENDOR_ADDRESS}/" ./implementation/migrations/externals.js
}

fetch_bonded_keep_vendor_address
set_bonded_keep_vendor_address
