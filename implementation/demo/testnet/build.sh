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

if ! [ -x "$(command -v seth)" ]; then
  echo 'Error: seth is not installed.' >&2
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

# a. Migrate
truffle migrate --reset

# b. Get addresses
KEEP_REGISTRY=$(cat build/contracts/KeepRegistry.json | jq -r ".networks.\"$NETWORK_ID\".address")
ECDSA_KEEP_FACTORY=$(cat build/contracts/ECDSAKeepFactory.json | jq -r ".networks.\"$NETWORK_ID\".address")

# c. Update config
cd ..
$SED -i -e "/ECDSAKeepFactory = /s/0x[a-fA-F0-9]\{0,40\}/$ECDSA_KEEP_FACTORY/" configs/config.toml


# -------
# 2. tBTC
# -------
cd $GOPATH/src/github.com/keep-network/tbtc/implementation

# a. Deploy/link external systems
cd scripts
./deploy_uniswap.sh
cd ..

$SED -i -e "/KeepRegistryAddress/s/0x[a-fA-F0-9]\{0,40\}/$KEEP_REGISTRY/" migrations/externals.js

# b. Migrate
truffle migrate --reset

# c. Get addresses
TBTC_SYSTEM=$(cat build/contracts/TBTCSystem.json | jq -r ".networks.\"$NETWORK_ID\".address")


# -------------------
# 3. tBTC maintainers
# -------------------
cd $GOPATH/src/github.com/keep-network/tbtc-maintainers

# a. Update config
$SED -i -e "/TBTCSystem = /s/[0x][a-fA-F0-9]\{0,40\}/$TBTC_SYSTEM/" configs/config.toml



# --------
# 4. Dapp
# --------

# a. Copy over the new contract artifacts, with their deployment info.
cd $GOPATH/src/github.com/keep-network/tbtc-dapp/client/src/eth
./copy-artifacts.sh



# ----------------
# 5. Fund accounts
# ----------------

# Now fund accounts. We have to fund:
# 1) keep tECDSA client - so it can register as a member, and submit public key
# 2) Metamask account for testing

# Setup seth, a CLI tool we can send ETH with
export ETH_RPC_URL=http://127.0.0.1:8545
export ETH_RPC_ACCOUNTS=yes
export ETH_FROM=$(seth rpc eth_coinbase)

# a. keep tECDSA client
seth send --value $(seth --to-wei 5 ether) 0xfb3106cc5af24a13d013db4a3efe711c11a0ccd1
# b. Metamask account
seth send --value $(seth --to-wei 5 ether) 0xa7b224751E9F023B9315726C99cA4AC4fb174dAE