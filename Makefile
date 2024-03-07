NETWORK ?= stage
AD_TAG ?= default
# One of patch, minor, or major
UPGRADE_TYPE ?= patch

ABI_DIR := pkg/register/ABIs
SRC := $(shell find . -type f -name '*.go') go.mod go.sum

VERSION_LDFLAG := -X main.Version=$(shell git rev-parse HEAD)
# Intentionally kept separate to allow dynamic versioning
#LDFLAGS := ""


audius-ctl: bin/audius-ctl-arm64 bin/audius-ctl-x86_64

bin/audius-ctl-arm64: $(SRC)
	@echo "Building arm audius-ctl..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-arm64 ./cmd/audius-ctl

bin/audius-ctl-x86_64: $(SRC)
	@echo "Building x86 audius-ctl..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-x86_64 ./cmd/audius-ctl

bin/audius-ctl-arm64-mac: $(SRC)
	@echo "Building arm audius-ctl..."
	GOOS=darwin GOARCH=arm64 go build -tags mac -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-arm64-mac ./cmd/audius-ctl

.PHONY: release-audius-ctl audius-ctl-production-build
release-audius-ctl:
	bash scripts/github_release.sh

audius-ctl-production-build: VERSION_LDFLAG := -X main.Version=$(shell bash scripts/get_new_version.sh $(UPGRADE_TYPE))
audius-ctl-production-build: clean audius-ctl

.PHONY: regen-abis
regen-abis:
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ERC20Detailed.json | jq '.abi' > $(ABI_DIR)/ERC20Detailed.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/Registry.json | jq '.abi' > $(ABI_DIR)/Registry.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ServiceProviderFactory.json | jq '.abi' > $(ABI_DIR)/ServiceProviderFactory.json

.PHONY: build-docker-local build-push-docker
build-docker-local:
	@echo "Building Docker image for local platform..."
	docker buildx build --load -t audius/audius-d:$(AD_TAG) .

build-push-docker:
	@echo "Building and pushing Docker images for all platforms..."
	docker buildx build --platform linux/amd64,linux/arm64 --push -t audius/audius-d:$(AD_TAG) .

.PHONY: install uninstall
install:
	bash scripts/install.sh

uninstall:
	bash scripts/uninstall.sh

.PHONY: clean
clean:
	rm -f bin/*
