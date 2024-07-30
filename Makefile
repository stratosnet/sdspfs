# Force Go Modules
GO111MODULE = on

GOCC ?= go
GOFLAGS ?=

# If set, override the install location for plugins
IPFS_PATH ?= $(HOME)/.ipfs
# Just to inform the user which kubo-version go.mod uses.
IPFS_VERSION = $(lastword $(shell $(GOCC) list -m github.com/ipfs/kubo))

# GOFLAGS += -asmflags=all="-trimpath=${GOPATH}" -gcflags=all="-trimpath=${GOPATH}"
GOFLAGS += -trimpath

.PHONY: install build

FORCE:

sdspfs.so: main/main.go go.mod
	CGO_ENABLED=1 $(GOCC) build $(GOFLAGS) -buildmode=plugin -o "$@" "$<"
	chmod +x "$@"

build: sdspfs.so
	@echo "Built against" $(IPFS_VERSION)

install: build
	mkdir -p "$(IPFS_PATH)/plugins/"
	cp -f sdspfs.so "$(IPFS_PATH)/plugins/sdspfs.so"