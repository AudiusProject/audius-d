# dot-slash-audius

because running audius should look like this
```
./audius
```

## build

build and copy the go binary to server
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audius main.go
scp audius ubuntu@stage-creator-10:~/audius
```

build and push the docker image
```
DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -t endliine/audius-docker-compose:linux .
docker push endliine/audius-docker-compose:linux
```

## run on instance

example running on stage cn 10
```
ssh stage-creator-10
```

stop the old
```
audius-cli down
```

run the new
```
./audius -c ~/path/to/override.env
```

down the new
```
./audius down
```
