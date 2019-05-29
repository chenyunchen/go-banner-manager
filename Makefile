# banner manager version
APPLICATION_VERSION = latest
BINARY = banner-manager

## command
GO = go

.PHONY: all
all: build

.PHONY: build
build:
	$(MAKE) src.build-server
	$(MAKE) src.build-client

.PHONY: src.build-server
src.build-server:
	$(GO) build -v -o ${BINARY}-server ./src/server/...

.PHONY: src.build-client
src.build-client:
	$(GO) build -v -o ${BINARY}-client ./src/client/...

.PHONY: dockerfile.build
dockerfile.build:
	docker build --tag mercari/banner-manager:$(APPLICATION_VERSION) .
