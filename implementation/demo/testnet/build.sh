# 
# Deploys and sets up config for a demo testnet
# 
#!/bin/bash
set -ex



SED="sed"

if [[ "$OSTYPE" == "darwin"* ]]; then
    if ! [ -x "$(command -v gsed)" ]; then
        # See also: https://unix.stackexchange.com/questions/13711/differences-between-sed-on-mac-osx-and-other-standard-sed
        echo 'Error: gsed is not installed.' >&2
        echo 'GNU sed is required due to difference in functionality from OSX/FreeBSD sed.'
        echo 'Install with: brew install gnu-sed'
        exit 1
    fi
    SED="gsed"
fi

if ! [ -x "$(command -v jq)" ]; then
  echo 'Error: jq is not installed.' >&2
  exit 1
fi


# Get the network ID from Truffle
# for later retrieving deployment details from Truffle artifacts
NETWORK_ID=$(curl -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}' http://localhost:8545 | jq -r .result)

cd $GOPATH

# --------------
# 1. Keep tECDSA
# --------------
cd $GOPATH/src/github.com/keep-network/keep-tecdsa/solidity
TECDSA_MIGRATION=$(truffle migrate --reset)

# The sed regex needs to match both:
# - Deploying KeepRegistry
# - Replacing KeepRegistry
# it then matches a range until the next Deploying, or the end of the migration log
# Then it's a simple extraction of the contract address
KEEP_REGISTRY=$(cat build/contracts/KeepRegistry.json | jq -r ".networks.\"$NETWORK_ID\".address")
ECDSA_KEEP_FACTORY=$(cat build/contracts/ECDSAKeepFactory.json | jq -r ".networks.\"$NETWORK_ID\".address")
cd ..
$SED -i -e "/ECDSAKeepFactory = /s/0x[a-fA-F0-9]\{0,40\}/$ECDSA_KEEP_FACTORY/" configs/config.toml


# -------
# 2. tBTC
# -------
cd $GOPATH/src/github.com/keep-network/tbtc/implementation

# Uncomment if you want to redeploy Uniswap (shouldn't be needed)
# cd scripts
# ./deploy_uniswap.sh
# cd ..

TBTC_MIGRATION=$(truffle migrate --reset)
TBTC_SYSTEM=$(cat build/contracts/TBTCSystem.json | jq -r ".networks.\"$NETWORK_ID\".address")


# -------------------
# 3. tBTC maintainers
# -------------------
cd $GOPATH/src/github.com/keep-network/tbtc-maintainers
$SED -i -e "/TBTCSystem = /s/[0x][a-fA-F0-9]\{0,40\}/$TBTC_SYSTEM/" configs/config.toml