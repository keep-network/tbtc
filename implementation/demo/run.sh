#!/bin/sh
set -e

echo '--- CREATE NEW DEPOSIT'
truffle exec 1_create_new_deposit.js

echo '--- FETCH PUBLIC KEY'

truffle exec 2_fetch_public_key.js

echo '--- SUBMIT FUNDING PROOF'
echo "Enter funding transaction ID and press [ENTER]: "
read txID

truffle exec 3_provide_funding_proof.js $txID 6 0
