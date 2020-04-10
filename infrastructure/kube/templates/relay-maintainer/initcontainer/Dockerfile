FROM node:11 AS runtime

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm ci

COPY ./env-template /tmp/env-template

COPY ./provision-relay-maintainer.js /tmp/provision-relay-maintainer.js

ENTRYPOINT ["node", "./provision-relay-maintainer.js"]
