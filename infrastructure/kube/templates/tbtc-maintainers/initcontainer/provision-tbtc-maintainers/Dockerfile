FROM node:11 AS runtime

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm install

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

