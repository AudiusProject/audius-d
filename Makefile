.PHONY: audius

audius: main.go
	@echo "Building Go binary..."
	CGO_ENABLED=0 go build -o bin/audius main.go

audius-linux: main.go
	@echo "Building Go binary for linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/audius-linux main.go

build-docker: audius-linux
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=stage -t audius/dot-slash:dev .

push-docker: audius-linux
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=stage -t audius/dot-slash:dev .

clean:
	rm -f bin/*
