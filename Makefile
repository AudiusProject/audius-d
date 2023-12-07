NETWORK ?= stage
TAG ?= latest

UI_DIR := web/ui
UI_ARTIFACT_DIR := pkg/gui/dist
UI_ARTIFACT := $(UI_ARTIFACT_DIR)/index.html
UI_SRC := $(shell find $(UI_DIR) -type f -not -path '$(UI_DIR)/node_modules/*')

ABI_DIR := pkg/register/ABIs
SRC := $(shell find . -type f -name '*.go') go.mod go.sum $(UI_ARTIFACT)

VERSION_FILE := .version.json
VERSION_LDFLAG := -X main.Version=$(shell git rev-parse HEAD)
# Intentionally kept separate to allow dynamic versioning
#LDFLAGS := ""


audius-ctl: bin/audius-ctl-arm bin/audius-ctl-x86

bin/audius-ctl-arm: $(SRC) $(UI_ARTIFACT)
	@echo "Building arm audius-ctl..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-arm ./cmd/audius-ctl

bin/audius-ctl-x86: $(SRC) $(UI_ARTIFACT)
	@echo "Building x86 audius-ctl..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-x86 ./cmd/audius-ctl

bin/audius-ctl-arm-mac: $(SRC) $(UI_ARTIFACT)
	@echo "Building arm audius-ctl..."
	GOOS=darwin GOARCH=arm64 go build -tags mac -ldflags "$(VERSION_LDFLAG) $(LDFLAGS)" -o bin/audius-ctl-arm ./cmd/audius-ctl

$(UI_ARTIFACT): $(UI_SRC)
	@echo "Building GUI..."
	cd $(UI_DIR) && npm i && npm run build

.PHONY: release-audius-ctl audius-ctl-versioned
release-audius-ctl: clean audius-ctl-versioned
	bash scripts/github_release.sh

audius-ctl-versioned: VERSION_LDFLAG := -X main.Version=$(shell jq -r .version $(VERSION_FILE))
audius-ctl-versioned: audius-ctl

.PHONY: regen-abis
regen-abis:
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ERC20Detailed.json | jq '.abi' > $(ABI_DIR)/ERC20Detailed.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/Registry.json | jq '.abi' > $(ABI_DIR)/Registry.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ServiceProviderFactory.json | jq '.abi' > $(ABI_DIR)/ServiceProviderFactory.json

.PHONY: build-docker push-docker build-push
build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

build-push: build-docker push-docker

.PHONY: clean
clean:
	rm -f bin/*
	rm -rf $(UI_ARTIFACT_DIR)
