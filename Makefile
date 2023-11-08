NETWORK ?= stage
TAG ?= latest

build-go:
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audius-d-x86 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.imageTag=$(TAG)" -o bin/audius-d-arm main.go

build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=$(NETWORK) -t audius/audius-docker-compose:$(TAG) .

build-push: build-docker push-docker

clean:
	rm -f bin/*
