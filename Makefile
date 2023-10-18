.PHONY: audius

build-go:
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/audius-x86 main.go
	CGO_ENABLED=0 GOARCH=arm64 go build -o bin/audius-arm main.go

build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=stage -t audius/audius-docker-compose:stage .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=stage -t audius/audius-docker-compose:stage .

build-push: build-docker push-docker

clean:
	rm -f bin/*
