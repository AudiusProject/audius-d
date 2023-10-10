.PHONY: build-go build-docker

all: build-go build-docker

# Build Go binary
build-go:
	@echo "Building Go binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audius-linux main.go
	CGO_ENABLED=0 go build -o audius main.go

# Build Docker image and push
build-docker: build-go
	@echo "Building Docker image..."
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build --build-arg NETWORK=stage -t endliine/audius-docker-compose:linux .
	rm audius-linux
	@echo "Pushing Docker image..."
	docker push endliine/audius-docker-compose:linux
