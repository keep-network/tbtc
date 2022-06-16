FROM node:13-alpine AS install

ENV LD_LIBRARY_PATH=/usr/local/lib/

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
    make \
	python3 \
    git

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm ci

FROM node:13-alpine

WORKDIR /tmp

COPY --from=install /tmp .

COPY ./relay-config.toml /tmp/relay-config.toml

COPY ./provision-relay.js /tmp/provision-relay.js

RUN mkdir -p /mnt/relay

ENTRYPOINT ["node", "./provision-relay.js"]
