.PHONY: audius

audius: main.go
	@echo "Building Go binary..."
	CGO_ENABLED=0 go build -o bin/audius main.go

build-docker:
	@echo "Building Docker image..."
	docker buildx build --load --build-arg NETWORK=stage BRANCH=as/dot-slash-audius -t audius/dot-slash:dev .

push-docker:
	@echo "Pushing Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --push --build-arg NETWORK=stage -t audius/dot-slash:dev .

clean:
	rm -f bin/*
