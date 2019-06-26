#!/bin/sh

echo '--- CREATE NEW DEPOSIT'
truffle exec 1_create_new_deposit.js

echo '--- FETCH PUBLIC KEY'
truffle exec 2_get_public_key.js

echo '--- SUBMIT FUNDING PROOF'
truffle exec 3_provide_funding_proof.js 4e3388fad9732d3c099fa95a98c9d2dd062e4a8cf7ad03f647eaccc478a79068 6 0
