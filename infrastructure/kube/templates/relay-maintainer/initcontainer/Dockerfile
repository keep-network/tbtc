FROM node:13-alpine AS install

RUN apk --no-cache add git

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm ci

FROM node:13-alpine

WORKDIR /tmp

COPY --from=install /tmp .

COPY ./env-template /tmp/env-template

COPY ./provision-relay-maintainer.js /tmp/provision-relay-maintainer.js

RUN mkdir -p /mnt/relay-maintainer

ENTRYPOINT ["node", "./provision-relay-maintainer.js"]
