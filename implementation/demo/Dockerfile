# Mount 3 volumes
# /Users/liamz/go/src/github.com/keep-network/keep-tecdsa/solidity
# /Users/liamz/go/src/github.com/keep-network/tbtc/implementation
# /Users/liamz/go/src/github.com/keep-network/tbtc-maintainers
# /Users/liamz/go/src/github.com/keep-network/keep-tecdsa

/Users/liamz/go/src/github.com/keep-network/keep-tecdsa/solidity
truffle migrate --reset

   Deploying 'KeepRegistry'
   ------------------------
   > transaction hash:    0xefa7d0a62decab2ae18ac08d08ee8a65a9f55ace29d7af97d5ae3c47e6f7913e
   > Blocks: 0            Seconds: 0
   > contract address:    0xd60c4A3d83BA92db16C35928cA9832f6F7b596fD
   > block number:        6
   > block timestamp:     1569931256
   > account:             0x563163755f7cd7F826ae1f62340D58D43F76AcB0
   > balance:             99.9300117
   > gas used:            671315
   > gas price:           20 gwei
   > value sent:          0 ETH
   > total cost:          0.0134263 ETH
KeepRegistry
0xd60c4A3d83BA92db16C35928cA9832f6F7b596fD

/Users/liamz/go/src/github.com/keep-network/tbtc/implementation
cd scripts
./deploy_uniswap.sh

cd ..
truffle migrate --reset

Deploying 'TBTCSystem'
   ----------------------
   > transaction hash:    0xe8ffffb7ad8926b597063827fcf1675c8ffc8b591c6eccad467cbdfdb6a59e9e
   > Blocks: 0            Seconds: 0
   > contract address:    0x8a4B9c09b15490fcaE3f1185Ae6aeB8f6bF3c702
TBTCSystem
0x3d545b44881881146BC31bfdE609bA5bDC414304

/Users/liamz/go/src/github.com/keep-network/tbtc-maintainers
code configs/config.toml




truffle exec demo/1_create_new_deposit.js

createNewDeposit transaction:  0x5c476aee344b6b5046b877cb696066268c811ebb4aeb5d04fedd54c06bd10006
new deposit deployed:  0x37539440173A37c4B6Bf56034E7718330afA6054


truffle exec demo/2_fetch_public_key.js 0x37539440173A37c4B6Bf56034E7718330afA6054






keep-tecdsa only tracks keeps created, not previously created
needs sender funds


Call getPublicKey for deposit [0x65041aDC9Cac6190EE3D3224084A09FECc7C265c]
retrieveSignerPubkey transaction:  0xc2acb3a1773ff47187a8a4635e260ac41961ed010700bd4663ff948bf1b73cde
Registered public key:
X: 0xabe7af4a777471598121e029dfe969d5197476655f212b99af4d3f1f706babb1
Y: 0xc727ebb6198875ee5d608b20055c32f52f12697e2e779f032e20248a73a3bdb2



truffle exec demo/3_provide_funding_proof.js 0xc2acb3a1773ff47187a8a4635e260ac41961ed010700bd4663ff948bf1b73cde 1


Using network 'development'.

Call getPublicKey for deposit [0x65041aDC9Cac6190EE3D3224084A09FECc7C265c]
retrieveSignerPubkey transaction:  0xc2acb3a1773ff47187a8a4635e260ac41961ed010700bd4663ff948bf1b73cde
Registered public key:
X: 0xabe7af4a777471598121e029dfe969d5197476655f212b99af4d3f1f706babb1
Y: 0xc727ebb6198875ee5d608b20055c32f52f12697e2e779f032e20248a73a3bdb2

