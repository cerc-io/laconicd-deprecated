FROM golang:alpine AS build-env

# Install minimum necessary dependencies,
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set up dependencies
ENV PACKAGES git build-base

# Set working directory for the build
WORKDIR /go/src/github.com/cerc-io/laconicd

# Add source files
COPY . .

# build binary
# RUN make build-linux
RUN COSMOS_BUILD_OPTIONS=badgerdb make build


# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates jq
WORKDIR /

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/cerc-io/laconicd/build/laconicd /usr/bin/laconicd

EXPOSE 26656 26657 1317 9090 8545 8546

# Run ethermintd by default
CMD ["laconicd","start","--gql-playground","--gql-server","--home","/laconic","--mode","validator","--db-backend","badgerdb"]