#!/bin/bash

DIR=../solidity

# Source contracts for binding generation.
# They should include the subdirectory prefix if any.
FILES="DepositLog.sol deposit/Deposit.sol"

SOLIDITY_DIR=$DIR SOLIDITY_FILES=$FILES make
