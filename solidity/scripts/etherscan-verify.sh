#!/bin/bash

set -e

echo "Verifying contracts on Etherscan..."

npx truffle run verify \
    BTCUtils \
    BytesLib \
    CheckBitcoinSigs \
    Deposit \
    DepositFactory \
    DepositFunding \
    DepositLiquidation \
    DepositRedemption \
    DepositStates \
    DepositUtils \
    ETHBTCPriceFeedMock \
    FeeRebateToken \
    FundingScript \
    KeepFactorySelection \
    Migrations \
    OutsourceDepositLogging \
    RedemptionScript \
    SatWeiPriceFeed \
    TBTCConstants \
    TBTCDepositToken \
    TBTCSystem \
    TBTCToken \
    ValidateSPV \
    VendingMachine \
    --network $TRUFFLE_NETWORK
