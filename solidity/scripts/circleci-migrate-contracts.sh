#!/bin/bash

set -e

if [[ -z $GOOGLE_PROJECT_NAME || -z $GOOGLE_PROJECT_ID || -z $BUILD_TAG || -z $GOOGLE_REGION || -z $GOOGLE_COMPUTE_ZONE_A || -z $TRUFFLE_NETWORK || -z $TENDERLY_TOKEN || -z $ETH_NETWORK_ID ]]; then
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
  echo "nohup kubectl port-forward svc/eth-tx-node 8545:8545 2>&1 > /dev/null &"
  echo "sleep 10s"
  nohup kubectl port-forward svc/eth-tx-node 8545:8545 2>&1 > /dev/null &
  sleep 10s
  echo ">>>>>>FINISH Port Forward eth-tx-node FINISH>>>>>>"

    echo "<<<<<<START Setting Contract Owner Key START<<<<<<"
  echo "export CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=\$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY"
  export CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY
  echo ">>>>>>FINISH Setting Contract Owner Key FINISH>>>>>>"

  echo "<<<<<<START Contract Migration START<<<<<<"
  cd /tmp/$BUILD_TAG/implementation
  npm ci
  ./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
  echo ">>>>>>FINISH Contract Migration FINISH>>>>>>"
  echo "<<<<<<START Tenderly Push START<<<<<<"
  tenderly login --authentication-method token --token $TENDERLY_TOKEN
  tenderly push --networks $ETH_NETWORK_ID --tag tbtc --tag $GOOGLE_PROJECT_NAME --tag $BUILD_TAG
  echo "<<<<<<FINISH Tenderly Push FINISH<<<<<<"
EOF

echo "<<<<<<START Contract Copy START<<<<<<"
echo "scp utilitybox:/tmp/$BUILD_TAG/implementation/build/contracts/* /tmp/tbtc/contracts"
scp utilitybox:/tmp/$BUILD_TAG/implementation/build/contracts/* /tmp/tbtc/contracts
echo ">>>>>>FINISH Contract Copy>>>>>>"

echo "<<<<<<START Migration Dir Cleanup START<<<<<<"
echo "ssh utilitybox rm -rf /tmp/$BUILD_TAG"
ssh utilitybox rm -rf /tmp/$BUILD_TAG
echo ">>>>>>FINISH Migration Dir Cleanup FINISH>>>>>>"
