FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 g++
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/Galactica-corp/galactica

COPY ../go.mod .
COPY ../go.sum .

# Install dependencies
RUN go mod download

# Add source files
COPY .. .

RUN make build-alpine

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/Galactica-corp/galactica/build/galacticad /usr/bin/galacticad
COPY --from=build-env /go/src/github.com/Galactica-corp/galactica/localnet/start-node.sh /root/start-node.sh

EXPOSE 26656 26657 1317 9090

CMD ["./start-node.sh"]
