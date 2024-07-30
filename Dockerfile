FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21 AS plugin

ARG TARGETOS TARGETARCH
ENV PLUGIN_DIR /sdspfs

RUN mkdir -p $PLUGIN_DIR

WORKDIR $PLUGIN_DIR

COPY go.mod go.sum $PLUGIN_DIR/
RUN go mod download
COPY . $PLUGIN_DIR

ARG MAKE_TARGET=build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make ${MAKE_TARGET}

FROM ipfs/kubo:v0.26.0

ENV PLUGIN_BINARY sdspfs.so
ENV PLUGIN_PATH /sdspfs/$PLUGIN_BINARY

RUN mkdir -p $IPFS_PATH/plugins

COPY --from=plugin --chown=ipfs:users $PLUGIN_PATH $IPFS_PATH/plugins/$PLUGIN_BINARY
