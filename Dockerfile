FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22.4 AS plugin

ARG TARGETOS TARGETARCH
ENV PLUGIN_DIR /plugin

# installing binary to replace it in kubo to prevent modules mismtach for plugin
# currently local builds failed and plugin will not work
# so error "plugin was built with a different version of package" will be occured
RUN cd /tmp && \
  wget https://dist.ipfs.tech/kubo/v0.29.0/kubo_v0.29.0_linux-amd64.tar.gz && \
  tar -xvzf kubo_v0.29.0_linux-amd64.tar.gz && \
  chmod +x kubo/ipfs && \
  mv kubo/ipfs /usr/local/bin

# NOTE: for this init dockerfile go1.22.4 is used in ipfs binary. As golang version not strict
# potentially build will not work. So we need to check if version of go matches, otherwise update this Dockerfile
# first line of golang version
RUN echo '#!/bin/bash\n' \
  'ipfs_go_version=$(ipfs version --all | grep "Golang version" | awk '"'"'{print $NF}'"'"')\n' \
  'system_go_version=$(go version | awk '"'"'{print $3}'"'"')\n' \
  'if [ "$ipfs_go_version" == "$system_go_version" ]; then\n' \
  '    echo "Go versions match: $system_go_version"\n' \
  'else\n' \
  '    echo "Go versions do not match."\n' \
  '    echo "IPFS Go version: $ipfs_go_version"\n' \
  '    echo "System Go version: $system_go_version"\n' \
  '    exit 1\n' \
  'fi' > /usr/local/bin/check_versions.sh && \
  chmod +x /usr/local/bin/check_versions.sh

RUN /usr/local/bin/check_versions.sh

RUN mkdir -p $PLUGIN_DIR

WORKDIR $PLUGIN_DIR

COPY go.mod go.sum $PLUGIN_DIR/
RUN go mod download
COPY . $PLUGIN_DIR

ARG MAKE_TARGET=build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make ${MAKE_TARGET}

FROM debian:stable-slim AS builder

FROM ipfs/kubo:v0.29.0

ARG PLUGIN_BINARY
ENV PLUGIN_PATH /plugin/$PLUGIN_BINARY

RUN mkdir -p $IPFS_PATH/plugins

RUN rm $IPFS_PATH/plugins/$PLUGIN_BINARY || :

COPY --from=builder --chown=root:root /usr/lib/x86_64-linux-gnu/libdl.so.2 /lib/libdl.so.2

COPY --from=plugin --chown=ipfs:users $PLUGIN_PATH $IPFS_PATH/plugins/$PLUGIN_BINARY
COPY --from=plugin --chown=root:root /usr/local/bin/ipfs /usr/local/bin/ipfs
