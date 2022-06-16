FROM golang:1.15.7-alpine3.13 AS gobuild

# Client Versioning.
ARG VERSION
ARG REVISION

# Environment variables.
ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=relay \
	APP_DIR=/go/src/github.com/keep-network/tbtc/relay \
	BIN_PATH=/usr/local/bin \
	LD_LIBRARY_PATH=/usr/local/lib/ \
	GO111MODULE=on

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
	make \
	nodejs \
	npm \
	python3 \
	git && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

# Install Solidity compiler.
COPY --from=ethereum/solc:0.5.17 /usr/bin/solc /usr/bin/solc

# Get gotestsum tool
RUN go get gotest.tools/gotestsum

# Configure working directory.
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

# Get dependencies.
COPY go.mod $APP_DIR/
COPY go.sum $APP_DIR/
RUN go mod download

# Install Solidity contracts.
COPY ./solidity $APP_DIR/solidity
RUN cd $APP_DIR/solidity && npm install

# Configure git to don't use unauthenticated protocol.
RUN git config --global url."https://".insteadOf git://

# Generate code.
COPY ./pkg/chain/ethereum/gen $APP_DIR/pkg/chain/ethereum/gen
RUN go generate ./...

# Copy app files.
COPY ./ $APP_DIR/

# Build the application.
RUN GOOS=linux go build \
	-ldflags "-X main.version=$VERSION -X main.revision=$REVISION" \
	-a -o $APP_NAME ./ && \
	mv $APP_NAME $BIN_PATH

# Configure runtime container.
FROM alpine:3.13

ENV APP_NAME=relay \
	BIN_PATH=/usr/local/bin

COPY --from=gobuild $BIN_PATH/$APP_NAME $BIN_PATH

# docker caches more when using CMD [] resulting in a faster build.
CMD []
