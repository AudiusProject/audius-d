NETWORK ?= stage
TAG ?= latest

build-go: build-gui
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audius-d-x86 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audius-d-arm main.go

build-gui:
	@echo "Building GUI..."
	cd ./gui/ui && npm run build

regen-abis:
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ERC20Detailed.json | jq '.abi' > ./register/ABIs/ERC20Detailed.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/Registry.json | jq '.abi' > ./register/ABIs/Registry.json
	curl -s https://raw.githubusercontent.com/AudiusProject/audius-protocol/main/packages/libs/src/eth-contracts/ABIs/ServiceProviderFactory.json | jq '.abi' > ./register/ABIs/ServiceProviderFactory.json

build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

build-push: build-docker push-docker

clean:
	rm -f bin/*
