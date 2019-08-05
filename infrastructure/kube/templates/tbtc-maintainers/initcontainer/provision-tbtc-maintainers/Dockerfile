FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55
RUN npm install concat-stream@2.0.0
RUN npm install toml@3.0.0
RUN npm install tomlify-j0.4@3.0.0

ARG CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
ENV CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS=$CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

ARG ETH_HOSTNAME
ENV ETH_HOSTNAME=$ETH_HOSTNAME

ARG ETH_RPC_PORT
ENV ETH_RPC_PORT=$ETH_RPC_PORT

ARG ETH_WS_PORT
ENV ETH_WS_PORT=$ETH_WS_PORT

ARG ETH_NETWORK_ID
ENV ETH_NETWORK_ID=$ETH_NETWORK_ID

COPY ./TBTCSystem.json /tmp/TBTCSystem.json

COPY ./tbtc-maintainers-template.toml /tmp/tbtc-maintainers-template.toml

COPY ./provision-tbtc-maintainers.js /tmp/provision-tbtc-maintainers.js

ENTRYPOINT ["node", "./provision-tbtc-maintainers.js"]

