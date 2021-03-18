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
	GO111MODULE=on

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
	make \
	git && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

# Get gotestsum tool
RUN go get gotest.tools/gotestsum

# Configure working directory.
RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR

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