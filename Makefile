.PHONY: audius

audius: main.go
	@echo "Building Go binary..."
	CGO_ENABLED=0 go build -o bin/audius main.go

build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=stage -t audius/dot-slash:dev .

build-docker-ci:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=stage --build-arg BRANCH=foundation -t audius/dot-slash:ci .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=stage -t audius/dot-slash:dev .

clean:
	rm -f bin/*
