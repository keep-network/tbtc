#!/bin/bash

set -e

if [[ -z $GOOGLE_PROJECT_NAME || -z $GOOGLE_PROJECT_ID || -z $BUILD_TAG || -z $GOOGLE_REGION || -z $GOOGLE_COMPUTE_ZONE_A || -z $TRUFFLE_NETWORK ]]; then
  echo "one or more required variables are undefined"
  exit 1
fi

UTILITYBOX_IP=$(gcloud compute instances --project $GOOGLE_PROJECT_ID describe $GOOGLE_PROJECT_NAME-utility-box --zone $GOOGLE_COMPUTE_ZONE_A --format json | jq .networkInterfaces[0].networkIP -r)

# Setup ssh environment
gcloud compute config-ssh --project $GOOGLE_PROJECT_ID -q
cat >> ~/.ssh/config << EOF
Host *
  StrictHostKeyChecking no

Host utilitybox
  HostName $UTILITYBOX_IP
  IdentityFile ~/.ssh/google_compute_engine
  ProxyCommand ssh -W %h:%p $GOOGLE_PROJECT_NAME-jumphost.$GOOGLE_COMPUTE_ZONE_A.$GOOGLE_PROJECT_ID
EOF

# Copy migration artifacts over
echo "<<<<<<START Prep Utility Box For Migration START<<<<<<"
echo "ssh utilitybox rm -rf /tmp/$BUILD_TAG"
echo "ssh utilitybox mkdir /tmp/$BUILD_TAG"
echo "scp -r ./implementation utilitybox:/tmp/$BUILD_TAG/"
ssh utilitybox rm -rf /tmp/$BUILD_TAG
ssh utilitybox mkdir /tmp/$BUILD_TAG
scp -r ./implementation utilitybox:/tmp/$BUILD_TAG/
echo ">>>>>>FINISH Prep Utility Box For Migration FINISH>>>>>>"

# Run migration
ssh utilitybox << EOF
  set -e
  echo "<<<<<<START Download Kube Creds START<<<<<<"
  echo "gcloud container clusters get-credentials $GOOGLE_PROJECT_NAME --region $GOOGLE_REGION --internal-ip --project=$GOOGLE_PROJECT_ID"
  gcloud container clusters get-credentials $GOOGLE_PROJECT_NAME --region $GOOGLE_REGION --internal-ip --project=$GOOGLE_PROJECT_ID
  echo ">>>>>>FINISH Download Kube Creds FINISH>>>>>>"

  echo "<<<<<<START Port Forward eth-tx-node START<<<<<<"
  echo "nohup timeout 900 kubectl port-forward svc/eth-tx-node 8545:8545 2>&1 > /dev/null &"
  echo "sleep 10s"
  nohup timeout 900 kubectl port-forward svc/eth-tx-node 8545:8545 2>&1 > /dev/null &
  sleep 10s
  echo ">>>>>>FINISH Port Forward eth-tx-node FINISH>>>>>>"

  echo "<<<<<<START Unlock Contract Owner ETH Account START<<<<<<"
  echo "geth --exec \"personal.unlockAccount(\"${CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS}\", \"${CONTRACT_OWNER_ETH_ACCOUNT_PASSWORD}\", 900)\" attach http://localhost:8545"
  geth --exec "personal.unlockAccount(\"${CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS}\", \"${CONTRACT_OWNER_ETH_ACCOUNT_PASSWORD}\", 900)" attach http://localhost:8545
  echo ">>>>>>FINISH Unlock Contract Owner ETH Account FINISH>>>>>>"

  echo "<<<<<<START Contract Migration START<<<<<<"
  cd /tmp/$BUILD_TAG/implementation

# TODO: Migrations fail with truffle version specified in package.json file. That's why we install dependencies manually here, bug: https://github.com/keep-network/tbtc/issues/231
npm install @keep-network/keep-ecdsa@0.1.1
npm install git+https://github.com/summa-tx/bitcoin-spv.git#v1.1.0-high-err
npm install bn-chai@1.0.1
npm install bn.js@4.11.8
npm install chai@4.2.0
npm install create-hash@1.2.0
npm install openzeppelin-solidity@2.3.0
npm install solc@0.5.10
npm install babel-polyfill@6.26.0
npm install babel-preset-es2015@6.18.0
npm install babel-preset-stage-2@6.18.0
npm install babel-preset-stage-3@6.17.0
npm install babel-register@6.26.0
npm install eth-gas-reporter@0.1.12
npm install ganache-cli@6.4.3
npm install truffle@5.0.7

  ./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
  echo ">>>>>>FINISH Contract Migration FINISH>>>>>>"
EOF

echo "<<<<<<START Contract Copy START<<<<<<"
echo "scp utilitybox:/tmp/$BUILD_TAG/implementation/build/contracts/* /tmp/tbtc/contracts"
scp utilitybox:/tmp/$BUILD_TAG/implementation/build/contracts/* /tmp/tbtc/contracts
echo ">>>>>>FINISH Contract Copy>>>>>>"

echo "<<<<<<START Migration Dir Cleanup START<<<<<<"
echo "ssh utilitybox rm -rf /tmp/$BUILD_TAG"
ssh utilitybox rm -rf /tmp/$BUILD_TAG
echo ">>>>>>FINISH Migration Dir Cleanup FINISH>>>>>>"
