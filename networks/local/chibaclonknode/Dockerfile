FROM golang:alpine AS build-env

# Install minimum necessary dependencies,
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set up dependencies
ENV PACKAGES git build-base

# Set working directory for the build
WORKDIR /go/src/github.com/vulcanize/chiba-clonk

# Add source files
COPY . .

# build binary
RUN make build-linux


# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates jq
WORKDIR /

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/vulcanize/chiba-clonk/build/chibaclonkd /usr/bin/chibaclonkd

EXPOSE 26656 26657 1317 9090 8545 8546

# Run ethermintd by default
CMD ["chibaclonkd","start","--gql-playground","--gql-server","--home","/chibaclonk"]