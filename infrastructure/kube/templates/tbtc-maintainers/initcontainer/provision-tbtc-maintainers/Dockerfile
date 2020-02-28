FROM node:11 AS runtime

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm install

COPY ./TBTCSystem.json /tmp/TBTCSystem.json

COPY ./tbtc-maintainers-template.toml /tmp/tbtc-maintainers-template.toml

COPY ./provision-tbtc-maintainers.js /tmp/provision-tbtc-maintainers.js

ENTRYPOINT ["node", "./provision-tbtc-maintainers.js"]
