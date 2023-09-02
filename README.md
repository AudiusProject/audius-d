# dot-slash-audius

because running audius should look like this
```
./audius
```

## run

on an x86 linux machine that has docker installed
```
curl -o audius https://raw.githubusercontent.com/AudiusProject/dot-slash-audius/main/audius && chmod +x audius
```

minimal required config, can mount at this location or pass via `-c` flag at runtime
```
# ~/.audius/audius.conf
creatorNodeEndpoint=
delegateOwnerWallet=
delegatePrivateKey=
spOwnerWallet=
```

run
```
./audius [-c audius.conf]
```

## build

```
# if building on arm, it is much faster to build go binary first
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o audius main.go

# then build docker image
DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build --build-arg NETWORK=stage -t endliine/audius-docker-compose:linux .
docker push endliine/audius-docker-compose:linux
```

## todo

- docker buildx manifests for multi arch
- make work for discovery nodes also
