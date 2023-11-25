# audius-d

the simpliest way to run and interact with an audius node.

## install

```
curl -sSL https://raw.githubusercontent.com/AudiusProject/audius-d/main/install.sh | sh
```

## quickstart

run a dev setup

```
mkdir ~/.audius && cp sample.audius.conf ~/.audius/audius.conf

audius-d
```

## run

**configure**

minimal required config, (default location `~/.audius/audius.conf`) or pass via `-c` flag at runtime

```
# creator-node audius.conf
creatorNodeEndpoint=
delegateOwnerWallet=
delegatePrivateKey=
spOwnerWallet=
```

```
# discovery-provider audius.conf
audius_discprov_url=
audius_delegate_owner_wallet=
audius_delegate_private_key=
```

**run**
```
audius-d [-c audius.conf]
```

## build

builds required go binaries that are (for now) committed to this repo on the `bin` branch by CI.

```
make
```

 ## dev

 for green https for localhost. add caddy to your host trusted certificate authorities.

 ```
 docker exec creator-node sh -c "docker exec caddy cat /data/caddy/pki/authorities/local/root.crt" > caddy-root.crt

 sudo security add-trusted-cert -d -r trustRoot -k "/Library/Keychains/System.keychain" "$(pwd)/caddy-root.crt"

 echo $(docker exec creator-node sh -c "docker exec caddy cat /data/caddy/pki/authorities/local/root.crt") | sudo security add-trusted-cert -d -r trustRoot -k "/Library/Keychains/System.keychain" /dev/stdin
 ```
 