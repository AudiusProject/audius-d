NETWORK ?= stage
TAG ?= latest

UI_DIR := web/ui
UI_ARTIFACTS_DIR := pkg/gui/dist
UI_ARTIFACTS := $(shell find $(UI_ARTIFACTS_DIR) -type f -name '*.js' -o -name '*.css') $(UI_ARTIFACTS_DIR)/index.html
UI_SRC := $(shell find $(UI_DIR) -type f -not -path '$(UI_DIR)/node_modules/*')

ABI_DIR := pkg/register/ABIs
SRC := $(shell find . -type f -name '*.go') go.mod go.sum $(UI_ARTIFACTS)


.PHONY: audiusctl
audiusctl: bin/audiusctl-arm bin/audiusctl-x86

bin/audiusctl-arm: $(SRC) $(UI_ARTIFACTS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audiusctl-arm ./cmd/audiusctl

bin/audiusctl-x86: $(SRC) $(UI_ARTIFACTS)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audiusctl-x86 ./cmd/audiusctl

$(UI_ARTIFACTS): $(UI_SRC)
	@echo "Building GUI..."
	cd $(UI_DIR) && npm i && npm run build

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

