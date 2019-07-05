#!/bin/sh
set -e

echo '--- CREATE NEW DEPOSIT'

output=$(truffle exec 1_create_new_deposit.js)
echo "$output"
depositAddress=$(echo "$output" | awk '/new deposit deployed:/ { print $4 }')

echo '--- FETCH PUBLIC KEY'

truffle exec 2_fetch_public_key.js $depositAddress

echo '--- SUBMIT FUNDING PROOF'
echo "Enter funding transaction ID and press [ENTER]: "
read txID

truffle exec 3_provide_funding_proof.js $depositAddress $txID 6
